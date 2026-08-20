<template>
  <div class="news-category-page">
    <art-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('content.newsCategory.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> {{ $t('content.newsCategory.addCategory') }}
          </el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="60" align="center" />
        <el-table-column prop="name" :label="$t('content.newsCategory.name')" min-width="150" />
        <el-table-column prop="sort_order" :label="$t('common.sortOrder')" width="80" align="center" />
        <el-table-column prop="news_count" :label="$t('content.newsCategory.newsCount')" width="80" align="center" />
        <el-table-column :label="$t('common.action')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? $t('content.newsCategory.editCategory') : $t('content.newsCategory.addCategory')" width="450px" @close="resetForm">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item :label="$t('common.name')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('content.newsCategory.enterName')" />
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
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const formData = reactive({ name: '', sort_order: 0 })
const rules: FormRules = { name: [{ required: true, message: () => $t('content.newsCategory.enterName'), trigger: 'blur' }] }

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/news/categories' })
    tableData.value = data?.list || data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  } finally {
    loading.value = false
  }
}

const resetForm = () => { formData.name = ''; formData.sort_order = 0; editingId.value = null; formRef.value?.resetFields() }
const handleAdd = () => { isEdit.value = false; resetForm(); dialogVisible.value = true }
const handleEdit = (row: any) => { isEdit.value = true; editingId.value = row.id; Object.assign(formData, { name: row.name, sort_order: row.sort_order || 0 }); dialogVisible.value = true }

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate(); submitting.value = true
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/news/categories/${editingId.value}`, data: formData }); ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/news/categories', data: formData }); ElMessage.success($t('common.addSuccess'))
    }
    dialogVisible.value = false; fetchList()
  } catch (error) { console.error('提交失败:', error) } finally { submitting.value = false }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('content.newsCategory.confirmDelete'), $t('common.tips'), { type: 'warning' })
    await request.del({ url: `/api/admin/news/categories/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList()
  } catch (error) { if (error !== 'cancel') console.error('删除失败:', error) }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.news-category-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
