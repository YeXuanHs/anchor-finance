<template>
  <div class="knowledge-page">
    <!-- Header -->
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <n-icon size="24" color="#fff"><AnchorOutline /></n-icon>
          </div>
          <span class="logo-text">锚点财务</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <a href="#" class="nav-link">公告</a>
          <router-link to="/knowledge" class="nav-link active">帮助中心</router-link>
        </nav>
        <div class="header-actions">
          <n-button text @click="$router.push('/login')">登录</n-button>
          <n-button type="primary" round size="small" @click="$router.push('/register')">免费注册</n-button>
        </div>
      </div>
    </header>

    <!-- Breadcrumb -->
    <div class="breadcrumb-bar">
      <div class="breadcrumb-inner">
        <n-breadcrumb>
          <n-breadcrumb-item @click="$router.push('/')">首页</n-breadcrumb-item>
          <n-breadcrumb-item>帮助中心</n-breadcrumb-item>
        </n-breadcrumb>
      </div>
    </div>

    <!-- Search Banner -->
    <div class="search-banner">
      <div class="search-inner">
        <h1 class="search-title">帮助中心</h1>
        <p class="search-subtitle">搜索您需要的帮助信息</p>
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索文章标题或关键词..."
          size="large"
          round
          clearable
          class="search-input"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Sidebar -->
        <aside class="sidebar">
          <div class="sidebar-card">
            <h3 class="sidebar-title">
              <n-icon size="18" color="#1890ff"><FolderOpenOutline /></n-icon>
              文章分类
            </h3>
            <n-tree
              :data="categoryTree"
              :selected-keys="selectedCategories"
              selectable
              block-line
              :render-suffix="renderCategorySuffix"
              @update:selected-keys="onCategorySelect"
            />
          </div>

          <!-- Hot Articles -->
          <div class="sidebar-card hot-articles">
            <h3 class="sidebar-title">
              <n-icon size="18" color="#ff7a45"><FlameOutline /></n-icon>
              热门文章
            </h3>
            <div class="hot-list">
              <div
                v-for="(article, index) in hotArticles"
                :key="article.id"
                class="hot-item"
                @click="toggleArticle(article.id)"
              >
                <span class="hot-rank" :class="{ 'top-3': index < 3 }">{{ index + 1 }}</span>
                <span class="hot-title">{{ article.title }}</span>
                <span class="hot-views">{{ formatViews(article.views) }}</span>
              </div>
            </div>
          </div>
        </aside>

        <!-- Article List -->
        <main class="article-list">
          <!-- Active Category Tags -->
          <div class="filter-tags">
            <n-tag
              v-for="cat in activeCategories"
              :key="cat"
              closable
              type="info"
              size="small"
              @close="removeCategory(cat)"
            >
              {{ cat }}
            </n-tag>
            <n-button
              v-if="activeCategories.length > 0"
              text
              type="primary"
              size="small"
              @click="clearCategories"
            >
              清除筛选
            </n-button>
          </div>

          <!-- Result Info -->
          <div class="result-info">
            <span>共 <strong>{{ filteredArticles.length }}</strong> 篇文章</span>
          </div>

          <!-- Article Cards -->
          <div class="articles-grid">
            <n-card
              v-for="article in paginatedArticles"
              :key="article.id"
              class="article-card"
              hoverable
              @click="toggleArticle(article.id)"
            >
              <div class="article-header">
                <n-tag :type="getCategoryTagType(article.category)" size="small" :bordered="false">
                  {{ article.category }}
                </n-tag>
                <n-icon
                  size="16"
                  class="expand-icon"
                  :class="{ expanded: expandedArticle === article.id }"
                >
                  <ChevronDownOutline />
                </n-icon>
              </div>
              <h3 class="article-title">{{ article.title }}</h3>
              <p class="article-summary">{{ article.summary }}</p>
              <div class="article-meta">
                <span class="meta-item">
                  <n-icon size="14"><EyeOutline /></n-icon>
                  {{ formatViews(article.views) }}
                </span>
                <span class="meta-item">
                  <n-icon size="14"><ThumbsUpOutline /></n-icon>
                  {{ article.helpful }}
                </span>
                <span class="meta-item">
                  <n-icon size="14"><TimeOutline /></n-icon>
                  {{ article.updateTime }}
                </span>
              </div>

              <!-- Expanded Detail -->
              <n-collapse v-if="expandedArticle === article.id" :default-expanded-names="['detail']">
                <n-collapse-item name="detail">
                  <div class="article-detail" v-html="article.content"></div>
                </n-collapse-item>
              </n-collapse>
            </n-card>
          </div>

          <!-- Empty State -->
          <div v-if="filteredArticles.length === 0" class="empty-state">
            <n-icon size="64" color="#c9cdd4"><DocumentTextOutline /></n-icon>
            <p>暂无匹配的文章</p>
            <n-button type="primary" @click="clearFilters">清除筛选</n-button>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="pagination-wrap">
            <n-pagination
              v-model:page="currentPage"
              :page-count="totalPages"
              :page-slot="7"
              show-quick-jumper
            />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import {
  AnchorOutline,
  SearchOutline,
  FolderOpenOutline,
  FlameOutline,
  ChevronDownOutline,
  EyeOutline,
  ThumbsUpOutline,
  TimeOutline,
  DocumentTextOutline,
  ServerOutline,
  ShieldCheckmarkOutline,
  CardOutline,
  SettingsOutline,
  HelpCircleOutline
} from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'
import type { TreeOption } from 'naive-ui'

