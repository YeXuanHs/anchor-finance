<template>
  <div class="affiliate-config-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsAffiliateConfig.title') }}</span>
          <el-button type="primary" @click="handleSave" :loading="submitLoading">
            <el-icon><Check /></el-icon>
            {{ $t('clientsAffiliateConfig.saveConfig') }}
          </el-button>
        </div>
      </template>

      <el-form :model="settingsForm" :rules="settingsRules" ref="settingsFormRef" label-width="140px" v-loading="loading">
        <el-divider content-position="left">{{ $t('clientsAffiliateConfig.basicConfig') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.enableAffiliate')" prop="is_open">
              <el-switch v-model="settingsForm.is_open" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.commissionType')" prop="type">
              <el-radio-group v-model="settingsForm.type">
                <el-radio :value="1">{{ $t('clientsAffiliateConfig.percentage') }}</el-radio>
                <el-radio :value="2">{{ $t('clientsAffiliateConfig.fixedAmount') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.commissionRate')" prop="ratio">
              <div class="input-with-unit">
                <el-input-number v-model="settingsForm.ratio" :min="0" :max="100" :precision="2" style="width: 100%" />
                <span class="unit">{{ settingsForm.type === 1 ? '%' : $t('clientsAffiliateConfig.yuan') }}</span>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.minWithdraw')" prop="min_amount">
              <div class="input-with-unit">
                <el-input-number v-model="settingsForm.min_amount" :min="0" :precision="2" style="width: 100%" />
                <span class="unit">{{ $t('clientsAffiliateConfig.yuan') }}</span>
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">{{ $t('clientsAffiliateConfig.ruleConfig') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.cookieDuration')" prop="expire">
              <div class="input-with-unit">
                <el-input-number v-model="settingsForm.expire" :min="1" :max="365" style="width: 100%" />
                <span class="unit">{{ $t('clientsAffiliateConfig.days') }}</span>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.settlementCycle')" prop="clear_type">
              <el-radio-group v-model="settingsForm.clear_type">
                <el-radio :value="1">{{ $t('clientsAffiliateConfig.instantSettlement') }}</el-radio>
                <el-radio :value="2">{{ $t('clientsAffiliateConfig.afterOrderComplete') }}</el-radio>
                <el-radio :value="3">{{ $t('clientsAffiliateConfig.monthlySettlement') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.settlementDelay')" prop="settlement_days">
              <div class="input-with-unit">
                <el-input-number v-model="settingsForm.settlement_days" :min="0" :max="365" style="width: 100%" />
                <span class="unit">{{ $t('clientsAffiliateConfig.days') }}</span>
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">{{ $t('clientsAffiliateConfig.promotionConfig') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.linkFormat')" prop="link_format">
              <el-input v-model="settingsForm.link_format" placeholder="https://example.com/?ref={code}" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.allowCustomCode')" prop="allow_custom_code">
              <el-switch v-model="settingsForm.allow_custom_code" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('clientsAffiliateConfig.registerBonus')" prop="register_bonus">
          <div class="input-with-unit">
            <el-input-number v-model="settingsForm.register_bonus" :min="0" :precision="2" style="width: 200px" />
            <span class="unit">{{ $t('clientsAffiliateConfig.yuan') }}</span>
          </div>
        </el-form-item>

        <el-divider content-position="left">{{ $t('clientsAffiliateConfig.notifyConfig') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.commissionNotice')" prop="commission_notice">
              <el-switch v-model="settingsForm.commission_notice" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('clientsAffiliateConfig.withdrawNotice')" prop="withdraw_notice">
              <el-switch v-model="settingsForm.withdraw_notice" />
            </el-form-item>
          </el-col>
        </el-row>
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
import { $t } from '@/locales'

const loading = ref(false)
const submitLoading = ref(false)
const settingsFormRef = ref<FormInstance>()

const settingsForm = reactive({
  is_open: 1, type: 1, ratio: 10, min_amount: 100, expire: 30, clear_type: 1, settlement_days: 0,
  link_format: 'https://example.com/?ref={code}', allow_custom_code: false, register_bonus: 10,
  commission_notice: true, withdraw_notice: true
})

const settingsRules: FormRules = {
  ratio: [{ required: true, message: $t('clientsAffiliateConfig.enterRate'), trigger: 'blur' }],
  min_amount: [{ required: true, message: $t('clientsAffiliateConfig.enterMinAmount'), trigger: 'blur' }]
}

const fetchSettings = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/affiliate/settings' }); Object.assign(settingsForm, data) } catch (e) { ElMessage.error($t('common.fetchFailed')) } finally { loading.value = false }
}

const handleSave = async () => {
  if (!settingsFormRef.value) return
  await settingsFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try { await request.put({ url: '/api/admin/affiliate/settings', params: settingsForm, showSuccessMessage: true }); ElMessage.success($t('common.saveSuccess')) } catch (e) { ElMessage.error($t('common.saveFailed')) } finally { submitLoading.value = false }
  })
}

onMounted(() => { fetchSettings() })
</script>

<style scoped lang="scss">
.affiliate-config-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.input-with-unit { display: flex; align-items: center; gap: 10px; .unit { color: #909399; font-size: 13px; } }
:deep(.el-divider__text) { font-weight: 600; color: #303133; }
</style>
