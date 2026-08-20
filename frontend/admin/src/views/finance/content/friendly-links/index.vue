<template>
  <div class="friendly-links-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.friendlyLinks.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.friendlyLinks.addLink') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.friendlyLinks.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.friendlyLinks.nameOrUrl')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.friendlyLinks.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.friendlyLinks.all')" clearable>
            <el-option :label="$t('finance.friendlyLinks.enabled')" :value="1" />
            <el-option :label="$t('finance.friendlyLinks.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.friendlyLinks.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.friendlyLinks.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('finance.friendlyLinks.name')" width="150" />
        <el-table-column prop="url" :label="$t('finance.friendlyLinks.linkUrl')" min-width="250" show-overflow-tooltip />
        <el-table-column prop="logo" label="Logo" width="100" align="center">
          <template #default="{ row }">
            <el-image v-if="row.logo" :src="row.logo" style="width: 40px; height: 40px" fit="contain" :preview-src-list="[row.logo]" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('finance.friendlyLinks.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('finance.friendlyLinks.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('finance.friendlyLinks.enabled') : $t('finance.friendlyLinks.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.friendlyLinks.createdAt')" width="180" />
        <el-table-column :label="$t('finance.friendlyLinks.actions')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.friendlyLinks.edit') }}</el-button>
            <el-popconfirm :title="$t('finance.friendlyLinks.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.friendlyLinks.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('finance.friendlyLinks.name')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('finance.friendlyLinks.enterLinkName')" />
        </el-form-item>
        <el-form-item :label="$t('finance.friendlyLinks.linkUrl')" prop="url">
          <el-input v-model="formData.url" :placeholder="$t('finance.friendlyLinks.enterLinkUrl')" />
        </el-form-item>
        <el-form-item label="Logo">
          <el-input v-model="formData.logo" :placeholder="$t('finance.friendlyLinks.enterLogoUrl')" />
        </el-form-item>
        <el-form-item :label="$t('finance.friendlyLinks.description')">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('finance.friendlyLinks.enterDescription')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('finance.friendlyLinks.sort')">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('finance.friendlyLinks.status')">
              <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('finance.friendlyLinks.enabled')" :inactive-text="$t('finance.friendlyLinks.disabled')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.friendlyLinks.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.friendlyLinks.confirm') }}</el-button>
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

defineOptions({ name: 'FriendlyLinksManage' })

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.friendlyLinks.addLink'))
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
  url: '',
  logo: '',
  description: '',
  sort: 0,
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: $t('finance.friendlyLinks.enterLinkName'), trigger: 'blur' }
  ],
  url: [
    { required: true, message: $t('finance.friendlyLinks.enterLinkUrl'), trigger: 'blur' },
    { type: 'url', message: $t('finance.friendlyLinks.invalidUrl'), trigger: 'blur' }
  ]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/friendly-links',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('finance.friendlyLinks.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

const resetForm = () => {
  formData.id = undefined
  formData.name = ''
  formData.url = ''
  formData.logo = ''
  formData.description = ''
  formData.sort = 0
  formData.status = 1
}

const handleAdd = () => {
  dialogTitle.value = $t('finance.friendlyLinks.addLink')
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.friendlyLinks.editLink')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/friendly-links/${row.id}` })
    ElMessage.success($t('finance.friendlyLinks.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('finance.friendlyLinks.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/friendly-links/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/friendly-links', params: { ...formData } })
      }
      ElMessage.success(formData.id ? $t('finance.friendlyLinks.updateSuccess') : $t('finance.friendlyLinks.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('finance.friendlyLinks.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.friendly-links-page {
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
