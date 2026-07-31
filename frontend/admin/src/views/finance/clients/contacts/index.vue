<template>
  <div class="page-container">
    <art-table :data="tableData" :loading="loading">
      <template #header>
        <el-input v-model="searchKey" placeholder="搜索联系人" style="width: 200px" />
      </template>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="user_id" label="用户ID" width="100" />
      <el-table-column prop="name" label="姓名" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column prop="phone" label="电话" width="140" />
      <el-table-column prop="is_default" label="默认" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.is_default" type="success" size="small">是</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link>编辑</el-button>
          <el-button type="danger" link>删除</el-button>
        </template>
      </el-table-column>
    </art-table>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import request from '@/utils/http'
const loading = ref(false)
const tableData = ref([])
const searchKey = ref('')
const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get({ url: '/api/admin/contacts' })
    tableData.value = data || []
  } finally { loading.value = false }
}
fetchData()
</script>
