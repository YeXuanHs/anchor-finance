<template>
  <div class="get-user-page page-container">
    <div class="page-header">
      <h2>获取用户</h2>
    </div>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户查询</span>
        </div>
      </template>
      <el-form :model="queryForm" label-width="100px">
        <el-form-item label="用户ID">
          <el-input v-model="queryForm.user_id" placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="queryForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="queryForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 20px;" v-if="userInfo">
      <template #header>
        <div class="card-header">
          <span>用户信息</span>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户ID">{{ userInfo.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ userInfo.email }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ userInfo.phone }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="userInfo.status === 1 ? 'success' : 'danger'">
            {{ userInfo.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ userInfo.created_at }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const queryForm = reactive({
  user_id: '',
  username: '',
  email: ''
})

const userInfo = ref(null)

const handleQuery = async () => {
  if (!queryForm.user_id && !queryForm.username && !queryForm.email) {
    ElMessage.warning('请输入至少一个查询条件')
    return
  }

  try {
    const { data } = await request.get('/admin/api/v1/users/search', { params: queryForm })
    userInfo.value = data.data || null
  } catch (error) {
    ElMessage.error('查询失败')
    userInfo.value = null
  }
}

const handleReset = () => {
  queryForm.user_id = ''
  queryForm.username = ''
  queryForm.email = ''
  userInfo.value = null
}
</script>

<style scoped>
.get-user-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
