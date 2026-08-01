<template>
  <div class="fund-config-page">
    <art-card title="资金设置" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- 充值设置 -->
        <el-tab-pane label="充值设置" name="recharge">
          <el-form :model="rechargeForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item label="启用充值功能">
              <el-switch v-model="rechargeForm.addfunds_enabled" :active-value="1" :inactive-value="0" />
              <div class="form-tip">开启后，客户可以在用户中心充值余额</div>
            </el-form-item>
            <el-form-item label="最小充值金额">
              <el-input-number v-model="rechargeForm.addfunds_minimum" :min="0" :precision="2" />
              <div class="form-tip">单笔充值的最小金额</div>
            </el-form-item>
            <el-form-item label="最大充值金额">
              <el-input-number v-model="rechargeForm.addfunds_maximum" :min="0" :precision="2" />
              <div class="form-tip">单笔充值的最大金额</div>
            </el-form-item>
            <el-form-item label="最高余额限制">
              <el-input-number v-model="rechargeForm.addfunds_maximum_balance" :min="0" :precision="2" />
              <div class="form-tip">客户余额上限，0表示不限制</div>
            </el-form-item>
            <el-form-item label="需要已激活订单">
              <el-switch v-model="rechargeForm.addfunds_require_active" :active-value="1" :inactive-value="0" />
              <div class="form-tip">开启后，客户需有已激活订单才能充值（防欺诈）</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveRecharge" :loading="saveLoading">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 财务设置 -->
        <el-tab-pane label="财务设置" name="finance">
          <el-form :model="financeForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item label="降级退款至余额">
              <el-switch v-model="financeForm.upgrade_down_product_config" :active-value="1" :inactive-value="0" />
              <div class="form-tip">产品降级时将差价退至客户余额</div>
            </el-form-item>
            <el-form-item label="自定义起始账单号">
              <el-switch v-model="financeForm.allow_custom_invoice_id" :active-value="1" :inactive-value="0" />
              <div class="form-tip">允许自定义账单编号起始值</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveFinance" :loading="saveLoading">保存设置</el-button>
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
    const res = await request.get({ url: '/api/admin/config/fund' })
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
    await request.put({ url: '/api/admin/config/fund/recharge', data: rechargeForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saveLoading.value = false
  }
}

const handleSaveFinance = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/config/fund/finance', data: financeForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error('保存失败')
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
