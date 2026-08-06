<template>
  <div class="document-container">
    <!-- Navbar -->
    <div class="doc-navbar">
      <button class="hamburger" @click="toggleSidebar" aria-label="Toggle menu">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 12h18M3 6h18M3 18h18" />
        </svg>
      </button>
      <router-link to="/" class="brand">锚点财务API文档</router-link>
      <div class="nav-links" :class="{ 'show': showMobileNav }">
        <a v-for="(data, key) in apiData" :key="key" :class="{ 'active': activePage === key }" @click="setPage(key as PageType)">
          {{ data.title }}
        </a>
      </div>
      <div class="search-box" ref="searchBox">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索API..."
          @input="handleSearch"
          @focus="showSearchResults = true"
          @keydown.escape="showSearchResults = false"
        >
        <span class="search-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
        </span>
        <div v-if="showSearchResults && filteredResults.length > 0" class="search-results">
          <div
            v-for="result in filteredResults"
            :key="result.id"
            class="search-result-item"
            @click="goToApi(result.page, result.id)"
          >
            <span class="method-badge" :class="getBadgeClass(result.method)">{{ result.method }}</span>
            <span class="result-title">{{ result.title }}</span>
            <span class="result-url">{{ result.url }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="doc-layout">
      <!-- Sidebar Overlay -->
      <div v-if="showSidebar" class="sidebar-overlay" @click="showSidebar = false"></div>

      <!-- Sidebar -->
      <div class="side-nav" :class="{ 'open': showSidebar }">
        <!-- Mobile Page Switch -->
        <div class="sidebar-page-nav">
          <a v-for="(data, key) in apiData" :key="key" :class="{ 'active': activePage === key }" @click="setPage(key as PageType)">
            {{ data.title }}
          </a>
        </div>

        <div class="panel-group">
          <div v-for="group in currentPageData?.groups" :key="group.name" class="panel">
            <div class="panel-heading" @click="toggleGroup(group.name)">
              <span class="panel-title">{{ group.name }}</span>
              <span class="arrow" :class="{ 'open': expandedGroups.has(group.name) }">
                <svg width="12" height="12" viewBox="0 0 12 12">
                  <path d="M3 4.5l3 3 3-3" stroke="currentColor" stroke-width="1.5" fill="none" />
                </svg>
              </span>
            </div>
            <transition name="slide">
              <div v-show="expandedGroups.has(group.name)" class="panel-collapse">
                <a
                  v-for="item in group.items"
                  :key="item.id"
                  :class="{ 'active': activeApiId === item.id }"
                  @click="setApi(item.id)"
                >
                  <span class="method-badge" :class="getBadgeClass(item.method)">{{ item.method }}</span>
                  <span class="item-title">{{ item.title }}</span>
                </a>
              </div>
            </transition>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="doc-content">
        <div v-if="currentApi" class="api-detail">
          <h2>{{ currentApi.title }}</h2>
          <div class="api-meta">
            <span class="method-badge lg" :class="getBadgeClass(currentApi.method)">{{ currentApi.method }}</span>
            <code class="api-url">{{ currentApi.url }}</code>
          </div>
          <p class="api-desc">{{ currentApi.desc }}</p>

          <!-- Request Params -->
          <div v-if="currentApi.reqParams.length > 0" class="params-section">
            <h3 @click="reqParamsExpanded = !reqParamsExpanded" class="toggle-header">
              请求参数
              <span class="toggle-icon">{{ reqParamsExpanded ? '−' : '+' }}</span>
            </h3>
            <div v-show="reqParamsExpanded" class="table-wrapper">
              <table>
                <thead>
                  <tr>
                    <th>参数名</th>
                    <th>类型</th>
                    <th>必填</th>
                    <th>说明</th>
                    <th>示例</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="param in currentApi.reqParams" :key="param.n">
                    <td><code>{{ param.n }}</code></td>
                    <td>{{ param.t }}</td>
                    <td>
                      <span :class="param.r === '必填' ? 'required' : 'optional'">{{ param.r }}</span>
                    </td>
                    <td>{{ param.d }}</td>
                    <td><code>{{ param.e }}</code></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Response Params -->
          <div v-if="currentApi.resParams.length > 0" class="params-section">
            <h3 @click="resParamsExpanded = !resParamsExpanded" class="toggle-header">
              返回参数
              <span class="toggle-icon">{{ resParamsExpanded ? '−' : '+' }}</span>
            </h3>
            <div v-show="resParamsExpanded" class="table-wrapper">
              <table>
                <thead>
                  <tr>
                    <th>参数名</th>
                    <th>类型</th>
                    <th>说明</th>
                    <th>示例</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="param in currentApi.resParams" :key="param.n">
                    <td><code>{{ param.n }}</code></td>
                    <td>{{ param.t }}</td>
                    <td>{{ param.d }}</td>
                    <td><code>{{ param.e }}</code></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Request Example -->
          <div class="example-section">
            <h3>请求示例</h3>
            <pre class="code-block"><code>{{ currentApi.reqExample }}</code></pre>
          </div>

          <!-- Response Example -->
          <div class="example-section">
            <h3>返回示例</h3>
            <pre class="code-block"><code>{{ formatJson(currentApi.resExample) }}</code></pre>
          </div>
        </div>
        <div v-else class="no-content">
          <p>请从左侧选择一个API</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { API_DATA, type PageType, type ApiItem, type ApiPage } from './document/apiData'

const route = useRoute()
const router = useRouter()

const apiData = API_DATA
const activePage = ref<PageType>('public')
const activeApiId = ref<string>('')
const showSidebar = ref(false)
const showMobileNav = ref(false)
const searchQuery = ref('')
const showSearchResults = ref(false)
const expandedGroups = ref(new Set<string>())
const reqParamsExpanded = ref(true)
const resParamsExpanded = ref(true)
const searchBox = ref<HTMLElement | null>(null)

const currentPageData = computed(() => apiData[activePage.value])

const currentApi = computed(() => {
  if (!currentPageData.value || !activeApiId.value) return null
  for (const group of currentPageData.value.groups) {
    const item = group.items.find((i: ApiItem) => i.id === activeApiId.value)
    if (item) return item
  }
  return null
})

const filteredResults = computed(() => {
  if (!searchQuery.value.trim()) return []
  const query = searchQuery.value.toLowerCase()
  const results: Array<ApiItem & { page: PageType }> = []

  for (const [page, data] of Object.entries(apiData) as [string, ApiPage][]) {
    for (const group of data.groups) {
      for (const item of group.items) {
        if (
          item.title.toLowerCase().includes(query) ||
          item.url.toLowerCase().includes(query) ||
          item.desc.toLowerCase().includes(query)
        ) {
          results.push({ ...item, page: page as PageType })
        }
      }
    }
  }
  return results.slice(0, 20)
})

function setPage(page: PageType) {
  activePage.value = page
  const firstItem = apiData[page].groups[0]?.items[0]
  if (firstItem) {
    setApi(firstItem.id)
  }
  showSidebar.value = false
  showMobileNav.value = false
}

function setApi(id: string) {
  activeApiId.value = id
  showSidebar.value = false
  updateUrl()
  // Auto-expand parent group
  for (const group of currentPageData.value?.groups || []) {
    if (group.items.some((i: ApiItem) => i.id === id)) {
      expandedGroups.value.add(group.name)
    }
  }
}

function toggleGroup(name: string) {
  if (expandedGroups.value.has(name)) {
    expandedGroups.value.delete(name)
  } else {
    expandedGroups.value.add(name)
  }
}

function toggleSidebar() {
  showSidebar.value = !showSidebar.value
}

function handleSearch() {
  showSearchResults.value = true
}

function goToApi(page: PageType, id: string) {
  activePage.value = page
  activeApiId.value = id
  searchQuery.value = ''
  showSearchResults.value = false
  updateUrl()
}

function updateUrl() {
  router.replace({
    path: '/document',
    query: { page: activePage.value, api: activeApiId.value }
  })
}

function getBadgeClass(method: string): string {
  const classes: Record<string, string> = {
    GET: 'badge-get',
    POST: 'badge-post',
    PUT: 'badge-put',
    DELETE: 'badge-delete',
    PATCH: 'badge-patch'
  }
  return classes[method] || 'badge-default'
}

function formatJson(str: string): string {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

function initFromUrl() {
  const page = (route.query.page as string) || 'public'
  const api = (route.query.api as string) || ''

  if (page in apiData) {
    activePage.value = page as PageType
  }

  if (api) {
    activeApiId.value = api
  } else {
    const firstItem = apiData[activePage.value].groups[0]?.items[0]
    if (firstItem) {
      activeApiId.value = firstItem.id
    }
  }

  // Auto-expand current group
  for (const group of apiData[activePage.value].groups) {
    if (group.items.some((i: ApiItem) => i.id === activeApiId.value)) {
      expandedGroups.value.add(group.name)
    }
  }
}

function handleClickOutside(e: MouseEvent) {
  if (searchBox.value && !searchBox.value.contains(e.target as Node)) {
    showSearchResults.value = false
  }
}

// Keyboard shortcut
function handleKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    const input = searchBox.value?.querySelector('input')
    input?.focus()
  }
}

