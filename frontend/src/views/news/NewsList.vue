<template>
  <div class="news-list-page">
    <SiteHeader />

    <!-- Hero Section -->
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">新闻资讯</h1>
        <p class="hero-desc">了解最新动态与行业资讯</p>
        <div class="search-wrapper">
          <el-input
            v-model="searchQuery"
            size="large"
            placeholder="搜索新闻..."
            :prefix-icon="Search"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
      </div>
    </section>

    <!-- Category Tabs -->
    <section class="category-tabs">
      <div class="container">
        <el-tabs v-model="activeCategory" @tab-change="handleCategoryChange">
          <el-tab-pane label="全部" name="all" />
          <el-tab-pane
            v-for="cat in categories"
            :key="cat.id"
            :label="cat.name"
            :name="String(cat.id)"
          />
        </el-tabs>
      </div>
    </section>

    <!-- News List -->
    <section class="section news-section">
      <div class="container">
        <el-row :gutter="24">
          <el-col :span="17">
            <div v-loading="loading">
              <div v-if="newsList.length > 0" class="news-grid">
                <div
                  v-for="item in newsList"
                  :key="item.id"
                  class="news-card"
                  @click="goToDetail(item.id)"
                >
                  <div class="news-cover" v-if="item.cover">
                    <img :src="item.cover" :alt="item.title" />
                  </div>
                  <div class="news-body">
                    <div class="news-meta">
                      <el-tag size="small" :type="getCategoryType(item.category_name)">
                        {{ item.category_name }}
                      </el-tag>
                      <span class="meta-date">
                        <el-icon><Calendar /></el-icon>
                        {{ item.created_at }}
                      </span>
                      <span class="meta-views">
                        <el-icon><View /></el-icon>
                        {{ item.views }}
                      </span>
                    </div>
                    <h3 class="news-title">{{ item.title }}</h3>
                    <p class="news-summary">{{ item.summary }}</p>
                    <span class="read-more">
                      阅读全文
                      <el-icon><ArrowRight /></el-icon>
                    </span>
                  </div>
                </div>
              </div>

              <el-empty v-else description="暂无新闻" />

              <div class="pagination-wrapper" v-if="total > pageSize">
                <el-pagination
                  v-model:current-page="currentPage"
                  v-model:page-size="pageSize"
                  :total="total"
                  :page-sizes="[10, 20, 30]"
                  layout="total, sizes, prev, pager, next, jumper"
                  @current-change="fetchNews"
                  @size-change="fetchNews"
                />
              </div>
            </div>
          </el-col>

          <!-- Sidebar -->
          <el-col :span="7">
            <div class="sidebar">
              <div class="sidebar-card">
                <h3 class="sidebar-title">热门新闻</h3>
                <div class="hot-list">
                  <div
                    v-for="(item, index) in hotNews"
                    :key="item.id"
                    class="hot-item"
                    @click="goToDetail(item.id)"
                  >
                    <span class="hot-rank" :class="{ 'top-3': index < 3 }">{{ index + 1 }}</span>
                    <span class="hot-title">{{ item.title }}</span>
                  </div>
                </div>
              </div>

              <div class="sidebar-card">
                <h3 class="sidebar-title">新闻分类</h3>
                <div class="category-list">
                  <div
                    v-for="cat in categories"
                    :key="cat.id"
                    class="category-item"
                    @click="goToCategory(cat.id)"
                  >
                    <span>{{ cat.name }}</span>
                    <el-badge :value="cat.count" :max="999" type="primary" />
                  </div>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Calendar, View, ArrowRight } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()
const searchQuery = ref('')
const activeCategory = ref('all')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

interface NewsItem {
  id: number
  title: string
  summary: string
  cover: string
  category_name: string
  created_at: string
  views: number
}

interface Category {
  id: number
  name: string
  count: number
}

const newsList = ref<NewsItem[]>([])
const hotNews = ref<NewsItem[]>([])
const categories = ref<Category[]>([])

