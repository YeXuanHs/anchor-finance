<template>
  <div class="rules-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>规则管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增规则
          </el-button>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border row-key="id" :tree-props="{ children: 'son' }">
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="title" label="规则标题" min-width="200" />
        <el-table-column prop="name" label="规则名称" min-width="200" />
        <el-table-column prop="url" label="URL" min-width="200" />
        <el-table-column prop="cn_name" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_display ? 'success' : 'info'" size="small">
              {{ row.cn_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该规则吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="规则标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入规则标题" />
        </el-form-item>
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="formData.name" placeholder="如：ADMIN_INDEX_INDEX" />
        </el-form-item>
        <el-form-item label="URL" prop="url">
          <el-input v-model="formData.url" placeholder="如：index/index" />
        </el-form-item>
        <el-form-item label="类型" prop="is_display">
          <el-radio-group v-model="formData.is_display">
            <el-radio :value="1">前台页面</el-radio>
            <el-radio :value="0">接口</el-radio>
          </el-radio-group>
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

const loading = ref(false)
const submitLoading = ref(false)

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('新增规则')
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  title: '',
  name: '',
  url: '',
  is_display: 0
})

const formRules: FormRules = {
  title: [{ required: true, message: '请输入规则标题', trigger: 'blur' }],
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入URL', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/rules' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取规则列表失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增规则'
  formData.id = undefined
  formData.title = ''
  formData.name = ''
  formData.url = ''
  formData.is_display = 0
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑规则'
  formData.id = row.id
  formData.title = row.title
  formData.name = row.name
  formData.url = row.url
  formData.is_display = row.is_display
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/rules/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
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
        await request.put({ url: `/api/admin/rules/${formData.id}`, params: { ...formData } })
      } else {
        await request.post({ url: '/api/admin/rules', params: { ...formData } })
      }
      ElMessage.success(formData.id ? '编辑成功' : '添加成功')
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.rules-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
