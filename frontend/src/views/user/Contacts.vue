<template>
  <div class="contacts-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>联系人管理</span>
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加联系人
          </el-button>
        </div>
      </template>

      <el-table :data="contacts" style="width: 100%">
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="电话" />
        <el-table-column prop="company" label="公司" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="editContact(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteContact(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="showAddDialog" :title="editingContact ? '编辑联系人' : '添加联系人'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="公司" prop="company">
          <el-input v-model="form.company" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const contacts = ref([])
const showAddDialog = ref(false)
const editingContact = ref<any>(null)
const formRef = ref()

const form = ref({
  id: 0,
  name: '',
  email: '',
  phone: '',
  company: '',
  remark: ''
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }]
}

const editContact = (contact: any) => {
  editingContact.value = contact
  form.value = { ...contact }
  showAddDialog.value = true
}

const deleteContact = async (contact: any) => {
  await ElMessageBox.confirm('确定删除该联系人？', '提示', { type: 'warning' })
    await request.delete(`/api/v1/contacts/${contact.id}`)
  ElMessage.success('删除成功')
}

const submitForm = async () => {
  if (editingContact.value) {
    await request.put(`/api/v1/contacts/${form.value.id}`, form.value)
  } else {
    await request.post('/api/v1/contacts', form.value)
  }
  ElMessage.success('保存成功')
  showAddDialog.value = false
}
</script>

<style scoped lang="scss">
.contacts-page {
  .page-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style>
