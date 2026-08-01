<template>
  <div class="ai-ticket-page">
    <div class="page-header">
      <div class="header-left">
        <h2>AI 工单自动回复</h2>
        <span class="subtitle">配置 AI 自动回复引擎和自动回复规则</span>
      </div>
    </div>

    <el-tabs v-model="activeTab">
      <!-- AI 模型配置 -->
      <el-tab-pane label="AI 模型配置" name="models">
        <div class="tab-header">
          <el-button type="primary" @click="showAddModel">添加 AI 模型</el-button>
        </div>
        <el-table :data="aiModels" v-loading="loadingModels" stripe>
          <el-table-column prop="provider" label="供应商" width="120" />
          <el-table-column prop="model" label="模型" width="200" />
          <el-table-column prop="api_endpoint" label="API地址" min-width="300" show-overflow-tooltip />
          <el-table-column prop="is_active" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
                {{ row.is_active ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" align="center">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="showEditModel(row)">编辑</el-button>
              <el-popconfirm title="确定删除？" @confirm="deleteModel(row.id)">
                <template #reference>
                  <el-button type="danger" link size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <el-dialog v-model="modelDialogVisible" :title="isEditModel ? '编辑 AI 模型' : '添加 AI 模型'" width="600px">
          <el-form :model="modelForm" label-width="120px">
            <el-form-item label="供应商" required>
              <el-select v-model="modelForm.provider">
                <el-option label="OpenAI" value="openai" />
                <el-option label="Claude" value="claude" />
                <el-option label="DeepSeek" value="deepseek" />
                <el-option label="自定义" value="custom" />
              </el-select>
            </el-form-item>
            <el-form-item label="模型" required>
              <el-input v-model="modelForm.model" placeholder="例如: gpt-4o, deepseek-chat" />
            </el-form-item>
            <el-form-item label="API 地址">
              <el-input v-model="modelForm.api_endpoint" placeholder="留空使用默认地址" />
            </el-form-item>
            <el-form-item label="API Key" required>
              <el-input v-model="modelForm.api_key" type="password" show-password placeholder="API 密钥" />
            </el-form-item>
            <el-form-item label="最大Token">
              <el-input-number v-model="modelForm.max_tokens" :min="100" :max="8000" />
            </el-form-item>
            <el-form-item label="温度">
              <el-slider v-model="modelForm.temperature" :min="0" :max="2" :step="0.1" show-input />
            </el-form-item>
            <el-form-item label="系统提示词">
              <el-input v-model="modelForm.system_prompt" type="textarea" :rows="4" placeholder="AI 的角色设定和行为指令" />
            </el-form-item>
            <el-form-item label="启用">
              <el-switch v-model="modelForm.is_active" />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="modelDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitModel">确定</el-button>
          </template>
        </el-dialog>
      </el-tab-pane>

      <!-- 自动回复配置 -->
      <el-tab-pane label="自动回复配置" name="config">
        <el-card>
          <el-form :model="autoReplyConfig" label-width="160px" style="max-width: 700px">
            <el-form-item label="启用自动回复">
              <el-switch v-model="autoReplyConfig.enabled" />
            </el-form-item>
            <el-form-item label="AI 模型">
              <el-select v-model="autoReplyConfig.ai_config_id" placeholder="选择 AI 模型">
                <el-option v-for="m in aiModels" :key="m.id" :label="`${m.provider} - ${m.model}`" :value="m.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="置信度阈值">
              <el-slider v-model="autoReplyConfig.confidence_threshold" :min="0.3" :max="1" :step="0.05" show-input />
              <div class="form-hint">低于此置信度的回复不会自动发送</div>
            </el-form-item>
            <el-form-item label="最大自动回复数">
              <el-input-number v-model="autoReplyConfig.max_auto_replies" :min="1" :max="10" />
              <div class="form-hint">同一工单最多自动回复次数</div>
            </el-form-item>
            <el-form-item label="回复延迟(秒)">
              <el-input-number v-model="autoReplyConfig.reply_delay" :min="0" :max="60" />
            </el-form-item>
            <el-form-item label="引用知识库">
              <el-switch v-model="autoReplyConfig.include_kb_content" />
            </el-form-item>
            <el-form-item label="知识库搜索数量" v-if="autoReplyConfig.include_kb_content">
              <el-input-number v-model="autoReplyConfig.kb_search_limit" :min="1" :max="20" />
            </el-form-item>
            <el-form-item label="适用部门ID">
              <el-input v-model="autoReplyConfig.dept_ids" placeholder="逗号分隔，留空表示全部" />
            </el-form-item>
            <el-form-item label="排除关键词">
              <el-input v-model="autoReplyConfig.exclude_keywords" placeholder="包含这些关键词的工单不自动回复，逗号分隔" />
            </el-form-item>
            <el-form-item label="添加AI声明">
              <el-switch v-model="autoReplyConfig.add_disclaimer" />
            </el-form-item>
            <el-form-item label="声明文本" v-if="autoReplyConfig.add_disclaimer">
              <el-input v-model="autoReplyConfig.disclaimer_text" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveAutoReplyConfig" :loading="savingConfig">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 自动回复日志 -->
      <el-tab-pane label="回复日志" name="logs">
        <el-table :data="logs" v-loading="loadingLogs" stripe>
          <el-table-column prop="ticket_id" label="工单ID" width="100" />
          <el-table-column prop="question" label="问题" min-width="200" show-overflow-tooltip />
          <el-table-column prop="answer" label="AI回复" min-width="300" show-overflow-tooltip />
          <el-table-column prop="confidence" label="置信度" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="row.confidence >= 0.7 ? 'success' : 'warning'" size="small">
                {{ (row.confidence * 100).toFixed(0) }}%
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="tokens_used" label="Token" width="80" align="center" />
          <el-table-column prop="accepted" label="接受" width="80" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.accepted === true" type="success" size="small">是</el-tag>
              <el-tag v-else-if="row.accepted === false" type="danger" size="small">否</el-tag>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="170" />
        </el-table>
      </el-tab-pane>

      <!-- 测试 -->
      <el-tab-pane label="测试" name="test">
        <el-card>
          <el-form :model="testForm" label-width="100px" style="max-width: 600px">
            <el-form-item label="工单主题">
              <el-input v-model="testForm.subject" placeholder="测试工单主题" />
            </el-form-item>
            <el-form-item label="工单内容">
              <el-input v-model="testForm.content" type="textarea" :rows="4" placeholder="测试工单内容" />
            </el-form-item>
            <el-form-item label="部门ID">
              <el-input-number v-model="testForm.dept_id" :min="0" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="runTest" :loading="testing">测试自动回复</el-button>
            </el-form-item>
          </el-form>
          <el-divider v-if="testResult" />
          <div v-if="testResult">
            <h4>测试结果</h4>
            <p><strong>是否回复：</strong>{{ testResult.should_reply ? '是' : '否' }}</p>
            <p><strong>置信度：</strong>{{ (testResult.confidence * 100).toFixed(1) }}%</p>
            <div v-if="testResult.reply">
              <h4>AI 回复：</h4>
              <div class="test-reply" v-html="testResult.reply"></div>
            </div>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const activeTab = ref('config')
const loadingModels = ref(false)
const loadingLogs = ref(false)
const savingConfig = ref(false)
const testing = ref(false)
const aiModels = ref<any[]>([])
const logs = ref<any[]>([])

const modelDialogVisible = ref(false)
const isEditModel = ref(false)
const editModelId = ref(0)
const modelForm = ref({
  provider: 'openai', model: '', api_endpoint: '', api_key: '',
  max_tokens: 2000, temperature: 0.7, system_prompt: '', is_active: true
})

const autoReplyConfig = ref({
  enabled: false, ai_config_id: null as number | null,
  confidence_threshold: 0.7, max_auto_replies: 3, reply_delay: 5,
  include_kb_content: true, kb_search_limit: 5,
  dept_ids: '', exclude_keywords: '',
  add_disclaimer: true, disclaimer_text: '此回复由AI生成，仅供参考。如需人工帮助请回复「转人工」'
})

const testForm = ref({ subject: '', content: '', dept_id: 0 })
const testResult = ref<any>(null)

const loadModels = async () => {
  loadingModels.value = true
  try {
    const res = await request.get('/api/admin/ai-ticket/configs')
    aiModels.value = res.data || []
  } catch (e) {}
  finally { loadingModels.value = false }
}

const loadAutoReplyConfig = async () => {
  try {
    const res = await request.get('/api/admin/ai-ticket/auto-reply-config')
    if (res.data) autoReplyConfig.value = { ...autoReplyConfig.value, ...res.data }
  } catch (e) {}
}

const loadLogs = async () => {
  loadingLogs.value = true
  try {
    const res = await request.get('/api/admin/ai-ticket/logs')
    logs.value = res.data?.items || []
  } catch (e) {}
  finally { loadingLogs.value = false }
}

const showAddModel = () => {
  isEditModel.value = false
  modelForm.value = { provider: 'openai', model: '', api_endpoint: '', api_key: '', max_tokens: 2000, temperature: 0.7, system_prompt: '', is_active: true }
  modelDialogVisible.value = true
}

const showEditModel = (row: any) => {
  isEditModel.value = true
  editModelId.value = row.id
  modelForm.value = { ...row, api_key: '' }
  modelDialogVisible.value = true
}

const submitModel = async () => {
  if (!modelForm.value.model) { ElMessage.warning('请输入模型名称'); return }
  try {
    const data = { ...modelForm.value }
    if (isEditModel.value && !data.api_key) delete data.api_key
    if (isEditModel.value) {
      await request.put(`/api/admin/ai-ticket/configs/${editModelId.value}`, data)
    } else {
      await request.post('/api/admin/ai-ticket/configs', data)
    }
    ElMessage.success('操作成功')
    modelDialogVisible.value = false
    loadModels()
  } catch (e) { ElMessage.error('操作失败') }
}

const deleteModel = async (id: number) => {
  try {
    await request.delete(`/api/admin/ai-ticket/configs/${id}`)
    ElMessage.success('删除成功')
    loadModels()
  } catch (e) { ElMessage.error('删除失败') }
}

const saveAutoReplyConfig = async () => {
  savingConfig.value = true
  try {
    await request.put('/api/admin/ai-ticket/auto-reply-config', autoReplyConfig.value)
    ElMessage.success('保存成功')
  } catch (e) { ElMessage.error('保存失败') }
  finally { savingConfig.value = false }
}

const runTest = async () => {
  if (!testForm.value.subject && !testForm.value.content) { ElMessage.warning('请输入测试内容'); return }
  testing.value = true
  testResult.value = null
  try {
    const res = await request.post('/api/admin/ai-ticket/test', testForm.value)
    testResult.value = res.data
  } catch (e) { ElMessage.error('测试失败') }
  finally { testing.value = false }
}

onMounted(() => {
  loadModels()
  loadAutoReplyConfig()
  loadLogs()
})
</script>

<style scoped>
.ai-ticket-page { padding: 20px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { margin: 0 0 4px 0; font-size: 20px; }
.subtitle { color: #909399; font-size: 14px; }
.tab-header { display: flex; gap: 12px; align-items: center; margin-bottom: 16px; }
.form-hint { color: #909399; font-size: 12px; margin-top: 4px; }
.test-reply { background: #f5f7fa; padding: 16px; border-radius: 8px; white-space: pre-wrap; }
</style>