watch(() => route.query, initFromUrl)

onMounted(() => {
  initFromUrl()
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.document-container {
  min-height: 100vh;
  background: #f5f7fa;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

/* Navbar */
.doc-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  padding: 0 20px;
  z-index: 1000;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
}

.hamburger {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  color: #333;
}

.brand {
  font-size: 18px;
  font-weight: 600;
  color: #409eff;
  text-decoration: none;
  margin-right: 40px;
}

.nav-links {
  display: flex;
  gap: 24px;
}

.nav-links a {
  color: #666;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  padding: 4px 0;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
}

.nav-links a:hover,
.nav-links a.active {
  color: #409eff;
  border-bottom-color: #409eff;
}

/* Search */
.search-box {
  position: relative;
  margin-left: auto;
}

.search-box input {
  width: 240px;
  height: 36px;
  padding: 0 36px 0 12px;
  border: 1px solid #dcdfe6;
  border-radius: 18px;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
}

.search-box input:focus {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64,158,255,0.2);
}

.search-icon {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #999;
  pointer-events: none;
}

.search-results {
  position: absolute;
  top: 100%;
  right: 0;
  width: 400px;
  max-height: 400px;
  overflow-y: auto;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  margin-top: 4px;
}

.search-result-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.search-result-item:hover {
  background: #f5f7fa;
}

.result-title {
  flex: 1;
  font-size: 14px;
  color: #333;
}

.result-url {
  font-size: 12px;
  color: #999;
  font-family: monospace;
}

/* Layout */
.doc-layout {
  display: flex;
  margin-top: 60px;
  min-height: calc(100vh - 60px);
}

/* Sidebar */
.side-nav {
  position: fixed;
  top: 60px;
  left: 0;
  width: 260px;
  height: calc(100vh - 60px);
  overflow-y: auto;
  background: #fff;
  border-right: 1px solid #e4e7ed;
  padding: 16px 0;
  z-index: 100;
}

.sidebar-page-nav {
  display: none;
  padding: 0 16px 12px;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 8px;
  gap: 8px;
}

.sidebar-page-nav a {
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s;
}

.sidebar-page-nav a.active {
  background: #ecf5ff;
  color: #409eff;
}

.panel-group {
  padding: 0 8px;
}

.panel {
  margin-bottom: 4px;
}

.panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.2s;
}

