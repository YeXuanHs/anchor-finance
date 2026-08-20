<template>
  <div class="link-cause-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('finance.linkCause.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('finance.linkCause.addCause') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('finance.linkCause.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('finance.linkCause.causeName')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('finance.linkCause.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('finance.linkCause.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 树形表格 -->
      <el-table
        :data="treeData"
        v-loading="loading"
        row-key="id"
        default-expand-all
        :tree-props="{ children: 'children' }"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('finance.linkCause.causeName')" min-width="200" />
        <el-table-column prop="description" :label="$t('finance.linkCause.description')" min-width="250" show-overflow-tooltip />
        <el-table-column prop="sort_order" :label="$t('finance.linkCause.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('finance.linkCause.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('finance.linkCause.enabled') : $t('finance.linkCause.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('finance.linkCause.createdAt')" width="180" />
        <el-table-column :label="$t('finance.linkCause.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('finance.linkCause.edit') }}</el-button>
            <el-button type="primary" link @click="handleAddChild(row)">{{ $t('finance.linkCause.addChild') }}</el-button>
            <el-popconfirm :title="$t('finance.linkCause.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('finance.linkCause.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('finance.linkCause.parentCause')" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="treeData"
            :props="{ label: 'name', value: 'id' } as any"
            :placeholder="$t('finance.linkCause.selectParentCause')"
            clearable
            check-strictly
          />
        </el-form-item>
        <el-form-item :label="$t('finance.linkCause.causeName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('finance.linkCause.enterCauseName')" />
        </el-form-item>
        <el-form-item :label="$t('finance.linkCause.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('finance.linkCause.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('finance.linkCause.sort')" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item :label="$t('finance.linkCause.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('finance.linkCause.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('finance.linkCause.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

interface LinkCause {
  id: number
  name: string
  description: string
  sort_order: number
  status: number
  parent_id: number | null
  created_at: string
  children?: LinkCause[]
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('finance.linkCause.addCause'))
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: ''
})

const treeData = ref<LinkCause[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  description: '',
  sort_order: 0,
  status: 1,
  parent_id: undefined as number | undefined
})

const formRules: FormRules = {
  name: [
    { required: true, message: $t('finance.linkCause.enterCauseName'), trigger: 'blur' },
    { min: 2, max: 50, message: $t('finance.linkCause.lengthLimit'), trigger: 'blur' }
  ]
}

const fetchTreeData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/link-causes/tree',
      params: searchForm
    })
    treeData.value = data || []
  } catch (error) {
    console.error('获取原因树失败:', error)
    ElMessage.error($t('finance.linkCause.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  fetchTreeData()
}

const handleReset = () => {
  searchForm.keyword = ''
  fetchTreeData()
}

const handleAdd = () => {
  dialogTitle.value = $t('finance.linkCause.addCause')
  formData.id = undefined
  formData.name = ''
  formData.description = ''
  formData.sort_order = 0
  formData.status = 1
  formData.parent_id = undefined
  dialogVisible.value = true
}

const handleAddChild = (row: any) => {
  dialogTitle.value = $t('finance.linkCause.addChildCause')
  formData.id = undefined
  formData.name = ''
  formData.description = ''
  formData.sort_order = 0
  formData.status = 1
  formData.parent_id = row.id
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('finance.linkCause.editCause')
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/link-causes/${row.id}`
    })
    ElMessage.success($t('finance.linkCause.deleteSuccess'))
    fetchTreeData()
  } catch (error) {
    ElMessage.error($t('finance.linkCause.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({
          url: `/api/admin/link-causes/${formData.id}`,
          params: { ...formData }
        })
      } else {
        await request.post({
          url: '/api/admin/link-causes',
          params: { ...formData }
        })
      }

      ElMessage.success(formData.id ? $t('finance.linkCause.updateSuccess') : $t('finance.linkCause.addSuccess'))
      dialogVisible.value = false
      fetchTreeData()
    } catch (error) {
      ElMessage.error($t('finance.linkCause.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => {
  fetchTreeData()
})
</script>

<style scoped lang="scss">
.link-cause-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}
</style>