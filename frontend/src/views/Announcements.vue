<template>
  <div class="announcements-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <n-icon :component="MegaphoneOutline" size="28" />
          公告中心
        </h1>
        <p class="page-subtitle">了解最新动态与重要通知</p>
      </div>
      <div class="header-actions">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索公告..."
          clearable
          style="width: 280px"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" color="#1890ff" />
          </template>
        </n-input>
      </div>
    </div>

    <div v-if="selectedAnnouncement" class="detail-section">
      <n-card class="detail-card" :bordered="false">
        <div class="detail-back" @click="selectedAnnouncement = null">
          <n-icon :component="ArrowBackOutline" size="18" />
          <span>返回列表</span>
        </div>
        <div class="detail-header">
          <h2 class="detail-title">{{ selectedAnnouncement.title }}</h2>
          <div class="detail-meta">
            <n-tag :type="getTagType(selectedAnnouncement.level)" size="small" :bordered="false">
              {{ selectedAnnouncement.levelLabel }}
            </n-tag>
            <span class="detail-time">
              <n-icon :component="TimeOutline" size="14" />
              {{ selectedAnnouncement.time }}
            </span>
          </div>
        </div>
        <n-divider />
        <div class="detail-content" v-html="selectedAnnouncement.content"></div>
      </n-card>
    </div>

    <div v-else class="list-section">
      <div v-if="filteredList.length > 0" class="announcement-list">
        <n-card
          v-for="item in paginatedList"
          :key="item.id"
          class="announcement-card"
          :bordered="false"
          hoverable
          @click="selectedAnnouncement = item"
        >
          <div class="card-header">
            <div class="card-title-row">
              <n-tag :type="getTagType(item.level)" size="small" :bordered="false">
                {{ item.levelLabel }}
              </n-tag>
              <h3 class="card-title">{{ item.title }}</h3>
            </div>
            <span class="card-time">
              <n-icon :component="TimeOutline" size="14" />
              {{ item.time }}
            </span>
          </div>
          <p class="card-summary">{{ item.summary }}</p>
          <div class="card-footer">
            <div class="card-tags">
              <n-tag
                v-for="tag in item.tags"
                :key="tag"
                size="tiny"
                :bordered="false"
                class="meta-tag"
              >
                {{ tag }}
              </n-tag>
            </div>
            <span class="read-more">
              查看详情
              <n-icon :component="ChevronForwardOutline" size="14" />
            </span>
          </div>
        </n-card>
      </div>

      <div v-else class="empty-state">
        <n-empty description="暂无公告" />
      </div>

      <div v-if="totalPages > 1" class="pagination-wrapper">
        <n-pagination
          v-model:page="currentPage"
          :page-count="totalPages"
          :page-slot="7"
          show-quick-jumper
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  MegaphoneOutline,
  SearchOutline,
  TimeOutline,
  ArrowBackOutline,
  ChevronForwardOutline
} from '@vicons/ionicons5'

const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = 6
const selectedAnnouncement = ref<Announcement | null>(null)

interface Announcement {
  id: number
  title: string
  summary: string
  content: string
  time: string
  level: string
  levelLabel: string
  tags: string[]
}

