<template>
  <div class="ai-shopping-page">
    <div class="page-header">
      <div class="header-left">
        <h2>AI 购物助手配置</h2>
        <span class="subtitle">配置 AI 购物助手和商品目录展示方式</span>
      </div>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 助手配置 -->
      <el-tab-pane label="购物助手配置" name="assistant">
        <el-card>
          <el-form :model="assistantConfig" label-width="140px" style="max-width: 700px">
            <el-form-item label="启用购物助手">
              <el-switch v-model="assistantConfig.enabled" />
            </el-form-item>
            <el-form-item label="AI 模型">
              <el-select v-model="assistantConfig.ai_config_id" placeholder="选择 AI 模型">
                <el-option v-for="m in aiModels" :key="m.id" :label="`${m.provider} - ${m.model}`" :value="m.id" />
              </el-select>
              <div class="form-hint">需要先在 AI 工单自动回复 中配置 AI 模型</div>
            </el-form-item>
            <el-form-item label="欢迎语">
              <el-input v-model="assistantConfig.welcome_message" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item label="系统提示词">
              <el-input v-model="assistantConfig.system_prompt" type="textarea" :rows="5" placeholder="AI 购物助手的角色设定" />
            </el-form-item>
            <el-form-item label="最大推荐数">
              <el-input-number v-model="assistantConfig.max_recommendations" :min="1" :max="20" />
            </el-form-item>
            <el-form-item label="展示价格">
              <el-switch v-model="assistantConfig.include_pricing" />
            </el-form-item>
            <el-form-item label="所有页面显示">
              <el-switch v-model="assistantConfig.show_on_all_pages" />
            </el-form-item>
            <el-form-item label="触发关键词">
              <el-input v-model="assistantConfig.trigger_keywords" placeholder="逗号分隔，留空表示始终显示" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveAssistantConfig" :loading="saving">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 商品目录配置 -->
      <el-tab-pane label="商品目录配置" name="catalog">
        <el-card>
          <el-form :model="catalogConfig" label-width="140px" style="max-width: 700px">
            <el-form-item label="布局样式">
              <el-radio-group v-model="catalogConfig.layout_style">
                <el-radio value="grid">网格</el-radio>
                <el-radio value="list">列表</el-radio>
                <el-radio value="card">卡片</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="显示筛选器">
              <el-switch v-model="catalogConfig.show_filters" />
            </el-form-item>
            <el-form-item label="显示产品对比">
              <el-switch v-model="catalogConfig.show_comparison" />
            </el-form-item>
            <el-form-item label="显示用户评价">
              <el-switch v-model="catalogConfig.show_reviews" />
            </el-form-item>
            <el-form-item label="显示技术规格">
              <el-switch v-model="catalogConfig.show_tech_specs" />
            </el-form-item>
            <el-form-item label="启用排序">
              <el-switch v-model="catalogConfig.enable_sort" />
            </el-form-item>
            <el-form-item label="默认排序" v-if="catalogConfig.enable_sort">
              <el-select v-model="catalogConfig.default_sort">
                <el-option label="推荐" value="recommend" />
                <el-option label="价格从低到高" value="price_asc" />
                <el-option label="价格从高到低" value="price_desc" />
                <el-option label="最新" value="newest" />
                <el-option label="最热门" value="popular" />
              </el-select>
            </el-form-item>
            <el-form-item label="每页产品数">
              <el-input-number v-model="catalogConfig.products_per_page" :min="6" :max="60" />
            </el-form-item>
            <el-form-item label="显示分类树">
              <el-switch v-model="catalogConfig.show_category_tree" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveCatalogConfig" :loading="savingCatalog">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const activeTab = ref('assistant')
const saving = ref(false)
const savingCatalog = ref(false)
const aiModels = ref<any[]>([])

const assistantConfig = ref({
  enabled: false, ai_config_id: null as number | null,
  welcome_message: '您好！我是AI购物助手，可以帮您推荐合适的产品和服务。请问您有什么需求？',
  system_prompt: '', max_recommendations: 5,
  include_pricing: true, show_on_all_pages: true, trigger_keywords: ''
})

const catalogConfig = ref({
  layout_style: 'grid', show_filters: true, show_comparison: true,
  show_reviews: true, show_tech_specs: true, enable_sort: true,
  default_sort: 'recommend', products_per_page: 12, show_category_tree: true
})

const loadModels = async () => {
  try {
    const res = await request.get('/api/admin/ai-ticket/configs')
    aiModels.value = res.data || []
  } catch (e) {}
}

const loadAssistantConfig = async () => {
  try {
    const res = await request.get('/api/admin/ai-shopping/config')
    if (res.data) assistantConfig.value = { ...assistantConfig.value, ...res.data }
  } catch (e) {}
}

const loadCatalogConfig = async () => {
  try {
    const res = await request.get('/api/admin/ai-shopping/catalog-config')
    if (res.data) catalogConfig.value = { ...catalogConfig.value, ...res.data }
  } catch (e) {}
}

const saveAssistantConfig = async () => {
  saving.value = true
  try {
    await request.put('/api/admin/ai-shopping/config', assistantConfig.value)
    ElMessage.success('保存成功')
  } catch (e) { ElMessage.error('保存失败') }
  finally { saving.value = false }
}

const saveCatalogConfig = async () => {
  savingCatalog.value = true
  try {
    await request.put('/api/admin/ai-shopping/catalog-config', catalogConfig.value)
    ElMessage.success('保存成功')
  } catch (e) { ElMessage.error('保存失败') }
  finally { savingCatalog.value = false }
}

onMounted(() => {
  loadModels()
  loadAssistantConfig()
  loadCatalogConfig()
})
</script>

<style scoped>
.ai-shopping-page { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { margin: 0 0 4px 0; font-size: 20px; }
.subtitle { color: #909399; font-size: 14px; }
.form-hint { color: #909399; font-size: 12px; margin-top: 4px; }
</style>
