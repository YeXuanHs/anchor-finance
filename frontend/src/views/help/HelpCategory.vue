<template>
  <div class="help-category-page">
    <SiteHeader />

    <!-- Breadcrumb -->
    <section class="breadcrumb-section">
      <div class="container">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">{{ $t('common.home') }}</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/help' }">{{ $t('helpCenter.title') }}</el-breadcrumb-item>
          <el-breadcrumb-item>{{ categoryName }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </section>

    <!-- Category Header -->
    <section class="category-header">
      <div class="container">
        <div class="header-content">
          <div class="header-icon">
            <el-icon :size="40"><component :is="categoryIcon" /></el-icon>
          </div>
          <div>
            <h1 class="category-title">{{ categoryName }}</h1>
            <p class="category-desc">{{ categoryDesc }}</p>
          </div>
        </div>
        <div class="search-wrapper">
          <el-input
            v-model="searchQuery"
            size="large"
            :placeholder="$t('helpCommon.searchInCategory')"
            :prefix-icon="Search"
            clearable
            @keyup.enter="handleSearch"
          />
        </div>
      </div>
    </section>

    <!-- Content -->
    <section class="section content-section">
      <div class="container">
        <el-row :gutter="24">
          <el-col :span="17">
            <div v-loading="loading">
              <!-- Sub-categories -->
              <div v-if="subCategories.length > 0" class="sub-categories">
                <div
                  v-for="sub in subCategories"
                  :key="sub.id"
                  class="sub-category-card"
                  @click="filterBySubCategory(sub.id)"
                >
                  <el-icon :size="24" color="#409eff"><Folder /></el-icon>
                  <div class="sub-info">
                    <h3>{{ sub.name }}</h3>
                    <span>{{ sub.count }} {{ $t('helpCommon.articles') }}</span>
                  </div>
                </div>
              </div>

              <!-- Article List -->
              <div class="article-list">
                <div
                  v-for="article in articles"
                  :key="article.id"
                  class="article-item"
                  @click="goToContent(article.id)"
                >
                  <div class="article-icon">
                    <el-icon><Document /></el-icon>
                  </div>
                  <div class="article-body">
                    <h3 class="article-title">{{ article.title }}</h3>
                    <p class="article-summary">{{ article.summary }}</p>
                    <div class="article-meta">
                      <span>
                        <el-icon><View /></el-icon>
                        {{ article.views }} {{ $t('helpCommon.readCount') }}
                      </span>
                      <span>
                        <el-icon><Clock /></el-icon>
                        {{ article.updated_at }}
                      </span>
                    </div>
                  </div>
                  <el-icon class="article-arrow"><ArrowRight /></el-icon>
                </div>
              </div>

              <el-empty v-if="!loading && articles.length === 0" :description="$t('helpCommon.noArticlesInCategory')" />

              <div class="pagination-wrapper" v-if="total > pageSize">
                <el-pagination
                  v-model:current-page="currentPage"
                  v-model:page-size="pageSize"
                  :total="total"
                  :page-sizes="[10, 20, 30]"
                  layout="total, sizes, prev, pager, next, jumper"
                  @current-change="fetchArticles"
                  @size-change="fetchArticles"
                />
              </div>
            </div>
          </el-col>

          <!-- Sidebar -->
          <el-col :span="7">
            <div class="sidebar">
              <div class="sidebar-card">
                <h3 class="sidebar-title">{{ $t('helpCommon.helpCategory') }}</h3>
                <div class="category-nav">
                  <div
                    v-for="cat in allCategories"
                    :key="cat.id"
                    class="nav-item"
                    :class="{ active: cat.id === Number(categoryId) }"
                    @click="goToCategory(cat.id)"
                  >
                    <el-icon><component :is="cat.icon || 'Folder'" /></el-icon>
                    <span>{{ cat.name }}</span>
                    <span class="count">{{ cat.count }}</span>
                  </div>
                </div>
              </div>

              <div class="sidebar-card">
                <h3 class="sidebar-title">{{ $t('helpCommon.needMoreHelp') }}</h3>
                <div class="help-actions">
                  <el-button type="primary" plain @click="$router.push('/user/tickets')">
                    <el-icon><Tickets /></el-icon>
                    {{ $t('helpCommon.submitTicket') }}
                  </el-button>
                  <el-button @click="$router.push('/help/search')">
                    <el-icon><Search /></el-icon>
                    {{ $t('helpCommon.searchHelp') }}
                  </el-button>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>

    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  Search, Document, View, Clock, ArrowRight, Folder, Tickets,
  ShoppingCart, Wallet, Setting, Connection
} from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const searchQuery = ref('')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeSubCategoryId = ref<number | null>(null)

const categoryId = computed(() => route.params.id as string)

interface Article {
  id: number
  title: string
  summary: string
  views: number
  updated_at: string
}

interface Category {
  id: number
  name: string
  count: number
  icon: string
  description: string
}

interface SubCategory {
  id: number
  name: string
  count: number
}

