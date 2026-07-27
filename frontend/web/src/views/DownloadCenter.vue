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
          <a href="#" class="nav-link">公告</a>
          <router-link to="/download" class="nav-link active">下载中心</router-link>
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
        <main class="file-list">
          <!-- Toolbar -->
          <div class="toolbar">
            <div class="toolbar-left">
              <n-input
                v-model:value="searchKeyword"
                placeholder="搜索文件名称..."
                size="small"
                clearable
                style="width: 260px;"
              >
                <template #prefix>
                  <n-icon :component="SearchOutline" />
                </template>
              </n-input>
            </div>
            <div class="toolbar-right">
              <span class="result-count">共 <strong>{{ filteredFiles.length }}</strong> 个文件</span>
            </div>
          </div>

          <!-- Data Table -->
          <n-card class="table-card" :bordered="false">
            <n-data-table
              v-if="filteredFiles.length > 0"
              :columns="columns"
              :data="paginatedFiles"
              :bordered="false"
              :single-line="false"
              striped
            />
            <n-empty v-else description="暂无匹配的文件" class="empty-state">
              <template #extra>
                <n-button type="primary" size="small" @click="clearFilters">清除筛选</n-button>
              </template>
            </n-empty>
          </n-card>

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
import {
  AnchorOutline,
  SearchOutline,
  FolderOpenOutline,
  ServerOutline,
  DesktopOutline,
  CodeSlashOutline,
  DocumentTextOutline,
  CloudDownloadOutline,
  ExtensionPuzzleOutline,
  DownloadOutline
} from '@vicons/ionicons5'
import { NButton, NIcon, NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

const searchKeyword = ref('')
const selectedGroup = ref('all')
const currentPage = ref(1)
const pageSize = 10

interface FileItem {
  id: number
  name: string
  category: string
  version: string
  size: string
  updateTime: string
  downloads: number
  description: string
  fileUrl: string
}

const files = ref<FileItem[]>([
  {
    id: 1,
    name: 'AnchorPanel-控制面板安装包',
    category: '面板程序',
    version: 'v3.2.1',
    size: '45.6 MB',
    updateTime: '2025-12-15',
    downloads: 8932,
    description: '服务器管理控制面板，支持一键部署、环境配置、站点管理等功能',
    fileUrl: '#'
  },
  {
    id: 2,
    name: 'AnchorPanel-Agent客户端',
    category: '面板程序',
    version: 'v3.2.0',
    size: '12.3 MB',
    updateTime: '2025-12-12',
    downloads: 6543,
    description: '配合控制面板使用的服务器端Agent程序',
    fileUrl: '#'
  },
  {
    id: 3,
    name: 'WordPress一键部署脚本',
    category: '部署脚本',
    version: 'v1.5.0',
    size: '2.1 MB',
    updateTime: '2025-12-10',
    downloads: 12456,
    description: '支持LNMP/LAMP环境自动配置WordPress站点',
    fileUrl: '#'
  },
  {
    id: 4,
    name: 'Node.js环境自动配置脚本',
    category: '部署脚本',
    version: 'v2.0.3',
    size: '1.8 MB',
    updateTime: '2025-12-08',
    downloads: 7821,
    description: '自动安装配置Node.js、PM2、Nginx反向代理',
    fileUrl: '#'
  },
  {
    id: 5,
    name: '服务器安全加固工具',
    category: '安全工具',
    version: 'v1.3.2',
    size: '8.7 MB',
    updateTime: '2025-12-05',
    downloads: 4567,
    description: '一键加固服务器安全配置，包括防火墙、SSH、端口管理',
    fileUrl: '#'
  },
  {
    id: 6,
    name: 'SSL证书自动申请工具',
    category: '安全工具',
    version: 'v1.1.0',
    size: '3.4 MB',
    updateTime: '2025-12-01',
    downloads: 3298,
    description: '基于Let\'s Encrypt自动申请和续期SSL证书',
    fileUrl: '#'
  },
  {
    id: 7,
    name: '数据备份恢复工具',
    category: '运维工具',
    version: 'v2.1.1',
    size: '15.2 MB',
    updateTime: '2025-11-28',
    downloads: 5678,
    description: '支持MySQL、PostgreSQL数据库的定时备份与一键恢复',
    fileUrl: '#'
  },
  {
    id: 8,
    name: '服务器迁移助手',
    category: '运维工具',
    version: 'v1.0.5',
    size: '6.9 MB',
    updateTime: '2025-11-25',
    downloads: 2345,
    description: '支持站点、数据库、邮件的一键跨服务器迁移',
    fileUrl: '#'
  },
  {
    id: 9,
    name: '产品API接口文档',
    category: '开发文档',
    version: 'v2.0',
    size: '4.5 MB',
    updateTime: '2025-11-20',
    downloads: 1890,
    description: '完整的RESTful API接口文档，含SDK示例代码',
    fileUrl: '#'
  },
  {
    id: 10,
    name: 'WebHook集成指南',
    category: '开发文档',
    version: 'v1.2',
    size: '2.3 MB',
    updateTime: '2025-11-18',
    downloads: 1234,
    description: 'WebHook事件订阅与回调处理开发指南',
    fileUrl: '#'
  },
  {
    id: 11,
    name: 'Windows远程连接工具',
    category: '客户端工具',
    version: 'v5.8.0',
    size: '28.6 MB',
    updateTime: '2025-11-15',
    downloads: 9876,
    description: '支持RDP、VNC、SSH多种协议的远程连接客户端',
    fileUrl: '#'
  },
  {
    id: 12,
    name: 'FTP文件传输工具',
    category: '客户端工具',
    version: 'v3.5.2',
    size: '18.4 MB',
    updateTime: '2025-11-10',
    downloads: 6543,
    description: '支持FTP/SFTP协议的文件传输管理工具',
    fileUrl: '#'
  }
])

const fileGroups = computed(() => {
  const categories = [...new Set(files.value.map(f => f.category))]
  const iconMap: Record<string, any> = {
    '面板程序': ServerOutline,
    '部署脚本': CodeSlashOutline,
    '安全工具': ExtensionPuzzleOutline,
    '运维工具': DesktopOutline,
    '开发文档': DocumentTextOutline,
    '客户端工具': CloudDownloadOutline
  }
  return [
    { key: 'all', label: '全部文件', icon: FolderOpenOutline, count: files.value.length },
    ...categories.map(cat => ({
      key: cat,
      label: cat,
      icon: iconMap[cat] || DocumentTextOutline,
      count: files.value.filter(f => f.category === cat).length
    }))
  ]
})

const filteredFiles = computed(() => {
  let list = [...files.value]

  if (selectedGroup.value !== 'all') {
    list = list.filter(f => f.category === selectedGroup.value)
  }

  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.trim().toLowerCase()
    list = list.filter(f =>
      f.name.toLowerCase().includes(keyword) ||
      f.description.toLowerCase().includes(keyword)
    )
  }

  return list
})

