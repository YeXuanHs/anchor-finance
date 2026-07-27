<template>
  <n-card :bordered="false" rounded>
    <template #header>
      <n-space justify="space-between" align="center">
        <span style="font-size: 18px; font-weight: 600">邮件模板管理</span>
        <n-button type="primary" @click="openCreate">
          <template #icon><n-icon><AddOutline /></n-icon></template>
          新建模板
        </n-button>
      </n-space>
    </template>

    <n-data-table
      :columns="columns"
      :data="tableData"
      :pagination="pagination"
      :bordered="false"
      striped
    />
  </n-card>

  <!-- 创建/编辑对话框 -->
  <n-modal
    v-model:show="showModal"
    preset="card"
    :title="isEditing ? '编辑模板' : '新建模板'"
    style="width: 720px"
    :bordered="false"
    :segmented="{ content: true, footer: true }"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="100"
    >
      <n-form-item label="模板名称" path="name">
        <n-input v-model:value="formData.name" placeholder="如：订单确认邮件" />
      </n-form-item>
      <n-form-item label="模板类型" path="type">
        <n-select v-model:value="formData.type" :options="typeOptions" placeholder="选择模板类型" />
      </n-form-item>
      <n-form-item label="邮件标题" path="subject">
        <n-input v-model:value="formData.subject" placeholder="如：您的订单 {{order_id}} 已确认" />
      </n-form-item>
      <n-form-item label="邮件内容" path="content">
        <n-input
          v-model:value="formData.content"
          type="textarea"
          :rows="10"
          placeholder="输入邮件 HTML 或纯文本内容"
        />
      </n-form-item>
      <n-form-item label="启用状态">
        <n-switch v-model:value="formData.enabled" />
      </n-form-item>
    </n-form>

    <n-collapse style="margin-top: 12px">
      <n-collapse-item title="可用模板变量" name="vars">
        <n-space vertical>
          <n-tag v-for="v in currentVars" :key="v.key" size="small" type="info">
            {{ v.key }}
          </n-tag>
          <n-text depth="3" style="font-size: 12px; margin-top: 4px">
            使用双花括号包裹变量，例如：{{ '{{username}}' }}
          </n-text>
        </n-space>
      </n-collapse-item>
    </n-collapse>

    <template #footer>
      <n-space justify="end">
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">保存</n-button>
      </n-space>
    </template>
  </n-modal>

  <!-- 预览对话框 -->
  <n-modal
    v-model:show="showPreview"
    preset="card"
    title="模板预览"
    style="width: 640px"
    :bordered="false"
  >
    <n-descriptions bordered :column="1" label-placement="left">
      <n-descriptions-item label="模板名称">{{ previewData.name }}</n-descriptions-item>
      <n-descriptions-item label="类型">
        <n-tag :type="getTypeTagType(previewData.type)" size="small">{{ getTypeLabel(previewData.type) }}</n-tag>
      </n-descriptions-item>
      <n-descriptions-item label="标题">{{ previewData.subject }}</n-descriptions-item>
      <n-descriptions-item label="状态">
        <n-tag :type="previewData.enabled ? 'success' : 'default'" size="small">
          {{ previewData.enabled ? '已启用' : '已禁用' }}
        </n-tag>
      </n-descriptions-item>
    </n-descriptions>
    <n-divider />
    <div class="preview-content" v-html="previewData.content" />
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h } from 'vue'
import { useMessage, NButton, NIcon, NSpace, NTag, NPopconfirm, NSwitch } from 'naive-ui'
import { AddOutline, CreateOutline, EyeOutline, TrashOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormRules } from 'naive-ui'

const message = useMessage()

// ---- 类型选项 ----
const typeOptions = [
  { label: '订单确认', value: 'order_confirm' },
  { label: '账单提醒', value: 'bill_reminder' },
  { label: '密码重置', value: 'password_reset' },
  { label: '工单回复', value: 'ticket_reply' },
  { label: '注册欢迎', value: 'welcome' },
  { label: '其他', value: 'other' },
]

