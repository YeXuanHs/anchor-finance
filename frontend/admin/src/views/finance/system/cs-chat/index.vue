<template>
  <div class="cs-chat-admin">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('csChat.title') }}</span>
          <el-tag type="info" size="small">anchor_cloud_finance_pro</el-tag>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('csChat.sessionManagement')" name="sessions">
          <el-form inline style="margin-bottom: 16px">
            <el-form-item :label="$t('csChat.status')">
              <el-select v-model="sessionFilter.status" clearable :placeholder="$t('csChat.all')" style="width: 120px">
                <el-option :label="$t('csChat.open')" value="open" />
                <el-option :label="$t('csChat.closed')" value="closed" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadSessions">{{ $t('csChat.query') }}</el-button>
            </el-form-item>
          </el-form>

          <el-table :data="sessions" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="client_name" :label="$t('csChat.visitorName')" width="120" />
            <el-table-column prop="mode" :label="$t('csChat.mode')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.mode === 'ai' ? 'success' : 'warning'" size="small">
                  {{ row.mode === 'ai' ? 'AI' : $t('csChat.human') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" :label="$t('csChat.status')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'open' ? 'primary' : 'info'" size="small">
                  {{ row.status === 'open' ? $t('csChat.open') : $t('csChat.closed') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="last_message" :label="$t('csChat.lastMessage')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="rating" :label="$t('csChat.rating')" width="100">
              <template #default="{ row }">
                <el-rate v-if="row.rating > 0" v-model="row.rating" disabled size="small" />
                <span v-else style="color: #999">{{ $t('csChat.notRated') }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('csChat.createdAt')" width="160" />
            <el-table-column :label="$t('csChat.actions')" width="140" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="openSession(row)">{{ $t('csChat.view') }}</el-button>
                <el-button size="small" type="danger" @click="closeSession(row)" v-if="row.status === 'open'">{{ $t('csChat.close') }}</el-button>
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

        <el-tab-pane :label="$t('csChat.aiConfig')" name="ai-config">
          <el-form :model="aiConfig" label-width="120px" style="max-width: 600px">
            <el-form-item :label="$t('csChat.apiEndpoint')">
              <el-input v-model="aiConfig.api_endpoint" placeholder="https://api.openai.com/v1" />
            </el-form-item>
            <el-form-item label="API Key">
              <el-input v-model="aiConfig.api_key" type="password" show-password placeholder="sk-..." />
            </el-form-item>
            <el-form-item :label="$t('csChat.model')">
              <el-input v-model="aiConfig.model" placeholder="gpt-3.5-turbo" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveAIConfig">{{ $t('csChat.save') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('csChat.appearance')" name="appearance">
          <el-form :model="appearance" label-width="120px" style="max-width: 600px">
            <el-form-item :label="$t('csChat.themeColor')">
              <el-color-picker v-model="appearance.theme_color" />
            </el-form-item>
            <el-form-item :label="$t('csChat.position')">
              <el-select v-model="appearance.position">
                <el-option :label="$t('csChat.bottomRight')" value="bottom-right" />
                <el-option :label="$t('csChat.bottomLeft')" value="bottom-left" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('csChat.welcomeMessage')">
              <el-input v-model="appearance.welcome_message" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item :label="$t('csChat.offlineMessage')">
              <el-input v-model="appearance.offline_message" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveAppearance">{{ $t('csChat.save') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('csChat.workingHours')" name="working-hours">
          <el-form :model="workingHours" label-width="120px" style="max-width: 600px">
            <el-form-item :label="$t('csChat.enableWorkingHours')">
              <el-switch v-model="workingHours.enabled" />
            </el-form-item>
            <el-form-item :label="$t('csChat.weekdays')">
              <el-select v-model="workingHours.weekdays" multiple :placeholder="$t('csChat.selectWeekdays')">
                <el-option v-for="d in weekdayOptions" :key="d.value" :label="d.label" :value="d.value" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('csChat.startTime')">
              <el-time-picker v-model="workingHours.start_time" format="HH:mm" value-format="HH:mm" />
            </el-form-item>
            <el-form-item :label="$t('csChat.endTime')">
              <el-time-picker v-model="workingHours.end_time" format="HH:mm" value-format="HH:mm" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveWorkingHours">{{ $t('csChat.save') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('csChat.stats')" name="stats">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-statistic :title="$t('csChat.totalSessions')" :value="stats.total_sessions || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic :title="$t('csChat.openSessions')" :value="stats.open_sessions || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic :title="$t('csChat.todaySessions')" :value="stats.today_sessions || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic :title="$t('csChat.avgRating')" :value="stats.avg_rating || 0" :precision="1" />
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="showChat" :title="$t('csChat.sessionDetail', { id: currentSession?.id })" width="700px" destroy-on-close>
      <div class="chat-messages" ref="chatBox">
        <div v-for="msg in messages" :key="msg.id" :class="['msg', msg.sender_type]">
          <div class="msg-bubble">
            <div class="msg-sender">{{ msg.sender_type === 'user' ? (currentSession?.client_name || $t('csChat.visitor')) : msg.sender_type === 'ai' ? $t('csChat.aiAssistant') : $t('csChat.staff') }}</div>
            <div class="msg-content">{{ msg.content }}</div>
            <div class="msg-time">{{ msg.created_at }}</div>
          </div>
        </div>
      </div>
      <div v-if="currentSession?.status === 'open'" style="margin-top: 16px; display: flex; gap: 8px">
        <el-input v-model="replyContent" :placeholder="$t('csChat.replyPlaceholder')" @keyup.enter="sendReply" />
        <el-button type="primary" @click="sendReply">{{ $t('csChat.send') }}</el-button>
        <el-button @click="transferToHuman" v-if="currentSession?.mode === 'ai'">{{ $t('csChat.transferToHuman') }}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const activeTab = ref('sessions')

// Sessions
const sessions = ref<any[]>([])
const sessionTotal = ref(0)
const sessionFilter = reactive({ status: '', page: 1 })
const showChat = ref(false)
const currentSession = ref<any>(null)
const messages = ref<any[]>([])
const replyContent = ref('')
const chatBox = ref<HTMLElement>()

// AI Config
const aiConfig = reactive({ api_endpoint: '', api_key: '', model: '' })

// Appearance
const appearance = reactive({ theme_color: '#409EFF', position: 'bottom-right', welcome_message: '', offline_message: '' })

// Working hours
const workingHours = reactive({ enabled: false, weekdays: ['monday','tuesday','wednesday','thursday','friday'], start_time: '09:00', end_time: '18:00' })

// Stats
const stats = reactive<any>({})

// Weekday options computed from i18n
const weekdayOptions = computed(() => [
  { label: $t('csChat.monday'), value: 'monday' },
  { label: $t('csChat.tuesday'), value: 'tuesday' },
  { label: $t('csChat.wednesday'), value: 'wednesday' },
  { label: $t('csChat.thursday'), value: 'thursday' },
  { label: $t('csChat.friday'), value: 'friday' },
  { label: $t('csChat.saturday'), value: 'saturday' },
  { label: $t('csChat.sunday'), value: 'sunday' }
])

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
  ElMessage.success($t('csChat.transferred'))
  openSession(currentSession.value)
}

const closeSession = async (row: any) => {
  await ElMessageBox.confirm($t('csChat.confirmClose'), $t('csChat.tip'))
  await request.post({ url: `/api/admin/cs/sessions/${row.id}/close` })
  ElMessage.success($t('csChat.closedSuccess'))
  loadSessions()
}

const loadAIConfig = async () => {
  const res = await request.get({ url: '/api/admin/cs/ai-config' })
  if (res) Object.assign(aiConfig, res)
}

const saveAIConfig = async () => {
  await request.put({ url: '/api/admin/cs/ai-config', params: aiConfig })
  ElMessage.success($t('csChat.saveSuccess'))
}

const loadAppearance = async () => {
  const res = await request.get({ url: '/api/admin/cs/appearance' })
  if (res) Object.assign(appearance, res)
}

const saveAppearance = async () => {
  await request.put({ url: '/api/admin/cs/appearance', params: appearance })
  ElMessage.success($t('csChat.saveSuccess'))
}

const loadWorkingHours = async () => {
  const res = await request.get({ url: '/api/admin/cs/working-hours' })
  if (res) Object.assign(workingHours, res)
}

const saveWorkingHours = async () => {
  await request.put({ url: '/api/admin/cs/working-hours', params: workingHours })
  ElMessage.success($t('csChat.saveSuccess'))
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
