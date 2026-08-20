<template>
  <div class="page-container">
    <art-card :title="$t('menuManage.title')" shadow="never">
      <template #header-extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ $t('menuManage.addMenu') }}
        </el-button>
      </template>

      <el-table :data="tableData" v-loading="loading" row-key="id" :tree-props="{ children: 'children' }">
        <el-table-column prop="name" :label="$t('menuManage.menuName')" min-width="200" />
        <el-table-column prop="url" :label="$t('menuManage.link')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="icon" :label="$t('menuManage.icon')" width="120" />
        <el-table-column prop="order" :label="$t('menuManage.sort')" width="100" />
        <el-table-column prop="menu_type" :label="$t('menuManage.type')" width="120">
          <template #default="{ row }">
            <el-tag>{{ getMenuTypeText(row.menu_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('menuManage.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? $t('menuManage.visible') : $t('menuManage.hidden') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('menuManage.actions')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">{{ $t('menuManage.edit') }}</el-button>
            <el-button size="small" @click="handleAddChild(row)">{{ $t('menuManage.addChild') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ $t('menuManage.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" label-width="120px">
        <el-form-item :label="$t('menuManage.menuName')" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item :label="$t('menuManage.link')">
          <el-input v-model="formData.url" />
        </el-form-item>
        <el-form-item :label="$t('menuManage.icon')">
          <el-input v-model="formData.icon" />
        </el-form-item>
        <el-form-item :label="$t('menuManage.sort')">
          <el-input-number v-model="formData.order" :min="0" />
        </el-form-item>
        <el-form-item :label="$t('menuManage.menuType')" required>
          <el-select v-model="formData.menu_type">
            <el-option :label="$t('menuManage.userCenter')" :value="1" />
            <el-option :label="$t('menuManage.wwwHeader')" :value="2" />
            <el-option :label="$t('menuManage.wwwFooter')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('menuManage.status')">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('menuManage.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ $t('menuManage.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formData = ref({
  id: null,
  name: '',
  url: '',
  icon: '',
  order: 0,
  parent_id: 0,
  menu_type: 1,
  status: 1
})

const getMenuTypeText = (type: number) => {
  const map: Record<number, string> = { 1: $t('menuManage.userCenter'), 2: $t('menuManage.wwwHeader'), 3: $t('menuManage.wwwFooter') }
  return map[type] || $t('menuManage.unknown')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/menus/tree' })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  dialogTitle.value = $t('menuManage.addMenu')
  formData.value = { id: null, name: '', url: '', icon: '', order: 0, parent_id: 0, menu_type: 1, status: 1 }
  dialogVisible.value = true
}

const handleAddChild = (row: any) => {
  dialogTitle.value = $t('menuManage.addChildMenu')
  formData.value = { id: null, name: '', url: '', icon: '', order: 0, parent_id: row.id, menu_type: row.menu_type, status: 1 }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('menuManage.editMenu')
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    if (formData.value.id) {
      await request.put({ url: `/api/admin/menus/${formData.value.id}`, params: formData.value })
    } else {
      await request.post({ url: '/api/admin/menus', params: formData.value })
    }
    ElMessage.success($t('menuManage.operationSuccess'))
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm($t('menuManage.confirmDelete'), $t('menuManage.tip'), { type: 'warning' })
  try {
    await request.del({ url: `/api/admin/menus/${row.id}` })
    ElMessage.success($t('menuManage.deleteSuccess'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
