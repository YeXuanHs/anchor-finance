<template>
  <div class="provisions-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('provisions.title') }}</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            {{ $t('provisions.addProvision') }}
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" :label="$t('provisions.id')" width="70" />
        <el-table-column prop="name" :label="$t('provisions.name')" min-width="150" />
        <el-table-column prop="module" :label="$t('provisions.module')" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.module }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="host" :label="$t('provisions.host')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('provisions.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('provisions.normal') : $t('provisions.abnormal') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_sync" :label="$t('provisions.lastSync')" width="180" />
        <el-table-column :label="$t('provisions.operations')" width="350" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleTest(row)">{{ $t('provisions.testConnection') }}</el-button>
            <el-button size="small" @click="handleEdit(row)">{{ $t('provisions.edit') }}</el-button>
            <el-button size="small" @click="handleViewButtons(row)">{{ $t('provisions.buttonManagement') }}</el-button>
            <el-button size="small" @click="handleViewMeta(row)">{{ $t('provisions.moduleInfo') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('provisions.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="650px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('provisions.name')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('provisions.inputName')" />
        </el-form-item>
        <el-form-item :label="$t('provisions.module')" prop="module">
          <el-select v-model="formData.module" :placeholder="$t('provisions.selectModule')" style="width: 100%" @change="handleModuleChange">
            <el-option v-for="mod in availableModules" :key="mod.name" :label="mod.name" :value="mod.name" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('provisions.hostAddress')" prop="host">
          <el-input v-model="formData.host" :placeholder="$t('provisions.inputHost')" />
        </el-form-item>
        <el-form-item :label="$t('provisions.port')">
          <el-input-number v-model="formData.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item :label="$t('provisions.username')">
          <el-input v-model="formData.username" :placeholder="$t('provisions.apiUsername')" />
        </el-form-item>
        <el-form-item :label="$t('provisions.password')">
          <el-input v-model="formData.password" type="password" show-password :placeholder="$t('provisions.apiPassword')" />
        </el-form-item>
        <!-- 动态配置项 -->
        <template v-if="moduleConfigOptions.length">
          <el-divider>{{ $t('provisions.moduleConfig') }}</el-divider>
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
        <el-button @click="dialogVisible = false">{{ $t('provisions.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('provisions.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 按钮管理对话框 -->
    <el-dialog v-model="buttonDialogVisible" :title="$t('provisions.buttonManageTitle')" width="700px">
      <div class="button-header">
        <el-button type="primary" size="small" @click="showButtonForm()">{{ $t('provisions.addButton') }}</el-button>
      </div>
      <el-table :data="buttonList" v-loading="buttonLoading" border size="small">
        <el-table-column prop="id" :label="$t('provisions.id')" width="70" />
        <el-table-column prop="name" :label="$t('provisions.buttonName')" width="150" />
        <el-table-column prop="func" :label="$t('provisions.funcName')" width="150" />
        <el-table-column prop="description" :label="$t('provisions.description')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="$t('provisions.operations')" width="200" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="showButtonForm(row)">{{ $t('provisions.editButton') }}</el-button>
            <el-button type="success" link @click="handleExecuteButton(row)">{{ $t('provisions.execute') }}</el-button>
            <el-popconfirm :title="$t('provisions.confirmDeleteButton')" @confirm="handleDeleteButton(row)">
              <template #reference><el-button type="danger" link>{{ $t('provisions.delete') }}</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 按钮表单对话框 -->
    <el-dialog v-model="buttonFormVisible" :title="buttonForm.id ? $t('provisions.editButtonText') : $t('provisions.addButtonText')" width="500px">
      <el-form :model="buttonForm" :rules="buttonRules" ref="buttonFormRef" label-width="100px">
        <el-form-item :label="$t('provisions.buttonName')" prop="name">
          <el-input v-model="buttonForm.name" :placeholder="$t('provisions.inputButtonName')" />
        </el-form-item>
        <el-form-item :label="$t('provisions.funcName')" prop="func">
          <el-input v-model="buttonForm.func" :placeholder="$t('provisions.inputFuncName')" />
        </el-form-item>
        <el-form-item :label="$t('provisions.description')">
          <el-input v-model="buttonForm.description" type="textarea" :rows="2" :placeholder="$t('provisions.buttonDesc')" />
        </el-form-item>
        <el-form-item :label="$t('provisions.sort')">
          <el-input-number v-model="buttonForm.order" :min="0" />
        </el-form-item>
        <el-form-item :label="$t('provisions.status')">
          <el-switch v-model="buttonForm.is_enabled" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="buttonFormVisible = false">{{ $t('provisions.cancel') }}</el-button>
        <el-button type="primary" @click="handleButtonSubmit" :loading="buttonSubmitLoading">{{ $t('provisions.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 模块信息对话框 -->
    <el-dialog v-model="metaDialogVisible" :title="$t('provisions.moduleInfoTitle')" width="600px">
      <div v-loading="metaLoading">
        <template v-if="metaData">
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('provisions.moduleName')">{{ metaData.name }}</el-descriptions-item>
            <el-descriptions-item :label="$t('provisions.version')">{{ metaData.version || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('provisions.author')">{{ metaData.author || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('provisions.description')" :span="2">{{ metaData.description || '-' }}</el-descriptions-item>
          </el-descriptions>
          <div v-if="metaData.HelpDoc" style="margin-top: 16px">
            <el-link :href="metaData.HelpDoc" target="_blank" type="primary">{{ $t('provisions.helpDoc') }}</el-link>
          </div>
        </template>
        <el-empty v-else :description="$t('provisions.noModuleInfo')" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const submitLoading = ref(false)
const buttonLoading = ref(false)
const buttonSubmitLoading = ref(false)
const metaLoading = ref(false)

const tableData = ref<any[]>([])
const availableModules = ref<any[]>([])

// 新增/编辑
const dialogVisible = ref(false)
const dialogTitle = ref($t('provisions.addProvisionTitle'))
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
const formRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('provisions.inputName'), trigger: 'blur' }],
  module: [{ required: true, message: $t('provisions.selectModule'), trigger: 'change' }],
  host: [{ required: true, message: $t('provisions.inputHost'), trigger: 'blur' }]
}))

// 按钮管理
const buttonDialogVisible = ref(false)
const buttonFormVisible = ref(false)
const buttonFormRef = ref<FormInstance>()
const buttonList = ref<any[]>([])
const currentProviderId = ref(0)
const buttonForm = reactive({ id: null as number | null, name: '', func: '', description: '', order: 0, is_enabled: 1 })
const buttonRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('provisions.inputButtonName'), trigger: 'blur' }],
  func: [{ required: true, message: $t('provisions.inputFuncName'), trigger: 'blur' }]
}))

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
  } catch { ElMessage.error($t('provisions.fetchListFailed')) } finally { loading.value = false }
}