const typeMap: Record<string, string> = {
  order_confirm: '订单确认',
  bill_reminder: '账单提醒',
  password_reset: '密码重置',
  ticket_reply: '工单回复',
  welcome: '注册欢迎',
  other: '其他',
}

const typeTagMap: Record<string, 'success' | 'warning' | 'error' | 'info' | 'default'> = {
  order_confirm: 'success',
  bill_reminder: 'warning',
  password_reset: 'error',
  ticket_reply: 'info',
  welcome: 'success',
  other: 'default',
}

// ---- 模板变量 ----
const varsMap: Record<string, { key: string }[]> = {
  order_confirm: [
    { key: '{{username}}' },
    { key: '{{order_id}}' },
    { key: '{{order_amount}}' },
    { key: '{{order_time}}' },
    { key: '{{order_items}}' },
  ],
  bill_reminder: [
    { key: '{{username}}' },
    { key: '{{bill_id}}' },
    { key: '{{bill_amount}}' },
    { key: '{{due_date}}' },
    { key: '{{service_name}}' },
  ],
  password_reset: [
    { key: '{{username}}' },
    { key: '{{reset_link}}' },
    { key: '{{expire_minutes}}' },
  ],
  ticket_reply: [
    { key: '{{username}}' },
    { key: '{{ticket_id}}' },
    { key: '{{reply_content}}' },
    { key: '{{reply_time}}' },
  ],
  welcome: [
    { key: '{{username}}' },
    { key: '{{site_name}}' },
    { key: '{{login_url}}' },
  ],
  other: [
    { key: '{{username}}' },
    { key: '{{site_name}}' },
  ],
}

// ---- 模拟数据 ----
interface TemplateRow {
  id: number
  name: string
  type: string
  subject: string
  content: string
  enabled: boolean
  updatedAt: string
}

const tableData = ref<TemplateRow[]>([
  {
    id: 1,
    name: '订单确认邮件',
    type: 'order_confirm',
    subject: '您的订单 {{order_id}} 已确认',
    content: '<h2>尊敬的 {{username}}，</h2><p>您的订单 <strong>{{order_id}}</strong> 已确认，金额 {{order_amount}} 元。</p><p>感谢您的支持！</p>',
    enabled: true,
    updatedAt: '2026-07-20 10:30',
  },
  {
    id: 2,
    name: '账单提醒',
    type: 'bill_reminder',
    subject: '账单 {{bill_id}} 即将到期',
    content: '<h2>{{username}} 您好，</h2><p>您的账单 <strong>{{bill_id}}</strong>（{{service_name}}）将于 {{due_date}} 到期，金额 {{bill_amount}} 元，请及时处理。</p>',
    enabled: true,
    updatedAt: '2026-07-18 14:20',
  },
  {
    id: 3,
    name: '密码重置',
    type: 'password_reset',
    subject: '重置您的密码',
    content: '<h2>{{username}} 您好，</h2><p>请点击以下链接重置密码（{{expire_minutes}} 分钟内有效）：</p><p><a href="{{reset_link}}">重置密码</a></p><p>如非本人操作，请忽略此邮件。</p>',
    enabled: true,
    updatedAt: '2026-07-15 09:00',
  },
  {
    id: 4,
    name: '工单回复通知',
    type: 'ticket_reply',
    subject: '工单 #{{ticket_id}} 有新回复',
    content: '<h2>{{username}} 您好，</h2><p>您的工单 <strong>#{{ticket_id}}</strong> 收到新回复：</p><blockquote>{{reply_content}}</blockquote><p>回复时间：{{reply_time}}</p>',
    enabled: false,
    updatedAt: '2026-07-10 16:45',
  },
  {
    id: 5,
    name: '注册欢迎邮件',
    type: 'welcome',
    subject: '欢迎加入 {{site_name}}',
    content: '<h2>Hi {{username}}，</h2><p>欢迎注册 {{site_name}}！</p><p><a href="{{login_url}}">点击此处登录</a></p>',
    enabled: true,
    updatedAt: '2026-07-05 11:00',
  },
])

