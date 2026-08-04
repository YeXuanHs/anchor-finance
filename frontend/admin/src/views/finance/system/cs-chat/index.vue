<template>
  <div class="cs-chat-admin">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>客服聊天系统</span>
          <el-tag type="info" size="small">anchor_cloud_finance_pro</el-tag>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <!-- 会话列表 -->
        <el-tab-pane label="会话管理" name="sessions">
          <el-form inline style="margin-bottom: 16px">
            <el-form-item label="状态">
              <el-select v-model="sessionFilter.status" clearable placeholder="全部" style="width: 120px">
                <el-option label="进行中" value="open" />
                <el-option label="已关闭" value="closed" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadSessions">查询</el-button>
            </el-form-item>
          </el-form>

          <el-table :data="sessions" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="client_name" label="访客名称" width="120" />
            <el-table-column prop="mode" label="模式" width="80">
              <template #default="{ row }">
                <el-tag :type="row.mode === 'ai' ? 'success' : 'warning'" size="small">
                  {{ row.mode === 'ai' ? 'AI' : '人工' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'open' ? 'primary' : 'info'" size="small">
                  {{ row.status === 'open' ? '进行中' : '已关闭' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="last_message" label="最后消息" min-width="200" show-overflow-tooltip />
            <el-table-column prop="rating" label="评分" width="100">
              <template #default="{ row }">
                <el-rate v-if="row.rating > 0" v-model="row.rating" disabled size="small" />
                <span v-else style="color: #999">未评价</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160" />
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="openSession(row)">查看</el-button>
                <el-button size="small" type="danger" @click="closeSession(row)" v-if="row.status === 'open'">关闭</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="sessionFilter.page"
            :page-size="20"
            :total="sessionTotal"
            layout="total, prev, pager, next"
            style="margin-top: 16px; justify-content: flex-end"
            @current-change="loadSessions"
          />
        </el-tab-pane>

        <!-- AI配置 -->
        <el-tab-pane label="AI 配置" name="ai-config">
          <el-form :model="aiConfig" label-width="120px" style="max-width: 600px">
            <el-form-item label="API 地址">
              <el-input v-model="aiConfig.api_endpoint" placeholder="https://api.openai.com/v1" />
            </el-form-item>
            <el-form-item label="API Key">
              <el-input v-model="aiConfig.api_key" type="password" show-password placeholder="sk-..." />
            </el-form-item>
            <el-form-item label="模型">
              <el-input v-model="aiConfig.model" placeholder="gpt-3.5-turbo" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveAIConfig">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 外观设置 -->
        <el-tab-pane label="外观设置" name="appearance">
          <el-form :model="appearance" label-width="120px" style="max-width: 600px">
            <el-form-item label="主题色">
              <el-color-picker v-model="appearance.theme_color" />
            </el-form-item>
            <el-form-item label="位置">
              <el-select v-model="appearance.position">
                <el-option label="右下角" value="bottom-right" />
                <el-option label="左下角" value="bottom-left" />
              </el-select>
            </el-form-item>
            <el-form-item label="欢迎语">
              <el-input v-model="appearance.welcome_message" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="离线提示">
              <el-input v-model="appearance.offline_message" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveAppearance">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 工作时间 -->
        <el-tab-pane label="工作时间" name="working-hours">
          <el-form :model="workingHours" label-width="120px" style="max-width: 600px">
            <el-form-item label="启用工作时间">
              <el-switch v-model="workingHours.enabled" />
            </el-form-item>
            <el-form-item label="工作日">
              <el-select v-model="workingHours.weekdays" multiple placeholder="选择工作日">
                <el-option v-for="d in ['周一','周二','周三','周四','周五','周六','周日']" :key="d" :label="d" :value="d" />
              </el-select>
            </el-form-item>
            <el-form-item label="开始时间">
              <el-time-picker v-model="workingHours.start_time" format="HH:mm" value-format="HH:mm" />
            </el-form-item>
            <el-form-item label="结束时间">
              <el-time-picker v-model="workingHours.end_time" format="HH:mm" value-format="HH:mm" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveWorkingHours">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 统计 -->
        <el-tab-pane label="统计" name="stats">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-statistic title="总会话数" :value="stats.total_sessions || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="进行中" :value="stats.open_sessions || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="今日会话" :value="stats.today_sessions || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="平均评分" :value="stats.avg_rating || 0" :precision="1" />
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 会话详情弹窗 -->
    <el-dialog v-model="showChat" :title="`会话 #${currentSession?.id}`" width="700px" destroy-on-close>
      <div class="chat-messages" ref="chatBox">
        <div v-for="msg in messages" :key="msg.id" :class="['msg', msg.sender_type]">
          <div class="msg-bubble">
            <div class="msg-sender">{{ msg.sender_type === 'user' ? (currentSession?.client_name || '访客') : msg.sender_type === 'ai' ? 'AI助手' : '客服' }}</div>
            <div class="msg-content">{{ msg.content }}</div>
            <div class="msg-time">{{ msg.created_at }}</div>
          </div>
        </div>
      </div>
      <div v-if="currentSession?.status === 'open'" style="margin-top: 16px; display: flex; gap: 8px">
        <el-input v-model="replyContent" placeholder="输入回复..." @keyup.enter="sendReply" />
        <el-button type="primary" @click="sendReply">发送</el-button>
        <el-button @click="transferToHuman" v-if="currentSession?.mode === 'ai'">转人工</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const activeTab = ref('sessions')

// 会话
const sessions = ref<any[]>([])
const sessionTotal = ref(0)
const sessionFilter = reactive({ status: '', page: 1 })
const showChat = ref(false)
const currentSession = ref<any>(null)
const messages = ref<any[]>([])
const replyContent = ref('')
const chatBox = ref<HTMLElement>()

// AI配置
const aiConfig = reactive({ api_endpoint: '', api_key: '', model: '' })

// 外观
const appearance = reactive({ theme_color: '#409EFF', position: 'bottom-right', welcome_message: '你好，请问有什么可以帮您？', offline_message: '当前为非工作时间，请留言。' })

// 工作时间
const workingHours = reactive({ enabled: false, weekdays: ['周一','周二','周三','周四','周五'], start_time: '09:00', end_time: '18:00' })

// 统计
const stats = reactive<any>({})

const loadSessions = async () => {
  const res = await request.get({ url: '/api/admin/cs/sessions', params: sessionFilter })
  sessions.value = res?.items || []
  sessionTotal.value = res?.total || 0
}

const openSession = async (row: any) => {
  currentSession.value = row
  const res = await request.get({ url: `/api/admin/cs/sessions/${row.id}` })
  messages.value = res?.messages || []
  currentSession.value = res?.session || row
  showChat.value = true
  await nextTick()
  if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight
}

const sendReply = async () => {
  if (!replyContent.value.trim()) return
  await request.post({ url: `/api/admin/cs/sessions/${currentSession.value.id}/reply`, params: { content: replyContent.value } })
  replyContent.value = ''
  openSession(currentSession.value)
}

const transferToHuman = async () => {
  await request.post({ url: `/api/admin/cs/sessions/${currentSession.value.id}/transfer` })
  ElMessage.success('已转人工')
  openSession(currentSession.value)
}

const closeSession = async (row: any) => {
  await ElMessageBox.confirm('确定关闭此会话？', '提示')
  await request.post({ url: `/api/admin/cs/sessions/${row.id}/close` })
  ElMessage.success('已关闭')
  loadSessions()
}

const loadAIConfig = async () => {
  const res = await request.get({ url: '/api/admin/cs/ai-config' })
  if (res) Object.assign(aiConfig, res)
}

const saveAIConfig = async () => {
  await request.put({ url: '/api/admin/cs/ai-config', params: aiConfig })
  ElMessage.success('保存成功')
}

const loadAppearance = async () => {
  const res = await request.get({ url: '/api/admin/cs/appearance' })
  if (res) Object.assign(appearance, res)
}

const saveAppearance = async () => {
  await request.put({ url: '/api/admin/cs/appearance', params: appearance })
  ElMessage.success('保存成功')
}

const loadWorkingHours = async () => {
  const res = await request.get({ url: '/api/admin/cs/working-hours' })
  if (res) Object.assign(workingHours, res)
}

const saveWorkingHours = async () => {
  await request.put({ url: '/api/admin/cs/working-hours', params: workingHours })
  ElMessage.success('保存成功')
}

const loadStats = async () => {
  const res = await request.get({ url: '/api/admin/cs/stats' })
  if (res) Object.assign(stats, res)
}

onMounted(() => {
  loadSessions()
  loadAIConfig()
  loadAppearance()
  loadWorkingHours()
  loadStats()
})
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.chat-messages { max-height: 400px; overflow-y: auto; padding: 16px; background: #f5f5f5; border-radius: 8px; }
.msg { margin-bottom: 12px; display: flex; }
.msg.user { justify-content: flex-end; }
.msg.ai, .msg.staff { justify-content: flex-start; }
.msg-bubble { max-width: 70%; padding: 10px 14px; border-radius: 12px; background: #fff; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.msg.user .msg-bubble { background: #409EFF; color: #fff; }
.msg-sender { font-size: 12px; color: #999; margin-bottom: 4px; }
.msg.user .msg-sender { color: rgba(255,255,255,0.7); }
.msg-content { white-space: pre-wrap; word-break: break-all; }
.msg-time { font-size: 11px; color: #bbb; margin-top: 4px; text-align: right; }
</style>
