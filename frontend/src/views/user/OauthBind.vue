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
import { ref } from 'vue'

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

const bind = (provider: any) => {
  // TODO: 跳转到第三方授权页面
}

const unbind = async (provider: any) => {
  // TODO: 解绑第三方账号
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
