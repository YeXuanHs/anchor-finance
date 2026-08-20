<template>
  <div class="client-track-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsTrack.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('clientsTrack.addRecord') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('clientsTrack.client')">
          <el-select
            v-model="searchForm.client_id"
            :placeholder="$t('clientsTrack.selectClient')"
            clearable
            filterable
            remote
            :remote-method="searchClients"
            :loading="clientSearchLoading"
          >
            <el-option
              v-for="client in clientOptions"
              :key="client.id"
              :label="client.username"
              :value="client.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsTrack.phoneCommunication')" value="phone" />
            <el-option :label="$t('clientsTrack.emailContact')" value="email" />
            <el-option :label="$t('clientsTrack.wechatQQ')" value="im" />
            <el-option :label="$t('clientsTrack.visit')" value="visit" />
            <el-option :label="$t('clientsTrack.other')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.dateRange')">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            :range-separator="$t('common.to')"
            :start-placeholder="$t('common.startDate')"
            :end-placeholder="$t('common.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="client_name" :label="$t('clientsTrack.client')" width="120" />
        <el-table-column prop="type" :label="$t('clientsTrack.trackType')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTrackTypeTag(row.type)" size="small">
              {{ getTrackTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" :label="$t('clientsTrack.trackContent')" min-width="250" show-overflow-tooltip />
        <el-table-column prop="operator" :label="$t('common.operator')" width="100" />
        <el-table-column prop="next_follow_at" :label="$t('clientsTrack.nextFollow')" width="180" />
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('common.action')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="primary" link @click="handleAddRemark(row)">{{ $t('clientsTrack.addRemark') }}</el-button>
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
        <el-form-item :label="$t('clientsTrack.client')" prop="client_id">
          <el-select
            v-model="formData.client_id"
            :placeholder="$t('clientsTrack.selectClient')"
            filterable
            remote
            :remote-method="searchClients"
            :loading="clientSearchLoading"
            style="width: 100%"
          >
            <el-option
              v-for="client in clientOptions"
              :key="client.id"
              :label="client.username"
              :value="client.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsTrack.trackType')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('clientsTrack.selectTrackType')" style="width: 100%">
            <el-option :label="$t('clientsTrack.phoneCommunication')" value="phone" />
            <el-option :label="$t('clientsTrack.emailContact')" value="email" />
            <el-option :label="$t('clientsTrack.wechatQQ')" value="im" />
            <el-option :label="$t('clientsTrack.visit')" value="visit" />
            <el-option :label="$t('clientsTrack.other')" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsTrack.trackContent')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="4" :placeholder="$t('clientsTrack.enterTrackContent')" />
        </el-form-item>
        <el-form-item :label="$t('clientsTrack.nextFollow')" prop="next_follow_at">
          <el-date-picker
            v-model="formData.next_follow_at"
            type="datetime"
            :placeholder="$t('clientsTrack.selectNextFollowTime')"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 添加备注对话框 -->
    <el-dialog v-model="remarkDialogVisible" :title="$t('clientsTrack.addRemark')" width="500px">
      <el-form :model="remarkForm" :rules="remarkRules" ref="remarkFormRef" label-width="80px">
        <el-form-item :label="$t('clientsTrack.remarkContent')" prop="remark">
          <el-input v-model="remarkForm.remark" type="textarea" :rows="4" :placeholder="$t('clientsTrack.enterRemarkContent')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="remarkDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleRemarkSubmit" :loading="remarkSubmitLoading">{{ $t('common.confirm') }}</el-button>
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

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)
const remarkSubmitLoading = ref(false)
const clientSearchLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  client_id: undefined as number | undefined,
  type: '',
  dateRange: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 客户选项
const clientOptions = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref($t('clientsTrack.addRecord'))
const formRef = ref<FormInstance>()

// 备注对话框
const remarkDialogVisible = ref(false)
const remarkFormRef = ref<FormInstance>()
const currentTrackId = ref<number>()

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  client_id: undefined as number | undefined,
  type: '',
  content: '',
  next_follow_at: ''
})

