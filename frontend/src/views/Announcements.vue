<template>
  <div class="announcements-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <n-icon :component="MegaphoneOutline" size="28" />
          公告中心
        </h1>
        <p class="page-subtitle">了解最新动态与重要通知</p>
      </div>
      <div class="header-actions">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索公告..."
          clearable
          style="width: 280px"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" color="#1890ff" />
          </template>
        </n-input>
      </div>
    </div>

    <div v-if="selectedAnnouncement" class="detail-section">
      <n-card class="detail-card" :bordered="false">
        <div class="detail-back" @click="selectedAnnouncement = null">
          <n-icon :component="ArrowBackOutline" size="18" />
          <span>返回列表</span>
        </div>
        <div class="detail-header">
          <h2 class="detail-title">{{ selectedAnnouncement.title }}</h2>
          <div class="detail-meta">
            <n-tag :type="getTagType(selectedAnnouncement.level)" size="small" :bordered="false">
              {{ selectedAnnouncement.levelLabel }}
            </n-tag>
            <span class="detail-time">
              <n-icon :component="TimeOutline" size="14" />
              {{ selectedAnnouncement.time }}
            </span>
          </div>
        </div>
        <n-divider />
        <div class="detail-content" v-html="selectedAnnouncement.content"></div>
      </n-card>
    </div>

    <div v-else class="list-section">
      <div v-if="loading" class="empty-state">
        <n-spin size="medium" />
      </div>
      <div v-else-if="filteredList.length > 0" class="announcement-list">
        <n-card
          v-for="item in paginatedList"
          :key="item.id"
          class="announcement-card"
          :bordered="false"
          hoverable
          @click="selectedAnnouncement = item"
        >
          <div class="card-header">
            <div class="card-title-row">
              <n-tag :type="getTagType(item.level)" size="small" :bordered="false">
                {{ item.levelLabel }}
              </n-tag>
              <h3 class="card-title">{{ item.title }}</h3>
            </div>
            <span class="card-time">
              <n-icon :component="TimeOutline" size="14" />
              {{ item.time }}
            </span>
          </div>
          <p class="card-summary">{{ item.summary }}</p>
          <div class="card-footer">
            <div class="card-tags">
              <n-tag
                v-for="tag in item.tags"
                :key="tag"
                size="tiny"
                :bordered="false"
                class="meta-tag"
              >
                {{ tag }}
              </n-tag>
            </div>
            <span class="read-more">
              查看详情
              <n-icon :component="ChevronForwardOutline" size="14" />
            </span>
          </div>
        </n-card>
      </div>

      <div v-else class="empty-state">
        <n-empty description="暂无公告" />
      </div>

      <div v-if="totalPages > 1" class="pagination-wrapper">
        <n-pagination
          v-model:page="currentPage"
          :page-count="totalPages"
          :page-slot="7"
          show-quick-jumper
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  MegaphoneOutline,
  SearchOutline,
  TimeOutline,
  ArrowBackOutline,
  ChevronForwardOutline
} from '@vicons/ionicons5'
import request from '@/utils/request'

const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = 6
const selectedAnnouncement = ref<Announcement | null>(null)
const loading = ref(false)

interface Announcement {
  id: number
  title: string
  summary: string
  content: string
  time: string
  level: string
  levelLabel: string
  tags: string[]
}

const announcements = ref<Announcement[]>([])

async function fetchAnnouncements() {
  loading.value = true
  try {
    const res = await request.get('/api/v2/announcements')
    announcements.value = res.data?.data || res.data || []
  } catch { /* ignore */ }
  loading.value = false
}

onMounted(() => { fetchAnnouncements() })

const tagTypeMap: Record<string, string> = {
  important: 'error',
  new: 'success',
  notice: 'info',
  activity: 'warning'
}

function getTagType(level: string) {
  return (tagTypeMap[level] || 'default') as 'error' | 'success' | 'info' | 'warning' | 'default'
}

const filteredList = computed(() => {
  if (!searchKeyword.value) return announcements.value
  const kw = searchKeyword.value.toLowerCase()
  return announcements.value.filter(
    a =>
      a.title.toLowerCase().includes(kw) ||
      a.summary.toLowerCase().includes(kw) ||
      a.tags.some(t => t.toLowerCase().includes(kw))
  )
})

const totalPages = computed(() => Math.ceil(filteredList.value.length / pageSize))

const paginatedList = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredList.value.slice(start, start + pageSize)
})

function handleSearch() {
  currentPage.value = 1
}
</script>

<style scoped>
.announcements-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 28px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 24px;
  font-weight: 700;
  color: #262626;
  margin: 0 0 6px 0;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.detail-section {
  animation: fadeIn 0.3s ease;
}

.detail-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(24, 144, 255, 0.06);
}

.detail-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #1890ff;
  cursor: pointer;
  font-size: 14px;
  margin-bottom: 20px;
  transition: color 0.3s;
}

.detail-back:hover {
  color: #40a9ff;
}

.detail-header {
  margin-bottom: 8px;
}

.detail-title {
  font-size: 22px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 12px 0;
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.detail-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #8c8c8c;
  font-size: 13px;
}

.detail-content {
  font-size: 15px;
  line-height: 1.8;
  color: #595959;
}

.detail-content :deep(p) {
  margin: 0 0 12px 0;
}

.detail-content :deep(ul) {
  padding-left: 20px;
  margin: 8px 0 16px 0;
}

.detail-content :deep(li) {
  margin-bottom: 6px;
}

.detail-content :deep(strong) {
  color: #262626;
}

.list-section {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.announcement-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(24, 144, 255, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;
}

.announcement-card:hover {
  box-shadow: 0 6px 24px rgba(24, 144, 255, 0.12);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-time {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #bfbfbf;
  font-size: 13px;
  flex-shrink: 0;
}

.card-summary {
  color: #595959;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 14px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.meta-tag {
  background: #f0f5ff;
  color: #1890ff;
}

.read-more {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #1890ff;
  font-size: 13px;
  font-weight: 500;
  transition: color 0.3s;
}

.announcement-card:hover .read-more {
  color: #40a9ff;
}

.empty-state {
  padding: 80px 0;
  text-align: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}
</style>
