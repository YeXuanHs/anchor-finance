<template>
  <n-card :bordered="false" rounded>
    <n-tabs v-model:value="activeTab" type="line" animated>
      <!-- 通知渠道配置 -->
      <n-tab-pane name="channels" tab="通知渠道">
        <n-grid :cols="3" :x-gap="16" style="margin-top: 20px">
          <!-- 邮件渠道 -->
          <n-gi>
            <n-card title="邮件通知" size="small" :bordered="false" style="border-radius: 12px">
              <n-form :model="channels.email" label-placement="left" label-width="80">
                <n-form-item label="启用">
                  <n-switch v-model:value="channels.email.enabled" />
                </n-form-item>
                <template v-if="channels.email.enabled">
                  <n-form-item label="发件人">
                    <n-input v-model:value="channels.email.fromName" placeholder="锚点财务" />
                  </n-form-item>
                  <n-form-item label="回复邮箱">
                    <n-input v-model:value="channels.email.replyTo" placeholder="support@example.com" />
                  </n-form-item>
                </template>
              </n-form>
              <template #footer>
                <n-space justify="end">
                  <n-button size="small" @click="testChannel('email')">发送测试</n-button>
                  <n-button size="small" type="primary" @click="saveChannel('email')">保存</n-button>
                </n-space>
              </template>
            </n-card>
          </n-gi>

          <!-- 短信渠道 -->
          <n-gi>
            <n-card title="短信通知" size="small" :bordered="false" style="border-radius: 12px">
              <n-form :model="channels.sms" label-placement="left" label-width="80">
                <n-form-item label="启用">
                  <n-switch v-model:value="channels.sms.enabled" />
                </n-form-item>
                <template v-if="channels.sms.enabled">
                  <n-form-item label="签名">
                    <n-input v-model:value="channels.sms.signName" placeholder="短信签名" />
                  </n-form-item>
                  <n-form-item label="测试手机">
                    <n-input v-model:value="channels.sms.testPhone" placeholder="13800138000" />
                  </n-form-item>
                </template>
              </n-form>
              <template #footer>
                <n-space justify="end">
                  <n-button size="small" @click="testChannel('sms')">发送测试</n-button>
                  <n-button size="small" type="primary" @click="saveChannel('sms')">保存</n-button>
                </n-space>
              </template>
            </n-card>
          </n-gi>

          <!-- 站内信渠道 -->
          <n-gi>
            <n-card title="站内信通知" size="small" :bordered="false" style="border-radius: 12px">
              <n-form :model="channels.site" label-placement="left" label-width="80">
                <n-form-item label="启用">
                  <n-switch v-model:value="channels.site.enabled" />
                </n-form-item>
                <template v-if="channels.site.enabled">
                  <n-form-item label="保留天数">
                    <n-input-number v-model:value="channels.site.retainDays" :min="1" :max="365" style="width: 100%" />
                  </n-form-item>
                  <n-form-item label="最大条数">
                    <n-input-number v-model:value="channels.site.maxCount" :min="10" :max="1000" style="width: 100%" />
                  </n-form-item>
                </template>
              </n-form>
              <template #footer>
                <n-space justify="end">
                  <n-button size="small" type="primary" @click="saveChannel('site')">保存</n-button>
                </n-space>
              </template>
            </n-card>
          </n-gi>
        </n-grid>
      </n-tab-pane>

      <!-- 事件通知开关 -->
      <n-tab-pane name="events" tab="事件通知">
        <n-data-table
          :columns="eventColumns"
          :data="eventData"
          :pagination="false"
          :bordered="false"
          striped
          style="margin-top: 16px"
        />
      </n-tab-pane>

      <!-- 发送记录 -->
      <n-tab-pane name="logs" tab="发送记录">
        <n-space style="margin-top: 16px; margin-bottom: 12px">
          <n-select
            v-model:value="logFilter.channel"
            :options="logChannelOptions"
            placeholder="全部渠道"
            clearable
            style="width: 140px"
          />
          <n-select
            v-model:value="logFilter.status"
            :options="logStatusOptions"
            placeholder="全部状态"
            clearable
            style="width: 140px"
          />
          <n-button type="primary" @click="refreshLogs">刷新</n-button>
        </n-space>
        <n-data-table
          :columns="logColumns"
          :data="filteredLogs"
          :pagination="logPagination"
          :bordered="false"
          striped
        />
      </n-tab-pane>
    </n-tabs>
  </n-card>

  <!-- 测试发送对话框 -->
  <n-modal
    v-model:show="showTestModal"
    preset="card"
    title="测试发送"
    style="width: 480px"
    :bordered="false"
  >
    <n-form :model="testForm" label-placement="left" label-width="80">
      <n-form-item label="渠道">
        <n-tag :type="testForm.channel === 'email' ? 'info' : testForm.channel === 'sms' ? 'warning' : 'default'">
          {{ testForm.channel === 'email' ? '邮件' : testForm.channel === 'sms' ? '短信' : '站内信' }}
        </n-tag>
      </n-form-item>
      <n-form-item :label="testForm.channel === 'email' ? '邮箱' : '手机号'" v-if="testForm.channel !== 'site'">
        <n-input v-model:value="testForm.target" :placeholder="testForm.channel === 'email' ? 'test@example.com' : '13800138000'" />
      </n-form-item>
    </n-form>
    <template #footer>
      <n-space justify="end">
        <n-button @click="showTestModal = false">取消</n-button>
        <n-button type="primary" :loading="testSending" @click="sendTest">发送</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h } from 'vue'
