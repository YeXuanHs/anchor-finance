<template>
  <div class="verification-page">
    <!-- Header Banner -->
    <div class="page-banner">
      <div class="banner-content">
        <h1 class="banner-title">实名认证</h1>
        <p class="banner-desc">完成实名认证，享受完整服务</p>
      </div>
      <div class="banner-status">
        <div class="status-icon" :class="statusClass">
          <n-icon :size="32" :component="statusIcon" />
        </div>
        <span class="status-text">{{ statusText }}</span>
      </div>
    </div>

    <!-- Status Alert -->
    <n-alert
      v-if="verificationStatus !== 'none'"
      :type="alertType"
      :title="alertTitle"
      closable
    >
      {{ alertMessage }}
    </n-alert>

    <!-- Steps -->
    <n-card class="section-card">
      <n-steps :current="currentStep" :status="stepStatus">
        <n-step title="选择认证类型" />
        <n-step title="填写认证信息" />
        <n-step title="上传证件照片" />
        <n-step title="提交审核" />
      </n-steps>
    </n-card>

    <!-- Verification Type Selection -->
    <n-card v-if="currentStep === 1" class="section-card" title="选择认证类型">
      <div class="type-selection">
        <div
          class="type-card"
          :class="{ active: verificationType === 'personal' }"
          @click="verificationType = 'personal'"
        >
          <div class="type-icon personal">
            <n-icon :size="40" :component="PersonOutline" />
          </div>
          <h3>个人认证</h3>
          <p>适用于个人用户，需提供身份证信息</p>
        </div>
        <div
          class="type-card"
          :class="{ active: verificationType === 'enterprise' }"
          @click="verificationType = 'enterprise'"
        >
          <div class="type-icon enterprise">
            <n-icon :size="40" :component="BusinessOutline" />
          </div>
          <h3>企业认证</h3>
          <p>适用于企业用户，需提供营业执照</p>
        </div>
      </div>
      <div class="step-actions">
        <n-button type="primary" round :disabled="!verificationType" @click="nextStep">
          下一步
        </n-button>
      </div>
    </n-card>

    <!-- Personal Verification Form -->
    <n-card v-if="currentStep === 2 && verificationType === 'personal'" class="section-card" title="个人信息">
      <n-form
        ref="personalFormRef"
        :model="personalForm"
        :rules="personalRules"
        label-placement="left"
        label-width="100"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="真实姓名" path="realName">
          <n-input v-model:value="personalForm.realName" placeholder="请输入真实姓名" />
        </n-form-item>
        <n-form-item label="身份证号" path="idCard">
          <n-input v-model:value="personalForm.idCard" placeholder="请输入18位身份证号码" maxlength="18" />
        </n-form-item>
      </n-form>
      <div class="step-actions">
        <n-button round @click="prevStep">上一步</n-button>
        <n-button type="primary" round @click="nextStep">下一步</n-button>
      </div>
    </n-card>

    <!-- Enterprise Verification Form -->
    <n-card v-if="currentStep === 2 && verificationType === 'enterprise'" class="section-card" title="企业信息">
      <n-form
        ref="enterpriseFormRef"
        :model="enterpriseForm"
        :rules="enterpriseRules"
        label-placement="left"
        label-width="120"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="企业名称" path="companyName">
          <n-input v-model:value="enterpriseForm.companyName" placeholder="请输入企业全称" />
        </n-form-item>
        <n-form-item label="统一社会信用代码" path="creditCode">
          <n-input v-model:value="enterpriseForm.creditCode" placeholder="请输入18位统一社会信用代码" maxlength="18" />
        </n-form-item>
        <n-form-item label="法人姓名" path="legalPerson">
          <n-input v-model:value="enterpriseForm.legalPerson" placeholder="请输入法人姓名" />
        </n-form-item>
      </n-form>
      <div class="step-actions">
        <n-button round @click="prevStep">上一步</n-button>
        <n-button type="primary" round @click="nextStep">下一步</n-button>
      </div>
    </n-card>

    <!-- Upload Documents -->
    <n-card v-if="currentStep === 3 && verificationType === 'personal'" class="section-card" title="上传身份证照片">
      <div class="upload-grid">
        <div class="upload-item">
          <h4>身份证正面（人像面）</h4>
          <n-upload
            :max="1"
            accept="image/*"
            list-type="image-card"
            @change="(options: any) => handleUpload('idFront', options)"
          >
            <div class="upload-placeholder">
              <n-icon :size="32" :component="CameraOutline" />
              <span>点击上传</span>
            </div>
          </n-upload>
        </div>
        <div class="upload-item">
          <h4>身份证反面（国徽面）</h4>
          <n-upload
            :max="1"
            accept="image/*"
            list-type="image-card"
            @change="(options: any) => handleUpload('idBack', options)"
          >
            <div class="upload-placeholder">
              <n-icon :size="32" :component="CameraOutline" />
              <span>点击上传</span>
            </div>
          </n-upload>
        </div>
      </div>
      <div class="upload-tips">
        <p>拍摄要求：</p>
        <ul>
          <li>请确保照片清晰，文字可辨认</li>
          <li>请勿遮挡或反光</li>
          <li>支持 JPG、PNG 格式，大小不超过 5MB</li>
        </ul>
      </div>
      <div class="step-actions">
        <n-button round @click="prevStep">上一步</n-button>
        <n-button type="primary" round :loading="submitting" @click="handleSubmit">提交认证</n-button>
      </div>
    </n-card>

    <n-card v-if="currentStep === 3 && verificationType === 'enterprise'" class="section-card" title="上传营业执照">
      <div class="upload-grid single">
        <div class="upload-item">
          <h4>营业执照照片</h4>
          <n-upload
            :max="1"
            accept="image/*"
            list-type="image-card"
            @change="(options: any) => handleUpload('license', options)"
          >
            <div class="upload-placeholder">
              <n-icon :size="32" :component="CameraOutline" />
              <span>点击上传</span>
            </div>
          </n-upload>
        </div>
      </div>
      <div class="upload-tips">
        <p>拍摄要求：</p>
        <ul>
          <li>请上传最新的营业执照副本</li>
          <li>确保照片清晰，信息完整可见</li>
          <li>支持 JPG、PNG 格式，大小不超过 5MB</li>
        </ul>
      </div>
      <div class="step-actions">
        <n-button round @click="prevStep">上一步</n-button>
        <n-button type="primary" round :loading="submitting" @click="handleSubmit">提交认证</n-button>
      </div>
    </n-card>

    <!-- Success -->
    <n-card v-if="currentStep === 4" class="section-card success-card">
      <div class="success-content">
        <div class="success-icon">
          <n-icon :size="64" :component="CheckmarkCircleOutline" color="#52c41a" />
        </div>
        <h2>认证申请已提交</h2>
        <p>我们将在 1-3 个工作日内完成审核，请耐心等待</p>
        <n-button type="primary" round @click="$router.push('/user/dashboard')">返回首页</n-button>
      </div>
    </n-card>

    <!-- Verification History -->
    <n-card v-if="verificationHistory.length" class="section-card" title="认证记录">
      <n-data-table
        :columns="historyColumns"
        :data="verificationHistory"
        :bordered="false"
        :single-line="false"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns, FormRules } from 'naive-ui'
