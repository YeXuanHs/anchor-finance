<template>
  <div class="client-care-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户关怀</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加规则
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="规则名称">
          <el-input v-model="searchForm.name" placeholder="请输入规则名称" clearable />
        </el-form-item>
        <el-form-item label="关怀类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="生日祝福" :value="1" />
            <el-option label="到期提醒" :value="2" />
            <el-option label="满意度回访" :value="3" />
            <el-option label="优惠通知" :value="4" />
            <el-option label="自定义" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="规则名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="type" label="关怀类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getCareTypeTag(row.type)" size="small">
              {{ getCareTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="send_method" label="发送方式" width="100" align="center">
          <template #default="{ row }">
            {{ row.send_method === 1 ? '邮件' : row.send_method === 2 ? '短信' : '站内信' }}
          </template>
        </el-table-column>
        <el-table-column prop="trigger_type" label="触发条件" width="120" align="center">
          <template #default="{ row }">
            {{ getTriggerText(row.trigger_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="send_count" label="已发送" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sent_at" label="最近发送" width="170" />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="220" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewRecords(row)">发送记录</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该规则吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
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
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入规则名称" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="关怀类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
                <el-option label="生日祝福" :value="1" />
                <el-option label="到期提醒" :value="2" />
                <el-option label="满意度回访" :value="3" />
                <el-option label="优惠通知" :value="4" />
                <el-option label="自定义" :value="5" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="发送方式" prop="send_method">
              <el-select v-model="formData.send_method" placeholder="请选择发送方式" style="width: 100%">
                <el-option label="邮件" :value="1" />
                <el-option label="短信" :value="2" />
                <el-option label="站内信" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="触发条件" prop="trigger_type">
              <el-select v-model="formData.trigger_type" placeholder="请选择触发条件" style="width: 100%">
                <el-option label="生日当天" :value="1" />
                <el-option label="服务到期前" :value="2" />
                <el-option label="订单完成后" :value="3" />
                <el-option label="手动发送" :value="4" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="提前天数" prop="advance_days" v-if="formData.trigger_type === 2">
              <el-input-number v-model="formData.advance_days" :min="1" :max="30" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="消息标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入消息标题" />
        </el-form-item>
        <el-form-item label="消息内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="6" placeholder="请输入消息内容，支持变量：{username} {date}" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 发送记录对话框 -->
    <el-dialog v-model="recordsVisible" title="发送记录" width="750px" destroy-on-close>
      <el-table :data="recordList" v-loading="recordLoading" style="width: 100%" border size="small">
        <el-table-column prop="id" label="ID" width="60" align="center" />
        <el-table-column prop="client_name" label="客户" width="100" />
        <el-table-column prop="send_method" label="发送方式" width="80" align="center">
          <template #default="{ row }">
            {{ row.send_method === 1 ? '邮件' : row.send_method === 2 ? '短信' : '站内信' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sent_at" label="发送时间" width="170" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
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
        <el-button @click="recordsVisible = false">关闭</el-button>
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
const dialogTitle = ref('添加规则')
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
    { required: true, message: '请输入规则名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择关怀类型', trigger: 'change' }
  ],
  send_method: [
    { required: true, message: '请选择发送方式', trigger: 'change' }
  ],
  trigger_type: [
    { required: true, message: '请选择触发条件', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入消息标题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入消息内容', trigger: 'blur' }
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
    1: '生日祝福',
    2: '到期提醒',
    3: '满意度回访',
    4: '优惠通知',
    5: '自定义'
  }
  return map[type] || '未知'
}

// 获取触发条件文本
const getTriggerText = (type: number) => {
  const map: Record<number, string> = {
    1: '生日当天',
    2: '服务到期前',
    3: '订单完成后',
    4: '手动发送'
  }
  return map[type] || '未知'
}

// 获取规则列表
const fetchRules = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/client-care',
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
    ElMessage.error('获取关怀规则失败')
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
  dialogTitle.value = '添加规则'
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
  dialogTitle.value = '编辑规则'
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
      url: `/api/admin/client-care/${currentRuleId.value}/records`,
      params: {
        page: recordPagination.page,
        page_size: recordPagination.page_size
      }
    })
    recordList.value = data.list || []
    recordPagination.total = data.total || 0
  } catch (error) {
    console.error('获取发送记录失败:', error)
    ElMessage.error('获取发送记录失败')
  } finally {
    recordLoading.value = false
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/client-care/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchRules()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/client-care/${formData.id}` : '/api/admin/client-care'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchRules()
    } catch (error) {
      ElMessage.error('操作失败')
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