const categoryTypeMap: Record<string, '' | 'success' | 'warning' | 'info' | 'danger'> = {
  '公司新闻': '',
  '行业动态': 'success',
  '产品更新': 'warning',
  '技术博客': 'info'
}

function getCategoryType(name: string) {
  return categoryTypeMap[name] || ''
}

const fetchNews = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (activeCategory.value !== 'all') {
      params.category_id = activeCategory.value
    }
    const { data } = await request.get('/api/v2/news', { params })
    if (data?.data) {
      newsList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取新闻列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchSidebar = async () => {
  try {
    const [hotRes, catRes] = await Promise.allSettled([
      request.get('/api/v2/news', { params: { limit: 8, sort: 'views' } }),
      request.get('/api/v2/news/categories')
    ])
    if (hotRes.status === 'fulfilled' && hotRes.value.data?.data) {
      hotNews.value = hotRes.value.data.data.list || hotRes.value.data.data
    }
    if (catRes.status === 'fulfilled' && catRes.value.data?.data) {
      categories.value = catRes.value.data.data
    }
  } catch (error) {
    console.error('获取侧边栏数据失败:', error)
  }
}

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push(`/news/search?q=${encodeURIComponent(searchQuery.value.trim())}`)
  }
}

const handleCategoryChange = () => {
  currentPage.value = 1
  fetchNews()
}

const goToDetail = (id: number) => {
  router.push(`/news/${id}`)
}

const goToCategory = (id: number) => {
  router.push(`/news/category/${id}`)
}

onMounted(() => {
  fetchNews()
  fetchSidebar()
})
</script>

<style scoped lang="scss">
.news-list-page {
  padding-top: 64px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }

  .hero-section {
    background: linear-gradient(135deg, #1a237e, #0d47a1);
    color: #fff;
    padding: 60px 0;
    text-align: center;

    .hero-title {
      font-size: 36px;
      font-weight: 700;
      margin: 0 0 12px;
    }

    .hero-desc {
      font-size: 16px;
      opacity: 0.8;
      margin: 0 0 32px;
    }

    .search-wrapper {
      max-width: 560px;
      margin: 0 auto;

      :deep(.el-input__wrapper) {
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }
    }
  }

  .category-tabs {
    background: #fff;
    border-bottom: 1px solid #e2e8f0;

    :deep(.el-tabs__header) {
      margin: 0;
    }

    :deep(.el-tabs__nav-wrap::after) {
      display: none;
    }
  }

  .section {
    padding: 32px 0 80px;
  }

  .news-section {
    background: #f5f7fa;
  }

  .news-grid {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .news-card {
    display: flex;
    background: #fff;
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
      transform: translateY(-2px);
    }

    .news-cover {
      width: 240px;
      flex-shrink: 0;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .news-body {
      flex: 1;
      padding: 20px 24px;
      display: flex;
      flex-direction: column;
    }

    .news-meta {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 12px;

      .meta-date,
      .meta-views {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #94a3b8;
      }
    }

    .news-title {
      font-size: 18px;
      font-weight: 600;
      color: #1a2332;
      margin: 0 0 8px;
      line-height: 1.4;
    }

    .news-summary {
      font-size: 14px;
      color: #64748b;
      line-height: 1.6;
      margin: 0;
      flex: 1;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .read-more {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      color: #409eff;
      font-size: 14px;
      font-weight: 500;
      margin-top: 12px;
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

  .hot-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .hot-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 10px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;

    &:hover {
      background: #f5f7fa;
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

      &.top-3 {
        background: linear-gradient(135deg, #409eff, #66b1ff);
        color: #fff;
      }
    }

    .hot-title {
      flex: 1;
      font-size: 14px;
      color: #4e5969;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .category-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .category-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 10px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;

    &:hover {
      background: #f5f7fa;
      color: #409eff;
    }

    span {
      font-size: 14px;
      color: #4e5969;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 32px;
  }
}

@media (max-width: 768px) {
  .news-list-page {
    .news-card {
      flex-direction: column;

      .news-cover {
        width: 100%;
        height: 180px;
      }
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
