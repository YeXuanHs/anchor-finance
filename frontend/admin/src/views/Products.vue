<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <el-tabs v-model="activeTab" class="header-tabs">
            <el-tab-pane label="产品列表" name="list" />
            <el-tab-pane label="产品组" name="groups" />
          </el-tabs>
          <div class="card-actions">
            <el-input
              v-if="activeTab === 'list'"
              v-model="searchKeyword"
              placeholder="搜索产品名称"
              clearable
              style="width: 220px"
              @clear="handleSearch"
              @keydown.enter="handleSearch"
            >
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-button type="primary" @click="activeTab === 'list' ? openProductModal() : openGroupModal()">
              <el-icon><Plus /></el-icon>{{ activeTab === 'list' ? '添加产品' : '添加产品组' }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- Product List -->
      <div v-if="activeTab === 'list'">
        <el-table :data="filteredProducts" v-loading="loading" stripe size="small">
          <el-table-column prop="id" label="ID" width="60" sortable />
          <el-table-column prop="name" label="名称" show-overflow-tooltip />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag size="small">{{ typeNameMap[row.type] || row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="价格(月)" width="110" sortable>
            <template #default="{ row }">
              <span style="font-weight: 600; color: #0056FF">¥{{ row.price }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-switch :model-value="row.status === 'active'" size="small" @change="handleToggleProduct(row)" />
            </template>
          </el-table-column>
          <el-table-column prop="stock" label="库存" width="80" sortable />
          <el-table-column label="操作" width="140" fixed="right">
            <template #default="{ row }">
              <el-button text type="primary" :icon="Edit" @click="openProductModal(row)" />
              <el-popconfirm :title="`确认${row.status === 'active' ? '下架' : '上架'}该产品？`" @confirm="handleToggleProduct(row)">
                <template #reference>
                  <el-button text type="warning" :icon="Bottom" />
                </template>
              </el-popconfirm>
              <el-popconfirm title="确认删除该产品？" @confirm="handleDeleteProduct(row.id)">
                <template #reference>
                  <el-button text type="danger" :icon="Delete" />
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- Product Groups -->
      <div v-if="activeTab === 'groups'">
        <el-table :data="groups" stripe size="small">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="name" label="产品组名称" />
          <el-table-column prop="description" label="描述" show-overflow-tooltip />
          <el-table-column prop="productCount" label="产品数量" width="100" />
          <el-table-column label="操作" width="140">
            <template #default="{ row }">
              <el-button text type="primary" :icon="Edit" @click="openGroupModal(row)" />
              <el-popconfirm title="确认删除该产品组？" @confirm="handleDeleteGroup(row.id)">
                <template #reference>
                  <el-button text type="danger" :icon="Delete" />
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- Product Dialog -->
    <el-dialog v-model="productModalVisible" :title="editingProduct ? '编辑产品' : '添加产品'" width="600px" destroy-on-close>
      <el-form ref="productFormRef" :model="productForm" :rules="productRules" label-width="100px">
        <el-form-item label="产品名称" prop="name">
          <el-input v-model="productForm.name" placeholder="请输入产品名称" />
        </el-form-item>
        <el-form-item label="产品类型" prop="type">
          <el-select v-model="productForm.type" placeholder="选择产品类型" style="width: 100%">
            <el-option v-for="o in typeOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属分组" prop="groupId">
          <el-select v-model="productForm.groupId" placeholder="选择产品组" style="width: 100%">
            <el-option v-for="o in groupOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格(月)" prop="price">
          <el-input-number v-model="productForm.price" :min="0" :precision="2" style="width: 100%">
            <template #prefix>¥</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="productForm.description" type="textarea" :rows="3" placeholder="产品描述" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="productForm.stock" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="productForm.status" active-value="active" inactive-value="inactive" active-text="上架" inactive-text="下架" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="productModalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleProductSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- Group Dialog -->
    <el-dialog v-model="groupModalVisible" :title="editingGroup ? '编辑产品组' : '添加产品组'" width="460px" destroy-on-close>
      <el-form ref="groupFormRef" :model="groupForm" :rules="groupRules" label-width="80px">
        <el-form-item label="组名" prop="name">
          <el-input v-model="groupForm.name" placeholder="请输入产品组名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="groupForm.description" type="textarea" :rows="2" placeholder="产品组描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="groupModalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleGroupSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, Edit, Delete, Bottom } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const activeTab = ref('list')
const loading = ref(false)
const submitting = ref(false)
const productModalVisible = ref(false)
const groupModalVisible = ref(false)
const searchKeyword = ref('')
const productFormRef = ref<FormInstance>()
const groupFormRef = ref<FormInstance>()
const editingProduct = ref<any>(null)
const editingGroup = ref<any>(null)

const typeOptions = [
  { label: '虚拟主机', value: 'hosting' },
  { label: '云服务器', value: 'vps' },
  { label: '域名', value: 'domain' },
  { label: 'SSL证书', value: 'ssl' },
  { label: '企业邮箱', value: 'email' },
]

const groups = ref([
  { id: 1, name: '虚拟主机', description: '各类虚拟主机产品', productCount: 3 },
  { id: 2, name: '云服务器', description: 'VPS / 云服务器产品', productCount: 2 },
  { id: 3, name: '域名服务', description: '域名注册与转入', productCount: 1 },
  { id: 4, name: '增值服务', description: 'SSL、CDN 等增值服务', productCount: 0 },
])

const groupOptions = computed(() => groups.value.map((g) => ({ label: g.name, value: g.id })))
const typeNameMap: Record<string, string> = { hosting: '虚拟主机', vps: '云服务器', domain: '域名', ssl: 'SSL证书', email: '企业邮箱' }

const products = ref([
  { id: 1, name: '基础版主机', type: 'hosting', groupId: 1, groupName: '虚拟主机', price: 299, stock: 100, status: 'active' },
  { id: 2, name: '高级版主机', type: 'hosting', groupId: 1, groupName: '虚拟主机', price: 599, stock: 50, status: 'active' },
  { id: 3, name: '企业版主机', type: 'hosting', groupId: 1, groupName: '虚拟主机', price: 1299, stock: 20, status: 'active' },
  { id: 4, name: '1核2G云服务器', type: 'vps', groupId: 2, groupName: '云服务器', price: 89, stock: 200, status: 'active' },
  { id: 5, name: '4核8G云服务器', type: 'vps', groupId: 2, groupName: '云服务器', price: 399, stock: 80, status: 'active' },
  { id: 6, name: '.com域名注册', type: 'domain', groupId: 3, groupName: '域名服务', price: 69, stock: 999, status: 'active' },
  { id: 7, name: 'DV SSL证书', type: 'ssl', groupId: 4, groupName: '增值服务', price: 199, stock: 500, status: 'inactive' },
])

const filteredProducts = computed(() => {
  if (!searchKeyword.value.trim()) return products.value
  const kw = searchKeyword.value.trim().toLowerCase()
  return products.value.filter((p) => p.name.toLowerCase().includes(kw))
})

const productForm = reactive({ name: '', type: '', groupId: null as number | null, price: 0, description: '', stock: 0, status: 'active' })
const productRules: FormRules = {
  name: { required: true, message: '请输入产品名称', trigger: 'blur' },
  type: { required: true, message: '请选择产品类型', trigger: 'change' },
  groupId: { required: true, message: '请选择产品组', trigger: 'change' },
  price: { required: true, message: '请输入价格', trigger: 'blur' },
}

const groupForm = reactive({ name: '', description: '' })
const groupRules: FormRules = { name: { required: true, message: '请输入产品组名称', trigger: 'blur' } }

function openProductModal(product?: any) {
  editingProduct.value = product || null
  if (product) {
    Object.assign(productForm, { name: product.name, type: product.type, groupId: product.groupId, price: product.price, description: '', stock: product.stock, status: product.status })
  } else {
    Object.assign(productForm, { name: '', type: '', groupId: null, price: 0, description: '', stock: 0, status: 'active' })
  }
  productModalVisible.value = true
}

function openGroupModal(group?: any) {
  editingGroup.value = group || null
  if (group) {
    Object.assign(groupForm, { name: group.name, description: group.description })
  } else {
    Object.assign(groupForm, { name: '', description: '' })
  }
  groupModalVisible.value = true
}

async function handleProductSubmit() {
  if (!productFormRef.value) return
  try { await productFormRef.value.validate() } catch { return }
  submitting.value = true
  try {
    ElMessage.success(editingProduct.value ? '产品更新成功' : '产品添加成功')
    productModalVisible.value = false
  } finally { submitting.value = false }
}

async function handleGroupSubmit() {
  if (!groupFormRef.value) return
  try { await groupFormRef.value.validate() } catch { return }
  submitting.value = true
  try {
    ElMessage.success(editingGroup.value ? '产品组更新成功' : '产品组添加成功')
    groupModalVisible.value = false
  } finally { submitting.value = false }
}

function handleToggleProduct(product: any) {
  product.status = product.status === 'active' ? 'inactive' : 'active'
  ElMessage.success(`产品「${product.name}」已${product.status === 'active' ? '上架' : '下架'}`)
}

function handleDeleteProduct(id: number) { products.value = products.value.filter((p) => p.id !== id); ElMessage.success('产品已删除') }
function handleDeleteGroup(id: number) { groups.value = groups.value.filter((g) => g.id !== id); ElMessage.success('产品组已删除') }
function handleSearch() {}
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-tabs {
  margin-bottom: -17px;
}

.card-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
