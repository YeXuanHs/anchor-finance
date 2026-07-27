<template>
  <div class="verification-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">实名认证</h1>
        <p class="page-desc">完成实名认证，享受更多服务和更高账户安全等级</p>
      </div>
      <div class="header-illustration">
        <svg viewBox="0 0 200 120" fill="none" width="200" height="120">
          <rect x="40" y="20" width="120" height="80" rx="8" fill="#fff" fill-opacity="0.2" />
          <rect x="55" y="35" width="40" height="30" rx="4" fill="#fff" fill-opacity="0.3" />
          <circle cx="75" cy="50" r="8" fill="#fff" fill-opacity="0.4" />
          <rect x="55" y="72" width="90" height="4" rx="2" fill="#fff" fill-opacity="0.3" />
          <rect x="55" y="82" width="60" height="4" rx="2" fill="#fff" fill-opacity="0.3" />
          <path d="M140 55l-20 0m0 0l5-5m-5 5l5 5" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill-opacity="0.6" />
          <circle cx="155" cy="75" r="18" fill="#fff" fill-opacity="0.15" />
          <path d="M148 75l4 4 8-8" stroke="#fff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill-opacity="0.7" />
        </svg>
      </div>
    </div>

    <!-- Status Alert -->
    <n-alert
      v-if="verificationStatus !== 'unverified'"
      :type="statusAlertType"
      :title="statusAlertTitle"
      :bordered="false"
      closable
    >
      {{ statusAlertDesc }}
    </n-alert>

    <!-- Verification Steps -->
    <n-card class="steps-card" :bordered="false">
      <n-steps :current="currentStep" :status="stepStatus">
        <n-step title="选择认证类型" description="个人认证或企业认证" />
        <n-step title="填写信息" description="提交认证资料" />
        <n-step title="审核中" description="等待平台审核" />
        <n-step title="认证完成" description="认证成功" />
      </n-steps>
    </n-card>

    <!-- Certification Type Selection -->
    <n-card v-if="currentStep === 1" class="form-card" title="选择认证类型" :bordered="false">
      <div class="type-selection">
        <div
          class="type-card"
          :class="{ active: certType === 'personal' }"
          @click="certType = 'personal'"
        >
          <div class="type-icon personal">
            <n-icon :size="40" :component="PersonOutline" />
          </div>
          <div class="type-info">
            <div class="type-title">个人认证</div>
            <div class="type-desc">适用于个人用户，需提供身份证信息</div>
          </div>
          <n-radio
            :checked="certType === 'personal'"
            @update:checked="certType = 'personal'"
          />
        </div>
        <div
          class="type-card"
          :class="{ active: certType === 'enterprise' }"
          @click="certType = 'enterprise'"
        >
          <div class="type-icon enterprise">
            <n-icon :size="40" :component="BusinessOutline" />
          </div>
          <div class="type-info">
            <div class="type-title">企业认证</div>
            <div class="type-desc">适用于企业用户，需提供营业执照信息</div>
          </div>
          <n-radio
            :checked="certType === 'enterprise'"
            @update:checked="certType = 'enterprise'"
          />
        </div>
      </div>
      <div class="form-actions">
        <n-button type="primary" round @click="goToStep2">
          下一步
        </n-button>
      </div>
    </n-card>

    <!-- Personal Certification Form -->
    <n-card v-if="currentStep === 2 && certType === 'personal'" class="form-card" title="个人实名认证" :bordered="false">
      <n-form
        ref="personalFormRef"
        :model="personalForm"
        :rules="personalRules"
        label-placement="left"
        label-width="100"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="真实姓名" path="realName">
          <n-input v-model:value="personalForm.realName" placeholder="请输入身份证上的真实姓名" />
        </n-form-item>
        <n-form-item label="身份证号" path="idCard">
          <n-input v-model:value="personalForm.idCard" placeholder="请输入18位身份证号码" maxlength="18" />
        </n-form-item>
        <n-form-item label="身份证正面" path="idCardFront">
          <n-upload
            :max="1"
            accept="image/*"
            :default-upload="false"
            @change="handleFrontChange"
          >
            <n-button>
              <template #icon>
                <n-icon :component="CloudUploadOutline" />
              </template>
              上传人像面照片
            </n-button>
          </n-upload>
          <div class="upload-tip">请上传身份证人像面照片，支持 JPG/PNG 格式</div>
        </n-form-item>
        <n-form-item label="身份证反面" path="idCardBack">
          <n-upload
            :max="1"
            accept="image/*"
            :default-upload="false"
            @change="handleBackChange"
          >
            <n-button>
              <template #icon>
                <n-icon :component="CloudUploadOutline" />
              </template>
              上传国徽面照片
            </n-button>
          </n-upload>
          <div class="upload-tip">请上传身份证国徽面照片，支持 JPG/PNG 格式</div>
        </n-form-item>
      </n-form>
      <div class="form-actions">
        <n-button round @click="currentStep = 1">上一步</n-button>
        <n-button type="primary" round :loading="submitting" @click="handleSubmitPersonal">
          提交认证
        </n-button>
      </div>
    </n-card>

    <!-- Enterprise Certification Form -->
    <n-card v-if="currentStep === 2 && certType === 'enterprise'" class="form-card" title="企业实名认证" :bordered="false">
      <n-form
        ref="enterpriseFormRef"
        :model="enterpriseForm"
        :rules="enterpriseRules"
        label-placement="left"
        label-width="120"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="企业名称" path="companyName">
          <n-input v-model:value="enterpriseForm.companyName" placeholder="请输入营业执照上的企业全称" />
        </n-form-item>
        <n-form-item label="统一社会信用代码" path="creditCode">
          <n-input v-model:value="enterpriseForm.creditCode" placeholder="请输入18位统一社会信用代码" maxlength="18" />
        </n-form-item>
        <n-form-item label="法人姓名" path="legalPerson">
          <n-input v-model:value="enterpriseForm.legalPerson" placeholder="请输入法定代表人姓名" />
        </n-form-item>
        <n-form-item label="法人身份证号" path="legalIdCard">
          <n-input v-model:value="enterpriseForm.legalIdCard" placeholder="请输入法定代表人身份证号码" maxlength="18" />
        </n-form-item>
        <n-form-item label="营业执照" path="businessLicense">
          <n-upload
            :max="1"
            accept="image/*"
            :default-upload="false"
            @change="handleLicenseChange"
          >
            <n-button>
              <template #icon>
                <n-icon :component="CloudUploadOutline" />
              </template>
              上传营业执照照片
            </n-button>
          </n-upload>
          <div class="upload-tip">请上传营业执照副本照片，支持 JPG/PNG 格式</div>
        </n-form-item>
        <n-form-item label="联系人" path="contactPerson">
          <n-input v-model:value="enterpriseForm.contactPerson" placeholder="请输入业务联系人姓名" />
        </n-form-item>
        <n-form-item label="联系电话" path="contactPhone">
          <n-input v-model:value="enterpriseForm.contactPhone" placeholder="请输入业务联系电话" />
        </n-form-item>
      </n-form>
      <div class="form-actions">
        <n-button round @click="currentStep = 1">上一步</n-button>
        <n-button type="primary" round :loading="submitting" @click="handleSubmitEnterprise">
          提交认证
        </n-button>
      </div>
    </n-card>

    <!-- Review Pending -->
    <n-card v-if="currentStep === 3" class="form-card" :bordered="false">
      <div class="review-pending">
        <div class="pending-icon">
          <n-icon :size="64" :component="TimeOutline" color="#fa8c16" />
        </div>
        <div class="pending-title">认证审核中</div>
        <div class="pending-desc">
          您的{{ certType === 'personal' ? '个人' : '企业' }}认证资料已提交，我们将在 1-3 个工作日内完成审核。
          审核结果将通过站内消息和邮件通知您。
        </div>
        <div class="pending-info">
          <div class="info-item">
            <span class="info-label">提交时间</span>
            <span class="info-value">{{ submitTime }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">认证类型</span>
            <span class="info-value">{{ certType === 'personal' ? '个人认证' : '企业认证' }}</span>
          </div>
        </div>
      </div>
    </n-card>

    <!-- Verified -->
    <n-card v-if="verificationStatus === 'verified'" class="form-card" :bordered="false">
      <div class="verified-info">
        <div class="verified-icon">
          <n-icon :size="64" :component="CheckmarkCircleOutline" color="#52c41a" />
        </div>
        <div class="verified-title">认证已通过</div>
        <div class="verified-desc">您的实名认证已通过审核，可以享受所有平台服务。</div>
        <div class="verified-details">
          <div class="detail-item">
            <span class="detail-label">认证类型</span>
            <span class="detail-value">{{ certType === 'personal' ? '个人认证' : '企业认证' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">{{ certType === 'personal' ? '真实姓名' : '企业名称' }}</span>
            <span class="detail-value">{{ certType === 'personal' ? personalForm.realName : enterpriseForm.companyName }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">认证时间</span>
            <span class="detail-value">2026-07-20</span>
          </div>
        </div>
      </div>
    </n-card>

    <!-- Verification History -->
    <n-card class="history-card" title="认证记录" :bordered="false">
      <n-data-table
        :columns="historyColumns"
        :data="historyRecords"
        :bordered="false"
        :single-line="false"
        size="small"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns, FormRules } from 'naive-ui'
import { NTag, NButton, useMessage } from 'naive-ui'
import {
  PersonOutline,
  BusinessOutline,
  CloudUploadOutline,
  TimeOutline,
  CheckmarkCircleOutline
} from '@vicons/ionicons5'

const message = useMessage()

const certType = ref<'personal' | 'enterprise'>('personal')
const currentStep = ref(1)
const stepStatus = ref<'process' | 'finish' | 'error' | 'wait'>('process')
const submitting = ref(false)
const verificationStatus = ref<'unverified' | 'pending' | 'verified' | 'failed'>('unverified')
const submitTime = ref('')

// Personal form
const personalFormRef = ref()
const personalForm = ref({
  realName: '',
  idCard: '',
  idCardFront: null as File | null,
  idCardBack: null as File | null
})

const personalRules: FormRules = {
  realName: { required: true, message: '请输入真实姓名', trigger: 'blur' },
  idCard: [
    { required: true, message: '请输入身份证号码', trigger: 'blur' },
    { pattern: /^\d{17}[\dXx]$/, message: '请输入有效的身份证号码', trigger: 'blur' }
  ]
}

// Enterprise form
const enterpriseFormRef = ref()
const enterpriseForm = ref({
  companyName: '',
  creditCode: '',
  legalPerson: '',
  legalIdCard: '',
  businessLicense: null as File | null,
  contactPerson: '',
  contactPhone: ''
})

const enterpriseRules: FormRules = {
  companyName: { required: true, message: '请输入企业名称', trigger: 'blur' },
  creditCode: [
    { required: true, message: '请输入统一社会信用代码', trigger: 'blur' },
    { pattern: /^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$/, message: '请输入有效的统一社会信用代码', trigger: 'blur' }
  ],
  legalPerson: { required: true, message: '请输入法定代表人姓名', trigger: 'blur' },
  legalIdCard: [
    { required: true, message: '请输入法定代表人身份证号码', trigger: 'blur' },
    { pattern: /^\d{17}[\dXx]$/, message: '请输入有效的身份证号码', trigger: 'blur' }
  ],
  contactPerson: { required: true, message: '请输入联系人姓名', trigger: 'blur' },
  contactPhone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码', trigger: 'blur' }
  ]
}

// Status alert
const statusAlertType = computed(() => {
  const map = {
    pending: 'warning' as const,
    verified: 'success' as const,
    failed: 'error' as const,
    unverified: 'info' as const
  }
  return map[verificationStatus.value]
})

const statusAlertTitle = computed(() => {
  const map = {
    pending: '认证审核中',
    verified: '认证已通过',
    failed: '认证未通过',
    unverified: ''
  }
  return map[verificationStatus.value]
})

const statusAlertDesc = computed(() => {
  const map = {
    pending: '您的认证资料正在审核中，请耐心等待。',
    verified: '您的实名认证已通过，可以正常使用所有功能。',
    failed: '您的认证未通过，请检查资料后重新提交。',
    unverified: ''
  }
  return map[verificationStatus.value]
})

// History records
interface HistoryRecord {
  id: string
  type: string
  typeName: string
  submitTime: string
  reviewTime: string
  status: 'approved' | 'pending' | 'rejected'
  statusText: string
  reason?: string
}

const historyRecords = ref<HistoryRecord[]>([
  {
    id: 'VER001',
    type: 'personal',
    typeName: '个人认证',
    submitTime: '2026-07-25 14:30',
    reviewTime: '2026-07-26 10:15',
    status: 'pending',
    statusText: '审核中'
  },
  {
    id: 'VER002',
    type: 'personal',
    typeName: '个人认证',
    submitTime: '2026-07-10 09:00',
    reviewTime: '2026-07-11 16:20',
    status: 'rejected',
    statusText: '未通过',
    reason: '身份证照片模糊，请重新上传'
  },
  {
    id: 'VER003',
    type: 'personal',
    typeName: '个人认证',
    submitTime: '2026-06-15 11:00',
    reviewTime: '2026-06-16 14:30',
    status: 'approved',
    statusText: '已通过'
  }
])

const historyColumns: DataTableColumns<HistoryRecord> = [
  { title: '编号', key: 'id', width: 100 },
  { title: '认证类型', key: 'typeName', width: 100 },
  { title: '提交时间', key: 'submitTime', width: 160 },
  { title: '审核时间', key: 'reviewTime', width: 160 },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const typeMap: Record<string, 'success' | 'warning' | 'error'> = {
        approved: 'success',
        pending: 'warning',
        rejected: 'error'
      }
      return h(
        NTag,
        { type: typeMap[row.status], size: 'small', round: true, bordered: false },
        { default: () => row.statusText }
      )
    }
  },
  {
    title: '备注',
    key: 'reason',
    ellipsis: { tooltip: true },
    render: (row) => h('span', { style: row.status === 'rejected' ? 'color: #f5222d' : '' }, row.reason || '-')
  }
]