const searchKeyword = ref('')
const selectedCategories = ref<string[]>([])
const expandedArticle = ref<number | null>(null)
const currentPage = ref(1)
const pageSize = 6

interface Article {
  id: number
  title: string
  category: string
  summary: string
  content: string
  views: number
  helpful: number
  updateTime: string
}

const articles = ref<Article[]>([
  {
    id: 1,
    title: '如何购买云服务器',
    category: '购买指南',
    summary: '详细介绍云服务器的购买流程，包括选择配置、下单支付、开通服务等步骤。',
    content: '<p>1. 登录账号后进入产品中心，选择适合的云服务器产品。</p><p>2. 根据业务需求选择机房线路、CPU、内存、带宽等配置。</p><p>3. 确认订单信息后完成支付，系统将在5分钟内自动开通服务。</p><p>4. 开通成功后，服务器信息将发送至您的邮箱和站内消息。</p>',
    views: 3562,
    helpful: 284,
    updateTime: '2025-12-15'
  },
  {
    id: 2,
    title: '云服务器远程连接教程',
    category: '使用教程',
    summary: 'Windows和Linux系统远程连接云服务器的详细步骤，包含SSH和远程桌面两种方式。',
    content: '<p><strong>Windows系统：</strong>使用远程桌面连接(mstsc)，输入服务器IP和端口，使用账号密码登录。</p><p><strong>Linux系统：</strong>使用SSH工具(如PuTTY、Xshell)，输入ssh root@服务器IP，使用密钥或密码认证。</p><p>首次登录建议修改默认密码并配置密钥登录以提升安全性。</p>',
    views: 2891,
    helpful: 236,
    updateTime: '2025-12-10'
  },
  {
    id: 3,
    title: '域名解析设置指南',
    category: '使用教程',
    summary: '域名DNS解析的完整配置教程，包括A记录、CNAME、MX记录等常用解析类型。',
    content: '<p>1. 登录控制台，进入域名管理页面。</p><p>2. 选择需要设置解析的域名，点击「解析设置」。</p><p>3. 添加A记录：将域名指向服务器IP地址。</p><p>4. 添加CNAME记录：将域名指向另一个域名。</p><p>5. 添加MX记录：用于邮箱服务配置。</p><p>解析生效时间一般为10分钟至24小时。</p>',
    views: 2145,
    helpful: 189,
    updateTime: '2025-12-08'
  },
  {
    id: 4,
    title: 'SSL证书申请与安装教程',
    category: '安全相关',
    summary: '免费SSL证书的申请流程以及在不同Web服务器上的安装配置方法。',
    content: '<p>1. 进入SSL证书管理页面，选择DV免费证书。</p><p>2. 填写域名信息，选择验证方式（DNS验证或文件验证）。</p><p>3. 完成验证后证书将自动签发。</p><p>4. 下载证书文件，根据服务器类型（Nginx/Apache/IIS）进行安装。</p><p>5. 配置完成后使用HTTPS访问测试。</p>',
    views: 1876,
    helpful: 152,
    updateTime: '2025-12-05'
  },
  {
    id: 5,
    title: '账单与支付方式说明',
    category: '购买指南',
    summary: '平台支持的支付方式、账单周期、续费提醒及发票申请等相关说明。',
    content: '<p><strong>支持的支付方式：</strong>支付宝、微信支付、银行卡转账、余额支付。</p><p><strong>账单周期：</strong>月付、季付、半年付、年付，周期越长优惠越大。</p><p><strong>续费提醒：</strong>服务到期前7天、3天、1天发送续费提醒。</p><p><strong>发票申请：</strong>支持增值税普通发票和专用发票，可在控制台申请。</p>',
    views: 1654,
    helpful: 128,
    updateTime: '2025-12-01'
  },
  {
    id: 6,
    title: '服务器备份与恢复指南',
    category: '使用教程',
    summary: '如何创建服务器快照、自动备份策略设置以及数据恢复操作步骤。',
    content: '<p>1. 进入服务器管理页面，选择「快照/备份」选项。</p><p>2. 手动创建快照：点击创建快照，输入名称和描述。</p><p>3. 设置自动备份：选择备份周期（每日/每周）和保留数量。</p><p>4. 数据恢复：选择快照点，点击恢复即可回滚到该时间点。</p><p>建议定期创建快照，重要操作前务必先备份。</p>',
    views: 1432,
    helpful: 115,
    updateTime: '2025-11-28'
  },
  {
    id: 7,
    title: '服务器安全加固建议',
    category: '安全相关',
    summary: '云服务器安全加固的最佳实践，包括防火墙配置、端口管理、入侵检测等。',
    content: '<p>1. 修改SSH默认端口，禁用root远程登录。</p><p>2. 配置防火墙规则，仅开放必要端口。</p><p>3. 安装fail2ban防止暴力破解。</p><p>4. 定期更新系统和软件补丁。</p><p>5. 启用登录日志审计和异常告警。</p><p>6. 使用密钥对替代密码认证。</p>',
    views: 1298,
    helpful: 103,
    updateTime: '2025-11-25'
  },
  {
    id: 8,
    title: '如何申请退款',
    category: '售后支持',
    summary: '退款政策说明、退款申请流程以及退款到账时间等常见问题解答。',
    content: '<p><strong>退款政策：</strong>新购产品5天内可申请无理由退款（已使用的流量费用除外）。</p><p>1. 进入订单管理页面，找到需要退款的订单。</p><p>2. 点击「申请退款」，选择退款原因并提交。</p><p>3. 客服将在24小时内审核退款申请。</p><p>4. 审核通过后，退款将在3-5个工作日内原路退回。</p>',
    views: 987,
    helpful: 76,
    updateTime: '2025-11-20'
  },
  {
    id: 9,
    title: '服务器性能优化指南',
    category: '使用教程',
    summary: '提升云服务器运行性能的实用技巧，涵盖系统调优、Web服务优化和缓存配置。',
    content: '<p><strong>系统层面：</strong>调整内核参数、优化文件系统、合理配置Swap。</p><p><strong>Web服务：</strong>启用Gzip压缩、配置浏览器缓存、优化数据库查询。</p><p><strong>缓存加速：</strong>使用Redis/Memcached缓存热点数据，配置CDN加速静态资源。</p><p><strong>监控排查：</strong>使用top/htop、vmstat、iostat等工具定位性能瓶颈。</p>',
    views: 876,
    helpful: 68,
    updateTime: '2025-11-18'
  },
  {
    id: 10,
    title: '工单提交与处理流程',
    category: '售后支持',
    summary: '如何提交技术支持工单、工单的处理流程以及紧急问题的联系方式。',
    content: '<p>1. 登录控制台，进入「工单中心」。</p><p>2. 选择问题分类（技术咨询/故障报备/账户问题）。</p><p>3. 详细描述问题并附上相关截图或日志。</p><p>4. 普通工单24小时内响应，紧急工单2小时内响应。</p><p>紧急问题可拨打7x24小时客服热线：400-XXX-XXXX。</p>',
    views: 765,
    helpful: 54,
    updateTime: '2025-11-15'
  }
])

