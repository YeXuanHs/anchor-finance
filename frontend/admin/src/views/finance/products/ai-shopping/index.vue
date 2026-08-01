<template>
  <div class="ai-shopping-admin">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>AI 购物助手</span>
          <el-tag type="info" size="small">mahiru_ai_shopping</el-tag>
        </div>
      </template>

      <el-form :model="config" label-width="140px" style="max-width: 700px">
        <el-divider content-position="left">AI 模型配置</el-divider>
        <el-form-item label="API 地址">
          <el-input v-model="config.api_endpoint" placeholder="https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="config.api_key" type="password" show-password />
        </el-form-item>
        <el-form-item label="模型">
          <el-input v-model="config.model" placeholder="gpt-3.5-turbo" />
        </el-form-item>

        <el-divider content-position="left">助手设置</el-divider>
        <el-form-item label="启用助手">
          <el-switch v-model="config.enabled" />
        </el-form-item>
        <el-form-item label="欢迎语">
          <el-input v-model="config.welcome_message" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="系统提示词">
          <el-input v-model="config.system_prompt" type="textarea" :rows="4" placeholder="你是一个专业的购物助手..." />
        </el-form-item>
        <el-form-item label="推荐商品数">
          <el-input-number v-model="config.max_recommendations" :min="1" :max="20" />
        </el-form-item>

        <el-divider content-position="left">商品目录展示</el-divider>
        <el-form-item label="显示商品目录">
          <el-switch v-model="config.catalog_enabled" />
        </el-form-item>
        <el-form-item label="目录布局">
          <el-radio-group v-model="config.catalog_layout">
            <el-radio value="grid">网格</el-radio>
            <el-radio value="list">列表</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="显示筛选">
          <el-switch v-model="config.show_filters" />
        </el-form-item>
        <el-form-item label="显示价格">
          <el-switch v-model="config.show_price" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveConfig">保存配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const config = reactive<any>({
  api_endpoint: '', api_key: '', model: '',
  enabled: false, welcome_message: '你好！我是购物助手，有什么可以帮您？',
  system_prompt: '你是一个专业的购物助手，帮助用户选择合适的商品。根据用户需求推荐产品，并提供详细的产品信息。',
  max_recommendations: 5, catalog_enabled: true, catalog_layout: 'grid',
  show_filters: true, show_price: true,
})

const loadConfig = async () => {
  const { data } = await request.get('/admin/ai-shopping/config')
  if (data) Object.assign(config, data)
}

const saveConfig = async () => {
  await request.put('/admin/ai-shopping/config', config)
  ElMessage.success('保存成功')
}

onMounted(loadConfig)
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
</style>
