<template>
  <div class="transfer-page">
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>{{ $t('transfer.title') }}</span>
          <el-button type="primary" @click="showDialog = true">{{ $t('transfer.initiateTransfer') }}</el-button>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('helpCommon.sentTransfers')" name="sent">
          <el-table :data="sentList">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="product_name" :label="$t('transfer.productName')" />
            <el-table-column prop="target_user" :label="$t('transfer.targetUser')" />
            <el-table-column prop="status" :label="$t('common.status')" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.operating')" width="150">
              <template #default="{ row }">
                <el-button v-if="row.status === 'pending'" type="danger" link @click="handleCancel(row)">{{ $t('helpCommon.revoke') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane :label="$t('helpCommon.receivedTransfers')" name="received">
          <el-table :data="receivedList">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="product_name" :label="$t('transfer.productName')" />
            <el-table-column prop="from_user" :label="$t('helpCommon.initiator')" />
            <el-table-column prop="status" :label="$t('common.status')" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.operating')" width="200">
              <template #default="{ row }">
                <template v-if="row.status === 'pending'">
                  <el-button type="success" link @click="handleAccept(row)">{{ $t('transfer.acceptTransfer') }}</el-button>
                  <el-button type="danger" link @click="handleReject(row)">{{ $t('transfer.rejectTransfer') }}</el-button>
                </template>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    <el-dialog v-model="showDialog" :title="$t('helpCommon.initiateProductTransfer')" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="$t('helpCommon.selectProduct')" required>
          <el-select v-model="form.product_id" style="width: 100%">
            <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('transfer.targetUser')" required>
          <el-input v-model="form.target_email" :placeholder="$t('transfer.pleaseEnterTargetEmail')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ $t('transfer.confirmTransfer') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/utils/request'

const { t } = useI18n()
const activeTab = ref('sent')
const sentList = ref([])
const receivedList = ref([])
const products = ref<any[]>([])
const showDialog = ref(false)
const form = ref({ product_id: '', target_email: '' })
const getStatusType = (s: string) => ({ pending: 'warning', accepted: 'success', rejected: 'danger', cancelled: 'info' }[s] || '')
const getStatusText = (s: string) => ({ pending: t('transfer.pending'), accepted: t('transfer.accepted'), rejected: t('transfer.rejected'), cancelled: t('helpCommon.cancelled') }[s] || s)
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
