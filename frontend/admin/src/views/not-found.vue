<template>
  <div class="not-found">
    <div class="content">
      <div class="error-code">404</div>
      <h1>页面未找到</h1>
      <p>抱歉，您访问的管理页面不存在</p>
      <el-button type="primary" @click="router.push('/')">
        <el-icon><HomeFilled /></el-icon>
        返回仪表盘
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled } from '@element-plus/icons-vue'

const router = useRouter()
</script>

<style scoped lang="scss">
.not-found {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;

  .content {
    text-align: center;

    .error-code {
      font-size: 120px;
      font-weight: 900;
      color: var(--el-color-primary);
      line-height: 1;
      margin-bottom: 16px;
    }

    h1 {
      font-size: 24px;
      color: #303133;
      margin-bottom: 8px;
    }

    p {
      font-size: 14px;
      color: #909399;
      margin-bottom: 24px;
    }
  }
}
</style>
