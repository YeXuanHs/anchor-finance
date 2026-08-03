<template>
  <div class="dcim-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>DCIM管理</span>
          <div>
            <el-button @click="handleRefreshAll" :loading="refreshAllLoading">
              <el-icon><Refresh /></el-icon>
              刷新全部状态
            </el-button>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增服务器
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.search" placeholder="服务器名称/主机名" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="服务器名称" min-width="150" />
        <el-table-column prop="hostname" label="主机名" min-width="180" />
        <el-table-column prop="server_num" label="运行中" width="80" align="center" />
        <el-table-column prop="api_status" label="API状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.api_status === 1 ? 'success' : 'danger'" size="small">
              {{ row.api_status === 1 ? '正常' : '异常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="350" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleDetail(row)">详情</el-button>
            <el-button type="success" link @click="handleRefreshStatus(row)">刷新状态</el-button>
            <el-popconfirm
              v-if="row.removable"
              title="确定删除该服务器吗？"
              @confirm="handleDelete(row)"
            >
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
          v-model:page-size="pagination.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="服务器名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入服务器名称" />
        </el-form-item>
        <el-form-item label="主机名/IP" prop="hostname">
          <el-input v-model="formData.hostname" placeholder="请输入主机名或IP地址" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input-number v-model="formData.port" :min="0" :max="65535" />
        </el-form-item>
        <el-form-item label="使用SSL">
          <el-switch v-model="formData.secure" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.disabled" :active-value="0" :inactive-value="1" active-text="启用" inactive-text="禁用" />
        </el-form-item>
        <el-form-item label="用户前缀">
          <el-input v-model="formData.user_prefix" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="服务器详情" width="700px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="服务器名称">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item label="主机名">{{ detailData.hostname }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ detailData.username }}</el-descriptions-item>
        <el-descriptions-item label="端口">{{ detailData.port }}</el-descriptions-item>
        <el-descriptions-item label="SSL">{{ detailData.secure ? '是' : '否' }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ detailData.disabled ? '禁用' : '启用' }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ detailData.bill_type === 'month' ? '月付' : '流量' }}</el-descriptions-item>
        <el-descriptions-item label="重装次数限制">{{ detailData.reinstall_times || '不限' }}</el-descriptions-item>
        <el-descriptions-item label="付费重装">{{ detailData.buy_times ? '启用' : '禁用' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)
const refreshAllLoading = ref(false)

const searchForm = reactive({
  search: ''
})

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0
})

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('新增服务器')
const formRef = ref<FormInstance>()

const detailVisible = ref(false)
const detailData = ref<any>(null)

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  hostname: '',
  username: '',
  password: '',
  port: 0,
  secure: 0,
  disabled: 0,
  user_prefix: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入服务器名称', trigger: 'blur' }],
  hostname: [{ required: true, message: '请输入主机名', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/dcim/servers',
      params: {
        page: pagination.page,
        limit: pagination.limit,
        search: searchForm.search || undefined
      }
    })
    tableData.value = data.list || []
    pagination.total = data.sum || 0
  } catch (error) {
    console.error('获取DCIM服务器列表失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.search = ''
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '新增服务器'
  formData.id = undefined
  formData.name = ''
  formData.hostname = ''
  formData.username = ''
  formData.password = ''
  formData.port = 0
  formData.secure = 0
  formData.disabled = 0
  formData.user_prefix = ''
  dialogVisible.value = true
}

const handleEdit = async (row: any) => {
  dialogTitle.value = '编辑服务器'
  try {
    const data = await request.get({ url: `/api/admin/dcim/servers/${row.id}` })
    Object.assign(formData, data)
    formData.password = ''
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取服务器详情失败')
  }
}

const handleDetail = async (row: any) => {
  try {
    const data = await request.get({ url: `/api/admin/dcim/servers/${row.id}` })
    detailData.value = data
    detailVisible.value = true
  } catch (error) {
    ElMessage.error('获取服务器详情失败')
  }
}

const handleRefreshStatus = async (row: any) => {
  try {
    const data = await request.post({ url: `/api/admin/dcim/servers/${row.id}/refresh` })
    ElMessage.success(data.msg || '刷新成功')
    fetchData()
  } catch (error) {
    ElMessage.error('刷新状态失败')
  }
}

const handleRefreshAll = async () => {
  refreshAllLoading.value = true
  try {
    // 后端无批量刷新接口，逐个刷新
    const list = tableData.value
    for (const row of list) {
      await request.post({ url: `/api/admin/dcim/servers/${row.id}/refresh` })
    }
    ElMessage.success('刷新完成')
    fetchData()
  } catch (error) {
    ElMessage.error('刷新失败')
  } finally {
    refreshAllLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.delete({ url: `/api/admin/dcim/servers/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
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
      if (formData.id) {
        await request.put({ url: `/api/admin/dcim/servers/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/dcim/servers', params: { ...formData } })
      }
      ElMessage.success(formData.id ? '编辑成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchData()
}

const handlePageChange = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.dcim-page {
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

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
