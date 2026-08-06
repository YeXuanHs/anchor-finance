<template>
  <div class="knowledge-page">
    <!-- Header -->
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <n-icon size="24" color="#fff"><BookmarkOutline /></n-icon>
          </div>
          <span class="logo-text">{{ $t('landing.brandName') }}</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">{{ $t('common.home') }}</router-link>
          <router-link to="/products" class="nav-link">{{ $t('menu.products') }}</router-link>
          <a href="#" class="nav-link">{{ $t('landing.announcement') }}</a>
          <router-link to="/knowledge" class="nav-link active">{{ $t('menu.knowledgeBase') }}</router-link>
        </nav>
        <div class="header-actions">
          <n-button text @click="$router.push('/login')">{{ $t('landing.login') }}</n-button>
          <n-button type="primary" round size="small" @click="$router.push('/register')">{{ $t('landing.freeRegister') }}</n-button>
        </div>
      </div>
    </header>

    <!-- Breadcrumb -->
    <div class="breadcrumb-bar">
      <div class="breadcrumb-inner">
        <n-breadcrumb>
          <n-breadcrumb-item @click="$router.push('/')">{{ $t('common.home') }}</n-breadcrumb-item>
          <n-breadcrumb-item>{{ $t('menu.knowledgeBase') }}</n-breadcrumb-item>
        </n-breadcrumb>
      </div>
    </div>

    <!-- Search Banner -->
    <div class="search-banner">
      <div class="search-inner">
        <h1 class="search-title">{{ $t('knowledgeBase.title') }}</h1>
        <p class="search-subtitle">{{ $t('knowledgeBase.searchHelp') }}</p>
        <n-input
          v-model:value="searchKeyword"
          :placeholder="$t('knowledgeBase.searchPlaceholder')"
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
              {{ $t('knowledgeBase.articleCategories') }}
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
              {{ $t('knowledgeBase.hotArticles') }}
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
              {{ $t('knowledgeBase.clearFilter') }}
            </n-button>
          </div>

          <!-- Result Info -->
          <div class="result-info">
            <span>{{ $t('knowledgeBase.totalArticles', { count: filteredArticles.length }) }}</span>
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
            <p>{{ $t('knowledgeBase.noMatchingArticles') }}</p>
            <n-button type="primary" @click="clearFilters">{{ $t('knowledgeBase.clearFilter') }}</n-button>
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
import { useI18n } from 'vue-i18n'
import {
  BookmarkOutline,
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

const { t } = useI18n()

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

const articles = ref<Article[]>([])
const loading = ref(false)

const fetchArticles = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/v1/help/articles')
    const result = await response.json()
    if (result.code === 0 && result.data?.items) {
      articles.value = result.data.items.map((item: any) => ({
        id: item.id,
        title: item.title,
        category: item.category_name || t('knowledgeBase.uncategorized'),
        summary: item.summary || item.content?.substring(0, 100) + '...',
        content: item.content,
        views: item.views || 0,
        helpful: item.helpful || 0,
        updateTime: item.updated_at?.split('T')[0] || ''
      }))
    }
  } catch (e) {
    console.error('Failed to fetch knowledge base articles:', e)
  } finally {
    loading.value = false
  }
}

fetchArticles()

const categoryTree = computed<TreeOption[]>(() => {
  const categories = [...new Set(articles.value.map(a => a.category))]
  return [
    {
      key: 'all',
      label: t('knowledgeBase.allArticles'),
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
