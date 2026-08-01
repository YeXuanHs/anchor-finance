<template>
  <div class="ticket-statistics-page">
    <art-card title="工单统计" shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单处理统计</span>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="fetchStatistics"
            style="width: 300px"
          />
        </div>
      </template>

      <!-- 概览卡片 -->
      <el-row :gutter="16" class="stat-cards">
        <el-col :span="6">
          <el-card shadow="hover">
            <el-statistic title="总工单数" :value="overview.total" />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <el-statistic title="本月工单" :value="overview.this_month" />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <el-statistic title="本年工单" :value="overview.this_year" />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <el-statistic title="平均评分" :value="overview.avg_score" :precision="1" />
          </el-card>
        </el-col>
      </el-row>

      <!-- 详细统计表格 -->
      <el-table :data="statisticsList" v-loading="loading" stripe border style="margin-top: 20px">
        <el-table-column prop="handler" label="处理人" />
        <el-table-column prop="ticket_count" label="处理工单数" width="120" sortable />
        <el-table-column prop="total_score" label="总分" width="100" sortable />
        <el-table-column label="评分分布">
          <el-table-column prop="score_1" label="1分" width="70" />
          <el-table-column prop="score_2" label="2分" width="70" />
          <el-table-column prop="score_3" label="3分" width="70" />
          <el-table-column prop="score_4" label="4分" width="70" />
          <el-table-column prop="score_5" label="5分" width="70" />
        </el-table-column>
      </el-table>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/http'

const loading = ref(false)
const dateRange = ref<[Date, Date] | null>(null)

const overview = reactive({
  total: 0,
  this_month: 0,
  this_year: 0,
  avg_score: 0
})

const statisticsList = ref<Array<{
  handler: string
  ticket_count: number
  total_score: number
  score_1: number
  score_2: number
  score_3: number
  score_4: number
  score_5: number
}>>([])

const fetchStatistics = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (dateRange.value) {
      params.start_time = dateRange.value[0].toISOString().split('T')[0]
      params.end_time = dateRange.value[1].toISOString().split('T')[0]
    }
    const res = await request.get({ url: '/api/admin/tickets/statistics', params })
    if (res?.data) {
      statisticsList.value = res.data.list || []
      Object.assign(overview, {
        total: res.data.sum || 0,
        this_month: res.data.this_month || 0,
        this_year: res.data.this_year || 0,
        avg_score: res.data.avg_score || 0
      })
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchStatistics())
</script>

<style scoped lang="scss">
.ticket-statistics-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.stat-cards {
  margin-bottom: 16px;
}
</style>
