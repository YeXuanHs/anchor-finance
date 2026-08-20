<template>
  <div class="ai-ticket-admin">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('aiAutoReply.title') }}</span>
          <el-tag type="info" size="small">mianyu_ai_ticket</el-tag>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <!-- 控制台配置 -->
        <el-tab-pane :label="$t('aiAutoReply.dashboard')" name="dashboard">
          <el-form :model="dashboard" label-width="140px" style="max-width: 700px">
            <el-divider content-position="left">{{ $t('aiAutoReply.modelConfig') }}</el-divider>
            <el-form-item :label="$t('aiAutoReply.apiEndpoint')">
              <el-input v-model="dashboard.api_endpoint" placeholder="https://api.openai.com/v1" />
            </el-form-item>
            <el-form-item label="API Key">
              <el-input v-model="dashboard.api_key" type="password" show-password />
            </el-form-item>
            <el-form-item :label="$t('aiAutoReply.model')">
              <el-input v-model="dashboard.model" placeholder="gpt-3.5-turbo" />
            </el-form-item>
            <el-divider content-position="left">{{ $t('aiAutoReply.autoReplySettings') }}</el-divider>
            <el-form-item :label="$t('aiAutoReply.enableAutoReply')">
              <el-switch v-model="dashboard.auto_reply_enabled" />
            </el-form-item>
            <el-form-item :label="$t('aiAutoReply.confidenceThreshold')">
              <el-slider v-model="dashboard.confidence_threshold" :min="0" :max="100" :format-tooltip="(v: number) => `${v}%`" />
            </el-form-item>
            <el-form-item :label="$t('aiAutoReply.maxDailyReplies')">
              <el-input-number v-model="dashboard.max_daily_replies" :min="0" :max="9999" />
            </el-form-item>
            <el-form-item :label="$t('aiAutoReply.transferAfterReply')">
              <el-switch v-model="dashboard.transfer_after_reply" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveDashboard">{{ $t('aiAutoReply.saveConfig') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 知识库 -->
        <el-tab-pane :label="$t('aiAutoReply.knowledge')" name="knowledge">
          <div style="margin-bottom: 16px; display: flex; justify-content: space-between">
            <el-input v-model="kbSearch" :placeholder="$t('aiAutoReply.searchKnowledge')" style="width: 300px" @keyup.enter="loadKnowledge" clearable />
            <div>
              <el-button @click="importDefault">{{ $t('aiAutoReply.importDefault') }}</el-button>
              <el-button type="primary" @click="showKBDialog = true">{{ $t('aiAutoReply.addEntry') }}</el-button>
            </div>
          </div>

          <el-table :data="knowledgeList" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="category" :label="$t('aiAutoReply.category')" width="120" />
            <el-table-column prop="question" :label="$t('aiAutoReply.question')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="hit_count" :label="$t('aiAutoReply.hitCount')" width="90" />
            <el-table-column prop="status" :label="$t('common.status')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.action')" width="140" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="editKB(row)">{{ $t('common.edit') }}</el-button>
                <el-button size="small" type="danger" @click="deleteKB(row)">{{ $t('common.delete') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 规则 -->
        <el-tab-pane :label="$t('aiAutoReply.automationRules')" name="rules">
          <div style="margin-bottom: 16px">
            <el-button type="primary" @click="showRuleDialog = true">{{ $t('aiAutoReply.addRule') }}</el-button>
          </div>
          <el-table :data="rules" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="name" :label="$t('aiAutoReply.ruleName')" width="150" />
            <el-table-column prop="match_type" :label="$t('aiAutoReply.matchType')" width="100">
              <template #default="{ row }">{{ { keyword: $t('aiAutoReply.keyword'), regex: $t('aiAutoReply.regex'), ai: $t('aiAutoReply.aiJudge') }[row.match_type as string] || row.match_type }}</template>
            </el-table-column>
            <el-table-column prop="match_value" :label="$t('aiAutoReply.matchValue')" min-width="150" show-overflow-tooltip />
            <el-table-column prop="action" :label="$t('aiAutoReply.action')" width="100">
              <template #default="{ row }">{{ { auto_reply: $t('aiAutoReply.autoReply'), transfer: $t('aiAutoReply.transfer'), close: $t('aiAutoReply.closeTicket'), tag: $t('aiAutoReply.tag') }[row.action as string] || row.action }}</template>
            </el-table-column>
            <el-table-column prop="priority" :label="$t('aiAutoReply.priority')" width="80" />
            <el-table-column prop="status" :label="$t('common.status')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.action')" width="140" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="editRule(row)">{{ $t('common.edit') }}</el-button>
                <el-button size="small" type="danger" @click="deleteRule(row)">{{ $t('common.delete') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 队列 -->
        <el-tab-pane :label="$t('aiAutoReply.processQueue')" name="queue">
          <el-row :gutter="20" style="margin-bottom: 16px">
            <el-col :span="6"><el-statistic :title="$t('aiAutoReply.pending')" :value="queueStats.pending || 0" /></el-col>
            <el-col :span="6"><el-statistic :title="$t('aiAutoReply.processing')" :value="queueStats.processing || 0" /></el-col>
            <el-col :span="6"><el-statistic :title="$t('aiAutoReply.completed')" :value="queueStats.completed || 0" /></el-col>
            <el-col :span="6"><el-statistic :title="$t('aiAutoReply.failed')" :value="queueStats.failed || 0" /></el-col>
          </el-row>
          <el-table :data="queue" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="ticket_id" :label="$t('aiAutoReply.ticketId')" width="80" />
            <el-table-column prop="status" :label="$t('common.status')" width="80">
              <template #default="{ row }">
                <el-tag :type="{ pending: 'warning', processing: 'primary', completed: 'success', failed: 'danger' }[row.status as string] as any" size="small">
                  {{ { pending: $t('aiAutoReply.pending'), processing: $t('aiAutoReply.processing'), completed: $t('aiAutoReply.completed'), failed: $t('aiAutoReply.failed') }[row.status as string] || row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="retry_count" :label="$t('aiAutoReply.retryCount')" width="90" />
            <el-table-column prop="error_message" :label="$t('aiAutoReply.errorMessage')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="created_at" :label="$t('aiAutoReply.createdAt')" width="160" />
          </el-table>
        </el-tab-pane>

        <!-- 日志 -->
        <el-tab-pane :label="$t('aiAutoReply.processLogs')" name="logs">
          <el-table :data="processLogs" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="ticket_id" :label="$t('aiAutoReply.ticketId')" width="80" />
            <el-table-column prop="action" :label="$t('aiAutoReply.action')" width="100" />
            <el-table-column prop="ai_response" :label="$t('aiAutoReply.aiResponse')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="confidence" :label="$t('aiAutoReply.confidence')" width="80">
              <template #default="{ row }">{{ row.confidence }}%</template>
            </el-table-column>
            <el-table-column prop="status" :label="$t('common.status')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('aiAutoReply.time')" width="160" />
          </el-table>
          <el-pagination
            v-model:current-page="logPage"
            :page-size="20"
            :total="logTotal"
            layout="total, prev, pager, next"
            style="margin-top: 16px; justify-content: flex-end"
            @current-change="loadProcessLogs"
          />
        </el-tab-pane>

        <!-- 测试 -->
        <el-tab-pane :label="$t('aiAutoReply.test')" name="test">
          <el-form :model="testForm" label-width="100px" style="max-width: 600px">
            <el-form-item :label="$t('aiAutoReply.ticketSubject')">
              <el-input v-model="testForm.subject" />
            </el-form-item>
            <el-form-item :label="$t('aiAutoReply.ticketContent')">
              <el-input v-model="testForm.content" type="textarea" :rows="4" />
            </el-form-item>
            <el-form-item :label="$t('aiAutoReply.deptId')">
              <el-input-number v-model="testForm.dept_id" :min="1" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testAutoReply">{{ $t('aiAutoReply.testAutoReply') }}</el-button>
            </el-form-item>
            <el-form-item v-if="testResult" :label="$t('aiAutoReply.testResult')">
              <el-card shadow="never">
                <pre style="white-space: pre-wrap; margin: 0">{{ JSON.stringify(testResult, null, 2) }}</pre>
              </el-card>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 知识库编辑弹窗 -->
    <el-dialog v-model="showKBDialog" :title="editingKB ? $t('aiAutoReply.editKnowledgeEntry') : $t('aiAutoReply.addKnowledgeEntry')" width="600px">
      <el-form :model="kbForm" label-width="100px">
        <el-form-item :label="$t('aiAutoReply.category')">
          <el-input v-model="kbForm.category" :placeholder="$t('aiAutoReply.categoryPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.question')">
          <el-input v-model="kbForm.question" />
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.answer')">
          <el-input v-model="kbForm.answer" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.keywordsPlaceholder')">
          <el-input v-model="kbForm.keywords" :placeholder="$t('aiAutoReply.commaSeparated')" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="kbForm.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showKBDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveKB">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- 规则编辑弹窗 -->
    <el-dialog v-model="showRuleDialog" :title="editingRule ? $t('aiAutoReply.editRule') : $t('aiAutoReply.addRuleTitle')" width="600px">
      <el-form :model="ruleForm" label-width="100px">
        <el-form-item :label="$t('aiAutoReply.ruleName')">
          <el-input v-model="ruleForm.name" />
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.matchType')">
          <el-select v-model="ruleForm.match_type">
            <el-option :label="$t('aiAutoReply.keyword')" value="keyword" />
            <el-option :label="$t('aiAutoReply.regex')" value="regex" />
            <el-option :label="$t('aiAutoReply.aiJudge')" value="ai" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.matchValue')">
          <el-input v-model="ruleForm.match_value" :placeholder="$t('aiAutoReply.matchValuePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.action')">
          <el-select v-model="ruleForm.action">
            <el-option :label="$t('aiAutoReply.autoReply')" value="auto_reply" />
            <el-option :label="$t('aiAutoReply.transfer')" value="transfer" />
            <el-option :label="$t('aiAutoReply.closeTicket')" value="close" />
            <el-option :label="$t('aiAutoReply.tag')" value="tag" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.replyContent')" v-if="ruleForm.action === 'auto_reply'">
          <el-input v-model="ruleForm.reply_content" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item :label="$t('aiAutoReply.priority')">
          <el-input-number v-model="ruleForm.priority" :min="0" :max="100" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="ruleForm.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRuleDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveRule">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const activeTab = ref('dashboard')

// 控制台配置
const dashboard = reactive<any>({ api_endpoint: '', api_key: '', model: '', auto_reply_enabled: false, confidence_threshold: 70, max_daily_replies: 100, transfer_after_reply: false })

// 知识库
const knowledgeList = ref<any[]>([])
const kbSearch = ref('')
const showKBDialog = ref(false)
const editingKB = ref<any>(null)
const kbForm = reactive<any>({ category: '', question: '', answer: '', keywords: '', status: 1 })

// 规则
const rules = ref<any[]>([])
const showRuleDialog = ref(false)
const editingRule = ref<any>(null)
const ruleForm = reactive<any>({ name: '', match_type: 'keyword', match_value: '', action: 'auto_reply', reply_content: '', priority: 0, status: 1 })

// 队列
const queue = ref<any[]>([])
const queueStats = reactive<any>({})

// 日志
const processLogs = ref<any[]>([])
const logPage = ref(1)
const logTotal = ref(0)

// 测试
const testForm = reactive({ subject: '', content: '', dept_id: 1 })
const testResult = ref<any>(null)

const loadDashboard = async () => {
  const res = await request.get({ url: '/api/admin/ai-ticket/dashboard' })
  if (res) Object.assign(dashboard, res)
}
const saveDashboard = async () => {
  await request.put({ url: '/api/admin/ai-ticket/dashboard', params: dashboard })
  ElMessage.success($t('common.saveSuccess'))
}

const loadKnowledge = async () => {
  const res = await request.get({ url: '/api/admin/ai-ticket/knowledge', params: { keyword: kbSearch.value } })
  knowledgeList.value = Array.isArray(res) ? res : []
}
const editKB = (row: any) => { editingKB.value = row; Object.assign(kbForm, row); showKBDialog.value = true }
const saveKB = async () => {
  if (editingKB.value) {
    await request.put({ url: `/api/admin/ai-ticket/knowledge/${editingKB.value.id}`, params: kbForm })
  } else {
    await request.post({ url: '/api/admin/ai-ticket/knowledge', params: kbForm })
  }
  ElMessage.success($t('common.saveSuccess'))
  showKBDialog.value = false
  editingKB.value = null
  Object.assign(kbForm, { category: '', question: '', answer: '', keywords: '', status: 1 })
  loadKnowledge()
}
const deleteKB = async (row: any) => {
  await ElMessageBox.confirm($t('common.confirmDelete'), $t('common.tips'))
  await request.del({ url: `/api/admin/ai-ticket/knowledge/${row.id}` })
  ElMessage.success($t('common.deleteSuccess'))
  loadKnowledge()
}
const importDefault = async () => {
  await request.post({ url: '/api/admin/ai-ticket/knowledge/import' })
  ElMessage.success($t('aiAutoReply.importComplete'))
  loadKnowledge()
}

const loadRules = async () => {
  const res = await request.get({ url: '/api/admin/ai-ticket/rules' })
  rules.value = Array.isArray(res) ? res : []
}
const editRule = (row: any) => { editingRule.value = row; Object.assign(ruleForm, row); showRuleDialog.value = true }
const saveRule = async () => {
  if (editingRule.value) {
    await request.put({ url: `/api/admin/ai-ticket/rules/${editingRule.value.id}`, params: ruleForm })
  } else {
    await request.post({ url: '/api/admin/ai-ticket/rules', params: ruleForm })
  }
  ElMessage.success($t('common.saveSuccess'))
  showRuleDialog.value = false
  editingRule.value = null
  Object.assign(ruleForm, { name: '', match_type: 'keyword', match_value: '', action: 'auto_reply', reply_content: '', priority: 0, status: 1 })
  loadRules()
}
const deleteRule = async (row: any) => {
  await ElMessageBox.confirm($t('common.confirmDelete'), $t('common.tips'))
  await request.del({ url: `/api/admin/ai-ticket/rules/${row.id}` })
  ElMessage.success($t('common.deleteSuccess'))
  loadRules()
}

const loadQueue = async () => {
  const res = await request.get({ url: '/api/admin/ai-ticket/queue' })
  queue.value = res?.items || []
  const stats = await request.get({ url: '/api/admin/ai-ticket/queue/stats' })
  if (stats) Object.assign(queueStats, stats)
}

const loadProcessLogs = async () => {
  const res = await request.get({ url: '/api/admin/ai-ticket/process-logs', params: { page: logPage.value } })
  processLogs.value = res?.items || []
  logTotal.value = res?.total || 0
}

const testAutoReply = async () => {
  const res = await request.post({ url: '/api/admin/ai-ticket/test', params: testForm })
  testResult.value = res
}

onMounted(() => {
  loadDashboard()
  loadKnowledge()
  loadRules()
  loadQueue()
  loadProcessLogs()
})
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
</style>
