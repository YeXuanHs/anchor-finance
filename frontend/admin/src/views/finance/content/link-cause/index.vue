<template>
  <div class="link-cause-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>关联原因管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加原因
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="原因名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
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
        <el-table-column prop="name" label="原因名称" min-width="200" />
        <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleAddChild(row)">添加子原因</el-button>
            <el-popconfirm title="确定删除该原因吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="父级原因" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="treeData"
            :props="{ label: 'name', value: 'id' }"
            placeholder="请选择父级原因（可不选）"
            clearable
            check-strictly
          />
        </el-form-item>
        <el-form-item label="原因名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入原因名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入原因描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
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
const dialogTitle = ref('添加原因')
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
    { required: true, message: '请输入原因名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
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
    ElMessage.error('获取原因树失败')
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
  dialogTitle.value = '添加原因'
  formData.id = undefined
  formData.name = ''
  formData.description = ''
  formData.sort_order = 0
  formData.status = 1
  formData.parent_id = undefined
  dialogVisible.value = true
}

const handleAddChild = (row: LinkCause) => {
  dialogTitle.value = '添加子原因'
  formData.id = undefined
  formData.name = ''
  formData.description = ''
  formData.sort_order = 0
  formData.status = 1
  formData.parent_id = row.id
  dialogVisible.value = true
}

const handleEdit = (row: LinkCause) => {
  dialogTitle.value = '编辑原因'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: LinkCause) => {
  try {
    await request.del({
      url: `/api/admin/link-causes/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchTreeData()
  } catch (error) {
    ElMessage.error('删除失败')
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

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchTreeData()
    } catch (error) {
      ElMessage.error('操作失败')
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