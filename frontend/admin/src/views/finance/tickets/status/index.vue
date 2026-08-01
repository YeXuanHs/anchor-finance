<template>
  <div class="ticket-status-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单状态管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增状态
          </el-button>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="状态名称" min-width="150">
          <template #default="{ row }">
            <el-tag :color="row.color" effect="dark" size="small" v-if="row.color">
              {{ row.title }}
            </el-tag>
            <span v-else>{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="color" label="颜色" width="100">
          <template #default="{ row }">
            <el-color-picker v-model="row.color" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="show_active" label="显示活跃" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.show_active ? 'success' : 'info'" size="small">
              {{ row.show_active ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="show_await" label="显示等待" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.show_await ? 'success' : 'info'" size="small">
              {{ row.show_await ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="auto_close" label="自动关闭" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.auto_close ? 'success' : 'info'" size="small">
              {{ row.auto_close ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="order" label="排序" width="80" align="center" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm
              v-if="row.id > 5"
              title="确定删除该状态吗？删除后关联工单将变为关闭状态。"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="状态标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入状态标题" :disabled="formData.id <= 5" />
        </el-form-item>
        <el-form-item label="状态颜色" prop="color">
          <el-color-picker v-model="formData.color" :disabled="formData.id <= 5" />
        </el-form-item>
        <el-form-item label="显示活跃">
          <el-switch v-model="formData.show_active" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="显示等待">
          <el-switch v-model="formData.show_await" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="自动关闭">
          <el-switch v-model="formData.auto_close" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="排序" prop="order">
          <el-input-number v-model="formData.order" :min="0" :max="999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
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

interface TicketStatus {
  id: number
  title: string
  color: string
  show_active: number
  show_await: number
  auto_close: number
  order: number
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新增状态')
const formRef = ref<FormInstance>()

const tableData = ref<TicketStatus[]>([])

const formData = reactive({
  id: 0,
  title: '',
  color: '',
  show_active: 0,
  show_await: 0,
  auto_close: 0,
  order: 0
})

const formRules: FormRules = {
  title: [
    { required: true, message: '请输入状态标题', trigger: 'blur' }
  ],
  order: [
    { type: 'number', message: '排序只能是数字', trigger: 'blur' }
  ]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/ticket-statuses' })
    tableData.value = data.statuses || data || []
  } catch (error) {
    console.error('获取工单状态列表失败:', error)
    ElMessage.error('获取工单状态列表失败')
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.id = 0
  formData.title = ''
  formData.color = ''
  formData.show_active = 0
  formData.show_await = 0
  formData.auto_close = 0
  formData.order = 0
}

const handleAdd = () => {
  dialogTitle.value = '新增状态'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: TicketStatus) => {
  dialogTitle.value = '编辑状态'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: TicketStatus) => {
  try {
    await request.del({
      url: `/api/admin/ticket-statuses/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({
          url: `/api/admin/ticket-statuses/${formData.id}`,
          params: formData
        })
      } else {
        await request.post({
          url: '/api/admin/ticket-statuses',
          params: formData
        })
      }

      ElMessage.success(formData.id ? '修改成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.ticket-status-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
