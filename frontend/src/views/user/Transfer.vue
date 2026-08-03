<template>
  <div class="transfer-page">
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>产品转移</span>
          <el-button type="primary" @click="showDialog = true">发起转移</el-button>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="我发起的" name="sent">
          <el-table :data="sentList">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="product_name" label="产品名称" />
            <el-table-column prop="target_user" label="目标用户" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button v-if="row.status === 'pending'" type="danger" link @click="handleCancel(row)">撤销</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="我收到的" name="received">
          <el-table :data="receivedList">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="product_name" label="产品名称" />
            <el-table-column prop="from_user" label="发起用户" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <template v-if="row.status === 'pending'">
                  <el-button type="success" link @click="handleAccept(row)">接受</el-button>
                  <el-button type="danger" link @click="handleReject(row)">拒绝</el-button>
                </template>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    <el-dialog v-model="showDialog" title="发起产品转移" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="选择产品" required>
          <el-select v-model="form.product_id" style="width: 100%">
            <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标用户" required>
          <el-input v-model="form.target_email" placeholder="输入目标用户邮箱" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确认转移</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import request from '@/utils/request'
const activeTab = ref('sent')
const sentList = ref([])
const receivedList = ref([])
const products = ref([])
const showDialog = ref(false)
const form = ref({ product_id: '', target_email: '' })
const getStatusType = (s: string) => ({ pending: 'warning', accepted: 'success', rejected: 'danger', cancelled: 'info' }[s] || '')
const getStatusText = (s: string) => ({ pending: '待处理', accepted: '已接受', rejected: '已拒绝', cancelled: '已撤销' }[s] || s)
const fetchData = async () => {
  const [s, r] = await Promise.all([
    request.get('/api/v1/product-diverts'),
    request.get('/api/v1/product-diverts/received')
  ])
  sentList.value = s.data || []
  receivedList.value = r.data || []
}
const handleCancel = async (row: any) => { await request.post(`/api/v1/product-diverts/${row.id}/cancel`); fetchData() }
const handleAccept = async (row: any) => { await request.post(`/api/v1/product-diverts/${row.id}/accept`); fetchData() }
const handleReject = async (row: any) => { await request.post(`/api/v1/product-diverts/${row.id}/reject`); fetchData() }
const handleSubmit = async () => { await request.post('/api/v1/product-diverts', form.value); showDialog.value = false; fetchData() }
fetchData()
</script>
