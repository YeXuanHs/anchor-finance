<template>
  <n-card :bordered="false" rounded>
    <n-tabs v-model:value="activeTab" type="line" animated>
      <!-- 通知渠道配置 -->
      <n-tab-pane name="channels" tab="通知渠道">
        <n-grid :cols="3" :x-gap="16" style="margin-top: 20px">
          <!-- 邮件渠道 -->
          <n-gi>
            <n-card title="邮件通知" size="small" :bordered="false">
              <n-space vertical>
                <n-space justify="space-between" align="center">
                  <n-text>启用邮件通知</n-text>
                  <n-switch v-model:value="channels.email.enabled" />
                </n-space>
                <template v-if="channels.email.enabled">
                  <n-form label-placement="top" size="small">
                    <n-form-item label="发件人名称">
                      <n-input v-model:value="channels.email.fromName" placeholder="锚点财务" />
                    </n-form-item>
                    <n-form-item label="发件人邮箱">
                      <n-input v-model:value="channels.email.fromEmail" placeholder="noreply@example.com" />
                    </n-form-item>
                  </n-form>
                  <n-button size="small" type="primary" @click="testChannel('email')">
                    发送测试邮件
                  </n-button>
                </template>
              </n-space>
            </n-card>
          </n-gi>

          <!-- 短信渠道 -->
          <n-gi>
            <n-card title="短信通知" size="small" :bordered="false">
              <n-space vertical>
                <n-space justify="space-between" align="center">
                  <n-text>启用短信通知</n-text>
                  <n-switch v-model:value="channels.sms.enabled" />
                </n-space>
                <template v-if="channels.sms.enabled">
                  <n-form label-placement="top" size="small">
                    <n-form-item label="短信签名">
                      <n-input v-model:value="channels.sms.signName" placeholder="短信签名" />
                    </n-form-item>
                    <n-form-item label="测试手机号">
                      <n-input v-model:value="channels.sms.testPhone" placeholder="13800138000" />
                    </n-form-item>
                  </n-form>
                  <n-button size="small" type="primary" @click="testChannel('sms')">
                    发送测试短信
                  </n-button>
                </template>
              </n-space>
            </n-card>
          </n-gi>

          <!-- 站内信渠道 -->
          <n-gi>
            <n-card title="站内信通知" size="small" :bordered="false">
              <n-space vertical>
                <n-space justify="space-between" align="center">
                  <n-text>启用站内信</n-text>
                  <n-switch v-model:value="channels.internal.enabled" />
                </n-space>
                <template v-if="channels.internal.enabled">
                  <n-form label-placement="top" size="small">
                    <n-form-item label="保留天数">
                      <n-input-number v-model:value="channels.internal.retainDays" :min="1" :max="365" style="width:100%" />
                    </n-form-item>
                    <n-form-item label="最大条数">
                      <n-input-number v-model:value="channels.internal.maxCount" :min="10" :max="1000" style="width:100%" />
                    </n-form-item>
                  </n-form>
                </template>
              </n-space>
            </n-card>
          </n-gi>
        </n-grid>
        <n-button type="primary" style="margin-top: 20px" @click="saveChannels">保存渠道配置</n-button>
      </n-tab-pane>

      <!-- 通知事件 -->
      <n-tab-pane name="events" tab="通知事件">
        <n-data-table
          :columns="eventColumns"
          :data="eventList"
          :bordered="false"
          :single-line="false"
          striped
          style="margin-top: 12px"
        />
      </n-tab-pane>

      <!-- 发送记录 -->
      <n-tab-pane name="logs" tab="发送记录">
        <n-space style="margin-top: 12px; margin-bottom: 12px">
          <n-select
            v-model:value="logFilter.channel"
            :options="channelFilterOptions"
            placeholder="渠道筛选"
            clearable
            style="width: 140px"
          />
          <n-select
            v-model:value="logFilter.status"
            :options="statusFilterOptions"
            placeholder="状态筛选"
            clearable
            style="width: 140px"
          />
          <n-button @click="refreshLogs">刷新</n-button>
        </n-space>
        <n-data-table
          :columns="logColumns"
          :data="filteredLogs"
          :bordered="false"
          :single-line="false"
          striped
        />
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h } from 'vue'
import { useMessage, NTag, NSwitch, NSpace, NButton } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

const message = useMessage()
const activeTab = ref('channels')

// ---- 通知渠道配置 ----
const channels = reactive({
  email: { enabled: true, fromName: '锚点财务', fromEmail: 'noreply@anchorfinance.com' },
  sms: { enabled: false, signName: '锚点财务', testPhone: '' },
  internal: { enabled: true, retainDays: 90, maxCount: 500 },
})

function saveChannels() {
  // TODO: API call
  message.success('渠道配置已保存')
}

function testChannel(type: string) {
  // TODO: API call
  const names: Record<string, string> = { email: '测试邮件', sms: '测试短信', internal: '站内信' }
  message.info(`${names[type] || '测试'}发送中...`)
}