const totalPages = computed(() => Math.ceil(filteredFiles.value.length / pageSize))

const paginatedFiles = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredFiles.value.slice(start, start + pageSize)
})

function clearFilters() {
  searchKeyword.value = ''
  selectedGroup.value = 'all'
  currentPage.value = 1
}

function formatDownloads(count: number): string {
  if (count >= 10000) return (count / 10000).toFixed(1) + 'w'
  if (count >= 1000) return (count / 1000).toFixed(1) + 'k'
  return count.toString()
}

function handleDownload(file: FileItem) {
  window.open(file.fileUrl, '_blank')
}

const columns: DataTableColumns<FileItem> = [
  {
    title: '文件名称',
    key: 'name',
    minWidth: 260,
    render(row) {
      return h('div', { class: 'file-name-cell' }, [
        h('div', { class: 'file-icon-wrap' }, [
          h(NIcon, { size: 20, color: '#1890ff' }, { default: () => h(DownloadOutline) })
        ]),
        h('div', { class: 'file-info' }, [
          h('div', { class: 'file-name' }, row.name),
          h('div', { class: 'file-desc' }, row.description)
        ])
      ])
    }
  },
  {
    title: '分类',
    key: 'category',
    width: 120,
    render(row) {
      return h(NTag, { type: 'info', size: 'small', bordered: false }, { default: () => row.category })
    }
  },
  {
    title: '版本',
    key: 'version',
    width: 100,
    render(row) {
      return h(NTag, { type: 'success', size: 'small', bordered: false }, { default: () => row.version })
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
    render(row) {
      return h('span', { class: 'download-count' }, formatDownloads(row.downloads))
    }
  },
  {
    title: '操作',
    key: 'action',
    width: 100,
    fixed: 'right',
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
        { default: () => '下载', icon: () => h(NIcon, null, { default: () => h(DownloadOutline) }) }
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

/* File List */
.file-list {
  flex: 1;
  min-width: 0;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.result-count {
  font-size: 14px;
  color: #86909c;
}

.result-count strong {
  color: #1890ff;
}

.table-card {
  border-radius: 12px;
}

:deep(.file-name-cell) {
  display: flex;
  align-items: center;
  gap: 12px;
}

:deep(.file-icon-wrap) {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

:deep(.file-info) {
  min-width: 0;
}

:deep(.file-name) {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 2px;
}

:deep(.file-desc) {
  font-size: 12px;
  color: #86909c;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.download-count) {
  font-weight: 600;
  color: #1890ff;
}

.empty-state {
  padding: 60px 0;
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
