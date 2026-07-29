<template>
  <div class="plugins-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="插件名称/作者" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="已启用" value="enabled" />
            <el-option label="已禁用" value="disabled" />
            <el-option label="未安装" value="uninstalled" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="全部" clearable>
            <el-option v-for="cat in categories" :key="cat.value" :label="cat.label" :value="cat.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>插件管理</h3>
        <el-button type="primary" @click="handleRefreshPlugins">
          <el-icon><Refresh /></el-icon>刷新列表
        </el-button>
      </div>

      <el-table :data="plugins" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="插件名称" min-width="150">
          <template #default="{ row }">
            <div class="plugin-name-cell">
              <el-icon :size="20" class="plugin-icon"><component :is="row.icon || 'Box'" /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="author" label="作者" width="120" />
        <el-table-column prop="category" label="分类" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getCategoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" show-overflow-tooltip min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              active-value="enabled"
              inactive-value="disabled"
              :disabled="row.status === 'uninstalled'"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="installed_at" label="安装时间" width="170" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleViewDetail(row)">详情</el-button>
            <el-button type="primary" link size="small" @click="handleConfig(row)" v-if="row.status !== 'uninstalled' && row.has_config">
              配置
            </el-button>
            <el-button type="warning" link size="small" @click="handleUpdate(row)" v-if="row.has_update">
              更新
            </el-button>
            <el-button type="danger" link size="small" @click="handleUninstall(row)" v-if="row.status !== 'uninstalled'">
              卸载
            </el-button>
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
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="插件详情" width="600px">
      <div class="plugin-detail" v-if="currentPlugin">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="插件名称">{{ currentPlugin.name }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{ currentPlugin.version }}</el-descriptions-item>
          <el-descriptions-item label="作者">{{ currentPlugin.author }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ getCategoryLabel(currentPlugin.category) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentPlugin.status)">{{ getStatusLabel(currentPlugin.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="安装时间">{{ currentPlugin.installed_at || '-' }}</el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ currentPlugin.description || '-' }}</el-descriptions-item>
          <el-descriptions-item label="主页" :span="2">
            <el-link type="primary" :href="currentPlugin.homepage" target="_blank" v-if="currentPlugin.homepage">
              {{ currentPlugin.homepage }}
            </el-link>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 配置对话框 -->
    <el-dialog v-model="configVisible" title="插件配置" width="650px">
      <el-form :model="configForm" label-width="120px" v-if="configFields.length">
        <template v-for="field in configFields" :key="field.key">
          <el-form-item :label="field.label" :required="field.required">
            <el-input v-if="field.type === 'text'" v-model="configForm[field.key]" :placeholder="field.placeholder" />
            <el-input v-else-if="field.type === 'textarea'" v-model="configForm[field.key]" type="textarea" :rows="3" :placeholder="field.placeholder" />
            <el-switch v-else-if="field.type === 'boolean'" v-model="configForm[field.key]" />
            <el-input-number v-else-if="field.type === 'number'" v-model="configForm[field.key]" :min="field.min" :max="field.max" />
            <el-select v-else-if="field.type === 'select'" v-model="configForm[field.key]" :placeholder="field.placeholder">
              <el-option v-for="opt in field.options" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
          </el-form-item>
        </template>
      </el-form>
      <el-empty description="该插件暂无配置项" v-else />
      <template #footer>
        <el-button @click="configVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveConfig" :loading="configSaving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Box } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface Plugin {
  id: number
  name: string
  version: string
  author: string
  category: string
  description: string
  icon: string
  status: string
  homepage: string
  has_config: boolean
  has_update: boolean
  installed_at: string
}

