<template>
  <div class="help-content-page">
    <SiteHeader />

    <!-- Breadcrumb -->
    <section class="breadcrumb-section">
      <div class="container">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/help' }">帮助中心</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: `/help/category/${article.category_id}` }">
            {{ article.category_name }}
          </el-breadcrumb-item>
          <el-breadcrumb-item>{{ article.title || '文章详情' }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </section>

    <!-- Content -->
    <section class="section content-section">
      <div class="container">
        <el-row :gutter="32">
          <el-col :span="17">
            <el-skeleton :loading="loading" animated :rows="12">
              <template #default>
                <article class="article-main">
                  <div class="article-header">
                    <el-tag type="info">{{ article.category_name }}</el-tag>
                    <h1 class="article-title">{{ article.title }}</h1>
                    <div class="article-meta">
                      <span><el-icon><View /></el-icon> {{ article.views || 0 }} 次阅读</span>
                      <span><el-icon><Clock /></el-icon> {{ article.updated_at }}</span>
                    </div>
                  </div>
                  <div class="article-content" v-html="article.content"></div>

                  <!-- Tags -->
                  <div class="article-tags" v-if="article.tags?.length">
                    <span class="tag-label">标签：</span>
                    <el-tag v-for="tag in article.tags" :key="tag" size="small" effect="plain">
                      {{ tag }}
                    </el-tag>
                  </div>

                  <!-- Feedback -->
                  <div class="article-feedback">
                    <p>这篇文章对您有帮助吗？</p>
                    <div class="feedback-actions">
                      <el-button
                        :type="feedback === 'yes' ? 'success' : 'default'"
                        @click="submitFeedback('yes')"
                      >
                        <el-icon style="margin-right: 4px;"><CircleCheck /></el-icon>
                        有帮助 ({{ article.helpful || 0 }})
                      </el-button>
                      <el-button
                        :type="feedback === 'no' ? 'danger' : 'default'"
                        @click="submitFeedback('no')"
                      >
                        <el-icon style="margin-right: 4px;"><CircleClose /></el-icon>
                        没帮助
                      </el-button>
                    </div>
                  </div>

                  <!-- Navigation -->
                  <div class="article-nav">
                    <div class="nav-item" v-if="prevArticle" @click="goToContent(prevArticle.id)">
                      <span class="nav-label">上一篇</span>
                      <span class="nav-title">{{ prevArticle.title }}</span>
                    </div>
                    <div class="nav-item nav-next" v-if="nextArticle" @click="goToContent(nextArticle.id)">
                      <span class="nav-label">下一篇</span>
                      <span class="nav-title">{{ nextArticle.title }}</span>
                    </div>
                  </div>
                </article>
              </template>
            </el-skeleton>
          </el-col>

          <!-- Sidebar -->
          <el-col :span="7">
            <aside class="sidebar">
              <!-- Table of Contents -->
              <div class="sidebar-card toc-card" v-if="article.toc?.length">
                <h3 class="sidebar-title">目录</h3>
                <div class="toc-list">
                  <div
                    v-for="(item, index) in article.toc"
                    :key="index"
                    class="toc-item"
                    :class="{ 'level-2': item.level === 2, 'level-3': item.level === 3 }"
                  >
                    <a :href="`#${item.id}`">{{ item.text }}</a>
                  </div>
                </div>
              </div>

              <!-- Related Articles -->
              <div class="sidebar-card" v-if="relatedArticles.length > 0">
                <h3 class="sidebar-title">相关文章</h3>
                <div class="related-list">
                  <div
                    v-for="item in relatedArticles"
                    :key="item.id"
                    class="related-item"
                    @click="goToContent(item.id)"
                  >
                    <el-icon><Document /></el-icon>
                    <span>{{ item.title }}</span>
                  </div>
                </div>
              </div>

              <!-- Help Actions -->
              <div class="sidebar-card">
                <h3 class="sidebar-title">需要更多帮助？</h3>
                <div class="help-actions">
                  <el-button type="primary" plain style="width: 100%;" @click="$router.push('/user/tickets')">
                    提交工单
                  </el-button>
                  <el-button style="width: 100%;" @click="$router.push('/help')">
                    返回帮助中心
                  </el-button>
                </div>
              </div>
            </aside>
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
import { View, Clock, CircleCheck, CircleClose, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const feedback = ref<string | null>(null)

interface TocItem {
  id: string
  text: string
  level: number
}

interface Article {
  id: number
  title: string
  content: string
  category_id: number
  category_name: string
  views: number
  helpful: number
  updated_at: string
  tags: string[]
  toc: TocItem[]
}

const article = ref<Article>({
  id: 0,
  title: '',
  content: '',
  category_id: 0,
  category_name: '',
  views: 0,
  helpful: 0,
  updated_at: '',
  tags: [],
  toc: []
})

const prevArticle = ref<{ id: number; title: string } | null>(null)
const nextArticle = ref<{ id: number; title: string } | null>(null)
const relatedArticles = ref<{ id: number; title: string }[]>([])

const fetchArticle = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  feedback.value = null

  try {
    const [articleRes, relatedRes] = await Promise.allSettled([
      request.get(`/api/v2/help/articles/${id}`),
      request.get(`/api/v2/help/articles/${id}/related`)
    ])

    if (articleRes.status === 'fulfilled' && articleRes.value.data?.data) {
      const data = articleRes.value.data.data
      article.value = data.article || data
      prevArticle.value = data.prev || null
      nextArticle.value = data.next || null
    }
    if (relatedRes.status === 'fulfilled' && relatedRes.value.data?.data) {
      relatedArticles.value = relatedRes.value.data.data
    }
  } catch (error) {
    console.error('获取文章失败:', error)
  } finally {
    loading.value = false
  }
}

const submitFeedback = async (type: string) => {
  if (feedback.value) return
  feedback.value = type
  try {
    await request.post(`/api/v2/help/articles/${article.value.id}/feedback`, {
      helpful: type === 'yes'
    })
    ElMessage.success('感谢您的反馈！')
  } catch (error) {
    console.error('提交反馈失败:', error)
  }
}

const goToContent = (id: number) => {
  router.push(`/help/content/${id}`)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(() => route.params.id, fetchArticle)
onMounted(fetchArticle)
</script>

<style scoped lang="scss">
.help-content-page {
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

  .section {
    padding: 32px 0 80px;
  }

  .content-section {
    background: #f5f7fa;
  }

  .article-main {
    background: #fff;
    border-radius: 16px;
    padding: 40px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  }

  .article-header {
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid #e2e8f0;
  }

  .article-title {
    font-size: 28px;
    font-weight: 700;
    color: #0f172a;
    margin: 12px 0 16px;
    line-height: 1.4;
  }

  .article-meta {
    display: flex;
    gap: 20px;
    color: #94a3b8;
    font-size: 14px;

    span {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }

  .article-content {
    font-size: 16px;
    line-height: 1.8;
    color: #334155;

    :deep(h2) {
      font-size: 22px;
      font-weight: 600;
      color: #0f172a;
      margin: 32px 0 16px;
    }

    :deep(h3) {
      font-size: 18px;
      font-weight: 600;
      color: #0f172a;
      margin: 24px 0 12px;
    }

    :deep(p) {
      margin: 0 0 16px;
    }

    :deep(ul),
    :deep(ol) {
      padding-left: 24px;
      margin: 0 0 16px;
    }

    :deep(li) {
      margin-bottom: 8px;
    }

    :deep(code) {
      background: #f1f5f9;
      padding: 2px 6px;
      border-radius: 4px;
      font-size: 14px;
      color: #e11d48;
    }

    :deep(pre) {
      background: #1e293b;
      color: #e2e8f0;
      padding: 16px;
      border-radius: 8px;
      overflow-x: auto;
      margin: 16px 0;
    }

    :deep(img) {
      max-width: 100%;
      border-radius: 8px;
      margin: 16px 0;
    }
  }

  .article-tags {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 24px;
    padding-top: 20px;
    border-top: 1px solid #e2e8f0;

    .tag-label {
      font-size: 14px;
      color: #64748b;
    }
  }

  .article-feedback {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid #e2e8f0;
    text-align: center;

    p {
      font-size: 15px;
      color: #64748b;
      margin: 0 0 16px;
    }
  }

  .feedback-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
  }

  .article-nav {
    display: flex;
    justify-content: space-between;
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid #e2e8f0;
    gap: 16px;
  }

  .nav-item {
    flex: 1;
    padding: 16px;
    background: #f8fafc;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.3s;
    max-width: 48%;

    &:hover {
      background: #ecf5ff;
    }

    &.nav-next {
      text-align: right;
    }

    .nav-label {
      display: block;
      font-size: 12px;
      color: #94a3b8;
      margin-bottom: 6px;
    }

    .nav-title {
      font-size: 14px;
      color: #334155;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
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
      color: #0f172a;
      margin: 0 0 16px;
      padding-bottom: 12px;
      border-bottom: 2px solid #409eff;
    }
  }

  .toc-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .toc-item {
    a {
      display: block;
      padding: 6px 10px;
      font-size: 14px;
      color: #4e5969;
      text-decoration: none;
      border-radius: 6px;
      transition: all 0.2s;

      &:hover {
        background: #f5f7fa;
        color: #409eff;
      }
    }

    &.level-2 a {
      padding-left: 10px;
      font-weight: 500;
    }

    &.level-3 a {
      padding-left: 24px;
      font-size: 13px;
    }
  }

  .related-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .related-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    background: #f8fafc;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    font-size: 14px;
    color: #334155;
    border: 1px solid #e2e8f0;

    &:hover {
      color: #409eff;
      border-color: #409eff;
      background: #ecf5ff;
    }

    .el-icon {
      flex-shrink: 0;
      color: #94a3b8;
    }

    span {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .help-actions {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
}

@media (max-width: 768px) {
  .help-content-page {
    .article-main {
      padding: 20px;
    }

    :deep(.el-col-17),
    :deep(.el-col-7) {
      width: 100%;
    }

    .sidebar {
      margin-top: 24px;
    }

    .article-nav {
      flex-direction: column;
    }

    .nav-item {
      max-width: 100%;
    }
  }
}
</style>
