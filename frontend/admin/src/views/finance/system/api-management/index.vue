<template>
  <div class="api-management-page">
    <art-card :title="$t('apiManagement.title')" shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('apiManagement.subtitle') }}</span>
          <el-button type="primary" @click="showAddDialog">
            <el-icon><Plus /></el-icon>
            {{ $t('apiManagement.addApi') }}
          </el-button>
        </div>
      </template>

      <el-alert :title="$t('apiManagement.description')" type="info" :closable="false" show-icon style="margin-bottom: 16px" />

      <el-table :data="apiList" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" :label="$t('apiManagement.username')" />
        <el-table-column prop="ip" :label="$t('apiManagement.ipWhitelist')" />
        <el-table-column prop="create_time" :label="$t('apiManagement.createTime')" width="180" />
        <el-table-column :label="$t('apiManagement.operations')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('apiManagement.delete') }}</el-button>
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
    <el-dialog v-model="addDialogVisible" :title="$t('apiManagement.addApiKey')" width="500px">
      <el-form :model="addForm" :rules="addRules" ref="addFormRef" label-width="100px">
        <el-form-item :label="$t('apiManagement.username')" prop="username">
          <el-input v-model="addForm.username" :placeholder="$t('apiManagement.enterUsername')" />
        </el-form-item>
        <el-form-item :label="$t('apiManagement.password')" prop="password">
          <el-input v-model="addForm.password" type="password" show-password :placeholder="$t('apiManagement.enterPassword')" />
        </el-form-item>
        <el-form-item :label="$t('apiManagement.ipWhitelist')" prop="ip">
          <el-input v-model="addForm.ip" :placeholder="$t('apiManagement.enterIpWhitelist')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleAdd" :loading="addLoading">{{ $t('common.confirm') }}</el-button>
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
import { $t } from '@/locales'

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
  username: [{ required: true, message: () => $t('apiManagement.enterUsername'), trigger: 'blur' }],
  password: [{ required: true, message: () => $t('apiManagement.enterPassword'), trigger: 'blur' }]
}

const fetchApiList = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/api-keys',
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
      await request.post({ url: '/api/admin/api-keys', data: addForm, showSuccessMessage: true })
      addDialogVisible.value = false
      fetchApiList()
    } catch (error) {
      ElMessage.error($t('apiManagement.addFailed'))
    } finally {
      addLoading.value = false
    }
  })
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('apiManagement.confirmDelete', { username: row.username }), $t('common.tips'))
    await request.del({ url: `/api/admin/api-keys/${row.id}`, showSuccessMessage: true })
    fetchApiList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error($t('apiManagement.deleteFailed'))
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