import { useMessage, NTag, NSwitch, NButton, NSpace, NPopconfirm } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

const message = useMessage()
const activeTab = ref('channels')

// ---- 渠道配置 ----
const channels = reactive({
  email: { enabled: true, fromName: '锚点财务', replyTo: 'support@anchorfinance.com' },
  sms: { enabled: false, signName: '锚点财务', testPhone: '' },
  site: { enabled: true, retainDays: 90, maxCount: 200 },
})

function saveChannel(type: string) {
  message.success(`${type === 'email' ? '邮件' : type === 'sms' ? '短信' : '站内信'}渠道已保存`)
}

// ---- 测试发送 ----
const showTestModal = ref(false)
const testSending = ref(false)
const testForm = reactive({ channel: '', target: '' })

function testChannel(type: string) {
  testForm.channel = type
  testForm.target = ''
  showTestModal.value = true
}

function sendTest() {
  if (testForm.channel !== 'site' && !testForm.target) {
    message.warning('请输入目标地址')
    return
  }
  testSending.value = true
  setTimeout(() => {
    testSending.value = false
    showTestModal.value = false
    message.success('测试消息已发送')
  }, 800)
}

// ---- 事件通知 ----
interface EventRow {
  id: number
  name: string
  description: string
  email: boolean
  sms: boolean
  site: boolean
}

const eventData = ref<EventRow[]>([
  { id: 1, name: '用户注册', description: '新用户完成注册时触发', email: true, sms: false, site: true },
  { id: 2, name: '订单创建', description: '用户创建新订单时触发', email: true, sms: true, site: true },
  { id: 3, name: '订单支付成功', description: '订单支付成功后触发', email: true, sms: true, site: true },
  { id: 4, name: '账单到期提醒', description: '账单到期前 N 天触发', email: true, sms: false, site: true },
  { id: 5, name: '账单逾期', description: '账单逾期未支付时触发', email: true, sms: true, site: true },
  { id: 6, name: '工单创建', description: '用户提交新工单时触发', email: false, sms: false, site: true },
  { id: 7, name: '工单回复', description: '工单收到新回复时触发', email: true, sms: false, site: true },
  { id: 8, name: '密码重置', description: '用户请求重置密码时触发', email: true, sms: false, site: false },
  { id: 9, name: '服务开通', description: '服务开通成功后触发', email: true, sms: true, site: true },
  { id: 10, name: '服务到期', description: '服务即将到期时触发', email: true, sms: true, site: true },
  { id: 11, name: '退款通知', description: '退款处理完成时触发', email: true, sms: false, site: true },
  { id: 12, name: '系统公告', description: '管理员发布公告时触发', email: false, sms: false, site: true },
])

