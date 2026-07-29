<template>
  <div class="menu-page page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>菜单管理</h3>
        <el-button type="primary" @click="showAddDialog = true">
          <el-icon><Plus /></el-icon>
          添加菜单
        </el-button>
      </div>

      <el-table :data="menus" style="width: 100%" v-loading="loading" row-key="id" default-expand-all>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="菜单名称" />
        <el-table-column prop="icon" label="图标" width="100" />
        <el-table-column prop="path" label="路由路径" />
        <el-table-column prop="permission" label="权限标识" />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="visible" label="显示">
          <template #default="{ row }">
            <el-tag :type="row.visible ? 'success' : 'info'" size="small">
              {{ row.visible ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editMenu(row)">编辑</el-button>
            <el-button type="primary" link @click="addChild(row)">添加子菜单</el-button>
            <el-button type="danger" link @click="deleteMenu(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showAddDialog" :title="editingMenu ? '编辑菜单' : '添加菜单'" width="500px">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="上级菜单">
          <el-tree-select v-model="formData.parent_id" :data="menuTree" placeholder="无（顶级菜单）" clearable check-strictly />
        </el-form-item>
        <el-form-item label="菜单名称">
          <el-input v-model="formData.title" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="formData.icon" placeholder="图标名称" />
        </el-form-item>
        <el-form-item label="路由路径">
          <el-input v-model="formData.path" placeholder="如: /users" />
        </el-form-item>
        <el-form-item label="权限标识">
          <el-input v-model="formData.permission" placeholder="如: users.index" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="是否显示">
          <el-switch v-model="formData.visible" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const menus = ref([])
const menuTree = ref([])
const showAddDialog = ref(false)
const editingMenu = ref(null)

const formData = ref({ parent_id: null, title: '', icon: '', path: '', permission: '', sort_order: 0, visible: true })

const buildMenuTree = (list: any[]) => {
  const map: Record<number, any> = {}
  const roots: any[] = []
  list.forEach(item => { map[item.id] = { ...item, children: [] } })
  list.forEach(item => {
    if (item.parent_id && map[item.parent_id]) {
      map[item.parent_id].children.push(map[item.id])
    } else {
      roots.push(map[item.id])
    }
  })
  return roots
}

const fetchMenus = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/menus')
    if (data?.data) {
      menus.value = data.data
      menuTree.value = buildMenuTree(data.data)
    }
  } catch {
    ElMessage.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

const editMenu = (menu: any) => {
  editingMenu.value = menu
  formData.value = { ...menu }
  showAddDialog.value = true
}

const addChild = (menu: any) => {
  editingMenu.value = null
  formData.value = { parent_id: menu.id, title: '', icon: '', path: '', permission: '', sort_order: 0, visible: true }
  showAddDialog.value = true
}

const deleteMenu = async (menu: any) => {
  try {
    await ElMessageBox.confirm(`确认删除菜单「${menu.title}」？`, '提示', { type: 'warning' })
    await request.delete(`/admin/menus/${menu.id}`)
    ElMessage.success('删除成功')
    fetchMenus()
  } catch {}
}

const handleSubmit = async () => {
  try {
    if (editingMenu.value) {
      await request.put(`/admin/menus/${editingMenu.value.id}`, formData.value)
      ElMessage.success('编辑成功')
    } else {
      await request.post('/admin/menus', formData.value)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    editingMenu.value = null
    formData.value = { parent_id: null, title: '', icon: '', path: '', permission: '', sort_order: 0, visible: true }
    fetchMenus()
  } catch {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  fetchMenus()
})
</script>

<style scoped lang="scss">
.menu-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
}
</style>
