<template>
  <div class="crm-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户CRM</span>
          <el-button type="primary" @click="handleAddRecord">
            <el-icon><Plus /></el-icon>
            添加记录
          </el-button>
        </div>
      </template>

      <!-- 客户信息卡片 -->
      <el-card shadow="hover" class="client-info-card">
        <el-descriptions :column="4" border>
          <el-descriptions-item label="客户名称">{{ clientInfo.username }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ clientInfo.email }}</el-descriptions-item>
          <el-descriptions-item label="手机">{{ clientInfo.phone }}</el-descriptions-item>
          <el-descriptions-item label="等级">{{ clientInfo.level_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="累计消费">¥{{ formatAmount(clientInfo.total_spent) }}</el-descriptions-item>
          <el-descriptions-item label="订单数">{{ clientInfo.order_count || 0 }}</el-descriptions-item>
          <el-descriptions-item label="最后登录">{{ clientInfo.last_login || '-' }}</el-descriptions-item>
          <el-descriptions-item label="客户来源">{{ clientInfo.source || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- Tab切换 -->
      <el-tabs v-model="activeTab" class="crm-tabs">
        <el-tab-pane label="跟进记录" name="follow">
          <!-- 搜索栏 -->
          <el-form :inline="true" :model="followSearch" class="search-form">
            <el-form-item label="类型">
              <el-select v-model="followSearch.type" placeholder="全部" clearable>
                <el-option label="电话" value="phone" />
                <el-option label="邮件" value="email" />
                <el-option label="拜访" value="visit" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchFollowRecords">搜索</el-button>
              <el-button @click="followSearch.type = ''; fetchFollowRecords()">重置</el-button>
            </el-form-item>
          </el-form>

          <!-- 时间线 -->
          <el-timeline>
            <el-timeline-item
              v-for="record in followRecords"
              :key="record.id"
              :timestamp="record.created_at"
              placement="top"
              :type="getTimelineType(record.type)"
            >
              <el-card shadow="hover" class="timeline-card">
                <div class="timeline-header">
                  <el-tag :type="getTypeTag(record.type)" size="small">{{ getTypeText(record.type) }}</el-tag>
                  <span class="operator">{{ record.operator }}</span>
                </div>
                <div class="timeline-content">{{ record.content }}</div>
                <div class="timeline-next" v-if="record.next_follow_at">
                  <el-icon><Clock /></el-icon>
                  下次跟进: {{ record.next_follow_at }}
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>

          <el-empty v-if="followRecords.length === 0" description="暂无跟进记录" />
        </el-tab-pane>

        <el-tab-pane label="商机记录" name="opportunity">
          <el-table :data="opportunityData" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="title" label="商机名称" min-width="150" />
            <el-table-column prop="amount" label="预计金额" width="120" align="right">
              <template #default="{ row }">¥{{ formatAmount(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="stage" label="阶段" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="getStageTag(row.stage)" size="small">{{ getStageText(row.stage) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="probability" label="成交概率" width="100" align="center">
              <template #default="{ row }">{{ row.probability }}%</template>
            </el-table-column>
            <el-table-column prop="expected_close_date" label="预计成交" width="120" />
            <el-table-column prop="created_at" label="创建时间" width="170" />
            <el-table-column label="操作" width="100" fixed="right" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEditOpportunity(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="合同记录" name="contract">
          <el-table :data="contractData" style="width: 100%" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="contract_no" label="合同编号" width="170" />
            <el-table-column prop="title" label="合同名称" min-width="150" />
            <el-table-column prop="amount" label="合同金额" width="120" align="right">
              <template #default="{ row }">¥{{ formatAmount(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getContractStatusTag(row.status)" size="small">{{ getContractStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="start_date" label="开始日期" width="120" />
            <el-table-column prop="end_date" label="结束日期" width="120" />
            <el-table-column label="操作" width="100" fixed="right" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleViewContract(row)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 添加跟进记录对话框 -->
    <el-dialog v-model="recordVisible" title="添加跟进记录" width="600px">
      <el-form :model="recordForm" :rules="recordRules" ref="recordFormRef" label-width="100px">
        <el-form-item label="跟进方式" prop="type">
          <el-select v-model="recordForm.type" placeholder="请选择">
            <el-option label="电话" value="phone" />
            <el-option label="邮件" value="email" />
            <el-option label="拜访" value="visit" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="跟进内容" prop="content">
          <el-input v-model="recordForm.content" type="textarea" :rows="4" placeholder="请输入跟进内容" />
        </el-form-item>
        <el-form-item label="下次跟进">
          <el-date-picker v-model="recordForm.next_follow_at" type="datetime" placeholder="选择日期时间" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="recordVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRecordSubmit" :loading="recordLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Plus, Clock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const clientId = route.params.id as string

const activeTab = ref('follow')
const recordLoading = ref(false)

const clientInfo = ref<any>({})
const followRecords = ref<any[]>([])
const opportunityData = ref<any[]>([])
const contractData = ref<any[]>([])

const followSearch = reactive({ type: '' })

const recordVisible = ref(false)
const recordFormRef = ref<FormInstance>()
const recordForm = reactive({
  type: 'phone',
  content: '',
  next_follow_at: ''
})

const recordRules: FormRules = {
  type: [{ required: true, message: '请选择跟进方式', trigger: 'change' }],
  content: [{ required: true, message: '请输入跟进内容', trigger: 'blur' }]
}

const formatAmount = (amount: number | undefined) =>
  amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'

const getTypeText = (type: string) => {
  const map: Record<string, string> = { phone: '电话', email: '邮件', visit: '拜访', other: '其他' }
  return map[type] || type
}

const getTypeTag = (type: string) => {
  const map: Record<string, string> = { phone: 'primary', email: 'success', visit: 'warning', other: 'info' }
  return (map[type] || 'info') as any
}

const getTimelineType = (type: string) => {
  const map: Record<string, string> = { phone: 'primary', email: 'success', visit: 'warning', other: '' }
  return map[type] || '' as any
}

const getStageText = (stage: string) => {
  const map: Record<string, string> = { initial: '初步接触', negotiation: '商务谈判', proposal: '方案报价', closed_won: '成交', closed_lost: '丢单' }
  return map[stage] || stage
}

const getStageTag = (stage: string) => {
  const map: Record<string, string> = { initial: 'info', negotiation: 'warning', proposal: 'primary', closed_won: 'success', closed_lost: 'danger' }
  return (map[stage] || 'info') as any
}

const getContractStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '草稿', 1: '生效中', 2: '已到期', 3: '已终止' }
  return map[status] || '未知'
}

const getContractStatusTag = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'success', 2: 'warning', 3: 'danger' }
  return (map[status] || 'info') as any
}

const fetchClientInfo = async () => {
  try {
    const data = await request.get({ url: `/api/admin/users/${clientId}` })
    clientInfo.value = data
  } catch (error) {
    console.error('获取客户信息失败:', error)
  }
}

const fetchFollowRecords = async () => {
  try {
    const params: any = { client_id: clientId }
    if (followSearch.type) params.type = followSearch.type
    const data = await request.get({ url: '/api/admin/crm/follow-records', params })
    followRecords.value = data.list || data || []
  } catch (error) {
    console.error('获取跟进记录失败:', error)
  }
}

const fetchOpportunities = async () => {
  try {
    const data = await request.get({ url: '/api/admin/crm/opportunities', params: { client_id: clientId } })
    opportunityData.value = data.list || data || []
  } catch (error) {
    console.error('获取商机记录失败:', error)
  }
}

const fetchContracts = async () => {
  try {
    const data = await request.get({ url: '/api/admin/crm/contracts', params: { client_id: clientId } })
    contractData.value = data.list || data || []
  } catch (error) {
    console.error('获取合同记录失败:', error)
  }
}

const handleAddRecord = () => {
  recordForm.type = 'phone'
  recordForm.content = ''
  recordForm.next_follow_at = ''
  recordVisible.value = true
}

const handleRecordSubmit = async () => {
  if (!recordFormRef.value) return
  await recordFormRef.value.validate(async (valid) => {
    if (!valid) return
    recordLoading.value = true
    try {
      await request.post({ url: '/api/admin/crm/follow-records', params: { ...recordForm, client_id: clientId } })
      ElMessage.success('添加成功')
      recordVisible.value = false
      fetchFollowRecords()
    } catch (error) {
      ElMessage.error('添加失败')
    } finally {
      recordLoading.value = false
    }
  })
}

const handleEditOpportunity = (row: any) => {
  ElMessage.info('编辑商机功能开发中...')
}

const handleViewContract = (row: any) => {
  ElMessage.info('查看合同功能开发中...')
}

onMounted(() => {
  fetchClientInfo()
  fetchFollowRecords()
  fetchOpportunities()
  fetchContracts()
})
</script>

<style scoped lang="scss">
.crm-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.client-info-card {
  margin-bottom: 20px;
}

.crm-tabs {
  margin-top: 20px;
}

.search-form {
  margin-bottom: 20px;
}

.timeline-card {
  .timeline-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;

    .operator {
      color: #909399;
      font-size: 13px;
    }
  }

  .timeline-content {
    color: #303133;
    line-height: 1.6;
  }

  .timeline-next {
    margin-top: 8px;
    color: #e6a23c;
    font-size: 13px;
    display: flex;
    align-items: center;
    gap: 4px;
  }
}
</style>
