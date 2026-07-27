<template>
  <div class="announcements-page">
    <div class="page-header">
      <h2 class="page-title">公告中心</h2>
      <p class="page-desc">了解最新动态和重要通知</p>
    </div>

    <!-- 搜索栏 -->
    <n-card class="search-card" :bordered="false">
      <div class="search-bar">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索公告标题或内容..."
          size="large"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" color="#1890ff" />
          </template>
        </n-input>
        <n-button type="primary" size="large" @click="handleSearch">
          搜索
        </n-button>
      </div>
    </n-card>

    <!-- 公告列表 -->
    <div v-if="filteredAnnouncements.length > 0" class="announcement-list">
      <n-card
        v-for="item in paginatedAnnouncements"
        :key="item.id"
        class="announcement-card"
        :bordered="false"
        hoverable
        @click="toggleDetail(item.id)"
      >
        <div class="announcement-header">
          <div class="announcement-title-row">
            <h3 class="announcement-title">{{ item.title }}</h3>
            <n-tag
              v-for="tag in item.tags"
              :key="tag"
              :type="getTagType(tag)"
              size="small"
              round
            >
              {{ tag }}
            </n-tag>
          </div>
          <span class="announcement-time">
            <n-icon :component="TimeOutline" size="14" />
            <n-time :time="new Date(item.publishAt)" format="yyyy-MM-dd HH:mm" />
          </span>
        </div>

        <p class="announcement-summary">{{ item.summary }}</p>

        <div v-if="expandedId === item.id" class="announcement-detail">
          <div class="detail-divider"></div>
          <div class="detail-content" v-html="item.content"></div>
        </div>

        <div class="announcement-footer">
          <span class="read-more">
            {{ expandedId === item.id ? '收起详情' : '查看详情' }}
            <n-icon :component="expandedId === item.id ? ChevronUpOutline : ChevronDownOutline" size="14" />
          </span>
        </div>
      </n-card>
    </div>

    <!-- 空状态 -->
    <n-card v-else class="empty-card" :bordered="false">
      <n-empty description="暂无公告" size="large">
        <template #extra>
          <n-button @click="searchKeyword = ''">清除搜索</n-button>
        </template>
      </n-empty>
    </n-card>

    <!-- 分页 -->
    <div v-if="filteredAnnouncements.length > 0" class="pagination-wrapper">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-slot="7"
        show-quick-jumper
      >
        <template #prefix="{ itemCount }">
          共 {{ filteredAnnouncements.length }} 条公告
        </template>
      </n-pagination>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  SearchOutline,
  TimeOutline,
  ChevronDownOutline,
  ChevronUpOutline
} from '@vicons/ionicons5'

const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = 5
const expandedId = ref<number | null>(null)

interface Announcement {
  id: number
  title: string
  summary: string
  content: string
  publishAt: string
  tags: string[]
}

const announcements = ref<Announcement[]>([
  {
    id: 1,
    title: '系统升级维护通知',
    summary: '为提升系统性能和用户体验，我们将于本周六凌晨2:00-6:00进行系统升级维护，届时服务将暂时不可用。',
    content: '<p>尊敬的用户：</p><p>为提升系统性能和用户体验，我们将于 <strong>2026年7月29日（周六）凌晨 2:00-6:00</strong> 进行系统升级维护。</p><p><strong>影响范围：</strong></p><ul><li>所有在线服务将暂时不可用</li><li>API接口将暂停响应</li><li>数据不会丢失，请放心</li></ul><p>维护完成后将恢复正常服务，请提前做好相关安排。感谢您的理解与支持！</p>',
    publishAt: '2026-07-27T08:00:00',
    tags: ['系统', '维护']
  },
  {
    id: 2,
    title: '新功能上线：智能报表分析',
    summary: '全新的智能报表分析功能已正式上线，支持多维度数据可视化和AI驱动的业务洞察。',
    content: '<p>我们很高兴地宣布，<strong>智能报表分析功能</strong>已正式上线！</p><p><strong>主要特性：</strong></p><ul><li>多维度数据可视化图表</li><li>AI驱动的业务趋势预测</li><li>自定义报表模板</li><li>一键导出PDF/Excel</li></ul><p>立即前往「报表中心」体验新功能！</p>',
    publishAt: '2026-07-25T10:00:00',
    tags: ['新功能', '更新']
  },
  {
    id: 3,
    title: '2026年端午节放假安排',
    summary: '根据国家法定节假日安排，2026年端午节放假时间为5月31日至6月2日，共3天。',
    content: '<p>尊敬的用户：</p><p>根据国家法定节假日安排，2026年端午节放假时间为 <strong>5月31日至6月2日</strong>，共3天。</p><p><strong>假期服务安排：</strong></p><ul><li>在线客服照常运行</li><li>技术支持响应时间可能延长</li><li>财务结算顺延至节后处理</li></ul><p>提前祝大家端午安康！</p>',
    publishAt: '2026-07-20T09:00:00',
    tags: ['公告']
  },
  {
    id: 4,
    title: '安全提醒：谨防钓鱼邮件',
    summary: '近期发现有不法分子冒充我司发送钓鱼邮件，请勿点击不明链接，保护好您的账户安全。',
    content: '<p>尊敬的用户：</p><p>近期我们发现有不法分子冒充我司名义发送钓鱼邮件，诱导用户点击恶意链接。</p><p><strong>请注意：</strong></p><ul><li>我司官方邮件域名仅为 @anchorfinance.com</li><li>不会通过邮件要求您提供密码</li><li>遇到可疑邮件请及时联系客服核实</li></ul><p>如已误点击，请立即修改密码并联系客服。</p>',
    publishAt: '2026-07-18T14:30:00',
    tags: ['安全', '重要']
  },
  {
    id: 5,
    title: 'API接口文档更新公告',
    summary: 'API接口文档已全面更新，新增WebSocket推送接口和批量操作接口说明。',
    content: '<p>开发者您好：</p><p>API接口文档已完成全面更新，主要变更包括：</p><ul><li>新增 WebSocket 实时数据推送接口</li><li>新增批量操作接口（批量下单、批量查询）</li><li>优化认证流程说明</li><li>新增错误码完整列表</li></ul><p>请前往「开发文档」查看最新版本。</p>',
    publishAt: '2026-07-15T11:00:00',
    tags: ['开发', '更新']
  },
  {
    id: 6,
    title: '用户协议与隐私政策更新',
    summary: '根据最新法规要求，我们更新了用户协议和隐私政策，新版本将于2026年8月1日起生效。',
    content: '<p>尊敬的用户：</p><p>根据最新法规要求，我们对用户协议和隐私政策进行了更新。</p><p><strong>主要变更：</strong></p><ul><li>完善数据收集和使用说明</li><li>增加用户数据删除权利说明</li><li>明确第三方数据共享范围</li></ul><p>新版本将于 <strong>2026年8月1日</strong> 起生效，继续使用服务即视为同意更新后的协议。</p>',
    publishAt: '2026-07-10T16:00:00',
    tags: ['公告', '重要']
  },
  {
    id: 7,
    title: '产品价格调整通知',
    summary: '自2026年8月起，部分产品价格将进行适度调整，已购买的服务不受影响。',
    content: '<p>尊敬的用户：</p><p>由于运营成本调整，自 <strong>2026年8月1日</strong> 起，部分产品价格将进行适度调整。</p><p><strong>调整说明：</strong></p><ul><li>已购买的服务在有效期内不受影响</li><li>续费将按新价格执行</li><li>年付用户享受额外折扣</li></ul><p>感谢您的理解与支持！</p>',
    publishAt: '2026-07-08T10:00:00',
    tags: ['公告']
  }
])

