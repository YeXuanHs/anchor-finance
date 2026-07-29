<template>
  <div class="plugin-store-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="搜索插件" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="全部" clearable>
            <el-option v-for="cat in categories" :key="cat.value" :label="cat.label" :value="cat.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-select v-model="searchForm.sort" placeholder="默认">
            <el-option label="最新发布" value="newest" />
            <el-option label="最多安装" value="popular" />
            <el-option label="评分最高" value="rating" />
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

    <!-- 分类标签 -->
    <div class="category-tabs">
      <el-tag
        v-for="cat in [{ label: '全部', value: '' }, ...categories]"
        :key="cat.value"
        :type="searchForm.category === cat.value ? '' : 'info'"
        :effect="searchForm.category === cat.value ? 'dark' : 'plain'"
        class="category-tag"
        @click="searchForm.category = cat.value; handleSearch()"
      >
        {{ cat.label }}
      </el-tag>
    </div>

    <div v-loading="loading">
      <div class="plugin-grid" v-if="plugins.length">
        <div class="plugin-card" v-for="plugin in plugins" :key="plugin.id">
          <div class="plugin-icon-wrapper">
            <el-icon :size="40"><component :is="plugin.icon || 'Box'" /></el-icon>
          </div>
          <div class="plugin-info">
            <div class="plugin-header">
              <h4>{{ plugin.name }}</h4>
              <el-tag size="small" v-if="plugin.is_official">官方</el-tag>
            </div>
            <p class="plugin-desc">{{ plugin.description }}</p>
            <div class="plugin-meta">
              <span><el-icon><User /></el-icon>{{ plugin.author }}</span>
              <span><el-icon><Download /></el-icon>{{ formatCount(plugin.downloads) }}</span>
              <span class="rating"><el-icon><Star /></el-icon>{{ plugin.rating?.toFixed(1) || '0.0' }}</span>
            </div>
            <div class="plugin-footer">
              <span class="version">v{{ plugin.version }}</span>
              <div class="actions">
                <el-button type="primary" size="small" @click="handleInstall(plugin)" v-if="!plugin.installed" :loading="plugin._installing">
                  安装
                </el-button>
                <el-button type="success" size="small" disabled v-else-if="plugin.installed && !plugin.has_update">
                  已安装
                </el-button>
                <el-button type="warning" size="small" @click="handleUpdate(plugin)" v-else-if="plugin.has_update" :loading="plugin._updating">
                  更新
                </el-button>
                <el-button size="small" @click="handleViewDetail(plugin)">详情</el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <el-empty description="暂无插件" v-else-if="!loading" />
    </div>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="插件详情" width="700px">
      <div class="store-detail" v-if="currentPlugin">
        <div class="detail-header">
          <el-icon :size="56"><component :is="currentPlugin.icon || 'Box'" /></el-icon>
          <div class="detail-title">
            <h3>{{ currentPlugin.name }} <el-tag size="small">v{{ currentPlugin.version }}</el-tag></h3>
            <p>作者：{{ currentPlugin.author }}</p>
          </div>
        </div>
        <el-divider />
        <el-descriptions :column="2" border>
          <el-descriptions-item label="分类">{{ getCategoryLabel(currentPlugin.category) }}</el-descriptions-item>
          <el-descriptions-item label="下载量">{{ formatCount(currentPlugin.downloads) }}</el-descriptions-item>
          <el-descriptions-item label="评分">{{ currentPlugin.rating?.toFixed(1) || '-' }}</el-descriptions-item>
          <el-descriptions-item label="大小">{{ formatSize(currentPlugin.size) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ currentPlugin.updated_at || '-' }}</el-descriptions-item>
          <el-descriptions-item label="兼容版本">{{ currentPlugin.compat_version || '-' }}</el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ currentPlugin.description || '-' }}</el-descriptions-item>
          <el-descriptions-item label="主页" :span="2">
            <el-link type="primary" :href="currentPlugin.homepage" target="_blank" v-if="currentPlugin.homepage">
              {{ currentPlugin.homepage }}
            </el-link>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleInstall(currentPlugin)" v-if="currentPlugin && !currentPlugin.installed">
          立即安装
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, User, Download, Star, Box } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface StorePlugin {
  id: number
  name: string
  version: string
  author: string
  category: string
  description: string
  icon: string
  homepage: string
  downloads: number
  rating: number
  size: number
  is_official: boolean
  installed: boolean
  has_update: boolean
  updated_at: string
  compat_version: string
  _installing?: boolean
  _updating?: boolean
}

