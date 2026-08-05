<template>
  <div class="provisions-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>模块供给管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            添加供给
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="module" label="模块" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.module }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="host" label="主机" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '异常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sync" label="最后同步" width="180" />
        <el-table-column label="操作" width="350" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleTest(row)">测试连接</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleViewButtons(row)">按钮管理</el-button>
            <el-button size="small" @click="handleViewMeta(row)">模块信息</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="模块" prop="module">
          <el-select v-model="formData.module" placeholder="请选择模块" style="width: 100%" @change="handleModuleChange">
            <el-option v-for="mod in availableModules" :key="mod.name" :label="mod.name" :value="mod.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="主机地址" prop="host">
          <el-input v-model="formData.host" placeholder="请输入主机地址" />
        </el-form-item>
        <el-form-item label="端口">
          <el-input-number v-model="formData.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="formData.username" placeholder="API用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="formData.password" type="password" show-password placeholder="API密码" />
        </el-form-item>
        <!-- 动态配置项 -->
        <template v-if="moduleConfigOptions.length">
          <el-divider>模块配置</el-divider>
          <el-form-item v-for="opt in moduleConfigOptions" :key="opt.name" :label="opt.description || opt.name">
            <el-select v-if="opt.type === 'dropdown'" v-model="formData.config[opt.name]" style="width: 100%">
              <el-option v-for="o in opt.options" :key="o.value" :label="o.name" :value="o.value" />
            </el-select>
            <el-input v-else-if="opt.type === 'text' || opt.type === 'password'" v-model="formData.config[opt.name]" :type="opt.type === 'password' ? 'password' : 'text'" show-password />
            <el-switch v-else-if="opt.type === 'yesno'" v-model="formData.config[opt.name]" active-value="on" inactive-value="" />
            <el-input-number v-else-if="opt.type === 'number'" v-model="formData.config[opt.name]" :min="0" />
            <el-input v-else v-model="formData.config[opt.name]" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 按钮管理对话框 -->
    <el-dialog v-model="buttonDialogVisible" title="按钮管理" width="700px">
      <div class="button-header">
        <el-button type="primary" size="small" @click="showButtonForm()">添加按钮</el-button>
      </div>
      <el-table :data="buttonList" v-loading="buttonLoading" border size="small">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="按钮名称" width="150" />
        <el-table-column prop="func" label="函数名" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="200" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="showButtonForm(row)">编辑</el-button>
            <el-button type="success" link @click="handleExecuteButton(row)">执行</el-button>
            <el-popconfirm title="确定删除该按钮吗？" @confirm="handleDeleteButton(row)">
              <template #reference><el-button type="danger" link>删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 按钮表单对话框 -->
    <el-dialog v-model="buttonFormVisible" :title="buttonForm.id ? '编辑按钮' : '添加按钮'" width="500px">
      <el-form :model="buttonForm" :rules="buttonRules" ref="buttonFormRef" label-width="100px">
        <el-form-item label="按钮名称" prop="name">
          <el-input v-model="buttonForm.name" placeholder="如：开机、关机、重启" />
        </el-form-item>
        <el-form-item label="函数名" prop="func">
          <el-input v-model="buttonForm.func" placeholder="如：on、off、reboot" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="buttonForm.description" type="textarea" :rows="2" placeholder="按钮功能描述" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="buttonForm.order" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="buttonForm.is_enabled" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="buttonFormVisible = false">取消</el-button>
        <el-button type="primary" @click="handleButtonSubmit" :loading="buttonSubmitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 模块信息对话框 -->
    <el-dialog v-model="metaDialogVisible" title="模块信息" width="600px">
      <div v-loading="metaLoading">
        <template v-if="metaData">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="模块名称">{{ metaData.name }}</el-descriptions-item>
            <el-descriptions-item label="版本">{{ metaData.version || '-' }}</el-descriptions-item>
            <el-descriptions-item label="作者">{{ metaData.author || '-' }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ metaData.description || '-' }}</el-descriptions-item>
          </el-descriptions>
          <div v-if="metaData.HelpDoc" style="margin-top: 16px">
            <el-link :href="metaData.HelpDoc" target="_blank" type="primary">查看帮助文档</el-link>
          </div>
        </template>
        <el-empty v-else description="暂无模块信息" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)
const buttonLoading = ref(false)
const buttonSubmitLoading = ref(false)
const metaLoading = ref(false)

const tableData = ref<any[]>([])
const availableModules = ref<any[]>([])

// 新增/编辑
const dialogVisible = ref(false)
const dialogTitle = ref('添加供给')
const formRef = ref<FormInstance>()
const formData = reactive({
  id: null as number | null,
  name: '',
  module: 'pve',
  host: '',
  port: 8006,
  username: '',
  password: '',
  config: {} as Record<string, any>
})
const moduleConfigOptions = ref<any[]>([])
const formRules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  module: [{ required: true, message: '请选择模块', trigger: 'change' }],
  host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }]
}

