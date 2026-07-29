<template>
  <div class="help-search-page">
    <SiteHeader />

    <!-- Search Header -->
    <section class="search-header">
      <div class="container">
        <h1 class="search-title">搜索帮助</h1>
        <p class="search-desc">输入关键词查找帮助文章</p>
        <div class="search-wrapper">
          <el-input
            v-model="searchQuery"
            size="large"
            placeholder="搜索问题或关键词..."
            :prefix-icon="Search"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
        <div class="search-tags" v-if="hotKeywords.length">
          <span class="tag-label">热门搜索：</span>
          <el-tag
            v-for="keyword in hotKeywords"
            :key="keyword"
            size="small"
            class="hot-tag"
            @click="searchByKeyword(keyword)"
          >
            {{ keyword }}
          </el-tag>
        </div>
      </div>
    </section>

    <!-- Results -->
    <section class="section results-section">
      <div class="container">
        <!-- Filter Bar -->
        <div class="filter-bar" v-if="searched">
          <div class="filter-left">
            <span class="results-info">
              共找到 <strong>{{ total }}</strong> 条关于 "<strong>{{ searchedQuery }}</strong>" 的结果
            </span>
          </div>
          <div class="filter-right">
            <el-select v-model="filterCategory" placeholder="全部分类" size="small" clearable style="width: 140px" @change="handleSearch">
              <el-option
                v-for="cat in categories"
                :key="cat.id"
                :label="cat.name"
                :value="cat.id"
              />
            </el-select>
            <el-select v-model="sortBy" size="small" style="width: 120px" @change="handleSearch">
              <el-option label="最相关" value="relevance" />
              <el-option label="最新" value="time" />
              <el-option label="最多浏览" value="views" />
            </el-select>
          </div>
        </div>

        <div v-loading="loading">
          <div v-if="articles.length > 0" class="article-list">
            <div
              v-for="article in articles"
              :key="article.id"
              class="article-card"
              @click="goToContent(article.id)"
            >
              <div class="article-icon">
                <el-icon :size="24"><Document /></el-icon>
              </div>
              <div class="article-body">
                <el-tag size="small" type="info" class="category-tag">{{ article.category_name }}</el-tag>
                <h3 class="article-title" v-html="highlightText(article.title)"></h3>
                <p class="article-summary" v-html="highlightText(article.summary)"></p>
                <div class="article-meta">
                  <span>
                    <el-icon><View /></el-icon>
                    {{ article.views }} 次阅读
                  </span>
                  <span>
                    <el-icon><Clock /></el-icon>
                    {{ article.updated_at }}
                  </span>
                  <span>
                    <el-icon><Star /></el-icon>
                    {{ article.helpful }} 人觉得有帮助
                  </span>
                </div>
              </div>
              <el-icon class="article-arrow"><ArrowRight /></el-icon>
            </div>
          </div>

          <el-empty v-else-if="searched" description="未找到相关帮助文章">
            <template #description>
              <p>未找到相关帮助文章</p>
              <p class="empty-suggest">请尝试更换关键词，或浏览帮助分类</p>
            </template>
            <el-button type="primary" @click="$router.push('/help')">浏览帮助中心</el-button>
          </el-empty>

          <div class="init-state" v-else>
            <el-icon :size="64" color="#c9cdd4"><Search /></el-icon>
            <p>输入关键词搜索帮助文章</p>
            <div class="quick-links">
              <h4>快速链接</h4>
              <div class="links-grid">
                <div
                  v-for="cat in categories"
                  :key="cat.id"
                  class="link-item"
                  @click="$router.push(`/help/category/${cat.id}`)"
                >
                  <el-icon><component :is="cat.icon || 'Folder'" /></el-icon>
                  <span>{{ cat.name }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="pagination-wrapper" v-if="total > pageSize">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[10, 20, 30]"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="fetchResults"
              @size-change="fetchResults"
            />
          </div>
        </div>
      </div>
    </section>

    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Search, Document, View, Clock, Star, ArrowRight } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const searchQuery = ref('')