.panel-heading:hover {
  background: #f5f7fa;
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.arrow {
  color: #999;
  transition: transform 0.3s;
}

.arrow.open {
  transform: rotate(0deg);
}

.arrow:not(.open) {
  transform: rotate(-90deg);
}

.panel-collapse {
  padding-left: 12px;
}

.panel-collapse a {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  color: #666;
  text-decoration: none;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.panel-collapse a:hover {
  background: #f5f7fa;
  color: #333;
}

.panel-collapse a.active {
  background: #ecf5ff;
  color: #409eff;
}

.item-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Sidebar Overlay */
.sidebar-overlay {
  display: none;
}

/* Content */
.doc-content {
  flex: 1;
  margin-left: 260px;
  padding: 24px;
  max-width: 900px;
}

.api-detail h2 {
  font-size: 24px;
  color: #333;
  margin-bottom: 16px;
}

.api-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.api-url {
  font-size: 15px;
  color: #666;
  background: #f5f7fa;
  padding: 6px 12px;
  border-radius: 6px;
  font-family: monospace;
}

.api-desc {
  color: #666;
  margin-bottom: 24px;
  line-height: 1.6;
}

/* Method Badges */
.method-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  min-width: 42px;
  text-align: center;
}

.method-badge.lg {
  padding: 4px 12px;
  font-size: 13px;
}

.badge-get {
  background: #e1f3d8;
  color: #67c23a;
}

.badge-post {
  background: #ecf5ff;
  color: #409eff;
}

.badge-put {
  background: #fdf6ec;
  color: #e6a23c;
}

.badge-delete {
  background: #fef0f0;
  color: #f56c6c;
}

.badge-patch {
  background: #f4f4f5;
  color: #909399;
}

.badge-default {
  background: #f4f4f5;
  color: #909399;
}

/* Params Section */
.params-section {
  margin-bottom: 24px;
}

.params-section h3,
.example-section h3 {
  font-size: 16px;
  color: #333;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.toggle-header {
  cursor: pointer;
  user-select: none;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.toggle-icon {
  font-size: 20px;
  color: #999;
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

th {
  background: #f5f7fa;
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  color: #333;
  border-bottom: 1px solid #ebeef5;
}

td {
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  color: #666;
}

td code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  color: #409eff;
}

.required {
  color: #f56c6c;
  font-weight: 500;
}

.optional {
  color: #909399;
}

/* Code Example */
.example-section {
  margin-bottom: 24px;
}

.code-block {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  font-size: 13px;
  line-height: 1.6;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
}

.code-block code {
  color: inherit;
  background: none;
}

/* No Content */
.no-content {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400px;
  color: #999;
  font-size: 16px;
}

/* Slide Transition */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
}

.slide-enter-to,
.slide-leave-from {
  max-height: 1000px;
  opacity: 1;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .hamburger {
    display: block;
  }

  .nav-links {
    display: none;
    position: absolute;
    top: 60px;
    left: 0;
    right: 0;
    background: #fff;
    flex-direction: column;
    padding: 16px;
    border-bottom: 1px solid #e4e7ed;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  }

  .nav-links.show {
    display: flex;
  }

  .search-box input {
    width: 160px;
  }

  .side-nav {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .side-nav.open {
    transform: translateX(0);
  }

  .sidebar-overlay {
    display: block;
    position: fixed;
    top: 60px;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0,0,0,0.3);
    z-index: 99;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.3s;
  }

  .sidebar-overlay.show {
    opacity: 1;
    pointer-events: auto;
  }

  .doc-content {
    margin-left: 0;
    padding: 16px;
  }

  .sidebar-page-nav {
    display: flex;
  }

  .search-results {
    width: calc(100vw - 32px);
    right: -16px;
  }
}
</style>
