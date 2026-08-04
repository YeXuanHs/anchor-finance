<template>
  <footer class="site-footer">
    <div class="container">
      <div class="footer-grid">
        <div class="footer-col footer-brand">
          <div class="footer-logo">
            <img src="/logo.png" alt="锚点财务" />
            <span>锚点财务</span>
          </div>
          <p class="footer-desc">高效、安全的财务管理系统，助力企业数字化转型</p>
          <div class="footer-contact">
            <p><el-icon><Phone /></el-icon> {{ siteSettings.contact_phone }}</p>
            <p><el-icon><Message /></el-icon> {{ siteSettings.contact_email }}</p>
          </div>
        </div>
        
        <div class="footer-col">
          <h4>产品服务</h4>
          <ul>
            <li><router-link to="/products">云服务器</router-link></li>
            <li><router-link to="/products">VPS主机</router-link></li>
            <li><router-link to="/products">独立服务器</router-link></li>
            <li><router-link to="/products">域名注册</router-link></li>
            <li><router-link to="/products">SSL证书</router-link></li>
          </ul>
        </div>
        
        <div class="footer-col">
          <h4>解决方案</h4>
          <ul>
            <li><router-link to="/solutions/game">游戏加速</router-link></li>
            <li><router-link to="/solutions/video">视频直播</router-link></li>
            <li><router-link to="/solutions/edu">在线教育</router-link></li>
            <li><router-link to="/solutions/ecommerce">电商平台</router-link></li>
            <li><router-link to="/solutions/security">安全防护</router-link></li>
          </ul>
        </div>
        
        <div class="footer-col">
          <h4>帮助支持</h4>
          <ul>
            <li><router-link to="/help">帮助中心</router-link></li>
            <li><router-link to="/knowledge-base">知识库</router-link></li>
            <li><router-link to="/downloads">下载中心</router-link></li>
            <li><router-link to="/contact">联系我们</router-link></li>
          </ul>
        </div>
        
        <div class="footer-col">
          <h4>关于我们</h4>
          <ul>
            <li><router-link to="/about">公司介绍</router-link></li>
            <li><router-link to="/news">新闻动态</router-link></li>
            <li><router-link to="/privacy">隐私政策</router-link></li>
            <li><router-link to="/terms">服务条款</router-link></li>
          </ul>
        </div>
      </div>
      
      <div class="footer-bottom">
        <div class="footer-bottom-left">
          <p>&copy; {{ currentYear }} 锚点财务 All Rights Reserved</p>
        </div>
        <div class="footer-bottom-right">
          <router-link to="/privacy">隐私政策</router-link>
          <span class="divider">|</span>
          <router-link to="/terms">服务条款</router-link>
          <span class="divider">|</span>
          <router-link to="/site-map">网站地图</router-link>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Phone, Message } from '@element-plus/icons-vue'
import request from '@/utils/request'

const currentYear = computed(() => new Date().getFullYear())

const siteSettings = ref({
  contact_phone: '400-000-0000',
  contact_email: 'support@anchorfinance.com',
  site_name: '锚点财务'
})

const fetchSiteSettings = async () => {
  try {
    const res = await request.get('/api/v1/site/settings')
    if (res?.data) {
      siteSettings.value = { ...siteSettings.value, ...res.data }
    }
  } catch {
    // Use defaults
  }
}

onMounted(() => {
  fetchSiteSettings()
})
</script>

<style scoped lang="scss">
.site-footer {
  background: #1a2332;
  color: rgba(255, 255, 255, 0.7);
  padding: 60px 0 0;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .footer-grid {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr 1fr 1fr;
    gap: 40px;
    padding-bottom: 40px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .footer-col {
    h4 {
      color: #fff;
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 20px;
    }
    
    ul {
      list-style: none;
      padding: 0;
      margin: 0;
      
      li {
        margin-bottom: 12px;
        
        a {
          color: rgba(255, 255, 255, 0.6);
          text-decoration: none;
          transition: color 0.3s;
          
          &:hover {
            color: #409eff;
          }
        }
      }
    }
  }
  
  .footer-brand {
    .footer-logo {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;
      
      img {
        width: 40px;
        height: 40px;
        border-radius: 8px;
      }
      
      span {
        font-size: 20px;
        font-weight: 700;
        color: #fff;
      }
    }
    
    .footer-desc {
      font-size: 14px;
      line-height: 1.6;
      margin-bottom: 20px;
    }
    
    .footer-contact {
      p {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 10px;
        font-size: 14px;
        
        .el-icon {
          color: #409eff;
        }
      }
    }
  }
  
  .footer-bottom {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 0;
    
    p {
      margin: 0;
      font-size: 13px;
    }
    
    .footer-bottom-right {
      display: flex;
      align-items: center;
      gap: 12px;
      
      a {
        color: rgba(255, 255, 255, 0.6);
        text-decoration: none;
        font-size: 13px;
        
        &:hover {
          color: #409eff;
        }
      }
      
      .divider {
        color: rgba(255, 255, 255, 0.3);
      }
    }
  }
}

@media (max-width: 768px) {
  .site-footer {
    .footer-grid {
      grid-template-columns: 1fr;
      gap: 30px;
    }
    
    .footer-bottom {
      flex-direction: column;
      gap: 10px;
      text-align: center;
    }
  }
}
</style>