const categoryTree = computed<TreeOption[]>(() => {
  const categories = [...new Set(articles.value.map(a => a.category))]
  return [
    {
      key: 'all',
      label: '全部文章',
      prefix: () => h(NIcon, { size: 16, color: '#1890ff' }, { default: () => h(DocumentTextOutline) }),
      suffix: () => h('span', { class: 'tree-count' }, articles.value.length),
      children: categories.map(cat => ({
        key: cat,
        label: cat,
        prefix: () => h(NIcon, { size: 16, color: '#86909c' }, { default: () => h(HelpCircleOutline) }),
        suffix: () => h('span', { class: 'tree-count' }, articles.value.filter(a => a.category === cat).length)
      }))
    }
  ]
})

function renderCategorySuffix({ option }: { option: TreeOption }) {
  return option.suffix ? option.suffix() : null
}

function onCategorySelect(keys: string[]) {
  const key = keys[0]
  if (!key || key === 'all') {
    selectedCategories.value = []
  } else {
    selectedCategories.value = [key]
  }
  currentPage.value = 1
}

const activeCategories = computed(() => selectedCategories.value)

function removeCategory(cat: string) {
  selectedCategories.value = selectedCategories.value.filter(c => c !== cat)
}

function clearCategories() {
  selectedCategories.value = []
}

