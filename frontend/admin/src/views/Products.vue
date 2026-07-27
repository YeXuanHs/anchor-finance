<template>
  <div>
    <n-card :bordered="false" rounded>
      <n-tabs v-model:value="activeTab" type="line" animated>
        <!-- Product List Tab -->
        <n-tab-pane name="list" tab="产品列表">
          <template #header-extra>
            <n-space style="padding: 12px 0 4px">
              <n-input
                v-model:value="searchKeyword"
                placeholder="搜索产品名称"
                clearable
                style="width: 240px"
                @clear="handleSearch"
                @keydown.enter="handleSearch"
              >
                <template #prefix>
                  <n-icon><SearchIcon /></n-icon>
                </template>
              </n-input>
              <n-button type="primary" @click="openProductModal()">
                <template #icon><n-icon><AddIcon /></n-icon></template>
                添加产品
              </n-button>
            </n-space>
          </template>

          <n-data-table
            :columns="productColumns"
            :data="filteredProducts"
            :loading="loading"
            :bordered="false"
            :row-key="(row: any) => row.id"
            size="small"
          />
        </n-tab-pane>

        <!-- Product Groups Tab -->
        <n-tab-pane name="groups" tab="产品组">
          <template #header-extra>
            <n-button type="primary" style="margin: 12px 0 4px" @click="openGroupModal()">
              <template #icon><n-icon><AddIcon /></n-icon></template>
              添加产品组
            </n-button>
          </template>

          <n-data-table
            :columns="groupColumns"
            :data="groups"
            :bordered="false"
            :row-key="(row: any) => row.id"
            size="small"
          />
        </n-tab-pane>
      </n-tabs>
    </n-card>

    <!-- Product Edit Modal -->
    <n-modal
      v-model:show="productModalVisible"
      preset="card"
      :title="editingProduct ? '编辑产品' : '添加产品'"
      style="width: 600px"
      :bordered="false"
      :segmented="{ content: true, footer: true }"
    >
      <n-form ref="productFormRef" :model="productForm" :rules="productRules" label-placement="left" label-width="100">
        <n-form-item label="产品名称" path="name">
          <n-input v-model:value="productForm.name" placeholder="请输入产品名称" />
        </n-form-item>
        <n-form-item label="产品类型" path="type">
          <n-select v-model:value="productForm.type" :options="typeOptions" placeholder="选择产品类型" />
        </n-form-item>
        <n-form-item label="所属分组" path="groupId">
          <n-select v-model:value="productForm.groupId" :options="groupOptions" placeholder="选择产品组" />
        </n-form-item>
        <n-form-item label="价格(月)" path="price">
          <n-input-number v-model:value="productForm.price" :min="0" :precision="2" style="width: 100%">
            <template #prefix>¥</template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input v-model:value="productForm.description" type="textarea" :rows="3" placeholder="产品描述" />
        </n-form-item>
        <n-form-item label="库存" path="stock">
          <n-input-number v-model:value="productForm.stock" :min="0" style="width: 100%" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <n-switch v-model:value="productForm.status" checked-value="active" unchecked-value="inactive">
            <template #checked>上架</template>
            <template #unchecked>下架</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="productModalVisible = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleProductSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- Group Modal -->
    <n-modal
      v-model:show="groupModalVisible"
      preset="card"
      :title="editingGroup ? '编辑产品组' : '添加产品组'"
      style="width: 460px"
      :bordered="false"
    >
      <n-form ref="groupFormRef" :model="groupForm" :rules="groupRules" label-placement="left" label-width="80">
        <n-form-item label="组名" path="name">
          <n-input v-model:value="groupForm.name" placeholder="请输入产品组名称" />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input v-model:value="groupForm.description" type="textarea" :rows="2" placeholder="产品组描述" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="groupModalVisible = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleGroupSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { h, ref, reactive, computed, defineComponent } from 'vue'
import {
  useMessage,
  NTag,
  NButton,
  NSwitch,
  NSpace,
  NPopconfirm,
  NTooltip,
  NIcon,
  type DataTableColumns,
  type FormInst,
  type FormRules,
} from 'naive-ui'
import {
  SearchOutline as SearchIcon,
  AddOutline as AddIcon,
  CreateOutline as EditIcon,
  TrashOutline as DeleteIcon,
} from '@vicons/ionicons5'

const message = useMessage()
const activeTab = ref('list')
const loading = ref(false)
const submitting = ref(false)
const productModalVisible = ref(false)
const groupModalVisible = ref(false)
const searchKeyword = ref('')
const productFormRef = ref<FormInst | null>(null)
const groupFormRef = ref<FormInst | null>(null)
const editingProduct = ref<any>(null)
const editingGroup = ref<any>(null)