const filteredAnnouncements = computed(() => {
  if (!searchKeyword.value.trim()) return announcements.value
  const keyword = searchKeyword.value.toLowerCase()
  return announcements.value.filter(
    item =>
      item.title.toLowerCase().includes(keyword) ||
      item.summary.toLowerCase().includes(keyword) ||
      item.tags.some(tag => tag.toLowerCase().includes(keyword))
  )
})

const totalPages = computed(() => Math.ceil(filteredAnnouncements.value.length / pageSize))

const paginatedAnnouncements = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredAnnouncements.value.slice(start, start + pageSize)
})

function handleSearch() {
  currentPage.value = 1
}

function toggleDetail(id: number) {
  expandedId.value = expandedId.value === id ? null : id
}

function getTagType(tag: string): 'success' | 'warning' | 'error' | 'info' | 'default' {
  const map: Record<string, 'success' | 'warning' | 'error' | 'info' | 'default'> = {
    '重要': 'error',
    '安全': 'warning',
    '系统': 'info',
    '维护': 'info',
    '新功能': 'success',
    '更新': 'success',
    '开发': 'default',
    '公告': 'default'
  }
  return map[tag] || 'default'
}
</script>

<style scoped>
.announcements-page {
  max-width: 860px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 4px;
}

.page-desc {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.search-card {
  margin-bottom: 20px;
  border-radius: 12px;
}

.search-bar {
  display: flex;
  gap: 12px;
}

.search-bar .n-input {
  flex: 1;
}

.search-bar .n-button {
  border-radius: 8px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border: none;
  min-width: 80px;
}

.search-bar .n-button:hover {
  background: linear-gradient(135deg, #40a9ff, #1890ff);
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.announcement-card {
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.announcement-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(24, 144, 255, 0.12);
}

.announcement-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 8px;
}

.announcement-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.announcement-title {
  font-size: 17px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.announcement-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #bfbfbf;
  font-size: 13px;
  white-space: nowrap;
  flex-shrink: 0;
}

.announcement-summary {
  color: #595959;
  font-size: 14px;
  line-height: 1.6;
  margin: 0;
}

.announcement-detail {
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    max-height: 0;
  }
  to {
    opacity: 1;
    max-height: 500px;
  }
}

.detail-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, #e8e8e8, transparent);
  margin: 16px 0;
}

.detail-content {
  color: #595959;
  font-size: 14px;
  line-height: 1.8;
}

.detail-content :deep(p) {
  margin: 8px 0;
}

.detail-content :deep(ul) {
  padding-left: 20px;
  margin: 8px 0;
}

.detail-content :deep(li) {
  margin: 4px 0;
}

.detail-content :deep(strong) {
  color: #1a1a1a;
}

.announcement-footer {
  margin-top: 12px;
}

.read-more {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #1890ff;
  font-size: 13px;
  cursor: pointer;
  transition: color 0.3s;
}

.read-more:hover {
  color: #40a9ff;
}

.empty-card {
  border-radius: 12px;
  padding: 60px 0;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding-bottom: 16px;
}
</style>
