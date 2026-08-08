<template>
  <div class="zjmf-api-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>智简魔方API对接</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增实例
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchForm.name" placeholder="实例名称" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="智简魔方" value="zjmf_api" />
            <el-option label="WHMCS" value="whmcs" />
            <el-option label="V10" value="v10" />
            <el-option label="手动" value="manual" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="连接正常" :value="1" />
            <el-option label="连接失败" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 实例列表 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="实例名称" min-width="150" />
        <el-table-column prop="type_zh" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type_zh || typeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hostname" label="接口地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="product_num" label="产品数" width="80" align="center" />
        <el-table-column prop="set_product_num" label="已配置产品" width="100" align="center" />
        <el-table-column prop="host_num" label="总主机数" width="90" align="center" />
        <el-table-column prop="active_host_num" label="活跃主机" width="90" align="center" />
        <el-table-column label="连接状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '异常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="contact_way" label="联系方式" width="120" show-overflow-tooltip />
        <el-table-column label="操作" width="350" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleTest(row)">测试连接</el-button>
            <el-button type="warning" link @click="handleSyncProducts(row)">同步产品</el-button>
            <el-button type="info" link @click="handleSummary(row)">汇总</el-button>
            <el-button type="primary" link @click="handleViewOrders(row)">订单</el-button>
            <el-popconfirm title="确定删除该实例吗？删除后不可恢复。" @confirm="handleDelete(row)">
              <template #reference><el-button type="danger" link>删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="实例名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入实例名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="智简魔方" value="zjmf_api" />
            <el-option label="WHMCS" value="whmcs" />
            <el-option label="V10" value="v10" />
            <el-option label="手动" value="manual" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="formData.type !== 'manual'" label="接口地址" prop="hostname">
          <el-input v-model="formData.hostname" placeholder="https://example.com" />
        </el-form-item>
        <el-form-item v-if="formData.type !== 'manual'" label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="API用户名" />
        </el-form-item>
        <el-form-item v-if="formData.type !== 'manual'" label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="API密码" />
        </el-form-item>
        <el-form-item label="联系方式">
          <el-input v-model="formData.contact_way" placeholder="联系方式（可选）" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.des" type="textarea" :rows="2" placeholder="描述信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 汇总对话框 -->
    <el-dialog v-model="summaryDialogVisible" title="下游汇总" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="总主机数">{{ summaryData.host_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="活跃主机">{{ summaryData.active_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="产品数">{{ summaryData.agent_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="API余额">
          {{ summaryData.credit?.prefix || '' }}{{ summaryData.credit?.balance || '0.00' }}{{ summaryData.credit?.suffix || '' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="summaryDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const router = useRouter()

const typeMap: Record<string, string> = { zjmf_api: '智简魔方', whmcs: 'WHMCS', v10: 'V10', manual: '手动' }

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const summaryDialogVisible = ref(false)
const dialogTitle = ref('新增实例')
const formRef = ref<FormInstance>()

const searchForm = reactive({ name: '', type: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const summaryData = ref<any>({})

const formData = reactive({
  id: null as number | null,
  name: '',
  type: 'zjmf_api',
  hostname: '',
  username: '',
  password: '',
  contact_way: '',
  des: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入实例名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  hostname: [{ required: true, message: '请输入接口地址', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.name) params.search = searchForm.name
    if (searchForm.type) params.type = searchForm.type
    if (searchForm.status !== undefined) params.status = searchForm.status
    const res = await request.get({ url: '/api/admin/zjmf-api', params })
    tableData.value = res?.list || res?.data?.list || res || []
    pagination.total = res?.total || res?.sum || 0
  } catch { ElMessage.error('获取实例列表失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.name = ''; searchForm.type = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  dialogTitle.value = '新增实例'
  Object.assign(formData, { id: null, name: '', type: 'zjmf_api', hostname: '', username: '', password: '', contact_way: '', des: '' })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑实例'
  Object.assign(formData, { id: row.id, name: row.name, type: row.type, hostname: row.hostname || '', username: row.username || '', password: '', contact_way: row.contact_way || '', des: row.des || '' })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/zjmf-api/${formData.id}`, data: formData, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/zjmf-api', data: formData, showSuccessMessage: true })
      }
      dialogVisible.value = false; fetchData()
    } catch { ElMessage.error('操作失败') } finally { submitLoading.value = false }
  })
}

const handleTest = async (row: any) => {
  try {
    const res = await request.post({ url: `/api/admin/zjmf-api/${row.id}/test` })
    ElMessage.success(res?.desc || res?.message || '连接测试成功')
    fetchData()
  } catch { ElMessage.error('连接测试失败') }
}

const handleSyncProducts = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定同步实例 "${row.name}" 的产品数据吗？`, '同步确认')
    const res = await request.post({ url: `/api/admin/zjmf-api/${row.id}/sync` })
    ElMessage.success(res?.message || '同步任务已提交')
    fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('同步失败') }
}

const handleSummary = async (row: any) => {
  try {
    const data = await request.get({ url: `/api/admin/zjmf-api/${row.id}/summary` })
    summaryData.value = data || {}
    summaryDialogVisible.value = true
  } catch { ElMessage.error('获取汇总失败') }
}

const handleViewOrders = (row: any) => {
  router.push({ path: '/order-list', query: { zjmf_api_id: row.id } })
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/zjmf-api/${row.id}` })
    ElMessage.success('删除成功'); fetchData()
  } catch { ElMessage.error('删除失败') }
}

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.zjmf-api-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 16px; .el-form-item { margin-bottom: 0; } }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
