<template>
  <div class="promo-code-page">
    <!-- 说明文字 -->
    <el-alert type="info" :closable="false" class="page-desc">
      <template #title>
        <span>优惠码可以用于给客户在购买产品时享受优惠折扣。</span>
      </template>
    </el-alert>

    <div class="action-bar">
      <el-button type="primary" @click="handleAdd">添加优惠码</el-button>
    </div>

    <!-- 优惠码表格 -->
    <el-table 
      :data="promoList" 
      v-loading="loading" 
      border 
      stripe
      empty-text="暂无数据"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" />
      
      <el-table-column prop="code" label="优惠码" width="140">
        <template #default="{ row }">
          <span class="promo-code">{{ row.code }}</span>
        </template>
      </el-table-column>
      
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="row.type === 'percentage' ? 'success' : 'primary'" size="small">
            {{ row.type === 'percentage' ? '百分比' : '固定金额' }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column prop="value" label="价值" width="100">
        <template #default="{ row }">
          {{ row.type === 'percentage' ? `${row.value}%` : `¥${row.value}` }}
        </template>
      </el-table-column>
      
      <el-table-column label="循环优惠" width="80">
        <template #default="{ row }">
          <el-icon v-if="row.recurring" color="var(--el-color-success)"><Check /></el-icon>
          <el-icon v-else color="var(--el-text-color-secondary)"><Close /></el-icon>
        </template>
      </el-table-column>
      
      <el-table-column label="已使用次数 / 最大使用次数" min-width="160">
        <template #default="{ row }">
          {{ row.used_count || 0 }} / {{ row.max_usage || '无限' }}
        </template>
      </el-table-column>
      
      <el-table-column prop="start_time" label="开始时间" width="170" />
      
      <el-table-column prop="expire_time" label="失效时间" width="170" />
      
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="warning" link @click="handleExpire(row)">立即过期</el-button>
          <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog 
      v-model="editDialogVisible" 
      :title="isEdit ? '编辑优惠码' : '添加优惠码'" 
      width="600px"
      destroy-on-close
    >
      <el-form :model="formData" label-position="top">
        <el-form-item label="优惠码" required>
          <el-input v-model="formData.code" placeholder="请输入优惠码" />
        </el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="formData.type" style="width: 100%">
            <el-option label="百分比" value="percentage" />
            <el-option label="固定金额" value="fixed" />
          </el-select>
        </el-form-item>
        <el-form-item label="价值" required>
          <el-input v-model="formData.value" placeholder="请输入价值" />
        </el-form-item>
        <el-form-item label="循环优惠">
          <el-switch v-model="formData.recurring" />
        </el-form-item>
        <el-form-item label="最大使用次数">
          <el-input-number v-model="formData.max_usage" :min="0" />
          <span class="form-tip">0表示无限次</span>
        </el-form-item>
        <el-form-item label="适用产品">
          <el-select v-model="formData.applied_products" multiple placeholder="留空表示所有产品" style="width: 100%">
            <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker v-model="formData.start_time" type="datetime" placeholder="选择开始时间" style="width: 100%" />
        </el-form-item>
        <el-form-item label="失效时间">
          <el-date-picker v-model="formData.expire_time" type="datetime" placeholder="选择失效时间" style="width: 100%" />
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
import { Check, Close } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const promoList = ref<any[]>([])
const products = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const editDialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const formData = ref({
  code: '',
  type: 'percentage',
  value: '',
  recurring: false,
  max_usage: 0,
  applied_products: [] as number[],
  start_time: '',
  expire_time: ''
})

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ 
      url: '/api/admin/promo-codes',
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    promoList.value = data?.list || data || []
    total.value = data?.total || promoList.value.length
  } catch (error) {
    console.error('fetch promo codes failed:', error)
  } finally {
    loading.value = false
  }
}

const fetchProducts = async () => {
  try {
    const data = await request.get({ url: '/api/admin/products' })
    products.value = data?.list || data || []
  } catch (error) {
    products.value = []
  }
}

const handleAdd = () => {
  isEdit.value = false
  editingId.value = null
  formData.value = {
    code: '', type: 'percentage', value: '', recurring: false,
    max_usage: 0, applied_products: [], start_time: '', expire_time: ''
  }
  editDialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  formData.value = {
    code: row.code,
    type: row.type,
    value: row.value?.toString() || '',
    recurring: row.recurring || false,
    max_usage: row.max_usage || 0,
    applied_products: row.applied_products || [],
    start_time: row.start_time || '',
    expire_time: row.expire_time || ''
  }
  editDialogVisible.value = true
}

const handleExpire = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要立即使此优惠码过期吗？', '确认', {
      confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning'
    })
    await request.post({ url: `/api/admin/promo-codes/${row.id}/expire` })
    ElMessage.success('已过期')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除优惠码 "${row.code}" 吗？`, 
      '确认删除', 
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await request.del({ url: `/api/admin/promo-codes/${row.id}` })
    ElMessage.success('已删除')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/promo-codes/${editingId.value}`, data: formData.value })
    } else {
      await request.post({ url: '/api/admin/promo-codes', data: formData.value })
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
  fetchProducts()
})
</script>

<style scoped lang="scss">
.promo-code-page {
  padding: 20px;
}

.page-desc {
  margin-bottom: 16px;
}

.action-bar {
  margin-bottom: 16px;
}

.promo-code {
  font-weight: 600;
  font-family: monospace;
}

.form-tip {
  margin-left: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