// 按钮管理
const buttonDialogVisible = ref(false)
const buttonFormVisible = ref(false)
const buttonFormRef = ref<FormInstance>()
const buttonList = ref<any[]>([])
const currentProviderId = ref(0)
const buttonForm = reactive({ id: null as number | null, name: '', func: '', description: '', order: 0, is_enabled: 1 })
const buttonRules: FormRules = {
  name: [{ required: true, message: '请输入按钮名称', trigger: 'blur' }],
  func: [{ required: true, message: '请输入函数名', trigger: 'blur' }]
}

// 模块信息
const metaDialogVisible = ref(false)
const metaData = ref<any>(null)

const fetchModules = async () => {
  try {
    const data = await request.get({ url: '/api/admin/provision-modules/modules' })
    availableModules.value = data || []
  } catch { /* ignore */ }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/provisions' })
    tableData.value = res.data || res || []
  } catch { ElMessage.error('获取供给列表失败') } finally { loading.value = false }
}

const handleCreate = () => {
  dialogTitle.value = '添加供给'
  Object.assign(formData, { id: null, name: '', module: 'pve', host: '', port: 8006, username: '', password: '', config: {} })
  moduleConfigOptions.value = []
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑供给'
  Object.assign(formData, { id: row.id, name: row.name, module: row.module, host: row.host, port: row.port || 8006, username: row.username || '', password: '', config: row.config || {} })
  moduleConfigOptions.value = []
  fetchModuleConfig(row.module)
  dialogVisible.value = true
}

const handleModuleChange = (module: string) => {
  formData.config = {}
  fetchModuleConfig(module)
}

const fetchModuleConfig = async (module: string) => {
  try {
    const data = await request.get({ url: `/api/admin/provision-modules/config`, params: { name: module } })
    moduleConfigOptions.value = data || []
  } catch { moduleConfigOptions.value = [] }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/provisions/${formData.id}`, data: formData, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/provisions', data: formData, showSuccessMessage: true })
      }
      dialogVisible.value = false; fetchData()
    } catch { ElMessage.error('操作失败') } finally { submitLoading.value = false }
  })
}

const handleTest = async (row: any) => {
  try {
    const res = await request.post({ url: `/api/admin/provisions/${row.id}/test` })
    ElMessage.success(res?.message || '连接成功')
  } catch { ElMessage.error('连接失败') }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定删除该供给？', '提示', { type: 'warning' })
  try { await request.del({ url: `/api/admin/provisions/${row.id}` }); ElMessage.success('删除成功'); fetchData() }
  catch { ElMessage.error('删除失败') }
}

// 按钮管理
const handleViewButtons = async (row: any) => {
  currentProviderId.value = row.id
  buttonDialogVisible.value = true
  fetchButtonList()
}

const fetchButtonList = async () => {
  buttonLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/provision-modules/${currentProviderId.value}/buttons` })
    buttonList.value = data || []
  } catch { ElMessage.error('获取按钮列表失败') } finally { buttonLoading.value = false }
}

const showButtonForm = (row?: any) => {
  if (row) { Object.assign(buttonForm, { id: row.id, name: row.name, func: row.func, description: row.description || '', order: row.order || 0, is_enabled: row.is_enabled ?? 1 }) }
  else { Object.assign(buttonForm, { id: null, name: '', func: '', description: '', order: 0, is_enabled: 1 }) }
  buttonFormVisible.value = true
}

const handleButtonSubmit = async () => {
  if (!buttonFormRef.value) return
  await buttonFormRef.value.validate(async (valid) => {
    if (!valid) return
    buttonSubmitLoading.value = true
    try {
      if (buttonForm.id) {
        await request.put({ url: `/api/admin/provision-modules/${currentProviderId.value}/buttons/${buttonForm.id}`, data: buttonForm, showSuccessMessage: true })
      } else {
        await request.post({ url: `/api/admin/provision-modules/${currentProviderId.value}/buttons`, data: buttonForm, showSuccessMessage: true })
      }
      buttonFormVisible.value = false; fetchButtonList()
    } catch { ElMessage.error('操作失败') } finally { buttonSubmitLoading.value = false }
  })
}

const handleDeleteButton = async (row: any) => {
  try { await request.del({ url: `/api/admin/provision-modules/${currentProviderId.value}/buttons/${row.id}` }); ElMessage.success('删除成功'); fetchButtonList() }
  catch { ElMessage.error('删除失败') }
}

const handleExecuteButton = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定执行按钮 "${row.name}" 吗？`, '执行确认')
    const res = await request.post({ url: `/api/admin/provision-modules/${currentProviderId.value}/buttons/${row.id}/execute` })
    ElMessage.success(res?.message || '执行成功')
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('执行失败') }
}

// 模块信息
const handleViewMeta = async (row: any) => {
  metaDialogVisible.value = true
  metaLoading.value = true
  metaData.value = null
  try {
    const data = await request.get({ url: `/api/admin/provision-modules/meta`, params: { name: row.module } })
    metaData.value = data
  } catch { /* ignore */ } finally { metaLoading.value = false }
}

onMounted(() => { fetchData(); fetchModules() })
</script>

<style scoped lang="scss">
.provisions-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.button-header { margin-bottom: 16px; display: flex; justify-content: flex-end; }
</style>