const pagination = reactive({ pageSize: 10 })

// ---- 表格列 ----
const columns: DataTableColumns<TemplateRow> = [
  { title: '模板名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 120,
    render(row) {
      return h(NTag, { type: typeTagMap[row.type] || 'default', size: 'small', bordered: false }, { default: () => typeMap[row.type] || row.type })
    },
  },
  { title: '邮件标题', key: 'subject', ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'enabled',
    width: 100,
    align: 'center',
    render(row) {
      return h(NSwitch, {
        value: row.enabled,
        size: 'small',
        onUpdateValue: (val: boolean) => { row.enabled = val; message.success(val ? '已启用' : '已禁用') },
      })
    },
  },
  { title: '更新时间', key: 'updatedAt', width: 160 },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    align: 'center',
    render(row) {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => openPreview(row) }, {
            icon: () => h(NIcon, null, { default: () => h(EyeOutline) }),
            default: () => '预览',
          }),
          h(NButton, { size: 'small', quaternary: true, type: 'warning', onClick: () => openEdit(row) }, {
            icon: () => h(NIcon, null, { default: () => h(CreateOutline) }),
            default: () => '编辑',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
            trigger: () => h(NButton, { size: 'small', quaternary: true, type: 'error' }, {
              icon: () => h(NIcon, null, { default: () => h(TrashOutline) }),
              default: () => '删除',
            }),
            default: () => '确认删除该模板？',
          }),
        ],
      })
    },
  },
]

// ---- 表单 ----
const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const formRef = ref()

const defaultForm = () => ({
  name: '',
  type: null as string | null,
  subject: '',
  content: '',
  enabled: true,
})

const formData = reactive(defaultForm())

const formRules: FormRules = {
  name: { required: true, message: '请输入模板名称', trigger: 'blur' },
  type: { required: true, message: '请选择模板类型', trigger: 'change' },
  subject: { required: true, message: '请输入邮件标题', trigger: 'blur' },
  content: { required: true, message: '请输入邮件内容', trigger: 'blur' },
}

const currentVars = computed(() => (formData.type && varsMap[formData.type]) || varsMap.other)

function openCreate() {
  isEditing.value = false
  editingId.value = null
  Object.assign(formData, defaultForm())
  showModal.value = true
}

function openEdit(row: TemplateRow) {
  isEditing.value = true
  editingId.value = row.id
  Object.assign(formData, { name: row.name, type: row.type, subject: row.subject, content: row.content, enabled: row.enabled })
  showModal.value = true
}

function handleSave() {
  formRef.value?.validate((errors: any) => {
    if (errors) return
    saving.value = true
    setTimeout(() => {
      if (isEditing.value && editingId.value !== null) {
        const idx = tableData.value.findIndex((t) => t.id === editingId.value)
        if (idx !== -1) {
          Object.assign(tableData.value[idx], { ...formData, updatedAt: new Date().toLocaleString() })
        }
        message.success('模板已更新')
      } else {
        const newId = Math.max(...tableData.value.map((t) => t.id), 0) + 1
        tableData.value.push({ id: newId, ...formData, updatedAt: new Date().toLocaleString() } as TemplateRow)
        message.success('模板已创建')
      }
      saving.value = false
      showModal.value = false
    }, 500)
  })
}

function handleDelete(id: number) {
  tableData.value = tableData.value.filter((t) => t.id !== id)
  message.success('模板已删除')
}

// ---- 预览 ----
const showPreview = ref(false)
const previewData = reactive({ name: '', type: '', subject: '', content: '', enabled: true })

function openPreview(row: TemplateRow) {
  Object.assign(previewData, row)
  showPreview.value = true
}

function getTypeLabel(type: string) { return typeMap[type] || type }
function getTypeTagType(type: string) { return typeTagMap[type] || 'default' }
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
.preview-content {
  padding: 16px;
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
  background: var(--n-color);
  line-height: 1.8;
}
.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3) {
  margin-top: 0;
}
</style>
