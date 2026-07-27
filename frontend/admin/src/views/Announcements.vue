<template>
  <n-card title="公告管理">
    <template #header-extra>
      <n-button type="primary" @click="openModal()">
        <template #icon><n-icon><AddIcon /></n-icon></template>
        发布公告
      </n-button>
    </template>

    <n-data-table :columns="columns" :data="announcements" :loading="loading" :bordered="false" />

    <!-- Announcement Modal -->
    <n-modal v-model:show="modalVisible" preset="card" :title="editing ? '编辑公告' : '发布公告'" style="width: 640px">
      <n-form ref="formRef" :model="formData" :rules="rules" label-placement="left" label-width="80">
        <n-form-item label="标题" path="title">
          <n-input v-model:value="formData.title" placeholder="公告标题" />
        </n-form-item>
        <n-form-item label="类型" path="type">
          <n-select v-model:value="formData.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="内容" path="content">
          <n-input v-model:value="formData.content" type="textarea" :rows="8" placeholder="公告内容，支持Markdown格式" />
        </n-form-item>
        <n-form-item label="置顶">
          <n-switch v-model:value="formData.pinned" />
        </n-form-item>
        <n-form-item label="状态">
          <n-radio-group v-model:value="formData.status">
            <n-radio value="published">立即发布</n-radio>
            <n-radio value="draft">存为草稿</n-radio>
          </n-radio-group>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="modalVisible = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </n-card>
</template>

<script setup lang="ts">
import { h, ref, reactive } from 'vue'
import { useMessage, NTag, NButton, NSpace, NPopconfirm, type DataTableColumns, type FormInst, type FormRules } from 'naive-ui'
import { AddOutline as AddIcon } from '@vicons/ionicons5'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const formRef = ref<FormInst | null>(null)
const editing = ref<any>(null)

const typeOptions = [
  { label: '通知', value: 'notice' },
  { label: '活动', value: 'promo' },
  { label: '维护', value: 'maintenance' },
  { label: '更新', value: 'update' },
]

const typeMap: Record<string, { label: string; type: 'success' | 'info' | 'warning' | 'error' }> = {
  notice: { label: '通知', type: 'info' },
  promo: { label: '活动', type: 'success' },
  maintenance: { label: '维护', type: 'warning' },
  update: { label: '更新', type: 'info' },
}

const formData = reactive({
  title: '',
  type: 'notice',
  content: '',
  pinned: false,
  status: 'published',
})

const rules: FormRules = {
  title: { required: true, message: '请输入公告标题', trigger: 'blur' },
  content: { required: true, message: '请输入公告内容', trigger: 'blur' },
}

const announcements = ref([
  { id: 1, title: '系统升级通知 - v2.0版本发布', type: 'update', status: 'published', pinned: true, views: 1250, createdAt: '2024-03-15 10:00', author: '管理员' },
  { id: 2, title: '春季促销活动 - 全场8折', type: 'promo', status: 'published', pinned: false, views: 890, createdAt: '2024-03-10 09:00', author: '管理员' },
  { id: 3, title: '3月18日凌晨服务器维护通知', type: 'maintenance', status: 'published', pinned: true, views: 2100, createdAt: '2024-03-08 15:00', author: '管理员' },
  { id: 4, title: '新用户注册优惠活动', type: 'promo', status: 'draft', pinned: false, views: 0, createdAt: '2024-03-07 11:00', author: '管理员' },
  { id: 5, title: '关于加强账户安全的通知', type: 'notice', status: 'published', pinned: false, views: 560, createdAt: '2024-03-01 08:30', author: '管理员' },
])

const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '标题',
    key: 'title',
    render: (row) =>
      h('span', [
        row.pinned ? h(NTag, { type: 'error', size: 'tiny', style: 'margin-right: 6px' }, { default: () => '置顶' }) : null,
        row.title,
      ]),
  },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row) => {
      const t = typeMap[row.type]
      return h(NTag, { type: t.type, size: 'small' }, { default: () => t.label })
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NTag, { type: row.status === 'published' ? 'success' : 'default', size: 'small', round: true }, { default: () => (row.status === 'published' ? '已发布' : '草稿') }),
  },
  { title: '浏览量', key: 'views', width: 80 },
  { title: '发布时间', key: 'createdAt', width: 160 },
  { title: '作者', key: 'author', width: 80 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => '编辑' }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
            default: () => '确认删除该公告？',
          }),
        ],
      }),
  },
]

function openModal(announcement?: any) {
  editing.value = announcement || null
  if (announcement) {
    Object.assign(formData, {
      title: announcement.title,
      type: announcement.type,
      content: '',
      pinned: announcement.pinned,
      status: announcement.status,
    })
  } else {
    Object.assign(formData, { title: '', type: 'notice', content: '', pinned: false, status: 'published' })
  }
  modalVisible.value = true
}

async function handleSubmit() {
  try { await formRef.value?.validate() } catch { return }
  submitting.value = true
  try {
    if (editing.value) {
      Object.assign(editing.value, { title: formData.title, type: formData.type, pinned: formData.pinned, status: formData.status })
      message.success('公告更新成功')
    } else {
      announcements.value.unshift({
        id: Date.now(),
        title: formData.title,
        type: formData.type,
        status: formData.status,
        pinned: formData.pinned,
        views: 0,
        createdAt: new Date().toLocaleString(),
        author: '管理员',
      })
      message.success('公告发布成功')
    }
    modalVisible.value = false
  } catch (err: any) {
    message.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleDelete(id: number) {
  announcements.value = announcements.value.filter((a) => a.id !== id)
  message.success('公告已删除')
}
</script>
