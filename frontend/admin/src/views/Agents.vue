<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">代理商管理</span>
          <div class="card-actions">
            <el-input v-model="searchKeyword" placeholder="搜索代理商名称/邮箱" clearable style="width: 240px">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-select v-model="filterStatus" placeholder="状态" clearable style="width: 110px">
              <el-option label="启用" value="active" /><el-option label="禁用" value="disabled" />
            </el-select>
            <el-button type="primary" @click="openModal()">
              <el-icon><Plus /></el-icon>添加代理商
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredAgents" stripe size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="代理商名称" width="140" />
        <el-table-column prop="email" label="邮箱" show-overflow-tooltip />
        <el-table-column prop="contactName" label="联系人" width="100" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column prop="commission" label="佣金比例" width="100">
          <template #default="{ row }"><span style="font-weight: 600; color: #0056FF">{{ row.commission }}%</span></template>
        </el-table-column>
        <el-table-column prop="totalSales" label="累计销售" width="120" sortable>
          <template #default="{ row }"><span style="font-weight: 600; color: #52c41a">¥{{ row.totalSales.toLocaleString() }}</span></template>
        </el-table-column>
        <el-table-column prop="userCount" label="用户数" width="80" sortable />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-switch :model-value="row.status === 'active'" size="small" @change="handleToggle(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" :icon="Edit" @click="openModal(row)" />
            <el-popconfirm title="确认删除该代理商？" @confirm="handleDelete(row.id)">
              <template #reference><el-button text type="danger" :icon="Delete" /></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="modalVisible" :title="editing ? '编辑代理商' : '添加代理商'" width="520px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="代理商名称" prop="name"><el-input v-model="formData.name" placeholder="请输入代理商名称" /></el-form-item>
        <el-form-item label="邮箱" prop="email"><el-input v-model="formData.email" placeholder="请输入邮箱" /></el-form-item>
        <el-form-item label="联系人"><el-input v-model="formData.contactName" placeholder="请输入联系人" /></el-form-item>
        <el-form-item label="电话"><el-input v-model="formData.phone" placeholder="请输入电话" /></el-form-item>
        <el-form-item label="佣金比例"><el-input-number v-model="formData.commission" :min="0" :max="100" :precision="1" style="width: 100%"><template #suffix>%</template></el-input-number></el-form-item>
        <el-form-item label="状态"><el-switch v-model="formData.status" active-value="active" inactive-value="disabled" active-text="启用" inactive-text="禁用" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const submitting = ref(false)
const modalVisible = ref(false)
const searchKeyword = ref('')
const filterStatus = ref<string | null>(null)
const formRef = ref<FormInstance>()
const editing = ref<any>(null)

const agents = ref([
  { id: 1, name: '华东代理商', email: 'east@example.com', contactName: '张经理', phone: '13800138001', commission: 15, totalSales: 256800, userCount: 156, status: 'active' },
  { id: 2, name: '华南代理商', email: 'south@example.com', contactName: '李经理', phone: '13800138002', commission: 12, totalSales: 189500, userCount: 98, status: 'active' },
  { id: 3, name: '华北代理商', email: 'north@example.com', contactName: '王经理', phone: '13800138003', commission: 10, totalSales: 134200, userCount: 67, status: 'active' },
  { id: 4, name: '西南代理商', email: 'west@example.com', contactName: '赵经理', phone: '13800138004', commission: 18, totalSales: 78900, userCount: 34, status: 'disabled' },
])

const filteredAgents = computed(() => {
  return agents.value.filter((a) => {
    if (searchKeyword.value.trim()) {
      const kw = searchKeyword.value.trim().toLowerCase()
      if (!a.name.toLowerCase().includes(kw) && !a.email.toLowerCase().includes(kw)) return false
    }
    if (filterStatus.value && a.status !== filterStatus.value) return false
    return true
  })
})

const formData = reactive({ name: '', email: '', contactName: '', phone: '', commission: 10, status: 'active' })
const rules: FormRules = {
  name: { required: true, message: '请输入代理商名称', trigger: 'blur' },
  email: { required: true, message: '请输入邮箱', trigger: 'blur' },
}

function openModal(agent?: any) {
  editing.value = agent || null
  if (agent) {
    Object.assign(formData, { name: agent.name, email: agent.email, contactName: agent.contactName, phone: agent.phone, commission: agent.commission, status: agent.status })
  } else {
    Object.assign(formData, { name: '', email: '', contactName: '', phone: '', commission: 10, status: 'active' })
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
      ElMessage.success('代理商更新成功')
    } else {
      agents.value.push({ id: Date.now(), ...formData, totalSales: 0, userCount: 0 })
      ElMessage.success('代理商添加成功')
    }
    modalVisible.value = false
  } finally { submitting.value = false }
}

function handleToggle(agent: any) {
  agent.status = agent.status === 'active' ? 'disabled' : 'active'
  ElMessage.success(`代理商「${agent.name}」已${agent.status === 'active' ? '启用' : '禁用'}`)
}

function handleDelete(id: number) {
  agents.value = agents.value.filter((a) => a.id !== id)
  ElMessage.success('代理商已删除')
}
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; }
.card-title { font-size: 16px; font-weight: 600; }
.card-actions { display: flex; align-items: center; gap: 12px; }
</style>
