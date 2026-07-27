<template>
  <div class="download-page">
    <!-- Header -->
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <n-icon size="24" color="#fff"><AnchorOutline /></n-icon>
          </div>
          <span class="logo-text">锚点财务</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <router-link to="/download" class="nav-link active">下载</router-link>
          <a href="#" class="nav-link">帮助</a>
        </nav>
        <div class="header-actions">
          <n-button text @click="$router.push('/login')">登录</n-button>
          <n-button type="primary" round size="small" @click="$router.push('/register')">免费注册</n-button>
        </div>
      </div>
    </header>

    <!-- Breadcrumb -->
    <div class="breadcrumb-bar">
      <div class="breadcrumb-inner">
        <n-breadcrumb>
          <n-breadcrumb-item @click="$router.push('/')">首页</n-breadcrumb-item>
          <n-breadcrumb-item>下载中心</n-breadcrumb-item>
        </n-breadcrumb>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <div class="content-inner">
        <!-- Sidebar -->
        <aside class="sidebar">
          <div class="sidebar-card">
            <h3 class="sidebar-title">
              <n-icon size="18" color="#1890ff"><FolderOpenOutline /></n-icon>
              文件分类
            </h3>
            <div class="group-list">
              <div
                v-for="group in fileGroups"
                :key="group.key"
                class="group-item"
                :class="{ active: selectedGroup === group.key }"
                @click="selectedGroup = group.key"
              >
                <n-icon size="18" :color="selectedGroup === group.key ? '#1890ff' : '#86909c'">
                  <component :is="group.icon" />
                </n-icon>
                <span class="group-label">{{ group.label }}</span>
                <span class="group-count">{{ group.count }}</span>
              </div>
            </div>
          </div>
        </aside>

        <!-- File List -->
        <main class="file-list-area">
          <!-- Toolbar -->
          <div class="toolbar">
            <div class="toolbar-left">
              <n-input
                v-model:value="searchKeyword"
                placeholder="搜索文件..."
                clearable
                style="width: 280px;"
              >
                <template #prefix>
                  <n-icon :component="SearchOutline" />
                </template>
              </n-input>
              <span class="result-count">共 <strong>{{ filteredFiles.length }}</strong> 个文件</span>
            </div>
          </div>

          <!-- Data Table -->
          <div v-if="filteredFiles.length > 0" class="table-wrap">
            <n-data-table
              :columns="columns"
              :data="paginatedFiles"
              :bordered="false"
              :single-line="false"
              striped
            />
          </div>

          <!-- Empty State -->
          <div v-else class="empty-state">
            <n-icon size="64" color="#c9cdd4"><CloudOfflineOutline /></n-icon>
            <p>暂无匹配的文件</p>
            <n-button type="primary" @click="resetFilters">重置筛选</n-button>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="pagination-wrap">
            <n-pagination
              v-model:page="currentPage"
              :page-count="totalPages"
              :page-slot="7"
              show-quick-jumper
            />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns } from 'naive-ui'
import { NButton, NTag, NIcon } from 'naive-ui'
import {
  AnchorOutline,
  FolderOpenOutline,
  SearchOutline,
  CloudOfflineOutline,
  DocumentOutline,
  CodeSlashOutline,
  SettingsOutline,
  ImageOutline,
  ArchiveOutline,
  DownloadOutline
} from '@vicons/ionicons5'

interface FileItem {
  id: number
  name: string
  category: string
  categoryKey: string
  version: string
  size: string
  updateTime: string
  downloads: number
  description: string
  icon: any
}

const searchKeyword = ref('')
const selectedGroup = ref('all')
const currentPage = ref(1)
const pageSize = 10

const fileGroups = computed(() => [
  { key: 'all', label: '全部文件', icon: FolderOpenOutline, count: files.value.length },
  { key: 'client', label: '客户端', icon: CodeSlashOutline, count: files.value.filter(f => f.categoryKey === 'client').length },
  { key: 'plugin', label: '插件扩展', icon: SettingsOutline, count: files.value.filter(f => f.categoryKey === 'plugin').length },
  { key: 'template', label: '模板主题', icon: ImageOutline, count: files.value.filter(f => f.categoryKey === 'template').length },
  { key: 'tool', label: '实用工具', icon: DocumentOutline, count: files.value.filter(f => f.categoryKey === 'tool').length },
  { key: 'sdk', label: 'SDK/API', icon: ArchiveOutline, count: files.value.filter(f => f.categoryKey === 'sdk').length }
])

