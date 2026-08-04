<template>
  <div class="email-suffixes-page">
    <div class="page-header">
      <div class="header-left">
        <h2>邮箱后缀白名单</h2>
        <span class="subtitle">限制用户注册时可使用的邮箱后缀</span>
      </div>
      <div class="header-actions">
        <el-button @click="handleImportDefaults" :loading="importing">导入默认后缀</el-button>
        <el-button type="primary" @click="showAddDialog">添加后缀</el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filter-bar">
      <el-input v-model="searchKeyword" placeholder="搜索后缀..." clearable style="width: 300px" @clear="loadData">
        <template #append>
          <el-button @click="loadData">搜索</el-button>
        </template>
      </el-input>
      <el-checkbox v-model="showInactive" @change="loadData">显示已禁用</el-checkbox>
    </div>

    <!-- 数据表格 -->
    <el-table :data="suffixList" v-loading="loading" stripe>
      <el-table-column prop="suffix" label="邮箱后缀" width="200">
        <template #default="{ row }">
          <span class="suffix-text">@{{ row.suffix }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="200" />
      <el-table-column prop="is_default" label="默认" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="is_active" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-switch v-model="row.is_active" @change="handleToggleActive(row)" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" align="center">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="showEditDialog(row)">编辑</el-button>
          <el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)">
            <template #reference>
              <el-button type="danger" link size="small">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑后缀' : '添加后缀'" width="500px">
      <el-form :model="formData" label-width="80px">
        <el-form-item label="后缀" required>
          <el-input v-model="formData.suffix" placeholder="例如: gmail.com" :disabled="isEdit">
            <template #prepend>@</template>
          </el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.remark" placeholder="例如: Google 邮箱" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

interface EmailSuffix {
  id: number
  suffix: string
  is_default: boolean
  is_active: boolean
  remark: string
}

const loading = ref(false)
const importing = ref(false)
const submitting = ref(false)
const suffixList = ref<EmailSuffix[]>([])
const searchKeyword = ref('')
const showInactive = ref(true)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const formData = ref({
  suffix: '',
  remark: ''
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/email-suffixes', params: { show_inactive: showInactive.value } })
    suffixList.value = res || []
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const showAddDialog = () => {
  isEdit.value = false
  formData.value = { suffix: '', remark: '' }
  dialogVisible.value = true
}

const showEditDialog = (row: EmailSuffix) => {
  isEdit.value = true
  editId.value = row.id
  formData.value = { suffix: row.suffix, remark: row.remark }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formData.value.suffix) {
    ElMessage.warning('请输入邮箱后缀')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await request.put(`/api/admin/email-suffixes/${editId.value}`, {
        remark: formData.value.remark
      })
    } else {
      await request.post({ url: '/api/admin/email-suffixes', params: formData.value })
    }
    ElMessage.success(isEdit.value ? '更新成功' : '添加成功')
    dialogVisible.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleToggleActive = async (row: EmailSuffix) => {
  try {
    await request.put(`/api/admin/email-suffixes/${row.id}`, {
      is_active: row.is_active
    })
  } catch (e) {
    row.is_active = !row.is_active
    ElMessage.error('更新失败')
  }
}

const handleDelete = async (id: number) => {
  try {
    await request.del({ url: `/api/admin/email-suffixes/${id}` })
    ElMessage.success('删除成功')
    loadData()
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

const handleImportDefaults = async () => {
  importing.value = true
  try {
    await request.post({ url: '/api/admin/email-suffixes/import-defaults' })
    ElMessage.success('默认后缀导入完成')
    loadData()
  } catch (e) {
    ElMessage.error('导入失败')
  } finally {
    importing.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.email-suffixes-page {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.header-left h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
}
.subtitle {
  color: #909399;
  font-size: 14px;
}
.filter-bar {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 16px;
}
.suffix-text {
  font-family: monospace;
  font-weight: 600;
  color: #409eff;
}
</style>
