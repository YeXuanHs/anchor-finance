package service

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/util"
)

// DcimCloudProvision 开通魔方云虚拟机
// 完整搬自zjmf DcimCloud.php:2570-2845 createAccount
// 职责：从DB查全量配置 → 在魔方云创建用户+云主机 → 回写service表
// 返回值：包含cloud_id/domain/ip/password等开通结果的map
func (c *DcimCloudClient) DcimCloudProvision(serviceID uint) (map[string]interface{}, error) {
	db := database.GetDB()

	// ============================================================
	// 第1步：查询service关联的全部数据
	// ============================================================
	var svc model.Service
	if err := db.First(&svc, serviceID).Error; err != nil {
		return nil, fmt.Errorf("服务不存在: %v", err)
	}

	var user model.User
	if err := db.First(&user, svc.UserID).Error; err != nil {
		return nil, fmt.Errorf("用户不存在: %v", err)
	}

	var product model.Product
	if err := db.First(&product, svc.ProductID).Error; err != nil {
		return nil, fmt.Errorf("产品不存在: %v", err)
	}

	var server model.Server
	if err := db.First(&server, svc.ServerID).Error; err != nil {
		return nil, fmt.Errorf("服务器不存在: %v", err)
	}
	if server.ServerType != "dcimcloud" {
		return nil, fmt.Errorf("服务器类型不是dcimcloud")
	}

	// ============================================================
	// 第2步：检查是否已开通（zjmf: if ($service['dcimid'] > 0)）
	// ============================================================
	if svc.DcimID > 0 {
		return map[string]interface{}{
			"status":    "error",
			"msg":       fmt.Sprintf("服务已开通，魔方云ID: %d", svc.DcimID),
			"cloud_id":  svc.DcimID,
			"domain":    svc.Domain,
			"username":  svc.Username,
			"service_id": svc.ID,
		}, nil
	}

	// ============================================================
	// 第3步：解析configoptions（zjmf: $service['configoptions']）
	// 从service.Config JSON字段 + product.ConfigOptions默认值
	// ============================================================
	configOptions := c.parseProvisionConfigOptions(svc.Config, product.ConfigOptions)

	// ============================================================
	// 第4步：在魔方云创建用户（zjmf:2593-2603）
	// 用email作为username，因为email唯一且稳定
	// ============================================================
	username := user.Email

	// 先检查用户是否已存在
	checkResp, err := c.CheckUser(username)
	if err != nil {
		return nil, fmt.Errorf("检查魔方云用户失败: %v", err)
	}

	var cloudUserID interface{}
	if checkResp["status"] == "success" {
		// 用户已存在，直接取ID
		if origin, ok := checkResp["origin"]; ok {
			switch v := origin.(type) {
			case map[string]interface{}:
				cloudUserID = v["id"]
			case float64:
				cloudUserID = uint(v)
			}
		}
	}

	// 用户不存在，创建新用户
	if cloudUserID == nil {
		createUserData := map[string]interface{}{
			"username": username,
			"email":    user.Email,
			"name":     user.Username,
			"password": user.PasswordHash, // 用用户密码哈希作为魔方云密码（实际zjmf也用原密码）
		}
		createResp, err := c.CreateUser(createUserData)
		if err != nil {
			return nil, fmt.Errorf("创建魔方云用户失败: %v", err)
		}
		if createResp["status"] != "success" {
			return nil, fmt.Errorf("创建魔方云用户失败: %v", createResp["msg"])
		}
		if origin, ok := createResp["origin"]; ok {
			switch v := origin.(type) {
			case map[string]interface{}:
				cloudUserID = v["id"]
			case float64:
				cloudUserID = uint(v)
			}
		}
	}

	if cloudUserID == nil {
		// 创建后再次查询以获取用户ID
		recheck, err := c.CheckUser(username)
		if err != nil {
			return nil, fmt.Errorf("查询魔方云用户ID失败: %v", err)
		}
		if origin, ok := recheck["origin"]; ok {
			switch v := origin.(type) {
			case map[string]interface{}:
				cloudUserID = v["id"]
			case float64:
				cloudUserID = uint(v)
			}
		}
	}

	if cloudUserID == nil {
		return nil, fmt.Errorf("无法获取魔方云用户ID")
	}

	// ============================================================
	// 第5步：生成随机root密码
	// ============================================================
	randomPass := c.generateRandomPassword(16)

	// ============================================================
	// 第6步：构建POST参数（zjmf:2631-2786）
	// 从configoptions逐个映射，全部参数对齐zjmf源码
	// ============================================================
	postData := map[string]interface{}{}

	// ---- 基础参数 ----
	postData["area"] = configOptions["area"]
	postData["node"] = configOptions["node"]
	postData["os"] = configOptions["os"]
	postData["cpu"] = configOptions["cpu"]

	// ---- 内存单位转换（zjmf:2661-2667）----
	// zjmf: if (strpos($params['memory'], 'GB') !== false) { $params['memory'] = str_replace('GB', '', $params['memory']) * 1024; }
	memory := configOptions["memory"]
	if memStr, ok := memory.(string); ok {
		memStr = strings.TrimSpace(memStr)
		if strings.HasSuffix(strings.ToUpper(memStr), "GB") {
			memStr = strings.TrimSuffix(strings.ToUpper(memStr), "GB")
			memStr = strings.TrimSpace(memStr)
			if memVal, err := strconv.ParseFloat(memStr, 64); err == nil {
				postData["memory"] = int(memVal * 1024)
			} else {
				postData["memory"] = memory
			}
		} else if strings.HasSuffix(strings.ToUpper(memStr), "MB") {
			memStr = strings.TrimSuffix(strings.ToUpper(memStr), "MB")
			memStr = strings.TrimSpace(memStr)
			if memVal, err := strconv.Atoi(memStr); err == nil {
				postData["memory"] = memVal
			} else {
				postData["memory"] = memory
			}
		} else {
			// 纯数字，假设已经是MB
			if memVal, err := strconv.Atoi(memStr); err == nil {
				postData["memory"] = memVal
			} else {
				postData["memory"] = memory
			}
		}
	} else if memNum, ok := memory.(float64); ok {
		// JSON数字，zjmf默认当MB传
		postData["memory"] = int(memNum)
	} else if memInt, ok := memory.(int); ok {
		postData["memory"] = memInt
	} else {
		postData["memory"] = memory
	}

	// ---- 带宽 ----
	postData["bw"] = configOptions["bw"]
	postData["in_bw"] = configOptions["in_bw"]
	postData["out_bw"] = configOptions["bw"] // zjmf: $params['out_bw'] = $params['bw']

	// ---- IP数量 ----
	postData["ip_num"] = configOptions["ip_num"]

	// ---- 流量限制 ----
	postData["flow_limit"] = configOptions["flow_limit"]

	// ---- 快照/备份数量 ----
	postData["snap_num"] = configOptions["snap_num"]
	postData["backup_num"] = configOptions["backup_num"]

	// ---- 流量计费方向（zjmf:2678-2690）----
	// 1=出方向, 2=入方向, 3=双向
	postData["traffic_type"] = 3 // 默认双向计费
	if ttRaw, ok := configOptions["traffic_type"]; ok {
		switch v := ttRaw.(type) {
		case string:
			tt := strings.ToLower(strings.TrimSpace(v))
			switch tt {
			case "out", "outbound", "出方向":
				postData["traffic_type"] = 1
			case "in", "inbound", "入方向":
				postData["traffic_type"] = 2
			case "all", "both", "双向":
				postData["traffic_type"] = 3
			default:
				if n, err := strconv.Atoi(v); err == nil && n >= 1 && n <= 3 {
					postData["traffic_type"] = n
				}
			}
		case float64:
			if int(v) >= 1 && int(v) <= 3 {
				postData["traffic_type"] = int(v)
			}
		case int:
			if v >= 1 && v <= 3 {
				postData["traffic_type"] = v
			}
		}
	}

	// ---- 磁盘配置（zjmf:2705-2734）----
	// system_disk_size 和 data_disk_size 可能是字符串 "50" 或配置选项key
	systemDiskSize := c.parseDiskSize(configOptions["system_disk_size"], configOptions)
	dataDiskSize := c.parseDiskSize(configOptions["data_disk_size"], configOptions)

	if systemDiskSize > 0 {
		postData["system_disk_size"] = systemDiskSize
	}
	if dataDiskSize > 0 {
		postData["data_disk_size"] = dataDiskSize
	}

	// ---- 磁盘驱动 ----
	if driver, ok := configOptions["disk_driver"]; ok && driver != nil && driver != "" {
		postData["disk_driver"] = driver
	}

	// ---- hostname（zjmf:2754）----
	hostname := svc.Domain
	if hostname == "" {
		hostname = fmt.Sprintf("vm-%d", serviceID)
	}
	postData["hostname"] = hostname

	// ---- rootpass（zjmf:2756）----
	postData["rootpass"] = randomPass

	// ---- client用户ID（zjmf:2758）----
	postData["client"] = cloudUserID

	// ---- NAT限制（zjmf:2762-2768）----
	if val, ok := configOptions["nat_acl_limit"]; ok && val != nil {
		postData["nat_acl_limit"] = val
	}
	if val, ok := configOptions["nat_web_limit"]; ok && val != nil {
		postData["nat_web_limit"] = val
	}

	// ---- 网络类型（zjmf:2770-2783）----
	// network_type: nat / vpc / public
	networkType := "public"
	if nt, ok := configOptions["network_type"]; ok && nt != nil && nt != "" {
		if ntStr, ok := nt.(string); ok {
			networkType = strings.ToLower(strings.TrimSpace(ntStr))
		}
	}
	postData["network_type"] = networkType

	// ============================================================
	// 第7步：VPC网络选择（zjmf:2643-2658）
	// 如果是vpc网络类型，查找该区域负载最低的VPC
	// ============================================================
	var vpcNetworkID interface{}
	if networkType == "vpc" {
		vpcResp, err := c.GetVpcNetworks()
		if err == nil && vpcResp["status"] == "success" {
			vpcList := extractList(vpcResp["origin"])
			if len(vpcList) > 0 {
				areaName := fmt.Sprintf("%v", configOptions["area"])
				// 先筛选当前区域的VPC，找使用量最小的
				var bestVpc map[string]interface{}
				bestCount := -1

				for _, vpc := range vpcList {
					vpcMap, ok := vpc.(map[string]interface{})
					if !ok {
						continue
					}
					// 检查区域匹配
					if vpcArea, ok := vpcMap["area"].(string); ok {
						if vpcArea != areaName {
							continue
						}
					}
					// 检查是否可用
					if vpcStatus, ok := vpcMap["status"].(string); ok {
						if vpcStatus != "available" && vpcStatus != "active" {
							continue
						}
					}
					// 获取当前使用数量（找负载最低的）
					currentCount := 0
					if cc, ok := vpcMap["current_count"]; ok {
						switch v := cc.(type) {
						case float64:
							currentCount = int(v)
						case int:
							currentCount = v
						}
					}
					if bestCount < 0 || currentCount < bestCount {
						bestCount = currentCount
						bestVpc = vpcMap
					}
				}

				// 如果没有匹配区域的，取第一个可用的
				if bestVpc == nil && len(vpcList) > 0 {
					for _, vpc := range vpcList {
						vpcMap, ok := vpc.(map[string]interface{})
						if !ok {
							continue
						}
						if vpcStatus, ok := vpcMap["status"].(string); ok {
							if vpcStatus == "available" || vpcStatus == "active" {
								bestVpc = vpcMap
								break
							}
						}
					}
				}

				if bestVpc != nil {
					vpcNetworkID = bestVpc["id"]
				}
			}
		}
	}

	if vpcNetworkID != nil {
		postData["vpc_network_id"] = vpcNetworkID
	}

	// ---- 网络子类型（可选）----
	if val, ok := configOptions["network_sub_type"]; ok && val != nil && val != "" {
		postData["network_sub_type"] = val
	}

	// ---- SSH密钥（可选）----
	if val, ok := configOptions["ssh_key_id"]; ok && val != nil && val != "" && val != "0" {
		postData["ssh_key_id"] = val
	}

	// ---- 其他可能的参数透传 ----
	if val, ok := configOptions["ip_group"]; ok && val != nil && val != "" {
		postData["ip_group"] = val
	}
	if val, ok := configOptions["store"]; ok && val != nil && val != "" {
		postData["store"] = val
	}
	if val, ok := configOptions["ipv6_num"]; ok && val != nil {
		postData["ipv6_num"] = val
	}

	// ============================================================
	// 第8步：调用POST /clouds创建云主机（zjmf:2787-2796）
	// ============================================================
	result, err := c.CreateAccount(postData)
	if err != nil {
		return nil, fmt.Errorf("调用魔方云创建云主机失败: %v", err)
	}

	// 检查创建是否成功
	if result["status"] != "success" {
		errMsg := "创建云主机失败"
		if msg, ok := result["msg"]; ok && msg != nil {
			errMsg = fmt.Sprintf("%v", msg)
		}
		return nil, fmt.Errorf("魔方云创建云主机失败: %s", errMsg)
	}

	// ============================================================
	// 第9步：解析创建结果（zjmf:2798-2820）
	// ============================================================
	origin := result["origin"]
	cloudID := extractUint(origin, "id")
	if cloudID == 0 {
		// 某些魔方云版本直接返回ID
		if idVal, ok := result["id"]; ok {
			cloudID = extractUintDirect(idVal)
		}
	}
	if cloudID == 0 {
		return nil, fmt.Errorf("创建成功但无法获取云主机ID")
	}

	domain := extractString(origin, "domain")
	if domain == "" {
		domain = extractString(origin, "hostname")
	}
	primaryIP := extractString(origin, "ip")
	if primaryIP == "" {
		primaryIP = extractString(origin, "primary_ip")
	}

	// ============================================================
	// 第10步：加密密码 + 更新service表（zjmf:2802-2818）
	// ============================================================
	// bcrypt哈希（用于本地密码验证）
	hashedPass, err := HashPassword(randomPass)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %v", err)
	}
	// AES加密（用于存储可逆密码，面板密码接口需要明文）
	encryptedPass, err := util.EncryptAES(randomPass)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	updateData := map[string]interface{}{
		"dcim_id":       cloudID,
		"domain":        domain,
		"username":      username,
		"password_hash": hashedPass,
		"status":        "active",
	}

	if err := db.Model(&model.Service{}).Where("id = ?", serviceID).Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新服务记录失败: %v", err)
	}

	// ============================================================
	// 第11步：保存面板密码（zjmf:2822-2830）
	// 调用魔方云 /clouds/{id}/panel/password 保存明文密码
	// ============================================================
	_, panelErr := c.SavePanelPass(cloudID, encryptedPass)
	if panelErr != nil {
		// 面板密码保存失败不影响开通，记录日志即可
		// 服务已经开通成功，只是面板密码保存失败
		fmt.Printf("[DcimCloudProvision] 保存面板密码失败 serviceID=%d cloudID=%d: %v\n", serviceID, cloudID, panelErr)
	}

	// ============================================================
	// 第12步：构建返回结果（zjmf:2832-2844）
	// ============================================================
	// 尝试获取云主机完整信息（获取IP等）
	var ipAddresses []string
	if primaryIP != "" {
		ipAddresses = append(ipAddresses, primaryIP)
	}

	// 如果没有从创建结果拿到IP，尝试从云主机详情获取
	if len(ipAddresses) == 0 {
		syncResp, syncErr := c.Sync(cloudID)
		if syncErr == nil && syncResp["status"] == "success" {
			if syncOrigin, ok := syncResp["origin"].(map[string]interface{}); ok {
				if ips, ok := syncOrigin["ip_list"].([]interface{}); ok {
					for _, ip := range ips {
						if ipStr, ok := ip.(string); ok && ipStr != "" {
							ipAddresses = append(ipAddresses, ipStr)
						}
					}
				}
				if len(ipAddresses) == 0 {
					if ipStr := extractString(syncOrigin, "ip"); ipStr != "" {
						ipAddresses = append(ipAddresses, ipStr)
					}
				}
				if len(ipAddresses) == 0 {
					if ipStr := extractString(syncOrigin, "primary_ip"); ipStr != "" {
						ipAddresses = append(ipAddresses, ipStr)
					}
				}
				// 更新domain（如果创建时没拿到）
				if domain == "" {
					if d := extractString(syncOrigin, "domain"); d != "" {
						domain = d
						db.Model(&model.Service{}).Where("id = ?", serviceID).Update("domain", d)
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"status":    "success",
		"msg":       "云主机开通成功",
		"cloud_id":  cloudID,
		"domain":    domain,
		"username":  username,
		"password":  randomPass,
		"ip":        primaryIP,
		"ips":       ipAddresses,
		"service_id": serviceID,
	}, nil
}

// ============================================================
// 内部辅助方法
// ============================================================

// parseProvisionConfigOptions 解析configoptions
// zjmf: $service['configoptions'] 可能来自service.Config和product.ConfigOptions
// 优先使用service.Config（用户下单时的选择），fallback到product.ConfigOptions（产品默认值）
func (c *DcimCloudClient) parseProvisionConfigOptions(serviceConfigJSON string, productConfigJSON string) map[string]interface{} {
	result := make(map[string]interface{})

	// 先解析product的默认configoptions
	if productConfigJSON != "" {
		var productOpts map[string]interface{}
		if err := json.Unmarshal([]byte(productConfigJSON), &productOpts); err == nil {
			for k, v := range productOpts {
				result[k] = v
			}
		}
	}

	// 再解析service的config（用户下单时的选择，覆盖产品默认值）
	if serviceConfigJSON != "" {
		var serviceData map[string]interface{}
		if err := json.Unmarshal([]byte(serviceConfigJSON), &serviceData); err == nil {
			// zjmf的service.config可能包含configoptions子字段
			if co, ok := serviceData["configoptions"]; ok {
				if coMap, ok := co.(map[string]interface{}); ok {
					for k, v := range coMap {
						result[k] = v
					}
				}
			}
			// 也可能直接就是configoptions的key-value（扁平结构）
			for k, v := range serviceData {
				if k != "configoptions" {
					// 如果值不是nil/空，覆盖
					if v != nil && v != "" {
						result[k] = v
					}
				}
			}
		}
	}

	return result
}

// parseDiskSize 解析磁盘大小配置
// zjmf:2705-2734 支持两种格式：
//   - 纯数字字符串 "50" → 直接返回50
//   - 配置选项名 "system_disk_size" → 从configOptions中查找对应值
//   - 带单位 "50GB" → 转换后返回
func (c *DcimCloudClient) parseDiskSize(diskVal interface{}, configOptions map[string]interface{}) int {
	if diskVal == nil {
		return 0
	}

	switch v := diskVal.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return 0
		}
		// 尝试直接解析数字
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		// 尝试解析带单位的值
		upper := strings.ToUpper(v)
		if strings.HasSuffix(upper, "GB") {
			numStr := strings.TrimSuffix(upper, "GB")
			numStr = strings.TrimSpace(numStr)
			if n, err := strconv.Atoi(numStr); err == nil {
				return n
			}
		}
		if strings.HasSuffix(upper, "MB") {
			numStr := strings.TrimSuffix(upper, "MB")
			numStr = strings.TrimSpace(numStr)
			if n, err := strconv.Atoi(numStr); err == nil {
				return n / 1024 // MB转GB
			}
		}
		// 可能是configOptions中的key名，递归查找
		if subVal, ok := configOptions[v]; ok && subVal != diskVal {
			return c.parseDiskSize(subVal, configOptions)
		}
		// 可能是数字字符串带小数点
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int(f)
		}
	}
	return 0
}

