<template>
  <div class="contact-page">
    <SiteHeader />
    
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="container">
        <h1 class="hero-title">联系我们</h1>
        <p class="hero-desc">我们期待与您的沟通，为您提供最优质的服务</p>
      </div>
    </section>
    
    <!-- Contact Content -->
    <section class="section contact-section">
      <div class="container">
        <div class="contact-grid">
          <!-- Contact Form -->
          <div class="contact-form-wrapper">
            <h2>发送消息</h2>
            <p>请填写以下表单，我们会尽快回复您</p>
            
            <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="您的姓名" prop="name">
                    <el-input v-model="form.name" placeholder="请输入您的姓名" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="电子邮箱" prop="email">
                    <el-input v-model="form.email" placeholder="请输入您的邮箱" />
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="联系电话" prop="phone">
                    <el-input v-model="form.phone" placeholder="请输入您的电话" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="咨询类型" prop="type">
                    <el-select v-model="form.type" placeholder="请选择咨询类型" style="width: 100%">
                      <el-option label="产品咨询" value="product" />
                      <el-option label="技术支持" value="support" />
                      <el-option label="商务合作" value="business" />
                      <el-option label="投诉建议" value="feedback" />
                      <el-option label="其他" value="other" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-form-item label="主题" prop="subject">
                <el-input v-model="form.subject" placeholder="请输入主题" />
              </el-form-item>
              
              <el-form-item label="详细内容" prop="content">
                <el-input v-model="form.content" type="textarea" :rows="6" placeholder="请输入详细内容" />
              </el-form-item>
              
              <el-form-item>
                <el-button type="primary" size="large" @click="submitForm" :loading="submitting">
                  <el-icon style="margin-right: 8px;"><Promotion /></el-icon>
                  提交消息
                </el-button>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- Contact Info -->
          <div class="contact-info">
            <div class="info-card">
              <div class="info-icon">
                <el-icon :size="28"><Location /></el-icon>
              </div>
              <h3>公司地址</h3>
              <p>{{ siteSettings.contact_address }}</p>
            </div>
            
            <div class="info-card">
              <div class="info-icon">
                <el-icon :size="28"><Phone /></el-icon>
              </div>
              <h3>联系电话</h3>
              <p>销售咨询：{{ siteSettings.sales_phone || siteSettings.contact_phone }}</p>
              <p>技术支持：{{ siteSettings.support_phone || siteSettings.contact_phone }}</p>
            </div>
            
            <div class="info-card">
              <div class="info-icon">
                <el-icon :size="28"><Message /></el-icon>
              </div>
              <h3>电子邮箱</h3>
              <p>销售咨询：{{ siteSettings.sales_email || siteSettings.contact_email }}</p>
              <p>技术支持：{{ siteSettings.contact_email }}</p>
            </div>
            
            <div class="info-card">
              <div class="info-icon">
                <el-icon :size="28"><Clock /></el-icon>
              </div>
              <h3>工作时间</h3>
              <p>{{ siteSettings.work_time }}</p>
              <p>节假日休息</p>
            </div>
            
            <!-- Social Media -->
            <div class="social-links">
              <h3>关注我们</h3>
              <div class="social-icons">
                <a href="#" class="social-icon" title="微信">
                  <svg viewBox="0 0 24 24" width="24" height="24"><path fill="currentColor" d="M9.5 4C5.36 4 2 6.69 2 10c0 1.89 1.08 3.56 2.78 4.66L4 17l2.5-1.18C7.61 16.27 8.54 16.5 9.5 16.5c.17 0 .33 0 .5-.02C9.68 15.83 9.5 15.18 9.5 14.5c0-3.04 2.69-5.5 6-5.5.17 0 .33 0 .5.02C15.27 5.94 12.64 4 9.5 4m-2.25 3.5a1 1 0 110 2 1 1 0 010-2m4.5 0a1 1 0 110 2 1 1 0 010-2M15.5 10c-3.04 0-5.5 2.24-5.5 5s2.46 5 5.5 5c.87 0 1.69-.2 2.44-.56L20 20.5l-.59-2.07C20.42 17.36 21 16.23 21 15c0-2.76-2.46-5-5.5-5m-2 4a1 1 0 110 2 1 1 0 010-2m4 0a1 1 0 110 2 1 1 0 010-2z"/></svg>
                </a>
                <a href="#" class="social-icon" title="微博">
                  <svg viewBox="0 0 24 24" width="24" height="24"><path fill="currentColor" d="M10.6 18.3c-3.6.4-6.7-1.3-6.9-3.7-.2-2.4 2.6-4.7 6.2-5.1 3.6-.4 6.7 1.3 6.9 3.7.2 2.4-2.6 4.7-6.2 5.1m8.6-11.8c-.3-1.1-1.3-1.7-2.3-1.4l-.3.1c-.2.1-.4 0-.5-.2l-.1-.3c-.4-1.3-1.7-2.1-3.1-1.8l-.3.1c-.2 0-.4-.1-.5-.3l-.1-.3c-.2-.8-.9-1.3-1.7-1.2h-.3c-.2 0-.4-.1-.4-.3l-.1-.3c-.2-.8-.9-1.3-1.7-1.2-1 .1-1.8.9-1.8 1.9 0 .1 0 .2.1.3l.1.3c.1.2 0 .4-.2.5l-.3.1c-1.3.4-2.1 1.7-1.8 3.1.1.3.1.5.1.8 0 .2-.1.4-.3.5l-.3.1c-1.1.3-1.7 1.3-1.4 2.3.3 1.1 1.3 1.7 2.3 1.4l.3-.1c.2-.1.4 0 .5.2l.1.3c.4 1.3 1.7 2.1 3.1 1.8l.3-.1c.2 0 .4.1.5.3l.1.3c.4 1.3 1.7 2.1 3.1 1.8 1.3-.4 2.1-1.7 1.8-3.1l-.1-.3c-.1-.2 0-.4.2-.5l.3-.1c1.1-.3 1.7-1.3 1.4-2.3M17 7.5c.6 0 1 .4 1 1s-.4 1-1 1-1-.4-1-1 .4-1 1-1m-5.5 10c-2.5 0-4.5-1.3-4.5-3s2-3 4.5-3 4.5 1.3 4.5 3-2 3-4.5 3m5.2-7.2c-.3-.1-.5-.3-.5-.6 0-.3.2-.5.5-.6 1.6-.5 2.8.2 2.8 1.2s-1.2 1.7-2.8 1.2l.5-.6c.2-.2.2-.5 0-.7l-.5.5z"/></svg>
                </a>
                <a href="#" class="social-icon" title="抖音">
                  <svg viewBox="0 0 24 24" width="24" height="24"><path fill="currentColor" d="M19.59 6.69a4.83 4.83 0 01-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 01-2.88 2.5 2.89 2.89 0 01-2.89-2.89 2.89 2.89 0 012.89-2.89c.28 0 .54.04.79.1v-3.4a6.37 6.37 0 00-.79-.05A6.34 6.34 0 003.15 15.2a6.34 6.34 0 0010.86 4.46v-7.15a8.16 8.16 0 005.58 2.2V11.3a4.85 4.85 0 01-3.77-1.38V6.69h3.77z"/></svg>
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Map Section -->
    <section class="section map-section">
      <div class="container">
        <div class="map-placeholder">
          <el-icon :size="48"><Location /></el-icon>
          <p>地图加载中...</p>
        </div>
      </div>
    </section>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Location, Phone, Message, Clock, Promotion } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const siteSettings = ref({
  contact_address: '',
  contact_phone: '',
  sales_phone: '',
  support_phone: '',
  contact_email: '',
  sales_email: '',
  work_time: ''
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

const formRef = ref()
const submitting = ref(false)

const form = reactive({
  name: '',
  email: '',
  phone: '',
  type: '',
  subject: '',
  content: ''
})

const rules = {
  name: [{ required: true, message: '请输入您的姓名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入您的邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择咨询类型', trigger: 'change' }],
  subject: [{ required: true, message: '请输入主题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入详细内容', trigger: 'blur' }]
}

const submitForm = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    
    await request.post('/api/v1/contact', form)
    ElMessage.success('消息已发送，我们会尽快回复您！')
    
    // Reset form
    formRef.value.resetFields()
  } catch (error) {
    if (error !== false) {
      ElMessage.error('发送失败，请稍后重试')
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped lang="scss">
.contact-page {
  padding-top: 64px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .hero-section {
    background: linear-gradient(135deg, #1a237e, #0d47a1);
    color: #fff;
    padding: 100px 0;
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
  
  .section {
    padding: 80px 0;
  }
  
  .contact-section {
    background: #fff;
    
    .contact-grid {
      display: grid;
      grid-template-columns: 1fr 400px;
      gap: 60px;
    }
    
    .contact-form-wrapper {
      h2 {
        font-size: 28px;
        font-weight: 700;
        color: #1a2332;
        margin: 0 0 8px;
      }
      
      p {
        font-size: 15px;
        color: #666;
        margin: 0 0 30px;
      }
    }
    
    .contact-info {
      .info-card {
        padding: 24px;
        background: #f5f7fa;
        border-radius: 12px;
        margin-bottom: 16px;
        
        .info-icon {
          width: 48px;
          height: 48px;
          background: linear-gradient(135deg, #409eff, #66b1ff);
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
          margin-bottom: 16px;
        }
        
        h3 {
          font-size: 16px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 8px;
        }
        
        p {
          font-size: 14px;
          color: #666;
          margin: 0 0 4px;
        }
      }
      
      .social-links {
        h3 {
          font-size: 16px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 16px;
        }
        
        .social-icons {
          display: flex;
          gap: 12px;
          
          .social-icon {
            width: 40px;
            height: 40px;
            background: #f5f7fa;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #666;
            transition: all 0.3s;
            
            &:hover {
              background: #409eff;
              color: #fff;
            }
          }
        }
      }
    }
  }
  
  .map-section {
    background: #f5f7fa;
    
    .map-placeholder {
      height: 400px;
      background: #e8ecf1;
      border-radius: 12px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      color: #909399;
      
      p {
        margin: 12px 0 0;
      }
    }
  }
}

@media (max-width: 768px) {
  .contact-page .contact-section .contact-grid {
    grid-template-columns: 1fr;
    gap: 40px;
  }
}
</style>
