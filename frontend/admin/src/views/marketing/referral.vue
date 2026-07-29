<template>
  <div class="referral-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="用户">
          <el-input v-model="searchForm.username" placeholder="推介人/被推介人" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待结算" value="pending" />
            <el-option label="已支付" value="paid" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>推荐返利</h3>
        <div>
          <el-button @click="showConfigDialog = true">
            <el-icon><Setting /></el-icon>
            返利配置
          </el-button>
          <el-button @click="showRuleDialog = true">
            <el-icon><List /></el-icon>
            返利规则
          </el-button>
          <el-button @click="exportData">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </div>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="referrer_name" label="推介人" width="120" />
        <el-table-column prop="referred_name" label="被推介人" width="120" />
        <el-table-column prop="referred_order" label="关联订单" width="160" />
        <el-table-column prop="order_amount" label="订单金额" width="120">
          <template #default="{ row }">
            <span>¥{{ row.order_amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="commission" label="返利金额" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ row.commission?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusTextMap[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="paid_at" label="支付时间" width="180" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button type="success" link @click="payCommission(row)" v-if="row.status === 'pending'">
              支付
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchData"
          @size-change="fetchData"
        />
      </div>
    </div>

    <el-dialog v-model="showConfigDialog" title="返利配置" width="550px">
      <el-form :model="configForm" label-width="120px">
        <el-form-item label="启用推荐返利">
          <el-switch v-model="configForm.enabled" />
        </el-form-item>
        <el-form-item label="返利类型">
          <el-radio-group v-model="configForm.type">
            <el-radio value="fixed">固定金额</el-radio>
            <el-radio value="percent">按比例</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="返利值">
          <el-input-number v-model="configForm.value" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">
            {{ configForm.type === 'percent' ? '%' : '元' }}
          </span>
        </el-form-item>
        <el-form-item label="结算方式">
          <el-select v-model="configForm.settle_type">
            <el-option label="订单完成后结算" value="after_order" />
            <el-option label="注册即结算" value="on_register" />
          </el-select>
        </el-form-item>
        <el-form-item label="返利层级">
          <el-input-number v-model="configForm.level" :min="1" :max="3" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">级</span>
        </el-form-item>
        <el-form-item label="最低提现">
          <el-input-number v-model="configForm.min_payout" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">元</span>
        </el-form-item>
        <el-form-item label="Cookie有效期">
          <el-input-number v-model="configForm.cookie_days" :min="1" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">天</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showConfigDialog = false">取消</el-button>
        <el-button type="primary" @click="saveConfig">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showRuleDialog" title="返利规则" width="700px">
      <div style="margin-bottom: 16px;">
        <el-button type="primary" size="small" @click="addRule">
          <el-icon><Plus /></el-icon>
          添加规则
        </el-button>
      </div>
      <el-table :data="ruleList" style="width: 100%">
        <el-table-column prop="name" label="规则名称" />
        <el-table-column prop="condition" label="触发条件" />
        <el-table-column prop="reward" label="奖励" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row, $index }">
            <el-button type="primary" link @click="editRule(row)">编辑</el-button>
            <el-button type="danger" link @click="ruleList.splice($index, 1)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="showRuleDialog = false">关闭</el-button>
        <el-button type="primary" @click="saveRules">保存规则</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showDetailDialog" title="返利详情" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="推介人">{{ detailData.referrer_name }}</el-descriptions-item>
        <el-descriptions-item label="被推介人">{{ detailData.referred_name }}</el-descriptions-item>
        <el-descriptions-item label="关联订单">{{ detailData.referred_order }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">¥{{ detailData.order_amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="返利金额">
          <span class="amount">¥{{ detailData.commission?.toFixed(2) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTypeMap[detailData.status]" size="small">
            {{ statusTextMap[detailData.status] || detailData.status }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="支付时间">{{ detailData.paid_at || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Setting, List, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const showConfigDialog = ref(false)
const showRuleDialog = ref(false)
const showDetailDialog = ref(false)
const detailData = ref<any>({})
const ruleList = ref<any[]>([])

const statusTypeMap: Record<string, string> = { pending: 'warning', paid: 'success', rejected: 'danger' }
const statusTextMap: Record<string, string> = { pending: '待结算', paid: '已支付', rejected: '已拒绝' }

const searchForm = ref({ username: '', status: '' })

const configForm = ref({
  enabled: true,
  type: 'percent',
  value: 10,
  settle_type: 'after_order',
  level: 1,
  min_payout: 50,
  cookie_days: 30
})

const handleSearch = () => { page.value = 1; fetchData() }

const resetSearch = () => {
  searchForm.value = { username: '', status: '' }
  handleSearch()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/referrals', {
      params: { page: page.value, page_size: pageSize.value, ...searchForm.value }
    })
    tableData.value = data.data || []
    total.value = data.total || 0
  } catch {} finally {
    loading.value = false
  }
}

const fetchConfig = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/referrals/config')
    if (data) configForm.value = { ...configForm.value, ...data }
  } catch {}
}

const fetchRules = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/referrals/rules')
    ruleList.value = data.data || []
  } catch {}
}

const viewDetail = (row: any) => {
  detailData.value = row
  showDetailDialog.value = true
}

const payCommission = async (row: any) => {
  try {
    await request.put(`/admin/api/v1/referrals/${row.id}/pay`)
    ElMessage.success('返利已支付')
    fetchData()
  } catch {}
}

const saveConfig = async () => {
  try {
    await request.put('/admin/api/v1/referrals/config', configForm.value)
    ElMessage.success('配置已保存')
    showConfigDialog.value = false
  } catch {}
}

const addRule = () => {
  ruleList.value.push({ name: '', condition: '', reward: '', status: 'active' })
}

const editRule = (row: any) => {
  // inline editing via table
}

const saveRules = async () => {
  try {
    await request.put('/admin/api/v1/referrals/rules', { rules: ruleList.value })
    ElMessage.success('规则已保存')
    showRuleDialog.value = false
  } catch {}
}

const exportData = () => {
  ElMessage.info('正在导出...')
}

onMounted(() => {
  fetchData()
  fetchConfig()
  fetchRules()
})
</script>

<style scoped lang="scss">
.referral-page {
  .search-bar { margin-bottom: 16px; }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .amount { color: var(--el-color-success); font-weight: 600; }
}
</style>
