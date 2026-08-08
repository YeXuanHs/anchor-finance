<template>
  <div class="supplier-list-page">
    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加供应商
          </el-button>
        </div>
        <div class="action-right">
          <el-button circle @click="fetchList">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="供应商名称" min-width="150" />
        <el-table-column prop="api_type" label="API类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getApiTypeTag(row.api_type)" size="small">
              {{ getApiTypeText(row.api_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="api_url" label="API地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'active'"
              :inactive-value="'disabled'"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="products_count" label="产品数" width="80" align="center" />
        <el-table-column prop="last_sync_at" label="最后同步" width="170" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="success" link size="small" @click="handleTest(row)">
              测试连接
            </el-button>
            <el-button type="warning" link size="small" @click="handleSync(row)">
              同步产品
            </el-button>
            <el-button type="info" link size="small" @click="handleViewProducts(row)">
              产品列表
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-form-item label="供应商名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入供应商名称" />
        </el-form-item>

        <el-form-item label="API类型" prop="api_type">
          <el-select v-model="formData.api_type" placeholder="请选择API类型" @change="handleApiTypeChange">
            <el-option label="手动管理" value="manual" />
            <el-option label="智简魔方(zjmf)" value="zjmf" />
            <el-option label="v10" value="v10" />
            <el-option label="锚点财务" value="anchor" />
          </el-select>
        </el-form-item>

        <template v-if="formData.api_type !== 'manual'">
          <el-form-item label="API地址" prop="api_url">
            <el-input v-model="formData.api_url" placeholder="https://api.example.com" />
          </el-form-item>

          <el-form-item label="API密钥" prop="api_key">
            <el-input v-model="formData.api_key" placeholder="请输入API密钥" show-password />
          </el-form-item>

          <el-form-item label="API密码" prop="api_password" v-if="formData.api_type === 'zjmf' || formData.api_type === 'v10'">
            <el-input v-model="formData.api_password" placeholder="请输入API密码" show-password />
          </el-form-item>
        </template>

        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])

// 弹窗
const dialogVisible = ref(false)
const dialogTitle = ref('添加供应商')
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref<number | null>(null)

// 表单数据
const formData = reactive({
  name: '',
  api_type: 'manual',
  api_url: '',
  api_key: '',
  api_password: '',
  description: '',
  status: 'active'
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入供应商名称', trigger: 'blur' }
  ],
  api_type: [
    { required: true, message: '请选择API类型', trigger: 'change' }
  ]
}

// API类型标签
const getApiTypeTag = (type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    manual: 'info',
    zjmf: 'primary',
    v10: 'success',
    anchor: 'warning'
  }
  return map[type] || 'info'
}

// API类型文本
const getApiTypeText = (type: string) => {
  const map: Record<string, string> = {
    manual: '手动管理',
    zjmf: '智简魔方',
    v10: 'v10',
    anchor: '锚点财务'
  }
  return map[type] || '未知'
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/suppliers' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取供应商列表失败:', error)
  } finally {
    loading.value = false
  }
}

// API类型变化
const handleApiTypeChange = (type: string) => {
  if (type === 'manual') {
    formData.api_url = ''
    formData.api_key = ''
    formData.api_password = ''
  }
}

// 添加供应商
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加供应商'
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑供应商
const handleEdit = (row: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑供应商'
  editingId.value = row.id
  Object.assign(formData, {
    name: row.name,
    api_type: row.api_type,
    api_url: row.api_url || '',
    api_key: row.api_key || '',
    api_password: '',
    description: row.description || '',
    status: row.status
  })
  dialogVisible.value = true
}

// 测试连接
const handleTest = async (row: any) => {
  try {
    const data = await request.post({ url: `/api/admin/suppliers/${row.id}/test` })
    if (data.success) {
      ElMessage.success('连接测试成功')
    } else {
      ElMessage.error(`连接测试失败: ${data.message}`)
    }
  } catch (error) {
    console.error('测试连接失败:', error)
  }
}

// 同步产品
const handleSync = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要同步供应商 "${row.name}" 的产品吗？`, '确认同步', {
      type: 'warning'
    })
    const data = await request.post({ url: `/api/admin/suppliers/${row.id}/sync` })
    ElMessage.success(`同步完成，共同步 ${data.synced_count || 0} 个产品`)
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('同步产品失败:', error)
    }
  }
}

// 查看产品
const handleViewProducts = (row: any) => {
  router.push(`/product-server?supplier_id=${row.id}`)
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/suppliers/${row.id}/status`, data: { status: row.status } })
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('更新状态失败:', error)
    fetchList()
  }
}

// 删除供应商
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除供应商 "${row.name}" 吗？此操作不可恢复。`, '确认删除', {
      type: 'warning'
    })
    await request.del({ url: `/api/admin/suppliers/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/suppliers/${editingId.value}`, data: formData })
      ElMessage.success('更新成功')
    } else {
      await request.post({ url: '/api/admin/suppliers', data: formData })
      ElMessage.success('添加成功')
    }

    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

// 重置表单
const resetForm = () => {
  formData.name = ''
  formData.api_type = 'manual'
  formData.api_url = ''
  formData.api_key = ''
  formData.api_password = ''
  formData.description = ''
  formData.status = 'active'
}

// 弹窗关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.supplier-list-page {
  padding: 16px;
}

.action-card {
  margin-bottom: 16px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}
</style>
