<template>
  <div class="resource-pool-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsResourcePool.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('clientsResourcePool.addPool') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsResourcePool.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('clientsResourcePool.poolNamePlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsResourcePool.ipPool')" :value="1" />
            <el-option :label="$t('clientsResourcePool.domainPool')" :value="2" />
            <el-option :label="$t('clientsResourcePool.certPool')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('clientsResourcePool.poolName')" min-width="150" />
        <el-table-column prop="type_text" :label="$t('common.type')" width="100" align="center" />
        <el-table-column prop="total_count" :label="$t('clientsResourcePool.totalCount')" width="100" align="center" />
        <el-table-column prop="used_count" :label="$t('clientsResourcePool.usedCount')" width="100" align="center" />
        <el-table-column prop="available_count" :label="$t('clientsResourcePool.availableCount')" width="100" align="center" />
        <el-table-column prop="usage_rate" :label="$t('clientsResourcePool.usageRate')" width="120" align="center">
          <template #default="{ row }">
            <el-progress :percentage="row.usage_rate || 0" :stroke-width="8" />
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('common.enabled') : $t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="170" />
        <el-table-column :label="$t('common.action')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="primary" link @click="handleViewDetail(row)">{{ $t('common.detail') }}</el-button>
            <el-popconfirm :title="$t('clientsResourcePool.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('common.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('clientsResourcePool.poolName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('clientsResourcePool.enterPoolName')" />
        </el-form-item>
        <el-form-item :label="$t('common.type')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('clientsResourcePool.selectType')">
            <el-option :label="$t('clientsResourcePool.ipPool')" :value="1" />
            <el-option :label="$t('clientsResourcePool.domainPool')" :value="2" />
            <el-option :label="$t('clientsResourcePool.certPool')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('clientsResourcePool.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('common.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const submitLoading = ref(false)

const searchForm = reactive({
  keyword: '',
  type: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref($t('clientsResourcePool.addPool'))
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  type: 1,
  description: '',
  status: 1
})

const formRules: FormRules = {
  name: [{ required: true, message: $t('clientsResourcePool.enterPoolName'), trigger: 'blur' }],
  type: [{ required: true, message: $t('clientsResourcePool.selectType'), trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/resource-pools',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error($t('common.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = $t('clientsResourcePool.addPool')
  formData.id = undefined
  formData.name = ''
  formData.type = 1
  formData.description = ''
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('clientsResourcePool.editPool')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleViewDetail = (row: any) => {
  dialogTitle.value = $t('clientsResourcePool.poolDetail')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/resource-pools/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('common.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/resource-pools/${formData.id}` : '/api/admin/resource-pools'
      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.resource-pool-page {
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
