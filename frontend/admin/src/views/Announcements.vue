<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">公告管理</span>
          <el-button type="primary" @click="openModal()">
            <el-icon><Plus /></el-icon>发布公告
          </el-button>
        </div>
      </template>

      <el-table :data="announcements" v-loading="loading" stripe size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column label="标题">
          <template #default="{ row }">
            <span>
              <el-tag v-if="row.pinned" type="danger" size="small" style="margin-right: 6px">置顶</el-tag>
              {{ row.title }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="typeMap[row.type]?.type as any" size="small">{{ typeMap[row.type]?.label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small" round>
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="浏览量" width="80" />
        <el-table-column prop="createdAt" label="发布时间" width="160" />
        <el-table-column prop="author" label="作者" width="80" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openModal(row)">编辑</el-button>
            <el-popconfirm title="确认删除该公告？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button text type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Announcement Dialog -->
    <el-dialog v-model="modalVisible" :title="editing ? '编辑公告' : '发布公告'" width="640px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="公告标题" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" style="width: 100%">
            <el-option v-for="o in typeOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="8" placeholder="公告内容，支持Markdown格式" />
        </el-form-item>
        <el-form-item label="置顶">
          <el-switch v-model="formData.pinned" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio value="published">立即发布</el-radio>
            <el-radio value="draft">存为草稿</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const formRef = ref<FormInstance>()
const editing = ref<any>(null)

const typeOptions = [
  { label: '通知', value: 'notice' }, { label: '活动', value: 'promo' },
  { label: '维护', value: 'maintenance' }, { label: '更新', value: 'update' },
]
const typeMap: Record<string, { label: string; type: string }> = {
  notice: { label: '通知', type: 'primary' }, promo: { label: '活动', type: 'success' },
  maintenance: { label: '维护', type: 'warning' }, update: { label: '更新', type: 'info' },
}

const formData = reactive({ title: '', type: 'notice', content: '', pinned: false, status: 'published' })
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

function openModal(announcement?: any) {
  editing.value = announcement || null
  if (announcement) {
    Object.assign(formData, { title: announcement.title, type: announcement.type, content: '', pinned: announcement.pinned, status: announcement.status })
  } else {
    Object.assign(formData, { title: '', type: 'notice', content: '', pinned: false, status: 'published' })
  }
  modalVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return
  try { await formRef.value.validate() } catch { return }
  submitting.value = true
  try {
    if (editing.value) {
      Object.assign(editing.value, { title: formData.title, type: formData.type, pinned: formData.pinned, status: formData.status })
      ElMessage.success('公告更新成功')
    } else {
      announcements.value.unshift({ id: Date.now(), title: formData.title, type: formData.type, status: formData.status, pinned: formData.pinned, views: 0, createdAt: new Date().toLocaleString(), author: '管理员' })
      ElMessage.success('公告发布成功')
    }
    modalVisible.value = false
  } finally { submitting.value = false }
}

function handleDelete(id: number) {
  announcements.value = announcements.value.filter((a) => a.id !== id)
  ElMessage.success('公告已删除')
}
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 16px; font-weight: 600; }
</style>
