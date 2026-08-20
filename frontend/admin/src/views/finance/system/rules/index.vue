<template>
  <div class="rules-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('ruleManage.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('ruleManage.addRule') }}
          </el-button>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border row-key="id" :tree-props="{ children: 'son' }">
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="title" :label="$t('ruleManage.ruleTitle')" min-width="200" />
        <el-table-column prop="name" :label="$t('ruleManage.ruleName')" min-width="200" />
        <el-table-column prop="url" label="URL" min-width="200" />
        <el-table-column prop="cn_name" :label="$t('ruleManage.type')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_display ? 'success' : 'info'" size="small">
              {{ row.cn_name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('ruleManage.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('ruleManage.enabled') : $t('ruleManage.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('ruleManage.operation')" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('ruleManage.edit') }}</el-button>
            <el-popconfirm :title="$t('ruleManage.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('ruleManage.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('ruleManage.ruleTitle')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('ruleManage.enterRuleTitle')" />
        </el-form-item>
        <el-form-item :label="$t('ruleManage.ruleName')" prop="name">
          <el-input v-model="formData.name" placeholder="ADMIN_INDEX_INDEX" />
        </el-form-item>
        <el-form-item label="URL" prop="url">
          <el-input v-model="formData.url" placeholder="index/index" />
        </el-form-item>
        <el-form-item :label="$t('ruleManage.type')" prop="is_display">
          <el-radio-group v-model="formData.is_display">
            <el-radio :value="1">{{ $t('ruleManage.frontendPage') }}</el-radio>
            <el-radio :value="0">{{ $t('ruleManage.api') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('ruleManage.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('ruleManage.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const loading = ref(false)
const submitLoading = ref(false)

const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref($t('ruleManage.addRule'))
const formRef = ref<FormInstance>()

const formData = reactive({
  id: undefined as number | undefined,
  title: '',
  name: '',
  url: '',
  is_display: 0
})

const formRules: FormRules = {
  title: [{ required: true, message: $t('ruleManage.enterRuleTitle'), trigger: 'blur' }],
  name: [{ required: true, message: $t('ruleManage.enterRuleName'), trigger: 'blur' }],
  url: [{ required: true, message: $t('ruleManage.enterUrl'), trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/rules' })
    tableData.value = data || []
  } catch (error) {
    console.error('获取规则列表失败:', error)
    ElMessage.error($t('ruleManage.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = $t('ruleManage.addRule')
  formData.id = undefined
  formData.title = ''
  formData.name = ''
  formData.url = ''
  formData.is_display = 0
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = $t('ruleManage.editRule')
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
    ElMessage.success($t('ruleManage.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('ruleManage.deleteFailed'))
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
      ElMessage.success(formData.id ? $t('ruleManage.editSuccess') : $t('ruleManage.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('ruleManage.operationFailed'))
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
