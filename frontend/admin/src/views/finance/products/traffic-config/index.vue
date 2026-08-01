<template>
  <div class="traffic-config-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>流量包配置</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加流量包
          </el-button>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="流量包名称" min-width="150" />
        <el-table-column prop="traffic_amount" label="流量额度" width="120" align="center">
          <template #default="{ row }">
            {{ formatTraffic(row.traffic_amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.price) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="validity_days" label="有效期" width="100" align="center">
          <template #default="{ row }">
            {{ row.validity_days }}天
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="is_hot" label="热销" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_hot" type="danger" size="small">热销</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleCopy(row)">复制</el-button>
            <el-popconfirm title="确定删除该流量包吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="流量包名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入流量包名称" />
        </el-form-item>
        <el-form-item label="流量额度(MB)" prop="traffic_amount">
          <el-input-number v-model="formData.traffic_amount" :min="1" />
          <span class="suffix-text">MB</span>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="formData.price" :min="0" :precision="2" />
          <span class="suffix-text">元</span>
        </el-form-item>
        <el-form-item label="有效期(天)" prop="validity_days">
          <el-input-number v-model="formData.validity_days" :min="1" />
          <span class="suffix-text">天</span>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="热销标记" prop="is_hot">
          <el-switch v-model="formData.is_hot" :active-value="1" :inactive-value="0" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('添加流量包')
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  traffic_amount: 1024,
  price: 0,
  validity_days: 30,
  description: '',
  sort: 0,
  is_hot: 0,
  status: 1
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入流量包名称', trigger: 'blur' }],
  traffic_amount: [{ required: true, message: '请输入流量额度', trigger: 'blur' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }],
  validity_days: [{ required: true, message: '请输入有效期', trigger: 'blur' }]
}

const formatAmount = (amount: number | undefined) =>
  amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'

const formatTraffic = (mb: number | undefined) => {
  if (!mb) return '0 MB'
  if (mb >= 1024) return `${(mb / 1024).toFixed(1)} GB`
  return `${mb} MB`
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/traffic-config' })
    tableData.value = data.list || data || []
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '添加流量包'
  formData.id = undefined
  formData.name = ''
  formData.traffic_amount = 1024
  formData.price = 0
  formData.validity_days = 30
  formData.description = ''
  formData.sort = 0
  formData.is_hot = 0
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑流量包'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleCopy = (row: any) => {
  dialogTitle.value = '复制流量包'
  Object.assign(formData, row)
  formData.id = undefined
  formData.name = row.name + '_副本'
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/traffic-config/${row.id}` })
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
      const url = formData.id ? `/api/admin/traffic-config/${formData.id}` : '/api/admin/traffic-config'
      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }
      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.traffic-config-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.suffix-text {
  margin-left: 8px;
  color: #909399;
}
</style>