const announcements = ref<Announcement[]>([
  {
    id: 1,
    title: '系统升级维护通知',
    summary: '为提升服务质量，系统将于本周六凌晨2:00-6:00进行升级维护，届时部分功能将暂停使用。',
    content: '<p>尊敬的用户：</p><p>为提升服务质量，系统将于<strong>2026年7月29日（周六）凌晨2:00-6:00</strong>进行升级维护。</p><p>维护期间，以下功能将暂停使用：</p><ul><li>账户充值与提现</li><li>交易记录查询</li><li>部分API接口</li></ul><p>请您提前做好相关安排，对此给您带来的不便，我们深表歉意。</p><p>感谢您的理解与支持！</p>',
    time: '2026-07-27',
    level: 'important',
    levelLabel: '重要',
    tags: ['系统维护', '服务通知']
  },
  {
    id: 2,
    title: '新增数字人民币支付方式',
    summary: '为丰富支付渠道，平台现已支持数字人民币支付，欢迎体验使用。',
    content: '<p>尊敬的用户：</p><p>我们很高兴地通知您，平台现已正式支持<strong>数字人民币</strong>支付方式。</p><p>您可以在充值页面选择"数字人民币"进行支付，享受更安全、便捷的支付体验。</p><p>如需帮助，请联系客服。</p>',
    time: '2026-07-25',
    level: 'new',
    levelLabel: '新功能',
    tags: ['支付', '新功能']
  },
  {
    id: 3,
    title: '2026年端午节假期服务安排',
    summary: '端午节期间客服服务时间调整，提现到账时间可能延迟，请提前做好资金安排。',
    content: '<p>尊敬的用户：</p><p>2026年端午节假期期间（6月19日-6月21日），服务安排如下：</p><ul><li>在线客服服务时间调整为 9:00-18:00</li><li>提现申请将在工作日统一处理</li><li>充值服务正常运行</li></ul><p>请提前做好资金安排，祝您节日愉快！</p>',
    time: '2026-07-20',
    level: 'notice',
    levelLabel: '通知',
    tags: ['假期', '服务调整']
  },
  {
    id: 4,
    title: '账户安全升级：双重验证功能上线',
    summary: '为保障账户安全，平台新增双重验证功能，建议所有用户开启。',
    content: '<p>尊敬的用户：</p><p>为进一步提升账户安全性，我们上线了<strong>双重验证（2FA）</strong>功能。</p><p>开启后，登录时除密码外还需输入动态验证码，有效防止账户被盗。</p><p>设置路径：个人中心 → 安全设置 → 双重验证</p><p>建议所有用户开启此功能。</p>',
    time: '2026-07-18',
    level: 'important',
    levelLabel: '重要',
    tags: ['安全', '新功能']
  },
  {
    id: 5,
    title: '费率调整公告',
    summary: '自8月1日起，部分服务费率将进行调整，详情请查看公告内容。',
    content: '<p>尊敬的用户：</p><p>自<strong>2026年8月1日</strong>起，以下服务费率将进行调整：</p><ul><li>提现手续费：由0.1%调整为0.08%</li><li>跨行转账费：由2元/笔调整为1元/笔</li></ul><p>此次调整旨在为用户提供更优惠的服务，感谢您的支持。</p>',
    time: '2026-07-15',
    level: 'notice',
    levelLabel: '通知',
    tags: ['费率', '服务调整']
  },
  {
    id: 7,
    title: '新用户专享福利活动',
    summary: '即日起至月底，新注册用户首次充值可享受额外赠送优惠，详情请查看活动规则。',
    content: '<p>尊敬的新用户：</p><p>为欢迎新用户加入，即日起至<strong>2026年7月31日</strong>，开展新用户专享活动：</p><ul><li>首次充值100元，赠送10元</li><li>首次充值500元，赠送60元</li><li>首次充值1000元，赠送150元</li></ul><p>赠送金额将在充值后24小时内到账，详情请查看活动页面。</p>',
    time: '2026-07-10',
    level: 'activity',
    levelLabel: '活动',
    tags: ['活动', '新用户']
  },
  {
    id: 8,
    title: '隐私政策更新说明',
    summary: '根据最新法规要求，我们更新了隐私政策条款，请您查阅。',
    content: '<p>尊敬的用户：</p><p>根据最新法规要求，我们对《隐私政策》进行了更新。</p><p>主要更新内容：</p><ul><li>明确了数据收集范围与用途</li><li>增加了用户数据删除权利说明</li><li>优化了第三方数据共享条款</li></ul><p>请您查阅最新版本的隐私政策，继续使用即表示您同意更新后的条款。</p>',
    time: '2026-07-05',
    level: 'notice',
    levelLabel: '通知',
    tags: ['隐私', '政策更新']
  }
])

const tagTypeMap: Record<string, string> = {
  important: 'error',
  new: 'success',
  notice: 'info',
  activity: 'warning'
}

function getTagType(level: string) {
  return (tagTypeMap[level] || 'default') as 'error' | 'success' | 'info' | 'warning' | 'default'
}

const filteredList = computed(() => {
  if (!searchKeyword.value) return announcements.value
  const kw = searchKeyword.value.toLowerCase()
  return announcements.value.filter(
    a =>
      a.title.toLowerCase().includes(kw) ||
      a.summary.toLowerCase().includes(kw) ||
      a.tags.some(t => t.toLowerCase().includes(kw))
  )
})

const totalPages = computed(() => Math.ceil(filteredList.value.length / pageSize))

const paginatedList = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredList.value.slice(start, start + pageSize)
})

function handleSearch() {
  currentPage.value = 1
}
</script>

<style scoped>
.announcements-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 28px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 24px;
  font-weight: 700;
  color: #262626;
  margin: 0 0 6px 0;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.detail-section {
  animation: fadeIn 0.3s ease;
}

.detail-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(24, 144, 255, 0.06);
}

.detail-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #1890ff;
  cursor: pointer;
  font-size: 14px;
  margin-bottom: 20px;
  transition: color 0.3s;
}

.detail-back:hover {
  color: #40a9ff;
}

.detail-header {
  margin-bottom: 8px;
}

.detail-title {
  font-size: 22px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 12px 0;
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.detail-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #8c8c8c;
  font-size: 13px;
}

.detail-content {
  font-size: 15px;
  line-height: 1.8;
  color: #595959;
}

.detail-content :deep(p) {
  margin: 0 0 12px 0;
}

.detail-content :deep(ul) {
  padding-left: 20px;
  margin: 8px 0 16px 0;
}

.detail-content :deep(li) {
  margin-bottom: 6px;
}

.detail-content :deep(strong) {
  color: #262626;
}

.list-section {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.announcement-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(24, 144, 255, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;
}

.announcement-card:hover {
  box-shadow: 0 6px 24px rgba(24, 144, 255, 0.12);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #bfbfbf;
  font-size: 13px;
  flex-shrink: 0;
}

.card-summary {
  color: #595959;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 14px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.meta-tag {
  background: #f0f5ff;
  color: #1890ff;
}

.read-more {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #1890ff;
  font-size: 13px;
  font-weight: 500;
  transition: color 0.3s;
}

.announcement-card:hover .read-more {
  color: #40a9ff;
}

.empty-state {
  padding: 80px 0;
  text-align: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}
</style>
