<template>
  <div class="sms-templates-page">
    <art-card title="短信模板管理" shadow="never">
      <template #header>
        <div class="card-header">
          <span>短信模板管理</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            创建模板
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 短信配置 -->
        <el-tab-pane label="短信配置" name="config">
          <el-form :model="smsConfig" label-width="140px" style="max-width: 600px" v-loading="configLoading">
            <el-form-item label="启用国内短信">
              <el-switch v-model="smsConfig.shd_allow_sms_send" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="国内短信供应商">
              <el-select v-model="smsConfig.sms_operator" placeholder="选择供应商">
                <el-option v-for="item in smsConfig.sms_operator_list" :key="item.label" :label="item.value" :value="item.label" />
              </el-select>
            </el-form-item>
            <el-form-item label="启用国际短信">
              <el-switch v-model="smsConfig.shd_allow_sms_send_global" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="国际短信供应商">
              <el-select v-model="smsConfig.sms_operator_global" placeholder="选择供应商">
                <el-option v-for="item in smsConfig.sms_operator_list" :key="item.label" :label="item.value" :value="item.label" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveConfig">保存配置</el-button>
              <el-button @click="handleTestSms">测试发送</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 模板列表 -->
        <el-tab-pane label="短信模板" name="templates">
          <div class="search-bar">
            <el-input v-model="searchForm.template_id" placeholder="模板ID" clearable style="width: 200px" />
            <el-input v-model="searchForm.title" placeholder="模板标题" clearable style="width: 200px" />
            <el-button type="primary" @click="fetchTemplates">搜索</el-button>
            <el-button @click="handleCheckStatus">刷新审核状态</el-button>
          </div>

          <el-table :data="templateList" v-loading="templateLoading" stripe border>
            <el-table-column prop="id" label="编号" width="80" />
            <el-table-column prop="template_id" label="模板ID" width="150" />
            <el-table-column prop="type" label="类型" width="120" />
            <el-table-column prop="title" label="模板标题" />
            <el-table-column prop="content" label="模板内容" show-overflow-tooltip />
            <el-table-column label="审核状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" @click="handleSubmitCheck(row)">提交审核</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.limit"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next"
            style="margin-top: 16px; justify-content: flex-end"
            @size-change="fetchTemplates"
            @current-change="fetchTemplates"
          />
        </el-tab-pane>
      </el-tabs>
    </art-card>

    <!-- 创建/编辑模板对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑模板' : '创建模板'" width="600px">
      <el-form :model="templateForm" :rules="templateRules" ref="templateFormRef" label-width="100px">
        <el-form-item label="短信供应商" prop="sms_operator">
          <el-select v-model="templateForm.sms_operator" placeholder="选择供应商">
            <el-option v-for="item in smsConfig.sms_operator_list" :key="item.label" :label="item.value" :value="item.label" />
          </el-select>
        </el-form-item>
        <el-form-item label="模板标题" prop="title">
          <el-input v-model="templateForm.title" placeholder="请输入模板标题" />
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <el-input v-model="templateForm.content" type="textarea" :rows="4" placeholder="请输入模板内容，变量用{变量名}表示" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="templateForm.remark" placeholder="可选备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveTemplate" :loading="saveLoading">确定</el-button>
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

const activeTab = ref('config')
const configLoading = ref(false)
const templateLoading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const templateFormRef = ref<FormInstance>()

const smsConfig = reactive({
  shd_allow_sms_send: 0,
  shd_allow_sms_send_global: 0,
  sms_operator: '',
  sms_operator_global: '',
  sms_operator_list: [] as Array<{ label: string; value: string }>
})

const searchForm = reactive({ template_id: '', title: '' })
const pagination = reactive({ page: 1, limit: 10, total: 0 })
const templateList = ref<any[]>([])

const templateForm = reactive({
  sms_operator: '',
  title: '',
  content: '',
  remark: '',
  range_type: 0
})

const templateRules: FormRules = {
  sms_operator: [{ required: true, message: '请选择供应商', trigger: 'change' }],
  title: [{ required: true, message: '请输入模板标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入模板内容', trigger: 'blur' }]
}

const getStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'warning', 2: 'success', 3: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '未提交', 1: '审核中', 2: '已通过', 3: '未通过' }
  return map[status] || '未知'
}

const fetchConfig = async () => {
  configLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/config/message/mobile' })
    if (res?.msg_config) Object.assign(smsConfig, res.msg_config)
  } catch (error) {
    console.error(error)
  } finally {
    configLoading.value = false
  }
}

const handleSaveConfig = async () => {
  try {
    await request.post({ url: '/api/admin/config/message/mobile', data: smsConfig, showSuccessMessage: true })
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const handleTestSms = async () => {
  try {
    const { value } = await ElMessageBox.prompt('请输入测试手机号', '测试短信')
    await request.post({ url: '/api/admin/config/message/sms/test', data: { phone: value } })
    ElMessage.success('测试短信已发送')
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('发送失败')
  }
}

const fetchTemplates = async () => {
  templateLoading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/config/message/template/list',
      params: { ...searchForm, page: pagination.page, limit: pagination.limit }
    })
    if (res) {
      templateList.value = res.templates || []
      pagination.total = res.total || 0
    }
  } catch (error) {
    console.error(error)
  } finally {
    templateLoading.value = false
  }
}

const handleCheckStatus = async () => {
  try {
    await request.post({ url: '/api/admin/config/message/template/status' })
    ElMessage.success('审核状态已刷新')
    fetchTemplates()
  } catch (error) {
    ElMessage.error('刷新失败')
  }
}

const showCreateDialog = async () => {
  isEdit.value = false
  editingId.value = null
  templateForm.sms_operator = smsConfig.sms_operator
  templateForm.title = ''
  templateForm.content = ''
  templateForm.remark = ''
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  templateForm.sms_operator = row.sms_operator
  templateForm.title = row.title
  templateForm.content = row.content
  templateForm.remark = ''
  dialogVisible.value = true
}

const handleSaveTemplate = async () => {
  if (!templateFormRef.value) return
  await templateFormRef.value.validate(async (valid) => {
    if (!valid) return
    saveLoading.value = true
    try {
      if (isEdit.value && editingId.value) {
        await request.put({ url: `/api/admin/config/message/template/${editingId.value}`, data: templateForm, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/config/message/template', data: templateForm, showSuccessMessage: true })
      }
      dialogVisible.value = false
      fetchTemplates()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      saveLoading.value = false
    }
  })
}

const handleSubmitCheck = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定提交模板 "${row.title}" 进行审核吗？`, '提示')
    await request.post({ url: '/api/admin/config/message/template/check', data: { ids: [row.id], type: row.sms_operator } })
    ElMessage.success('已提交审核')
    fetchTemplates()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('提交失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定删除模板 "${row.title}" 吗？`, '提示')
    await request.delete({ url: '/api/admin/config/message/template', data: { ids: [row.id], type: row.sms_operator }, showSuccessMessage: true })
    fetchTemplates()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleTabChange = (tab: string) => {
  if (tab === 'templates') fetchTemplates()
  if (tab === 'config') fetchConfig()
}

onMounted(() => fetchConfig())
</script>

<style scoped lang="scss">
.sms-templates-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.search-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
