<template>
  <div class="notification-templates-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('notificationTemplates.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('notificationTemplates.addTemplate') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('notificationTemplates.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('notificationTemplates.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('notificationTemplates.type')">
          <el-select v-model="searchForm.type" :placeholder="$t('notificationTemplates.all')" clearable>
            <el-option :label="$t('notificationTemplates.email')" value="email" />
            <el-option :label="$t('notificationTemplates.sms')" value="sms" />
            <el-option :label="$t('notificationTemplates.internal')" value="internal" />
            <el-option :label="$t('notificationTemplates.wechat')" value="wechat" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('notificationTemplates.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('notificationTemplates.all')" clearable>
            <el-option :label="$t('notificationTemplates.enabled')" :value="1" />
            <el-option :label="$t('notificationTemplates.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('notificationTemplates.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('notificationTemplates.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" :label="$t('notificationTemplates.id')" width="80" align="center" />
        <el-table-column prop="name" :label="$t('notificationTemplates.templateName')" width="180" />
        <el-table-column prop="event" :label="$t('notificationTemplates.triggerEvent')" width="200" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('notificationTemplates.typeColumn')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">{{ typeTextMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" :label="$t('notificationTemplates.subject')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('notificationTemplates.statusColumn')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('notificationTemplates.enabled') : $t('notificationTemplates.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" :label="$t('notificationTemplates.updateTime')" width="180" />
        <el-table-column :label="$t('notificationTemplates.operations')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('notificationTemplates.edit') }}</el-button>
            <el-button type="warning" link @click="handleResetEvent(row)">{{ $t('notificationTemplates.resetEvent') }}</el-button>
            <el-popconfirm :title="$t('notificationTemplates.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('notificationTemplates.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="isEditing ? $t('notificationTemplates.editTemplateTitle') : $t('notificationTemplates.addTemplateTitle')" width="800px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('notificationTemplates.templateName')" prop="name">
              <el-input v-model="formData.name" :placeholder="$t('notificationTemplates.inputTemplateName')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('notificationTemplates.typeColumn')" prop="type">
              <el-select v-model="formData.type" :placeholder="$t('notificationTemplates.selectType')" style="width: 100%">
                <el-option :label="$t('notificationTemplates.email')" value="email" />
                <el-option :label="$t('notificationTemplates.sms')" value="sms" />
                <el-option :label="$t('notificationTemplates.internal')" value="internal" />
                <el-option :label="$t('notificationTemplates.wechat')" value="wechat" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('notificationTemplates.triggerEvent')" prop="event">
          <el-input v-model="formData.event" :placeholder="$t('notificationTemplates.inputTriggerEvent')" />
        </el-form-item>
        <el-form-item :label="$t('notificationTemplates.subject')" prop="subject">
          <el-input v-model="formData.subject" :placeholder="$t('notificationTemplates.subject')" />
        </el-form-item>
        <el-form-item label="Content" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="10" :placeholder="$t('notificationTemplates.inputContent')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('notificationTemplates.statusColumn')">
              <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('notificationTemplates.enabled')" :inactive-text="$t('notificationTemplates.disabled')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Sort">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('notificationTemplates.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('notificationTemplates.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'NotificationTemplatesManage' })

const typeTextMap = computed<Record<string, string>>(() => ({
  email: $t('notificationTemplates.email'),
  sms: $t('notificationTemplates.sms'),
  internal: $t('notificationTemplates.internal'),
  wechat: $t('notificationTemplates.wechat')
}))
const typeTagMap: Record<string, any> = { email: 'primary', sms: 'success', internal: 'info', wechat: 'warning' }

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()

const isEditing = computed(() => !!formData.id)

const searchForm = reactive({
  keyword: '',
  type: '',
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
  type: 'email',
  event: '',
  subject: '',
  content: '',
  status: 1,
  sort: 0
})

const formRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('notificationTemplates.inputTemplateName'), trigger: 'blur' }],
  type: [{ required: true, message: $t('notificationTemplates.selectType'), trigger: 'change' }],
  event: [{ required: true, message: $t('notificationTemplates.inputTriggerEvent'), trigger: 'blur' }],
  content: [{ required: true, message: $t('notificationTemplates.inputContent'), trigger: 'blur' }]
}))

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/notifications/templates',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('notificationTemplates.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { keyword: '', type: '', status: undefined }); handleSearch() }

const resetForm = () => {
  formData.id = undefined
  formData.name = ''
  formData.type = 'email'
  formData.event = ''
  formData.subject = ''
  formData.content = ''
  formData.status = 1
  formData.sort = 0
}

const handleAdd = () => {
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleResetEvent = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('notificationTemplates.resetConfirm', { name: row.name }), $t('notificationTemplates.confirmResetEvent'))
    await request.put({ url: `/api/admin/notifications/templates/${row.id}/reset-event` })
    ElMessage.success($t('notificationTemplates.eventResetSuccess'))
    fetchData()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error($t('notificationTemplates.resetFailed'))
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/notifications/templates/${row.id}` })
    ElMessage.success($t('notificationTemplates.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('notificationTemplates.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/notifications/templates/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/notifications/templates', params: { ...formData } })
      }
      ElMessage.success(formData.id ? $t('notificationTemplates.updateSuccess') : $t('notificationTemplates.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('notificationTemplates.operationFailed'))
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
.notification-templates-page {
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