// 备注表单
const remarkForm = reactive({
  remark: ''
})

// 表单验证规则
const formRules: FormRules = {
  client_id: [
    { required: true, message: $t('clientsTrack.selectClient'), trigger: 'change' }
  ],
  type: [
    { required: true, message: $t('clientsTrack.selectTrackType'), trigger: 'change' }
  ],
  content: [
    { required: true, message: $t('clientsTrack.enterTrackContent'), trigger: 'blur' }
  ]
}

// 备注验证规则
const remarkRules: FormRules = {
  remark: [
    { required: true, message: $t('clientsTrack.enterRemarkContent'), trigger: 'blur' }
  ]
}

// 获取跟踪类型文本
const getTrackTypeText = (type: string) => {
  const map: Record<string, string> = {
    phone: $t('clientsTrack.phoneCommunication'),
    email: $t('clientsTrack.emailContact'),
    im: $t('clientsTrack.wechatQQ'),
    visit: $t('clientsTrack.visit'),
    other: $t('clientsTrack.other')
  }
  return map[type] || $t('clientsTrack.unknown')
}

// 获取跟踪类型标签类型
const getTrackTypeTag = (type: string) => {
  const map: Record<string, any> = {
    phone: 'primary',
    email: 'success',
    im: 'warning',
    visit: 'danger',
    other: 'info'
  }
  return map[type] || 'info'
}

// 搜索客户
const searchClients = async (query: string) => {
  if (!query) return
  clientSearchLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/users',
      params: { keyword: query, page_size: 20 }
    })
    clientOptions.value = data.list || []
  } catch (error) {
    console.error('搜索客户失败:', error)
  } finally {
    clientSearchLoading.value = false
  }
}

// 获取跟踪列表
const fetchTracks = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      client_id: searchForm.client_id,
      type: searchForm.type
    }

    if (searchForm.dateRange && searchForm.dateRange.length === 2) {
      params.start_date = searchForm.dateRange[0]
      params.end_date = searchForm.dateRange[1]
    }

    const data = await request.get({
      url: '/api/admin/client-tracks',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取跟踪列表失败:', error)
    ElMessage.error($t('clientsTrack.fetchFailed'))
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchTracks()
}

// 重置
const handleReset = () => {
  searchForm.client_id = undefined
  searchForm.type = ''
  searchForm.dateRange = []
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = $t('clientsTrack.addRecord')
  formData.id = undefined
  formData.client_id = undefined
  formData.type = ''
  formData.content = ''
  formData.next_follow_at = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = $t('clientsTrack.editRecord')
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 添加备注
const handleAddRemark = (row: any) => {
  currentTrackId.value = row.id
  remarkForm.remark = ''
  remarkDialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/client-tracks/${formData.id}` : '/api/admin/client-tracks'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess'))
      dialogVisible.value = false
      fetchTracks()
    } catch (error) {
      ElMessage.error($t('common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

// 提交备注
const handleRemarkSubmit = async () => {
  if (!remarkFormRef.value) return

  await remarkFormRef.value.validate(async (valid) => {
    if (!valid) return

    remarkSubmitLoading.value = true
    try {
      await request.put({
        url: `/api/admin/client-tracks/${currentTrackId.value}`,
        params: { remark: remarkForm.remark }
      })
      ElMessage.success($t('clientsTrack.remarkAddSuccess'))
      remarkDialogVisible.value = false
      fetchTracks()
    } catch (error) {
      ElMessage.error($t('clientsTrack.remarkAddFailed'))
    } finally {
      remarkSubmitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchTracks()
}

// 页码变化
const handlePageChange = () => {
  fetchTracks()
}

onMounted(() => {
  fetchTracks()
})
</script>

<style scoped lang="scss">
.client-track-page {
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
