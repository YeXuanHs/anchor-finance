<template>
  <div class="news-search-page">
    <SiteHeader />

    <!-- Search Header -->
    <section class="search-header">
      <div class="container">
        <h1 class="search-title">新闻搜索</h1>
        <div class="search-wrapper">
          <el-input
            v-model="searchQuery"
            size="large"
            placeholder="输入关键词搜索新闻..."
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
        <div class="results-header" v-if="searched">
          <span class="results-info">
            共找到 <strong>{{ total }}</strong> 条关于 "<strong>{{ searchedQuery }}</strong>" 的新闻
          </span>
          <el-select v-model="sortBy" size="small" style="width: 120px" @change="handleSearch">
            <el-option label="最新发布" value="time" />
            <el-option label="最多浏览" value="views" />
            <el-option label="最多点赞" value="likes" />
          </el-select>
        </div>

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
                  <el-tag size="small" type="info">{{ item.category_name }}</el-tag>
                  <span class="meta-date">
                    <el-icon><Calendar /></el-icon>
                    {{ item.created_at }}
                  </span>
                  <span class="meta-views">
                    <el-icon><View /></el-icon>
                    {{ item.views }}
                  </span>
                </div>
                <h3 class="news-title" v-html="highlightText(item.title)"></h3>
                <p class="news-summary" v-html="highlightText(item.summary)"></p>
              </div>
            </div>
          </div>

          <el-empty v-else-if="searched" description="未找到相关新闻" />

          <div class="init-state" v-else>
            <el-icon :size="64" color="#c9cdd4"><Search /></el-icon>
            <p>输入关键词搜索新闻</p>
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
import { Search, Calendar, View } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const searchQuery = ref('')
const searchedQuery = ref('')
const searched = ref(false)
const loading = ref(false)
const sortBy = ref('time')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const hotKeywords = ref(['系统更新', '安全公告', '新功能', '维护通知', '优惠活动'])

interface NewsItem {
  id: number
  title: string
  summary: string
  cover: string
  category_name: string
  created_at: string
  views: number
}

const newsList = ref<NewsItem[]>([])

const fetchResults = async () => {
  if (!searchQuery.value.trim()) return

  loading.value = true
  searched.value = true
  searchedQuery.value = searchQuery.value.trim()

  try {
    const { data } = await request.get('/api/v1/news', {
      params: {
        q: searchedQuery.value,
        sort: sortBy.value,
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    if (data?.data) {
      newsList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('搜索新闻失败:', error)
  } finally {
    loading.value = false
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

const goToDetail = (id: number) => {
  router.push(`/news/${id}`)
}

onMounted(async () => {
  // Fetch hot news keywords
  try {
    const { data } = await request.get('/api/v1/news', { params: { limit: 5, sort: 'views' } })
    if (data?.data?.list?.length) {
      hotKeywords.value = data.data.list.map((n: any) => n.title)
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
.news-search-page {
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

  .results-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .results-info {
      font-size: 15px;
      color: #64748b;

      strong {
        color: #409eff;
      }
    }
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

      :deep(mark) {
        background: #fef08a;
        color: #92400e;
        padding: 0 2px;
        border-radius: 2px;
      }
    }

    .news-summary {
      font-size: 14px;
      color: #64748b;
      line-height: 1.6;
      margin: 0;
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
  }

  .init-state {
    text-align: center;
    padding: 80px 0;

    p {
      font-size: 16px;
      color: #94a3b8;
      margin: 16px 0 0;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 32px;
  }
}

@media (max-width: 768px) {
  .news-search-page {
    .news-card {
      flex-direction: column;

      .news-cover {
        width: 100%;
        height: 180px;
      }
    }

    .results-header {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;
    }
  }
}
</style>