function clearFilters() {
  searchKeyword.value = ''
  selectedCategories.value = []
  currentPage.value = 1
}

const hotArticles = computed(() => {
  return [...articles.value].sort((a, b) => b.views - a.views).slice(0, 5)
})

const filteredArticles = computed(() => {
  let list = [...articles.value]

  if (selectedCategories.value.length > 0) {
    list = list.filter(a => selectedCategories.value.includes(a.category))
  }

  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.trim().toLowerCase()
    list = list.filter(a =>
      a.title.toLowerCase().includes(keyword) || a.summary.toLowerCase().includes(keyword)
    )
  }

  return list
})

const totalPages = computed(() => Math.ceil(filteredArticles.value.length / pageSize))

const paginatedArticles = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredArticles.value.slice(start, start + pageSize)
})

function toggleArticle(id: number) {
  expandedArticle.value = expandedArticle.value === id ? null : id
}

function formatViews(views: number): string {
  if (views >= 10000) return (views / 10000).toFixed(1) + 'w'
  if (views >= 1000) return (views / 1000).toFixed(1) + 'k'
  return views.toString()
}

function getCategoryTagType(category: string): 'info' | 'success' | 'warning' | 'error' {
  const map: Record<string, 'info' | 'success' | 'warning' | 'error'> = {
    '购买指南': 'info',
    '使用教程': 'success',
    '安全相关': 'warning',
    '售后支持': 'error'
  }
  return map[category] || 'info'
}
</script>

<style scoped>
.knowledge-page {
  min-height: 100vh;
  background: #f7f8fa;
}

/* Header */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: #fff;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.06);
  z-index: 100;
}

.header-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  width: 34px;
  height: 34px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
}

.nav-links {
  display: flex;
  gap: 32px;
}

.nav-link {
  color: #4e5969;
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-link:hover,
.nav-link.active {
  color: #1890ff;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Breadcrumb */
.breadcrumb-bar {
  background: #fff;
  border-bottom: 1px solid #f0f1f5;
  margin-top: 64px;
}

.breadcrumb-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 14px 24px;
}

/* Search Banner */
.search-banner {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  padding: 48px 24px;
}

.search-inner {
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
}

.search-title {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
}

.search-subtitle {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 24px;
}

.search-input {
  max-width: 500px;
}

/* Main Content */
.main-content {
  padding: 24px 0 40px;
}

.content-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  gap: 24px;
}

/* Sidebar */
.sidebar {
  width: 280px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.sidebar-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.tree-count) {
  font-size: 12px;
  color: #c9cdd4;
  background: #f2f3f5;
  padding: 1px 8px;
  border-radius: 10px;
  margin-left: 8px;
}

/* Hot Articles */
.hot-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.hot-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.hot-item:hover {
  background: #f7f8fa;
}

.hot-rank {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  background: #f2f3f5;
  color: #86909c;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.hot-rank.top-3 {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  color: #fff;
}

.hot-title {
  flex: 1;
  font-size: 13px;
  color: #4e5969;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hot-views {
  font-size: 12px;
  color: #c9cdd4;
  flex-shrink: 0;
}

/* Article List */
.article-list {
  flex: 1;
  min-width: 0;
}

.filter-tags {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.result-info {
  font-size: 14px;
  color: #86909c;
  margin-bottom: 16px;
}

.result-info strong {
  color: #1890ff;
}

.articles-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article-card {
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 12px;
}

.article-card:hover {
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.1);
}

.article-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.expand-icon {
  transition: transform 0.3s;
  color: #c9cdd4;
}

.expand-icon.expanded {
  transform: rotate(180deg);
  color: #1890ff;
}

.article-title {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 8px;
}

.article-summary {
  font-size: 13px;
  color: #86909c;
  line-height: 1.6;
  margin-bottom: 12px;
}

.article-meta {
  display: flex;
  gap: 20px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #c9cdd4;
}

.article-detail {
  padding: 16px;
  margin-top: 12px;
  background: #f7f8fa;
  border-radius: 8px;
  font-size: 14px;
  color: #4e5969;
  line-height: 1.8;
}

.article-detail :deep(p) {
  margin-bottom: 8px;
}

.article-detail :deep(strong) {
  color: #1d2129;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 80px 0;
  color: #c9cdd4;
}

.empty-state p {
  margin: 16px 0 24px;
  font-size: 15px;
}

/* Pagination */
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

/* Responsive */
@media (max-width: 1024px) {
  .sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
}
</style>