const files = ref<FileItem[]>([
  {
    id: 1,
    name: 'AnchorFinance-Desktop-Client.exe',
    category: '客户端',
    categoryKey: 'client',
    version: 'v3.2.1',
    size: '68.5 MB',
    updateTime: '2026-07-20',
    downloads: 12580,
    description: 'Windows桌面客户端，支持账务管理、报表导出等功能',
    icon: CodeSlashOutline
  },
  {
    id: 2,
    name: 'AnchorFinance-Mac-Client.dmg',
    category: '客户端',
    categoryKey: 'client',
    version: 'v3.2.1',
    size: '72.3 MB',
    updateTime: '2026-07-20',
    downloads: 8920,
    description: 'macOS桌面客户端，适配Apple Silicon芯片',
    icon: CodeSlashOutline
  },
  {
    id: 3,
    name: 'WechatPay-Plugin-v2.1.0.zip',
    category: '插件扩展',
    categoryKey: 'plugin',
    version: 'v2.1.0',
    size: '2.8 MB',
    updateTime: '2026-07-18',
    downloads: 5634,
    description: '微信支付对接插件，支持JSAPI/Native/H5支付',
    icon: SettingsOutline
  },
  {
    id: 4,
    name: 'Alipay-Plugin-v2.0.5.zip',
    category: '插件扩展',
    categoryKey: 'plugin',
    version: 'v2.0.5',
    size: '3.1 MB',
    updateTime: '2026-07-15',
    downloads: 4892,
    description: '支付宝对接插件，支持当面付/手机网站支付',
    icon: SettingsOutline
  },
  {
    id: 5,
    name: 'Admin-Theme-Blue.zip',
    category: '模板主题',
    categoryKey: 'template',
    version: 'v1.5.0',
    size: '15.2 MB',
    updateTime: '2026-07-12',
    downloads: 3456,
    description: '蓝色主题后台管理模板，清新简洁风格',
    icon: ImageOutline
  },
  {
    id: 6,
    name: 'Admin-Theme-Dark.zip',
    category: '模板主题',
    categoryKey: 'template',
    version: 'v1.3.0',
    size: '16.8 MB',
    updateTime: '2026-07-10',
    downloads: 2789,
    description: '暗黑主题后台管理模板，护眼深色风格',
    icon: ImageOutline
  },
  {
    id: 7,
    name: 'DataExport-Tool-v1.2.0.exe',
    category: '实用工具',
    categoryKey: 'tool',
    version: 'v1.2.0',
    size: '12.5 MB',
    updateTime: '2026-07-08',
    downloads: 2345,
    description: '数据导出工具，支持Excel/CSV/PDF格式',
    icon: DocumentOutline
  },
  {
    id: 8,
    name: 'Database-Backup-Tool-v1.0.3.exe',
    category: '实用工具',
    categoryKey: 'tool',
    version: 'v1.0.3',
    size: '8.9 MB',
    updateTime: '2026-07-05',
    downloads: 1987,
    description: '数据库备份工具，支持定时自动备份',
    icon: DocumentOutline
  },
  {
    id: 9,
    name: 'AnchorFinance-PHP-SDK-v2.3.0.zip',
    category: 'SDK/API',
    categoryKey: 'sdk',
    version: 'v2.3.0',
    size: '1.2 MB',
    updateTime: '2026-07-22',
    downloads: 3210,
    description: 'PHP SDK，支持Composer安装，完善的API封装',
    icon: ArchiveOutline
  },
  {
    id: 10,
    name: 'AnchorFinance-Python-SDK-v1.8.0.tar.gz',
    category: 'SDK/API',
    categoryKey: 'sdk',
    version: 'v1.8.0',
    size: '0.9 MB',
    updateTime: '2026-07-20',
    downloads: 2678,
    description: 'Python SDK，支持pip安装，完整类型提示',
    icon: ArchiveOutline
  },
  {
    id: 11,
    name: 'AnchorFinance-Java-SDK-v1.5.0.jar',
    category: 'SDK/API',
    categoryKey: 'sdk',
    version: 'v1.5.0',
    size: '2.4 MB',
    updateTime: '2026-07-18',
    downloads: 2145,
    description: 'Java SDK，Maven中央仓库可直接引用',
    icon: ArchiveOutline
  },
  {
    id: 12,
    name: 'Stripe-Plugin-v1.2.0.zip',
    category: '插件扩展',
    categoryKey: 'plugin',
    version: 'v1.2.0',
    size: '3.5 MB',
    updateTime: '2026-07-16',
    downloads: 1876,
    description: 'Stripe国际支付插件，支持信用卡/Apple Pay',
    icon: SettingsOutline
  }
])

