<template>
  <div class="payments-page page-container">
    <div class="art-card">
      <div class="table-header"><h3>支付方式</h3><el-button type="primary" @click="showDialog=true; resetForm()"><el-icon><Plus /></el-icon>添加支付方式</el-button></div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="code" label="标识" width="120" />
        <el-table-column prop="icon" label="图标" width="80"><template #default="{row}"><img v-if="row.icon" :src="row.icon" style="height:24px" /></template></el-table-column>
        <el-table-column prop="is_enabled" label="状态" width="100">
          <template #default="{row}"><el-switch v-model="row.is_enabled" @change="toggleStatus(row)" /></template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="100" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{row}"><el-button link type="primary" @click="editRow(row)">编辑</el-button><el-button link type="danger" @click="deleteRow(row.id)">删除</el-button></template>
        </el-table-column>
      </el-table>
    </div>
    <el-dialog v-model="showDialog" :title="form.id?'编辑支付方式':'添加支付方式'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="标识"><el-input v-model="form.code" :disabled="!!form.id" /></el-form-item>
        <el-form-item label="图标URL"><el-input v-model="form.icon" /></el-form-item>
        <el-form-item label="手续费率"><el-input v-model="form.fee_rate" placeholder="0.006" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.is_enabled" /></el-form-item>
        <el-divider content-position="left">网关配置（JSON）</el-divider>
        <el-form-item label="配置"><el-input v-model="form.config" type="textarea" :rows="6" placeholder='{"app_id":"","app_secret":""}' /></el-form-item>
      </el-form>
      <template #footer><el-button @click="showDialog=false">取消</el-button><el-button type="primary" @click="saveForm">保存</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const list = ref<any[]>([])
const showDialog = ref(false)
const form = ref<any>({ id: 0, name: '', code: '', icon: '', fee_rate: 0, sort_order: 0, is_enabled: true, config: '{}' })

const resetForm = () => { form.value = { id: 0, name: '', code: '', icon: '', fee_rate: 0, sort_order: 0, is_enabled: true, config: '{}' } }

const fetchData = async () => {
  loading.value = true
  try { const { data } = await request.get('/api/admin/payment-gateways'); list.value = data.data || [] } catch {} finally { loading.value = false }
}

const editRow = (row: any) => { form.value = { ...row, config: typeof row.config === 'object' ? JSON.stringify(row.config) : row.config }; showDialog.value = true }

const saveForm = async () => {
  try {
    if (form.value.id) { await request.put(`/api/admin/payment-gateways/${form.value.id}`, form.value) }
    else { await request.post('/api/admin/payment-gateways', form.value) }
    ElMessage.success('保存成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('保存失败') }
}

const toggleStatus = async (row: any) => {
  try { await request.put(`/api/admin/payment-gateways/${row.id}`, { is_enabled: row.is_enabled }) } catch { row.is_enabled = !row.is_enabled }
}

const deleteRow = async (id: number) => {
  await ElMessageBox.confirm('确定删除？', '确认')
  try { await request.delete(`/api/admin/payment-gateways/${id}`); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
</style>
