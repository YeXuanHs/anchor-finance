<template>
  <div class="affiliate-config-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>推介计划配置</span>
          <el-button type="primary" @click="handleSaveConfig" :loading="saveLoading">
            <el-icon><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>

      <!-- 配置表单 -->
      <el-form :model="configForm" :rules="configRules" ref="configFormRef" label-width="160px" class="config-form">
        <el-divider content-position="left">基础配置</el-divider>

        <el-form-item label="启用推介计划" prop="is_enabled">
          <el-switch v-model="configForm.is_enabled" :active-value="1" :inactive-value="0" />
        </el-form-item>

        <el-form-item label="佣金类型" prop="commission_type">
          <el-radio-group v-model="configForm.commission_type">
            <el-radio value="percentage">百分比</el-radio>
            <el-radio value="fixed">固定金额</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="佣金比例/金额" prop="commission_value">
          <el-input-number v-model="configForm.commission_value" :min="0" :max="100" :precision="2" />
          <span class="suffix-text">{{ configForm.commission_type === 'percentage' ? '%' : '元' }}</span>
        </el-form-item>

        <el-form-item label="最低提现金额" prop="min_withdraw">
          <el-input-number v-model="configForm.min_withdraw" :min="0" :precision="2" />
          <span class="suffix-text">元</span>
        </el-form-item>

        <el-divider content-position="left">规则配置</el-divider>

        <el-form-item label="Cookie有效期" prop="cookie_days">
          <el-input-number v-model="configForm.cookie_days" :min="1" :max="365" />
          <span class="suffix-text">天</span>
        </el-form-item>

        <el-form-item label="佣金结算周期" prop="settlement_cycle">
          <el-select v-model="configForm.settlement_cycle" placeholder="请选择">
            <el-option label="即时结算" value="instant" />
            <el-option label="订单完成后" value="completed" />
            <el-option label="每月结算" value="monthly" />
          </el-select>
        </el-form-item>

        <el-form-item label="结算延迟天数" prop="settlement_delay">
          <el-input-number v-model="configForm.settlement_delay" :min="0" :max="90" />
          <span class="suffix-text">天</span>
        </el-form-item>

        <el-divider content-position="left">推广配置</el-divider>

        <el-form-item label="推广链接格式" prop="link_format">
          <el-input v-model="configForm.link_format" placeholder="例如: https://example.com?ref={affiliate_code}" />
        </el-form-item>

        <el-form-item label="允许自定义推广码" prop="allow_custom_code">
          <el-switch v-model="configForm.allow_custom_code" :active-value="1" :inactive-value="0" />
        </el-form-item>

        <el-form-item label="推荐注册奖励" prop="register_bonus">
          <el-input-number v-model="configForm.register_bonus" :min="0" :precision="2" />
          <span class="suffix-text">元</span>
        </el-form-item>

        <el-divider content-position="left">通知配置</el-divider>

        <el-form-item label="佣金到账通知" prop="notify_commission">
          <el-switch v-model="configForm.notify_commission" :active-value="1" :inactive-value="0" />
        </el-form-item>

        <el-form-item label="提现申请通知" prop="notify_withdraw">
          <el-switch v-model="configForm.notify_withdraw" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const saveLoading = ref(false)
const configFormRef = ref<FormInstance>()

const configForm = reactive({
  is_enabled: 1,
  commission_type: 'percentage',
  commission_value: 10,
  min_withdraw: 100,
  cookie_days: 30,
  settlement_cycle: 'completed',
  settlement_delay: 7,
  link_format: 'https://example.com?ref={affiliate_code}',
  allow_custom_code: 0,
  register_bonus: 0,
  notify_commission: 1,
  notify_withdraw: 1
})

const configRules: FormRules = {
  commission_value: [{ required: true, message: '请输入佣金比例/金额', trigger: 'blur' }],
  min_withdraw: [{ required: true, message: '请输入最低提现金额', trigger: 'blur' }],
  cookie_days: [{ required: true, message: '请输入Cookie有效期', trigger: 'blur' }]
}

const fetchConfig = async () => {
  try {
    const data = await request.get({ url: '/api/admin/affiliate-config' })
    Object.assign(configForm, data)
  } catch (error) {
    console.error('获取配置失败:', error)
  }
}

const handleSaveConfig = async () => {
  if (!configFormRef.value) return
  await configFormRef.value.validate(async (valid) => {
    if (!valid) return
    saveLoading.value = true
    try {
      await request.put({ url: '/api/admin/affiliate-config', params: configForm })
      ElMessage.success('保存成功')
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      saveLoading.value = false
    }
  })
}

onMounted(() => { fetchConfig() })
</script>

<style scoped lang="scss">
.affiliate-config-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-form {
  max-width: 800px;

  .suffix-text {
    margin-left: 8px;
    color: #909399;
  }
}
</style>
