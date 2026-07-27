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
          <router-link to="/knowledge" class="nav-link active">帮助</router-link>
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
        <p class="search-desc">搜索常见问题，快速找到解决方案</p>
        <div class="search-box">
          <n-input
            v-model:value="searchKeyword"
            placeholder="输入关键词搜索帮助文章..."
            size="large"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <n-icon :component="SearchOutline" />
            </template>
          </n-input>
          <n-button type="primary" size="large" @click="handleSearch">搜索</n-button>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Sidebar -->
        <aside class="sidebar">
          <div class="sidebar-card">
            <h3 class="sidebar-title">
              <n-icon size="18" color="#1890ff"><ListOutline /></n-icon>
              文章分类
            </h3>
            <n-tree
              :data="categoryTree"
              :selected-keys="selectedCategoryKeys"
              selectable
              block-line
              @update:selected-keys="handleCategorySelect"
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
                @click="openArticle(article)"
              >
                <span class="hot-rank" :class="{ top: index < 3 }">{{ index + 1 }}</span>
                <span class="hot-title">{{ article.title }}</span>
              </div>
            </div>
          </div>
        </aside>

        <!-- Article List -->
        <main class="article-list">
          <!-- Toolbar -->
          <div class="toolbar">
            <div class="toolbar-left">
              <span class="result-count">共 <strong>{{ filteredArticles.length }}</strong> 篇文章</span>
              <n-tag v-if="activeCategoryLabel" type="info" size="small" :bordered="false" closable @close="clearCategory">
                {{ activeCategoryLabel }}
              </n-tag>
            </div>
            <div class="toolbar-right">
              <n-select
                v-model:value="sortBy"
                :options="sortOptions"
                size="small"
                style="width: 140px;"
                placeholder="排序方式"
              />
            </div>
          </div>

          <!-- Article Cards -->
          <div class="articles-grid">
            <div
              v-for="article in paginatedArticles"
              :key="article.id"
              class="article-card"
              @click="toggleArticle(article)"
            >
              <div class="article-header">
                <div class="article-meta">
                  <n-tag :type="getCategoryTagType(article.category)" size="small" :bordered="false">
                    {{ article.category }}
                  </n-tag>
                  <span class="article-views">
                    <n-icon size="14"><EyeOutline /></n-icon>
                    {{ article.views }}
                  </span>
                  <span class="article-helpful">
                    <n-icon size="14"><ThumbsUpOutline /></n-icon>
                    {{ article.helpful }}
                  </span>
                </div>
              </div>
              <h3 class="article-title">{{ article.title }}</h3>
              <p class="article-summary">{{ article.summary }}</p>

              <!-- Expanded Content -->
              <n-collapse v-if="expandedArticle === article.id">
                <n-collapse-item :name="article.id" :title="null">
                  <div class="article-content" v-html="article.content"></div>
                </n-collapse-item>
              </n-collapse>

              <div class="article-footer">
                <span class="article-date">
                  <n-icon size="14"><TimeOutline /></n-icon>
                  {{ article.date }}
                </span>
                <n-button text type="primary" size="small">
                  {{ expandedArticle === article.id ? '收起' : '展开阅读' }}
                  <template #icon>
                    <n-icon>
                      <component :is="expandedArticle === article.id ? ChevronUpOutline : ChevronDownOutline" />
                    </n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="filteredArticles.length === 0" class="empty-state">
            <n-empty description="暂无匹配的文章" />
            <n-button type="primary" class="empty-btn" @click="resetFilters">重置筛选</n-button>
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
import { ref, computed } from 'vue'
import type { TreeOption } from 'naive-ui'
import {
  AnchorOutline,
  ListOutline,
  SearchOutline,
  FlameOutline,
  EyeOutline,
  ThumbsUpOutline,
  TimeOutline,
  ChevronUpOutline,
  ChevronDownOutline,
  ServerOutline,
  GlobeOutline,
  ShieldCheckmarkOutline,
  CardOutline,
  SettingsOutline
} from '@vicons/ionicons5'

const searchKeyword = ref('')
const selectedCategoryKeys = ref<string[]>(['all'])
const sortBy = ref('default')
const currentPage = ref(1)
const pageSize = 6
const expandedArticle = ref<number | null>(null)