// ---- 通知事件列表 ----
const eventList = ref([
  { id: '1', name: 'user_register', label: '用户注册', desc: '新用户完成注册时触发', email: true, sms: false, internal: true },
  { id: '2', name: 'order_created', label: '订单创建', desc: '用户下单成功时触发', email: true, sms: true, internal: true },
  { id: '3', name: 'order_paid', label: '订单支付', desc: '订单支付成功时触发', email: true, sms: true, internal: true },
  { id: '4', name: 'order_expired', label: '订单到期', desc: '服务即将到期时触发', email: true, sms: false, internal: true },
  { id: '5', name: 'bill_created', label: '账单生成', desc: '新账单生成时触发', email: true, sms: false, internal: true },
  { id: '6', name: 'bill_overdue', label: '账单逾期', desc: '账单逾期未付时触发', email: true, sms: true, internal: true },
  { id: '7', name: 'ticket_reply', label: '工单回复', desc: '工单有新回复时触发', email: true, sms: false, internal: true },
  { id: '8', name: 'ticket_closed', label: '工单关闭', desc: '工单被关闭时触发', email: true, sms: false, internal: true },
  { id: '9', name: 'password_reset', label: '密码重置', desc: '用户请求重置密码时触发', email: true, sms: false, internal: false },
  { id: '10', name: 'server_down', label: '服务器宕机', desc: '监控到服务器不可用时触发', email: true, sms: true, internal: true },
])

const eventColumns: DataTableColumns<any> = [
  { title: '事件名称', key: 'label', width: 140 },
  { title: '描述', key: 'desc', ellipsis: { tooltip: true } },
  {
    title: '邮件通知',
    key: 'email',
    width: 100,
    align: 'center',
    render(row) {
      return h(NSwitch, {
        value: row.email,
        size: 'small',
        onUpdateValue: (val: boolean) => {
          row.email = val
          message.success(`邮件通知已${val ? '开启' : '关闭'}`)
        },
      })
    },
  },
  {
    title: '短信通知',
    key: 'sms',
    width: 100,
    align: 'center',
    render(row) {
      return h(NSwitch, {
        value: row.sms,
        size: 'small',
        onUpdateValue: (val: boolean) => {
          row.sms = val
          message.success(`短信通知已${val ? '开启' : '关闭'}`)
        },
      })
    },
  },
  {
    title: '站内信',
    key: 'internal',
    width: 100,
    align: 'center',
    render(row) {
      return h(NSwitch, {
        value: row.internal,
        size: 'small',
        onUpdateValue: (val: boolean) => {
          row.internal = val
          message.success(`站内信通知已${val ? '开启' : '关闭'}`)
        },
      })
    },
  },
]

// ---- 发送记录 ----
const logFilter = reactive({ channel: null as string | null, status: null as string | null })

const channelFilterOptions = [
  { label: '邮件', value: 'email' },
  { label: '短信', value: 'sms' },
  { label: '站内信', value: 'internal' },
]

const statusFilterOptions = [
  { label: '成功', value: 'success' },
  { label: '失败', value: 'failed' },
  { label: '待发送', value: 'pending' },
]

const logList = ref([
  { id: '1', time: '2026-07-27 09:30:12', channel: 'email', event: '订单支付', recipient: 'user@example.com', status: 'success' },
  { id: '2', time: '2026-07-27 09:28:05', channel: 'sms', event: '订单创建', recipient: '138****8000', status: 'success' },
  { id: '3', time: '2026-07-27 09:25:00', channel: 'internal', event: '工单回复', recipient: '用户 张三', status: 'success' },
  { id: '4', time: '2026-07-27 09:20:33', channel: 'email', event: '账单逾期', recipient: 'billing@example.com', status: 'failed' },
  { id: '5', time: '2026-07-27 08:15:22', channel: 'email', event: '用户注册', recipient: 'newuser@example.com', status: 'success' },
  { id: '6', time: '2026-07-26 18:00:00', channel: 'sms', event: '服务器宕机', recipient: '139****0000', status: 'success' },
])

const filteredLogs = computed(() => {
  return logList.value.filter(l => {
    if (logFilter.channel && l.channel !== logFilter.channel) return false
    if (logFilter.status && l.status !== logFilter.status) return false
    return true
  })
})

function refreshLogs() {
  // TODO: API call
  message.success('发送记录已刷新')
}

const channelLabel: Record<string, string> = { email: '邮件', sms: '短信', internal: '站内信' }

const logColumns: DataTableColumns<any> = [
  { title: '时间', key: 'time', width: 170 },
  {
    title: '渠道',
    key: 'channel',
    width: 80,
    render(row) {
      const typeMap: Record<string, string> = { email: 'info', sms: 'warning', internal: 'success' }
      return h(NTag, { size: 'small', type: (typeMap[row.channel] as any) || 'default' }, { default: () => channelLabel[row.channel] || row.channel })
    },
  },
  { title: '事件', key: 'event', width: 120 },
  { title: '接收方', key: 'recipient', ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'status',
    width: 80,
    align: 'center',
    render(row) {
      const map: Record<string, { label: string; type: string }> = {
        success: { label: '成功', type: 'success' },
        failed: { label: '失败', type: 'error' },
        pending: { label: '待发送', type: 'warning' },
      }
      const info = map[row.status] || { label: row.status, type: 'default' }
      return h(NTag, { size: 'small', type: info.type as any }, { default: () => info.label })
    },
  },
]
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
</style>
