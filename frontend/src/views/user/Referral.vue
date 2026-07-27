<template>
  <div class="referral-page">
    <div class="page-header">
      <h1 class="page-title">推介计划</h1>
    </div>

    <el-card shadow="never" class="banner-card">
      <div class="banner-content">
        <div class="banner-text">
          <h2>邀请好友，赚取返利</h2>
          <p>每成功邀请一位好友注册并消费，您将获得 <strong>10%</strong> 的消费返利</p>
        </div>
        <div class="banner-stats">
          <div class="banner-stat">
            <span class="stat-num">{{ totalReferrals }}</span>
            <span class="stat-label">已邀请人数</span>
          </div>
          <div class="banner-stat">
            <span class="stat-num">¥{{ totalEarnings }}</span>
            <span class="stat-label">累计返利</span>
          </div>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="code-card">
      <template #header>
        <span class="card-title">我的推荐码</span>
      </template>
      <div class="code-area">
        <el-input v-model="referralCode" readonly size="large" class="code-input">
          <template #append>
            <el-button @click="handleCopy">
              <el-icon><CopyDocument /></el-icon>复制
            </el-button>
          </template>
        </el-input>
        <p class="code-hint">将推荐码分享给好友，好友注册时填写即可</p>
      </div>
      <div class="share-links">
        <span class="share-label">分享链接：</span>
        <el-input v-model="shareLink" readonly size="small" style="flex: 1;">
          <template #append>
            <el-button @click="handleCopyLink">复制链接</el-button>
          </template>
        </el-input>
      </div>
    </el-card>

    <el-card shadow="never" class="records-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">邀请记录</span>
        </div>
      </template>
      <el-table :data="referralRecords" stripe style="width: 100%" empty-text="暂无邀请记录">
        <el-table-column prop="username" label="用户" min-width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="28" class="record-avatar">{{ row.username.charAt(0) }}</el-avatar>
              <span>{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="registerDate" label="注册时间" width="130" />
        <el-table-column prop="spent" label="消费金额" width="120">
          <template #default="{ row }">
            <span class="amount-text">¥{{ row.spent }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="commission" label="返利金额" width="120">
          <template #default="{ row }">
            <span class="commission-text">¥{{ row.commission }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'settled' ? 'success' : 'warning'" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card shadow="never" class="rules-card">
      <template #header>
        <span class="card-title">返利规则</span>
      </template>
      <el-timeline>
        <el-timeline-item v-for="(rule, i) in rules" :key="i" :type="i === 0 ? 'primary' : ''">
          {{ rule }}
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'

const referralCode = ref('AF2026ABC')
const shareLink = ref('https://anchorfin.com/register?ref=AF2026ABC')
const totalReferrals = ref(12)
const totalEarnings = ref('1,280.00')

const referralRecords = ref([
  { username: 'user_wang', registerDate: '2026-07-20', spent: '2,397.00', commission: '239.70', status: 'settled', statusText: '已结算' },
  { username: 'user_li', registerDate: '2026-07-15', spent: '899.00', commission: '89.90', status: 'settled', statusText: '已结算' },
  { username: 'user_zhao', registerDate: '2026-07-10', spent: '199.00', commission: '19.90', status: 'pending', statusText: '待结算' },
  { username: 'user_qian', registerDate: '2026-06-28', spent: '699.00', commission: '69.90', status: 'settled', statusText: '已结算' }
])

const rules = [
  '好友通过您的推荐码或链接注册即绑定推荐关系',
  '好友每次消费后，您将获得消费金额 10% 的返利',
  '返利金额将在好友付款后 7 个工作日内结算到您的账户余额',
  '返利金额可用于平台内所有产品消费，不可提现',
  '推荐关系一旦绑定，永久有效'
]

function handleCopy() {
  navigator.clipboard?.writeText(referralCode.value)
  ElMessage.success('推荐码已复制')
}

function handleCopyLink() {
  navigator.clipboard?.writeText(shareLink.value)
  ElMessage.success('分享链接已复制')
}
</script>

<style scoped>
.referral-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 900px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.banner-card {
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, #1a3a5c 0%, #0056FF 100%);
}

.banner-card :deep(.el-card__body) { padding: 32px; }

.banner-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.banner-text h2 { font-size: 22px; font-weight: 700; color: #fff; margin: 0 0 8px 0; }
.banner-text p { font-size: 14px; color: rgba(255,255,255,0.75); margin: 0; }
.banner-text strong { color: #ffd54f; }

.banner-stats { display: flex; gap: 32px; }

.banner-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-num { font-size: 28px; font-weight: 700; color: #fff; }
.stat-label { font-size: 13px; color: rgba(255,255,255,0.65); }

.code-card, .records-card, .rules-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.card-title { font-size: 15px; font-weight: 600; color: #303133; }
.card-header { display: flex; align-items: center; justify-content: space-between; }

.code-area { display: flex; flex-direction: column; gap: 12px; }

.code-input :deep(.el-input__inner) {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 18px;
  font-weight: 600;
  color: #0056FF;
  letter-spacing: 2px;
}

.code-hint { font-size: 13px; color: #909399; margin: 0; }

.share-links {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e8ecf1;
}

.share-label { font-size: 14px; color: #606266; white-space: nowrap; }

.user-cell { display: flex; align-items: center; gap: 8px; }

.record-avatar {
  background: linear-gradient(135deg, #0056FF, #4080FF);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
}

.amount-text { font-weight: 600; color: #303133; }
.commission-text { font-weight: 600; color: #fa8c16; }

.records-card :deep(.el-table th.el-table__cell) { background: #fafbfc; font-weight: 600; }
.rules-card :deep(.el-timeline) { padding-left: 0; }

@media (max-width: 768px) {
  .banner-content { flex-direction: column; gap: 20px; text-align: center; }
}
</style>
