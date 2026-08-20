<template>
  <div class="plugins-page">
    <!-- 页面标题 -->
    <h2 class="page-title">插件列表</h2>

    <!-- 说明文字 -->
    <el-alert type="info" :closable="false" class="page-desc">
      <template #title>
        <span>管理独立的功能扩展插件，可以启用、禁用、配置或卸载插件。</span>
      </template>
    </el-alert>

    <!-- 插件卡片网格 -->
    <div class="plugin-grid" v-loading="loading">
      <div 
        v-for="plugin in plugins" 
        :key="plugin.id" 
        class="plugin-card"
        :class="{ 'plugin-disabled': !plugin.is_enabled }"
      >
        <div class="plugin-header">
          <div class="plugin-icon">
            <el-icon :size="32" color="var(--el-color-primary)">
              <component :is="plugin.icon || 'Grid'" />
            </el-icon>
          </div>
          <div class="plugin-info">
            <h3 class="plugin-name">{{ plugin.title }}</h3>
            <p class="plugin-author">{{ plugin.author || '未知作者' }}</p>
          </div>
          <div class="plugin-status">
            <el-tag :type="plugin.is_enabled ? 'success' : 'info'" size="small">
              {{ plugin.is_enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </div>
        </div>
        
        <div class="plugin-body">
          <p class="plugin-desc">{{ plugin.description || '暂无描述' }}</p>
          <div class="plugin-meta">
            <span>版本: {{ plugin.version || '1.0' }}</span>
          </div>
        </div>
        
        <div class="plugin-footer">
          <el-switch 
            v-model="plugin.is_enabled" 
            @change="handleToggle(plugin)"
          />
          <div class="plugin-actions">
            <el-button size="small" type="primary" link @click="handleConfigure(plugin)">配置</el-button>
            <el-button size="small" type="danger" link @click="handleUninstall(plugin)">卸载</el-button>
          </div>
        </div>
      </div>
      
      <el-empty v-if="!loading && plugins.length === 0" description="暂无插件" />
    </div>

    <!-- 配置弹窗 -->
    <el-dialog 
      v-model="configDialogVisible" 
      :title="`${currentPlugin?.title} - 配置`" 
      width="600px"
      destroy-on-close
    >
      <el-form :model="configForm" label-position="top">
        <el-form-item label="启用状态">
          <el-switch v-model="configForm.is_enabled" />
        </el-form-item>
        <template v-if="currentPlugin?.config_fields">
          <el-form-item 
            v-for="field in currentPlugin.config_fields" 
            :key="field.key" 
            :label="field.label"
          >
            <el-input 
              v-if="field.type === 'text' || field.type === 'password'"
              v-model="configForm[field.key]" 
              :type="field.type" 
              :placeholder="field.placeholder" 
              show-password
            />
            <el-switch v-else-if="field.type === 'switch'" v-model="configForm[field.key]" />
            <el-select v-else-if="field.type === 'select'" v-model="configForm[field.key]" style="width: 100%">
              <el-option v-for="opt in field.options" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
          </el-form-item>
        </template>
        <el-form-item label="配置说明" v-if="!currentPlugin?.config_fields">
          <el-input v-model="configForm.description" type="textarea" :rows="3" placeholder="请输入配置说明" />
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
import { Grid } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const plugins = ref<any[]>([])
const configDialogVisible = ref(false)
const currentPlugin = ref<any>(null)
const configForm = ref<any>({})
const savingConfig = ref(false)

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/addons-plugins' })
    plugins.value = data || []
  } catch (error) {
    console.error('fetch addons plugins failed:', error)
  } finally {
    loading.value = false
  }
}

const handleToggle = async (plugin: any) => {
  const original = plugin.is_enabled
  try {
    await request.post({ 
      url: `/api/admin/addons-plugins/${plugin.id}/toggle`, 
      data: { enabled: plugin.is_enabled } 
    })
    ElMessage.success(plugin.is_enabled ? '已启用' : '已禁用')
  } catch (error) {
    plugin.is_enabled = original
    ElMessage.error('操作失败')
  }
}

const handleConfigure = async (plugin: any) => {
  currentPlugin.value = plugin
  try {
    const data = await request.get({ url: `/api/admin/addons-plugins/${plugin.id}/config` })
    configForm.value = data || { is_enabled: plugin.is_enabled, description: '' }
  } catch (error) {
    configForm.value = { is_enabled: plugin.is_enabled, description: '' }
  }
  configDialogVisible.value = true
}

const handleUninstall = async (plugin: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要卸载 "${plugin.title}" 吗？卸载后相关数据将被清除。`, 
      '确认卸载', 
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await request.post({ url: `/api/admin/addons-plugins/${plugin.id}/uninstall` })
    ElMessage.success('已卸载')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('卸载失败')
    }
  }
}

const handleSaveConfig = async () => {
  savingConfig.value = true
  try {
    await request.post({ 
      url: `/api/admin/addons-plugins/${currentPlugin.value.id}/config`, 
      data: configForm.value 
    })
    ElMessage.success('保存成功')
    configDialogVisible.value = false
    fetchList()
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
.plugins-page {
  padding: 20px;
}

.page-title {
  margin: 0 0 16px 0;
  font-size: 20px;
  font-weight: 600;
}

.page-desc {
  margin-bottom: 16px;
}

.plugin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.plugin-card {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 16px;
  transition: all 0.3s;
  
  &:hover {
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  }
  
  &.plugin-disabled {
    opacity: 0.7;
  }
}

.plugin-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.plugin-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-color-primary-light-9);
  border-radius: 8px;
  flex-shrink: 0;
}

.plugin-info {
  flex: 1;
  min-width: 0;
}

.plugin-name {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.plugin-author {
  margin: 0;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.plugin-body {
  margin-bottom: 12px;
}

.plugin-desc {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: var(--el-text-color-regular);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.plugin-meta {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.plugin-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid var(--el-border-color-lighter);
}

.plugin-actions {
  display: flex;
  gap: 8px;
}
</style>
