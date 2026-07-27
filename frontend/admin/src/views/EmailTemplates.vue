<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">邮件模板管理</span>
          <el-button type="primary" @click="openModal()">
            <el-icon><Plus /></el-icon>添加模板
          </el-button>
        </div>
      </template>

      <el-table :data="templates" stripe size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="模板名称" width="180" />
        <el-table-column prop="code" label="模板代码" width="160">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="邮件主题" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" width="160" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openModal(row)">编辑</el-button>
            <el-button text type="info" @click="handleTest(row)">测试</el-button>
            <el-popconfirm title="确认删除该模板？" @confirm="handleDelete(row.id)">
              <template #reference><el-button text type="danger">删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="modalVisible" :title="editing ? '编辑模板' : '添加模板'" width="720px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="模板名称" prop="name"><el-input v-model="formData.name" placeholder="请输入模板名称" /></el-form-item>
        <el-form-item label="模板代码" prop="code"><el-input v-model="formData.code" placeholder="如: verify_email, order_notify" /></el-form-item>
        <el-form-item label="邮件主题" prop="subject"><el-input v-model="formData.subject" placeholder="邮件主题" /></el-form-item>
        <el-form-item label="邮件内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="12" placeholder="支持HTML，可用变量: {{username}}, {{code}}, {{order_no}}" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.status" active-value="active" inactive-value="inactive" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const submitting = ref(false)
const modalVisible = ref(false)
const formRef = ref<FormInstance>()
const editing = ref<any>(null)

const templates = ref([
  { id: 1, name: '注册验证码', code: 'verify_email', subject: '【智简魔方】邮箱验证码', status: 'active', updatedAt: '2026-07-20 10:00', content: '<p>尊敬的 {{username}}，您的验证码是 {{code}}，5分钟内有效。</p>' },
  { id: 2, name: '订单确认通知', code: 'order_confirm', subject: '【智简魔方】订单确认 - {{order_no}}', status: 'active', updatedAt: '2026-07-18 14:30', content: '<p>您的订单 {{order_no}} 已确认，金额 ¥{{amount}}。</p>' },
  { id: 3, name: '密码重置', code: 'reset_password', subject: '【智简魔方】密码重置', status: 'active', updatedAt: '2026-07-15 09:00', content: '<p>请点击链接重置密码: {{reset_link}}</p>' },
  { id: 4, name: '工单回复通知', code: 'ticket_reply', subject: '【智简魔方】工单 {{ticket_no}} 有新回复', status: 'active', updatedAt: '2026-07-10 16:20', content: '<p>您的工单 {{ticket_no}} 有新回复，请登录查看。</p>' },
  { id: 5, name: '账单到期提醒', code: 'invoice_due', subject: '【智简魔方】账单到期提醒', status: 'inactive', updatedAt: '2026-07-05 11:00', content: '<p>您的账单 {{invoice_no}} 即将到期，请及时支付。</p>' },
])

const formData = reactive({ name: '', code: '', subject: '', content: '', status: 'active' })
const rules: FormRules = {
  name: { required: true, message: '请输入模板名称', trigger: 'blur' },
  code: { required: true, message: '请输入模板代码', trigger: 'blur' },
  subject: { required: true, message: '请输入邮件主题', trigger: 'blur' },
  content: { required: true, message: '请输入邮件内容', trigger: 'blur' },
}

function openModal(template?: any) {
  editing.value = template || null
  if (template) {
    Object.assign(formData, { name: template.name, code: template.code, subject: template.subject, content: template.content, status: template.status })
  } else {
    Object.assign(formData, { name: '', code: '', subject: '', content: '', status: 'active' })
  }
  modalVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return
  try { await formRef.value.validate() } catch { return }
  submitting.value = true
  try {
    if (editing.value) {
      Object.assign(editing.value, formData)
      ElMessage.success('模板更新成功')
    } else {
      templates.value.push({ id: Date.now(), ...formData, updatedAt: new Date().toLocaleString() })
      ElMessage.success('模板添加成功')
    }
    modalVisible.value = false
  } finally { submitting.value = false }
}

function handleTest(template: any) { ElMessage.success(`测试邮件已发送: ${template.name}`) }
function handleDelete(id: number) { templates.value = templates.value.filter((t) => t.id !== id); ElMessage.success('模板已删除') }
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 16px; font-weight: 600; }
</style>
