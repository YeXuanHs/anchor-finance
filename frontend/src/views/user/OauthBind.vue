<template>
  <div class="oauth-bind-page">
    <el-card>
      <template #header><span>第三方账号绑定</span></template>
      <div class="bind-list">
        <div v-for="item in providers" :key="item.name" class="bind-item">
          <div class="bind-info">
            <img :src="`/assets/oauth/${item.name}.svg`" class="bind-icon" />
            <span>{{ item.title }}</span>
          </div>
          <div class="bind-action">
            <el-tag v-if="item.bound" type="success">已绑定</el-tag>
            <el-button v-else type="primary" size="small" @click="handleBind(item)">绑定</el-button>
            <el-button v-if="item.bound" type="danger" size="small" link @click="handleUnbind(item)">解绑</el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import request from '@/utils/request'
const providers = ref<any[]>([])
const fetchData = async () => {
  const { data } = await request.get('/api/v1/oauth/providers')
  providers.value = data || []
}
const handleBind = async (item: any) => {
  const { data } = await request.get(`/api/v1/oauth/${item.name}`)
  if (data?.data?.url) window.location.href = data.data.url
}
const handleUnbind = async (item: any) => {
  await request.post('/api/v1/oauth/unbind', { provider: item.name })
  fetchData()
}
fetchData()
</script>
<style scoped>
.bind-list { display: flex; flex-direction: column; gap: 16px; }
.bind-item { display: flex; justify-content: space-between; align-items: center; padding: 12px; border: 1px solid #eee; border-radius: 8px; }
.bind-info { display: flex; align-items: center; gap: 12px; }
.bind-icon { width: 32px; height: 32px; }
</style>
