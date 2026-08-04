<template>
  <div class="enterprise-verification">
    <div class="page-header">
      <h2>企业认证</h2>
      <p>完成企业认证可享受更多权益</p>
    </div>
    <el-card v-if="status === 'pending'">
      <el-result icon="info" title="认证审核中" sub-title="您的企业认证申请正在审核中，预计1-3个工作日内完成">
        <template #extra><el-button type="primary" @click="status = 'none'">重新提交</el-button></template>
      </el-result>
    </el-card>
    <el-card v-else-if="status === 'approved'">
      <el-result icon="success" title="已认证" sub-title="您的企业认证已通过">
        <template #extra>
          <div class="cert-info">
            <p>企业名称：{{ certInfo.company_name }}</p>
            <p>统一社会信用代码：{{ certInfo.credit_code }}</p>
          </div>
        </template>
      </el-result>
    </el-card>
    <el-card v-else>
      <el-form :model="form" label-width="140px" :rules="rules" ref="formRef">
        <el-form-item label="企业名称" prop="company_name">
          <el-input v-model="form.company_name" placeholder="请输入企业全称" />
        </el-form-item>
        <el-form-item label="统一社会信用代码" prop="credit_code">
          <el-input v-model="form.credit_code" placeholder="18位统一社会信用代码" />
        </el-form-item>
        <el-form-item label="法人姓名" prop="legal_person">
          <el-input v-model="form.legal_person" placeholder="请输入法人姓名" />
        </el-form-item>
        <el-form-item label="法人身份证号" prop="legal_id_card">
          <el-input v-model="form.legal_id_card" placeholder="请输入法人身份证号" />
        </el-form-item>
        <el-form-item label="营业执照" prop="business_license">
          <el-upload action="/api/upload" :on-success="handleUpload" accept="image/*" :limit="1">
            <el-button>上传营业执照</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item label="联系人" prop="contact_name">
          <el-input v-model="form.contact_name" placeholder="请输入联系人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="contact_phone">
          <el-input v-model="form.contact_phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">提交认证</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const status = ref('none') // none/pending/approved
const submitting = ref(false)
const certInfo = ref<any>({})
const formRef = ref()
const form = ref({ company_name: '', credit_code: '', legal_person: '', legal_id_card: '', business_license: '', contact_name: '', contact_phone: '' })
const rules = {
  company_name: [{ required: true, message: '请输入企业名称', trigger: 'blur' }],
  credit_code: [{ required: true, message: '请输入统一社会信用代码', trigger: 'blur' }],
  legal_person: [{ required: true, message: '请输入法人姓名', trigger: 'blur' }]
}

const handleUpload = (res: any) => { form.value.business_license = res.data?.url || '' }

const handleSubmit = async () => {
  submitting.value = true
  try {
    await request.post('/api/v2/user/certification/enterprise', form.value)
    ElMessage.success('提交成功')
    status.value = 'pending'
  } catch { ElMessage.error('提交失败') } finally { submitting.value = false }
}

const fetchStatus = async () => {
  try {
    const { data } = await request.get('/api/v2/user/certification/enterprise')
    if (data.data) { status.value = data.data.status; certInfo.value = data.data }
  } catch {}
}

onMounted(fetchStatus)
</script>
<style scoped lang="scss">
.enterprise-verification { .page-header { margin-bottom: 24px; h2 { font-size: 20px; color: #1a365d; } p { color: #6b7280; margin-top: 4px; } } }
.cert-info { text-align: left; p { color: #4b5563; line-height: 2; } }
</style>
