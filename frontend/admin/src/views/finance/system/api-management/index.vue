<template>
  <div class="api-management-page">
    <art-card title="API管理" shadow="never">
      <template #header>
        <div class="card-header">
          <span>对外API管理</span>
          <el-button type="primary" @click="showAddDialog">
            <el-icon><Plus /></el-icon>
            添加API
          </el-button>
        </div>
      </template>

      <el-alert title="管理对外API接口的访问密钥，控制第三方系统的接入权限。" type="info" :closable="false" show-icon style="margin-bottom: 16px" />

      <el-table :data="apiList" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="ip" label="IP白名单" />
        <el-table-column prop="create_time" label="创建时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchApiList"
        @current-change="fetchApiList"
      />
    </art-card>

    <!-- 添加API对话框 -->
    <el-dialog v-model="addDialogVisible" title="添加API密钥" width="500px">
      <el-form :model="addForm" :rules="addRules" ref="addFormRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="addForm.username" placeholder="请输入API用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="addForm.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="IP白名单" prop="ip">
          <el-input v-model="addForm.ip" placeholder="请输入允许访问的IP地址，多个用逗号分隔" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAdd" :loading="addLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

interface ApiItem {
  id: number
  username: string
  ip: string
  create_time: string
}

const loading = ref(false)
const addLoading = ref(false)
const addDialogVisible = ref(false)
const addFormRef = ref<FormInstance>()
const apiList = ref<ApiItem[]>([])

const pagination = reactive({ page: 1, limit: 10, total: 0 })
const addForm = reactive({ username: '', password: '', ip: '' })

const addRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const fetchApiList = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/api',
      params: { page: pagination.page, limit: pagination.limit }
    })
    if (res) {
      apiList.value = res.list || []
      pagination.total = res.sum || 0
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const showAddDialog = () => {
  addForm.username = ''
  addForm.password = ''
  addForm.ip = ''
  addDialogVisible.value = true
}

const handleAdd = async () => {
  if (!addFormRef.value) return
  await addFormRef.value.validate(async (valid) => {
    if (!valid) return
    addLoading.value = true
    try {
      await request.post({ url: '/api/admin/api', data: addForm, showSuccessMessage: true })
      addDialogVisible.value = false
      fetchApiList()
    } catch (error) {
      ElMessage.error('添加失败')
    } finally {
      addLoading.value = false
    }
  })
}

const handleDelete = async (row: ApiItem) => {
  try {
    await ElMessageBox.confirm(`确定删除API "${row.username}" 吗？`, '提示')
    await request.delete({ url: '/api/admin/api', data: { id: row.id }, showSuccessMessage: true })
    fetchApiList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchApiList())
</script>

<style scoped lang="scss">
.api-management-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
