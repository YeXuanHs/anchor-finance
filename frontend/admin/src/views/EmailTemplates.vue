<template>
  <n-card :bordered="false" rounded>
    <template #header>
      <n-space justify="space-between" align="center">
        <span>邮件模板管理</span>
        <n-button type="primary" @click="openCreate">
          <template #icon><n-icon><AddIcon /></n-icon></template>
          新建模板
        </n-button>
      </n-space>
    </template>

    <!-- 模板列表 -->
    <n-data-table
      :columns="columns"
      :data="templateList"
      :bordered="false"
      :single-line="false"
      striped
    />

    <!-- 创建/编辑对话框 -->
    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="editingId ? '编辑模板' : '新建模板'"
      style="width: 680px"
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
          <n-select
            v-model:value="formData.type"
            :options="templateTypeOptions"
            placeholder="选择模板类型"
          />
        </n-form-item>
        <n-form-item label="邮件标题" path="subject">
          <n-input v-model:value="formData.subject" placeholder="如：您的订单 {{order_id}} 已确认" />
        </n-form-item>
        <n-form-item label="邮件内容" path="content">
          <n-input
            v-model:value="formData.content"
            type="textarea"
            :rows="10"
            placeholder="输入邮件 HTML 内容..."
          />
        </n-form-item>
        <n-form-item label="启用状态">
          <n-switch v-model:value="formData.enabled" />
        </n-form-item>
      </n-form>

      <!-- 变量说明 -->
      <n-divider>可用模板变量</n-divider>
      <n-space vertical size="small">
        <n-text depth="3">在标题和内容中使用以下变量，发送时会自动替换：</n-text>
        <n-space>
          <n-tag v-for="v in currentVars" :key="v.var" size="small" type="info">
            {{ v.var }} - {{ v.desc }}
          </n-tag>
        </n-space>
      </n-space>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" @click="handleSave">保存</n-button>
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
      <n-space vertical>
        <n-text strong>标题：</n-text>
        <n-text>{{ previewData.subject }}</n-text>
        <n-divider />
        <n-text strong>内容：</n-text>
        <div class="preview-content" v-html="previewData.content" />
      </n-space>
      <template #footer>
        <n-button @click="showPreview = false">关闭</n-button>
      </template>
    </n-modal>
  </n-card>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h } from 'vue'
import { useMessage, NButton, NIcon, NSpace, NTag, NSwitch, NPopconfirm } from 'naive-ui'
import { AddOutline as AddIcon } from '@vicons/ionicons5'
import type { DataTableColumns, FormRules } from 'naive-ui'

const message = useMessage()
const showModal = ref(false)
const showPreview = ref(false)
const editingId = ref<string | null>(null)

// ---- 模板类型选项 ----
const templateTypeOptions = [
  { label: '订单确认', value: 'order_confirm' },
  { label: '账单提醒', value: 'bill_remind' },
  { label: '密码重置', value: 'password_reset' },
  { label: '工单回复', value: 'ticket_reply' },
  { label: '账户激活', value: 'account_activate' },
  { label: '通用通知', value: 'general' },
]

// ---- 模板变量映射 ----
const variableMap: Record<string, { var: string; desc: string }[]> = {
  order_confirm: [
    { var: '{{username}}', desc: '用户名' },
    { var: '{{order_id}}', desc: '订单号' },
    { var: '{{amount}}', desc: '金额' },
    { var: '{{product}}', desc: '产品名' },
    { var: '{{expire_date}}', desc: '到期时间' },
  ],
  bill_remind: [
    { var: '{{username}}', desc: '用户名' },
    { var: '{{bill_id}}', desc: '账单号' },
    { var: '{{amount}}', desc: '应付金额' },
    { var: '{{due_date}}', desc: '到期日' },
  ],
  password_reset: [
    { var: '{{username}}', desc: '用户名' },
    { var: '{{reset_link}}', desc: '重置链接' },
    { var: '{{expire_minutes}}', desc: '有效期(分)' },
  ],
  ticket_reply: [
    { var: '{{username}}', desc: '用户名' },
    { var: '{{ticket_id}}', desc: '工单号' },
    { var: '{{reply_content}}', desc: '回复内容' },
  ],
  account_activate: [
    { var: '{{username}}', desc: '用户名' },
    { var: '{{activate_link}}', desc: '激活链接' },
  ],
  general: [
    { var: '{{username}}', desc: '用户名' },
    { var: '{{content}}', desc: '通知内容' },
  ],
}