const searchedQuery = ref('')
const searched = ref(false)
const loading = ref(false)
const sortBy = ref('relevance')
const filterCategory = ref<number | ''>('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const hotKeywords = ref(['如何注册', '购买产品', '续费', '退款', '工单', '发票'])

interface Article {
  id: number
  title: string
  summary: string
  category_name: string
  views: number
  helpful: number
  updated_at: string
}

interface Category {
  id: number
  name: string
  count: number
  icon: string
}

const articles = ref<Article[]>([])
const categories = ref<Category[]>([])

const fetchResults = async () => {
  if (!searchQuery.value.trim()) return

  loading.value = true
  searched.value = true
  searchedQuery.value = searchQuery.value.trim()

  try {
    const params: Record<string, any> = {
      q: searchedQuery.value,
      sort: sortBy.value,
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (filterCategory.value) {
      params.category_id = filterCategory.value
    }
    const { data } = await request.get('/api/v1/help/search', { params })
    if (data?.data) {
      articles.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('搜索帮助文章失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const { data } = await request.get('/api/v1/help/categories')
    if (data?.data) {
      categories.value = data.data
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  router.replace({ query: { q: searchQuery.value.trim() } })
  fetchResults()
}

const searchByKeyword = (keyword: string) => {
  searchQuery.value = keyword
  handleSearch()
}

const highlightText = (text: string) => {
  if (!searchedQuery.value) return text
  const regex = new RegExp(`(${searchedQuery.value})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

const goToContent = (id: number) => {
  router.push(`/help/content/${id}`)
}

onMounted(async () => {
  fetchCategories()
  // Fetch hot keywords
  try {
    const { data } = await request.get('/api/v1/help/articles/hot', { params: { limit: 6 } })
    if (data?.data?.list?.length) {
      hotKeywords.value = data.data.list.map((a: any) => a.title || a.keyword)
    }
  } catch (e) { /* keep defaults */ }
  const q = route.query.q as string
  if (q) {
    searchQuery.value = q
    fetchResults()
  }
})
</script>

<style scoped lang="scss">
.help-search-page {
  padding-top: 64px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }

  .search-header {
    background: linear-gradient(135deg, #1a237e, #0d47a1);
    color: #fff;
    padding: 48px 0;
    text-align: center;

    .search-title {
      font-size: 32px;
      font-weight: 700;
      margin: 0 0 8px;
    }

    .search-desc {
      font-size: 16px;
      opacity: 0.8;
      margin: 0 0 24px;
    }

    .search-wrapper {
      max-width: 600px;
      margin: 0 auto 16px;

      :deep(.el-input__wrapper) {
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }
    }

    .search-tags {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      flex-wrap: wrap;

      .tag-label {
        font-size: 14px;
        opacity: 0.8;
      }

      .hot-tag {
        cursor: pointer;
        background: rgba(255, 255, 255, 0.15);
        border: none;
        color: #fff;

        &:hover {
          background: rgba(255, 255, 255, 0.25);
        }
      }
    }
  }

  .section {
    padding: 32px 0 80px;
  }

  .results-section {
    background: #f5f7fa;
  }

  .filter-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 12px;

    .filter-left .results-info {
      font-size: 15px;
      color: #64748b;

      strong {
        color: #409eff;
      }
    }

    .filter-right {
      display: flex;
      gap: 8px;
    }
  }

  .article-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .article-card {
    display: flex;
    align-items: flex-start;
    gap: 16px;
    padding: 24px;
    background: #fff;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
      transform: translateY(-2px);
    }

    .article-icon {
      width: 48px;
      height: 48px;
      background: #ecf5ff;
      border-radius: 12px;
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

    .category-tag {
      margin-bottom: 8px;
    }

    .article-title {
      font-size: 18px;
      font-weight: 600;
      color: #1a2332;
      margin: 0 0 8px;

      :deep(mark) {
        background: #fef08a;
        color: #92400e;
        padding: 0 2px;
        border-radius: 2px;
      }
    }

    .article-summary {
      font-size: 14px;
      color: #64748b;
      line-height: 1.6;
      margin: 0 0 12px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;

      :deep(mark) {
        background: #fef08a;
        color: #92400e;
        padding: 0 2px;
        border-radius: 2px;
      }
    }

    .article-meta {
      display: flex;
      gap: 16px;

      span {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #94a3b8;
      }
    }

    .article-arrow {
      color: #c9cdd4;
      flex-shrink: 0;
      margin-top: 16px;
    }
  }

  .init-state {
    text-align: center;
    padding: 60px 0;

    > p {
      font-size: 16px;
      color: #94a3b8;
      margin: 16px 0 32px;
    }

    .quick-links {
      max-width: 600px;
      margin: 0 auto;

      h4 {
        font-size: 15px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 16px;
      }

      .links-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 12px;
      }

      .link-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 14px 16px;
        background: #fff;
        border-radius: 10px;
        cursor: pointer;
        transition: all 0.3s;
        border: 1px solid #e2e8f0;

        &:hover {
          border-color: #409eff;
          color: #409eff;
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.1);
        }

        .el-icon {
          color: #409eff;
        }

        span {
          font-size: 14px;
          color: #4e5969;
        }
      }
    }
  }

  .empty-suggest {
    font-size: 14px;
    color: #94a3b8;
    margin: 4px 0 0;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 32px;
  }
}

@media (max-width: 768px) {
  .help-search-page {
    .filter-bar {
      flex-direction: column;
      align-items: flex-start;
    }

    .article-card {
      flex-direction: column;

      .article-arrow {
        display: none;
      }
    }

    .init-state .links-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
