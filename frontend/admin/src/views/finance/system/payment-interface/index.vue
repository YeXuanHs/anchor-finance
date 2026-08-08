<template>
  <div class="payment-interface-page">
    <!-- 支付接口卡片网格 -->
    <div class="payment-grid" v-loading="loading">
      <el-card v-for="item in paymentList" :key="item.id" shadow="hover" class="payment-card" :class="{ 'is-enabled': item.enabled }">
        <div class="payment-header">
          <div class="payment-icon">
            <el-icon :size="32" :color="item.enabled ? 'var(--el-color-primary)' : '#86909C'"><CreditCard /></el-icon>
          </div>
          <el-switch v-model="item.enabled" @change="handleToggle(item)" />
        </div>
        <div class="payment-body">
          <h3>{{ item.name }}</h3>
          <p>{{ item.description }}</p>
        </div>
        <div class="payment-footer">
          <el-button type="primary" size="small" @click="handleConfigure(item)">配置</el-button>
        </div>
      </el-card>
    </div>

    <!-- 配置弹窗 -->
    <el-dialog v-model="configDialogVisible" :title="`${currentItem?.name} - 配置`" width="600px">
      <el-form v-if="currentItem?.config_fields" ref="configFormRef" :model="configFormData" label-width="120px">
        <el-form-item v-for="field in currentItem.config_fields" :key="field.key" :label="field.label" :prop="field.key">
          <el-input v-if="field.type === 'text' || field.type === 'password'" v-model="configFormData[field.key]" :type="field.type" :placeholder="field.placeholder" show-password />
          <el-switch v-else-if="field.type === 'switch'" v-model="configFormData[field.key]" />
          <el-select v-else-if="field.type === 'select'" v-model="configFormData[field.key]">
            <el-option v-for="opt in field.options" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
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
import { ElMessage } from 'element-plus'
import { CreditCard } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const paymentList = ref<any[]>([])
const configDialogVisible = ref(false)
const currentItem = ref<any>(null)
const configFormData = ref<Record<string, any>>({})
const configFormRef = ref<FormInstance>()
const savingConfig = ref(false)

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/payment-gateways' })
    paymentList.value = data || []
  } catch (error) {
    console.error('获取支付接口失败:', error)
  } finally {
    loading.value = false
  }
}

const handleToggle = async (item: any) => {
  try {
    await request.post({ url: `/api/admin/payment-gateways/${item.id}/toggle`, data: { enabled: item.enabled } })
    ElMessage.success(item.enabled ? '已启用' : '已禁用')
  } catch (error) {
    console.error('切换状态失败:', error)
    item.enabled = !item.enabled
  }
}

const handleConfigure = async (item: any) => {
  currentItem.value = item
  try {
    const data = await request.get({ url: `/api/admin/payment-gateways/${item.id}/config` })
    configFormData.value = data || {}
  } catch (error) {
    configFormData.value = {}
  }
  configDialogVisible.value = true
}

const handleSaveConfig = async () => {
  if (!configFormRef.value) return
  try {
    await configFormRef.value.validate()
    savingConfig.value = true
    await request.post({ url: `/api/admin/payment-gateways/${currentItem.value.id}/config`, data: configFormData.value })
    ElMessage.success('配置保存成功')
    configDialogVisible.value = false
  } catch (error) {
    console.error('保存配置失败:', error)
  } finally {
    savingConfig.value = false
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.payment-interface-page { padding: 16px; }
.payment-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.payment-card { transition: all 0.3s; &:hover { transform: translateY(-4px); } &.is-enabled { border-color: var(--el-color-primary-light-5); } }
.payment-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; }
.payment-icon { width: 48px; height: 48px; border-radius: 8px; background: #F2F3F5; display: flex; align-items: center; justify-content: center; }
.payment-body { margin-bottom: 12px; h3 { margin: 0 0 8px; font-size: 16px; } p { margin: 0; font-size: 13px; color: #86909C; } }
</style>
