<template>
  <div class="help-article-page">
    <!-- Breadcrumb -->
    <section class="breadcrumb-section">
      <div class="container">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">{{ $t('common.home') }}</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/help' }">{{ $t('helpCenter.title') }}</el-breadcrumb-item>
          <el-breadcrumb-item>{{ article.title || $t('helpCommon.articleDetail') }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </section>

    <!-- 文章内容 -->
    <section class="section article-section">
      <div class="container">
        <div class="article-layout">
          <el-skeleton :loading="loading" animated :rows="12">
            <template #default>
              <article class="article-main">
                <div class="article-header">
                  <el-tag type="info">{{ article.category_name }}</el-tag>
                  <h1 class="article-title">{{ article.title }}</h1>
                  <div class="article-meta">
                    <span><el-icon><View /></el-icon> {{ article.views || 0 }} {{ $t('helpCommon.readCount') }}</span>
                    <span><el-icon><Clock /></el-icon> {{ article.updated_at || '2025-01-01' }}</span>
                  </div>
                </div>
                <div class="article-content" v-html="article.content"></div>

                <!-- 有用评价 -->
                <div class="article-feedback">
                  <p>{{ $t('helpCommon.isHelpful') }}</p>
                  <div class="feedback-actions">
                    <el-button :type="feedback === 'yes' ? 'success' : 'default'" @click="submitFeedback('yes')">
                      <el-icon style="margin-right: 4px;"><CircleCheck /></el-icon>
                      {{ $t('helpCommon.helpfulYes') }} ({{ article.helpful || 0 }})
                    </el-button>
                    <el-button :type="feedback === 'no' ? 'danger' : 'default'" @click="submitFeedback('no')">
                      <el-icon style="margin-right: 4px;"><CircleClose /></el-icon>
                      {{ $t('helpCommon.helpfulNo') }}
                    </el-button>
                  </div>
                </div>
              </article>
            </template>
          </el-skeleton>

          <!-- 相关文章 -->
          <aside class="article-sidebar" v-if="relatedArticles.length > 0">
            <h3 class="sidebar-title">{{ $t('helpCommon.relatedArticles') }}</h3>
            <div class="related-list">
              <div
                v-for="item in relatedArticles"
                :key="item.id"
                class="related-item"
                @click="$router.push(`/help/article/${item.id}`)"
              >
                <el-icon><Document /></el-icon>
                <span>{{ item.title }}</span>
              </div>
            </div>
          </aside>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { View, Clock, CircleCheck, CircleClose, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const { t } = useI18n()
const route = useRoute()
const loading = ref(true)
const feedback = ref<string | null>(null)

interface Article {
  id: number
  title: string
  content: string
  category_name: string
  views: number
  helpful: number
  updated_at: string
}

const article = ref<Article>({
  id: 0,
  title: '',
  content: '',
  category_name: '',
  views: 0,
  helpful: 0,
  updated_at: ''
})

const relatedArticles = ref<{ id: number; title: string }[]>([])

const fetchArticle = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const [articleRes, relatedRes] = await Promise.allSettled([
      request.get(`/api/v1/help/articles/${id}`),
      request.get(`/api/v1/help/articles/${id}/related`)
    ])

    if (articleRes.status === 'fulfilled' && articleRes.value.data?.data) {
      article.value = articleRes.value.data.data
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
    await request.post(`/api/v1/help/articles/${article.value.id}/feedback`, { helpful: type === 'yes' })
    ElMessage.success(t('helpCenter.thankFeedback'))
  } catch (error) {
    console.error('提交反馈失败:', error)
  }
}

watch(() => route.params.id, fetchArticle)
onMounted(fetchArticle)
</script>

<style scoped lang="scss">
.help-article-page {
  min-height: 100vh;
  background: #f8fafc;
}

.breadcrumb-section {
  background: #fff;
  padding: 80px 0 16px;
  border-bottom: 1px solid #e2e8f0;
}

.section {
  padding: 40px 0 80px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.article-layout {
  display: grid;
  grid-template-columns: 1fr 280px;
  gap: 32px;

  @media (max-width: 1024px) {
    grid-template-columns: 1fr;
  }
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

  :deep(ul), :deep(ol) {
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

.article-feedback {
  margin-top: 40px;
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

// 侧边栏
.sidebar-title {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
  margin: 0 0 16px;
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
  padding: 12px;
  background: #fff;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
  color: #334155;
  border: 1px solid #e2e8f0;

  &:hover {
    color: #1a56db;
    border-color: #1a56db;
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
</style>