import { NTag } from 'naive-ui'
import {
  PersonOutline,
  BusinessOutline,
  CameraOutline,
  CheckmarkCircleOutline,
  ShieldCheckmarkOutline,
  TimeOutline,
  CloseCircleOutline,
  EllipseOutline
} from '@vicons/ionicons5'
import { useMessage } from 'naive-ui'

const message = useMessage()

type VerificationStatus = 'none' | 'pending' | 'verified' | 'failed'
type VerificationType = '' | 'personal' | 'enterprise'

const verificationStatus = ref<VerificationStatus>('none')
const verificationType = ref<VerificationType>('')
const currentStep = ref(1)
const stepStatus = ref<'process' | 'finish' | 'error'>('process')
const submitting = ref(false)

const personalForm = ref({
  realName: '',
  idCard: ''
})

const enterpriseForm = ref({
  companyName: '',
  creditCode: '',
  legalPerson: ''
})

const personalRules: FormRules = {
  realName: { required: true, message: '请输入真实姓名', trigger: 'blur' },
  idCard: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { pattern: /^\d{17}[\dXx]$/, message: '请输入正确的18位身份证号码', trigger: 'blur' }
  ]
}

const enterpriseRules: FormRules = {
  companyName: { required: true, message: '请输入企业名称', trigger: 'blur' },
  creditCode: [
    { required: true, message: '请输入统一社会信用代码', trigger: 'blur' },
    { pattern: /^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$/, message: '请输入正确的统一社会信用代码', trigger: 'blur' }
  ],
  legalPerson: { required: true, message: '请输入法人姓名', trigger: 'blur' }
}

const statusClass = computed(() => {
  const map: Record<VerificationStatus, string> = {
    none: 'unverified',
    pending: 'pending',
    verified: 'verified',
    failed: 'failed'
  }
  return map[verificationStatus.value]
})

const statusIcon = computed(() => {
  const map: Record<VerificationStatus, any> = {
    none: EllipseOutline,
    pending: TimeOutline,
    verified: ShieldCheckmarkOutline,
    failed: CloseCircleOutline
  }
  return map[verificationStatus.value]
})

const statusText = computed(() => {
  const map: Record<VerificationStatus, string> = {
    none: '未认证',
    pending: '审核中',
    verified: '已认证',
    failed: '认证失败'
  }
  return map[verificationStatus.value]
})

const alertType = computed(() => {
  const map = { pending: 'info', verified: 'success', failed: 'error', none: 'info' }
  return map[verificationStatus.value] as 'info' | 'success' | 'error'
})

