<template>
  <div class="news-detail-page">
    <SiteHeader />
    
    <section class="content-section">
      <div class="container">
        <el-row :gutter="40">
          <!-- Main Content -->
          <el-col :span="16">
            <div class="article-wrapper" v-loading="loading">
              <div class="article-header">
                <div class="article-meta">
                  <span class="meta-item">
                    <el-icon><Calendar /></el-icon>
                    {{ article.created_at }}
                  </span>
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ article.views }} 次阅读
                  </span>
                  <span class="meta-item">
                    <el-icon><Folder /></el-icon>
                    {{ article.category }}
                  </span>
                </div>
                <h1 class="article-title">{{ article.title }}</h1>
              </div>
              
              <div class="article-content" v-html="article.content"></div>
              
              <!-- Tags -->
              <div class="article-tags" v-if="article.tags?.length">
                <span class="tag-label">标签：</span>
                <el-tag v-for="tag in article.tags" :key="tag" size="small">{{ tag }}</el-tag>
              </div>
              
              <!-- Share -->
              <div class="article-share">
                <span>分享到：</span>
                <el-button circle size="small">
                  <el-icon><Share /></el-icon>
                </el-button>
              </div>
              
              <!-- Navigation -->
              <div class="article-nav">
                <div class="nav-prev" v-if="prevArticle" @click="goToArticle(prevArticle.id)">
                  <span class="nav-label">上一篇</span>
                  <span class="nav-title">{{ prevArticle.title }}</span>
                </div>
                <div class="nav-next" v-if="nextArticle" @click="goToArticle(nextArticle.id)">
                  <span class="nav-label">下一篇</span>
                  <span class="nav-title">{{ nextArticle.title }}</span>
                </div>
              </div>
            </div>
          </el-col>
          
          <!-- Sidebar -->
          <el-col :span="8">
            <div class="sidebar">
              <!-- Recent News -->
              <div class="sidebar-card">
                <h3 class="sidebar-title">最新文章</h3>
                <div class="recent-list">
                  <div class="recent-item" v-for="item in recentNews" :key="item.id" @click="goToArticle(item.id)">
                    <h4>{{ item.title }}</h4>
                    <span>{{ item.created_at }}</span>
                  </div>
                </div>
              </div>
              
              <!-- Categories -->
              <div class="sidebar-card">
                <h3 class="sidebar-title">文章分类</h3>
                <div class="category-list">
                  <div class="category-item" v-for="cat in categories" :key="cat.id" @click="goToCategory(cat.id)">
                    <span>{{ cat.name }}</span>
                    <span class="count">{{ cat.count }}</span>
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
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Calendar, View, Folder, Share } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const article = ref<any>({})
const prevArticle = ref<any>(null)
const nextArticle = ref<any>(null)
const recentNews = ref<any[]>([])
const categories = ref<any[]>([])

const fetchArticle = async (id: string) => {
  loading.value = true
  try {
    const { data } = await request.get(`/api/v1/news/${id}`)
    if (data?.data) {
      article.value = data.data.article || {}
      prevArticle.value = data.data.prev || null
      nextArticle.value = data.data.next || null
    }
  } catch (error) {
    console.error('获取文章失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchSidebar = async () => {
  try {
    const [recentRes, catRes] = await Promise.all([
      request.get('/api/v1/news?limit=5'),
      request.get('/api/v1/news/categories')
    ])
    
    if (recentRes.data?.data) recentNews.value = recentRes.data.data
    if (catRes.data?.data) categories.value = catRes.data.data
  } catch (error) {
    console.error('获取侧边栏数据失败:', error)
  }
}

const goToArticle = (id: number) => {
  router.push(`/news/${id}`)
}

const goToCategory = (id: number) => {
  router.push(`/news?category=${id}`)
}

watch(() => route.params.id, (newId) => {
  if (newId) fetchArticle(newId as string)
})

onMounted(() => {
  const id = route.params.id as string
  if (id) fetchArticle(id)
  fetchSidebar()
})
</script>

<style scoped lang="scss">
.news-detail-page {
  padding-top: 64px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .content-section {
    padding: 40px 0 80px;
    background: #f5f7fa;
  }
  
  .article-wrapper {
    background: #fff;
    border-radius: 12px;
    padding: 40px;
    
    .article-header {
      margin-bottom: 30px;
      padding-bottom: 20px;
      border-bottom: 1px solid #eee;
      
      .article-meta {
        display: flex;
        gap: 20px;
        margin-bottom: 16px;
        
        .meta-item {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: #909399;
        }
      }
      
      .article-title {
        font-size: 28px;
        font-weight: 700;
        color: #1a2332;
        margin: 0;
        line-height: 1.4;
      }
    }
    
    .article-content {
      font-size: 15px;
      line-height: 1.8;
      color: #333;
      
      :deep(img) {
        max-width: 100%;
        border-radius: 8px;
      }
    }
    
    .article-tags {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-top: 30px;
      padding-top: 20px;
      border-top: 1px solid #eee;
      
      .tag-label {
        font-size: 14px;
        color: #666;
      }
    }
    
    .article-share {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-top: 20px;
      
      span {
        font-size: 14px;
        color: #666;
      }
    }
    
    .article-nav {
      display: flex;
      justify-content: space-between;
      margin-top: 40px;
      padding-top: 20px;
      border-top: 1px solid #eee;
      
      .nav-prev,
      .nav-next {
        cursor: pointer;
        padding: 16px;
        background: #f5f7fa;
        border-radius: 8px;
        max-width: 45%;
        
        &:hover {
          background: #ecf5ff;
        }
        
        .nav-label {
          display: block;
          font-size: 12px;
          color: #909399;
          margin-bottom: 4px;
        }
        
        .nav-title {
          font-size: 14px;
          color: #333;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }
      }
      
      .nav-next {
        text-align: right;
      }
    }
  }
  
  .sidebar {
    .sidebar-card {
      background: #fff;
      border-radius: 12px;
      padding: 24px;
      margin-bottom: 20px;
      
      .sidebar-title {
        font-size: 16px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 16px;
        padding-bottom: 12px;
        border-bottom: 2px solid #409eff;
      }
      
      .recent-list {
        .recent-item {
          padding: 12px 0;
          border-bottom: 1px solid #f0f0f0;
          cursor: pointer;
          
          &:hover h4 {
            color: #409eff;
          }
          
          &:last-child {
            border-bottom: none;
          }
          
          h4 {
            font-size: 14px;
            font-weight: 500;
            color: #333;
            margin: 0 0 4px;
            transition: color 0.3s;
          }
          
          span {
            font-size: 12px;
            color: #909399;
          }
        }
      }
      
      .category-list {
        .category-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 10px 0;
          border-bottom: 1px solid #f0f0f0;
          cursor: pointer;
          
          &:hover {
            color: #409eff;
          }
          
          &:last-child {
            border-bottom: none;
          }
          
          .count {
            background: #f5f7fa;
            padding: 2px 8px;
            border-radius: 10px;
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .news-detail-page {
    .el-col-16,
    .el-col-8 {
      width: 100%;
    }
    
    .article-wrapper {
      padding: 20px;
      
      .article-header .article-title {
        font-size: 22px;
      }
    }
  }
}
</style>
