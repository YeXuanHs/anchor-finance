<template>
  <div class="fund-config-page">
    <art-card :title="$t('page.fundConfig.title')" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- 充值设置 -->
        <el-tab-pane :label="$t('page.fundConfig.rechargeSettings')" name="recharge">
          <el-form :model="rechargeForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item :label="$t('page.fundConfig.enableRecharge')">
              <el-switch v-model="rechargeForm.addfunds_enabled" :active-value="1" :inactive-value="0" />
              <div class="form-tip">{{ $t('page.fundConfig.enableRechargeTip') }}</div>
            </el-form-item>
            <el-form-item :label="$t('page.fundConfig.minRecharge')">
              <el-input-number v-model="rechargeForm.addfunds_minimum" :min="0" :precision="2" />
              <div class="form-tip">{{ $t('page.fundConfig.minRechargeTip') }}</div>
            </el-form-item>
            <el-form-item :label="$t('page.fundConfig.maxRecharge')">
              <el-input-number v-model="rechargeForm.addfunds_maximum" :min="0" :precision="2" />
              <div class="form-tip">{{ $t('page.fundConfig.maxRechargeTip') }}</div>
            </el-form-item>
            <el-form-item :label="$t('page.fundConfig.maxBalance')">
              <el-input-number v-model="rechargeForm.addfunds_maximum_balance" :min="0" :precision="2" />
              <div class="form-tip">{{ $t('page.fundConfig.maxBalanceTip') }}</div>
            </el-form-item>
            <el-form-item :label="$t('page.fundConfig.requireActive')">
              <el-switch v-model="rechargeForm.addfunds_require_active" :active-value="1" :inactive-value="0" />
              <div class="form-tip">{{ $t('page.fundConfig.requireActiveTip') }}</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveRecharge" :loading="saveLoading">{{ $t('page.fundConfig.saveSettings') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 财务设置 -->
        <el-tab-pane :label="$t('page.fundConfig.financeSettings')" name="finance">
          <el-form :model="financeForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item :label="$t('page.fundConfig.downgradeRefund')">
              <el-switch v-model="financeForm.upgrade_down_product_config" :active-value="1" :inactive-value="0" />
              <div class="form-tip">{{ $t('page.fundConfig.downgradeRefundTip') }}</div>
            </el-form-item>
            <el-form-item :label="$t('page.fundConfig.customInvoiceId')">
              <el-switch v-model="financeForm.allow_custom_invoice_id" :active-value="1" :inactive-value="0" />
              <div class="form-tip">{{ $t('page.fundConfig.customInvoiceIdTip') }}</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveFinance" :loading="saveLoading">{{ $t('page.fundConfig.saveSettings') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const activeTab = ref('recharge')
const loading = ref(false)
const saveLoading = ref(false)

const rechargeForm = reactive({
  addfunds_enabled: 0,
  addfunds_minimum: 0,
  addfunds_maximum: 0,
  addfunds_maximum_balance: 0,
  addfunds_require_active: 0
})

const financeForm = reactive({
  upgrade_down_product_config: 0,
  allow_custom_invoice_id: 0
})

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/config/recharge' })
    if (res?.data) {
      Object.assign(rechargeForm, res.data.recharge || {})
      Object.assign(financeForm, res.data.finance || {})
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSaveRecharge = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/config/recharge', data: rechargeForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error($t('page.fundConfig.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

const handleSaveFinance = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/config/recharge', data: financeForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error($t('page.fundConfig.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

onMounted(() => fetchConfig())
</script>

<style scoped lang="scss">
.fund-config-page {
  padding: 20px;
}
.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}
</style>
