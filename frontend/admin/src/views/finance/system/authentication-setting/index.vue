<template>
  <div class="authentication-setting-page">
    <h2>实名认证</h2>
    <p class="subtitle">管理实名认证插件，配置认证方式</p>

    <!-- 插件卡片网格 -->
    <div class="plugin-grid">
      <el-card v-for="plugin in plugins" :key="plugin.id" shadow="hover" class="plugin-card">
        <div class="plugin-header">
          <div class="plugin-icon">
            <el-icon :size="32"><User /></el-icon>
          </div>
          <el-switch v-model="plugin.is_enabled" @change="handleToggle(plugin)" />
        </div>
        <h3>{{ plugin.title }}</h3>
        <p>{{ plugin.description }}</p>
        <div class="plugin-actions">
          <el-button type="primary" link @click="handleConfig(plugin)">配置</el-button>
          <el-button v-if="!plugin.is_system" type="danger" link @click="handleUninstall(plugin)">卸载</el-button>
        </div>
      </el-card>
    </div>

    <!-- 配置对话框 -->
    <el-dialog v-model="configVisible" :title="`${currentPlugin?.title} - 配置`" width="600px">
      <el-form :model="configForm" label-width="120px">
        <!-- 支付宝实名认证配置 -->
        <template v-if="currentPlugin?.name === 'AlipayCertify'">
          <el-form-item label="AppID">
            <el-input v-model="configForm.appid" placeholder="请输入支付宝AppID" />
          </el-form-item>
          <el-form-item label="应用私钥">
            <el-input v-model="configForm.private_key" type="textarea" :rows="3" placeholder="请输入应用私钥" />
          </el-form-item>
          <el-form-item label="支付宝公钥">
            <el-input v-model="configForm.alipay_public_key" type="textarea" :rows="3" placeholder="请输入支付宝公钥" />
          </el-form-item>
        </template>

        <!-- 人工审核配置 -->
        <template v-if="currentPlugin?.name === 'ManualVerify'">
          <el-form-item label="审核说明">
            <el-input v-model="configForm.description" type="textarea" :rows="3" placeholder="显示给用户的审核说明" />
          </el-form-item>
          <el-form-item label="允许证件类型">
            <el-checkbox-group v-model="configForm.id_types">
              <el-checkbox label="id_card">身份证</el-checkbox>
              <el-checkbox label="passport">护照</el-checkbox>
              <el-checkbox label="driver_license">驾驶证</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </template>

        <!-- 通用配置 -->
        <el-form-item label="备注">
          <el-input v-model="configForm.remark" type="textarea" :rows="2" placeholder="管理员备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const plugins = ref<any[]>([])
const configVisible = ref(false)
const currentPlugin = ref<any>(null)
const configForm = reactive<any>({
  appid: '',
  private_key: '',
  alipay_public_key: '',
  description: '',
  id_types: ['id_card'],
  remark: ''
})

// 获取插件列表
const fetchPlugins = async () => {
  try {
    const data = await request.get({ url: '/api/admin/plugins', params: { type: 'certification' } })
    plugins.value = data || []
  } catch (error) {
    console.error('获取插件列表失败:', error)
  }
}

// 切换启用状态
const handleToggle = async (plugin: any) => {
  try {
    await request.put({
      url: `/api/admin/plugins/${plugin.id}/toggle`,
      data: { is_enabled: plugin.is_enabled }
    })
    ElMessage.success(plugin.is_enabled ? '已启用' : '已禁用')
  } catch (error) {
    plugin.is_enabled = !plugin.is_enabled
    ElMessage.error('操作失败')
  }
}

// 打开配置
const handleConfig = (plugin: any) => {
  currentPlugin.value = plugin
  Object.assign(configForm, plugin.config || {})
  configVisible.value = true
}

// 保存配置
const handleSaveConfig = async () => {
  try {
    await request.put({
      url: `/api/admin/plugins/${currentPlugin.value.id}/config`,
      data: configForm
    })
    ElMessage.success('配置保存成功')
    configVisible.value = false
    fetchPlugins()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 卸载插件
const handleUninstall = async (plugin: any) => {
  try {
    await ElMessageBox.confirm(`确定要卸载插件 "${plugin.title}" 吗？`, '卸载确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.del({ url: `/api/admin/plugins/${plugin.id}` })
    ElMessage.success('卸载成功')
    fetchPlugins()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('卸载失败')
    }
  }
}

onMounted(() => {
  fetchPlugins()
})
</script>

<style scoped lang="scss">
.authentication-setting-page {
  padding: 20px;

  h2 {
    margin: 0 0 8px;
    font-size: 18px;
  }

  .subtitle {
    color: var(--el-text-color-secondary);
    margin: 0 0 24px;
  }
}

.plugin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.plugin-card {
  .plugin-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .plugin-icon {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    background: var(--el-color-primary-light-9);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--el-color-primary);
  }

  h3 {
    margin: 0 0 8px;
    font-size: 16px;
  }

  p {
    color: var(--el-text-color-secondary);
    font-size: 14px;
    margin: 0 0 16px;
    min-height: 40px;
  }

  .plugin-actions {
    display: flex;
    gap: 12px;
  }
}
</style>
