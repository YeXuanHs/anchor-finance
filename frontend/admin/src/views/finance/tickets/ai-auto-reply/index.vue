<template>
  <div class="ai-ticket-admin">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>AI 工单自动回复</span>
          <el-tag type="info" size="small">mianyu_ai_ticket</el-tag>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <!-- 控制台配置 -->
        <el-tab-pane label="控制台配置" name="dashboard">
          <el-form :model="dashboard" label-width="140px" style="max-width: 700px">
            <el-divider content-position="left">AI 模型配置</el-divider>
            <el-form-item label="API 地址">
              <el-input v-model="dashboard.api_endpoint" placeholder="https://api.openai.com/v1" />
            </el-form-item>
            <el-form-item label="API Key">
              <el-input v-model="dashboard.api_key" type="password" show-password />
            </el-form-item>
            <el-form-item label="模型">
              <el-input v-model="dashboard.model" placeholder="gpt-3.5-turbo" />
            </el-form-item>
            <el-divider content-position="left">自动回复设置</el-divider>
            <el-form-item label="启用自动回复">
              <el-switch v-model="dashboard.auto_reply_enabled" />
            </el-form-item>
            <el-form-item label="置信度阈值">
              <el-slider v-model="dashboard.confidence_threshold" :min="0" :max="100" :format-tooltip="(v: number) => `${v}%`" />
            </el-form-item>
            <el-form-item label="每日最大回复数">
              <el-input-number v-model="dashboard.max_daily_replies" :min="0" :max="9999" />
            </el-form-item>
            <el-form-item label="回复后转人工">
              <el-switch v-model="dashboard.transfer_after_reply" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveDashboard">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 知识库 -->
        <el-tab-pane label="知识库" name="knowledge">
          <div style="margin-bottom: 16px; display: flex; justify-content: space-between">
            <el-input v-model="kbSearch" placeholder="搜索知识库..." style="width: 300px" @keyup.enter="loadKnowledge" clearable />
            <div>
              <el-button @click="importDefault">导入默认知识库</el-button>
              <el-button type="primary" @click="showKBDialog = true">新增条目</el-button>
            </div>
          </div>

          <el-table :data="knowledgeList" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="category" label="分类" width="120" />
            <el-table-column prop="question" label="问题" min-width="200" show-overflow-tooltip />
            <el-table-column prop="hit_count" label="命中次数" width="90" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="editKB(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteKB(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 规则 -->
        <el-tab-pane label="自动化规则" name="rules">
          <div style="margin-bottom: 16px">
            <el-button type="primary" @click="showRuleDialog = true">新增规则</el-button>
          </div>
          <el-table :data="rules" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="name" label="规则名称" width="150" />
            <el-table-column prop="match_type" label="匹配方式" width="100">
              <template #default="{ row }">{{ { keyword: '关键词', regex: '正则', ai: 'AI判断' }[row.match_type as string] || row.match_type }}</template>
            </el-table-column>
            <el-table-column prop="match_value" label="匹配值" min-width="150" show-overflow-tooltip />
            <el-table-column prop="action" label="动作" width="100">
              <template #default="{ row }">{{ { auto_reply: '自动回复', transfer: '转人工', close: '关闭工单', tag: '打标签' }[row.action as string] || row.action }}</template>
            </el-table-column>
            <el-table-column prop="priority" label="优先级" width="80" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="editRule(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteRule(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 队列 -->
        <el-tab-pane label="处理队列" name="queue">
          <el-row :gutter="20" style="margin-bottom: 16px">
            <el-col :span="6"><el-statistic title="待处理" :value="queueStats.pending || 0" /></el-col>
            <el-col :span="6"><el-statistic title="处理中" :value="queueStats.processing || 0" /></el-col>
            <el-col :span="6"><el-statistic title="已完成" :value="queueStats.completed || 0" /></el-col>
            <el-col :span="6"><el-statistic title="失败" :value="queueStats.failed || 0" /></el-col>
          </el-row>
          <el-table :data="queue" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="ticket_id" label="工单ID" width="80" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="{ pending: 'warning', processing: 'primary', completed: 'success', failed: 'danger' }[row.status as string] as any" size="small">
                  {{ { pending: '待处理', processing: '处理中', completed: '已完成', failed: '失败' }[row.status as string] || row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="retry_count" label="重试次数" width="90" />
            <el-table-column prop="error_message" label="错误信息" min-width="200" show-overflow-tooltip />
            <el-table-column prop="created_at" label="创建时间" width="160" />
          </el-table>
        </el-tab-pane>

        <!-- 日志 -->
        <el-tab-pane label="处理日志" name="logs">
          <el-table :data="processLogs" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="ticket_id" label="工单ID" width="80" />
            <el-table-column prop="action" label="动作" width="100" />
            <el-table-column prop="ai_response" label="AI回复" min-width="200" show-overflow-tooltip />
            <el-table-column prop="confidence" label="置信度" width="80">
              <template #default="{ row }">{{ row.confidence }}%</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="160" />
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
        <el-tab-pane label="测试" name="test">
          <el-form :model="testForm" label-width="100px" style="max-width: 600px">
            <el-form-item label="工单标题">
              <el-input v-model="testForm.subject" />
            </el-form-item>
            <el-form-item label="工单内容">
              <el-input v-model="testForm.content" type="textarea" :rows="4" />
            </el-form-item>
            <el-form-item label="部门ID">
              <el-input-number v-model="testForm.dept_id" :min="1" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testAutoReply">测试自动回复</el-button>
            </el-form-item>
            <el-form-item v-if="testResult" label="测试结果">
              <el-card shadow="never">
                <pre style="white-space: pre-wrap; margin: 0">{{ JSON.stringify(testResult, null, 2) }}</pre>
              </el-card>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 知识库编辑弹窗 -->
    <el-dialog v-model="showKBDialog" :title="editingKB ? '编辑知识条目' : '新增知识条目'" width="600px">
      <el-form :model="kbForm" label-width="100px">
        <el-form-item label="分类">
          <el-input v-model="kbForm.category" placeholder="如：账号、支付、产品" />
        </el-form-item>
        <el-form-item label="问题">
          <el-input v-model="kbForm.question" />
        </el-form-item>
        <el-form-item label="答案">
          <el-input v-model="kbForm.answer" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="kbForm.keywords" placeholder="逗号分隔" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="kbForm.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showKBDialog = false">取消</el-button>
        <el-button type="primary" @click="saveKB">保存</el-button>
      </template>
    </el-dialog>

    <!-- 规则编辑弹窗 -->
    <el-dialog v-model="showRuleDialog" :title="editingRule ? '编辑规则' : '新增规则'" width="600px">
      <el-form :model="ruleForm" label-width="100px">
        <el-form-item label="规则名称">
          <el-input v-model="ruleForm.name" />
        </el-form-item>
        <el-form-item label="匹配方式">
          <el-select v-model="ruleForm.match_type">
            <el-option label="关键词" value="keyword" />
            <el-option label="正则" value="regex" />
            <el-option label="AI判断" value="ai" />
          </el-select>
        </el-form-item>
        <el-form-item label="匹配值">
          <el-input v-model="ruleForm.match_value" placeholder="关键词或正则表达式" />
        </el-form-item>
        <el-form-item label="动作">
          <el-select v-model="ruleForm.action">
            <el-option label="自动回复" value="auto_reply" />
            <el-option label="转人工" value="transfer" />
            <el-option label="关闭工单" value="close" />
            <el-option label="打标签" value="tag" />
          </el-select>
        </el-form-item>
        <el-form-item label="回复内容" v-if="ruleForm.action === 'auto_reply'">
          <el-input v-model="ruleForm.reply_content" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="ruleForm.priority" :min="0" :max="100" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="ruleForm.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRuleDialog = false">取消</el-button>
        <el-button type="primary" @click="saveRule">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
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
  const { data } = await request.get('/admin/ai-ticket/dashboard')
  if (data) Object.assign(dashboard, data)
}
const saveDashboard = async () => {
  await request.put('/admin/ai-ticket/dashboard', dashboard)
  ElMessage.success('保存成功')
}

const loadKnowledge = async () => {
  const { data } = await request.get('/admin/ai-ticket/knowledge', { params: { keyword: kbSearch.value } })
  knowledgeList.value = Array.isArray(data) ? data : []
}
const editKB = (row: any) => { editingKB.value = row; Object.assign(kbForm, row); showKBDialog.value = true }
const saveKB = async () => {
  if (editingKB.value) {
    await request.put(`/admin/ai-ticket/knowledge/${editingKB.value.id}`, kbForm)
  } else {
    await request.post('/admin/ai-ticket/knowledge', kbForm)
  }
  ElMessage.success('保存成功')
  showKBDialog.value = false
  editingKB.value = null
  Object.assign(kbForm, { category: '', question: '', answer: '', keywords: '', status: 1 })
  loadKnowledge()
}
const deleteKB = async (row: any) => {
  await ElMessageBox.confirm('确定删除？', '提示')
  await request.delete(`/admin/ai-ticket/knowledge/${row.id}`)
  ElMessage.success('删除成功')
  loadKnowledge()
}
const importDefault = async () => {
  await request.post('/admin/ai-ticket/knowledge/import')
  ElMessage.success('导入完成')
  loadKnowledge()
}

const loadRules = async () => {
  const { data } = await request.get('/admin/ai-ticket/rules')
  rules.value = Array.isArray(data) ? data : []
}
const editRule = (row: any) => { editingRule.value = row; Object.assign(ruleForm, row); showRuleDialog.value = true }
const saveRule = async () => {
  if (editingRule.value) {
    await request.put(`/admin/ai-ticket/rules/${editingRule.value.id}`, ruleForm)
  } else {
    await request.post('/admin/ai-ticket/rules', ruleForm)
  }
  ElMessage.success('保存成功')
  showRuleDialog.value = false
  editingRule.value = null
  Object.assign(ruleForm, { name: '', match_type: 'keyword', match_value: '', action: 'auto_reply', reply_content: '', priority: 0, status: 1 })
  loadRules()
}
const deleteRule = async (row: any) => {
  await ElMessageBox.confirm('确定删除？', '提示')
  await request.delete(`/admin/ai-ticket/rules/${row.id}`)
  ElMessage.success('删除成功')
  loadRules()
}

const loadQueue = async () => {
  const { data } = await request.get('/admin/ai-ticket/queue')
  queue.value = data?.items || []
  const { data: stats } = await request.get('/admin/ai-ticket/queue/stats')
  if (stats) Object.assign(queueStats, stats)
}

const loadProcessLogs = async () => {
  const { data } = await request.get('/admin/ai-ticket/process-logs', { params: { page: logPage.value } })
  processLogs.value = data?.items || []
  logTotal.value = data?.total || 0
}

const testAutoReply = async () => {
  const { data } = await request.post('/admin/ai-ticket/test', testForm)
  testResult.value = data
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