const handleCreate = () => {
  dialogTitle.value = $t('provisions.addProvisionTitle')
  Object.assign(formData, { id: null, name: '', module: 'pve', host: '', port: 8006, username: '', password: '', config: {} })
  moduleConfigOptions.value = []
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('provisions.editProvisionTitle')
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
    } catch { ElMessage.error($t('provisions.operationFailed')) } finally { submitLoading.value = false }
  })
}

const handleTest = async (row: any) => {
  try {
    const res = await request.post({ url: `/api/admin/provisions/${row.id}/test` })
    ElMessage.success(res?.message || $t('provisions.connectionSuccess'))
  } catch { ElMessage.error($t('provisions.connectionFailed')) }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('provisions.confirmDelete'), $t('provisions.tips'), { type: 'warning' })
  try { await request.del({ url: `/api/admin/provisions/${row.id}` }); ElMessage.success($t('provisions.deleteSuccess')); fetchData() }
  catch { ElMessage.error($t('provisions.deleteFailed')) }
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
  } catch { ElMessage.error($t('provisions.fetchButtonsFailed')) } finally { buttonLoading.value = false }
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
    } catch { ElMessage.error($t('provisions.operationFailed')) } finally { buttonSubmitLoading.value = false }
  })
}

const handleDeleteButton = async (row: any) => {
  try { await request.del({ url: `/api/admin/provision-modules/${currentProviderId.value}/buttons/${row.id}` }); ElMessage.success($t('provisions.deleteSuccess')); fetchButtonList() }
  catch { ElMessage.error($t('provisions.deleteFailed')) }
}

const handleExecuteButton = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('provisions.confirmExecute', { name: row.name }), $t('provisions.executeConfirm'))
    const res = await request.post({ url: `/api/admin/provision-modules/${currentProviderId.value}/buttons/${row.id}/execute` })
    ElMessage.success(res?.message || $t('provisions.executeSuccess'))
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('provisions.executeFailed')) }
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
