<template>
  <div class="client-care-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.clientCare.pageTitle') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.clientCare.addRule') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.clientCare.ruleName')">
          <el-input v-model="searchForm.name" :placeholder="$t('finance.clientCare.enterRuleName')" clearable />
        </el-form-item>
        <el-form-item :label="$t('finance.clientCare.careType')">
          <el-select v-model="searchForm.type" :placeholder="$t('finance.clientCare.all')" clearable>
            <el-option :label="$t('finance.clientCare.typeBirthday')" :value="1" />
            <el-option :label="$t('finance.clientCare.typeExpiry')" :value="2" />
            <el-option :label="$t('finance.clientCare.typeSatisfaction')" :value="3" />
            <el-option :label="$t('finance.clientCare.typePromotion')" :value="4" />
            <el-option :label="$t('finance.clientCare.typeCustom')" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('finance.clientCare.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('finance.clientCare.all')" clearable>
            <el-option :label="$t('finance.clientCare.enabled')" :value="1" />
            <el-option :label="$t('finance.clientCare.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.clientCare.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.clientCare.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" :label="$t('finance.clientCare.ruleName')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('finance.clientCare.careType')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getCareTypeTag(row.type)" size="small">
              {{ getCareTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="send_method" :label="$t('finance.clientCare.sendMethod')" width="100" align="center">
          <template #default="{ row }">
            {{ row.send_method === 1 ? $t('finance.clientCare.email') : row.send_method === 2 ? $t('finance.clientCare.sms') : $t('finance.clientCare.notice') }}
          </template>
        </el-table-column>
        <el-table-column prop="trigger_type" :label="$t('finance.clientCare.triggerCondition')" width="120" align="center">
          <template #default="{ row }">
            {{ getTriggerText(row.trigger_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="send_count" :label="$t('finance.clientCare.sentCount')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('finance.clientCare.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? $t('finance.clientCare.enabled') : $t('finance.clientCare.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sent_at" :label="$t('finance.clientCare.lastSent')" width="170" />
        <el-table-column prop="created_at" :label="$t('finance.clientCare.createdAt')" width="170" />
        <el-table-column :label="$t('finance.clientCare.actions')" width="220" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewRecords(row)">{{ $t('finance.clientCare.sendRecords') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.clientCare.edit') }}</el-button>
            <el-popconfirm :title="$t('finance.clientCare.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.clientCare.delete') }}</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('finance.clientCare.ruleName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('finance.clientCare.enterRuleName')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('finance.clientCare.careType')" prop="type">
              <el-select v-model="formData.type" :placeholder="$t('finance.clientCare.selectType')" style="width: 100%">
                <el-option :label="$t('finance.clientCare.typeBirthday')" :value="1" />
                <el-option :label="$t('finance.clientCare.typeExpiry')" :value="2" />
                <el-option :label="$t('finance.clientCare.typeSatisfaction')" :value="3" />
                <el-option :label="$t('finance.clientCare.typePromotion')" :value="4" />
                <el-option :label="$t('finance.clientCare.typeCustom')" :value="5" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('finance.clientCare.sendMethod')" prop="send_method">
              <el-select v-model="formData.send_method" :placeholder="$t('finance.clientCare.selectSendMethod')" style="width: 100%">
                <el-option :label="$t('finance.clientCare.email')" :value="1" />
                <el-option :label="$t('finance.clientCare.sms')" :value="2" />
                <el-option :label="$t('finance.clientCare.notice')" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('finance.clientCare.triggerCondition')" prop="trigger_type">
              <el-select v-model="formData.trigger_type" :placeholder="$t('finance.clientCare.selectTrigger')" style="width: 100%">
                <el-option :label="$t('finance.clientCare.triggerBirthday')" :value="1" />
                <el-option :label="$t('finance.clientCare.triggerExpiry')" :value="2" />
                <el-option :label="$t('finance.clientCare.triggerOrderComplete')" :value="3" />
                <el-option :label="$t('finance.clientCare.triggerManual')" :value="4" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('finance.clientCare.advanceDays')" prop="advance_days" v-if="formData.trigger_type === 2">
              <el-input-number v-model="formData.advance_days" :min="1" :max="30" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('finance.clientCare.messageTitle')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('finance.clientCare.enterMessageTitle')" />
        </el-form-item>
        <el-form-item :label="$t('finance.clientCare.messageContent')" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="6" :placeholder="$t('finance.clientCare.enterMessageContent')" />
        </el-form-item>
        <el-form-item :label="$t('finance.clientCare.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item :label="$t('finance.clientCare.remark')" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="2" :placeholder="$t('finance.clientCare.enterRemark')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.clientCare.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.clientCare.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 发送记录对话框 -->
    <el-dialog v-model="recordsVisible" :title="$t('finance.clientCare.sendRecords')" width="750px" destroy-on-close>
      <el-table :data="recordList" v-loading="recordLoading" style="width: 100%" border size="small">
        <el-table-column prop="id" label="ID" width="60" align="center" />
        <el-table-column prop="client_name" :label="$t('finance.clientCare.client')" width="100" />
        <el-table-column prop="send_method" :label="$t('finance.clientCare.sendMethod')" width="80" align="center">
          <template #default="{ row }">
            {{ row.send_method === 1 ? $t('finance.clientCare.email') : row.send_method === 2 ? $t('finance.clientCare.sms') : $t('finance.clientCare.notice') }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('finance.clientCare.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('finance.clientCare.sendSuccess') : $t('finance.clientCare.sendFailed') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sent_at" :label="$t('finance.clientCare.sentAt')" width="170" />
        <el-table-column prop="remark" :label="$t('finance.clientCare.remark')" min-width="150" show-overflow-tooltip />
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="recordPagination.page"
          v-model:page-size="recordPagination.page_size"
          :page-sizes="[10, 20, 50]"
          :total="recordPagination.total"
          layout="total, sizes, prev, next"
          @size-change="handleRecordSizeChange"
          @current-change="handleRecordPageChange"
        />
      </div>
      <template #footer>
        <el-button @click="recordsVisible = false">{{ $t('finance.clientCare.close') }}</el-button>
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
const recordLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  name: '',
  type: undefined as number | undefined,
  status: undefined as number | undefined
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.clientCare.addRule'))
const formRef = ref<FormInstance>()

// 发送记录对话框
const recordsVisible = ref(false)
const recordList = ref<any[]>([])
const recordPagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})
const currentRuleId = ref(0)

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  type: 1,
  send_method: 1,
  trigger_type: 1,
  advance_days: 7,
  title: '',
  content: '',
  status: 1,
  remark: ''
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: $t('finance.clientCare.enterRuleName'), trigger: 'blur' }
  ],
  type: [
    { required: true, message: $t('finance.clientCare.selectCareType'), trigger: 'change' }
  ],
  send_method: [
    { required: true, message: $t('finance.clientCare.selectSendMethod'), trigger: 'change' }
  ],
  trigger_type: [
    { required: true, message: $t('finance.clientCare.selectTrigger'), trigger: 'change' }
  ],
  title: [
    { required: true, message: $t('finance.clientCare.enterMessageTitle'), trigger: 'blur' }
  ],
  content: [
    { required: true, message: $t('finance.clientCare.enterMessageContent'), trigger: 'blur' }
  ]
}

// 关怀类型标签颜色
const getCareTypeTag = (type: number) => {
  const map: Record<number, string> = {
    1: 'danger',
    2: 'warning',
    3: 'primary',
    4: 'success',
    5: 'info'
  }
  return (map[type] || 'info') as any
}

// 获取关怀类型文本
const getCareTypeText = (type: number) => {
  const map: Record<number, string> = {
    1: $t('finance.clientCare.typeBirthday'),
    2: $t('finance.clientCare.typeExpiry'),
    3: $t('finance.clientCare.typeSatisfaction'),
    4: $t('finance.clientCare.typePromotion'),
    5: $t('finance.clientCare.typeCustom')
  }
  return map[type] || $t('finance.clientCare.unknown')
}

// 获取触发条件文本
const getTriggerText = (type: number) => {
  const map: Record<number, string> = {
    1: $t('finance.clientCare.triggerBirthday'),
    2: $t('finance.clientCare.triggerExpiry'),
    3: $t('finance.clientCare.triggerOrderComplete'),
    4: $t('finance.clientCare.triggerManual')
  }
  return map[type] || $t('finance.clientCare.unknown')
}

// 获取规则列表
const fetchRules = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/client-care/rules',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        name: searchForm.name || undefined,
        type: searchForm.type,
        status: searchForm.status
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取关怀规则失败:', error)
    ElMessage.error($t('finance.clientCare.fetchRulesFailed'))
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchRules()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.type = undefined
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = $t('finance.clientCare.addRule')
  formData.id = undefined
  formData.name = ''
  formData.type = 1
  formData.send_method = 1
  formData.trigger_type = 1
  formData.advance_days = 7
  formData.title = ''
  formData.content = ''
  formData.status = 1
  formData.remark = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.clientCare.editRule')
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 查看发送记录
const handleViewRecords = async (row: any) => {
  currentRuleId.value = row.id
  recordPagination.page = 1
  await fetchRecords()
  recordsVisible.value = true
}

// 获取发送记录
const fetchRecords = async () => {
  recordLoading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/client-care/logs',
      params: {
        rule_id: currentRuleId.value,
        page: recordPagination.page,
        page_size: recordPagination.page_size
      }
    })
    recordList.value = data.list || []
    recordPagination.total = data.total || 0
  } catch (error) {
    console.error('获取发送记录失败:', error)
    ElMessage.error($t('finance.clientCare.fetchRecordsFailed'))
  } finally {
    recordLoading.value = false
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/client-care/rules/${row.id}`
    })
    ElMessage.success($t('finance.clientCare.deleteSuccess'))
    fetchRules()
  } catch (error) {
    ElMessage.error($t('finance.clientCare.deleteFailed'))
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/client-care/rules/${formData.id}` : '/api/admin/client-care/rules'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? $t('finance.clientCare.updateSuccess') : $t('finance.clientCare.addSuccess'))
      dialogVisible.value = false
      fetchRules()
    } catch (error) {
      ElMessage.error($t('finance.clientCare.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchRules()
}

// 页码变化
const handlePageChange = () => {
  fetchRules()
}

// 发送记录分页大小变化
const handleRecordSizeChange = () => {
  recordPagination.page = 1
  fetchRecords()
}

// 发送记录页码变化
const handleRecordPageChange = () => {
  fetchRecords()
}

onMounted(() => {
  fetchRules()
})
</script>

<style scoped lang="scss">
.client-care-page {
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