// File handlers
function handleFrontChange(options: { fileList: any[] }) {
  if (options.fileList.length > 0) {
    personalForm.value.idCardFront = options.fileList[0].file
  }
}

function handleBackChange(options: { fileList: any[] }) {
  if (options.fileList.length > 0) {
    personalForm.value.idCardBack = options.fileList[0].file
  }
}

function handleLicenseChange(options: { fileList: any[] }) {
  if (options.fileList.length > 0) {
    enterpriseForm.value.businessLicense = options.fileList[0].file
  }
}

// Step navigation
function goToStep2() {
  currentStep.value = 2
}

// Submit handlers
function handleSubmitPersonal() {
  personalFormRef.value?.validate((errors: any) => {
    if (errors) {
      message.warning('请完善必填信息')
      return
    }
    submitting.value = true
    setTimeout(() => {
      submitting.value = false
      verificationStatus.value = 'pending'
      currentStep.value = 3
      submitTime.value = new Date().toLocaleString('zh-CN')
      message.success('个人认证资料已提交')
    }, 1500)
  })
}

function handleSubmitEnterprise() {
  enterpriseFormRef.value?.validate((errors: any) => {
    if (errors) {
      message.warning('请完善必填信息')
      return
    }
    submitting.value = true
    setTimeout(() => {
      submitting.value = false
      verificationStatus.value = 'pending'
      currentStep.value = 3
      submitTime.value = new Date().toLocaleString('zh-CN')
      message.success('企业认证资料已提交')
    }, 1500)
  })
}
</script>

