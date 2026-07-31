<template>
  <div class="upstream-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>上游管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加上游
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="名称/地址" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="手动对接" value="manual" />
            <el-option label="V10系统" value="v10" />
            <el-option label="智简魔方" value="zjmf" />
            <el-option label="锚点财务" value="anchor" />
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
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="上游名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="type" label="接口类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type]" size="small">
              {{ typeLabelMap[row.type] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="接口地址" min-width="220" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="test_status" label="连接状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.test_status === 'success'" type="success" size="small">正常</el-tag>
            <el-tag v-else-if="row.test_status === 'failed'" type="danger" size="small">失败</el-tag>
            <el-tag v-else type="info" size="small">未测试</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sync_at" label="最后同步" width="170">
          <template #default="{ row }">
            {{ row.last_sync_at || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleTest(row)" :loading="row._testing">
              测试连接
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该上游接口吗？" @confirm="handleDelete(row)">
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item label="上游名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入上游名称" />
        </el-form-item>
        <el-form-item label="接口类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择接口类型" style="width: 100%" @change="onTypeChange">
            <el-option label="手动对接" value="manual" />
            <el-option label="V10系统" value="v10" />
            <el-option label="智简魔方" value="zjmf" />
            <el-option label="锚点财务" value="anchor" />
          </el-select>
        </el-form-item>
        <el-form-item label="接口地址" prop="url" v-if="formData.type !== 'manual'">
          <el-input v-model="formData.url" placeholder="https://example.com/api" />
        </el-form-item>
        <el-form-item label="API密钥" prop="api_key" v-if="formData.type !== 'manual'">
          <el-input v-model="formData.api_key" placeholder="请输入API密钥" show-password />
        </el-form-item>
        <el-form-item label="API密码" prop="api_password" v-if="formData.type !== 'manual'">
          <el-input v-model="formData.api_password" placeholder="请输入API密码（可选）" show-password />
        </el-form-item>
        <el-form-item label="同步间隔" prop="sync_interval" v-if="formData.type !== 'manual'">
          <el-input-number
            v-model="formData.sync_interval"
            :min="5"
            :max="1440"
            :step="5"
            controls-position="right"
          />
          <span class="form-tip">分钟</span>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注信息（可选）" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
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

const typeLabelMap: Record<string, string> = {
  manual: '手动对接',
  v10: 'V10系统',
  zjmf: '智简魔方',
  anchor: '锚点财务'
}

const typeTagMap: Record<string, string> = {
  manual: 'info',
  v10: 'warning',
  zjmf: '',
  anchor: 'success'
}

const loading = ref(false)
const submitLoading = ref(false)

const searchForm = reactive({
  keyword: '',
  type: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('添加上游')
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  type: '',
  url: '',
  api_key: '',
  api_password: '',
  sync_interval: 30,
  remark: '',
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入上游名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择接口类型', trigger: 'change' }
  ],
  url: [
    { required: true, message: '请输入接口地址', trigger: 'blur' }
  ],
  api_key: [
    { required: true, message: '请输入API密钥', trigger: 'blur' }
  ]
}

const onTypeChange = () => {
  if (formData.type === 'manual') {
    formData.url = ''
    formData.api_key = ''
    formData.api_password = ''
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/upstream',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取上游列表失败:', error)
    ElMessage.error('获取上游列表失败')
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
  searchForm.type = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '添加上游'
  formData.id = undefined
  formData.name = ''
  formData.type = ''
  formData.url = ''
  formData.api_key = ''
  formData.api_password = ''
  formData.sync_interval = 30
  formData.remark = ''
  formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑上游'
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    type: row.type,
    url: row.url || '',
    api_key: row.api_key || '',
    api_password: '',
    sync_interval: row.sync_interval || 30,
    remark: row.remark || '',
    status: row.status
  })
  dialogVisible.value = true
}

const handleTest = async (row: any) => {
  row._testing = true
  try {
    const data = await request.post({
      url: `/api/admin/upstream/${row.id}/test`
    })
    row.test_status = data.status === 'success' ? 'success' : 'failed'
    ElMessage.success(data.message || '测试完成')
  } catch (error) {
    row.test_status = 'failed'
    ElMessage.error('连接测试失败')
  } finally {
    row._testing = false
  }
}

const handleStatusChange = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/upstream/${row.id}`,
      params: { status: row.status }
    })
    ElMessage.success('状态已更新')
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('更新状态失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/upstream/${row.id}`
    })
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
      if (formData.type === 'manual') {
        delete submitData.url
        delete submitData.api_key
        delete submitData.api_password
        delete submitData.sync_interval
      }
      if (!submitData.api_password) {
        delete submitData.api_password
      }

      const url = formData.id ? `/api/admin/upstream/${formData.id}` : '/api/admin/upstream'

      if (formData.id) {
        await request.put({ url, params: submitData })
      } else {
        await request.post({ url, params: submitData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
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
.upstream-page {
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

.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
</style>
