<template>
  <div class="oauth-bind-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>第三方账号绑定</span>
        </div>
      </template>

      <div class="oauth-list">
        <div
          v-for="provider in providers"
          :key="provider.name"
          class="oauth-item"
        >
          <div class="oauth-info">
            <img :src="provider.icon" :alt="provider.label" class="oauth-icon" />
            <div class="oauth-text">
              <div class="oauth-name">{{ provider.label }}</div>
              <div class="oauth-status" v-if="provider.bound">
                已绑定：{{ provider.account }}
              </div>
              <div class="oauth-status unbound" v-else>未绑定</div>
            </div>
          </div>
          <el-button
            v-if="provider.bound"
            type="danger"
            @click="unbind(provider)"
          >
            解绑
          </el-button>
          <el-button
            v-else
            type="primary"
            @click="bind(provider)"
          >
            绑定
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const providers = ref([
  {
    name: 'wechat',
    label: '微信',
    icon: '/assets/oauth/wechat.svg',
    bound: false,
    account: ''
  },
  {
    name: 'qq',
    label: 'QQ',
    icon: '/assets/oauth/qq.svg',
    bound: false,
    account: ''
  },
  {
    name: 'github',
    label: 'GitHub',
    icon: '/assets/oauth/github.svg',
    bound: false,
    account: ''
  },
  {
    name: 'google',
    label: 'Google',
    icon: '/assets/oauth/google.svg',
    bound: false,
    account: ''
  }
])

onMounted(async () => {
  try {
    const { data } = await request.get('/api/v1/oauth/providers')
    if (data?.data?.list?.length) {
      providers.value = data.data.list.map((p: any) => ({
        name: p.slug || p.name,
        label: p.display_name || p.name,
        icon: p.icon || `/assets/oauth/${p.slug || p.name}.svg`,
        bound: false,
        account: ''
      }))
    }
    // Fetch bound accounts
    const { data: bindData } = await request.get('/api/v1/oauth/accounts')
    if (bindData?.data) {
      const boundMap = new Map(bindData.data.map((b: any) => [b.provider, b]))
      providers.value.forEach(p => {
        const bound = boundMap.get(p.name)
        if (bound) {
          p.bound = true
          p.account = bound.username || bound.email || ''
        }
      })
    }
  } catch (e) {
    console.error('Failed to fetch OAuth providers:', e)
  }
})

const bind = (provider: any) => {
  window.location.href = `/api/v1/oauth/${provider.name}?action=bind`
}

const unbind = async (provider: any) => {
  await request.post(`/api/v1/oauth/${provider.name}/unbind`)
  provider.bound = false
  provider.account = ''
  ElMessage.success('解绑成功')
}
</script>

<style scoped lang="scss">
.oauth-bind-page {
  .oauth-list {
    .oauth-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px;
      border: 1px solid #ebeef5;
      border-radius: 8px;
      margin-bottom: 15px;

      &:hover {
        border-color: #409eff;
      }

      .oauth-info {
        display: flex;
        align-items: center;
        gap: 15px;

        .oauth-icon {
          width: 40px;
          height: 40px;
        }

        .oauth-name {
          font-weight: bold;
          margin-bottom: 4px;
        }

        .oauth-status {
          color: #67c23a;
          font-size: 14px;

          &.unbound {
            color: #909399;
          }
        }
      }
    }
  }
}
</style>