const alertTitle = computed(() => {
  const map = {
    pending: '认证审核中',
    verified: '认证已通过',
    failed: '认证未通过',
    none: ''
  }
  return map[verificationStatus.value]
})

const alertMessage = computed(() => {
  const map = {
    pending: '您的认证申请正在审核中，预计 1-3 个工作日内完成。',
    verified: '恭喜您，实名认证已通过！现在可以享受完整服务。',
    failed: '很抱歉，您的认证申请未通过。请检查信息是否正确后重新提交。',
    none: ''
  }
  return map[verificationStatus.value]
})

interface HistoryRecord {
  id: number
  type: string
  submitDate: string
  reviewDate: string
  status: 'approved' | 'pending' | 'rejected'
  statusText: string
  remark: string
}

const verificationHistory = ref<HistoryRecord[]>([])

const historyColumns: DataTableColumns<HistoryRecord> = [
  { title: '认证类型', key: 'type', width: 100 },
  { title: '提交时间', key: 'submitDate', width: 120 },
  { title: '审核时间', key: 'reviewDate', width: 120 },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) =>
      h(
        NTag,
        {
          type: row.status === 'approved' ? 'success' : row.status === 'pending' ? 'warning' : 'error',
          size: 'small',
          round: true,
          bordered: false
        },
        { default: () => row.statusText }
      )
  },
  { title: '备注', key: 'remark', ellipsis: { tooltip: true } }
]

function nextStep() {
  if (currentStep.value < 4) {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

function handleUpload(field: string, options: any) {
  // Handle file upload
  console.log(`Upload ${field}:`, options)
}

async function handleSubmit() {
  submitting.value = true
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500))
    verificationStatus.value = 'pending'
    currentStep.value = 4
    message.success('认证申请已提交')
  } catch {
    message.error('提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.verification-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ==================== Banner ==================== */
.page-banner {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 50%, #0050b3 100%);
  border-radius: 12px;
  padding: 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  overflow: hidden;
  position: relative;
}

.page-banner::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.banner-content {
  position: relative;
  z-index: 1;
}

.banner-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
}

.banner-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
}

.banner-status {
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
  z-index: 1;
}

.status-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-icon.unverified {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.status-icon.pending {
  background: rgba(250, 173, 20, 0.3);
  color: #fad414;
}

.status-icon.verified {
  background: rgba(82, 196, 26, 0.3);
  color: #52c41a;
}

.status-icon.failed {
  background: rgba(255, 77, 79, 0.3);
  color: #ff4d4f;
}

.status-text {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

/* ==================== Section Card ==================== */
.section-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

.section-card :deep(.n-card-header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.section-card :deep(.n-card__content) {
  padding: 20px;
}

/* ==================== Type Selection ==================== */
.type-selection {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 24px;
}

.type-card {
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  padding: 32px 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.type-card:hover {
  border-color: #91caff;
  background: #f0f7ff;
}

.type-card.active {
  border-color: #1890ff;
  background: #e6f4ff;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.15);
}

.type-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.type-icon.personal {
  background: linear-gradient(135deg, #e6f4ff, #bae0ff);
  color: #1890ff;
}

.type-icon.enterprise {
  background: linear-gradient(135deg, #f6ffed, #d9f7be);
  color: #52c41a;
}

.type-card h3 {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 8px 0;
}

.type-card p {
  font-size: 14px;
  color: #8c8c8c;
  margin: 0;
}

/* ==================== Upload Grid ==================== */
.upload-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 16px;
}

.upload-grid.single {
  grid-template-columns: 1fr;
  max-width: 400px;
}

.upload-item h4 {
  font-size: 15px;
  font-weight: 500;
  color: #262626;
  margin: 0 0 12px 0;
}

.upload-placeholder {
  width: 148px;
  height: 148px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: #fafafa;
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
  color: #8c8c8c;
  font-size: 14px;
}

.upload-tips {
  background: #f0f7ff;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 24px;
}

.upload-tips p {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
  margin: 0 0 8px 0;
}

.upload-tips ul {
  margin: 0;
  padding-left: 20px;
}

.upload-tips li {
  font-size: 13px;
  color: #595959;
  line-height: 1.8;
}

/* ==================== Step Actions ==================== */
.step-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

/* ==================== Success Card ==================== */
.success-card {
  text-align: center;
}

.success-content {
  padding: 32px 0;
}

.success-icon {
  margin-bottom: 24px;
}

.success-content h2 {
  font-size: 22px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 8px 0;
}

.success-content p {
  font-size: 14px;
  color: #8c8c8c;
  margin: 0 0 32px 0;
}

/* ==================== Responsive ==================== */
@media (max-width: 768px) {
  .page-banner {
    flex-direction: column;
    gap: 24px;
    text-align: center;
  }

  .type-selection {
    grid-template-columns: 1fr;
  }

  .upload-grid {
    grid-template-columns: 1fr;
  }

  .upload-grid.single {
    max-width: 100%;
  }
}
</style>