const loading = ref(false)
const plugins = ref<StorePlugin[]>([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const detailVisible = ref(false)
const currentPlugin = ref<StorePlugin | null>(null)

const categories = [
  { label: '支付', value: 'payment' },
  { label: '通知', value: 'notification' },
  { label: '安全', value: 'security' },
  { label: '工具', value: 'tool' },
  { label: '主题', value: 'theme' },
  { label: '集成', value: 'integration' }
]

const searchForm = ref({ keyword: '', category: '', sort: '' })

const getCategoryLabel = (val: string) => categories.find(c => c.value === val)?.label || val

const formatCount = (num: number) => {
  if (!num) return '0'
  if (num >= 10000) return (num / 10000).toFixed(1) + 'w'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
  return String(num)
}

const formatSize = (bytes: number) => {
  if (!bytes) return '-'
  if (bytes >= 1048576) return (bytes / 1048576).toFixed(1) + ' MB'
  return (bytes / 1024).toFixed(0) + ' KB'
}

const fetchPlugins = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/plugins/store', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        ...searchForm.value
      }
    })
    plugins.value = (data.data?.list || []).map((p: StorePlugin) => ({ ...p, _installing: false, _updating: false }))
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取插件商店失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchPlugins() }
const resetSearch = () => { searchForm.value = { keyword: '', category: '', sort: '' }; handleSearch() }
const handleSizeChange = (val: number) => { pageSize.value = val; fetchPlugins() }
const handlePageChange = (val: number) => { currentPage.value = val; fetchPlugins() }

const handleInstall = async (plugin: StorePlugin | null) => {
  if (!plugin) return
  try {
    await ElMessageBox.confirm(`确认安装插件「${plugin.name}」？`, '安装确认', { type: 'info' })
    plugin._installing = true
    await request.post(`/admin/api/v1/plugins/store/${plugin.id}/install`)
    ElMessage.success('安装成功')
    plugin.installed = true
    plugin.has_update = false
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('安装失败')
  } finally {
    plugin._installing = false
  }
}

const handleUpdate = async (plugin: StorePlugin) => {
  try {
    plugin._updating = true
    await request.post(`/admin/api/v1/plugins/store/${plugin.id}/update`)
    ElMessage.success('更新成功')
    plugin.has_update = false
  } catch {
    ElMessage.error('更新失败')
  } finally {
    plugin._updating = false
  }
}

const handleViewDetail = (plugin: StorePlugin) => {
  currentPlugin.value = plugin
  detailVisible.value = true
}

onMounted(() => { fetchPlugins() })
</script>

<style scoped lang="scss">
.plugin-store-page {
  .category-tabs {
    margin-bottom: 20px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    .category-tag { cursor: pointer; transition: all 0.2s; &:hover { opacity: 0.8; } }
  }
  .plugin-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
    margin-bottom: 20px;
  }
  .plugin-card {
    background: var(--bg-card);
    border-radius: var(--border-radius);
    padding: 20px;
    display: flex;
    gap: 16px;
    box-shadow: var(--shadow-sm);
    transition: all 0.3s ease;
    &:hover { transform: translateY(-2px); box-shadow: var(--shadow-md); }
    .plugin-icon-wrapper {
      width: 64px; height: 64px; flex-shrink: 0;
      background: var(--primary-bg, rgba(64,158,255,0.1));
      border-radius: 12px;
      display: flex; align-items: center; justify-content: center;
      color: var(--primary-color, #409eff);
    }
    .plugin-info {
      flex: 1; min-width: 0;
      .plugin-header { display: flex; align-items: center; gap: 8px; margin-bottom: 6px;
        h4 { margin: 0; font-size: 15px; font-weight: 600; }
      }
      .plugin-desc { color: var(--text-secondary); font-size: 13px; margin: 0 0 8px;
        display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
      }
      .plugin-meta { display: flex; gap: 12px; color: var(--text-secondary); font-size: 12px; margin-bottom: 12px;
        span { display: flex; align-items: center; gap: 3px; }
        .rating { color: #e6a23c; }
      }
      .plugin-footer { display: flex; justify-content: space-between; align-items: center;
        .version { color: var(--text-secondary); font-size: 12px; }
        .actions { display: flex; gap: 6px; }
      }
    }
  }
  .pagination { display: flex; justify-content: flex-end; }
  .store-detail {
    .detail-header { display: flex; align-items: center; gap: 16px;
      .detail-title { h3 { margin: 0 0 4px; font-size: 18px; } p { margin: 0; color: var(--text-secondary); font-size: 13px; } }
    }
  }
}
</style>
