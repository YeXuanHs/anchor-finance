<template>
  <div class="security-config-page">
    <art-card title="安全配置" shadow="never">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 实名认证设置 -->
        <el-tab-pane label="实名认证设置" name="settings">
          <el-form :model="settingForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item label="开启实名认证">
              <el-switch v-model="settingForm.certifi_open" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="认证方式">
              <el-checkbox-group v-model="settingForm.certifi_select">
                <el-checkbox v-for="item in certifiTypes" :key="item.value" :label="item.value">{{ item.name }}</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="允许上传证件">
              <el-switch v-model="settingForm.certifi_is_upload" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="未实名暂停产品">
              <el-switch v-model="settingForm.certifi_is_stop" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="暂停期限(天)" v-if="settingForm.certifi_is_stop">
              <el-input-number v-model="settingForm.certifi_stop_day" :min="1" :max="365" />
            </el-form-item>
            <el-form-item label="同步实名到姓名">
              <el-switch v-model="settingForm.certifi_realname" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="绑定手机一致性">
              <el-switch v-model="settingForm.certifi_isbindphone" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="自动发送认证提醒">
              <el-switch v-model="settingForm.artificial_auto_send_msg" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveSettings" :loading="saveLoading">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 接口设置 -->
        <el-tab-pane label="接口设置" name="api">
          <el-form :model="apiForm" label-width="160px" style="max-width: 700px" v-loading="loading">
            <el-form-item label="认证接口类型">
              <el-select v-model="apiForm.certifi_type" placeholder="选择认证接口">
                <el-option v-for="item in apiTypes" :key="item.value" :label="item.name" :value="item.value" />
              </el-select>
            </el-form-item>
            <template v-if="apiForm.certifi_type === 'ali'">
              <el-form-item label="APP ID">
                <el-input v-model="apiForm.certifi_app_id" placeholder="支付宝AppID" />
              </el-form-item>
              <el-form-item label="支付宝公钥">
                <el-input v-model="apiForm.certifi_alipay_public_key" type="textarea" :rows="3" placeholder="支付宝公钥" />
              </el-form-item>
              <el-form-item label="商户私钥">
                <el-input v-model="apiForm.certifi_merchant_private_key" type="textarea" :rows="3" placeholder="商户私钥" />
              </el-form-item>
              <el-form-item label="认证方式">
                <el-select v-model="apiForm.certifi_alipay_biz_code" placeholder="选择认证方式">
                  <el-option v-for="item in alipayBizCodes" :key="item.value" :label="item.name" :value="item.value" />
                </el-select>
              </el-form-item>
            </template>
            <template v-if="apiForm.certifi_type === 'three' || apiForm.certifi_type === 'phonethree'">
              <el-form-item label="三要素类型">
                <el-select v-model="apiForm.certifi_three_type" placeholder="选择类型">
                  <el-option v-for="item in threeTypes" :key="item.value" :label="item.name" :value="item.value" />
                </el-select>
              </el-form-item>
              <el-form-item label="AppCode">
                <el-input v-model="apiForm.certifi_appcode" placeholder="请输入AppCode" />
              </el-form-item>
            </template>
            <el-form-item>
              <el-button type="primary" @click="handleSaveApi" :loading="saveLoading">保存设置</el-button>
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

const activeTab = ref('settings')
const loading = ref(false)
const saveLoading = ref(false)

const certifiTypes = ref<Array<{ name: string; value: string }>>([])
const apiTypes = ref<Array<{ name: string; value: string }>>([])
const threeTypes = ref([
  { name: '两要素', value: 'two' },
  { name: '三要素', value: 'three' },
  { name: '四要素', value: 'four' }
])
const alipayBizCodes = ref([
  { name: '快捷认证', value: 'SMART_FACE' },
  { name: '人脸识别', value: 'FACE' },
  { name: '身份证识别', value: 'CERT_PHOTO' },
  { name: '人脸+身份证', value: 'CERT_PHOTO_FACE' }
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
    ElMessage.error('保存失败')
  } finally {
    saveLoading.value = false
  }
}

const handleSaveApi = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/config/certifi', data: apiForm, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saveLoading.value = false
  }
}

const handleTabChange = (tab: string) => {
  if (tab === 'api') fetchConfig()
}

onMounted(() => fetchConfig())
</script>

<style scoped lang="scss">
.security-config-page {
  padding: 20px;
}
</style>
