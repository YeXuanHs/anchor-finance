<template>
  <div class="terms-page">
    <SiteHeader />
    
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">服务条款</h1>
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
  { id: 'acceptance', title: '条款接受' },
  { id: 'services', title: '服务说明' },
  { id: 'account', title: '账户注册' },
  { id: 'payment', title: '付款条款' },
  { id: 'usage', title: '使用规范' },
  { id: 'intellectual', title: '知识产权' },
  { id: 'limitation', title: '责任限制' },
  { id: 'termination', title: '服务终止' },
  { id: 'governing', title: '适用法律' },
  { id: 'changes', title: '条款变更' }
])

const sections = ref([
  {
    id: 'acceptance',
    title: '条款接受',
    content: '<p>欢迎使用锚点财务提供的服务。在使用我们的服务之前，请仔细阅读并理解以下服务条款。</p><p>通过访问或使用我们的服务，您同意受本服务条款的约束。如果您不同意这些条款，请不要使用我们的服务。</p>'
  },
  {
    id: 'services',
    title: '服务说明',
    content: '<p>锚点财务提供以下服务：</p><ul><li>云服务器租用和管理</li><li>VPS主机服务</li><li>独立服务器租用</li><li>域名注册服务</li><li>SSL证书服务</li><li>技术支持服务</li></ul><p>我们保留随时修改、暂停或终止任何服务的权利。</p>'
  },
  {
    id: 'account',
    title: '账户注册',
    content: '<p>使用我们的服务需要注册账户。您同意：</p><ul><li>提供真实、准确、完整的注册信息</li><li>维护并及时更新您的账户信息</li><li>保护您的账户安全，不得与他人共享账户</li><li>对您账户下的所有活动负责</li></ul><p>如发现未经授权使用您账户的情况，请立即通知我们。</p>'
  },
  {
    id: 'payment',
    title: '付款条款',
    content: '<p>关于付款，您同意：</p><ul><li>按时支付所有费用</li><li>提供有效的付款方式</li><li>承担因付款延迟产生的额外费用</li></ul><p>所有费用均不含税，除非另有说明。我们保留调整价格的权利，但会提前通知您。</p>'
  },
  {
    id: 'usage',
    title: '使用规范',
    content: '<p>使用我们的服务时，您不得：</p><ul><li>违反任何适用的法律法规</li><li>侵犯他人的知识产权</li><li>传播恶意软件或进行网络攻击</li><li>发送垃圾邮件或进行欺诈活动</li><li>干扰或破坏服务的正常运行</li><li>超出合理的使用限制</li></ul><p>违反使用规范可能导致服务暂停或终止。</p>'
  },
  {
    id: 'intellectual',
    title: '知识产权',
    content: '<p>我们的服务、商标、标识和内容均受知识产权法保护。未经我们明确授权，您不得：</p><ul><li>复制、修改或分发我们的内容</li><li>使用我们的商标或标识</li><li>反向工程或试图获取源代码</li></ul>'
  },
  {
    id: 'limitation',
    title: '责任限制',
    content: '<p>在法律允许的最大范围内：</p><ul><li>我们不对任何间接、附带或后果性损害负责</li><li>我们的总责任不超过您在过去12个月内支付的费用</li><li>我们不对因不可抗力导致的服务中断负责</li></ul>'
  },
  {
    id: 'termination',
    title: '服务终止',
    content: '<p>我们可能在以下情况下终止您的服务：</p><ul><li>您违反本服务条款</li><li>您未按时支付费用</li><li>法律要求或政府命令</li><li>我们决定停止提供服务</li></ul><p>服务终止后，您需要备份并删除存储在我们服务器上的数据。</p>'
  },
  {
    id: 'governing',
    title: '适用法律',
    content: '<p>本服务条款受中华人民共和国法律管辖。任何争议应首先通过友好协商解决。协商不成的，应提交北京仲裁委员会仲裁解决。</p>'
  },
  {
    id: 'changes',
    title: '条款变更',
    content: '<p>我们保留随时修改本服务条款的权利。修改后的条款将在本页面发布。重大变更时，我们会通过邮件或服务内通知的方式告知您。</p><p>继续使用我们的服务即表示您同意修改后的条款。</p>'
  }
])
</script>

<style scoped lang="scss">
.terms-page {
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
  .terms-page .content-section .content-wrapper {
    grid-template-columns: 1fr;
  }
  
  .content-toc {
    position: static !important;
  }
}
</style>
