<template>
  <div class="plugins-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部插件" name="all" />
      <el-tab-pane label="认证插件" name="authentication" />
      <el-tab-pane label="支付接口" name="payment" />
      <el-tab-pane label="实名认证" name="verification" />
      <el-tab-pane label="短信接口" name="sms" />
      <el-tab-pane label="其他插件" name="other" />
    </el-tabs>

    <!-- 插件卡片网格 -->
    <div class="plugins-grid" v-loading="loading">
      <el-card
        v-for="plugin in filteredPlugins"
        :key="plugin.id"
        shadow="hover"
        class="plugin-card"
        :class="{ 'is-enabled': plugin.enabled }"
      >
        <div class="plugin-header">
          <div class="plugin-icon">
            <el-icon :size="32" :color="plugin.enabled ? 'var(--el-color-primary)' : '#86909C'">
              <component :is="plugin.icon || 'Setting'" />
            </el-icon>
          </div>
          <div class="plugin-status">
            <el-switch
              v-model="plugin.enabled"
              @change="handleTogglePlugin(plugin)"
            />
          </div>
        </div>
        <div class="plugin-body">
          <h3 class="plugin-name">{{ plugin.name }}</h3>
          <p class="plugin-desc">{{ plugin.description }}</p>
          <div class="plugin-version">v{{ plugin.version }}</div>
        </div>
        <div class="plugin-footer">
          <el-button
            type="primary"
            size="small"
            @click="handleConfigure(plugin)"
          >
            配置
          </el-button>
          <el-button
            v-if="plugin.has_update"
            type="warning"
            size="small"
            @click="handleUpdate(plugin)"
          >
            更新
          </el-button>
        </div>
      </el-card>

      <!-- 空状态 -->
      <el-empty v-if="filteredPlugins.length === 0" description="暂无插件" />
    </div>

    <!-- 配置弹窗 -->
    <el-dialog
      v-model="configDialogVisible"
      :title="`${currentPlugin?.name} - 配置`"
      width="600px"
    >
      <el-form
        v-if="currentPlugin?.config_fields"
        ref="configFormRef"
        :model="configFormData"
        label-width="120px"
      >
        <el-form-item
          v-for="field in currentPlugin.config_fields"
          :key="field.key"
          :label="field.label"
          :prop="field.key"
          :rules="field.required ? [{ required: true, message: `请输入${field.label}`, trigger: 'blur' }] : []"
        >
          <!-- 输入框 -->
          <el-input
            v-if="field.type === 'text' || field.type === 'password'"
            v-model="configFormData[field.key]"
            :type="field.type"
            :placeholder="field.placeholder || `请输入${field.label}`"
            show-password
          />

          <!-- 文本域 -->
          <el-input
            v-else-if="field.type === 'textarea'"
            v-model="configFormData[field.key]"
            type="textarea"
            :rows="3"
            :placeholder="field.placeholder || `请输入${field.label}`"
          />

          <!-- 下拉选择 -->
          <el-select
            v-else-if="field.type === 'select'"
            v-model="configFormData[field.key]"
            :placeholder="`请选择${field.label}`"
          >
            <el-option
              v-for="option in field.options"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>

          <!-- 开关 -->
          <el-switch
            v-else-if="field.type === 'switch'"
            v-model="configFormData[field.key]"
          />

          <!-- 数字 -->
          <el-input-number
            v-else-if="field.type === 'number'"
            v-model="configFormData[field.key]"
            :min="field.min"
            :max="field.max"
          />
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
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Setting } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const plugins = ref<any[]>([])
const activeTab = ref('all')

// 配置弹窗
const configDialogVisible = ref(false)
const currentPlugin = ref<any>(null)
const configFormData = ref<Record<string, any>>({})
const configFormRef = ref<FormInstance>()
const savingConfig = ref(false)

// 过滤后的插件
const filteredPlugins = computed(() => {
  if (activeTab.value === 'all') return plugins.value
  return plugins.value.filter((p) => p.category === activeTab.value)
})

// 标签页切换
const handleTabChange = () => {
  // 切换标签页时不需要重新加载，使用 computed 过滤
}

// 获取插件列表
const fetchPlugins = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/plugins' })
    plugins.value = data || []
  } catch (error) {
    console.error('获取插件列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 切换插件状态
const handleTogglePlugin = async (plugin: any) => {
  try {
    await request.post({
      url: `/api/admin/plugins/${plugin.id}/toggle`,
      data: { enabled: plugin.enabled }
    })
    ElMessage.success(plugin.enabled ? '插件已启用' : '插件已禁用')
  } catch (error) {
    console.error('切换插件状态失败:', error)
    plugin.enabled = !plugin.enabled
  }
}

// 配置插件
const handleConfigure = async (plugin: any) => {
  currentPlugin.value = plugin

  // 获取当前配置
  try {
    const data = await request.get({ url: `/api/admin/plugins/${plugin.id}/config` })
    configFormData.value = data || {}
  } catch (error) {
    configFormData.value = {}
    console.error('获取插件配置失败:', error)
  }

  configDialogVisible.value = true
}

// 保存配置
const handleSaveConfig = async () => {
  if (!configFormRef.value) return

  try {
    await configFormRef.value.validate()
    savingConfig.value = true

    await request.post({
      url: `/api/admin/plugins/${currentPlugin.value.id}/config`,
      data: configFormData.value
    })
    ElMessage.success('配置保存成功')
    configDialogVisible.value = false
  } catch (error) {
    console.error('保存配置失败:', error)
  } finally {
    savingConfig.value = false
  }
}

// 更新插件
const handleUpdate = async (plugin: any) => {
  try {
    await ElMessageBox.confirm(`确定要更新插件 "${plugin.name}" 吗？`, '确认更新', {
      type: 'warning'
    })
    await request.post({ url: `/api/admin/plugins/${plugin.id}/update` })
    ElMessage.success('插件更新成功')
    fetchPlugins()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新插件失败:', error)
    }
  }
}

onMounted(() => {
  fetchPlugins()
})
</script>

<style scoped lang="scss">
.plugins-page {
  padding: 16px;
}

.plugins-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.plugin-card {
  transition: all 0.3s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  &.is-enabled {
    border-color: var(--el-color-primary-light-5);
  }

  :deep(.el-card__body) {
    padding: 16px;
  }
}

.plugin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.plugin-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: #F2F3F5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.plugin-body {
  margin-bottom: 12px;
}

.plugin-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1D2129;
}

.plugin-desc {
  margin: 0 0 8px 0;
  font-size: 13px;
  color: #86909C;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.plugin-version {
  font-size: 12px;
  color: #C0C4CC;
}

.plugin-footer {
  display: flex;
  gap: 8px;
}
</style>
