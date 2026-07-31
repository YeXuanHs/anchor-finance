<template>
  <div class="product-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>产品分组</span>
          <el-button type="primary" @click="handleAdd()">
            <el-icon><Plus /></el-icon>
            添加分组
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="分组名称">
          <el-input v-model="searchForm.name" placeholder="请输入分组名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 树形表格 -->
      <el-table
        :data="filteredTableData"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        default-expand-all
      >
        <el-table-column prop="name" label="分组名称" min-width="200" />
        <el-table-column prop="code" label="分组编码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="product_count" label="产品数量" width="100" align="center" />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right" align="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleAdd(row)">添加子分组</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该分组吗？删除后子分组也会被删除。" @confirm="handleDelete(row)">
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
        <el-form-item label="上级分组" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="groupTreeData"
            :props="{ value: 'id', label: 'name', children: 'children' }"
            placeholder="无（顶级分组）"
            clearable
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="分组编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入分组编码" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入分组描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'ProductGroups' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  name: ''
})

// 表格数据
const tableData = ref([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加分组')
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: undefined as number | undefined,
  parent_id: undefined as number | undefined,
  name: '',
  code: '',
  description: '',
  sort: 0,
  status: 1
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入分组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入分组编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '只能包含字母、数字和下划线', trigger: 'blur' }
  ]
}

// 过滤后的表格数据
const filteredTableData = computed(() => {
  if (!searchForm.name) {
    return tableData.value
  }
  return filterTree(tableData.value, searchForm.name.toLowerCase())
})

// 递归过滤树形数据
const filterTree = (data: any[], keyword: string): any[] => {
  return data.filter(item => {
    const nameMatch = item.name.toLowerCase().includes(keyword)
    const childrenMatch = item.children && filterTree(item.children, keyword).length > 0
    return nameMatch || childrenMatch
  }).map(item => {
    if (item.children) {
      return {
        ...item,
        children: filterTree(item.children, keyword)
      }
    }
    return item
  })
}

// 分组树形数据（用于选择器）
const groupTreeData = computed(() => {
  return buildTreeSelectData(tableData.value, formData.id)
})

// 构建树形选择器数据（排除当前编辑的节点及其子节点）
const buildTreeSelectData = (data: any[], excludeId?: number): any[] => {
  return data
    .filter(item => item.id !== excludeId)
    .map(item => ({
      id: item.id,
      name: item.name,
      children: item.children ? buildTreeSelectData(item.children, excludeId) : []
    }))
}

// 获取分组列表
const fetchGroups = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/product-groups'
    })
    tableData.value = data || []
  } catch (error) {
    console.error('获取分组列表失败:', error)
    ElMessage.error('获取分组列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  // 前端过滤，无需重新请求
}

// 重置
const handleReset = () => {
  searchForm.name = ''
}

// 添加
const handleAdd = (parent?: any) => {
  dialogTitle.value = parent ? `添加子分组 - ${parent.name}` : '添加分组'
  formData.id = undefined
  formData.parent_id = parent?.id || undefined
  formData.name = ''
  formData.code = ''
  formData.description = ''
  formData.sort = 0
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑分组'
  formData.id = row.id
  formData.parent_id = row.parent_id
  formData.name = row.name
  formData.code = row.code
  formData.description = row.description
  formData.sort = row.sort
  formData.status = row.status
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/product-groups/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchGroups()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const url = formData.id ? `/api/admin/product-groups/${formData.id}` : '/api/admin/product-groups'

      if (formData.id) {
        await request.put({ url, params: formData })
      } else {
        await request.post({ url, params: formData })
      }

      ElMessage.success(formData.id ? '更新成功' : '添加成功')
      dialogVisible.value = false
      fetchGroups()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => {
  fetchGroups()
})
</script>

<style scoped lang="scss">
.product-groups-page {
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
