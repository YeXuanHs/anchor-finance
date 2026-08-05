<template>
  <div class="promo-plan-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>推广计划</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加计划
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="计划名称" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="CPS" value="cps" />
            <el-option label="CPA" value="cpa" />
            <el-option label="CPC" value="cpc" />
            <el-option label="混合" value="mixed" />
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
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="计划名称" min-width="180" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="commission_rate" label="佣金比例" width="120">
          <template #default="{ row }">
            {{ (row.commission_rate * 100).toFixed(1) }}%
          </template>
        </el-table-column>
        <el-table-column prop="click_count" label="点击数" width="100" />
        <el-table-column prop="conversion_count" label="转化数" width="100" />
        <el-table-column prop="total_commission" label="累计佣金" width="120">
          <template #default="{ row }">
            <span class="text-primary">¥{{ row.total_commission?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleGenerateLink(row)">推广链接</el-button>
            <el-button type="info" link @click="handleViewStats(row)">效果统计</el-button>
            <el-popconfirm title="确定删除该推广计划吗？" @confirm="handleDelete(row)">
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
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="计划名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入计划名称" />
        </el-form-item>
        <el-form-item label="推广类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择推广类型" style="width: 100%">
            <el-option label="CPS (按销售分成)" value="cps" />
            <el-option label="CPA (按行为付费)" value="cpa" />
            <el-option label="CPC (按点击付费)" value="cpc" />
            <el-option label="混合模式" value="mixed" />
          </el-select>
        </el-form-item>
        <el-form-item label="佣金比例" prop="commission_rate">
          <el-input-number v-model="formData.commission_rate" :min="0" :max="1" :step="0.01" :precision="2" style="width: 100%" />
          <div class="form-tip">输入0-1之间的小数，如0.15表示15%</div>
        </el-form-item>
        <el-form-item label="计划描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入计划描述" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 推广链接对话框 -->
    <el-dialog v-model="linkDialogVisible" title="推广链接" width="600px">
      <el-form label-width="100px">
        <el-form-item label="计划名称">
          <el-input :value="currentPlan.name" disabled />
        </el-form-item>
        <el-form-item label="推广链接">
          <el-input v-model="promoLink" readonly>
            <template #append>
              <el-button @click="handleCopyLink">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="短链接">
          <el-input v-model="shortLink" readonly>
            <template #append>
              <el-button @click="handleCopyShortLink">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="二维码">
          <div class="qrcode-container">
            <el-image :src="qrcodeUrl" style="width: 160px; height: 160px" fit="contain" />
          </div>
        </el-form-item>
      </el-form>
    </el-dialog>

    <!-- 效果统计对话框 -->
    <el-dialog v-model="statsDialogVisible" title="推广效果统计" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="计划名称">{{ currentPlan.name }}</el-descriptions-item>
        <el-descriptions-item label="推广类型">{{ getTypeText(currentPlan.type) }}</el-descriptions-item>
        <el-descriptions-item label="总点击数">{{ statsData.click_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="总转化数">{{ statsData.conversion_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="转化率">{{ statsData.conversion_rate || '0.00' }}%</el-descriptions-item>
        <el-descriptions-item label="累计佣金">¥{{ statsData.total_commission?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="本月点击">{{ statsData.month_click_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="本月转化">{{ statsData.month_conversion_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="本月佣金">¥{{ statsData.month_commission?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="活跃推广员">{{ statsData.active_promoters || 0 }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'PromoPlan' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  type: undefined as string | undefined,
  status: undefined as number | undefined
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加推广计划')
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  type: 'cps',
  commission_rate: 0.1,
  description: '',
  status: 1
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入计划名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择推广类型', trigger: 'change' }
  ],
  commission_rate: [
    { required: true, message: '请输入佣金比例', trigger: 'blur' }
  ]
}

// 推广链接对话框
const linkDialogVisible = ref(false)
const currentPlan = ref<any>({})
const promoLink = ref('')
const shortLink = ref('')
const qrcodeUrl = ref('')

// 效果统计对话框
const statsDialogVisible = ref(false)
const statsData = ref<any>({})

// 获取类型标签
const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    cps: 'success',
    cpa: 'primary',
    cpc: 'warning',
    mixed: 'info'
  }
  return map[type] || 'info'
}

// 获取类型文本
const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    cps: 'CPS',
    cpa: 'CPA',
    cpc: 'CPC',
    mixed: '混合'
  }
  return map[type] || '未知'
}

// 获取推广计划列表
const fetchPlans = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/promo-plans',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取推广计划列表失败:', error)
    ElMessage.error('获取推广计划列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchPlans()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = undefined
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = '添加推广计划'
  formData.id = undefined
  formData.name = ''
  formData.type = 'cps'
  formData.commission_rate = 0.1
  formData.description = ''
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑推广计划'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 生成推广链接
const handleGenerateLink = async (row: any) => {
  currentPlan.value = row
  try {
    const data = await request.get({
      url: `/api/admin/promo-plans/${row.id}/link`
    })
    promoLink.value = data.link || ''
    shortLink.value = data.short_link || ''
    qrcodeUrl.value = data.qrcode_url || ''
    linkDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取推广链接失败')
  }
}

// 复制推广链接
const handleCopyLink = () => {
  navigator.clipboard.writeText(promoLink.value)
  ElMessage.success('推广链接已复制')
}

// 复制短链接
const handleCopyShortLink = () => {
  navigator.clipboard.writeText(shortLink.value)
  ElMessage.success('短链接已复制')
}

// 查看效果统计
const handleViewStats = async (row: any) => {
  currentPlan.value = row
  try {
    const data = await request.get({
      url: `/api/admin/promo-plans/${row.id}/stats`
    })
    statsData.value = data || {}
    statsDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取统计数据失败')
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/promo-plans/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchPlans()
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
      const url = formData.id ? `/api/admin/promo-plans/${formData.id}` : '/api/admin/promo-plans'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchPlans()
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
  fetchPlans()
}

// 页码变化
const handlePageChange = () => {
  fetchPlans()
}

onMounted(() => {
  fetchPlans()
})
</script>

<style scoped lang="scss">
.promo-plan-page {
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

.text-primary {
  color: #409eff;
  font-weight: 600;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.qrcode-container {
  text-align: center;
  padding: 10px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}
</style>
