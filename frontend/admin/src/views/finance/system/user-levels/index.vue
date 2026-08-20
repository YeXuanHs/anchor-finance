<template>
  <div class="user-levels-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.userLevels.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.userLevels.addLevel') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.userLevels.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.userLevels.levelName')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.userLevels.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.userLevels.all')" clearable>
            <el-option :label="$t('finance.userLevels.enabled')" :value="1" />
            <el-option :label="$t('finance.userLevels.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.userLevels.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.userLevels.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('finance.userLevels.levelName')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="level" :label="$t('finance.userLevels.levelValue')" width="100" align="center" />
        <el-table-column prop="discount" :label="$t('finance.userLevels.discount')" width="100" align="center">
          <template #default="{ row }">
            {{ row.discount }}%
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('finance.userLevels.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="is_enabled" :label="$t('finance.userLevels.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              :active-value="1"
              :inactive-value="0"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.userLevels.createdAt')" width="170" />
        <el-table-column :label="$t('finance.userLevels.actions')" width="160" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.userLevels.edit') }}</el-button>
            <el-popconfirm :title="$t('finance.userLevels.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.userLevels.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('finance.userLevels.levelName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('finance.userLevels.enterLevelName')" />
        </el-form-item>
        <el-form-item :label="$t('finance.userLevels.levelValue')" prop="level">
          <el-input-number v-model="formData.level" :min="1" :max="999" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('finance.userLevels.discount')" prop="discount">
          <el-input-number v-model="formData.discount" :min="1" :max="100" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('finance.userLevels.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('finance.userLevels.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('finance.userLevels.status')" prop="is_enabled">
          <el-switch v-model="formData.is_enabled" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.userLevels.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.userLevels.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.userLevels.addLevel'))
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  level: 1,
  discount: 100,
  description: '',
  is_enabled: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: $t('finance.userLevels.enterLevelName'), trigger: 'blur' }
  ],
  level: [
    { required: true, message: $t('finance.userLevels.enterLevelValue'), trigger: 'blur' }
  ],
  discount: [
    { required: true, message: $t('finance.userLevels.enterDiscount'), trigger: 'blur' }
  ]
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/user-levels',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchForm.keyword || undefined,
        status: searchForm.status
      }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('Failed to fetch user levels:', error)
    ElMessage.error($t('finance.userLevels.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('finance.userLevels.addLevel')
  formData.id = undefined
  formData.name = ''
  formData.level = 1
  formData.discount = 100
  formData.description = ''
  formData.is_enabled = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.userLevels.editLevel')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/user-levels/${row.id}`,
      params: { is_enabled: row.is_enabled }
    })
    ElMessage.success(row.is_enabled ? $t('finance.userLevels.enabledSuccess') : $t('finance.userLevels.disabledSuccess'))
  } catch (error) {
    row.is_enabled = row.is_enabled ? 0 : 1
    ElMessage.error($t('finance.userLevels.operationFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/user-levels/${row.id}` })
    ElMessage.success($t('finance.userLevels.deleteSuccess'))
    fetchList()
  } catch (error) {
    ElMessage.error($t('finance.userLevels.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/user-levels/${formData.id}`, params: formData })
      } else {
        await request.post({ url: '/api/admin/user-levels', params: formData })
      }
      ElMessage.success(formData.id ? $t('finance.userLevels.updateSuccess') : $t('finance.userLevels.addSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      ElMessage.error($t('finance.userLevels.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchList()
}

const handlePageChange = () => {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.user-levels-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
