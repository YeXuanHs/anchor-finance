<template>
  <div class="credit-record">
    <!-- 搜索筛选 -->
    <div class="filter-bar">
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        style="width: 300px"
      />
      <el-select v-model="filterType" placeholder="类型" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="信用增加" value="increase" />
        <el-option label="信用扣减" value="decrease" />
        <el-option label="信用恢复" value="restore" />
      </el-select>
      <el-select v-model="filterSource" placeholder="来源" clearable style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="按时付款" value="payment" />
        <el-option label="订单取消" value="cancel" />
        <el-option label="系统调整" value="system" />
        <el-option label="违规处罚" value="violation" />
      </el-select>
      <el-button type="primary" @click="handleSearch">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 信用概览 -->
    <div class="credit-overview">
      <el-card class="credit-score-card" shadow="hover">
        <div class="score-content">
          <div class="score-circle">
            <el-progress
              type="circle"
              :percentage="creditScore"
              :width="120"
              :stroke-width="10"
              :color="scoreColor"
            >
              <template #default="{ percentage }">
                <div class="score-text">
                  <span class="score-number">{{ percentage }}</span>
                  <span class="score-label">信用分</span>
                </div>
              </template>
            </el-progress>
          </div>
          <div class="score-info">
            <div class="score-level">
              信用等级：<el-tag :type="levelTagType" size="large">{{ creditLevel }}</el-tag>
            </div>
            <div class="score-desc">{{ creditDesc }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 数据表格 -->
    <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
      <el-table-column prop="id" label="记录ID" width="120" show-overflow-tooltip />
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="typeTagType(row.type)" size="small">
            {{ typeLabel(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="change" label="变动分数" width="120">
        <template #default="{ row }">
          <span :class="['change-text', row.change > 0 ? 'positive' : 'negative']">
            {{ row.change > 0 ? '+' : '' }}{{ row.change }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="before_score" label="变动前" width="100" />
      <el-table-column prop="after_score" label="变动后" width="100" />
      <el-table-column prop="source" label="来源" width="120">
        <template #default="{ row }">
          <span>{{ sourceLabel(row.source) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="180" show-overflow-tooltip />
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSearch"
      @current-change="handleSearch"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface CreditItem {
  id: string
  created_at: string
  type: string
  change: number
  before_score: number
  after_score: number
  source: string
  remark: string
}

const loading = ref(false)
const dateRange = ref<string[]>([])
const filterType = ref('')
const filterSource = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const creditScore = ref(85)
const creditLevel = computed(() => {
  if (creditScore.value >= 90) return '优秀'
  if (creditScore.value >= 70) return '良好'
  if (creditScore.value >= 60) return '一般'
  return '较差'
})

const levelTagType = computed(() => {
  if (creditScore.value >= 90) return 'success'
  if (creditScore.value >= 70) return 'primary'
  if (creditScore.value >= 60) return 'warning'
  return 'danger'
})

const creditDesc = computed(() => {
  if (creditScore.value >= 90) return '您的信用极好，享受优先处理和专属优惠'
  if (creditScore.value >= 70) return '您的信用良好，大部分服务正常享受'
  if (creditScore.value >= 60) return '您的信用一般，建议保持良好的交易记录'
  return '您的信用较差，部分服务可能受限'
})

const scoreColor = computed(() => {
  if (creditScore.value >= 90) return '#52c41a'
  if (creditScore.value >= 70) return '#1890ff'
  if (creditScore.value >= 60) return '#faad14'
  return '#ff4d4f'
})

const tableData = ref<CreditItem[]>([
  { id: 'CR001', created_at: '2026-07-27 10:00:00', type: 'increase', change: 5, before_score: 80, after_score: 85, source: 'payment', remark: '订单按时付款奖励' },
  { id: 'CR002', created_at: '2026-07-25 15:30:00', type: 'decrease', change: -10, before_score: 90, after_score: 80, source: 'cancel', remark: '频繁取消订单' },
  { id: 'CR003', created_at: '2026-07-20 09:00:00', type: 'increase', change: 3, before_score: 87, after_score: 90, source: 'payment', remark: '连续按时付款奖励' },
  { id: 'CR004', created_at: '2026-07-15 14:20:00', type: 'restore', change: 5, before_score: 82, after_score: 87, source: 'system', remark: '系统信用恢复' },
  { id: 'CR005', created_at: '2026-07-10 11:45:00', type: 'decrease', change: -20, before_score: 100, after_score: 80, source: 'violation', remark: '违规操作处罚' },
])

const typeLabel = (type: string) => {
  const map: Record<string, string> = { increase: '信用增加', decrease: '信用扣减', restore: '信用恢复' }
  return map[type] || type
}

const typeTagType = (type: string) => {
  const map: Record<string, string> = { increase: 'success', decrease: 'danger', restore: 'warning' }
  return map[type] || 'info'
}

const sourceLabel = (source: string) => {
  const map: Record<string, string> = { payment: '按时付款', cancel: '订单取消', system: '系统调整', violation: '违规处罚' }
  return map[source] || source
}

const handleSearch = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/credit/logs', { params: { page: currentPage.value, page_size: pageSize.value } })
    tableData.value = res.data.data.list
    total.value = res.data.data.total
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  dateRange.value = []
  filterType.value = ''
  filterSource.value = ''
  currentPage.value = 1
  handleSearch()
}

defineExpose({ handleSearch })
</script>

<style scoped lang="scss">
.credit-record {
  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }

  .credit-overview {
    margin-bottom: 20px;
  }

  .credit-score-card {
    border-radius: 12px;
    border: none;
  }

  .score-content {
    display: flex;
    align-items: center;
    gap: 40px;
  }

  .score-circle {
    flex-shrink: 0;
  }

  .score-text {
    display: flex;
    flex-direction: column;
    align-items: center;

    .score-number {
      font-size: 32px;
      font-weight: 700;
      color: #262626;
      line-height: 1;
    }

    .score-label {
      font-size: 12px;
      color: #8c8c8c;
      margin-top: 4px;
    }
  }

  .score-info {
    .score-level {
      font-size: 16px;
      margin-bottom: 12px;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .score-desc {
      font-size: 14px;
      color: #595959;
      line-height: 1.6;
    }
  }

  .change-text {
    font-weight: 600;

    &.positive {
      color: #52c41a;
    }

    &.negative {
      color: #ff4d4f;
    }
  }

  .el-pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }
}
</style>
