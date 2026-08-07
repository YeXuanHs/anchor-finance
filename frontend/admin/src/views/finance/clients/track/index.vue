<template>
  <div class="client-track-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户跟踪</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加跟踪记录
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="客户">
          <el-select
            v-model="searchForm.client_id"
            placeholder="选择客户"
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
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="电话沟通" value="phone" />
            <el-option label="邮件往来" value="email" />
            <el-option label="微信/QQ" value="im" />
            <el-option label="上门拜访" value="visit" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="client_name" label="客户" width="120" />
        <el-table-column prop="type" label="跟踪类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTrackTypeTag(row.type)" size="small">
              {{ getTrackTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="跟踪内容" min-width="250" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="next_follow_at" label="下次跟进" width="180" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleAddRemark(row)">添加备注</el-button>
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
        <el-form-item label="客户" prop="client_id">
          <el-select
            v-model="formData.client_id"
            placeholder="选择客户"
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
        <el-form-item label="跟踪类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择跟踪类型" style="width: 100%">
            <el-option label="电话沟通" value="phone" />
            <el-option label="邮件往来" value="email" />
            <el-option label="微信/QQ" value="im" />
            <el-option label="上门拜访" value="visit" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="跟踪内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="4" placeholder="请输入跟踪内容" />
        </el-form-item>
        <el-form-item label="下次跟进" prop="next_follow_at">
          <el-date-picker
            v-model="formData.next_follow_at"
            type="datetime"
            placeholder="选择下次跟进时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加备注对话框 -->
    <el-dialog v-model="remarkDialogVisible" title="添加备注" width="500px">
      <el-form :model="remarkForm" :rules="remarkRules" ref="remarkFormRef" label-width="80px">
        <el-form-item label="备注内容" prop="remark">
          <el-input v-model="remarkForm.remark" type="textarea" :rows="4" placeholder="请输入备注内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="remarkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRemarkSubmit" :loading="remarkSubmitLoading">确定</el-button>
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
const dialogTitle = ref('添加跟踪记录')
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
    { required: true, message: '请选择客户', trigger: 'change' }
  ],
  type: [
    { required: true, message: '请选择跟踪类型', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入跟踪内容', trigger: 'blur' }
  ]
}

// 备注验证规则
const remarkRules: FormRules = {
  remark: [
    { required: true, message: '请输入备注内容', trigger: 'blur' }
  ]
}

// 获取跟踪类型文本
const getTrackTypeText = (type: string) => {
  const map: Record<string, string> = {
    phone: '电话沟通',
    email: '邮件往来',
    im: '微信/QQ',
    visit: '上门拜访',
    other: '其他'
  }
  return map[type] || '未知'
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
    ElMessage.error('获取跟踪列表失败')
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
  dialogTitle.value = '添加跟踪记录'
  formData.id = undefined
  formData.client_id = undefined
  formData.type = ''
  formData.content = ''
  formData.next_follow_at = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑跟踪记录'
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

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchTracks()
    } catch (error) {
      ElMessage.error('操作失败')
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
      ElMessage.success('备注添加成功')
      remarkDialogVisible.value = false
      fetchTracks()
    } catch (error) {
      ElMessage.error('添加备注失败')
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
