<template>
  <div class="general-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>常规设置</span>
        </div>
      </template>

      <el-form
        :model="formData"
        :rules="formRules"
        ref="formRef"
        label-width="120px"
        v-loading="loading"
        style="max-width: 800px"
      >
        <!-- 网站信息 -->
        <el-divider content-position="left">网站信息</el-divider>
        <el-form-item label="网站名称" prop="site_name">
          <el-input v-model="formData.site_name" placeholder="请输入网站名称" />
        </el-form-item>
        <el-form-item label="网站描述" prop="site_description">
          <el-input v-model="formData.site_description" type="textarea" :rows="3" placeholder="请输入网站描述" />
        </el-form-item>
        <el-form-item label="网站Logo" prop="site_logo">
          <el-input v-model="formData.site_logo" placeholder="请输入Logo URL地址" />
        </el-form-item>
        <el-form-item label="favicon" prop="site_favicon">
          <el-input v-model="formData.site_favicon" placeholder="请输入favicon URL地址" />
        </el-form-item>

        <!-- 联系方式 -->
        <el-divider content-position="left">联系方式</el-divider>
        <el-form-item label="联系电话" prop="contact_phone">
          <el-input v-model="formData.contact_phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="联系邮箱" prop="contact_email">
          <el-input v-model="formData.contact_email" placeholder="请输入联系邮箱" />
        </el-form-item>
        <el-form-item label="客服QQ" prop="contact_qq">
          <el-input v-model="formData.contact_qq" placeholder="请输入客服QQ号" />
        </el-form-item>
        <el-form-item label="客服微信" prop="contact_wechat">
          <el-input v-model="formData.contact_wechat" placeholder="请输入客服微信号" />
        </el-form-item>
        <el-form-item label="联系地址" prop="contact_address">
          <el-input v-model="formData.contact_address" placeholder="请输入联系地址" />
        </el-form-item>

        <!-- 公司信息 -->
        <el-divider content-position="left">公司信息</el-divider>
        <el-form-item label="公司名称" prop="company_name">
          <el-input v-model="formData.company_name" placeholder="请输入公司名称" />
        </el-form-item>
        <el-form-item label="公司简称" prop="company_short_name">
          <el-input v-model="formData.company_short_name" placeholder="请输入公司简称" />
        </el-form-item>
        <el-form-item label="统一社会信用代码" prop="company_code">
          <el-input v-model="formData.company_code" placeholder="请输入统一社会信用代码" />
        </el-form-item>
        <el-form-item label="公司地址" prop="company_address">
          <el-input v-model="formData.company_address" placeholder="请输入公司注册地址" />
        </el-form-item>
        <el-form-item label="法人代表" prop="company_legal_person">
          <el-input v-model="formData.company_legal_person" placeholder="请输入法人代表" />
        </el-form-item>

        <!-- 备案信息 -->
        <el-divider content-position="left">备案信息</el-divider>
        <el-form-item label="ICP备案号" prop="icp_number">
          <el-input v-model="formData.icp_number" placeholder="请输入ICP备案号" />
        </el-form-item>
        <el-form-item label="公安备案号" prop="police_number">
          <el-input v-model="formData.police_number" placeholder="请输入公安备案号" />
        </el-form-item>
        <el-form-item label="版权信息" prop="copyright">
          <el-input v-model="formData.copyright" placeholder="请输入版权信息" />
        </el-form-item>

        <!-- 其他设置 -->
        <el-divider content-position="left">其他设置</el-divider>
        <el-form-item label="网站状态" prop="site_status">
          <el-switch v-model="formData.site_status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="关闭提示" prop="site_close_tips" v-if="formData.site_status === 0">
          <el-input v-model="formData.site_close_tips" type="textarea" :rows="3" placeholder="网站关闭时显示的提示信息" />
        </el-form-item>
        <el-form-item label="统计代码" prop="analytics_code">
          <el-input v-model="formData.analytics_code" type="textarea" :rows="4" placeholder="请输入网站统计代码" />
        </el-form-item>
        <el-form-item label="自定义CSS" prop="custom_css">
          <el-input v-model="formData.custom_css" type="textarea" :rows="4" placeholder="请输入自定义CSS代码" />
        </el-form-item>
        <el-form-item label="页脚HTML" prop="footer_html">
          <el-input v-model="formData.footer_html" type="textarea" :rows="4" placeholder="请输入页脚自定义HTML" />
        </el-form-item>

        <!-- 提交 -->
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">保存设置</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

interface GeneralConfig {
  site_name: string
  site_description: string
  site_logo: string
  site_favicon: string
  contact_phone: string
  contact_email: string
  contact_qq: string
  contact_wechat: string
  contact_address: string
  company_name: string
  company_short_name: string
  company_code: string
  company_address: string
  company_legal_person: string
  icp_number: string
  police_number: string
  copyright: string
  site_status: number
  site_close_tips: string
  analytics_code: string
  custom_css: string
  footer_html: string
}

const loading = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive<GeneralConfig>({
  site_name: '',
  site_description: '',
  site_logo: '',
  site_favicon: '',
  contact_phone: '',
  contact_email: '',
  contact_qq: '',
  contact_wechat: '',
  contact_address: '',
  company_name: '',
  company_short_name: '',
  company_code: '',
  company_address: '',
  company_legal_person: '',
  icp_number: '',
  police_number: '',
  copyright: '',
  site_status: 1,
  site_close_tips: '',
  analytics_code: '',
  custom_css: '',
  footer_html: ''
})

const formRules: FormRules = {
  site_name: [
    { required: true, message: '请输入网站名称', trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

// 获取配置
const fetchConfig = async () => {
  loading.value = true
  try {
    const data = await request.get<GeneralConfig>({
      url: '/api/admin/config/general'
    })
    if (data) {
      Object.assign(formData, data)
    }
  } catch (error) {
    console.error('获取常规设置失败:', error)
    ElMessage.error('获取常规设置失败')
  } finally {
    loading.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.put({
        url: '/api/admin/config/general',
        params: { ...formData },
        showSuccessMessage: true
      })
      ElMessage.success('保存成功')
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 重置
const handleReset = () => {
  fetchConfig()
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.general-settings-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: var(--el-text-color-primary);
}
</style>
