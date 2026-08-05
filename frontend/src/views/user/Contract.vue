<template>
  <div class="contract-page">
    <el-card>
      <template #header><span>{{ $t('contractPage.title') }}</span></template>
      <el-table :data="contracts" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('contractPage.contractTitle')" />
        <el-table-column prop="status" :label="$t('contractPage.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'signed' ? 'success' : 'warning'">{{ row.status === 'signed' ? $t('contractPage.signed') : $t('contractPage.pendingSign') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('contractPage.createTime')" width="180" />
        <el-table-column :label="$t('contractPage.operating')" width="200">
          <template #default="{ row }">
            <el-button type="primary" link>{{ $t('contractPage.view') }}</el-button>
            <el-button v-if="row.status !== 'signed'" type="success" link>{{ $t('contractPage.sign') }}</el-button>
            <el-button type="info" link>{{ $t('contractPage.downloadPdf') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import request from '@/utils/request'
const loading = ref(false)
const contracts = ref([])
const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/contracts')
    contracts.value = data || []
  } finally { loading.value = false }
}
fetchData()
</script>
