<template>
  <div class="page-container">
    <art-card shadow="never">
      <template #header>
        <div class="card-header"><span>{{ $t('siteSettings.title') }}</span></div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('siteSettings.official')" name="official">
          <el-form :model="officialForm" label-width="120px" style="max-width: 800px;">
            <el-form-item :label="$t('siteSettings.phone')"><el-input v-model="officialForm.phone" :placeholder="$t('siteSettings.enterPhone')" /></el-form-item>
            <el-form-item label="QQ"><el-input v-model="officialForm.qq" placeholder="QQ" /></el-form-item>
            <el-form-item :label="$t('siteSettings.address')"><el-input v-model="officialForm.address" :placeholder="$t('siteSettings.enterAddress')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.icp')"><el-input v-model="officialForm.icp" :placeholder="$t('siteSettings.enterIcp')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.coordinates')"><el-input v-model="officialForm.coordinates" :placeholder="$t('siteSettings.enterCoordinates')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.logo')"><el-input v-model="officialForm.logo" :placeholder="$t('siteSettings.enterLogo')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.keywords')"><el-input v-model="officialForm.keywords" :placeholder="$t('siteSettings.enterKeywords')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.description')"><el-input v-model="officialForm.description" type="textarea" :rows="3" :placeholder="$t('siteSettings.enterDescription')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.companyIntro')"><el-input v-model="officialForm.company_intro" type="textarea" :rows="4" :placeholder="$t('siteSettings.enterCompanyIntro')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.header')"><el-input v-model="officialForm.header" type="textarea" :rows="3" :placeholder="$t('siteSettings.enterHeader')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.footer')"><el-input v-model="officialForm.footer" type="textarea" :rows="3" :placeholder="$t('siteSettings.enterFooter')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.enableLoginHeaderFooter')"><el-switch v-model="officialForm.enable_login_header_footer" /></el-form-item>
            <el-form-item :label="$t('siteSettings.loginHeader')"><el-input v-model="officialForm.login_header" type="textarea" :rows="2" :placeholder="$t('siteSettings.enterLoginHeader')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.loginFooter')"><el-input v-model="officialForm.login_footer" type="textarea" :rows="2" :placeholder="$t('siteSettings.enterLoginFooter')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.widget')"><el-input v-model="officialForm.widget" type="textarea" :rows="2" :placeholder="$t('siteSettings.enterWidget')" /></el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSave" :loading="saving">{{ $t('common.save') }}</el-button>
              <el-button @click="handleReset">{{ $t('common.cancel') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('siteSettings.memberCenter')" name="member">
          <el-form :model="memberForm" label-width="120px" style="max-width: 800px; margin-top: 20px;">
            <el-form-item :label="$t('siteSettings.memberLogo')"><el-input v-model="memberForm.logo" :placeholder="$t('siteSettings.enterMemberLogo')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.memberName')"><el-input v-model="memberForm.name" :placeholder="$t('siteSettings.enterMemberName')" /></el-form-item>
            <el-form-item :label="$t('siteSettings.defaultLanguage')">
              <el-select v-model="memberForm.language" :placeholder="$t('siteSettings.selectLanguage')" style="width: 100%">
                <el-option :label="$t('siteSettings.langZhCn')" value="zh-CN" />
                <el-option :label="$t('siteSettings.langZhTw')" value="zh-TW" />
                <el-option label="English" value="en" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('siteSettings.allowRegister')"><el-switch v-model="memberForm.allow_register" /></el-form-item>
            <el-form-item :label="$t('siteSettings.emailVerify')"><el-switch v-model="memberForm.email_verify" /></el-form-item>
            <el-form-item :label="$t('siteSettings.phoneVerify')"><el-switch v-model="memberForm.phone_verify" /></el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveMember" :loading="savingMember">{{ $t('common.save') }}</el-button>
              <el-button @click="fetchMemberSettings">{{ $t('common.cancel') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const activeTab = ref('official')
const saving = ref(false)
const savingMember = ref(false)

const officialForm = ref({ phone: '', qq: '', address: '', icp: '', coordinates: '', logo: '', keywords: '', description: '', company_intro: '', header: '', footer: '', enable_login_header_footer: true, login_header: '', login_footer: '', widget: '' })
const memberForm = ref({ logo: '', name: '', language: 'zh-CN', allow_register: true, email_verify: false, phone_verify: false })

const fetchSettings = async () => { try { const res = await request.get({ url: '/api/admin/site-settings' }); if (res) officialForm.value = { ...officialForm.value, ...res } } catch {} }
const fetchMemberSettings = async () => { try { const res = await request.get({ url: '/api/admin/member-settings' }); if (res) memberForm.value = { ...memberForm.value, ...res } } catch {} }
const handleSave = async () => { saving.value = true; try { await request.put({ url: '/api/admin/site-settings', data: officialForm.value }); ElMessage.success($t('common.saveSuccess')) } catch {} finally { saving.value = false } }
const handleSaveMember = async () => { savingMember.value = true; try { await request.put({ url: '/api/admin/member-settings', data: memberForm.value }); ElMessage.success($t('common.saveSuccess')) } catch {} finally { savingMember.value = false } }
const handleReset = () => { fetchSettings() }

onMounted(() => { fetchSettings(); fetchMemberSettings() })
</script>

<style scoped lang="scss">
.page-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