<style scoped>
.verification-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ==================== Page Header ==================== */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 50%, #0050b3 100%);
  border-radius: 12px;
  padding: 32px;
  position: relative;
  overflow: hidden;
}

.page-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.header-content {
  position: relative;
  z-index: 1;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
}

.page-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
}

.header-illustration {
  position: relative;
  z-index: 1;
}

/* ==================== Steps Card ==================== */
.steps-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

/* ==================== Form Card ==================== */
.form-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

/* ==================== Type Selection ==================== */
.type-selection {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.type-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 32px 24px;
  border-radius: 12px;
  border: 2px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.type-card:hover {
  border-color: #d0e4ff;
  background: #f0f5ff;
}

.type-card.active {
  border-color: #1890ff;
  background: #e6f0ff;
}

.type-card.active::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: #1890ff;
  border-radius: 12px 12px 0 0;
}

.type-icon {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.type-icon.personal {
  background: linear-gradient(135deg, #e6f0ff, #bae0ff);
  color: #1890ff;
}

.type-icon.enterprise {
  background: linear-gradient(135deg, #f0f5ff, #d6e4ff);
  color: #2f54eb;
}

.type-info {
  text-align: center;
}

.type-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
}

.type-desc {
  font-size: 13px;
  color: #8c8c8c;
}

/* ==================== Form Actions ==================== */
.form-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

/* ==================== Upload Tips ==================== */
.upload-tip {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 8px;
}

/* ==================== Review Pending ==================== */
.review-pending {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 24px;
  text-align: center;
}

.pending-icon {
  margin-bottom: 24px;
}

.pending-title {
  font-size: 20px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 12px;
}

.pending-desc {
  font-size: 14px;
  color: #8c8c8c;
  max-width: 480px;
  line-height: 1.6;
  margin-bottom: 32px;
}

.pending-info {
  display: flex;
  gap: 48px;
  padding: 24px;
  background: #fafafa;
  border-radius: 12px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-label {
  font-size: 13px;
  color: #8c8c8c;
}

.info-value {
  font-size: 15px;
  font-weight: 500;
  color: #262626;
}

/* ==================== Verified Info ==================== */
.verified-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 24px;
  text-align: center;
}

.verified-icon {
  margin-bottom: 24px;
}

.verified-title {
  font-size: 20px;
  font-weight: 600;
  color: #52c41a;
  margin-bottom: 12px;
}

.verified-desc {
  font-size: 14px;
  color: #8c8c8c;
  max-width: 480px;
  line-height: 1.6;
  margin-bottom: 32px;
}

.verified-details {
  display: flex;
  gap: 48px;
  padding: 24px;
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 12px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-label {
  font-size: 13px;
  color: #8c8c8c;
}

.detail-value {
  font-size: 15px;
  font-weight: 500;
  color: #262626;
}

/* ==================== History Card ==================== */
.history-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

/* ==================== Responsive ==================== */
@media (max-width: 1200px) {
  .type-selection {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header-illustration {
    display: none;
  }

  .pending-info,
  .verified-details {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
