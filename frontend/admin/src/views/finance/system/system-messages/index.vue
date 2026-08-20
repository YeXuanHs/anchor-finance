<template>
  <div class="system-messages-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('systemMessages.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('systemMessages.sendMessage') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('systemMessages.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('systemMessages.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('systemMessages.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('systemMessages.all')" clearable>
            <el-option :label="$t('systemMessages.systemNotification')" value="system" />
            <el-option :label="$t('systemMessages.activityNotification')" value="activity" />
            <el-option :label="$t('systemMessages.securityAlert')" value="security" />
            <el-option :label="$t('systemMessages.orderNotification')" value="order" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('systemMessages.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('systemMessages.all')" clearable>
            <el-option :label="$t('systemMessages.unread')" value="unread" />
            <el-option :label="$t('systemMessages.read')" value="read" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('systemMessages.dateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('systemMessages.to')"
            :start-placeholder="$t('systemMessages.startDate')"
            :end-placeholder="$t('systemMessages.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('systemMessages.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('systemMessages.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" :label="$t('systemMessages.id')" width="80" align="center" />
        <el-table-column prop="title" :label="$t('systemMessages.titleColumn')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('systemMessages.type')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">{{ typeTextMap[row.type as keyof typeof typeTextMap] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target" :label="$t('systemMessages.target')" width="120" />
        <el-table-column prop="is_read" :label="$t('systemMessages.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_read ? 'info' : 'danger'" size="small">
              {{ row.is_read ? $t('systemMessages.read') : $t('systemMessages.unread') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('systemMessages.sendTime')" width="180" />
        <el-table-column :label="$t('systemMessages.operations')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">{{ $t('systemMessages.detail') }}</el-button>
            <el-button v-if="!row.is_read" type="success" link @click="handleMarkRead(row)">{{ $t('systemMessages.markRead') }}</el-button>
            <el-popconfirm :title="$t('systemMessages.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('systemMessages.delete') }}</el-button>
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

    <!-- 发送消息对话框 -->
    <el-dialog v-model="dialogVisible" :title="$t('systemMessages.sendSystemMessage')" width="700px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('systemMessages.messageTitle')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('systemMessages.inputTitle')" />
        </el-form-item>
        <el-form-item :label="$t('systemMessages.type')" prop="type">
          <el-select v-model="formData.type" :placeholder="$t('systemMessages.selectType')" style="width: 100%">
            <el-option :label="$t('systemMessages.systemNotification')" value="system" />
            <el-option :label="$t('systemMessages.activityNotification')" value="activity" />
            <el-option :label="$t('systemMessages.securityAlert')" value="security" />
            <el-option :label="$t('systemMessages.orderNotification')" value="order" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('systemMessages.target')">
          <el-select v-model="formData.target" :placeholder="$t('systemMessages.allUsers')" clearable style="width: 100%">
            <el-option :label="$t('systemMessages.allUsers')" value="all" />
            <el-option :label="$t('systemMessages.specificUsers')" value="specific" />
            <el-option :label="$t('systemMessages.specificGroup')" value="group" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="formData.target === 'specific'" :label="$t('systemMessages.userId')">
          <el-input v-model="formData.user_ids" :placeholder="$t('systemMessages.userIdPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('systemMessages.content')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="8" :placeholder="$t('systemMessages.inputContent')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('systemMessages.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('systemMessages.send') }}</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="$t('systemMessages.messageDetail')" width="600px">
      <el-descriptions :column="1" border v-if="detailData">
        <el-descriptions-item :label="$t('systemMessages.id')">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('systemMessages.titleColumn')">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item :label="$t('systemMessages.type')">
          <el-tag :type="typeTagMap[detailData.type]" size="small">{{ typeTextMap[detailData.type as keyof typeof typeTextMap] || detailData.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('systemMessages.target')">{{ detailData.target }}</el-descriptions-item>
        <el-descriptions-item :label="$t('systemMessages.status')">
          <el-tag :type="detailData.is_read ? 'info' : 'danger'" size="small">
            {{ detailData.is_read ? $t('systemMessages.read') : $t('systemMessages.unread') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('systemMessages.content')">
          <div style="white-space: pre-wrap;">{{ detailData.content }}</div>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('systemMessages.sendTime')">{{ detailData.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('systemMessages.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'SystemMessagesManage' })

const typeTextMap = computed(() => ({
  system: $t('systemMessages.systemNotification'),
  activity: $t('systemMessages.activityNotification'),
  security: $t('systemMessages.securityAlert'),
  order: $t('systemMessages.orderNotification')
}))
const typeTagMap: Record<string, any> = { system: 'primary', activity: 'success', security: 'danger', order: 'warning' }

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  type: '',
  status: '',
  date_range: [] as string[]
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formData = reactive({
  title: '',
  type: 'system',
  target: 'all',
  user_ids: '',
  content: ''
})

const formRules = computed<FormRules>(() => ({
  title: [{ required: true, message: $t('systemMessages.titleRequired'), trigger: 'blur' }],
  type: [{ required: true, message: $t('systemMessages.typeRequired'), trigger: 'change' }],
  content: [{ required: true, message: $t('systemMessages.contentRequired'), trigger: 'blur' }]
}))

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword || undefined,
      type: searchForm.type || undefined,
      status: searchForm.status || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/system-messages', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('systemMessages.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', type: '', status: '', date_range: [] }); handleSearch() }

const handleAdd = () => {
  Object.assign(formData, { title: '', type: 'system', target: 'all', user_ids: '', content: '' })
  dialogVisible.value = true
}

const handleDetail = (row: any) => { detailData.value = row; detailVisible.value = true }

const handleMarkRead = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/system-messages/${row.id}/read` })
    ElMessage.success($t('systemMessages.markSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('systemMessages.operationFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/system-messages/${row.id}` })
    ElMessage.success($t('systemMessages.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('systemMessages.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await request.post({ url: '/api/admin/system-messages', params: { ...formData } })
      ElMessage.success($t('systemMessages.sendSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('systemMessages.sendFailed'))
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
.system-messages-page {
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