const sortOptions = [
  { label: '默认排序', value: 'default' },
  { label: '最多浏览', value: 'views' },
  { label: '最有帮助', value: 'helpful' },
  { label: '最新发布', value: 'newest' }
]

const categoryTree: TreeOption[] = [
  {
    key: 'all',
    label: '全部文章'
  },
  {
    key: 'quickstart',
    label: '快速入门',
    children: [
      { key: 'quickstart-register', label: '注册与认证' },
      { key: 'quickstart-first', label: '首次购买' }
    ]
  },
  {
    key: 'product',
    label: '产品使用',
    children: [
      { key: 'product-ecs', label: '云服务器' },
      { key: 'product-vps', label: 'VPS' },
      { key: 'product-domain', label: '域名管理' }
    ]
  },
  {
    key: 'billing',
    label: '费用与账单',
    children: [
      { key: 'billing-pay', label: '支付方式' },
      { key: 'billing-refund', label: '退款说明' }
    ]
  },
  {
    key: 'security',
    label: '安全相关',
    children: [
      { key: 'security-ssl', label: 'SSL证书' },
      { key: 'security-ddos', label: 'DDoS防护' }
    ]
  },
  {
    key: 'faq',
    label: '常见问题'
  }
]

const activeCategoryLabel = computed(() => {
  const key = selectedCategoryKeys.value[0]
  if (!key || key === 'all') return ''
  const find = (nodes: TreeOption[]): string | undefined => {
    for (const n of nodes) {
      if (n.key === key) return n.label as string
      if (n.children) {
        const r = find(n.children)
        if (r) return r
      }
    }
  }
  return find(categoryTree) ?? ''
})

interface Article {
  id: number
  title: string
  category: string
  categoryKey: string
  summary: string
  content: string
  views: number
  helpful: number
  date: string
}

