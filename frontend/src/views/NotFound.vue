<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1>页面未找到</h1>
      <p>抱歉，您访问的页面不存在或已被移除</p>
      <div class="actions">
        <el-button type="primary" size="large" @click="goHome">
          <el-icon><HomeFilled /></el-icon>
          返回首页
        </el-button>
        <el-button size="large" @click="goBack">
          <el-icon><Back /></el-icon>
          返回上页
        </el-button>
      </div>
      <div class="help-links">
        <p>您可能在找：</p>
        <div class="links">
          <router-link to="/products">产品中心</router-link>
          <router-link to="/knowledge-base">知识库</router-link>
          <router-link to="/help">帮助中心</router-link>
          <router-link to="/contact">联系我们</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled, Back } from '@element-plus/icons-vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}
</script>

<style scoped lang="scss">
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  padding: 20px;
}

.not-found-content {
  text-align: center;
  max-width: 500px;

  .error-code {
    font-size: 150px;
    font-weight: 900;
    background: linear-gradient(135deg, var(--el-color-primary) 0%, #409eff 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    line-height: 1;
    margin-bottom: 20px;
    text-shadow: none;
  }

  h1 {
    font-size: 28px;
    color: #303133;
    margin-bottom: 12px;
    font-weight: 600;
  }

  p {
    font-size: 16px;
    color: #909399;
    margin-bottom: 32px;
  }

  .actions {
    display: flex;
    gap: 16px;
    justify-content: center;
    margin-bottom: 48px;
  }

  .help-links {
    p {
      font-size: 14px;
      color: #909399;
      margin-bottom: 12px;
    }

    .links {
      display: flex;
      gap: 24px;
      justify-content: center;
      flex-wrap: wrap;

      a {
        color: var(--el-color-primary);
        text-decoration: none;
        font-size: 14px;
        transition: opacity 0.2s;

        &:hover {
          opacity: 0.8;
          text-decoration: underline;
        }
      }
    }
  }
}
</style>
