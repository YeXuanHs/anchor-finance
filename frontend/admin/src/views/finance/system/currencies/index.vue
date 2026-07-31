<template>
  <div class="currencies-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>货币管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加货币
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="货币名称/代码" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="货币代码" width="100">
          <template #default="{ row }">
            <el-tag effect="plain">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="货币名称" min-width="120" />
        <el-table-column prop="symbol" label="符号" width="80" align="center" />
        <el-table-column prop="exchange_rate" label="汇率" width="120" align="right">
          <template #default="{ row }">
            {{ row.exchange_rate?.toFixed(4) || '1.0000' }}
          </template>
        </el-table-column>
        <el-table-column prop="is_default" label="默认货币" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_default ? 'success' : 'info'" size="small">
              {{ row.is_default ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              type="success"
              link
              @click="handleSetDefault(row)"
              :disabled="row.is_default"
            >
              设为默认
            </el-button>
            <el-popconfirm title="确定删除该货币吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link :disabled="row.is_default">删除</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="货币代码" prop="code">
          <el-input v-model="formData.code" placeholder="如 USD、CNY、EUR" :disabled="!!formData.id" />
        </el-form-item>
        <el-form-item label="货币名称" prop="name">
          <el-input v-model="formData.name" placeholder="如 美元、人民币" />
        </el-form-item>
        <el-form-item label="货币符号" prop="symbol">
          <el-input v-model="formData.symbol" placeholder="如 $、¥、€" />
        </el-form-item>
        <el-form-item label="汇率" prop="exchange_rate">
          <el-input-number
            v-model="formData.exchange_rate"
            :min="0"
            :precision="4"
            :step="0.0001"
            style="width: 100%"
          />
          <div class="form-tip">相对于默认货币的汇率，默认货币汇率为 1.0000</div>
        </el-form-item>
        <el-form-item label="小数位数" prop="decimal_places">
          <el-input-number v-model="formData.decimal_places" :min="0" :max="8" style="width: 100%" />
        </el-form-item>
        <el-form-item label="货币前缀" prop="prefix">
          <el-input v-model="formData.prefix" placeholder="如 ¥" />
        </el-form-item>
        <el-form-item label="货币后缀" prop="suffix">
          <el-input v-model="formData.suffix" placeholder="如 元" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

interface Currency {
  id: number
  code: string
  name: string
  symbol: string
  exchange_rate: number
  decimal_places: number
  prefix: string
  suffix: string
  is_default: boolean
  status: number
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加货币')
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<Currency[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  code: '',
  name: '',
  symbol: '',
  exchange_rate: 1,
  decimal_places: 2,
  prefix: '',
  suffix: '',
  status: 1
})

const formRules: FormRules = {
  code: [
    { required: true, message: '请输入货币代码', trigger: 'blur' },
    { min: 2, max: 5, message: '长度在 2 到 5 个字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入货币名称', trigger: 'blur' }
  ],
  symbol: [
    { required: true, message: '请输入货币符号', trigger: 'blur' }
  ],
  exchange_rate: [
    { required: true, message: '请输入汇率', trigger: 'blur' }
  ]
}

// 获取货币列表
const fetchCurrencies = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/currencies',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取货币列表失败:', error)
    ElMessage.error('获取货币列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchCurrencies()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加货币'
  formData.id = undefined
  formData.code = ''
  formData.name = ''
  formData.symbol = ''
  formData.exchange_rate = 1
  formData.decimal_places = 2
  formData.prefix = ''
  formData.suffix = ''
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: Currency) => {
  dialogTitle.value = '编辑货币'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 设为默认
const handleSetDefault = async (row: Currency) => {
  try {
    await ElMessageBox.confirm(`确定将 ${row.name} 设为默认货币吗？`, '提示', {
      type: 'warning'
    })
    await request.put({
      url: `/api/admin/currencies/${row.id}/default`,
      showSuccessMessage: true
    })
    ElMessage.success('设置成功')
    fetchCurrencies()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('设置失败')
    }
  }
}

// 删除
const handleDelete = async (row: Currency) => {
  try {
    await request.del({
      url: `/api/admin/currencies/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchCurrencies()
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
          url: `/api/admin/currencies/${formData.id}`,
          params: { ...formData },
          showSuccessMessage: true
        })
      } else {
        await request.post({
          url: '/api/admin/currencies',
          params: { ...formData },
          showSuccessMessage: true
        })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchCurrencies()
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
  fetchCurrencies()
}

// 页码变化
const handlePageChange = () => {
  fetchCurrencies()
}

onMounted(() => {
  fetchCurrencies()
})
</script>

<style scoped lang="scss">
.currencies-page {
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
</style>
