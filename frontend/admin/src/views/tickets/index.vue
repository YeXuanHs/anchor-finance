<template>
  <div class="tickets-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="工单号">
          <el-input v-model="searchForm.ticket_no" placeholder="工单号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待处理" value="open" />
            <el-option label="处理中" value="processing" />
            <el-option label="已回复" value="replied" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="searchForm.priority" placeholder="全部" clearable>
            <el-option label="高" value="high" />
            <el-option label="中" value="medium" />
            <el-option label="低" value="low" />
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
        <h3>工单列表</h3>
      </div>
      
      <el-table :data="tickets" style="width: 100%" v-loading="loading">
        <el-table-column prop="ticket_no" label="工单号" width="150" />
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="user" label="用户" />
        <el-table-column prop="department" label="部门" />
        <el-table-column prop="priority" label="优先级">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)" size="small">
              {{ row.priority }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <span class="status-tag" :class="row.status">
              {{ getStatusText(row.status) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewTicket(row)">处理</el-button>
            <el-button type="danger" link @click="closeTicket(row)" v-if="row.status !== 'closed'">
              关闭
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const tickets = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const searchForm = ref({
  ticket_no: '',
  status: '',
  priority: ''
})

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    open: '待处理',
    processing: '处理中',
    replied: '已回复',
    closed: '已关闭'
  }
  return map[status] || status
}

const getPriorityType = (priority: string) => {
  const map: Record<string, string> = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return map[priority] || 'info'
}

const handleSearch = () => {
  currentPage.value = 1
}

const resetSearch = () => {
  searchForm.value = { ticket_no: '', status: '', priority: '' }
  handleSearch()
}

const viewTicket = (ticket: any) => {
  router.push(`/tickets/${ticket.ticket_no}`)
}

const closeTicket = async (ticket: any) => {
  // TODO: 关闭工单
}
</script>

<style scoped lang="scss">
.tickets-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
    }
  }
  
  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
