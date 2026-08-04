<template>
  <div class="news-category-page">
    <SiteHeader />

    <!-- Breadcrumb -->
    <section class="breadcrumb-section">
      <div class="container">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/news' }">新闻资讯</el-breadcrumb-item>
          <el-breadcrumb-item>{{ categoryName }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </section>

    <!-- Category Header -->
    <section class="category-header">
      <div class="container">
        <h1 class="category-title">{{ categoryName }}</h1>
        <p class="category-desc">{{ categoryDesc }}</p>
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

              <el-empty v-else description="该分类暂无新闻" />

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
                <h3 class="sidebar-title">新闻分类</h3>
                <div class="category-list">
                  <div
                    v-for="cat in categories"
                    :key="cat.id"
                    class="category-item"
                    :class="{ active: cat.id === Number(categoryId) }"
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Calendar, View, ArrowRight } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const categoryId = computed(() => route.params.id as string)

interface NewsItem {
  id: number
  title: string
  summary: string
  cover: string
  created_at: string
  views: number
}

interface Category {
  id: number
  name: string
  count: number
  description: string
}

const newsList = ref<NewsItem[]>([])
const categories = ref<Category[]>([])

const currentCategory = computed(() =>
  categories.value.find(c => c.id === Number(categoryId.value))
)

const categoryName = computed(() => currentCategory.value?.name || '新闻分类')
const categoryDesc = computed(() => currentCategory.value?.description || '')

const fetchNews = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/news', {
      params: {
        category_id: categoryId.value,
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    if (data?.data) {
      newsList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取分类新闻失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const { data } = await request.get('/api/v2/news/categories')
    if (data?.data) {
      categories.value = data.data
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const goToDetail = (id: number) => {
  router.push(`/news/${id}`)
}

const goToCategory = (id: number) => {
  router.push(`/news/category/${id}`)
}

watch(() => route.params.id, () => {
  currentPage.value = 1
  fetchNews()
})

onMounted(() => {
  fetchCategories()
  fetchNews()
})
</script>

<style scoped lang="scss">
.news-category-page {
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
    text-align: center;

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
      width: 220px;
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
      gap: 16px;
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

  .category-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .category-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
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
  .news-category-page {
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