// generateRandomPassword 生成指定长度的随机密码
// 包含大小写字母和数字，确保至少各含一个
func (c *DcimCloudClient) generateRandomPassword(length int) string {
	if length < 4 {
		length = 4
	}

	const (
		lowerLetters = "abcdefghijklmnopqrstuvwxyz"
		upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits       = "0123456789"
		specialChars = "!@#$%^&*"
		allChars     = lowerLetters + upperLetters + digits + specialChars
	)

	// 使用crypto/rand确保密码学安全随机
	result := make([]byte, length)

	// 确保至少包含每种字符各一个
	result[0] = lowerLetters[cryptoRandInt(len(lowerLetters))]
	result[1] = upperLetters[cryptoRandInt(len(upperLetters))]
	result[2] = digits[cryptoRandInt(len(digits))]
	result[3] = specialChars[cryptoRandInt(len(specialChars))]

	// 填充剩余位置
	for i := 4; i < length; i++ {
		result[i] = allChars[cryptoRandInt(len(allChars))]
	}

	// Fisher-Yates洗牌打乱顺序
	for i := length - 1; i > 0; i-- {
		j := cryptoRandInt(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// cryptoRandInt 使用crypto/rand生成 [0, max) 的随机整数
func cryptoRandInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// 极端fallback，实际不会发生
		return 0
	}
	return int(n.Int64())
}

// ============================================================
// 数据提取辅助函数
// ============================================================

// extractList 从origin提取列表（可能是[]interface{}或map带list/data子字段）
func extractList(origin interface{}) []interface{} {
	if origin == nil {
		return nil
	}
	switch v := origin.(type) {
	case []interface{}:
		return v
	case map[string]interface{}:
		if list, ok := v["list"].([]interface{}); ok {
			return list
		}
		if data, ok := v["data"].([]interface{}); ok {
			return data
		}
		if list, ok := v["items"].([]interface{}); ok {
			return list
		}
	}
	return nil
}

// extractUint 从嵌套map中安全提取uint值
func extractUint(data interface{}, key string) uint {
	if data == nil {
		return 0
	}
	if m, ok := data.(map[string]interface{}); ok {
		if val, ok := m[key]; ok {
			return extractUintDirect(val)
		}
	}
	return 0
}

// extractUintDirect 将任意类型转换为uint
func extractUintDirect(val interface{}) uint {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return uint(v)
	case int:
		return uint(v)
	case uint:
		return v
	case int64:
		return uint(v)
	case uint64:
		return uint(v)
	case string:
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint(n)
		}
	}
	return 0
}

// extractString 从嵌套map中安全提取string值
func extractString(data interface{}, key string) string {
	if data == nil {
		return ""
	}
	if m, ok := data.(map[string]interface{}); ok {
		if val, ok := m[key]; ok && val != nil {
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}