const filteredFiles = computed(() => {
  let list = [...files.value]

  if (selectedGroup.value !== 'all') {
    list = list.filter(f => f.categoryKey === selectedGroup.value)
  }

  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.trim().toLowerCase()
    list = list.filter(
      f => f.name.toLowerCase().includes(keyword) || f.description.toLowerCase().includes(keyword)
    )
  }

  return list
})

const totalPages = computed(() => Math.ceil(filteredFiles.value.length / pageSize))

const paginatedFiles = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredFiles.value.slice(start, start + pageSize)
})

function resetFilters() {
  selectedGroup.value = 'all'
  searchKeyword.value = ''
  currentPage.value = 1
}

function handleDownload(file: FileItem) {
  // TODO: implement actual download logic
  console.log('Downloading:', file.name)
}

const columns: DataTableColumns<FileItem> = [
  {
    title: '文件名',
    key: 'name',
    minWidth: 280,
    render(row) {
      return h('div', { style: 'display: flex; align-items: center; gap: 10px;' }, [
        h(NIcon, { size: 20, color: '#1890ff' }, { default: () => h(row.icon) }),
        h('div', null, [
          h('div', { style: 'font-weight: 600; color: #1d2129; font-size: 14px;' }, row.name),
          h('div', { style: 'font-size: 12px; color: #86909c; margin-top: 2px;' }, row.description)
        ])
      ])
    }
  },
  {
    title: '分类',
    key: 'category',
    width: 100,
    render(row) {
      return h(NTag, { size: 'small', bordered: false, type: 'info' }, { default: () => row.category })
    }
  },
  {
    title: '版本',
    key: 'version',
    width: 90,
    render(row) {
      return h(NTag, { size: 'small', bordered: false, type: 'success' }, { default: () => row.version })
    }
  },
  {
    title: '大小',
    key: 'size',
    width: 100
  },
  {
    title: '更新时间',
    key: 'updateTime',
    width: 120
  },
  {
    title: '下载次数',
    key: 'downloads',
    width: 100,
    sorter: (a, b) => a.downloads - b.downloads,
    render(row) {
      return h('span', { style: 'color: #1890ff; font-weight: 500;' }, row.downloads.toLocaleString())
    }
  },
  {
    title: '操作',
    key: 'action',
    width: 100,
    render(row) {
      return h(
        NButton,
        {
          type: 'primary',
          size: 'small',
          round: true,
          onClick: (e: Event) => {
            e.stopPropagation()
            handleDownload(row)
          }
        },
        {
          icon: () => h(NIcon, null, { default: () => h(DownloadOutline) }),
          default: () => '下载'
        }
      )
    }
  }
]
</script>

<style scoped>
.download-page {
  min-height: 100vh;
  background: #f7f8fa;
}

/* Header */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: #fff;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.06);
  z-index: 100;
}

.header-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  width: 34px;
  height: 34px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
}

.nav-links {
  display: flex;
  gap: 32px;
}

.nav-link {
  color: #4e5969;
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-link:hover,
.nav-link.active {
  color: #1890ff;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Breadcrumb */
.breadcrumb-bar {
  background: #fff;
  border-bottom: 1px solid #f0f1f5;
  margin-top: 64px;
}

.breadcrumb-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 14px 24px;
}

/* Main Content */
.main-content {
  padding: 24px 0 40px;
}

.content-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  gap: 24px;
}

/* Sidebar */
.sidebar {
  width: 240px;
  flex-shrink: 0;
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.sidebar-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.group-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #4e5969;
}

.group-item:hover {
  background: #f2f3f5;
  color: #1d2129;
}

.group-item.active {
  background: #e6f7ff;
  color: #1890ff;
}

.group-label {
  flex: 1;
  font-weight: 500;
}

.group-count {
  font-size: 12px;
  color: #c9cdd4;
  background: #f2f3f5;
  padding: 1px 8px;
  border-radius: 10px;
}

.group-item.active .group-count {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
}

/* File List Area */
.file-list-area {
  flex: 1;
  min-width: 0;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.result-count {
  font-size: 14px;
  color: #86909c;
  white-space: nowrap;
}

.result-count strong {
  color: #1890ff;
}

/* Table */
.table-wrap {
  background: #fff;
  border-radius: 12px;
  padding: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 80px 0;
  color: #c9cdd4;
  background: #fff;
  border-radius: 12px;
}

.empty-state p {
  margin: 16px 0 24px;
  font-size: 15px;
}

/* Pagination */
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

/* Responsive */
@media (max-width: 1024px) {
  .sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
}
</style>
