<template>
  <div class="help-cate-page">
    <SiteHeader />
    <div class="hero-section">
      <div class="container">
        <h1>帮助中心</h1>
        <p>浏览帮助分类，快速找到您需要的答案</p>
        <div class="search-box">
          <el-input v-model="searchKeyword" placeholder="搜索帮助文档..." size="large" clearable @keyup.enter="handleSearch">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
            <template #append>
              <el-button @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
      </div>
    </div>

    <div class="container content-section">
      <el-row :gutter="24">
        <el-col :span="16">
          <div class="categories-grid">
            <div v-for="cat in categories" :key="cat.id" class="category-card" @click="goCategory(cat.id)">
              <div class="cat-icon" :style="{ background: cat.color }">
                <el-icon :size="32"><component :is="cat.icon" /></el-icon>
              </div>
              <div class="cat-info">
                <h3>{{ cat.name }}</h3>
                <p>{{ cat.description }}</p>
                <span class="article-count">{{ cat.article_count }} 篇文章</span>
              </div>
              <el-icon class="arrow"><ArrowRight /></el-icon>
            </div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="sidebar">
            <div class="sidebar-card">
              <h3>热门问题</h3>
              <div v-for="item in hotQuestions" :key="item.id" class="hot-item" @click="goArticle(item.id)">
                <span class="hot-rank" :class="{ top: item.rank <= 3 }">{{ item.rank }}</span>
                <span class="hot-title">{{ item.title }}</span>
              </div>
            </div>
            <div class="sidebar-card">
              <h3>常见标签</h3>
              <div class="tags">
                <el-tag v-for="tag in tags" :key="tag" class="tag-item" @click="searchByTag(tag)">{{ tag }}</el-tag>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, ArrowRight, Document, Setting, Monitor, CreditCard, ChatDotRound, Star } from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()
const searchKeyword = ref('')

const categories = ref<any[]>([])
const hotQuestions = ref<any[]>([])
const tags = ref<string[]>(['服务器', '域名', '支付', '退款', '续费', '备案', 'SSL', 'CDN', '安全', '工单'])

const fetchCategories = async () => {
  try {
    const res = await request.get('/api/v1/help/categories')
    categories.value = res.data?.data || [
      { id: 1, name: '服务器相关', description: '云服务器、VPS、独立服务器常见问题', icon: 'Monitor', color: '#409eff', article_count: 45 },
      { id: 2, name: '域名相关', description: '域名注册、解析、备案等问题', icon: 'Document', color: '#67c23a', article_count: 32 },
      { id: 3, name: '支付与账单', description: '充值、支付、退款、发票相关', icon: 'CreditCard', color: '#e6a23c', article_count: 28 },
      { id: 4, name: '账户管理', description: '注册、登录、安全设置等', icon: 'Setting', color: '#909399', article_count: 20 },
      { id: 5, name: '工单与售后', description: '提交工单、售后服务', icon: 'ChatDotRound', color: '#f56c6c', article_count: 15 },
      { id: 6, name: '产品使用', description: '产品使用教程和技巧', icon: 'Star', color: '#00d4ff', article_count: 38 },
    ]
  } catch { /* use defaults */ }
}

const fetchHotQuestions = async () => {
  try {
    const res = await request.get('/api/v1/help/hot')
    hotQuestions.value = res.data?.data || [
      { id: 1, title: '如何重置服务器密码？', rank: 1 },
      { id: 2, title: '域名解析如何设置？', rank: 2 },
      { id: 3, title: '如何申请退款？', rank: 3 },
      { id: 4, title: '服务器无法连接怎么办？', rank: 4 },
      { id: 5, title: '如何续费产品？', rank: 5 },
    ]
  } catch { /* use defaults */ }
}

const goCategory = (id: number) => router.push(`/help/category/${id}`)
const goArticle = (id: number) => router.push(`/help/article/${id}`)
const handleSearch = () => { if (searchKeyword.value) router.push({ path: '/help/search', query: { q: searchKeyword.value } }) }
const searchByTag = (tag: string) => router.push({ path: '/help/search', query: { q: tag } })

onMounted(() => { fetchCategories(); fetchHotQuestions() })
</script>

<style scoped lang="scss">
.help-cate-page { background: #f5f7fa; min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a3a5c 0%, #0d2137 100%); padding: 60px 0; text-align: center; color: #fff;
  h1 { font-size: 36px; margin-bottom: 12px; }
  p { font-size: 16px; color: rgba(255,255,255,0.7); margin-bottom: 32px; }
  .search-box { max-width: 600px; margin: 0 auto; }
}
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.content-section { padding: 40px 20px; }
.categories-grid { display: flex; flex-direction: column; gap: 16px; }
.category-card { background: #fff; border-radius: 12px; padding: 24px; display: flex; align-items: center; gap: 20px; cursor: pointer; transition: all 0.3s; border: 1px solid #e4e7ed;
  &:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.1); transform: translateY(-2px); }
  .cat-icon { width: 64px; height: 64px; border-radius: 16px; display: flex; align-items: center; justify-content: center; color: #fff; flex-shrink: 0; }
  .cat-info { flex: 1; h3 { font-size: 18px; margin-bottom: 6px; } p { font-size: 14px; color: #606266; margin-bottom: 4px; } .article-count { font-size: 12px; color: #909399; } }
  .arrow { color: #c0c4cc; }
}
.sidebar { position: sticky; top: 80px; }
.sidebar-card { background: #fff; border-radius: 12px; padding: 20px; margin-bottom: 16px; border: 1px solid #e4e7ed;
  h3 { font-size: 16px; margin-bottom: 16px; color: #303133; }
  .hot-item { display: flex; align-items: center; gap: 12px; padding: 10px 0; cursor: pointer; border-bottom: 1px solid #f0f2f5;
    &:last-child { border-bottom: none; }
    &:hover .hot-title { color: var(--el-color-primary); }
    .hot-rank { width: 24px; height: 24px; border-radius: 6px; background: #f0f2f5; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 600; color: #909399; &.top { background: var(--el-color-primary); color: #fff; } }
    .hot-title { font-size: 14px; color: #606266; }
  }
  .tags { display: flex; flex-wrap: wrap; gap: 8px; .tag-item { cursor: pointer; } }
}
</style>