const articles = ref<Article[]>([
  {
    id: 1,
    title: '如何注册并完成实名认证',
    category: '快速入门',
    categoryKey: 'quickstart',
    summary: '本教程将引导您完成账号注册和实名认证流程，认证后可享受更多功能和优惠。',
    content: '<p><strong>第一步：</strong>访问官网首页，点击右上角"免费注册"按钮。</p><p><strong>第二步：</strong>填写邮箱、手机号和密码，完成基础信息注册。</p><p><strong>第三步：</strong>进入控制台，选择"账号设置" → "实名认证"，上传身份证正反面照片。</p><p><strong>第四步：</strong>等待审核（通常1-2个工作日），认证完成后即可购买产品。</p>',
    views: 3420,
    helpful: 286,
    date: '2026-07-15'
  },
  {
    id: 2,
    title: '云服务器购买与初始化配置指南',
    category: '云服务器',
    categoryKey: 'product-ecs',
    summary: '详细介绍如何选择合适的云服务器配置、购买流程以及初始化环境设置。',
    content: '<p>本文将帮助您了解如何选购最适合您业务需求的云服务器。</p><p><strong>选择地域：</strong>根据目标用户分布选择最近的机房，如面向大陆用户推荐香港节点。</p><p><strong>选择配置：</strong>个人建站推荐1核2G起步，企业应用推荐2核4G以上。</p><p><strong>系统选择：</strong>支持CentOS、Ubuntu、Debian、Windows Server等主流系统。</p><p><strong>安全设置：</strong>购买后请立即修改默认密码，配置安全组规则。</p>',
    views: 2856,
    helpful: 198,
    date: '2026-07-10'
  },
  {
    id: 3,
    title: '域名注册与DNS解析设置教程',
    category: '域名管理',
    categoryKey: 'product-domain',
    summary: '从域名注册到DNS解析配置，手把手教您完成域名相关操作。',
    content: '<p><strong>域名注册：</strong>在产品中心选择"域名注册"，搜索您想要的域名，加入购物车并完成支付。</p><p><strong>DNS解析：</strong>进入控制台 → 域名管理 → DNS解析，添加A记录、CNAME记录等。</p><p><strong>注意事项：</strong>域名解析生效通常需要10分钟至24小时，取决于TTL设置和DNS缓存。</p>',
    views: 2140,
    helpful: 165,
    date: '2026-07-08'
  },
  {
    id: 4,
    title: '支持哪些支付方式及充值说明',
    category: '费用与账单',
    categoryKey: 'billing-pay',
    summary: '了解平台支持的所有支付方式，以及余额充值、发票申请等常见问题。',
    content: '<p><strong>支付方式：</strong>目前支持支付宝、微信支付、银行卡转账、USDT等支付方式。</p><p><strong>余额充值：</strong>控制台 → 财务中心 → 余额充值，最低充值金额10元。</p><p><strong>发票申请：</strong>完成实名认证后，可在"财务中心"申请增值税普通发票或专用发票。</p>',
    views: 1890,
    helpful: 142,
    date: '2026-07-05'
  },
  {
    id: 5,
    title: 'SSL证书申请与安装全流程',
    category: 'SSL证书',
    categoryKey: 'security-ssl',
    summary: '免费SSL证书申请和付费证书购买后的安装部署完整指南。',
    content: '<p><strong>免费证书：</strong>DV SSL证书支持免费申请，10分钟内签发。</p><p><strong>申请流程：</strong>控制台 → SSL证书 → 申请证书 → 选择域名验证方式（DNS验证或文件验证）。</p><p><strong>安装部署：</strong>证书签发后，下载对应服务器格式的证书文件，按文档配置到Nginx/Apache/IIS。</p>',
    views: 1654,
    helpful: 128,
    date: '2026-07-01'
  },
  {
    id: 6,
    title: 'VPS与云服务器有什么区别',
    category: '常见问题',
    categoryKey: 'faq',
    summary: '从技术架构、性能、价格等方面对比VPS和云服务器的差异，帮您做出最佳选择。',
    content: '<p><strong>架构差异：</strong>VPS基于KVM虚拟化技术，资源共享但性能隔离；云服务器采用分布式架构，资源完全独享。</p><p><strong>性能对比：</strong>云服务器支持弹性伸缩、热迁移，高可用性更强；VPS适合中小规模应用。</p><p><strong>价格差异：</strong>VPS价格更低，适合个人和小型项目；云服务器适合企业级应用。</p>',
    views: 3100,
    helpful: 245,
    date: '2026-06-28'
  },
  {
    id: 7,
    title: '如何申请退款及退款政策说明',
    category: '退款说明',
    categoryKey: 'billing-refund',
    summary: '了解退款条件、退款流程和退款到账时间等详细说明。',
    content: '<p><strong>退款条件：</strong>购买后5天内可申请无理由退款（域名、SSL证书等特殊产品除外）。</p><p><strong>退款流程：</strong>控制台 → 工单中心 → 提交退款工单，说明退款原因。</p><p><strong>到账时间：</strong>审核通过后，退款将在3-5个工作日内原路返回。</p>',
    views: 980,
    helpful: 76,
    date: '2026-06-25'
  },
  {
    id: 8,
    title: 'DDoS高防产品使用指南',
    category: 'DDoS防护',
    categoryKey: 'security-ddos',
    summary: '了解DDoS高防产品的配置方法和最佳实践，保障业务安全稳定运行。',
    content: '<p><strong>产品概述：</strong>DDoS高防可为云服务器、独立服务器提供最高T级防护能力。</p><p><strong>接入方式：</strong>支持域名接入（CNAME方式）和IP接入（直接绑定高防IP）。</p><p><strong>防护配置：</strong>控制台 → DDoS高防 → 防护规则，可配置CC防护、黑白名单等。</p>',
    views: 1320,
    helpful: 98,
    date: '2026-06-20'
  },
  {
    id: 9,
    title: 'VPS系统重装与快照管理',
    category: 'VPS',
    categoryKey: 'product-vps',
    summary: '学习如何重装VPS操作系统、创建和恢复快照，以及常用运维操作。',
    content: '<p><strong>系统重装：</strong>控制台 → VPS管理 → 更多 → 重装系统，选择新系统镜像并设置密码。</p><p><strong>快照功能：</strong>快照可在"备份管理"中创建，建议在重大操作前创建快照以便回滚。</p><p><strong>注意事项：</strong>重装系统会清空所有数据，请提前备份重要文件。</p>',
    views: 1560,
    helpful: 112,
    date: '2026-06-18'
  }
])

const hotArticles = computed(() => {
  return [...articles.value].sort((a, b) => b.views - a.views).slice(0, 6)
})

