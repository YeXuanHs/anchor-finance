<template>
  <div class="banners-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>轮播图管理</h3>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增轮播图
        </el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="140" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="90">
          <template #default="{ row }">
            <el-tag size="small">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="image_url" label="图片" width="100">
          <template #default="{ row }">
            <el-image
              v-if="row.image_url"
              :src="row.image_url"
              :preview-src-list="[row.image_url]"
              fit="cover"
              style="width: 60px; height: 36px; border-radius: 4px; cursor: pointer;"
            />
            <span v-else class="text-secondary">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="link" label="链接" min-width="150" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑轮播图' : '新增轮播图'"
      width="650px"
      destroy-on-close
    >
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入轮播图标题" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="2" placeholder="轮播图描述" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
          </el-select>
        </el-form-item>
        <el-form-item label="媒体URL" prop="media_url">
          <el-input v-model="formData.media_url" placeholder="图片或视频的URL地址" />
        </el-form-item>
        <el-form-item label="图片预览" v-if="formData.media_url && formData.type === 'image'">
          <el-image
            :src="formData.media_url"
            fit="contain"
            style="max-width: 300px; max-height: 150px; border: 1px solid var(--el-border-color-lighter); border-radius: 4px;"
          />
        </el-form-item>
        <el-form-item label="链接">
          <el-input v-model="formData.link" placeholder="点击跳转的链接地址" />
        </el-form-item>
        <el-form-item label="按钮文字">
          <el-input v-model="formData.button_text" placeholder="如: 立即查看" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="定时发布">
          <el-date-picker
            v-model="formData.schedule"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const list = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

const formData = reactive({
  title: '',
  description: '',
  type: 'image',
  media_url: '',
  link: '',
  button_text: '',
  sort: 0,
  status: 'active',
  schedule: null as string[] | null
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  media_url: [{ required: true, message: '请输入媒体URL', trigger: 'blur' }],
  sort: [{ required: true, message: '请输入排序', trigger: 'blur' }]
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = { image: '图片', video: '视频' }
  return map[type] || type
}

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/content/banners', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    list.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
    title: '',
    description: '',
    type: 'image',
    media_url: '',
    link: '',
    button_text: '',
    sort: 0,
    status: 'active',
    schedule: null
  })
  editingId.value = null
}

const handleAdd = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, {
    title: row.title,
    description: row.description || '',
    type: row.type,
    media_url: row.media_url || row.image_url || '',
    link: row.link || '',
    button_text: row.button_text || '',
    sort: row.sort || 0,
    status: row.status,
    schedule: row.start_time && row.end_time ? [row.start_time, row.end_time] : null
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    const payload: any = { ...formData }
    if (payload.schedule && payload.schedule.length === 2) {
      payload.start_time = payload.schedule[0]
      payload.end_time = payload.schedule[1]
    }
    delete payload.schedule

    if (isEdit.value && editingId.value) {
      await request.put(`/api/admin/content/banners/${editingId.value}`, payload)
      ElMessage.success('更新成功')
    } else {
      await request.post('/api/admin/content/banners', payload)
      ElMessage.success('新增成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch {
    // handled by interceptor
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定要删除轮播图「${row.title}」吗？`, '警告', { type: 'error' })
  try {
    await request.delete(`/api/admin/content/banners/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* handled */ }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.banners-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
    }
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  .text-secondary {
    color: var(--el-text-color-secondary);
  }
}
</style>
