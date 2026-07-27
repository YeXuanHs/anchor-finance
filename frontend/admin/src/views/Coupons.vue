<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">优惠券管理</span>
          <div class="card-actions">
            <el-input v-model="searchKeyword" placeholder="搜索优惠券名称" clearable style="width: 220px">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-select v-model="filterType" placeholder="优惠券类型" clearable style="width: 130px">
              <el-option v-for="o in typeOptions" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
            <el-select v-model="filterStatus" placeholder="状态" clearable style="width: 110px">
              <el-option v-for="o in statusOptions" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
            <el-button type="primary" @click="openCouponModal()">
              <el-icon><Plus /></el-icon>创建优惠券
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredCoupons" v-loading="loading" stripe size="small">
        <el-table-column prop="id" label="ID" width="60" sortable />
        <el-table-column prop="name" label="名称" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="90">
          <template #default="{ row }">
            <el-tag :type="row.type === 'fixed' ? 'success' : row.type === 'percent' ? 'warning' : 'primary'" size="small">
              {{ typeNameMap[row.type] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" label="面值" width="90" sortable>
          <template #default="{ row }">
            <span style="font-weight: 600; color: #52c41a">{{ row.type === 'percent' ? row.value + '%' : '¥' + row.value }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用条件" width="120">
          <template #default="{ row }">{{ row.minAmount > 0 ? `满¥${row.minAmount}可用` : '无门槛' }}</template>
        </el-table-column>
        <el-table-column label="有效期" width="200">
          <template #default="{ row }">{{ row.startDate }} ~ {{ row.endDate }}</template>
        </el-table-column>
        <el-table-column label="已使用" width="100" sortable>
          <template #default="{ row }">{{ row.usedCount }} / {{ row.totalCount }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-switch :model-value="row.status === 'active'" size="small" :disabled="row.status === 'expired'" @change="handleToggleCoupon(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="130" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" :icon="Edit" @click="openCouponModal(row)" />
            <el-popconfirm title="确认删除该优惠券？" @confirm="handleDeleteCoupon(row.id)">
              <template #reference>
                <el-button text type="danger" :icon="Delete" />
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Coupon Dialog -->
    <el-dialog v-model="couponModalVisible" :title="editingCoupon ? '编辑优惠券' : '创建优惠券'" width="640px" destroy-on-close>
      <el-form ref="couponFormRef" :model="couponForm" :rules="couponRules" label-width="100px">
        <el-form-item label="优惠券名称" prop="name">
          <el-input v-model="couponForm.name" placeholder="请输入优惠券名称" />
        </el-form-item>
        <el-form-item label="优惠券类型" prop="type">
          <el-select v-model="couponForm.type" placeholder="选择优惠券类型" style="width: 100%">
            <el-option v-for="o in typeOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="面值" prop="value">
          <el-input-number v-model="couponForm.value" :min="0" :precision="2" style="width: 100%">
            <template #prefix>{{ couponForm.type === 'percent' ? '%' : '¥' }}</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="使用条件">
          <el-input-number v-model="couponForm.minAmount" :min="0" :precision="2" style="width: 100%">
            <template #prefix>¥</template>
            <template #suffix>满减门槛</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="发放总量" prop="totalCount">
          <el-input-number v-model="couponForm.totalCount" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="有效期">
          <el-date-picker v-model="couponForm.dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="couponForm.status" active-value="active" inactive-value="inactive" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="couponModalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleCouponSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const couponModalVisible = ref(false)
const searchKeyword = ref('')
const filterType = ref<string | null>(null)
const filterStatus = ref<string | null>(null)
const couponFormRef = ref<FormInstance>()
const editingCoupon = ref<any>(null)

const typeOptions = [
  { label: '满减', value: 'full' },
  { label: '折扣', value: 'percent' },
  { label: '固定金额', value: 'fixed' },
]

const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
  { label: '已过期', value: 'expired' },
]

const typeNameMap: Record<string, string> = { full: '满减', percent: '折扣', fixed: '固定金额' }

const coupons = ref([
  { id: 1, name: '新用户专享券', type: 'fixed', value: 50, minAmount: 100, totalCount: 1000, usedCount: 356, status: 'active', startDate: '2026-01-01', endDate: '2026-12-31' },
  { id: 2, name: '满500减100', type: 'full', value: 100, minAmount: 500, totalCount: 500, usedCount: 128, status: 'active', startDate: '2026-03-01', endDate: '2026-09-30' },
  { id: 3, name: '8折优惠券', type: 'percent', value: 20, minAmount: 200, totalCount: 800, usedCount: 445, status: 'active', startDate: '2026-01-15', endDate: '2026-06-30' },
  { id: 4, name: '双十一狂欢券', type: 'fixed', value: 100, minAmount: 300, totalCount: 2000, usedCount: 2000, status: 'expired', startDate: '2025-11-01', endDate: '2025-11-30' },
  { id: 5, name: '年中大促折扣', type: 'percent', value: 15, minAmount: 150, totalCount: 600, usedCount: 89, status: 'inactive', startDate: '2026-06-01', endDate: '2026-06-30' },
  { id: 6, name: 'VIP专属券', type: 'fixed', value: 200, minAmount: 1000, totalCount: 100, usedCount: 23, status: 'active', startDate: '2026-01-01', endDate: '2026-12-31' },
])

const filteredCoupons = computed(() => {
  return coupons.value.filter((c) => {
    if (searchKeyword.value.trim() && !c.name.toLowerCase().includes(searchKeyword.value.trim().toLowerCase())) return false
    if (filterType.value && c.type !== filterType.value) return false
    if (filterStatus.value && c.status !== filterStatus.value) return false
    return true
  })
})

const couponForm = reactive({ name: '', type: '' as string, value: 0, minAmount: 0, totalCount: 100, dateRange: null as any, status: 'active' })
const couponRules: FormRules = {
  name: { required: true, message: '请输入优惠券名称', trigger: 'blur' },
  type: { required: true, message: '请选择优惠券类型', trigger: 'change' },
  value: { required: true, message: '请输入面值', trigger: 'blur' },
  totalCount: { required: true, message: '请输入发放总量', trigger: 'blur' },
}

function openCouponModal(coupon?: any) {
  editingCoupon.value = coupon || null
  if (coupon) {
    Object.assign(couponForm, { name: coupon.name, type: coupon.type, value: coupon.value, minAmount: coupon.minAmount, totalCount: coupon.totalCount, dateRange: null, status: coupon.status })
  } else {
    Object.assign(couponForm, { name: '', type: '', value: 0, minAmount: 0, totalCount: 100, dateRange: null, status: 'active' })
  }
  couponModalVisible.value = true
}

async function handleCouponSubmit() {
  if (!couponFormRef.value) return
  try { await couponFormRef.value.validate() } catch { return }
  submitting.value = true
  try {
    ElMessage.success(editingCoupon.value ? '优惠券更新成功' : '优惠券创建成功')
    couponModalVisible.value = false
  } finally { submitting.value = false }
}

function handleToggleCoupon(coupon: any) {
  if (coupon.status === 'expired') return
  coupon.status = coupon.status === 'active' ? 'inactive' : 'active'
  ElMessage.success(`优惠券「${coupon.name}」已${coupon.status === 'active' ? '启用' : '禁用'}`)
}

function handleDeleteCoupon(id: number) {
  coupons.value = coupons.value.filter((c) => c.id !== id)
  ElMessage.success('优惠券已删除')
}
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 16px; font-weight: 600; }
.card-actions { display: flex; align-items: center; gap: 12px; }
</style>
