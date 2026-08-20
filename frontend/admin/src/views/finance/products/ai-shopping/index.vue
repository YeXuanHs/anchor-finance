<template>
  <div class="ai-shopping-admin">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('aiShopping.title') }}</span>
          <el-tag type="info" size="small">mahiru_ai_shopping</el-tag>
        </div>
      </template>

      <el-form :model="config" label-width="140px" style="max-width: 700px">
        <el-divider content-position="left">{{ $t('aiShopping.modelConfig') }}</el-divider>
        <el-form-item :label="$t('aiShopping.apiEndpoint')">
          <el-input v-model="config.api_endpoint" placeholder="https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="config.api_key" type="password" show-password />
        </el-form-item>
        <el-form-item :label="$t('aiShopping.model')">
          <el-input v-model="config.model" placeholder="gpt-3.5-turbo" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('aiShopping.assistantSettings') }}</el-divider>
        <el-form-item :label="$t('aiShopping.enableAssistant')">
          <el-switch v-model="config.enabled" />
        </el-form-item>
        <el-form-item :label="$t('aiShopping.welcomeMessage')">
          <el-input v-model="config.welcome_message" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="$t('aiShopping.systemPrompt')">
          <el-input v-model="config.system_prompt" type="textarea" :rows="4" :placeholder="$t('aiShopping.systemPromptPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('aiShopping.maxRecommendations')">
          <el-input-number v-model="config.max_recommendations" :min="1" :max="20" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('aiShopping.catalogDisplay') }}</el-divider>
        <el-form-item :label="$t('aiShopping.showCatalog')">
          <el-switch v-model="config.catalog_enabled" />
        </el-form-item>
        <el-form-item :label="$t('aiShopping.catalogLayout')">
          <el-radio-group v-model="config.catalog_layout">
            <el-radio value="grid">{{ $t('aiShopping.grid') }}</el-radio>
            <el-radio value="list">{{ $t('aiShopping.list') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('aiShopping.showFilters')">
          <el-switch v-model="config.show_filters" />
        </el-form-item>
        <el-form-item :label="$t('aiShopping.showPrice')">
          <el-switch v-model="config.show_price" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveConfig">{{ $t('aiShopping.saveConfig') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const config = reactive<any>({
  api_endpoint: '', api_key: '', model: '',
  enabled: false, welcome_message: '你好！我是购物助手，有什么可以帮您？',
  system_prompt: '你是一个专业的购物助手，帮助用户选择合适的商品。根据用户需求推荐产品，并提供详细的产品信息。',
  max_recommendations: 5, catalog_enabled: true, catalog_layout: 'grid',
  show_filters: true, show_price: true,
})

const loadConfig = async () => {
  const res = await request.get({ url: '/api/admin/ai-shopping/config' })
  if (res) Object.assign(config, res)
}

const saveConfig = async () => {
  await request.put({ url: '/api/admin/ai-shopping/config', params: config })
  ElMessage.success($t('common.saveSuccess'))
}

onMounted(loadConfig)
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
</style>
