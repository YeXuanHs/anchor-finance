<template>
  <div class="config-options-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.configOptions.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('page.configOptions.addConfig') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('page.configOptions.group')">
          <el-select v-model="searchForm.group" :placeholder="$t('page.configOptions.all')" clearable>
            <el-option v-for="g in groups" :key="g" :label="g" :value="g" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('page.configOptions.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('page.configOptions.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('page.configOptions.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('page.configOptions.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="key" :label="$t('page.configOptions.configKey')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="value" :label="$t('page.configOptions.configValue')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="group" :label="$t('page.configOptions.group')" width="120" />
        <el-table-column prop="description" :label="$t('page.configOptions.description')" min-width="150" show-overflow-tooltip />
        <el-table-column :label="$t('page.configOptions.actions')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('page.configOptions.edit') }}</el-button>
            <el-popconfirm :title="$t('page.configOptions.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('page.configOptions.delete') }}</el-button>
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
        <el-form-item :label="$t('page.configOptions.configKey')" prop="key">
          <el-input v-model="formData.key" :placeholder="$t('page.configOptions.pleaseConfigKey')" />
        </el-form-item>
        <el-form-item :label="$t('page.configOptions.configValue')" prop="value">
          <el-input v-model="formData.value" type="textarea" :rows="3" :placeholder="$t('page.configOptions.pleaseConfigValue')" />
        </el-form-item>
        <el-form-item :label="$t('page.configOptions.group')" prop="group">
          <el-input v-model="formData.group" :placeholder="$t('page.configOptions.pleaseGroup')" />
        </el-form-item>
        <el-form-item :label="$t('page.configOptions.description')">
          <el-input v-model="formData.description" :placeholder="$t('page.configOptions.pleaseDescription')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('page.configOptions.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('page.configOptions.confirm') }}</el-button>
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
const dialogTitle = ref($t('page.configOptions.addConfig'))
const formRef = ref<FormInstance>()

const groups = ref<string[]>([])
const tableData = ref<any[]>([])

const searchForm = reactive({
  group: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const formData = reactive({
  id: undefined as number | undefined,
  key: '',
  value: '',
  group: '',
  description: ''
})

const formRules: FormRules = {
  key: [{ required: true, message: $t('page.configOptions.pleaseConfigKey'), trigger: 'blur' }],
  group: [{ required: true, message: $t('page.configOptions.pleaseGroup'), trigger: 'blur' }]
}

const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/config-options/groups-list' })
    groups.value = data || []
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/config-options/search-page',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('page.configOptions.fetchConfigFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.group = ''; searchForm.keyword = ''; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = $t('page.configOptions.addConfig')
  formData.id = undefined; formData.key = ''; formData.value = ''; formData.group = ''; formData.description = ''
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('page.configOptions.editConfig')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: '/api/admin/config-options/options', params: { id: row.id } })
    ElMessage.success($t('page.configOptions.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('page.configOptions.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.post({ url: '/api/admin/config-options/edit-config', params: formData })
      } else {
        await request.post({ url: '/api/admin/config-options/add-options', params: formData })
      }
      ElMessage.success(formData.id ? $t('page.configOptions.updateSuccess') : $t('page.configOptions.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('page.configOptions.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => { fetchGroups(); fetchData() })
</script>

<style scoped lang="scss">
.config-options-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
