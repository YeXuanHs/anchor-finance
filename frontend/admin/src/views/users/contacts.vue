<template>
  <div class="contacts-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="姓名/邮箱/手机" clearable />
        </el-form-item>
        <el-form-item label="用户">
          <el-input v-model="searchForm.username" placeholder="关联用户名" clearable />
        </el-form-item>
        <el-form-item label="验证状态">
          <el-select v-model="searchForm.verified" placeholder="全部" clearable>
            <el-option label="已验证" :value="true" />
            <el-option label="未验证" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>联系方式列表</h3>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加联系方式
        </el-button>
      </div>

      <el-table :data="contacts" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user.username" label="关联用户" width="120" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
        <el-table-column prop="phone" label="手机" width="130" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_primary" label="主联系人" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_primary ? 'success' : 'info'" size="small">
              {{ row.is_primary ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="email_verified" label="邮箱验证" width="100">
          <template #default="{ row }">
            <el-tag :type="row.email_verified ? 'success' : 'warning'" size="small">
              {{ row.email_verified ? '已验证' : '未验证' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone_verified" label="手机验证" width="100">
          <template #default="{ row }">
            <el-tag :type="row.phone_verified ? 'success' : 'warning'" size="small">
              {{ row.phone_verified ? '已验证' : '未验证' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editContact(row)">编辑</el-button>
            <el-button type="primary" link @click="sendVerify(row)" v-if="!row.email_verified && row.email">发验证</el-button>
            <el-button type="danger" link @click="deleteContact(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchContacts"
          @current-change="fetchContacts"
        />
      </div>
    </div>

    <el-dialog v-model="showFormDialog" :title="editingItem ? '编辑联系方式' : '添加联系方式'" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="关联用户" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" :disabled="!!editingItem" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="formData.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="formData.type" placeholder="请选择类型">
            <el-option label="管理员" value="admin" />
            <el-option label="技术" value="tech" />
            <el-option label="财务" value="finance" />
            <el-option label="销售" value="sales" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="主联系人">
          <el-switch v-model="formData.is_primary" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFormDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const contacts = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showFormDialog = ref(false)
const editingItem = ref<any>(null)
const formRef = ref<FormInstance>()

const searchForm = ref({ keyword: '', username: '', verified: '' })
const formData = ref({ username: '', name: '', email: '', phone: '', type: 'other', is_primary: false })
const formRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = { admin: '管理员', tech: '技术', finance: '财务', sales: '销售', other: '其他' }
  return map[type] || type
}

const fetchContacts = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/users/contacts', { params })
    contacts.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取联系方式列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchContacts() }
const resetSearch = () => { searchForm.value = { keyword: '', username: '', verified: '' }; handleSearch() }

const openAddDialog = () => {
  editingItem.value = null
  formData.value = { username: '', name: '', email: '', phone: '', type: 'other', is_primary: false }
  showFormDialog.value = true
}

const editContact = (contact: any) => {
  editingItem.value = contact
  formData.value = { username: contact.user?.username || '', name: contact.name, email: contact.email || '', phone: contact.phone || '', type: contact.type || 'other', is_primary: contact.is_primary || false }
  showFormDialog.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    if (editingItem.value) {
      await request.put(`/admin/api/v1/users/contacts/${editingItem.value.id}`, formData.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/admin/api/v1/users/contacts', formData.value)
      ElMessage.success('添加成功')
    }
    showFormDialog.value = false
    fetchContacts()
  } catch {
    ElMessage.error(editingItem.value ? '更新失败' : '添加失败')
  } finally {
    submitLoading.value = false
  }
}

const deleteContact = async (contact: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该联系方式吗？', '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/users/contacts/${contact.id}`)
    ElMessage.success('删除成功')
    fetchContacts()
  } catch {}
}

const sendVerify = async (contact: any) => {
  try {
    await request.post(`/admin/api/v1/users/contacts/${contact.id}/send-verify`)
    ElMessage.success('验证邮件已发送')
  } catch {
    ElMessage.error('发送失败')
  }
}

onMounted(fetchContacts)
</script>

<style scoped lang="scss">
.contacts-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
