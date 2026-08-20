<template>
  <div class="product-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('productGroup.title') }}</span>
          <el-button type="primary" @click="handleAdd()"><el-icon><Plus /></el-icon>{{ $t('productGroup.addGroup') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('productGroup.groupName')">
          <el-input v-model="searchForm.name" :placeholder="$t('productGroup.enterGroupName')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="filteredTableData" v-loading="loading" style="width: 100%" row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" default-expand-all>
        <el-table-column prop="name" :label="$t('productGroup.groupName')" min-width="200" />
        <el-table-column prop="code" :label="$t('productGroup.groupCode')" width="150" />
        <el-table-column prop="description" :label="$t('common.description')" min-width="200" />
        <el-table-column prop="product_count" :label="$t('productGroup.productCount')" width="100" align="center" />
        <el-table-column prop="sort" :label="$t('productGroup.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('productGroup.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
        <el-table-column :label="$t('productGroup.operations')" width="250" fixed="right" align="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleAdd(row)">{{ $t('productGroup.addChild') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-popconfirm :title="$t('productGroup.confirmDelete')" @confirm="handleDelete(row)"><template #reference><el-button type="danger" link>{{ $t('common.delete') }}</el-button></template></el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('productGroup.parentGroup')" prop="parent_id">
          <el-tree-select v-model="formData.parent_id" :data="groupTreeData" :props="{ value: 'id', label: 'name', children: 'children' } as any" :placeholder="$t('productGroup.topGroup')" clearable check-strictly style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('productGroup.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('productGroup.enterGroupName')" />
        </el-form-item>
        <el-form-item :label="$t('productGroup.groupCode')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('productGroup.enterGroupCode')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('productGroup.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('productGroup.sort')" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('productGroup.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'ProductGroups' })

const loading = ref(false)
const submitLoading = ref(false)
const searchForm = reactive({ name: '' })
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref($t('productGroup.addGroup'))
const formRef = ref<FormInstance>()

const formData = reactive({ id: undefined as number | undefined, parent_id: undefined as number | undefined, name: '', code: '', description: '', sort: 0, status: 1 })

const formRules: FormRules = {
  name: [{ required: true, message: () => $t('productGroup.enterGroupName'), trigger: 'blur' }],
  code: [{ required: true, message: () => $t('productGroup.enterGroupCode'), trigger: 'blur' }]
}

const filteredTableData = computed(() => {
  if (!searchForm.name) return tableData.value
  return filterTree(tableData.value, searchForm.name.toLowerCase())
})

const filterTree = (data: any[], keyword: string): any[] => {
  return data.filter(item => {
    const nameMatch = item.name.toLowerCase().includes(keyword)
    const childrenMatch = item.children && filterTree(item.children, keyword).length > 0
    return nameMatch || childrenMatch
  }).map(item => item.children ? { ...item, children: filterTree(item.children, keyword) } : item)
}

const groupTreeData = computed(() => buildTreeSelectData(tableData.value, formData.id))
const buildTreeSelectData = (data: any[], excludeId?: number): any[] => data.filter(item => item.id !== excludeId).map(item => ({ id: item.id, name: item.name, children: item.children ? buildTreeSelectData(item.children, excludeId) : [] }))

const fetchGroups = async () => {
  loading.value = true
  try { const data = await request.get({ url: '/api/admin/product-groups' }); tableData.value = data || [] } catch { ElMessage.error($t('productGroup.fetchFailed')) } finally { loading.value = false }
}

const handleSearch = () => {}
const handleReset = () => { searchForm.name = '' }

const handleAdd = (parent?: any) => {
  dialogTitle.value = parent ? `${$t('productGroup.addChild')} - ${parent.name}` : $t('productGroup.addGroup')
  formData.id = undefined; formData.parent_id = parent?.id || undefined; formData.name = ''; formData.code = ''; formData.description = ''; formData.sort = 0; formData.status = 1; dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('productGroup.editGroup')
  formData.id = row.id; formData.parent_id = row.parent_id; formData.name = row.name; formData.code = row.code; formData.description = row.description; formData.sort = row.sort; formData.status = row.status; dialogVisible.value = true
}

const handleDelete = async (row: any) => { try { await request.del({ url: `/api/admin/product-groups/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchGroups() } catch { ElMessage.error($t('common.deleteFailed')) } }

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/product-groups/${formData.id}` : '/api/admin/product-groups'
      if (formData.id) await request.put({ url, params: formData })
      else await request.post({ url, params: formData })
      ElMessage.success(formData.id ? $t('common.updateSuccess') : $t('common.addSuccess')); dialogVisible.value = false; fetchGroups()
    } catch { ElMessage.error($t('common.operateFailed')) } finally { submitLoading.value = false }
  })
}

onMounted(() => { fetchGroups() })
</script>

<style scoped lang="scss">
.product-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
</style>
