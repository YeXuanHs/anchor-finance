<template>
  <div class="products-page">
    <el-card class="admin-card" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- Product List Tab -->
        <el-tab-pane label="产品列表" name="list">
          <div class="tab-header">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索产品名称"
              clearable
              style="width: 240px"
              :prefix-icon="Search"
            />
            <el-button type="primary" :icon="Plus" @click="openProductModal()">添加产品</el-button>
          </div>

          <el-table :data="filteredProducts" style="width: 100%" v-loading="loading" size="default">
            <el-table-column prop="id" label="ID" width="60" sortable />
            <el-table-column prop="name" label="名称" show-overflow-tooltip />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag size="small">{{ typeNameMap[row.type] || row.type }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="price" label="价格(月)" width="110" sortable>
              <template #default="{ row }">
                <span class="price">¥{{ row.price }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  active-value="active"
                  inactive-value="inactive"
                  @change="handleToggleProduct(row)"
                />
              </template>
            </el-table-column>
            <el-table-column prop="stock" label="库存" width="80" sortable />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link :icon="Edit" @click="openProductModal(row)">编辑</el-button>
                <el-popconfirm
                  title="确认删除该产品？"
                  @confirm="handleDeleteProduct(row.id)"
                >
                  <template #reference>
                    <el-button type="danger" link :icon="Delete">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Product Groups Tab -->
        <el-tab-pane label="产品组" name="groups">
          <div class="tab-header">
            <el-button type="primary" :icon="Plus" @click="openGroupModal()">添加产品组</el-button>
          </div>

          <el-table :data="groups" style="width: 100%" size="default">
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="name" label="产品组名称" />
            <el-table-column prop="description" label="描述" show-overflow-tooltip />
            <el-table-column prop="productCount" label="产品数量" width="100" />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link :icon="Edit" @click="openGroupModal(row)">编辑</el-button>
                <el-popconfirm
                  title="确认删除该产品组？"
                  @confirm="handleDeleteGroup(row.id)"
                >
                  <template #reference>
                    <el-button type="danger" link :icon="Delete">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Product Edit Modal -->
    <el-dialog
      v-model="productModalVisible"
      :title="editingProduct ? '编辑产品' : '添加产品'"
      width="600px"
    >
      <el-form ref="productFormRef" :model="productForm" :rules="productRules" label-width="100px">
        <el-form-item label="产品名称" prop="name">
          <el-input v-model="productForm.name" placeholder="请输入产品名称" />
        </el-form-item>
        <el-form-item label="产品类型" prop="type">
          <el-select v-model="productForm.type" placeholder="选择产品类型" style="width: 100%">
            <el-option v-for="opt in typeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属分组" prop="groupId">
          <el-select v-model="productForm.groupId" placeholder="选择产品组" style="width: 100%">
            <el-option v-for="opt in groupOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格(月)" prop="price">
          <el-input-number v-model="productForm.price" :min="0" :precision="2" style="width: 100%">
            <template #prefix>¥</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="productForm.description" type="textarea" :rows="3" placeholder="产品描述" />
        </el-form-item>
        <el-form-item label="库存" prop="stock">
          <el-input-number v-model="productForm.stock" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="productForm.status"
            active-value="active"
            inactive-value="inactive"
            active-text="上架"
            inactive-text="下架"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="productModalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleProductSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- Group Modal -->
    <el-dialog
      v-model="groupModalVisible"
      :title="editingGroup ? '编辑产品组' : '添加产品组'"
      width="460px"
    >
      <el-form ref="groupFormRef" :model="groupForm" :rules="groupRules" label-width="80px">
        <el-form-item label="组名" prop="name">
          <el-input v-model="groupForm.name" placeholder="请输入产品组名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
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
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

definePageMeta({
  layout: 'admin',
})

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

// Options
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

const typeNameMap: Record<string, string> = {
  hosting: '虚拟主机',
  vps: '云服务器',
  domain: '域名',
  ssl: 'SSL证书',
  email: '企业邮箱',
}

// Products
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

// Product Form
const productForm = reactive({
  name: '',
  type: null as string | null,
  groupId: null as number | null,
  price: 0,
  description: '',
  stock: 0,
  status: 'active',
})

const productRules: FormRules = {
  name: [{ required: true, message: '请输入产品名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择产品类型', trigger: 'change' }],
  groupId: [{ required: true, message: '请选择产品组', trigger: 'change' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }],
}

// Group Form
const groupForm = reactive({ name: '', description: '' })
const groupRules: FormRules = {
  name: [{ required: true, message: '请输入产品组名称', trigger: 'blur' }],
}

// Actions
function openProductModal(product?: any) {
  editingProduct.value = product || null
  if (product) {
    Object.assign(productForm, {
      name: product.name,
      type: product.type,
      groupId: product.groupId,
      price: product.price,
      description: '',
      stock: product.stock,
      status: product.status,
    })
  } else {
    Object.assign(productForm, { name: '', type: null, groupId: null, price: 0, description: '', stock: 0, status: 'active' })
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
  await productFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      ElMessage.success(editingProduct.value ? '产品更新成功' : '产品添加成功')
      productModalVisible.value = false
    } catch (err: any) {
      ElMessage.error(err.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

async function handleGroupSubmit() {
  if (!groupFormRef.value) return
  await groupFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      ElMessage.success(editingGroup.value ? '产品组更新成功' : '产品组添加成功')
      groupModalVisible.value = false
    } catch (err: any) {
      ElMessage.error(err.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

function handleToggleProduct(product: any) {
  ElMessage.success(`产品「${product.name}」已${product.status === 'active' ? '上架' : '下架'}`)
}

function handleDeleteProduct(id: number) {
  products.value = products.value.filter((p) => p.id !== id)
  ElMessage.success('产品已删除')
}

function handleDeleteGroup(id: number) {
  groups.value = groups.value.filter((g) => g.id !== id)
  ElMessage.success('产品组已删除')
}
</script>

<style scoped>
.products-page {
  padding: 0;
}

.tab-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.price {
  font-weight: 600;
  color: #409eff;
}
</style>
