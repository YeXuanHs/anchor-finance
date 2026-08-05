<template>
  <div class="help-page">
    <SiteHeader />
    
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">{{ $t('helpCenter.title') }}</h1>
        <p class="hero-desc">{{ $t('helpCommon.foundAnswer') }}</p>
        
        <!-- Search -->
        <div class="search-wrapper">
          <el-input
            v-model="searchQuery"
            size="large"
            :placeholder="$t('helpCenter.searchPlaceholder')"
            :prefix-icon="Search"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">{{ $t('common.search') }}</el-button>
            </template>
          </el-input>
        </div>
      </div>
    </section>
    
    <!-- Hot Questions -->
    <section class="section hot-questions">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('helpCenter.hotQuestions') }}</h2>
          <p>{{ $t('helpCommon.commonQuestions') }}</p>
        </div>
        
        <div class="question-list">
          <div class="question-item" v-for="(item, index) in hotQuestions" :key="index" @click="goToArticle(item.id)">
            <div class="question-icon">
              <el-icon><QuestionFilled /></el-icon>
            </div>
            <div class="question-content">
              <h3>{{ item.title }}</h3>
              <p>{{ item.desc }}</p>
            </div>
            <el-icon class="question-arrow"><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Categories -->
    <section class="section categories-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('helpCenter.questionCategories') }}</h2>
          <p>{{ $t('helpCommon.searchByCategory') }}</p>
        </div>
        
        <div class="categories-grid">
          <div class="category-card" v-for="(cat, index) in categories" :key="index" @click="goToCategory(cat.id)">
            <div class="category-icon">
              <el-icon :size="32"><component :is="cat.icon" /></el-icon>
            </div>
            <h3>{{ cat.name }}</h3>
            <p>{{ cat.desc }}</p>
            <span class="article-count">{{ cat.count }} {{ $t('helpCommon.articles') }}</span>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Quick Links -->
    <section class="section quick-links">
      <div class="container">
        <div class="links-grid">
          <div class="link-card">
            <el-icon :size="40"><Tickets /></el-icon>
            <h3>{{ $t('helpCommon.submitTicket') }}</h3>
            <p>{{ $t('helpCommon.ticketDesc') }}</p>
            <el-button type="primary" @click="$router.push('/user/tickets')">{{ $t('helpCommon.submitTicket') }}</el-button>
          </div>
          <div class="link-card">
            <el-icon :size="40"><ChatDotRound /></el-icon>
            <h3>{{ $t('helpCommon.onlineService') }}</h3>
            <p>{{ $t('helpCommon.onlineServiceDesc') }}</p>
            <el-button type="primary">{{ $t('helpCommon.contactService') }}</el-button>
          </div>
          <div class="link-card">
            <el-icon :size="40"><Document /></el-icon>
            <h3>{{ $t('helpCommon.knowledgeBase') }}</h3>
            <p>{{ $t('helpCommon.knowledgeBaseDesc') }}</p>
            <el-button type="primary" @click="$router.push('/knowledge-base')">{{ $t('helpCommon.viewDocs') }}</el-button>
          </div>
          <div class="link-card">
            <el-icon :size="40"><VideoCamera /></el-icon>
            <h3>{{ $t('helpCommon.videoTutorial') }}</h3>
            <p>{{ $t('helpCommon.videoTutorialDesc') }}</p>
            <el-button type="primary">{{ $t('helpCommon.watchVideo') }}</el-button>
          </div>
        </div>
      </div>
    </section>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import request from '@/utils/request'
import { 
  Search, QuestionFilled, ArrowRight, Tickets, ChatDotRound, 
  Document, VideoCamera, ShoppingCart, Wallet, Setting, Connection
} from '@element-plus/icons-vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'

const { t } = useI18n()
const router = useRouter()
const searchQuery = ref('')

const hotQuestions = ref([
  { id: 1, title: '如何注册账户？', desc: '了解如何快速注册并激活您的账户' },
  { id: 2, title: '如何购买产品？', desc: '产品购买流程详解' },
  { id: 3, title: '如何续费产品？', desc: '产品续费方式和注意事项' },
  { id: 4, title: '如何提交工单？', desc: '遇到问题如何获取技术支持' },
  { id: 5, title: '如何使用优惠券？', desc: '优惠券使用方法和规则' }
])