const loading = ref(false)
const plugins = ref<Plugin[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const detailVisible = ref(false)
const configVisible = ref(false)
const configSaving = ref(false)
const currentPlugin = ref<Plugin | null>(null)
const configFields = ref<any[]>([])
const configForm = ref<Record<string, any>>({})

const categories = [
  { label: '支付', value: 'payment' },
  { label: '通知', value: 'notification' },
  { label: '安全', value: 'security' },
  { label: '工具', value: 'tool' },
  { label: '主题', value: 'theme' },
  { label: '集成', value: 'integration' }
]

const searchForm = ref({ keyword: '', status: '', category: '' })

const getCategoryLabel = (val: string) => categories.find(c => c.value === val)?.label || val
const getStatusType = (status: string) => {
  const map: Record<string, string> = { enabled: 'success', disabled: 'info', uninstalled: 'danger' }
  return map[status] || 'info'
}
const getStatusLabel = (status: string) => {
  const map: Record<string, string> = { enabled: '已启用', disabled: '已禁用', uninstalled: '未安装' }
  return map[status] || status
}

const fetchPlugins = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/plugins', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        ...searchForm.value
      }
    })
    plugins.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取插件列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchPlugins() }
const resetSearch = () => { searchForm.value = { keyword: '', status: '', category: '' }; handleSearch() }
const handleSizeChange = (val: number) => { pageSize.value = val; fetchPlugins() }
const handlePageChange = (val: number) => { currentPage.value = val; fetchPlugins() }

const handleToggleStatus = async (plugin: Plugin) => {
  try {
    await request.put(`/admin/api/v1/plugins/${plugin.id}/status`, { status: plugin.status })
    ElMessage.success(`插件已${plugin.status === 'enabled' ? '启用' : '禁用'}`)
  } catch {
    plugin.status = plugin.status === 'enabled' ? 'disabled' : 'enabled'
    ElMessage.error('操作失败')
  }
}

const handleViewDetail = (plugin: Plugin) => {
  currentPlugin.value = plugin
  detailVisible.value = true
}

const handleConfig = async (plugin: Plugin) => {
  currentPlugin.value = plugin
  try {
    const { data } = await request.get(`/admin/api/v1/plugins/${plugin.id}/config`)
    configFields.value = data.data?.fields || []
    configForm.value = data.data?.values || {}
    configVisible.value = true
  } catch {
    ElMessage.error('获取配置失败')
  }
}

const handleSaveConfig = async () => {
  configSaving.value = true
  try {
    await request.put(`/admin/api/v1/plugins/${currentPlugin.value?.id}/config`, configForm.value)
    ElMessage.success('配置已保存')
    configVisible.value = false
  } catch {
    ElMessage.error('保存配置失败')
  } finally {
    configSaving.value = false
  }
}

const handleUpdate = async (plugin: Plugin) => {
  try {
    await ElMessageBox.confirm(`确认更新插件「${plugin.name}」？`, '更新确认', { type: 'info' })
    loading.value = true
    await request.post(`/admin/api/v1/plugins/${plugin.id}/update`)
    ElMessage.success('插件更新成功')
    fetchPlugins()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('更新失败')
  } finally {
    loading.value = false
  }
}

const handleUninstall = async (plugin: Plugin) => {
  try {
    await ElMessageBox.confirm(`确认卸载插件「${plugin.name}」？卸载后数据将被清除。`, '卸载确认', {
      type: 'warning',
      confirmButtonText: '确认卸载'
    })
    loading.value = true
    await request.delete(`/admin/api/v1/plugins/${plugin.id}`)
    ElMessage.success('插件已卸载')
    fetchPlugins()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('卸载失败')
  } finally {
    loading.value = false
  }
}

const handleRefreshPlugins = () => { fetchPlugins() }

onMounted(() => { fetchPlugins() })
</script>

<style scoped lang="scss">
.plugins-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .plugin-name-cell {
    display: flex;
    align-items: center;
    gap: 8px;
    .plugin-icon { color: var(--primary-color); }
  }
  .plugin-detail {
    :deep(.el-descriptions) { margin-top: 8px; }
  }
}
</style>
