<template>
  <div class="data-migration-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>数据迁移</span>
        </div>
      </template>

      <!-- 迁移任务列表 -->
      <div class="section">
        <div class="section-header">
          <h3>迁移任务</h3>
          <el-button type="primary" size="small" @click="handleCreateTask">创建迁移任务</el-button>
        </div>
        <el-table :data="tasks" v-loading="loading" style="width: 100%" border>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="任务名称" min-width="150" />
          <el-table-column prop="source_type" label="数据源" width="120" />
          <el-table-column prop="target_type" label="目标" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="statusTypeMap[row.status]" size="small">
                {{ statusLabelMap[row.status] }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="progress" label="进度" width="200">
            <template #default="{ row }">
              <el-progress :percentage="row.progress" :status="row.status === 'completed' ? 'success' : undefined" />
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button v-if="row.status === 'pending'" type="success" link @click="handleStart(row)">开始</el-button>
              <el-button v-if="row.status === 'running'" type="warning" link @click="handlePause(row)">暂停</el-button>
              <el-button type="primary" link @click="handleViewLog(row)">日志</el-button>
              <el-popconfirm v-if="row.status !== 'running'" title="确定删除吗？" @confirm="handleDelete(row)">
                <template #reference>
                  <el-button type="danger" link>删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 创建任务对话框 -->
    <el-dialog v-model="dialogVisible" title="创建迁移任务" width="600px">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="任务名称" required>
          <el-input v-model="formData.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="数据源类型" required>
          <el-select v-model="formData.source_type" placeholder="请选择">
            <el-option label="WHMCS" value="whmcs" />
            <el-option label="Blesta" value="blesta" />
            <el-option label="CSV文件" value="csv" />
            <el-option label="数据库" value="database" />
          </el-select>
        </el-form-item>
        <el-form-item label="数据源配置">
          <el-input v-model="formData.source_config" type="textarea" :rows="4" placeholder="JSON格式配置" />
        </el-form-item>
        <el-form-item label="迁移类型">
          <el-checkbox-group v-model="formData.migration_types">
            <el-checkbox value="users">用户</el-checkbox>
            <el-checkbox value="products">产品</el-checkbox>
            <el-checkbox value="orders">订单</el-checkbox>
            <el-checkbox value="tickets">工单</el-checkbox>
            <el-checkbox value="invoices">账单</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">创建</el-button>
      </template>
    </el-dialog>

    <!-- 日志查看对话框 -->
    <el-dialog v-model="logVisible" title="迁移日志" width="700px">
      <pre style="max-height: 500px; overflow-y: auto; background: #f5f7fa; padding: 16px; border-radius: 8px; font-size: 13px; line-height: 1.6; white-space: pre-wrap; word-break: break-all;">{{ logData }}</pre>
      <template #footer>
        <el-button @click="logVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<!-- TODO: 完善数据迁移页面功能
  1. 实现实际的迁移逻辑
  2. 添加迁移进度实时更新
  3. 添加迁移日志查看
  4. 添加迁移回滚功能
  5. 添加数据映射配置
-->
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const statusTypeMap: Record<string, string> = {
  pending: 'info',
  running: 'warning',
  completed: 'success',
  failed: 'danger',
  paused: 'info'
}

const statusLabelMap: Record<string, string> = {
  pending: '待执行',
  running: '执行中',
  completed: '已完成',
  failed: '失败',
  paused: '已暂停'
}

const loading = ref(false)
const dialogVisible = ref(false)
const tasks = ref<any[]>([])

const formData = reactive({
  name: '',
  source_type: 'whmcs',
  source_config: '',
  migration_types: ['users'] as string[]
})

const fetchTasks = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/system/data-migration/tasks' })
    tasks.value = data || []
  } catch (error) {
    ElMessage.error('获取迁移任务失败')
  } finally {
    loading.value = false
  }
}

const handleCreateTask = () => {
  formData.name = ''
  formData.source_type = 'whmcs'
  formData.source_config = ''
  formData.migration_types = ['users']
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formData.name) { ElMessage.warning('请输入任务名称'); return }
  try {
    await request.post({ url: '/api/admin/system/data-migration/tasks', params: { ...formData } })
    ElMessage.success('任务创建成功')
    dialogVisible.value = false
    fetchTasks()
  } catch (error) {
    ElMessage.error('创建失败')
  }
}

const handleStart = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/system/data-migration/tasks/${row.id}/start` })
    ElMessage.success('任务已开始')
    fetchTasks()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handlePause = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/system/data-migration/tasks/${row.id}/pause` })
    ElMessage.success('任务已暂停')
    fetchTasks()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleViewLog = (row: any) => {
  logData.value = row.logs || '暂无日志数据'
  logVisible.value = true
}

const logVisible = ref(false)
const logData = ref('')

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/system/data-migration/tasks/${row.id}` })
    ElMessage.success('删除成功')
    fetchTasks()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => { fetchTasks() })
</script>

<style scoped lang="scss">
.data-migration-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.section { margin-top: 24px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.section-header h3 { margin: 0; font-size: 16px; font-weight: 600; }
</style>