const articles = ref<Article[]>([])
const allCategories = ref<Category[]>([])
const subCategories = ref<SubCategory[]>([])

const currentCategory = computed(() =>
  allCategories.value.find(c => c.id === Number(categoryId.value))
)

const categoryName = computed(() => currentCategory.value?.name || t('helpCommon.helpCategory'))
const categoryDesc = computed(() => currentCategory.value?.description || '')
const categoryIcon = computed(() => currentCategory.value?.icon || 'Folder')

const fetchArticles = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      category_id: categoryId.value,
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (activeSubCategoryId.value) {
      params.sub_category_id = activeSubCategoryId.value
    }
    if (searchQuery.value.trim()) {
      params.q = searchQuery.value.trim()
    }
    const { data } = await request.get('/api/v1/help/articles', { params })
    if (data?.data) {
      articles.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取帮助文章失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const { data } = await request.get('/api/v1/help/categories')
    if (data?.data) {
      allCategories.value = data.data
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const fetchSubCategories = async () => {
  try {
    const { data } = await request.get(`/api/v1/help/categories/${categoryId.value}/sub`)
    if (data?.data) {
      subCategories.value = data.data
    }
  } catch (error) {
    console.error('获取子分类失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  activeSubCategoryId.value = null
  fetchArticles()
}

const filterBySubCategory = (id: number) => {
  activeSubCategoryId.value = id
  currentPage.value = 1
  fetchArticles()
}

const goToContent = (id: number) => {
  router.push(`/help/content/${id}`)
}

const goToCategory = (id: number) => {
  router.push(`/help/category/${id}`)
}

watch(() => route.params.id, () => {
  currentPage.value = 1
  activeSubCategoryId.value = null
  searchQuery.value = ''
  fetchArticles()
  fetchSubCategories()
})

onMounted(() => {
  fetchCategories()
  fetchSubCategories()
  fetchArticles()
})
</script>

<style scoped lang="scss">
.help-category-page {
  padding-top: 64px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }

  .breadcrumb-section {
    background: #fff;
    padding: 16px 0;
    border-bottom: 1px solid #e2e8f0;
  }

  .category-header {
    background: linear-gradient(135deg, #1a237e, #0d47a1);
    color: #fff;
    padding: 40px 0;

    .header-content {
      display: flex;
      align-items: center;
      gap: 20px;
      margin-bottom: 24px;

      .header-icon {
        width: 72px;
        height: 72px;
        background: rgba(255, 255, 255, 0.15);
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .category-title {
      font-size: 32px;
      font-weight: 700;
      margin: 0 0 8px;
    }

    .category-desc {
      font-size: 16px;
      opacity: 0.8;
      margin: 0;
    }

    .search-wrapper {
      max-width: 500px;

      :deep(.el-input__wrapper) {
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }
    }
  }

  .section {
    padding: 32px 0 80px;
  }

  .content-section {
    background: #f5f7fa;
  }

  .sub-categories {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    margin-bottom: 24px;
  }

  .sub-category-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: #fff;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
      transform: translateY(-2px);
    }

    .sub-info {
      h3 {
        font-size: 16px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 4px;
      }

      span {
        font-size: 13px;
        color: #94a3b8;
      }
    }
  }

  .article-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .article-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: #fff;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
      transform: translateX(4px);
    }

    .article-icon {
      width: 44px;
      height: 44px;
      background: #ecf5ff;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #409eff;
      flex-shrink: 0;
    }

    .article-body {
      flex: 1;
      min-width: 0;
    }

    .article-title {
      font-size: 16px;
      font-weight: 600;
      color: #1a2332;
      margin: 0 0 6px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .article-summary {
      font-size: 13px;
      color: #64748b;
      margin: 0 0 8px;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .article-meta {
      display: flex;
      gap: 16px;

      span {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12px;
        color: #94a3b8;
      }
    }

    .article-arrow {
      color: #c9cdd4;
      flex-shrink: 0;
    }
  }

  .sidebar {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .sidebar-card {
    background: #fff;
    border-radius: 12px;
    padding: 20px;

    .sidebar-title {
      font-size: 16px;
      font-weight: 600;
      color: #1a2332;
      margin: 0 0 16px;
      padding-bottom: 12px;
      border-bottom: 2px solid #409eff;
    }
  }

  .category-nav {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover,
    &.active {
      background: #ecf5ff;
      color: #409eff;
    }

    span {
      flex: 1;
      font-size: 14px;
      color: #4e5969;
    }

    .count {
      flex: none;
      font-size: 12px;
      color: #94a3b8;
      background: #f5f7fa;
      padding: 2px 8px;
      border-radius: 10px;
    }
  }

  .help-actions {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 32px;
  }
}

@media (max-width: 768px) {
  .help-category-page {
    .sub-categories {
      grid-template-columns: 1fr;
    }

    :deep(.el-col-17),
    :deep(.el-col-7) {
      width: 100%;
    }

    .sidebar {
      margin-top: 24px;
    }
  }
}
</style>
