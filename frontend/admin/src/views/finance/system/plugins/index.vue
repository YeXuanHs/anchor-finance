<template>
  <div class="plugins-page">
    <!-- 顶部操作栏 -->
    <div class="page-header">
      <div class="header-left">
        <h2>插件管理</h2>
        <span class="subtitle">管理系统功能插件，包括邮件、短信、实名认证、支付、OAuth、服务器模块等</span>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新增插件
      </el-button>
    </div>

    <!-- 分类筛选 -->
    <el-tabs v-model="activeType" @tab-change="fetchData">
      <el-tab-pane label="全部" name="" />
      <el-tab-pane label="邮件" name="mail" />
      <el-tab-pane label="短信" name="sms" />
      <el-tab-pane label="实名认证" name="certification" />
      <el-tab-pane label="支付网关" name="gateway" />
      <el-tab-pane label="OAuth登录" name="oauth" />
      <el-tab-pane label="服务器模块" name="server" />
      <el-tab-pane label="扩展插件" name="addon" />
    </el-tabs>

    <!-- 插件列表 -->
    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="插件名称" min-width="150">
          <template #default="{ row }">
            <div class="plugin-info">
              <span class="plugin-title">{{ row.title }}</span>
              <span class="plugin-version" v-if="row.version">v{{ row.version }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="标识" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="author" label="作者" width="100" />
        <el-table-column label="系统内置" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_system" type="info" size="small">内置</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              @change="handleToggleStatus(row)"
              :disabled="row.is_system && !row.is_enabled"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="info" link size="small" @click="handleConfig(row)">配置</el-button>
            <el-popconfirm 
              v-if="!row.is_system"
              title="确定删除该插件？" 
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑插件' : '新增插件'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="插件类型" prop="type">
          <el-select v-model="formData.type" placeholder="选择插件类型" :disabled="isEdit">
            <el-option 
              v-for="(label, value) in pluginTypes" 
              :key="value" 
              :label="label" 
              :value="value" 
            />
          </el-select>
        </el-form-item>

        <el-form-item label="标识名" prop="name">
          <el-input v-model="formData.name" placeholder="如 Smtp" :disabled="isEdit" />
          <div class="form-hint">系统内部标识，创建后不可修改</div>
        </el-form-item>

        <el-form-item label="显示名称" prop="title">
          <el-input v-model="formData.title" placeholder="如 SMTP邮件" />
        </el-form-item>

        <el-form-item label="作者">
          <el-input v-model="formData.author" placeholder="作者名称" />
        </el-form-item>

        <el-form-item label="版本号">
          <el-input v-model="formData.version" placeholder="1.0" />
        </el-form-item>

        <el-form-item label="帮助文档">
          <el-input v-model="formData.help_url" placeholder="https://..." />
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="插件描述" />
        </el-form-item>

        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" />
          <div class="form-hint">数值越小越靠前</div>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="formData.is_enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      :title="`${currentPlugin?.title || ''} - 配置`"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form label-width="120px">
        <!-- SMTP配置 -->
        <template v-if="currentPlugin?.name === 'Smtp'">
          <el-form-item label="SMTP服务器">
            <el-input v-model="configForm.host" placeholder="smtp.example.com" />
          </el-form-item>
          <el-form-item label="端口">
            <el-input-number v-model="configForm.port" :min="1" :max="65535" />
          </el-form-item>
          <el-form-item label="加密方式">
            <el-select v-model="configForm.encryption">
              <el-option label="SSL" value="ssl" />
              <el-option label="TLS" value="tls" />
              <el-option label="无" value="" />
            </el-select>
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="configForm.username" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="configForm.password" type="password" show-password />
          </el-form-item>
          <el-form-item label="发件人地址">
            <el-input v-model="configForm.from_address" />
          </el-form-item>
          <el-form-item label="发件人名称">
            <el-input v-model="configForm.from_name" />
          </el-form-item>
        </template>

        <!-- 赛邮短信配置 -->
        <template v-if="currentPlugin?.name === 'Submail'">
          <el-form-item label="AppID">
            <el-input v-model="configForm.appid" />
          </el-form-item>
          <el-form-item label="AppKey">
            <el-input v-model="configForm.appkey" show-password />
          </el-form-item>
          <el-form-item label="短信签名">
            <el-input v-model="configForm.app_sign" />
          </el-form-item>
        </template>

        <!-- 短信宝配置 -->
        <template v-if="currentPlugin?.name === 'Smsbao'">
          <el-form-item label="用户名">
            <el-input v-model="configForm.user" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="configForm.pass" type="password" show-password />
          </el-form-item>
          <el-form-item label="短信签名">
            <el-input v-model="configForm.sign" />
          </el-form-item>
        </template>

        <!-- 微信实名认证配置 -->
        <template v-if="currentPlugin?.name === 'Wechat'">
          <el-form-item label="SecretID">
            <el-input v-model="configForm.secret_id" />
          </el-form-item>
          <el-form-item label="SecretKey">
            <el-input v-model="configForm.secret_key" show-password />
          </el-form-item>
        </template>

        <!-- 阿里云实名认证配置 -->
        <template v-if="currentPlugin?.name === 'Idcsmartali'">
          <el-form-item label="AccessKeyID">
            <el-input v-model="configForm.access_key_id" />
          </el-form-item>
          <el-form-item label="AccessKeySecret">
            <el-input v-model="configForm.access_key_secret" show-password />
          </el-form-item>
        </template>

        <!-- 微信登录配置 -->
        <template v-if="currentPlugin?.name === 'Weixin' && currentPlugin?.type === 'oauth'">
          <el-form-item label="AppID">
            <el-input v-model="configForm.app_id" />
          </el-form-item>
          <el-form-item label="AppSecret">
            <el-input v-model="configForm.app_secret" show-password />
          </el-form-item>
          <el-form-item label="回调地址">
            <el-input v-model="configForm.callback_url" placeholder="https://example.com/oauth/weixin/callback" />
          </el-form-item>
        </template>

        <!-- QQ登录配置 -->
        <template v-if="currentPlugin?.name === 'QQ'">
          <el-form-item label="AppID">
            <el-input v-model="configForm.app_id" />
          </el-form-item>
          <el-form-item label="AppKey">
            <el-input v-model="configForm.app_key" show-password />
          </el-form-item>
          <el-form-item label="回调地址">
            <el-input v-model="configForm.callback_url" />
          </el-form-item>
        </template>

        <!-- ProxmoxVE配置 -->
        <template v-if="currentPlugin?.name === 'ProxmoxVE'">
          <el-form-item label="面板地址">
            <el-input v-model="configForm.panel" placeholder="https://pve.example.com:8006/" />
          </el-form-item>
          <el-form-item label="节点名称">
            <el-input v-model="configForm.node" placeholder="pve" />
          </el-form-item>
          <el-form-item label="API端口">
            <el-input-number v-model="configForm.port" :min="1" :max="65535" />
          </el-form-item>
        </template>

        <!-- 通用配置（其他插件） -->
        <template v-if="isGenericConfig">
          <el-form-item label="配置JSON">
            <el-input v-model="configFormRaw" type="textarea" :rows="12" placeholder='{"key": "value"}' />
            <div class="form-hint">请输入有效的JSON格式配置</div>
          </el-form-item>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveConfig" :loading="submitting">保存配置</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/http'

interface Plugin {
  id: number
  name: string
  title: string
  type: string
  description: string
  author: string
  version: string
  help_url: string
  config: string
  is_system: boolean
  is_enabled: boolean
  sort_order: number
}

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const configDialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const activeType = ref('')
const tableData = ref<Plugin[]>([])
const currentPlugin = ref<Plugin | null>(null)
const pluginTypes = ref<Record<string, string>>({})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const formData = reactive({
  name: '',
  title: '',
  type: 'mail',
  description: '',
  author: '',
  version: '1.0',
  help_url: '',
  config: '',
  is_system: false,
  is_enabled: false,
  sort_order: 0
})

const configForm = reactive<Record<string, any>>({})
const configFormRaw = ref('')

const formRules = {
  name: [{ required: true, message: '请输入标识名', trigger: 'blur' }],
  title: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择插件类型', trigger: 'change' }]
}

// 特定插件的配置表单
const specificConfigs = ['Smtp', 'Submail', 'Smsbao', 'Wechat', 'Idcsmartali', 'Weixin', 'QQ', 'ProxmoxVE']

// 是否使用通用配置
const isGenericConfig = computed(() => {
  return !specificConfigs.includes(currentPlugin.value?.name || '')
})

// 类型标签样式
const getTypeTag = (type: string) => {
  const map: Record<string, any> = {
    mail: '',
    sms: 'success',
    certification: 'warning',
    gateway: 'danger',
    oauth: 'info',
    server: '',
    addon: 'success'
  }
  return map[type] || ''
}

// 类型显示文本
const getTypeLabel = (type: string) => {
  return pluginTypes.value[type] || type
}

// 获取插件类型
const fetchPluginTypes = async () => {
  try {
    const { data } = await request.get({ url: '/api/admin/plugins/types' })
    pluginTypes.value = data || {}
  } catch (error) {
    console.error('获取插件类型失败:', error)
  }
}

// 获取列表数据
const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get({
      url: '/api/admin/plugins',
      params: {
        type: activeType.value,
        page: pagination.page,
        page_size: pagination.pageSize
      }
    })

    tableData.value = data?.list || data?.items || data || []
    pagination.total = data?.total || tableData.value.length
  } catch (error) {
    console.error('获取插件列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 打开新增对话框
const handleAdd = () => {
  isEdit.value = false
  editId.value = 0
  Object.assign(formData, {
    name: '',
    title: '',
    type: 'mail',
    description: '',
    author: '',
    version: '1.0',
    help_url: '',
    config: '',
    is_system: false,
    is_enabled: false,
    sort_order: 0
  })
  dialogVisible.value = true
}

// 打开编辑对话框
const handleEdit = (row: Plugin) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(formData, {
    name: row.name,
    title: row.title,
    type: row.type,
    description: row.description,
    author: row.author,
    version: row.version,
    help_url: row.help_url,
    config: row.config,
    is_system: row.is_system,
    is_enabled: row.is_enabled,
    sort_order: row.sort_order
  })
  dialogVisible.value = true
}

// 打开配置对话框
const handleConfig = (row: Plugin) => {
  currentPlugin.value = row
  configFormRaw.value = row.config || '{}'
  
  // 清空配置表单
  Object.keys(configForm).forEach(key => delete configForm[key])
  
  // 解析配置
  if (row.config) {
    try {
      const config = JSON.parse(row.config)
      Object.assign(configForm, config)
    } catch (e) {
      console.error('解析配置失败:', e)
    }
  }
  
  configDialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEdit.value) {
      await request.put({
        url: `/api/admin/plugins/${editId.value}`,
        data: formData
      })
      ElMessage.success('更新成功')
    } else {
      await request.post({
        url: '/api/admin/plugins',
        data: formData
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 保存配置
const handleSaveConfig = async () => {
  submitting.value = true
  try {
    let configStr = ''
    
    // 特定插件使用configForm对象
    if (!isGenericConfig.value) {
      configStr = JSON.stringify(configForm)
    } else {
      // 验证JSON格式
      try {
        JSON.parse(configFormRaw.value)
        configStr = configFormRaw.value
      } catch {
        ElMessage.error('请输入有效的JSON格式')
        return
      }
    }

    await request.put({
      url: `/api/admin/plugins/${currentPlugin.value?.id}/config`,
      data: { config: configStr }
    })
    ElMessage.success('配置已保存')
    configDialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '保存配置失败')
  } finally {
    submitting.value = false
  }
}

// 切换状态
const handleToggleStatus = async (row: Plugin) => {
  try {
    const endpoint = row.is_enabled 
      ? `/api/admin/plugins/${row.id}/enable`
      : `/api/admin/plugins/${row.id}/disable`
    
    await request.post({ url: endpoint })
    ElMessage.success('状态已更新')
  } catch (error) {
    row.is_enabled = !row.is_enabled
    ElMessage.error('更新状态失败')
  }
}

// 删除
const handleDelete = async (row: Plugin) => {
  try {
    await request.del({
      url: `/api/admin/plugins/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

onMounted(() => {
  fetchPluginTypes()
  fetchData()
})
</script>

<style scoped>
.plugins-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  color: #909399;
  font-size: 14px;
}

.table-card {
  background: #fff;
}

.plugin-info {
  display: flex;
  flex-direction: column;
}

.plugin-title {
  font-weight: 500;
}

.plugin-version {
  font-size: 12px;
  color: #909399;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.form-hint {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
</style>
