<template>
  <div class="crons-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>定时任务管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加任务
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="任务名称" clearable />
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="任务名称" min-width="150" />
        <el-table-column prop="command" label="执行命令" min-width="200" show-overflow-tooltip />
        <el-table-column prop="cron_expression" label="Cron表达式" width="140">
          <template #default="{ row }">
            <el-tag effect="plain" size="small">{{ row.cron_expression }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_run_at" label="上次执行" width="180" />
        <el-table-column prop="last_run_status" label="执行结果" width="100">
          <template #default="{ row }">
            <el-tag
              v-if="row.last_run_status !== undefined"
              :type="row.last_run_status === 1 ? 'success' : 'danger'"
              size="small"
            >
              {{ row.last_run_status === 1 ? '成功' : '失败' }}
            </el-tag>
            <span v-else class="empty-text">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="next_run_at" label="下次执行" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleRun(row)">执行</el-button>
            <el-button type="warning" link @click="handleViewLogs(row)">日志</el-button>
            <el-popconfirm title="确定删除该任务吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="执行命令" prop="command">
          <el-input v-model="formData.command" placeholder="请输入执行命令" />
        </el-form-item>
        <el-form-item label="Cron表达式" prop="cron_expression">
          <el-input v-model="formData.cron_expression" placeholder="如: 0 */6 * * *" />
          <div class="form-tip">格式：分 时 日 月 周，例如 "0 2 * * *" 表示每天凌晨2点</div>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入任务描述" />
        </el-form-item>
        <el-form-item label="超时时间(秒)" prop="timeout">
          <el-input-number v-model="formData.timeout" :min="0" :max="86400" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="任务执行日志" width="800px" top="5vh">
      <el-table :data="logData" v-loading="logLoading" style="width: 100%" max-height="400">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="started_at" label="开始时间" width="180" />
        <el-table-column prop="finished_at" label="结束时间" width="180" />
        <el-table-column prop="duration" label="耗时" width="100">
          <template #default="{ row }">
            {{ row.duration ? `${row.duration}s` : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="output" label="输出" min-width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewOutput(row)">查看输出</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 输出详情对话框 -->
    <el-dialog v-model="outputDialogVisible" title="执行输出" width="700px">
      <pre class="output-content">{{ currentOutput }}</pre>
      <template #footer>
        <el-button @click="outputDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

interface Cron {
  id: number
  name: string
  command: string
  cron_expression: string
  description: string
  timeout: number
  status: number
  last_run_at: string
  last_run_status: number
  next_run_at: string
}

interface CronLog {
  id: number
  started_at: string
  finished_at: string
  duration: number
  status: number
  output: string
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加任务')
const formRef = ref<FormInstance>()

const logDialogVisible = ref(false)
const logLoading = ref(false)
const logData = ref<CronLog[]>([])

const outputDialogVisible = ref(false)
const currentOutput = ref('')

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<Cron[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  command: '',
  cron_expression: '',
  description: '',
  timeout: 300,
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' }
  ],
  command: [
    { required: true, message: '请输入执行命令', trigger: 'blur' }
  ],
  cron_expression: [
    { required: true, message: '请输入Cron表达式', trigger: 'blur' }
  ]
}

// 获取任务列表
const fetchCrons = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/cron-tasks',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取定时任务列表失败:', error)
    ElMessage.error('获取定时任务列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchCrons()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加任务'
  formData.id = undefined
  formData.name = ''
  formData.command = ''
  formData.cron_expression = ''
  formData.description = ''
  formData.timeout = 300
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: Cron) => {
  dialogTitle.value = '编辑任务'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 手动执行
const handleRun = async (row: Cron) => {
  try {
    await ElMessageBox.confirm(`确定立即执行任务「${row.name}」吗？`, '提示', {
      type: 'warning'
    })
    await request.post({
      url: `/api/admin/cron-tasks/${row.id}/run`,
      showSuccessMessage: true
    })
    ElMessage.success('任务已触发执行')
    fetchCrons()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('执行失败')
    }
  }
}

// 查看日志
const handleViewLogs = async (row: Cron) => {
  logDialogVisible.value = true
  logLoading.value = true
  try {
    const data = await request.get({
      url: `/api/admin/cron-tasks/${row.id}/logs`
    })
    logData.value = data || []
  } catch (error) {
    console.error('获取日志失败:', error)
    ElMessage.error('获取日志失败')
  } finally {
    logLoading.value = false
  }
}

// 查看输出
const handleViewOutput = (row: CronLog) => {
  currentOutput.value = row.output || '无输出内容'
  outputDialogVisible.value = true
}

// 删除
const handleDelete = async (row: Cron) => {
  try {
    await request.del({
      url: `/api/admin/cron-tasks/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchCrons()
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
      if (formData.id) {
        await request.put({
          url: `/api/admin/cron-tasks/${formData.id}`,
          params: { ...formData },
          showSuccessMessage: true
        })
      } else {
        await request.post({
          url: '/api/admin/cron-tasks',
          params: { ...formData },
          showSuccessMessage: true
        })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchCrons()
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
  fetchCrons()
}

// 页码变化
const handlePageChange = () => {
  fetchCrons()
}

onMounted(() => {
  fetchCrons()
})
</script>

<style scoped lang="scss">
.crons-page {
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

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin-top: 4px;
}

.empty-text {
  color: var(--el-text-color-placeholder);
  font-size: 13px;
}

.output-content {
  max-height: 400px;
  padding: 16px;
  overflow-y: auto;
  background-color: var(--el-fill-color-light);
  border-radius: 6px;
  font-family: Consolas, Monaco, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
