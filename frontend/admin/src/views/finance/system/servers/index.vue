<template>
  <div class="servers-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>供给配置（服务器模块）</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增模块
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="模块名称/主机地址" clearable />
        </el-form-item>
        <el-form-item label="模块类型">
          <el-select v-model="searchForm.module" placeholder="全部" clearable>
            <el-option label="cPanel" value="cpanel" />
            <el-option label="Plesk" value="plesk" />
            <el-option label="DirectAdmin" value="directadmin" />
            <el-option label="虚拟主机" value="whm" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="模块名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="module" label="模块类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ row.module }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hostname" label="主机地址" min-width="180" show-overflow-tooltip />
        <el-table-column prop="max_accounts" label="最大账户" width="100" align="center" />
        <el-table-column label="当前账户" width="100" align="center">
          <template #default="{ row }">
            <span :class="{ 'warning-text': row.current_count >= row.max_accounts * 0.9 }">
              {{ row.current_count }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="is_enabled" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              :active-value="1"
              :inactive-value="0"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleTest(row)">测试连接</el-button>
            <el-popconfirm title="确定删除该模块吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item label="模块名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入模块名称" />
        </el-form-item>
        <el-form-item label="模块类型" prop="module">
          <el-select v-model="formData.module" placeholder="请选择模块类型" style="width: 100%">
            <el-option label="cPanel" value="cpanel" />
            <el-option label="Plesk" value="plesk" />
            <el-option label="DirectAdmin" value="directadmin" />
            <el-option label="WHM" value="whm" />
          </el-select>
        </el-form-item>
        <el-form-item label="主机地址" prop="hostname">
          <el-input v-model="formData.hostname" placeholder="请输入主机地址" />
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input-number v-model="formData.port" :min="1" :max="65535" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="最大账户" prop="max_accounts">
          <el-input-number v-model="formData.max_accounts" :min="1" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="is_enabled">
          <el-switch v-model="formData.is_enabled" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新增模块')
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  module: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  module: '',
  hostname: '',
  port: 2087,
  username: '',
  password: '',
  max_accounts: 100,
  is_enabled: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入模块名称', trigger: 'blur' }
  ],
  module: [
    { required: true, message: '请选择模块类型', trigger: 'change' }
  ],
  hostname: [
    { required: true, message: '请输入主机地址', trigger: 'blur' }
  ]
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/server-modules',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchForm.keyword || undefined,
        module: searchForm.module || undefined,
        status: searchForm.status
      }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取服务器模块列表失败:', error)
    ElMessage.error('获取服务器模块列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.module = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '新增模块'
  formData.id = undefined
  formData.name = ''
  formData.module = ''
  formData.hostname = ''
  formData.port = 2087
  formData.username = ''
  formData.password = ''
  formData.max_accounts = 100
  formData.is_enabled = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑模块'
  Object.assign(formData, { ...row, password: '' })
  dialogVisible.value = true
}

const handleTest = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/server-modules/${row.id}/test` })
    ElMessage.success('连接测试成功')
  } catch (error) {
    ElMessage.error('连接测试失败')
  }
}

const handleToggleStatus = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/server-modules/${row.id}`,
      params: { is_enabled: row.is_enabled }
    })
    ElMessage.success(row.is_enabled ? '已启用' : '已禁用')
  } catch (error) {
    row.is_enabled = row.is_enabled ? 0 : 1
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/server-modules/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const submitData: any = { ...formData }
      if (formData.id && !formData.password) {
        delete submitData.password
      }

      if (formData.id) {
        await request.put({ url: `/api/admin/server-modules/${formData.id}`, params: submitData })
      } else {
        await request.post({ url: '/api/admin/server-modules', params: submitData })
      }
      ElMessage.success(formData.id ? '更新成功' : '新增成功')
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchList()
}

const handlePageChange = () => {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.servers-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.warning-text {
  color: #e6a23c;
  font-weight: 600;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
