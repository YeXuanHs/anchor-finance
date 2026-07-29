<template>
  <div class="sitemap-page">
    <SiteHeader />
    
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">网站地图</h1>
        <p class="hero-desc">快速找到您需要的页面</p>
      </div>
    </section>
    
    <section class="content-section">
      <div class="container">
        <div class="sitemap-grid">
          <!-- Products -->
          <div class="sitemap-card">
            <h3>产品服务</h3>
            <ul>
              <li><router-link to="/products">全部产品</router-link></li>
              <li><router-link to="/products?group=cloud">云服务器</router-link></li>
              <li><router-link to="/products?group=vps">VPS主机</router-link></li>
              <li><router-link to="/products?group=dedicated">独立服务器</router-link></li>
              <li><router-link to="/products?group=domain">域名注册</router-link></li>
              <li><router-link to="/products?group=ssl">SSL证书</router-link></li>
            </ul>
          </div>
          
          <!-- Solutions -->
          <div class="sitemap-card">
            <h3>解决方案</h3>
            <ul>
              <li><router-link to="/solutions">全部方案</router-link></li>
              <li><router-link to="/solutions/video">视频直播</router-link></li>
              <li><router-link to="/solutions/game">游戏加速</router-link></li>
              <li><router-link to="/solutions/ecommerce">电商平台</router-link></li>
              <li><router-link to="/solutions/edu">在线教育</router-link></li>
              <li><router-link to="/solutions/security">安全防护</router-link></li>
            </ul>
          </div>
          
          <!-- Support -->
          <div class="sitemap-card">
            <h3>帮助支持</h3>
            <ul>
              <li><router-link to="/help">帮助中心</router-link></li>
              <li><router-link to="/knowledge-base">知识库</router-link></li>
              <li><router-link to="/downloads">下载中心</router-link></li>
              <li><router-link to="/contact">联系我们</router-link></li>
            </ul>
          </div>
          
          <!-- About -->
          <div class="sitemap-card">
            <h3>关于我们</h3>
            <ul>
              <li><router-link to="/about">公司介绍</router-link></li>
              <li><router-link to="/news">新闻动态</router-link></li>
              <li><router-link to="/privacy">隐私政策</router-link></li>
              <li><router-link to="/terms">服务条款</router-link></li>
            </ul>
          </div>
          
          <!-- User Center -->
          <div class="sitemap-card">
            <h3>用户中心</h3>
            <ul>
              <li><router-link to="/user/dashboard">控制台</router-link></li>
              <li><router-link to="/user/products">我的产品</router-link></li>
              <li><router-link to="/user/orders">我的订单</router-link></li>
              <li><router-link to="/user/tickets">工单管理</router-link></li>
              <li><router-link to="/user/profile">个人资料</router-link></li>
              <li><router-link to="/user/security">安全设置</router-link></li>
            </ul>
          </div>
          
          <!-- Account -->
          <div class="sitemap-card">
            <h3>账户相关</h3>
            <ul>
              <li><router-link to="/login">登录</router-link></li>
              <li><router-link to="/register">注册</router-link></li>
              <li><router-link to="/forgot-password">忘记密码</router-link></li>
            </ul>
          </div>
        </div>
      </div>
    </section>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
</script>

<style scoped lang="scss">
.sitemap-page {
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
      margin: 0;
    }
  }
  
  .content-section {
    padding: 60px 0;
    background: #f5f7fa;
    
    .sitemap-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 30px;
    }
    
    .sitemap-card {
      background: #fff;
      border-radius: 12px;
      padding: 32px;
      
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 20px;
        padding-bottom: 12px;
        border-bottom: 2px solid #409eff;
      }
      
      ul {
        list-style: none;
        padding: 0;
        margin: 0;
        
        li {
          margin-bottom: 12px;
          
          a {
            font-size: 14px;
            color: #555;
            text-decoration: none;
            transition: color 0.3s;
            
            &:hover {
              color: #409eff;
            }
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .sitemap-page .sitemap-grid {
    grid-template-columns: 1fr;
  }
}
</style>
