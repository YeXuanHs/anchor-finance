<template>
  <div class="maintenance-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>维护通知</span>
        </div>
      </template>

      <div class="notice-list">
        <div
          v-for="notice in notices"
          :key="notice.id"
          class="notice-item"
          @click="viewDetail(notice)"
        >
          <div class="notice-header">
            <el-tag :type="getLevelType(notice.level)" size="small">
              {{ getLevelText(notice.level) }}
            </el-tag>
            <span class="notice-time">{{ notice.created_at }}</span>
          </div>
          <h3 class="notice-title">{{ notice.title }}</h3>
          <p class="notice-summary">{{ notice.summary }}</p>
          <div class="notice-meta">
            <span v-if="notice.affected_services" class="meta-item">
              <el-icon><Monitor /></el-icon>
              影响服务：{{ notice.affected_services }}
            </span>
            <span v-if="notice.start_time" class="meta-item">
              <el-icon><Clock /></el-icon>
              预计时间：{{ notice.start_time }} ~ {{ notice.end_time }}
            </span>
          </div>
        </div>

        <el-empty v-if="notices.length === 0" description="暂无维护通知" />
      </div>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadNotices"
          @current-change="loadNotices"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="currentNotice?.title" width="680px">
      <div v-if="currentNotice" class="notice-detail">
        <div class="detail-meta">
          <el-tag :type="getLevelType(currentNotice.level)">
            {{ getLevelText(currentNotice.level) }}
          </el-tag>
          <span class="meta-time">发布时间：{{ currentNotice.created_at }}</span>
        </div>

        <div class="detail-section" v-if="currentNotice.affected_services">
          <h4>影响服务</h4>
          <p>{{ currentNotice.affected_services }}</p>
        </div>

        <div class="detail-section" v-if="currentNotice.start_time">
          <h4>维护时间</h4>
          <p>{{ currentNotice.start_time }} ~ {{ currentNotice.end_time }}</p>
        </div>

        <div class="detail-section">
          <h4>详细说明</h4>
          <div class="detail-content" v-html="currentNotice.content"></div>
        </div>

        <div class="detail-section" v-if="currentNotice.actions">
          <h4>用户操作建议</h4>
          <ul class="action-list">
            <li v-for="(action, index) in currentNotice.actions" :key="index">
              {{ action }}
            </li>
          </ul>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Monitor, Clock } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface MaintenanceNotice {
  id: number
  title: string
  summary: string
  content: string
  level: string
  created_at: string
  start_time?: string
  end_time?: string
  affected_services?: string
  actions?: string[]
}

const notices = ref<MaintenanceNotice[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const detailVisible = ref(false)
const currentNotice = ref<MaintenanceNotice | null>(null)

const getLevelType = (level: string) => {
  const map: Record<string, string> = {
    urgent: 'danger',
    important: 'warning',
    normal: 'info'
  }
  return map[level] || 'info'
}

const getLevelText = (level: string) => {
  const map: Record<string, string> = {
    urgent: '紧急',
    important: '重要',
    normal: '普通'
  }
  return map[level] || level
}

const viewDetail = (notice: MaintenanceNotice) => {
  currentNotice.value = notice
  detailVisible.value = true
}

const loadNotices = async () => {
  const res = await request.get('/api/v1/maintenance/notices', { params: { page: currentPage.value, page_size: pageSize.value } })
  notices.value = res.data.data.list
  total.value = res.data.data.total
}

onMounted(() => {
  loadNotices()
})
</script>

<style scoped lang="scss">
.maintenance-page {
  .notice-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .notice-item {
    padding: 20px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      border-color: #409eff;
      box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
    }

    .notice-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .notice-time {
        font-size: 13px;
        color: #909399;
      }
    }

    .notice-title {
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 8px 0;
    }

    .notice-summary {
      font-size: 14px;
      color: #606266;
      margin: 0 0 12px 0;
      line-height: 1.6;
    }

    .notice-meta {
      display: flex;
      gap: 24px;
      font-size: 13px;
      color: #909399;

      .meta-item {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
  }

  .pagination-wrap {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}

.notice-detail {
  .detail-meta {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 24px;

    .meta-time {
      font-size: 14px;
      color: #909399;
    }
  }

  .detail-section {
    margin-bottom: 20px;

    h4 {
      font-size: 15px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 8px 0;
    }

    p {
      font-size: 14px;
      color: #606266;
      line-height: 1.8;
      margin: 0;
    }

    .detail-content {
      font-size: 14px;
      color: #606266;
      line-height: 1.8;
    }

    .action-list {
      margin: 0;
      padding-left: 20px;

      li {
        font-size: 14px;
        color: #606266;
        line-height: 1.8;
      }
    }
  }
}
</style>
