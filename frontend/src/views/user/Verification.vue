<template>
  <div class="verification-page">
    <div class="page-header">
      <h1 class="page-title">实名认证</h1>
    </div>

    <el-card shadow="never" class="status-card">
      <div class="status-content">
        <el-icon :size="48" :color="verified ? '#52c41a' : '#fa8c16'">
          <component :is="verified ? CircleCheck : Warning" />
        </el-icon>
        <div class="status-info">
          <h3>{{ verified ? '已实名认证' : '未实名认证' }}</h3>
          <p>{{ verified ? '您的身份已通过验证，可正常使用所有服务' : '完成实名认证后可解锁更多功能' }}</p>
        </div>
        <el-tag v-if="verified" type="success" size="large" effect="light" round>已认证</el-tag>
        <el-tag v-else type="warning" size="large" effect="light" round>待认证</el-tag>
      </div>
    </el-card>

    <template v-if="!verified">
      <el-card shadow="never" class="form-card">
        <template #header>
          <span class="card-title">身份信息</span>
        </template>
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
          label-position="left"
          class="verify-form"
        >
          <el-form-item label="真实姓名" prop="realName">
            <el-input v-model="form.realName" placeholder="请输入身份证上的姓名" />
          </el-form-item>
          <el-form-item label="证件类型" prop="idType">
            <el-select v-model="form.idType" placeholder="请选择证件类型" style="width: 100%;">
              <el-option label="身份证" value="idcard" />
              <el-option label="护照" value="passport" />
              <el-option label="港澳通行证" value="hk_pass" />
            </el-select>
          </el-form-item>
          <el-form-item label="证件号码" prop="idNumber">
            <el-input v-model="form.idNumber" placeholder="请输入证件号码" />
          </el-form-item>

          <el-divider>证件照片上传</el-divider>

          <div class="upload-area">
            <el-form-item label="证件正面" prop="idFront">
              <el-upload
                class="id-upload"
                drag
                :auto-upload="false"
                :limit="1"
                accept="image/*"
                @change="(file: any) => handleUpload('front', file)"
              >
                <el-icon :size="40"><Upload /></el-icon>
                <div class="upload-text">点击或拖拽上传证件正面照</div>
                <template #tip>
                  <div class="upload-tip">支持 JPG/PNG，大小不超过 5MB</div>
                </template>
              </el-upload>
            </el-form-item>

            <el-form-item label="证件反面" prop="idBack">
              <el-upload
                class="id-upload"
                drag
                :auto-upload="false"
                :limit="1"
                accept="image/*"
                @change="(file: any) => handleUpload('back', file)"
              >
                <el-icon :size="40"><Upload /></el-icon>
                <div class="upload-text">点击或拖拽上传证件反面照</div>
                <template #tip>
                  <div class="upload-tip">支持 JPG/PNG，大小不超过 5MB</div>
                </template>
              </el-upload>
            </el-form-item>
          </div>

          <el-form-item>
            <el-button type="primary" size="large" :loading="submitting" @click="handleSubmit">
              提交认证
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card shadow="never" class="tips-card">
        <template #header>
          <span class="card-title">认证须知</span>
        </template>
        <ul class="tips-list">
          <li>实名认证信息仅用于身份验证，我们将严格保护您的隐私</li>
          <li>请确保上传的证件照片清晰可辨，信息完整</li>
          <li>认证审核通常在 1-3 个工作日内完成</li>
          <li>认证通过后不可修改，如需变更请联系客服</li>
          <li>每个账号只能绑定一个身份信息</li>
        </ul>
      </el-card>
    </template>

    <template v-else>
      <el-card shadow="never" class="info-card">
        <template #header>
          <span class="card-title">认证信息</span>
        </template>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="真实姓名">张 * </el-descriptions-item>
          <el-descriptions-item label="证件类型">身份证</el-descriptions-item>
          <el-descriptions-item label="证件号码">110***********1234</el-descriptions-item>
          <el-descriptions-item label="认证时间">2026-01-15 14:30:00</el-descriptions-item>
          <el-descriptions-item label="认证状态">
            <el-tag type="success" size="small" effect="light" round>已通过</el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { CircleCheck, Warning, Upload } from '@element-plus/icons-vue'
import request from '@/utils/request'

const verified = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  realName: '',
  idType: 'idcard',
  idNumber: '',
  idFront: null as File | null,
  idBack: null as File | null
})

const rules: FormRules = {
  realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  idType: [{ required: true, message: '请选择证件类型', trigger: 'change' }],
  idNumber: [{ required: true, message: '请输入证件号码', trigger: 'blur' }]
}

function handleUpload(side: string, file: any) {
  if (side === 'front') form.idFront = file.raw
  else form.idBack = file.raw
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    const formData = new FormData()
    formData.append('realName', form.realName)
    formData.append('idType', form.idType)
    formData.append('idNumber', form.idNumber)
    if (form.idFront) formData.append('idFront', form.idFront)
    if (form.idBack) formData.append('idBack', form.idBack)
    await request.post('/api/v2/certification/submit', formData)
    ElMessage.success('认证信息已提交，请等待审核')
  } catch (e: any) { ElMessage.error(e?.message || '提交失败，请重试') } finally { submitting.value = false }
}
</script>

<style scoped>
.verification-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.status-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.status-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.status-info h3 { font-size: 18px; font-weight: 600; color: #303133; margin: 0 0 4px 0; }
.status-info p { font-size: 14px; color: #909399; margin: 0; }

.form-card, .info-card, .tips-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.card-title { font-size: 15px; font-weight: 600; color: #303133; }
.verify-form { max-width: 560px; }

.upload-area {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
}

.id-upload :deep(.el-upload-dragger) {
  width: 240px;
  height: 160px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.upload-text { font-size: 13px; color: #606266; }
.upload-tip { font-size: 12px; color: #c0c4cc; }

.tips-list {
  padding-left: 18px;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tips-list li {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .status-content { flex-direction: column; text-align: center; }
  .upload-area { flex-direction: column; }
  .id-upload :deep(.el-upload-dragger) { width: 100%; }
}
</style>
