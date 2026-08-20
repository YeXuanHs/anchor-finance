<template>
  <div class="third-login-page">
    <!-- 页面标题 -->
    <h2 class="page-title">第三方登录</h2>

    <!-- 说明文字 -->
    <div class="page-desc">
      <span>选择开启第三方登录方式</span>
      <el-link type="primary" href="https://www.idcsmart.com/wiki_list/653.html" target="_blank" :underline="false">帮助文档</el-link>
    </div>

    <div class="action-bar">
      <el-button type="primary" @click="handleGetMore">获取更多接口</el-button>
    </div>

    <!-- 第三方登录列表 -->
    <el-table 
      :data="loginMethods" 
      v-loading="loading" 
      border 
      stripe
      empty-text="暂无数据"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" />
      
      <el-table-column prop="display_name" label="名称" min-width="150">
        <template #default="{ row }">
          <span class="plugin-name">{{ row.display_name }}</span>
        </template>
      </el-table-column>
      
      <el-table-column prop="developer" label="开发者" min-width="120" />
      
      <el-table-column prop="current_version" label="当前版本" width="100">
        <template #default="{ row }">
          {{ row.current_version || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column prop="latest_version" label="最新版本" width="100">
        <template #default="{ row }">
          {{ row.latest_version || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch 
            v-if="row.installed"
            v-model="row.is_enabled" 
            @change="handleToggle(row)"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>
      
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <template v-if="row.installed">
            <el-button 
              size="small" 
              :type="row.is_enabled ? 'warning' : 'success'" 
              link
              @click="handleToggle(row)"
            >
              {{ row.is_enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" link @click="handleUninstall(row)">卸载</el-button>
            <el-button size="small" type="primary" link @click="handleConfigure(row)">配置</el-button>
          </template>
          <template v-else>
            <el-button size="small" type="primary" @click="handleInstall(row)">安装</el-button>
          </template>
        </template>
      </el-table-column>
    </el-table>

    <!-- 配置弹窗 -->
    <el-dialog 
      v-model="configDialogVisible" 
      :title="`${currentItem?.display_name} - 配置`" 
      width="600px"
      destroy-on-close
    >
      <el-form ref="configFormRef" :model="configFormData" label-position="top">
        <el-form-item label="App ID" prop="app_id">
          <el-input v-model="configFormData.app_id" placeholder="请输入App ID" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input v-model="configFormData.app_secret" type="password" placeholder="请输入App Secret" show-password />
        </el-form-item>
        <el-form-item label="回调地址" prop="callback_url">
          <el-input v-model="configFormData.callback_url" disabled />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingConfig" @click="handleSaveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const loginMethods = ref<any[]>([])
const configDialogVisible = ref(false)
const currentItem = ref<any>(null)
const configFormData = ref<any>({
  app_id: '',
  app_secret: '',
  callback_url: ''
})
const configFormRef = ref<FormInstance>()
const savingConfig = ref(false)

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/oauth-providers' })
    loginMethods.value = data || []
  } catch (error) {
    console.error('fetch oauth providers failed:', error)
  } finally {
    loading.value = false
  }
}

const handleGetMore = () => {
  window.open('https://www.idcsmart.com/wiki_list/653.html', '_blank')
}

const handleToggle = async (item: any) => {
  const original = item.is_enabled
  try {
    await request.post({ 
      url: `/api/admin/oauth-providers/${item.id}/toggle`, 
      data: { enabled: item.is_enabled } 
    })
    ElMessage.success(item.is_enabled ? '已启用' : '已禁用')
  } catch (error) {
    item.is_enabled = original
    ElMessage.error('操作失败')
  }
}

const handleUninstall = async (item: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要卸载 "${item.display_name}" 吗？`, 
      '确认卸载', 
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await request.post({ url: `/api/admin/oauth-providers/${item.id}/uninstall` })
    ElMessage.success('已卸载')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('卸载失败')
    }
  }
}

const handleInstall = async (item: any) => {
  try {
    await request.post({ url: `/api/admin/oauth-providers/${item.id}/install` })
    ElMessage.success('已安装')
    fetchList()
  } catch (error) {
    ElMessage.error('安装失败')
  }
}

const handleConfigure = async (item: any) => {
  currentItem.value = item
  try {
    const data = await request.get({ url: `/api/admin/oauth-providers/${item.id}/config` })
    configFormData.value = data || { app_id: '', app_secret: '', callback_url: '' }
  } catch (error) {
    configFormData.value = { app_id: '', app_secret: '', callback_url: '' }
  }
  configDialogVisible.value = true
}

const handleSaveConfig = async () => {
  if (!configFormRef.value) return
  try {
    await configFormRef.value.validate()
    savingConfig.value = true
    await request.post({ 
      url: `/api/admin/oauth-providers/${currentItem.value.id}/config`, 
      data: configFormData.value 
    })
    ElMessage.success('保存成功')
    configDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    savingConfig.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.third-login-page {
  padding: 20px;
}

.page-title {
  margin: 0 0 16px 0;
  font-size: 20px;
  font-weight: 600;
}

.page-desc {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  color: var(--el-text-color-secondary);
}

.action-bar {
  margin-bottom: 16px;
}

.plugin-name {
  font-weight: 500;
}
</style>