function toggleEvent(row: EventRow, field: 'email' | 'sms' | 'site', val: boolean) {
  row[field] = val
  message.success(`已${val ? '开启' : '关闭'}「${row.name}」的${field === 'email' ? '邮件' : field === 'sms' ? '短信' : '站内信'}通知`)
}

const eventColumns: DataTableColumns<EventRow> = [
  { title: '事件名称', key: 'name', width: 160 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  {
    title: '邮件通知',
    key: 'email',
    width: 110,
    align: 'center',
    render(row) {
      return h(NSwitch, { value: row.email, size: 'small', onUpdateValue: (v: boolean) => toggleEvent(row, 'email', v) })
    },
  },
  {
    title: '短信通知',
    key: 'sms',
    width: 110,
    align: 'center',
    render(row) {
      return h(NSwitch, { value: row.sms, size: 'small', onUpdateValue: (v: boolean) => toggleEvent(row, 'sms', v) })
    },
  },
  {
    title: '站内信通知',
    key: 'site',
    width: 110,
    align: 'center',
    render(row) {
      return h(NSwitch, { value: row.site, size: 'small', onUpdateValue: (v: boolean) => toggleEvent(row, 'site', v) })
    },
  },
]

// ---- 发送记录 ----
interface LogRow {
  id: number
  channel: string
  event: string
  target: string
  status: string
  time: string
}

const logFilter = reactive({ channel: null as string | null, status: null as string | null })

const logChannelOptions = [
  { label: '邮件', value: 'email' },
  { label: '短信', value: 'sms' },
  { label: '站内信', value: 'site' },
]

const logStatusOptions = [
  { label: '成功', value: 'success' },
  { label: '失败', value: 'failed' },
  { label: '待发送', value: 'pending' },
]

const logData = ref<LogRow[]>([
  { id: 1, channel: 'email', event: '订单支付成功', target: 'user@example.com', status: 'success', time: '2026-07-27 09:32' },
  { id: 2, channel: 'sms', event: '账单到期提醒', target: '138****8000', status: 'success', time: '2026-07-27 09:00' },
  { id: 3, channel: 'site', event: '工单回复', target: '用户 #1024', status: 'success', time: '2026-07-26 18:15' },
  { id: 4, channel: 'email', event: '密码重置', target: 'admin@example.com', status: 'failed', time: '2026-07-26 14:20' },
  { id: 5, channel: 'email', event: '用户注册', target: 'newuser@example.com', status: 'success', time: '2026-07-26 10:05' },
  { id: 6, channel: 'sms', event: '订单创建', target: '139****1234', status: 'pending', time: '2026-07-26 09:50' },
  { id: 7, channel: 'site', event: '系统公告', target: '全体用户', status: 'success', time: '2026-07-25 16:00' },
  { id: 8, channel: 'email', event: '服务到期', target: 'vip@example.com', status: 'success', time: '2026-07-25 11:30' },
])

const filteredLogs = computed(() => {
  return logData.value.filter((l) => {
    if (logFilter.channel && l.channel !== logFilter.channel) return false
    if (logFilter.status && l.status !== logFilter.status) return false
    return true
  })
})

const logPagination = reactive({ pageSize: 10 })

const channelLabel: Record<string, string> = { email: '邮件', sms: '短信', site: '站内信' }
const channelTagType: Record<string, 'info' | 'warning' | 'default'> = { email: 'info', sms: 'warning', site: 'default' }
const statusLabel: Record<string, string> = { success: '成功', failed: '失败', pending: '待发送' }
const statusTagType: Record<string, 'success' | 'error' | 'warning'> = { success: 'success', failed: 'error', pending: 'warning' }

const logColumns: DataTableColumns<LogRow> = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '渠道',
    key: 'channel',
    width: 90,
    render(row) {
      return h(NTag, { type: channelTagType[row.channel], size: 'small', bordered: false }, { default: () => channelLabel[row.channel] })
    },
  },
  { title: '事件', key: 'event', width: 160 },
  { title: '目标', key: 'target', ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(NTag, { type: statusTagType[row.status], size: 'small', bordered: false }, { default: () => statusLabel[row.status] })
    },
  },
  { title: '时间', key: 'time', width: 160 },
]

function refreshLogs() {
  message.success('记录已刷新')
}
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
</style>
