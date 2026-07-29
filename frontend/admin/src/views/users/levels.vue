<template>
  <div class="user-levels-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="等级名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>用户等级</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加等级
        </el-button>
      </div>

      <el-table :data="levels" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="等级名称" width="150" />
        <el-table-column prop="level" label="等级值" width="100" align="center" />
        <el-table-column prop="discount" label="折扣" width="100" align="center">
          <template #default="{ row }">
            <span>{{ row.discount ? `${row.discount}%` : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="min_spending" label="最低消费" width="120" align="right">
          <template #default="{ row }">
            <span>¥{{ row.min_spending?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="auto_upgrade" label="自动升级" width="100">
          <template #default="{ row }">
            <el-tag :type="row.auto_upgrade ? 'success' : 'info'" size="small">
              {{ row.auto_upgrade ? '开启' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_count" label="用户数" width="100" align="center" />
        <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editLevel(row)">编辑</el-button>
            <el-button type="primary" link @click="openUpgradeRule(row)">升级规则</el-button>
            <el-button type="danger" link @click="deleteLevel(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchLevels"
          @current-change="fetchLevels"
        />
      </div>
    </div>

    <el-dialog v-model="showFormDialog" :title="editingItem ? '编辑等级' : '添加等级'" width="550px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item label="等级名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入等级名称" />
        </el-form-item>
        <el-form-item label="等级值" prop="level">
          <el-input-number v-model="formData.level" :min="1" :max="999" />
          <span class="form-tip">数值越大级别越高</span>
        </el-form-item>
        <el-form-item label="折扣">
          <el-input-number v-model="formData.discount" :min="0" :max="100" />
          <span style="margin-left: 8px;">%</span>
        </el-form-item>
        <el-form-item label="最低消费">
          <el-input-number v-model="formData.min_spending" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="自动升级">
          <el-switch v-model="formData.auto_upgrade" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFormDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showUpgradeDialog" :title="`自动升级规则 - ${currentLevel?.name}`" width="550px">
      <el-alert type="info" :closable="false" style="margin-bottom: 16px;">
        设置满足条件后用户自动升级到此等级
      </el-alert>
      <el-form :model="upgradeForm" label-width="120px">
        <el-form-item label="累计消费满">
          <el-input-number v-model="upgradeForm.total_spending" :min="0" :precision="2" />
          <span style="margin-left: 8px;">元</span>
        </el-form-item>
        <el-form-item label="订单数满">
          <el-input-number v-model="upgradeForm.total_orders" :min="0" />
          <span style="margin-left: 8px;">单</span>
        </el-form-item>
        <el-form-item label="注册天数满">
          <el-input-number v-model="upgradeForm.register_days" :min="0" />
          <span style="margin-left: 8px;">天</span>
        </el-form-item>
        <el-form-item label="满足条件">
          <el-radio-group v-model="upgradeForm.condition_type">
            <el-radio value="all">全部满足</el-radio>
            <el-radio value="any">任一满足</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUpgradeDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSaveUpgradeRule">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const levels = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showFormDialog = ref(false)
const showUpgradeDialog = ref(false)
const editingItem = ref<any>(null)
const currentLevel = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = ref({ keyword: '' })
const formData = ref({ name: '', level: 1, discount: 0, min_spending: 0, auto_upgrade: false, description: '' })
const formRules = {
  name: [{ required: true, message: '请输入等级名称', trigger: 'blur' }],
  level: [{ required: true, message: '请输入等级值', trigger: 'blur' }]
}
const upgradeForm = ref({ total_spending: 0, total_orders: 0, register_days: 0, condition_type: 'all' })

const fetchLevels = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/users/levels', { params })
    levels.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取等级列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchLevels() }
const resetSearch = () => { searchForm.value = { keyword: '' }; handleSearch() }

const openAddDialog = () => {
  editingItem.value = null
  formData.value = { name: '', level: 1, discount: 0, min_spending: 0, auto_upgrade: false, description: '' }
  showFormDialog.value = true
}

const editLevel = (level: any) => {
  editingItem.value = level
  formData.value = { name: level.name, level: level.level, discount: level.discount || 0, min_spending: level.min_spending || 0, auto_upgrade: level.auto_upgrade || false, description: level.description || '' }
  showFormDialog.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (editingItem.value) {
      await request.put(`/admin/api/v1/users/levels/${editingItem.value.id}`, formData.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/users/levels', formData.value)
      ElMessage.success('添加成功')
    }
    showFormDialog.value = false
    fetchLevels()
  } catch {
    ElMessage.error(editingItem.value ? '更新失败' : '添加失败')
  } finally {
    submitLoading.value = false
  }
}

const deleteLevel = async (level: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除等级「${level.name}」吗？`, '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/users/levels/${level.id}`)
    ElMessage.success('删除成功')
    fetchLevels()
  } catch {}
}

const openUpgradeRule = async (level: any) => {
  currentLevel.value = level
  try {
    const { data } = await request.get(`/admin/api/v1/users/levels/${level.id}/upgrade-rule`)
    if (data.data) {
      upgradeForm.value = { total_spending: data.data.total_spending || 0, total_orders: data.data.total_orders || 0, register_days: data.data.register_days || 0, condition_type: data.data.condition_type || 'all' }
    } else {
      upgradeForm.value = { total_spending: 0, total_orders: 0, register_days: 0, condition_type: 'all' }
    }
  } catch {
    upgradeForm.value = { total_spending: 0, total_orders: 0, register_days: 0, condition_type: 'all' }
  }
  showUpgradeDialog.value = true
}

const handleSaveUpgradeRule = async () => {
  submitLoading.value = true
  try {
    await request.put(`/admin/api/v1/users/levels/${currentLevel.value.id}/upgrade-rule`, upgradeForm.value)
    ElMessage.success('升级规则保存成功')
    showUpgradeDialog.value = false
  } catch {
    ElMessage.error('保存失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(fetchLevels)
</script>

<style scoped lang="scss">
.user-levels-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .form-tip { margin-left: 8px; color: var(--el-text-color-secondary); font-size: 12px; }
}
</style>
