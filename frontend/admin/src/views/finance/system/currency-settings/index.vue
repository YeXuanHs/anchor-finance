<template>
  <div class="currency-settings-page">
    <!-- 说明文字 -->
    <el-alert type="info" :closable="false" class="page-desc">
      <template #title>
        <span>您可以通过不同的货币销售您的产品。客户可以在网站上选择他们合适的货币进行支付购买。</span>
      </template>
    </el-alert>

    <div class="action-bar">
      <el-button type="primary" @click="handleUpdateRate">汇率更新</el-button>
      <el-button type="warning" @click="handleUpdatePrice">价格更新</el-button>
    </div>

    <!-- 货币表格 -->
    <el-table 
      :data="currencyList" 
      v-loading="loading" 
      border 
      stripe
      empty-text="暂无数据"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" />
      
      <el-table-column prop="code" label="货币代码" width="120">
        <template #default="{ row }">
          <span class="currency-code">{{ row.code }}</span>
        </template>
      </el-table-column>
      
      <el-table-column prop="prefix" label="前缀（例如：￥）" width="140" />
      
      <el-table-column prop="suffix" label="后缀" width="100" />
      
      <el-table-column prop="decimal_places" label="格式" width="80" />
      
      <el-table-column prop="exchange_rate" label="汇率" width="120">
        <template #default="{ row }">
          {{ row.exchange_rate?.toFixed(5) || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button 
            size="small" 
            :type="row.is_default ? 'info' : 'primary'" 
            link
            :disabled="row.is_default"
            @click="handleSetDefault(row)"
          >
            {{ row.is_default ? '默认货币' : '设为默认货币' }}
          </el-button>
          <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" link :disabled="row.is_default" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 编辑弹窗 -->
    <el-dialog 
      v-model="editDialogVisible" 
      :title="isEdit ? '编辑货币' : '添加货币'" 
      width="500px"
      destroy-on-close
    >
      <el-form :model="formData" label-position="left" label-width="120px">
        <el-form-item label="货币代码" required>
          <el-input v-model="formData.code" placeholder="如 CNY" />
        </el-form-item>
        <el-form-item label="前缀">
          <el-input v-model="formData.prefix" placeholder="如 ￥" />
        </el-form-item>
        <el-form-item label="后缀">
          <el-input v-model="formData.suffix" placeholder="如 元" />
        </el-form-item>
        <el-form-item label="小数位数">
          <el-input-number v-model="formData.decimal_places" :min="0" :max="4" />
        </el-form-item>
        <el-form-item label="汇率">
          <el-input v-model="formData.exchange_rate" placeholder="如 1.00000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const currencyList = ref<any[]>([])
const editDialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const formData = ref({
  code: '',
  prefix: '',
  suffix: '',
  decimal_places: 3,
  exchange_rate: '1.00000'
})

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/currencies' })
    currencyList.value = data || []
  } catch (error) {
    console.error('fetch currencies failed:', error)
  } finally {
    loading.value = false
  }
}

const handleUpdateRate = async () => {
  try {
    await request.post({ url: '/api/admin/currencies/update-rates' })
    ElMessage.success('汇率更新成功')
    fetchList()
  } catch (error) {
    ElMessage.error('汇率更新失败')
  }
}

const handleUpdatePrice = async () => {
  try {
    await request.post({ url: '/api/admin/currencies/update-prices' })
    ElMessage.success('价格更新成功')
  } catch (error) {
    ElMessage.error('价格更新失败')
  }
}

const handleSetDefault = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/currencies/${row.id}/set-default` })
    ElMessage.success('已设为默认货币')
    fetchList()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  formData.value = {
    code: row.code,
    prefix: row.prefix || '',
    suffix: row.suffix || '',
    decimal_places: row.decimal_places || 3,
    exchange_rate: row.exchange_rate?.toString() || '1.00000'
  }
  editDialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除货币 "${row.code}" 吗？`, 
      '确认删除', 
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await request.del({ url: `/api/admin/currencies/${row.id}` })
    ElMessage.success('已删除')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/currencies/${editingId.value}`, data: formData.value })
    } else {
      await request.post({ url: '/api/admin/currencies', data: formData.value })
    }
    ElMessage.success('保存成功')
    editDialogVisible.value = false
    fetchList()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.currency-settings-page {
  padding: 20px;
}

.page-desc {
  margin-bottom: 16px;
}

.action-bar {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
}

.currency-code {
  font-weight: 600;
}
</style>
