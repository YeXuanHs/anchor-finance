<template>
  <div class="general-settings page-container">
    <div class="art-card">
      <h3>基本设置</h3>
      <el-form :model="form" label-width="120px">
        <el-form-item label="站点名称">
          <el-input v-model="form.site_name" />
        </el-form-item>
        <el-form-item label="站点URL">
          <el-input v-model="form.site_url" />
        </el-form-item>
        <el-form-item label="站点描述">
          <el-input v-model="form.site_description" type="textarea" />
        </el-form-item>
        <el-form-item label="Logo">
          <el-upload
            action="#"
            :show-file-list="false"
            :auto-upload="false"
          >
            <img v-if="form.site_logo" :src="form.site_logo" class="logo-preview" />
            <el-button v-else>选择图片</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item label="开放注册">
          <el-switch v-model="form.enable_registration" />
        </el-form-item>
        <el-form-item label="默认货币">
          <el-select v-model="form.currency">
            <el-option label="人民币 (CNY)" value="CNY" />
            <el-option label="美元 (USD)" value="USD" />
            <el-option label="欧元 (EUR)" value="EUR" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveSettings">保存设置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const form = ref({
  site_name: '锚点财务',
  site_url: 'http://localhost:8080',
  site_description: '高效、安全的财务管理平台',
  site_logo: '/logo.png',
  enable_registration: true,
  currency: 'CNY'
})

const saveSettings = async () => {
  // TODO: 保存设置
  ElMessage.success('保存成功')
}
</script>

<style scoped lang="scss">
.general-settings {
  .art-card {
    max-width: 800px;
    
    h3 {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 24px;
      padding-bottom: 16px;
      border-bottom: 1px solid var(--border-color);
    }
    
    .logo-preview {
      width: 100px;
      height: 100px;
      object-fit: contain;
      border: 1px solid var(--border-color);
      border-radius: 8px;
    }
  }
}
</style>
