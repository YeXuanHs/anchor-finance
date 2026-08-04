<template>
  <div class="aup-page">
    <SiteHeader />
    <section class="hero-section">
      <div class="container">
        <h1>可接受使用政策</h1>
        <p class="hero-desc">最后更新：2024年1月1日</p>
      </div>
    </section>
    <section class="content-section">
      <div class="container">
        <div class="content-main">
          <div class="content-block" v-for="(s, i) in sections" :key="i" :id="s.id">
            <h2>{{ s.title }}</h2>
            <div v-html="s.content"></div>
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
const sections = [
  { id: 'general', title: '总则', content: '<p>本可接受使用政策（AUP）规定了使用我们服务时必须遵守的规则。使用我们的服务即表示您同意遵守本政策。</p>' },
  { id: 'prohibited', title: '禁止行为', content: '<p>以下行为严格禁止：</p><ul><li>发送垃圾邮件或未经请求的批量邮件</li><li>传播恶意软件、病毒或有害代码</li><li>进行DDoS攻击或网络扫描</li><li>托管钓鱼网站或欺诈内容</li><li>侵犯他人知识产权</li><li>存储或传播非法内容</li><li>未经授权访问他人系统或数据</li><li>运行加密货币挖矿程序</li></ul>' },
  { id: 'content', title: '内容政策', content: '<p>您不得在我们的服务器上存储或传播以下内容：</p><ul><li>违反国家法律法规的内容</li><li>侵犯版权的内容</li><li>色情或儿童不宜内容</li><li>暴力或仇恨言论</li><li>虚假或误导性信息</li></ul>' },
  { id: 'security', title: '安全要求', content: '<p>您有责任维护您账户的安全：</p><ul><li>使用强密码并定期更换</li><li>不在公共设备上保存登录凭据</li><li>及时报告安全漏洞</li><li>保持系统和软件更新</li></ul>' },
  { id: 'enforcement', title: '执行措施', content: '<p>违反本政策可能导致以下后果：</p><ul><li>警告通知</li><li>临时暂停服务</li><li>永久终止服务</li><li>报告给相关执法部门</li></ul>' },
  { id: 'reporting', title: '举报', content: '<p>如发现违反本政策的行为，请通过用户中心提交工单或查看网站底部联系方式进行举报。</p>' },
  { id: 'changes', title: '政策变更', content: '<p>我们保留随时修改本政策的权利。重大变更将通过邮件或站内公告通知用户。</p>' }
]
</script>
<style scoped lang="scss">
.aup-page { min-height: 100vh; }
.hero-section { background: linear-gradient(135deg, #1a365d 0%, #2563eb 100%); padding: 80px 0 60px; text-align: center; color: #fff;
  h1 { font-size: 36px; font-weight: 700; margin-bottom: 12px; }
  .hero-desc { font-size: 16px; opacity: 0.8; }
}
.content-section { padding: 60px 0; background: #f8fafc; }
.container { max-width: 900px; margin: 0 auto; padding: 0 20px; }
.content-main { background: #fff; border-radius: 12px; padding: 40px; box-shadow: 0 2px 12px rgba(0,0,0,0.06); }
.content-block { margin-bottom: 32px;
  h2 { font-size: 22px; color: #1a365d; margin-bottom: 16px; padding-bottom: 8px; border-bottom: 2px solid #e5e7eb; }
  :deep(p) { color: #4b5563; line-height: 1.8; margin-bottom: 12px; }
  :deep(ul) { padding-left: 20px; color: #4b5563; line-height: 2; }
}
</style>