// ---- Options ----
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

// ---- Products ----
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

// ---- Product Form ----
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
  name: { required: true, message: '请输入产品名称', trigger: 'blur' },
  type: { required: true, message: '请选择产品类型', trigger: 'change' },
  groupId: { required: true, type: 'number', message: '请选择产品组', trigger: 'change' },
  price: { required: true, type: 'number', message: '请输入价格', trigger: 'blur' },
}

// ---- Group Form ----
const groupForm = reactive({ name: '', description: '' })
const groupRules: FormRules = {
  name: { required: true, message: '请输入产品组名称', trigger: 'blur' },
}

// ---- Product Table Columns ----
const productColumns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60, sorter: (a, b) => a.id - b.id },
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: (row) => h(NTag, { size: 'small', round: true, bordered: false, type: 'info' }, { default: () => typeNameMap[row.type] || row.type }),
  },
  {
    title: '价格(月)',
    key: 'price',
    width: 110,
    sorter: (a, b) => a.price - b.price,
    render: (row) => h('span', { style: 'font-weight:600;color:#1890ff' }, `¥${row.price}`),
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NSwitch, {
        value: row.status === 'active',
        size: 'small',
        onUpdateValue: () => handleToggleProduct(row),
      }),
  },
  { title: '库存', key: 'stock', width: 80, sorter: (a, b) => a.stock - b.stock },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) =>
      h(NSpace, { size: 4 }, {
        default: () => [
          h(NTooltip, {}, {
            trigger: () =>
              h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => openProductModal(row) }, {
                icon: () => h(NIcon, null, { default: () => h(EditIcon) }),
              }),
            default: () => '编辑',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleToggleProduct(row) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(NButton, { size: 'small', quaternary: true, type: 'warning' }, {
                    icon: () => h(NIcon, null, { default: () => h(DownIcon) }),
                  }),
                default: () => row.status === 'active' ? '下架' : '上架',
              }),
            default: () => `确认${row.status === 'active' ? '下架' : '上架'}该产品？`,
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteProduct(row.id) }, {
            trigger: () =>
              h(NTooltip, {}, {
                trigger: () =>
                  h(NButton, { size: 'small', quaternary: true, type: 'error' }, {
                    icon: () => h(NIcon, null, { default: () => h(DeleteIcon) }),
                  }),
                default: () => '删除',
              }),
            default: () => `确认删除产品「${row.name}」？`,
          }),
        ],
      }),
  },
]

// ---- Group Table Columns ----
const groupColumns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '产品组名称', key: 'name' },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  { title: '产品数量', key: 'productCount', width: 100 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) =>
      h(NSpace, { size: 4 }, {
        default: () => [
          h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => openGroupModal(row) }, {
            icon: () => h(NIcon, null, { default: () => h(EditIcon) }),
            default: () => '编辑',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteGroup(row.id) }, {
            trigger: () =>
              h(NButton, { size: 'small', quaternary: true, type: 'error' }, {
                icon: () => h(NIcon, null, { default: () => h(DeleteIcon) }),
                default: () => '删除',
              }),
            default: () => `确认删除产品组「${row.name}」？`,
          }),
        ],
      }),
  },
]

// DownArrow icon for "下架"
const DownIcon = defineComponent({
  render: () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 512 512', fill: 'currentColor' }, [
    h('path', { d: 'M256 464l128-128H320V256h-32v80H128l128 128zm0-400v80h-32V144H128l128-128 128 128H320V64h-64z' }),
  ]),
})

// ---- Actions ----
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
  try { await productFormRef.value?.validate() } catch { return }
  submitting.value = true
  try {
    // TODO: API call
    message.success(editingProduct.value ? '产品更新成功' : '产品添加成功')
    productModalVisible.value = false
  } catch (err: any) {
    message.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleGroupSubmit() {
  try { await groupFormRef.value?.validate() } catch { return }
  submitting.value = true
  try {
    // TODO: API call
    message.success(editingGroup.value ? '产品组更新成功' : '产品组添加成功')
    groupModalVisible.value = false
  } catch (err: any) {
    message.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleToggleProduct(product: any) {
  product.status = product.status === 'active' ? 'inactive' : 'active'
  message.success(`产品「${product.name}」已${product.status === 'active' ? '上架' : '下架'}`)
}

function handleDeleteProduct(id: number) {
  products.value = products.value.filter((p) => p.id !== id)
  message.success('产品已删除')
}

function handleDeleteGroup(id: number) {
  groups.value = groups.value.filter((g) => g.id !== id)
  message.success('产品组已删除')
}

function handleSearch() {
  // filter is reactive via computed
}
</script>
