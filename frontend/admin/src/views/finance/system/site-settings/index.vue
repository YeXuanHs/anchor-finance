<template>
  <div class="page-container">
    <art-card title="站点设置" shadow="never">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本设置" name="basic">
          <el-form :model="basicForm" label-width="120px" style="max-width: 700px">
            <el-form-item label="站点名称">
              <el-input v-model="basicForm.site_name" />
            </el-form-item>
            <el-form-item label="站点URL">
              <el-input v-model="basicForm.site_url" />
            </el-form-item>
            <el-form-item label="站点描述">
              <el-input v-model="basicForm.site_description" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item label="Logo">
              <el-input v-model="basicForm.logo" />
            </el-form-item>
            <el-form-item label="ICP备案">
              <el-input v-model="basicForm.icp" />
            </el-form-item>
            <el-form-item label="版权信息">
              <el-input v-model="basicForm.copyright" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveBasic">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="联系信息" name="contact">
          <el-form :model="contactForm" label-width="120px" style="max-width: 700px">
            <el-form-item label="联系电话">
              <el-input v-model="contactForm.contact_phone" />
            </el-form-item>
            <el-form-item label="联系邮箱">
              <el-input v-model="contactForm.contact_email" />
            </el-form-item>
            <el-form-item label="联系地址">
              <el-input v-model="contactForm.contact_address" />
            </el-form-item>
            <el-form-item label="QQ">
              <el-input v-model="contactForm.qq" />
            </el-form-item>
            <el-form-item label="微信">
              <el-input v-model="contactForm.wechat" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveContact">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="主题设置" name="theme">
          <el-form :model="themeForm" label-width="120px" style="max-width: 700px">
            <el-form-item label="主题色">
              <el-color-picker v-model="themeForm.primary_color" />
            </el-form-item>
            <el-form-item label="圆角大小">
              <el-slider v-model="themeForm.border_radius" :min="0" :max="24" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveTheme">保存</el-button>
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

const activeTab = ref('basic')
const basicForm = ref({
  site_name: '',
  site_url: '',
  site_description: '',
  logo: '',
  icp: '',
  copyright: ''
})
const contactForm = ref({
  contact_phone: '',
  contact_email: '',
  contact_address: '',
  qq: '',
  wechat: ''
})
const themeForm = ref({
  primary_color: '#1890ff',
  border_radius: 12
})

const fetchSettings = async () => {
  try {
    const { data } = await request.get('/admin/site-settings')
    if (data?.data) {
      Object.assign(basicForm.value, data.data)
      Object.assign(contactForm.value, data.data)
      Object.assign(themeForm.value, data.data)
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSaveBasic = async () => {
  try {
    await request.put('/admin/site-settings', basicForm.value)
    ElMessage.success('保存成功')
  } catch (error) {
    console.error(error)
  }
}

const handleSaveContact = async () => {
  try {
    await request.put('/admin/site-settings', contactForm.value)
    ElMessage.success('保存成功')
  } catch (error) {
    console.error(error)
  }
}

const handleSaveTheme = async () => {
  try {
    await request.put('/admin/site-settings', themeForm.value)
    ElMessage.success('保存成功')
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchSettings())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
