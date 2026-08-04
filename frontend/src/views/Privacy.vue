<template>
  <div class="privacy-page">
    <SiteHeader />
    
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">隐私政策</h1>
        <p class="hero-desc">最后更新：2024年1月1日</p>
      </div>
    </section>
    
    <section class="content-section">
      <div class="container">
        <div class="content-wrapper">
          <div class="content-toc">
            <h3>目录</h3>
            <ul>
              <li v-for="(item, index) in toc" :key="index">
                <a :href="`#${item.id}`">{{ item.title }}</a>
              </li>
            </ul>
          </div>
          
          <div class="content-main">
            <div class="content-block" v-for="(section, index) in sections" :key="index" :id="section.id">
              <h2>{{ section.title }}</h2>
              <div v-html="section.content"></div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'

const toc = ref([
  { id: 'introduction', title: '引言' },
  { id: 'collection', title: '信息收集' },
  { id: 'usage', title: '信息使用' },
  { id: 'sharing', title: '信息共享' },
  { id: 'security', title: '信息安全' },
  { id: 'cookies', title: 'Cookie使用' },
  { id: 'rights', title: '您的权利' },
  { id: 'children', title: '未成年人保护' },
  { id: 'changes', title: '政策变更' },
  { id: 'contact', title: '联系我们' }
])

const sections = ref([
  {
    id: 'introduction',
    title: '引言',
    content: '<p>欢迎使用我们的服务（以下简称"我们"）。我们非常重视您的隐私保护。本隐私政策旨在向您说明我们如何收集、使用、存储和保护您的个人信息。</p><p>请在使用我们的服务前仔细阅读本隐私政策。使用我们的服务即表示您同意本隐私政策中描述的数据实践。</p>'
  },
  {
    id: 'collection',
    title: '信息收集',
    content: '<p>我们可能收集以下类型的信息：</p><ul><li><strong>账户信息：</strong>注册时提供的姓名、邮箱、手机号等</li><li><strong>支付信息：</strong>用于处理付款的银行卡、支付宝等信息</li><li><strong>使用数据：</strong>您使用我们服务的方式、时间和频率</li><li><strong>设备信息：</strong>设备型号、操作系统、浏览器类型等</li><li><strong>日志信息：</strong>IP地址、访问时间、访问页面等</li></ul>'
  },
  {
    id: 'usage',
    title: '信息使用',
    content: '<p>我们使用收集的信息用于：</p><ul><li>提供、维护和改进我们的服务</li><li>处理您的交易和发送相关通知</li><li>发送技术通知、更新和安全警报</li><li>响应您的评论、问题和客户服务请求</li><li>监控和分析趋势、使用情况和活动</li><li>检测、调查和防止欺诈行为</li></ul>'
  },
  {
    id: 'sharing',
    title: '信息共享',
    content: '<p>我们不会出售您的个人信息。我们可能在以下情况下共享您的信息：</p><ul><li><strong>服务提供商：</strong>与帮助我们运营业务的第三方服务提供商共享</li><li><strong>法律要求：</strong>根据法律法规、法律程序或政府要求</li><li><strong>保护权利：</strong>为保护我们的权利、隐私、安全或财产</li><li><strong>业务转让：</strong>在合并、收购或资产出售的情况下</li></ul>'
  },
  {
    id: 'security',
    title: '信息安全',
    content: '<p>我们采取合理的技术和组织措施来保护您的个人信息，包括：</p><ul><li>使用加密技术保护数据传输</li><li>实施访问控制机制</li><li>定期进行安全审计</li><li>对员工进行隐私保护培训</li></ul><p>但请注意，没有任何安全措施是完全无懈可击的。</p>'
  },
  {
    id: 'cookies',
    title: 'Cookie使用',
    content: '<p>我们使用Cookie和类似技术来：</p><ul><li>记住您的偏好设置</li><li>分析服务使用情况</li><li>提供个性化体验</li><li>保障服务安全</li></ul><p>您可以通过浏览器设置管理Cookie偏好。</p>'
  },
  {
    id: 'rights',
    title: '您的权利',
    content: '<p>根据适用法律，您可能享有以下权利：</p><ul><li>访问和获取您的个人信息副本</li><li>更正不准确的个人信息</li><li>删除您的个人信息</li><li>限制或反对处理您的个人信息</li><li>数据可携带权</li><li>撤回同意</li></ul><p>如需行使上述权利，请联系我们。</p>'
  },
  {
    id: 'children',
    title: '未成年人保护',
    content: '<p>我们的服务不面向16岁以下的未成年人。我们不会故意收集16岁以下未成年人的个人信息。如果我们发现已收集了未成年人的个人信息，我们将采取措施删除相关信息。</p>'
  },
  {
    id: 'changes',
    title: '政策变更',
    content: '<p>我们可能会不时更新本隐私政策。更新后的政策将在本页面发布，并注明最新修订日期。重大变更时，我们会通过邮件或服务内通知的方式告知您。</p>'
  },
  {
    id: 'contact',
    title: '联系我们',
    content: '<p>如果您对本隐私政策有任何疑问，请通过网站底部的联系方式联系我们。</p>'
  }
])
</script>

<style scoped lang="scss">
.privacy-page {
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
      font-size: 16px;
      opacity: 0.8;
      margin: 0;
    }
  }
  
  .content-section {
    padding: 60px 0;
    background: #f5f7fa;
    
    .content-wrapper {
      display: grid;
      grid-template-columns: 250px 1fr;
      gap: 40px;
    }
    
    .content-toc {
      background: #fff;
      border-radius: 12px;
      padding: 24px;
      position: sticky;
      top: 80px;
      height: fit-content;
      
      h3 {
        font-size: 16px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 16px;
      }
      
      ul {
        list-style: none;
        padding: 0;
        margin: 0;
        
        li {
          margin-bottom: 8px;
          
          a {
            font-size: 14px;
            color: #666;
            text-decoration: none;
            
            &:hover {
              color: #409eff;
            }
          }
        }
      }
    }
    
    .content-main {
      background: #fff;
      border-radius: 12px;
      padding: 40px;
      
      .content-block {
        margin-bottom: 40px;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        h2 {
          font-size: 24px;
          font-weight: 700;
          color: #1a2332;
          margin: 0 0 16px;
          padding-bottom: 12px;
          border-bottom: 2px solid #409eff;
        }
        
        :deep(p) {
          font-size: 15px;
          line-height: 1.8;
          color: #555;
          margin: 0 0 16px;
        }
        
        :deep(ul) {
          padding-left: 20px;
          margin: 0 0 16px;
          
          li {
            font-size: 15px;
            line-height: 1.8;
            color: #555;
            margin-bottom: 8px;
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .privacy-page .content-section .content-wrapper {
    grid-template-columns: 1fr;
  }
  
  .content-toc {
    position: static !important;
  }
}
</style>