const currentVars = computed(() => variableMap[formData.type] || variableMap.general)

// ---- 模板列表数据 ----
const templateList = ref([
  {
    id: '1',
    name: '订单确认邮件',
    type: 'order_confirm',
    subject: '您的订单 {{order_id}} 已确认',
    content: '<p>尊敬的 {{username}}，您的订单 {{order_id}} 已确认，金额 {{amount}} 元。</p>',
    enabled: true,
  },
  {
    id: '2',
    name: '账单到期提醒',
    type: 'bill_remind',
    subject: '账单 {{bill_id}} 即将到期',
    content: '<p>尊敬的 {{username}}，您的账单 {{bill_id}} 将于 {{due_date}} 到期，请及时支付 {{amount}} 元。</p>',
    enabled: true,
  },
  {
    id: '3',
    name: '密码重置邮件',
    type: 'password_reset',
    subject: '密码重置请求',
    content: '<p>尊敬的 {{username}}，请点击以下链接重置密码（{{expire_minutes}} 分钟内有效）：<a href="{{reset_link}}">重置密码</a></p>',
    enabled: true,
  },
  {
    id: '4',
    name: '工单回复通知',
    type: 'ticket_reply',
    subject: '工单 {{ticket_id}} 有新回复',
    content: '<p>尊敬的 {{username}}，您的工单 {{ticket_id}} 有新回复：{{reply_content}}</p>',
    enabled: false,
  },
])

// ---- 表单 ----
const formData = reactive({
  name: '',
  type: 'order_confirm' as string,
  subject: '',
  content: '',
  enabled: true,
})

const formRules: FormRules = {
  name: { required: true, message: '请输入模板名称', trigger: 'blur' },
  type: { required: true, message: '请选择模板类型', trigger: 'change' },
  subject: { required: true, message: '请输入邮件标题', trigger: 'blur' },
  content: { required: true, message: '请输入邮件内容', trigger: 'blur' },
}

// ---- 预览 ----
const previewData = reactive({ subject: '', content: '' })

// ---- 表格列定义 ----
const columns: DataTableColumns<any> = [
  { title: '模板名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 120,
    render(row) {
      const opt = templateTypeOptions.find(o => o.value === row.type)
      return h(NTag, { size: 'small', type: 'info' }, { default: () => opt?.label || row.type })
    },
  },
  { title: '标题', key: 'subject', ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'enabled',
    width: 100,
    align: 'center',
    render(row) {
      return h(NSwitch, {
        value: row.enabled,
        size: 'small',
        onUpdateValue: (val: boolean) => {
          row.enabled = val
          message.success(val ? '已启用' : '已禁用')
        },
      })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    align: 'center',
    render(row) {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => handlePreview(row) }, { default: () => '预览' }),
          h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => openEdit(row) }, { default: () => '编辑' }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
            trigger: () => h(NButton, { size: 'small', quaternary: true, type: 'error' }, { default: () => '删除' }),
            default: () => '确定删除该模板？',
          }),
        ],
      })
    },
  },
]

// ---- 操作方法 ----
function openCreate() {
  editingId.value = null
  Object.assign(formData, { name: '', type: 'order_confirm', subject: '', content: '', enabled: true })
  showModal.value = true
}

function openEdit(row: any) {
  editingId.value = row.id
  Object.assign(formData, { name: row.name, type: row.type, subject: row.subject, content: row.content, enabled: row.enabled })
  showModal.value = true
}

function handlePreview(row: any) {
  previewData.subject = row.subject
  previewData.content = row.content
  showPreview.value = true
}

function handleSave() {
  if (!formData.name || !formData.subject || !formData.content) {
    message.warning('请填写完整信息')
    return
  }
  if (editingId.value) {
    const idx = templateList.value.findIndex(t => t.id === editingId.value)
    if (idx >= 0) {
      Object.assign(templateList.value[idx], { ...formData })
    }
    message.success('模板已更新')
  } else {
    templateList.value.push({
      id: String(Date.now()),
      ...JSON.parse(JSON.stringify(formData)),
    })
    message.success('模板已创建')
  }
  showModal.value = false
}

function handleDelete(id: string) {
  templateList.value = templateList.value.filter(t => t.id !== id)
  message.success('模板已删除')
}
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
  min-height: 120px;
}
</style>
