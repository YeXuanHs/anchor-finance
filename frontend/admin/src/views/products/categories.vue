<template>
  <div class="categories-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>产品分类</h3>
        <el-button type="primary" @click="showDialog = true; resetForm()">
          <el-icon><Plus /></el-icon>添加分类
        </el-button>
      </div>
      <el-table :data="treeData" v-loading="loading" row-key="id" default-expand-all :tree-props="{ children: 'children' }">
        <el-table-column prop="name" label="分类名称" min-width="200" />
        <el-table-column prop="code" label="标识" width="120" />
        <el-table-column prop="icon" label="图标" width="80">
          <template #default="{ row }">
            <img v-if="row.icon" :src="row.icon" style="height: 24px" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="product_count" label="产品数" width="100" />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editRow(row)">编辑</el-button>
            <el-button link type="primary" @click="addChild(row)">添加子分类</el-button>
            <el-button link type="danger" @click="deleteRow(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showDialog" :title="form.id ? '编辑分类' : '添加分类'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="分类名称"><el-input v-model="form.name" placeholder="请输入分类名称" /></el-form-item>
        <el-form-item label="标识"><el-input v-model="form.code" placeholder="英文标识" :disabled="!!form.id" /></el-form-item>
        <el-form-item label="上级分类">
          <el-tree-select
            v-model="form.parent_id"
            :data="treeData"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="无（顶级分类）"
            clearable
            check-strictly
          />
        </el-form-item>
        <el-form-item label="图标URL"><el-input v-model="form.icon" placeholder="图标URL" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const list = ref<any[]>([])
const showDialog = ref(false)
const form = ref<any>({ id: 0, name: '', code: '', parent_id: null, icon: '', description: '', sort_order: 0, status: 1 })

const resetForm = () => { form.value = { id: 0, name: '', code: '', parent_id: null, icon: '', description: '', sort_order: 0, status: 1 } }

const buildTree = (items: any[], parentId: number | null = null): any[] => {
  return items
    .filter(item => item.parent_id === parentId)
    .map(item => ({ ...item, children: buildTree(items, item.id) }))
    .sort((a, b) => a.sort_order - b.sort_order)
}

const treeData = computed(() => buildTree(list.value))

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/product-categories', { params: { page_size: 999 } })
    list.value = data.data || []
  } catch {} finally { loading.value = false }
}

const editRow = (row: any) => { form.value = { ...row }; showDialog.value = true }

const addChild = (row: any) => {
  resetForm()
  form.value.parent_id = row.id
  showDialog.value = true
}

const saveForm = async () => {
  try {
    if (form.value.id) { await request.put(`/admin/api/v1/product-categories/${form.value.id}`, form.value) }
    else { await request.post('/admin/api/v1/product-categories', form.value) }
    ElMessage.success('保存成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('保存失败') }
}

const deleteRow = async (id: number) => {
  await ElMessageBox.confirm('确定删除该分类？其子分类将一并删除。', '确认')
  try { await request.delete(`/admin/api/v1/product-categories/${id}`); ElMessage.success('删除成功'); fetchData() } catch { ElMessage.error('删除失败') }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
</style>
