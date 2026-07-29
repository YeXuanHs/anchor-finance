<template>
  <div class="currencies-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>货币管理</h3>
        <el-button type="primary" @click="showDialog = true; resetForm()">
          <el-icon><Plus /></el-icon>添加货币
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="代码" width="100">
          <template #default="{ row }">
            <span style="font-weight: 600;">{{ row.code }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="symbol" label="符号" width="80" />
        <el-table-column label="汇率" width="120">
          <template #default="{ row }">
            <span>{{ row.exchange_rate?.toFixed(4) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_default" label="默认货币" width="110">
          <template #default="{ row }">
            <el-tag :type="row.is_default ? 'success' : 'info'" size="small">{{ row.is_default ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.status" active-value="active" inactive-value="disabled" @change="toggleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editRow(row)">编辑</el-button>
            <el-button link type="warning" @click="showRateDialog = true; rateForm.currency_id = row.id; rateForm.code = row.code; rateForm.exchange_rate = row.exchange_rate">调整汇率</el-button>
            <el-button link type="success" v-if="!row.is_default" @click="setDefault(row.id)">设为默认</el-button>
            <el-button link type="danger" @click="deleteRow(row.id)" :disabled="row.is_default">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showDialog" :title="form.id ? '编辑货币' : '添加货币'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="货币代码">
          <el-input v-model="form.code" placeholder="如: USD, CNY, EUR" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="货币名称"><el-input v-model="form.name" placeholder="如: 美元, 人民币" /></el-form-item>
        <el-form-item label="货币符号"><el-input v-model="form.symbol" placeholder="如: $, ¥, €" /></el-form-item>
        <el-form-item label="汇率">
          <el-input-number v-model="form.exchange_rate" :min="0" :precision="4" />
          <span style="margin-left: 8px; color: var(--text-secondary); font-size: 12px;">相对于默认货币</span>
        </el-form-item>
        <el-form-item label="小数位数"><el-input-number v-model="form.decimal_places" :min="0" :max="6" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.status" active-value="active" inactive-value="disabled" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showRateDialog" title="调整汇率" width="400px">
      <el-form :model="rateForm" label-width="80px">
        <el-form-item label="货币">{{ rateForm.code }}</el-form-item>
        <el-form-item label="当前汇率">{{ rateForm.exchange_rate?.toFixed(4) }}</el-form-item>
        <el-form-item label="新汇率">
          <el-input-number v-model="rateForm.new_rate" :min="0" :precision="4" style="width: 200px" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRateDialog = false">取消</el-button>
        <el-button type="primary" @click="updateRate">确认修改</el-button>
      </template>
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
const showRateDialog = ref(false)
const form = ref<any>({ id: 0, code: '', name: '', symbol: '', exchange_rate: 1, decimal_places: 2, status: 'active' })
const rateForm = ref({ currency_id: 0, code: '', exchange_rate: 0, new_rate: 1 })

const resetForm = () => { form.value = { id: 0, code: '', name: '', symbol: '', exchange_rate: 1, decimal_places: 2, status: 'active' } }

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/currencies')
    list.value = data.data || data || []
  } catch {} finally { loading.value = false }
}

const editRow = (row: any) => { form.value = { ...row }; showDialog.value = true }

const saveForm = async () => {
  try {
    if (form.value.id) { await request.put(`/admin/api/v1/currencies/${form.value.id}`, form.value) }
    else { await request.post('/admin/api/v1/currencies', form.value) }
    ElMessage.success('保存成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('保存失败') }
}

const toggleStatus = async (row: any) => {
  try { await request.put(`/admin/api/v1/currencies/${row.id}`, { status: row.status }) } catch { row.status = row.status === 'active' ? 'disabled' : 'active' }
}

const setDefault = async (id: number) => {
  await ElMessageBox.confirm('确定将该货币设为默认？所有汇率将基于此货币重新计算。', '确认')
  try { await request.put(`/admin/api/v1/currencies/${id}/default`); ElMessage.success('设置成功'); fetchData() } catch { ElMessage.error('设置失败') }
}

const deleteRow = async (id: number) => {
  await ElMessageBox.confirm('确定删除该货币？', '确认')
  try { await request.delete(`/admin/api/v1/currencies/${id}`); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}

const updateRate = async () => {
  try {
    await request.put(`/admin/api/v1/currencies/${rateForm.value.currency_id}/rate`, { exchange_rate: rateForm.value.new_rate })
    ElMessage.success('汇率更新成功'); showRateDialog.value = false; fetchData()
  } catch { ElMessage.error('汇率更新失败') }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
</style>