const filteredArticles = computed(() => {
  let list = [...articles.value]

  // Category filter
  const key = selectedCategoryKeys.value[0]
  if (key && key !== 'all') {
    list = list.filter(a => a.categoryKey === key)
  }

  // Search filter
  if (searchKeyword.value.trim()) {
    const kw = searchKeyword.value.trim().toLowerCase()
    list = list.filter(a =>
      a.title.toLowerCase().includes(kw) ||
      a.summary.toLowerCase().includes(kw) ||
      a.category.toLowerCase().includes(kw)
    )
  }

  // Sort
  if (sortBy.value === 'views') {
    list.sort((a, b) => b.views - a.views)
  } else if (sortBy.value === 'helpful') {
    list.sort((a, b) => b.helpful - a.helpful)
  } else if (sortBy.value === 'newest') {
    list.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
  }

  return list
})

const totalPages = computed(() => Math.ceil(filteredArticles.value.length / pageSize))

const paginatedArticles = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredArticles.value.slice(start, start + pageSize)
})

function getCategoryTagType(category: string): 'success' | 'info' | 'warning' | 'error' | 'default' {
  const map: Record<string, 'success' | 'info' | 'warning' | 'error'> = {
    '快速入门': 'success',
    '云服务器': 'info',
    'VPS': 'info',
    '域名管理': 'warning',
    '费用与账单': 'default' as any,
    '退款说明': 'default' as any,
    'SSL证书': 'success',
    'DDoS防护': 'error',
    '常见问题': 'warning'
  }
  return (map[category] ?? 'info')
}

function handleCategorySelect(keys: string[]) {
  selectedCategoryKeys.value = keys.length ? keys : ['all']
  currentPage.value = 1
  expandedArticle.value = null
}

function clearCategory() {
  selectedCategoryKeys.value = ['all']
  currentPage.value = 1
}

function handleSearch() {
  currentPage.value = 1
  expandedArticle.value = null
}

function toggleArticle(article: Article) {
  expandedArticle.value = expandedArticle.value === article.id ? null : article.id
}

function openArticle(article: Article) {
  expandedArticle.value = article.id
  currentPage.value = 1
  // Scroll to article section
  const el = document.querySelector('.article-list')
  if (el) el.scrollIntoView({ behavior: 'smooth' })
}

function resetFilters() {
  searchKeyword.value = ''
  selectedCategoryKeys.value = ['all']
  sortBy.value = 'default'
  currentPage.value = 1
  expandedArticle.value = null
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
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  padding: 48px 0;
}

.search-inner {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 24px;
  text-align: center;
}

.search-title {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 10px;
}

.search-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: 28px;
}

.search-box {
  display: flex;
  gap: 12px;
}

.search-box :deep(.n-input) {
  flex: 1;
  border-radius: 8px;
}

.search-box .n-button {
  border-radius: 8px;
  padding: 0 28px;
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

/* Hot Articles */
.hot-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.hot-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.hot-item:hover {
  background: #f2f3f5;
}

.hot-rank {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  background: #f2f3f5;
  color: #86909c;
  flex-shrink: 0;
}

.hot-rank.top {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  color: #fff;
}

.hot-title {
  font-size: 13px;
  color: #4e5969;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Article List */
.article-list {
  flex: 1;
  min-width: 0;
}

/* Toolbar */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.result-count {
  font-size: 14px;
  color: #86909c;
}

.result-count strong {
  color: #1890ff;
}

/* Articles Grid */
.articles-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #f0f1f5;
}

.article-card:hover {
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.08);
  border-color: rgba(24, 144, 255, 0.2);
}

.article-header {
  margin-bottom: 12px;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 14px;
}

.article-views,
.article-helpful {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #c9cdd4;
}

.article-title {
  font-size: 17px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 8px;
}

.article-summary {
  font-size: 14px;
  color: #86909c;
  line-height: 1.6;
  margin-bottom: 12px;
}

.article-content {
  font-size: 14px;
  color: #4e5969;
  line-height: 1.8;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 8px;
  margin-bottom: 12px;
}

.article-content :deep(p) {
  margin-bottom: 8px;
}

.article-content :deep(p:last-child) {
  margin-bottom: 0;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 14px;
  border-top: 1px solid #f2f3f5;
}

.article-date {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #c9cdd4;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 80px 0;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.empty-btn {
  margin-top: 20px;
}

/* Pagination */
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 24px;
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

  .search-title {
    font-size: 24px;
  }

  .search-box {
    flex-direction: column;
  }
}
</style>
