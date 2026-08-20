<template>
  <div class="security-config-page">
    <art-card :title="$t('securityConfig.title')" shadow="never">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 实名认证设置 -->
        <el-tab-pane :label="$t('securityConfig.certSettings')" name="settings">
          <el-form :model="settingForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item :label="$t('securityConfig.enableCert')">
              <el-switch v-model="settingForm.certifi_open" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item :label="$t('securityConfig.certMethod')">
              <el-checkbox-group v-model="settingForm.certifi_select">
                <el-checkbox v-for="item in certifiTypes" :key="item.value" :label="item.value">{{ item.name }}</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item :label="$t('securityConfig.allowUpload')">
              <el-switch v-model="settingForm.certifi_is_upload" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item :label="$t('securityConfig.pauseOnUncertified')">
              <el-switch v-model="settingForm.certifi_is_stop" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item :label="$t('securityConfig.pauseDays')" v-if="settingForm.certifi_is_stop">
              <el-input-number v-model="settingForm.certifi_stop_day" :min="1" :max="365" />
            </el-form-item>
            <el-form-item :label="$t('securityConfig.syncToName')">
              <el-switch v-model="settingForm.certifi_realname" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item :label="$t('securityConfig.bindPhoneConsistency')">
              <el-switch v-model="settingForm.certifi_isbindphone" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item :label="$t('securityConfig.autoSendReminder')">
              <el-switch v-model="settingForm.artificial_auto_send_msg" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveSettings" :loading="saveLoading">{{ $t('securityConfig.saveSettings') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 接口设置 -->
        <el-tab-pane :label="$t('securityConfig.apiSettings')" name="api">
          <el-form :model="apiForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item :label="$t('securityConfig.certApiType')">
              <el-select v-model="apiForm.certifi_type" :placeholder="$t('securityConfig.selectCertApi')">
                <el-option v-for="item in apiTypes" :key="item.value" :label="item.name" :value="item.value" />
              </el-select>
            </el-form-item>
            <template v-if="apiForm.certifi_type === 'ali'">
              <el-form-item label="APP ID">
                <el-input v-model="apiForm.certifi_app_id" :placeholder="$t('securityConfig.alipayAppIdPlaceholder')" />
              </el-form-item>
              <el-form-item :label="$t('securityConfig.alipayPublicKey')">
                <el-input v-model="apiForm.certifi_alipay_public_key" type="textarea" :rows="3" :placeholder="$t('securityConfig.alipayPublicKey')" />
              </el-form-item>
              <el-form-item :label="$t('securityConfig.merchantPrivateKey')">
                <el-input v-model="apiForm.certifi_merchant_private_key" type="textarea" :rows="3" :placeholder="$t('securityConfig.merchantPrivateKey')" />
              </el-form-item>
              <el-form-item :label="$t('securityConfig.certMethodLabel')">
                <el-select v-model="apiForm.certifi_alipay_biz_code" :placeholder="$t('securityConfig.selectCertMethod')">
                  <el-option v-for="item in alipayBizCodes" :key="item.value" :label="item.name" :value="item.value" />
                </el-select>
              </el-form-item>
            </template>
            <template v-if="apiForm.certifi_type === 'three' || apiForm.certifi_type === 'phonethree'">
              <el-form-item :label="$t('securityConfig.threeElementType')">
                <el-select v-model="apiForm.certifi_three_type" :placeholder="$t('securityConfig.selectType')">
                  <el-option v-for="item in threeTypes" :key="item.value" :label="item.name" :value="item.value" />
                </el-select>
              </el-form-item>
              <el-form-item label="AppCode">
                <el-input v-model="apiForm.certifi_appcode" :placeholder="$t('securityConfig.inputAppCode')" />
              </el-form-item>
            </template>
            <el-form-item>
              <el-button type="primary" @click="handleSaveApi" :loading="saveLoading">{{ $t('securityConfig.saveSettings') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const activeTab = ref('settings')
const loading = ref(false)
const saveLoading = ref(false)

const certifiTypes = ref<Array<{ name: string; value: string }>>([])
const apiTypes = ref<Array<{ name: string; value: string }>>([])
const threeTypes = computed(() => [
  { name: $t('securityConfig.twoFactor'), value: 'two' },
  { name: $t('securityConfig.threeFactor'), value: 'three' },
  { name: $t('securityConfig.fourFactor'), value: 'four' }
])
const alipayBizCodes = computed(() => [
  { name: $t('securityConfig.quickCert'), value: 'SMART_FACE' },
  { name: $t('securityConfig.faceRecognition'), value: 'FACE' },
  { name: $t('securityConfig.idCardRecognition'), value: 'CERT_PHOTO' },
  { name: $t('securityConfig.faceIdCard'), value: 'CERT_PHOTO_FACE' }
])

const settingForm = reactive({
  certifi_open: 0,
  certifi_select: [] as string[],
  certifi_is_upload: 0,
  certifi_is_stop: 0,
  certifi_stop_day: 7,
  certifi_realname: 0,
  certifi_isbindphone: 0,
  artificial_auto_send_msg: 0
})

const apiForm = reactive({
  certifi_type: '',
  certifi_app_id: '',
  certifi_alipay_public_key: '',
  certifi_merchant_private_key: '',
  certifi_alipay_biz_code: '',
  certifi_three_type: '',
  certifi_appcode: ''
})

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/config/certifi' })
    if (res?.data) {
      const data = res.data
      Object.assign(settingForm, {
        certifi_open: parseInt(data.certifi_open) || 0,
        certifi_select: data.certifi_select ? data.certifi_select.split(',') : [],
        certifi_is_upload: parseInt(data.certifi_is_upload) || 0,
        certifi_is_stop: parseInt(data.certifi_is_stop) || 0,
        certifi_stop_day: parseInt(data.certifi_stop_day) || 7,
        certifi_realname: parseInt(data.certifi_realname) || 0,
        certifi_isbindphone: parseInt(data.certifi_isbindphone) || 0,
        artificial_auto_send_msg: parseInt(data.artificial_auto_send_msg) || 0
      })
    }
    if (res?.types) apiTypes.value = res.types
    if (res?.certifi_select_all) {
      certifiTypes.value = Object.entries(res.certifi_select_all).map(([value, name]) => ({ name: name as string, value }))
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSaveSettings = async () => {
  saveLoading.value = true
  try {
    await request.put({
      url: '/api/admin/config/certifi',
      data: { ...settingForm, certifi_select: settingForm.certifi_select.join(',') },
      showSuccessMessage: true
    })
  } catch (error) {
    ElMessage.error($t('securityConfig.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

const handleSaveApi = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/config/certifi', data: apiForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error($t('securityConfig.saveFailed'))
  } finally {
    saveLoading.value = false
  }
}

const handleTabChange = (tab: string | number) => {
  if (tab === 'api') fetchConfig()
}

onMounted(() => fetchConfig())
</script>

<style scoped lang="scss">
.security-config-page {
  padding: 20px;
}
</style>
