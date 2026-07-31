<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>维护模式</span>
          <el-switch v-model="enabled" active-text="开启" inactive-text="关闭" @change="handleToggle" />
        </div>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item label="维护标题">
          <el-input v-model="form.title" placeholder="系统维护中" />
        </el-form-item>
        <el-form-item label="维护说明">
          <el-input v-model="form.message" type="textarea" :rows="4" placeholder="系统正在进行维护，请稍后再试..." />
        </el-form-item>
        <el-form-item label="允许IP">
          <div v-for="(ip, index) in form.allowed_ips" :key="index" style="display: flex; gap: 10px; margin-bottom: 10px">
            <el-input v-model="form.allowed_ips[index]" placeholder="IP地址" />
            <el-button type="danger" @click="form.allowed_ips.splice(index, 1)">删除</el-button>
          </div>
          <el-button @click="form.allowed_ips.push('')">添加IP</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave">保存</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import request from '@/utils/http'
const enabled = ref(false)
const form = ref({ title: '', message: '', allowed_ips: [] as string[] })
const handleToggle = async (val: boolean) => {
  await request.post({ url: `/api/admin/maintenance/${val ? 'enable' : 'disable'}` })
}
const handleSave = async () => {
  await request.put({ url: '/api/admin/maintenance/settings', data: form.value })
}
</script>
