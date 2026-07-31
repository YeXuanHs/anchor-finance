<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <span>用户偏好设置</span>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item label="默认语言">
          <el-select v-model="form.language">
            <el-option label="简体中文" value="zh-CN" />
            <el-option label="English" value="en-US" />
          </el-select>
        </el-form-item>
        <el-form-item label="时区">
          <el-select v-model="form.timezone">
            <el-option label="Asia/Shanghai" value="Asia/Shanghai" />
            <el-option label="UTC" value="UTC" />
          </el-select>
        </el-form-item>
        <el-form-item label="通知方式">
          <el-checkbox-group v-model="form.notifications">
            <el-checkbox label="email">邮件</el-checkbox>
            <el-checkbox label="sms">短信</el-checkbox>
            <el-checkbox label="wechat">微信</el-checkbox>
          </el-checkbox-group>
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
const form = ref({ language: 'zh-CN', timezone: 'Asia/Shanghai', notifications: ['email'] })
const handleSave = async () => {
  await request.put({ url: '/api/admin/user-tastes', data: form.value })
}
</script>
