<template>
  <div class="server-groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('serverGroups.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('serverGroups.addGroup') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('serverGroups.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('serverGroups.keywordPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('serverGroups.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('serverGroups.all')" clearable>
            <el-option :label="$t('serverGroups.enabled')" :value="1" />
            <el-option :label="$t('serverGroups.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('serverGroups.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('serverGroups.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" :label="$t('serverGroups.id')" width="70" align="center" />
        <el-table-column prop="name" :label="$t('serverGroups.groupName')" min-width="180" />
        <el-table-column prop="code" :label="$t('serverGroups.groupCode')" width="150" />
        <el-table-column prop="description" :label="$t('serverGroups.description')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="server_count" :label="$t('serverGroups.serverCount')" width="110" align="center" />
        <el-table-column prop="load_balance_type" :label="$t('serverGroups.loadBalanceStrategy')" width="130" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ loadBalanceMap[row.load_balance_type as keyof typeof loadBalanceMap] || row.load_balance_type || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('serverGroups.sort')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('serverGroups.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? $t('serverGroups.enabled') : $t('serverGroups.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('serverGroups.createTime')" width="180" />
        <el-table-column :label="$t('serverGroups.operations')" width="250" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleManageServers(row)">{{ $t('serverGroups.manageServers') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('serverGroups.edit') }}</el-button>
            <el-popconfirm :title="$t('serverGroups.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('serverGroups.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item :label="$t('serverGroups.groupName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('serverGroups.inputGroupName')" />
        </el-form-item>
        <el-form-item :label="$t('serverGroups.groupCode')" prop="code">
          <el-input v-model="formData.code" :placeholder="$t('serverGroups.inputGroupCode')" />
        </el-form-item>
        <el-form-item :label="$t('serverGroups.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('serverGroups.inputDescription')" />
        </el-form-item>
        <el-form-item :label="$t('serverGroups.loadBalanceStrategy')" prop="load_balance_type">
          <el-select v-model="formData.load_balance_type" :placeholder="$t('serverGroups.selectLoadBalance')" style="width: 100%">
            <el-option :label="$t('serverGroups.roundRobin')" value="round_robin" />
            <el-option :label="$t('serverGroups.weightedRoundRobin')" value="weighted_round_robin" />
            <el-option :label="$t('serverGroups.leastConnections')" value="least_connections" />
            <el-option :label="$t('serverGroups.ipHash')" value="ip_hash" />
            <el-option :label="$t('serverGroups.random')" value="random" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('serverGroups.sort')" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('serverGroups.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" :active-text="$t('serverGroups.enabled')" :inactive-text="$t('serverGroups.disabled')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('serverGroups.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('serverGroups.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 管理服务器对话框 -->
    <el-dialog v-model="serverDialogVisible" :title="$t('serverGroups.manageRelatedServers')" width="750px" destroy-on-close>
      <div class="server-transfer">
        <el-transfer
          v-model="selectedServerIds"
          :data="allServers"
          :titles="[$t('serverGroups.optionalServers'), $t('serverGroups.linkedServers')]"
          :props="{ key: 'id', label: 'name' }"
          filterable
          :filter-placeholder="$t('serverGroups.searchServer')"
        />
      </div>
      <template #footer>
        <el-button @click="serverDialogVisible = false">{{ $t('serverGroups.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveServers" :loading="serverSaving">{{ $t('serverGroups.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

defineOptions({ name: 'ServerGroups' })

const loadBalanceMap = computed(() => ({
  round_robin: $t('serverGroups.roundRobin'),
  weighted_round_robin: $t('serverGroups.weightedRoundRobin'),
  least_connections: $t('serverGroups.leastConnections'),
  ip_hash: $t('serverGroups.ipHash'),
  random: $t('serverGroups.random')
}))

const loading = ref(false)
const submitLoading = ref(false)
const serverSaving = ref(false)
const dialogVisible = ref(false)
const serverDialogVisible = ref(false)
const formRef = ref<FormInstance>()
const currentGroupId = ref<number>(0)

const dialogTitle = computed(() => formData.id ? $t('serverGroups.editGroupTitle') : $t('serverGroups.addGroupTitle'))

const searchForm = reactive({ keyword: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const allServers = ref<any[]>([])
const selectedServerIds = ref<number[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  code: '',
  description: '',
  load_balance_type: 'round_robin',
  sort: 0,
  status: 1
})

const formRules = computed<FormRules>(() => ({
  name: [
    { required: true, message: $t('serverGroups.inputGroupName'), trigger: 'blur' },
    { min: 2, max: 50, message: $t('serverGroups.groupNameLength'), trigger: 'blur' }
  ],
  code: [
    { required: true, message: $t('serverGroups.inputGroupCode'), trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: $t('serverGroups.codePattern'), trigger: 'blur' }
  ]
}))

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/server-groups',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('serverGroups.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const fetchAllServers = async () => {
  try {
    const data = await request.get({ url: '/api/admin/config-servers/server-list', params: { page_size: 999 } })
    allServers.value = data.list || data || []
  } catch (error) {
    console.error('获取服务器列表失败:', error)
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = undefined; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleAdd = () => {
  formData.id = undefined; formData.name = ''; formData.code = ''
  formData.description = ''; formData.load_balance_type = 'round_robin'
  formData.sort = 0; formData.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/server-groups/${row.id}` })
    ElMessage.success($t('serverGroups.deleteSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('serverGroups.deleteFailed'))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/server-groups/${formData.id}`, params: formData })
      } else {
        await request.post({ url: '/api/admin/server-groups', params: formData })
      }
      ElMessage.success(formData.id ? $t('serverGroups.updateSuccess') : $t('serverGroups.addSuccess'))
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('serverGroups.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

const handleManageServers = async (row: any) => {
  currentGroupId.value = row.id
  try {
    const data = await request.get({ url: `/api/admin/server-groups/${row.id}/servers` })
    selectedServerIds.value = (data.list || data || []).map((s: any) => s.id || s)
    serverDialogVisible.value = true
  } catch (error) {
    ElMessage.error($t('serverGroups.fetchLinkedFailed'))
  }
}

const handleSaveServers = async () => {
  serverSaving.value = true
  try {
    await request.put({
      url: `/api/admin/server-groups/${currentGroupId.value}/servers`,
      params: { server_ids: selectedServerIds.value }
    })
    ElMessage.success($t('serverGroups.linkSuccess'))
    serverDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error($t('serverGroups.saveFailed'))
  } finally {
    serverSaving.value = false
  }
}

onMounted(() => { fetchData(); fetchAllServers() })
</script>

<style scoped lang="scss">
.server-groups-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.server-transfer { display: flex; justify-content: center; }
</style>
