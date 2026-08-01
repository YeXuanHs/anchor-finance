<template>
  <div class="cs-widgets-page">
    <div class="page-header">
      <div class="header-left">
        <h2>客服悬浮窗管理</h2>
        <span class="subtitle">管理前台客服入口和悬浮窗样式</span>
      </div>
      <div class="header-actions">
        <el-button @click="showSettingsDialog">全局设置</el-button>
        <el-button type="primary" @click="showAddDialog">添加客服入口</el-button>
      </div>
    </div>

    <!-- 客服列表 -->
    <el-table :data="widgetList" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" width="150" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag>{{ typeMap[row.type] || row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="content" label="内容/账号" min-width="200" />
      <el-table-column prop="url" label="跳转链接" min-width="200" show-overflow-tooltip />
      <el-table-column prop="sort_order" label="排序" width="80" align="center" />
      <el-table-column prop="is_active" label="状态" width="80" align="center">
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
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑客服入口' : '添加客服入口'" width="600px">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="formData.name" placeholder="例如: 在线QQ客服" />
        </el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="formData.type" placeholder="选择类型">
            <el-option label="QQ" value="qq" />
            <el-option label="微信" value="wechat" />
            <el-option label="Telegram" value="telegram" />
            <el-option label="电话" value="phone" />
            <el-option label="邮箱" value="email" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容/账号">
          <el-input v-model="formData.content" placeholder="QQ号/微信号/手机号等" />
        </el-form-item>
        <el-form-item label="跳转链接">
          <el-input v-model="formData.url" placeholder="点击后跳转的链接" />
        </el-form-item>
        <el-form-item label="图标URL">
          <el-input v-model="formData.icon" placeholder="自定义图标地址（可选）" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 全局设置对话框 -->
    <el-dialog v-model="settingsDialogVisible" title="悬浮窗全局设置" width="600px">
      <el-form :model="settingsForm" label-width="120px">
        <el-form-item label="启用悬浮窗">
          <el-switch v-model="settingsForm.enabled" />
        </el-form-item>
        <el-form-item label="显示位置">
          <el-radio-group v-model="settingsForm.position">
            <el-radio value="right">右侧</el-radio>
            <el-radio value="left">左侧</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="距底部距离">
          <el-input-number v-model="settingsForm.offset_bottom" :min="0" :max="500" />
        </el-form-item>
        <el-form-item label="主题色">
          <el-color-picker v-model="settingsForm.theme_color" />
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="settingsForm.title" />
        </el-form-item>
        <el-form-item label="欢迎语">
          <el-input v-model="settingsForm.welcome_text" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="工作时间">
          <el-input v-model="settingsForm.working_hours" placeholder="例如: 周一至周五 9:00-18:00" />
        </el-form-item>
        <el-form-item label="移动端显示">
          <el-switch v-model="settingsForm.show_on_mobile" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="settingsDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveSettings" :loading="savingSettings">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const typeMap: Record<string, string> = {
  qq: 'QQ', wechat: '微信', telegram: 'Telegram',
  phone: '电话', email: '邮箱', custom: '自定义'
}

const loading = ref(false)
const submitting = ref(false)
const savingSettings = ref(false)
const widgetList = ref<any[]>([])
const dialogVisible = ref(false)
const settingsDialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const formData = ref({
  name: '', type: 'qq', content: '', url: '', icon: '', sort_order: 0
})

const settingsForm = ref({
  enabled: true, position: 'right', offset_bottom: 100,
  theme_color: '#1890ff', title: '联系客服',
  welcome_text: '您好，请问有什么可以帮助您？',
  working_hours: '', show_on_mobile: true
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/admin/cs-widgets')
    widgetList.value = res.data || []
  } catch (e) { ElMessage.error('加载失败') }
  finally { loading.value = false }
}

const showAddDialog = () => {
  isEdit.value = false
  formData.value = { name: '', type: 'qq', content: '', url: '', icon: '', sort_order: 0 }
  dialogVisible.value = true
}

const showEditDialog = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formData.value.name) { ElMessage.warning('请输入名称'); return }
  submitting.value = true
  try {
    if (isEdit.value) {
      await request.put(`/api/admin/cs-widgets/${editId.value}`, formData.value)
    } else {
      await request.post('/api/admin/cs-widgets', formData.value)
    }
    ElMessage.success('操作成功')
    dialogVisible.value = false
    loadData()
  } catch (e) { ElMessage.error('操作失败') }
  finally { submitting.value = false }
}

const handleToggleActive = async (row: any) => {
  try {
    await request.put(`/api/admin/cs-widgets/${row.id}`, { is_active: row.is_active })
  } catch (e) { row.is_active = !row.is_active; ElMessage.error('更新失败') }
}

const handleDelete = async (id: number) => {
  try {
    await request.delete(`/api/admin/cs-widgets/${id}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (e) { ElMessage.error('删除失败') }
}

const showSettingsDialog = async () => {
  try {
    const res = await request.get('/api/admin/cs-widgets/settings')
    if (res.data) settingsForm.value = { ...settingsForm.value, ...res.data }
  } catch (e) {}
  settingsDialogVisible.value = true
}

const handleSaveSettings = async () => {
  savingSettings.value = true
  try {
    await request.put('/api/admin/cs-widgets/settings', settingsForm.value)
    ElMessage.success('保存成功')
    settingsDialogVisible.value = false
  } catch (e) { ElMessage.error('保存失败') }
  finally { savingSettings.value = false }
}

onMounted(loadData)
</script>

<style scoped>
.cs-widgets-page { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.header-left h2 { margin: 0 0 4px 0; font-size: 20px; }
.subtitle { color: #909399; font-size: 14px; }
</style>
