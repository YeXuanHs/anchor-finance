<template>
  <div class="news-page">
    <art-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('content.newsList.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> {{ $t('content.newsList.addNews') }}
          </el-button>
        </div>
      </template>

      <p class="description">{{ $t('content.newsList.description') }}</p>

      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="60" align="center" />
        <el-table-column prop="title" :label="$t('content.newsList.titleColumn')" min-width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ row.title }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="category_name" :label="$t('content.newsList.category')" width="120" />
        <el-table-column prop="created_at" :label="$t('content.newsList.publishTime')" width="170" />
        <el-table-column :label="$t('content.newsList.hidden')" width="80" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.is_hidden" :active-value="1" :inactive-value="0" @change="handleToggleHidden(row)" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('content.newsList.manage')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination" v-if="total > 0">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="total" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" @size-change="fetchList" @current-change="fetchList" />
      </div>
    </art-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? $t('content.newsList.editNews') : $t('content.newsList.addNews')" width="700px" @close="resetForm">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item :label="$t('content.newsList.titleColumn')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('content.newsList.enterTitle')" />
        </el-form-item>
        <el-form-item :label="$t('content.newsList.category')" prop="category_id">
          <el-select v-model="formData.category_id" :placeholder="$t('content.newsList.selectCategory')" style="width: 100%">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('content.newsList.content')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="10" :placeholder="$t('content.newsList.enterContent')" />
        </el-form-item>
        <el-form-item :label="$t('common.sortOrder')" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(100)
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const categories = ref<{ id: number; name: string }[]>([])

const formData = reactive({ title: '', category_id: null as number | null, content: '', sort_order: 0 })
const rules: FormRules = {
  title: [{ required: true, message: () => $t('content.newsList.enterTitle'), trigger: 'blur' }],
  category_id: [{ required: true, message: () => $t('content.newsList.selectCategory'), trigger: 'change' }],
  content: [{ required: true, message: () => $t('content.newsList.enterContent'), trigger: 'blur' }]
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/news', params: { page: currentPage.value, page_size: pageSize.value } })
    tableData.value = data?.list || data || []
    total.value = data?.total || 0
  } catch (error) { console.error('获取新闻列表失败:', error) } finally { loading.value = false }
}

const fetchCategories = async () => {
  try {
    const data = await request.get({ url: '/api/admin/news-categories' })
    categories.value = data?.list || data || []
  } catch (error) { console.error('获取分类失败:', error) }
}

const resetForm = () => { formData.title = ''; formData.category_id = null; formData.content = ''; formData.sort_order = 0; editingId.value = null; formRef.value?.resetFields() }
const handleAdd = () => { isEdit.value = false; resetForm(); dialogVisible.value = true }
const handleEdit = (row: any) => {
  isEdit.value = true; editingId.value = row.id
  Object.assign(formData, { title: row.title, category_id: row.category_id, content: row.content || '', sort_order: row.sort_order || 0 })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate(); submitting.value = true
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/news/${editingId.value}`, data: formData }); ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/news', data: formData }); ElMessage.success($t('common.addSuccess'))
    }
    dialogVisible.value = false; fetchList()
  } catch (error) { console.error('提交失败:', error) } finally { submitting.value = false }
}

const handleToggleHidden = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/news/${row.id}`, data: { is_hidden: row.is_hidden } }); ElMessage.success($t('common.updateSuccess'))
  } catch (error) { console.error('更新失败:', error); fetchList() }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('content.newsList.confirmDelete'), $t('common.tips'), { type: 'warning' })
    await request.del({ url: `/api/admin/news/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList()
  } catch (error) { if (error !== 'cancel') console.error('删除失败:', error) }
}

onMounted(() => { fetchList(); fetchCategories() })
</script>

<style scoped lang="scss">
.news-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.description { color: #666; font-size: 14px; margin-bottom: 16px; }
.pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
</style>