const categories = ref([
  { id: 1, name: '账户相关', icon: 'Connection', desc: '注册、登录、账户设置', count: 15 },
  { id: 2, name: '产品服务', icon: 'ShoppingCart', desc: '购买、续费、升级', count: 23 },
  { id: 3, name: '财务相关', icon: 'Wallet', desc: '付款、退款、发票', count: 18 },
  { id: 4, name: '技术支持', icon: 'Setting', desc: '服务器、网络、安全', count: 32 }
])

const handleSearch = () => {
  if (searchQuery.value) {
    router.push(`/help/search?q=${searchQuery.value}`)
  }
}

const goToArticle = (id: number) => {
  router.push(`/help/article/${id}`)
}

const goToCategory = (id: number) => {
  router.push(`/help/category/${id}`)
}

onMounted(async () => {
  try {
    const [catRes, hotRes] = await Promise.all([
      request.get('/api/v1/help/categories'),
      request.get('/api/v1/help/articles/hot', { params: { limit: 6 } })
    ])
    if (catRes?.data?.data?.list?.length) {
      categories.value = catRes.data.data.list
    }
    if (hotRes?.data?.data?.list?.length) {
      hotQuestions.value = hotRes.data.data.list
    }
  } catch (e) {
    console.error('Failed to fetch help data:', e)
  }
})
</script>

<style scoped lang="scss">
.help-page {
  padding-top: 64px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .hero-section {
    background: linear-gradient(135deg, #1a237e, #0d47a1);
    color: #fff;
    padding: 80px 0;
    text-align: center;
    
    .hero-title {
      font-size: 42px;
      font-weight: 700;
      margin: 0 0 16px;
    }
    
    .hero-desc {
      font-size: 18px;
      opacity: 0.8;
      margin: 0 0 40px;
    }
    
    .search-wrapper {
      max-width: 600px;
      margin: 0 auto;
      
      :deep(.el-input__wrapper) {
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }
    }
  }
  
  .section {
    padding: 80px 0;
    
    .section-header {
      text-align: center;
      margin-bottom: 50px;
      
      h2 {
        font-size: 32px;
        font-weight: 700;
        color: #1a2332;
        margin: 0 0 12px;
      }
      
      p {
        font-size: 16px;
        color: #666;
        margin: 0;
      }
    }
  }
  
  .hot-questions {
    background: #fff;
    
    .question-list {
      max-width: 800px;
      margin: 0 auto;
    }
    
    .question-item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 20px;
      background: #f5f7fa;
      border-radius: 12px;
      margin-bottom: 12px;
      cursor: pointer;
      transition: all 0.3s;
      
      &:hover {
        background: #ecf5ff;
        transform: translateX(5px);
      }
      
      .question-icon {
        width: 40px;
        height: 40px;
        background: linear-gradient(135deg, #409eff, #66b1ff);
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
        flex-shrink: 0;
      }
      
      .question-content {
        flex: 1;
        
        h3 {
          font-size: 16px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 4px;
        }
        
        p {
          font-size: 14px;
          color: #666;
          margin: 0;
        }
      }
      
      .question-arrow {
        color: #909399;
      }
    }
  }
  
  .categories-section {
    background: #f5f7fa;
    
    .categories-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 24px;
    }
    
    .category-card {
      padding: 32px;
      background: #fff;
      border-radius: 12px;
      text-align: center;
      cursor: pointer;
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
      }
      
      .category-icon {
        width: 64px;
        height: 64px;
        margin: 0 auto 20px;
        background: linear-gradient(135deg, #409eff, #66b1ff);
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
      }
      
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 8px;
      }
      
      p {
        font-size: 14px;
        color: #666;
        margin: 0 0 12px;
      }
      
      .article-count {
        font-size: 12px;
        color: #409eff;
      }
    }
  }
  
  .quick-links {
    background: #fff;
    
    .links-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 24px;
    }
    
    .link-card {
      padding: 32px;
      background: #f5f7fa;
      border-radius: 12px;
      text-align: center;
      
      .el-icon {
        color: #409eff;
        margin-bottom: 20px;
      }
      
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 8px;
      }
      
      p {
        font-size: 14px;
        color: #666;
        margin: 0 0 20px;
      }
    }
  }
}

@media (max-width: 768px) {
  .help-page {
    .categories-grid,
    .links-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
