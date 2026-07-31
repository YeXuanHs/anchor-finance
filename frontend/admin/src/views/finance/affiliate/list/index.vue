<template>
  <div class="affiliate-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>推介用户列表</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户名/邮箱/手机号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" :value="1" />
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
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="referral_code" label="推介码" width="120" />
        <el-table-column prop="commission_rate" label="佣金比例" width="110" align="center">
          <template #default="{ row }">
            <span class="commission-text">{{ row.commission_rate }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_commission" label="累计佣金" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.total_commission) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="referral_count" label="推介人数" width="100" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-button type="primary" link @click="handleAdjustRate(row)">调整比例</el-button>
            <el-popconfirm
              :title="`确定${row.status === 1 ? '禁用' : '启用'}该推介用户吗？`"
              @confirm="handleToggleStatus(row)"
            >
              <template #reference>
                <el-button :type="row.status === 1 ? 'danger' : 'success'" link>
                  {{ row.status === 1 ? '禁用' : '启用' }}
                </el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="推介用户详情" width="750px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ detailData.username }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ detailData.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ detailData.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="推介码">
          <el-tag type="info">{{ detailData.referral_code }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="佣金比例">
          <span class="commission-text">{{ detailData.commission_rate }}%</span>
        </el-descriptions-item>
        <el-descriptions-item label="累计佣金">
          <span class="amount-text">¥{{ formatAmount(detailData.total_commission) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="推介人数">{{ detailData.referral_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailData.status === 1 ? 'success' : 'danger'" size="small">
            {{ detailData.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="上级用户" :span="2">
          {{ detailData.parent_username || '无（顶级推介人）' }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '无' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 下级客户列表 -->
      <div v-if="detailData.subordinates?.length" class="subordinate-section">
        <el-divider content-position="left">推介客户</el-divider>
        <el-table :data="detailData.subordinates" size="small" border>
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="email" label="邮箱" />
          <el-table-column prop="order_amount" label="订单金额" align="right">
            <template #default="{ row }">
              ¥{{ formatAmount(row.order_amount) }}
            </template>
          </el-table-column>
          <el-table-column prop="commission" label="产生佣金" align="right">
            <template #default="{ row }">
              <span class="amount-text">¥{{ formatAmount(row.commission) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="注册时间" width="170" />
        </el-table>
      </div>

      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 调整佣金比例对话框 -->
    <el-dialog v-model="rateDialogVisible" title="调整佣金比例" width="450px" destroy-on-close>
      <el-form :model="rateForm" :rules="rateFormRules" ref="rateFormRef" label-width="100px">
        <el-form-item label="当前用户">
          <el-input :model-value="rateForm.username" disabled />
        </el-form-item>
        <el-form-item label="当前比例">
          <el-tag type="info">{{ rateForm.current_rate }}%</el-tag>
        </el-form-item>
        <el-form-item label="新比例" prop="commission_rate">
          <el-input-number
            v-model="rateForm.commission_rate"
            :min="0"
            :max="100"
            :precision="1"
            :step="0.5"
            style="width: 200px"
          />
          <span style="margin-left: 8px; color: #909399">%</span>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="rateForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入调整原因（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRateSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: '',
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

// 详情对话框
const detailVisible = ref(false)
const detailData = ref<any>({})

// 佣金比例调整对话框
const rateDialogVisible = ref(false)
const rateFormRef = ref<FormInstance>()
const rateForm = reactive({
  id: 0,
  username: '',
  current_rate: 0,
  commission_rate: 0,
  remark: ''
})
const rateFormRules: FormRules = {
  commission_rate: [
    { required: true, message: '请输入佣金比例', trigger: 'blur' }
  ]
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取推介列表
const fetchAffiliates = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/affiliates',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchForm.keyword || undefined,
        status: searchForm.status
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取推介列表失败:', error)
    ElMessage.error('获取推介列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchAffiliates()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

// 查看详情
const handleViewDetail = async (row: any) => {
  try {
    const data = await request.get({
      url: `/api/admin/affiliates/${row.id}`
    })
    detailData.value = data
    detailVisible.value = true
  } catch (error) {
    detailData.value = { ...row }
    detailVisible.value = true
  }
}

// 调整佣金比例
const handleAdjustRate = (row: any) => {
  rateForm.id = row.id
  rateForm.username = row.username
  rateForm.current_rate = row.commission_rate
  rateForm.commission_rate = row.commission_rate
  rateForm.remark = ''
  rateDialogVisible.value = true
}

// 提交佣金比例调整
const handleRateSubmit = async () => {
  if (!rateFormRef.value) return

  await rateFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.put({
        url: `/api/admin/affiliates/${rateForm.id}`,
        params: {
          commission_rate: rateForm.commission_rate,
          remark: rateForm.remark || undefined
        }
      })
      ElMessage.success('佣金比例调整成功')
      rateDialogVisible.value = false
      fetchAffiliates()
    } catch (error) {
      ElMessage.error('调整失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/affiliates/${row.id}`,
      params: { status: row.status === 1 ? 0 : 1 }
    })
    ElMessage.success(row.status === 1 ? '已禁用' : '已启用')
    fetchAffiliates()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchAffiliates()
}

// 页码变化
const handlePageChange = () => {
  fetchAffiliates()
}

onMounted(() => {
  fetchAffiliates()
})
</script>

<style scoped lang="scss">
.affiliate-list-page {
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

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.commission-text {
  font-weight: 600;
  color: #409eff;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.subordinate-section {
  margin-top: 16px;
}
</style>
