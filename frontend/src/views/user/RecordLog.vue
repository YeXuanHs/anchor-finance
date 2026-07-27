<template>
  <div class="record-log-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
          <div class="header-actions">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 300px"
            />
            <el-input
              v-model="searchKeyword"
              placeholder="搜索日志..."
              style="width: 200px"
            />
          </div>
        </div>
      </template>

      <el-table :data="logs" style="width: 100%" v-loading="loading">
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column prop="ip" label="IP地址" width="150" />
        <el-table-column prop="description" label="操作描述" />
        <el-table-column prop="user" label="操作人" width="120" />
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const loading = ref(false)
const logs = ref([])
const dateRange = ref([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
</script>

<style scoped lang="scss">
.record-log-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }
}
</style>
