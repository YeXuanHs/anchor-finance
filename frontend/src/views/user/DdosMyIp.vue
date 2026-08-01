<template>
  <div class="ddos-myip-page">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(239, 68, 68, 0.1); color: #ef4444;">
          <el-icon :size="24"><View /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">总IP数量</div>
          <div class="stat-value">{{ stats.total || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(16, 185, 129, 0.1); color: #10b981;">
          <el-icon :size="24"><Shield /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">防护中</div>
          <div class="stat-value">{{ stats.protected || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(245, 158, 11, 0.1); color: #f59e0b;">
          <el-icon :size="24"><Warning /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">今日攻击</div>
          <div class="stat-value">{{ stats.attacks_today || 0 }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: rgba(59, 130, 246, 0.1); color: #3b82f6;">
          <el-icon :size="24"><DataLine /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">清洗流量</div>
          <div class="stat-value">{{ stats.cleaned_traffic || '0 Gbps' }}</div>
        </div>
      </div>
    </div>

    <!-- IP列表 -->
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>我的IP地址</span>
          <el-button type="primary" @click="showAddIpDialog">
            <el-icon><Plus /></el-icon>
            添加IP
          </el-button>
        </div>
      </template>

      <el-table :data="ipList" style="width: 100%" v-loading="loading">
        <el-table-column prop="ip" label="IP地址" width="160" />
        <el-table-column prop="domain" label="绑定域名" min-width="150" />
        <el-table-column prop="protection_level" label="防护等级" width="120">
          <template #default="{ row }">
            <el-tag :type="getProtectionType(row.protection_level)">
              {{ getProtectionText(row.protection_level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="bandwidth" label="防护带宽" width="120" />
        <el-table-column prop="status" label="防护状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '防护中' : '已关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="attack_count" label="累计攻击" width="100" />
        <el-table-column prop="created_at" label="添加时间" width="120" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewRules(row)">清洗规则</el-button>
            <el-button size="small" :type="row.status === 'active' ? 'warning' : 'success'" @click="toggleProtection(row)">
              {{ row.status === 'active' ? '关闭防护' : '开启防护' }}
            </el-button>
            <el-button size="small" type="danger" @click="deleteIp(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加IP对话框 -->
    <el-dialog v-model="addIpVisible" title="添加IP地址" width="500px">
      <el-form :model="addIpForm" :rules="addIpRules" ref="addIpFormRef" label-width="100px">
        <el-form-item label="IP地址" prop="ip">
          <el-input v-model="addIpForm.ip" placeholder="请输入IP地址" />
        </el-form-item>
        <el-form-item label="绑定域名" prop="domain">
          <el-input v-model="addIpForm.domain" placeholder="请输入绑定域名（可选）" />
        </el-form-item>
        <el-form-item label="防护等级" prop="protection_level">
          <el-select v-model="addIpForm.protection_level" placeholder="请选择防护等级">
            <el-option label="基础防护" value="basic" />
            <el-option label="高级防护" value="advanced" />
            <el-option label="旗舰防护" value="premium" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addIpVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAddIp" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 清洗规则对话框 -->
    <el-dialog v-model="rulesVisible" title="清洗规则配置" width="700px">
      <div class="rules-header">
        <span>当前IP：{{ currentIp?.ip }}</span>
        <el-button type="primary" size="small" @click="addRule">添加规则</el-button>
      </div>
      <el-table :data="cleanRules" style="width: 100%; margin-top: 16px;">
        <el-table-column prop="name" label="规则名称" />
        <el-table-column prop="type" label="规则类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ getRuleTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="动作" width="100">
          <template #default="{ row }">
            <el-tag :type="row.action === 'block' ? 'danger' : 'warning'">
              {{ row.action === 'block' ? '拦截' : '放行' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.status" :active-value="'active'" :inactive-value="'inactive'" @change="updateRuleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="editRule(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="deleteRule(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, View, Shield, Warning, DataLine } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitting = ref(false)
const ipList = ref([])
const cleanRules = ref([])
const addIpVisible = ref(false)
const rulesVisible = ref(false)
const currentIp = ref<any>(null)
const addIpFormRef = ref<FormInstance>()

const stats = ref({
  total: 0,
  protected: 0,
  attacks_today: 0,
  cleaned_traffic: '0 Gbps'
})

const addIpForm = ref({
  ip: '',
  domain: '',
  protection_level: 'basic'
})

const addIpRules = {
  ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
  protection_level: [{ required: true, message: '请选择防护等级', trigger: 'change' }]
}

const getProtectionType = (level: string) => {
  const map: Record<string, string> = {
    basic: 'info',
    advanced: 'warning',
    premium: 'danger'
  }
  return map[level] || 'info'
}

const getProtectionText = (level: string) => {
  const map: Record<string, string> = {
    basic: '基础防护',
    advanced: '高级防护',
    premium: '旗舰防护'
  }
  return map[level] || level
}

const getRuleTypeText = (type: string) => {
  const map: Record<string, string> = {
    cc: 'CC防护',
    syn: 'SYN防护',
    udp: 'UDP防护',
    http: 'HTTP防护',
    ip: 'IP黑名单'
  }
  return map[type] || type
}

const showAddIpDialog = () => {
  addIpForm.value = { ip: '', domain: '', protection_level: 'basic' }
  addIpVisible.value = true
}

const submitAddIp = async () => {
  if (!addIpFormRef.value) return
  await addIpFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await request.post('/api/v1/user/ddos/ips', addIpForm.value)
      ElMessage.success('IP添加成功')
      addIpVisible.value = false
      loadIpList()
    } catch (e: any) {
      ElMessage.error(e?.message || '添加失败')
    } finally {
      submitting.value = false
    }
  })
}

const viewRules = async (ip: any) => {
  currentIp.value = ip
  try {
    const { data } = await request.get(`/api/v1/user/ddos/ips/${ip.id}/rules`)
    cleanRules.value = data?.data || []
  } catch {
    cleanRules.value = []
  }
  rulesVisible.value = true
}

const toggleProtection = async (ip: any) => {
  const action = ip.status === 'active' ? '关闭' : '开启'
  try {
    await ElMessageBox.confirm(`确定要${action}该IP的防护吗？`, '确认操作', { type: 'warning' })
    await request.post(`/api/v1/user/ddos/ips/${ip.id}/toggle`)
    ElMessage.success(`防护已${action}`)
    loadIpList()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '操作失败')
  }
}

const deleteIp = async (ip: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该IP吗？删除后防护将立即停止。', '确认删除', { type: 'error' })
    await request.delete(`/api/v1/user/ddos/ips/${ip.id}`)
    ElMessage.success('IP已删除')
    loadIpList()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '删除失败')
  }
}

const addRule = () => {
  // TODO: 添加规则逻辑
  ElMessage.info('添加规则功能开发中')
}

const editRule = (rule: any) => {
  // TODO: 编辑规则逻辑
  ElMessage.info('编辑规则功能开发中')
}

const deleteRule = async (rule: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该规则吗？', '确认删除', { type: 'warning' })
    await request.delete(`/api/v1/user/ddos/rules/${rule.id}`)
    ElMessage.success('规则已删除')
    if (currentIp.value) viewRules(currentIp.value)
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '删除失败')
  }
}

const updateRuleStatus = async (rule: any) => {
  try {
    await request.put(`/api/v1/user/ddos/rules/${rule.id}`, { status: rule.status })
    ElMessage.success('规则状态已更新')
  } catch (e: any) {
    ElMessage.error(e?.message || '更新失败')
    if (currentIp.value) viewRules(currentIp.value)
  }
}

const loadIpList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/user/ddos/ips')
    ipList.value = data?.data?.list || data?.data?.items || data?.data || []
    // 更新统计
    stats.value = {
      total: ipList.value.length,
      protected: ipList.value.filter((ip: any) => ip.status === 'active').length,
      attacks_today: data?.data?.stats?.attacks_today || 0,
      cleaned_traffic: data?.data?.stats?.cleaned_traffic || '0 Gbps'
    }
  } catch {
    ipList.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadIpList()
})
</script>

<style scoped lang="scss">
.ddos-myip-page {
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    margin-bottom: 16px;

    @media (max-width: 992px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  .stat-card {
    background: #fff;
    border-radius: 12px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .stat-info {
      .stat-label {
        font-size: 13px;
        color: #86868b;
        margin-bottom: 4px;
      }

      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: #1d1d1f;
      }
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .rules-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
  }
}
</style>
