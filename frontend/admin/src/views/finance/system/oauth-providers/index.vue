<template>
  <div class="oauth-providers-page">
    <div class="page-header">
      <div class="header-left">
        <h2>{{ $t('oauthProvider.title') }}</h2>
        <span class="subtitle">{{ $t('oauthProvider.subtitle') }}</span>
      </div>
    </div>

    <!-- 平台分类 -->
    <el-tabs v-model="activeRegion">
      <el-tab-pane :label="$t('oauthProvider.all')" name="" />
      <el-tab-pane :label="$t('oauthProvider.domesticPlatforms')" name="cn" />
      <el-tab-pane :label="$t('oauthProvider.overseasPlatforms')" name="us" />
    </el-tabs>

    <!-- 平台列表 -->
    <el-row :gutter="20" class="provider-grid">
      <el-col :span="6" v-for="provider in filteredProviders" :key="provider.name">
        <el-card class="provider-card" :class="{ 'is-enabled': provider.is_enabled }">
          <div class="provider-header">
            <img :src="getProviderIcon(provider.name)" class="provider-icon" @error="handleIconError($event, provider.name)" />
            <div class="provider-info">
              <h3>{{ provider.title }}</h3>
              <el-tag :type="(provider.region === 'cn' ? '' : 'info') as any" size="small">
                {{ provider.region === 'cn' ? $t('oauthProvider.domestic') : $t('oauthProvider.overseas') }}
              </el-tag>
            </div>
            <el-switch
              v-model="provider.is_enabled"
              @change="handleToggle(provider)"
            />
          </div>
          <div class="provider-desc">{{ provider.description }}</div>
          <div class="provider-actions">
            <el-button type="primary" link size="small" @click="handleConfig(provider)">
              {{ $t('oauthProvider.config') }}
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      :title="$t('oauthProvider.dialogTitle', { name: currentProvider?.title || '' })"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form label-width="120px">
        <!-- 通用配置 -->
        <el-form-item label="App ID / Key">
          <el-input v-model="configForm.app_id" :placeholder="$t('oauthProvider.appIdPlaceholder')" />
        </el-form-item>
        <el-form-item label="App Secret">
          <el-input v-model="configForm.app_secret" type="password" show-password :placeholder="$t('oauthProvider.appSecretPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('oauthProvider.callbackUrl')">
          <el-input v-model="configForm.callback_url" :placeholder="getCallbackPlaceholder(currentProvider?.name)" />
          <div class="form-hint">{{ $t('oauthProvider.callbackHint') }}</div>
        </el-form-item>

        <!-- 特定平台配置 -->
        <template v-if="currentProvider?.name === 'wechat'">
          <el-divider content-position="left">{{ $t('oauthProvider.wechatConfig') }}</el-divider>
          <el-form-item :label="$t('oauthProvider.platformType')">
            <el-select v-model="configForm.platform">
              <el-option :label="$t('oauthProvider.websiteApp')" value="website" />
              <el-option :label="$t('oauthProvider.mobileApp')" value="mobile" />
              <el-option :label="$t('oauthProvider.mpApp')" value="mp" />
            </el-select>
          </el-form-item>
        </template>

        <template v-if="currentProvider?.name === 'dingtalk'">
          <el-divider content-position="left">{{ $t('oauthProvider.dingtalkConfig') }}</el-divider>
          <el-form-item :label="$t('oauthProvider.appType')">
            <el-select v-model="configForm.app_type">
              <el-option :label="$t('oauthProvider.internalApp')" value="internal" />
              <el-option :label="$t('oauthProvider.thirdPartyEnterpriseApp')" value="third" />
              <el-option :label="$t('oauthProvider.thirdPartyPersonalApp')" value="personal" />
            </el-select>
          </el-form-item>
        </template>

        <template v-if="currentProvider?.name === 'apple'">
          <el-divider content-position="left">{{ $t('oauthProvider.appleConfig') }}</el-divider>
          <el-form-item label="Team ID">
            <el-input v-model="configForm.team_id" />
          </el-form-item>
          <el-form-item label="Key ID">
            <el-input v-model="configForm.key_id" />
          </el-form-item>
          <el-form-item label="Private Key">
            <el-input v-model="configForm.private_key" type="textarea" :rows="4" />
          </el-form-item>
        </template>

        <template v-if="currentProvider?.name === 'telegram'">
          <el-divider content-position="left">{{ $t('oauthProvider.telegramConfig') }}</el-divider>
          <el-form-item label="Bot Token">
            <el-input v-model="configForm.bot_token" :placeholder="$t('oauthProvider.botFatherHint')" />
          </el-form-item>
          <el-form-item label="Bot Username">
            <el-input v-model="configForm.bot_username" placeholder="MyBot" />
          </el-form-item>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="configDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveConfig" :loading="submitting">{{ $t('oauthProvider.saveConfig') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

interface OAuthProvider {
  id: number
  name: string
  title: string
  description: string
  region: string
  icon: string
  config: string
  is_enabled: boolean
}

// 支持的24个平台
const getSupportedPlatforms = () => [
  // 国内平台
  { name: 'wechat', title: $t('oauthProvider.platforms.wechat.title'), region: 'cn', description: $t('oauthProvider.platforms.wechat.description') },
  { name: 'qq', title: $t('oauthProvider.platforms.qq.title'), region: 'cn', description: $t('oauthProvider.platforms.qq.description') },
  { name: 'weibo', title: $t('oauthProvider.platforms.weibo.title'), region: 'cn', description: $t('oauthProvider.platforms.weibo.description') },
  { name: 'alipay', title: $t('oauthProvider.platforms.alipay.title'), region: 'cn', description: $t('oauthProvider.platforms.alipay.description') },
  { name: 'baidu', title: $t('oauthProvider.platforms.baidu.title'), region: 'cn', description: $t('oauthProvider.platforms.baidu.description') },
  { name: 'gitee', title: $t('oauthProvider.platforms.gitee.title'), region: 'cn', description: $t('oauthProvider.platforms.gitee.description') },
  { name: 'dingtalk', title: $t('oauthProvider.platforms.dingtalk.title'), region: 'cn', description: $t('oauthProvider.platforms.dingtalk.description') },
  { name: 'feishu', title: $t('oauthProvider.platforms.feishu.title'), region: 'cn', description: $t('oauthProvider.platforms.feishu.description') },
  { name: 'csdn', title: $t('oauthProvider.platforms.csdn.title'), region: 'cn', description: $t('oauthProvider.platforms.csdn.description') },
  { name: 'oschina', title: $t('oauthProvider.platforms.oschina.title'), region: 'cn', description: $t('oauthProvider.platforms.oschina.description') },
  { name: 'tencent_cloud', title: $t('oauthProvider.platforms.tencent_cloud.title'), region: 'cn', description: $t('oauthProvider.platforms.tencent_cloud.description') },
  { name: 'aliyun', title: $t('oauthProvider.platforms.aliyun.title'), region: 'cn', description: $t('oauthProvider.platforms.aliyun.description') },
  
  // 海外平台
  { name: 'google', title: $t('oauthProvider.platforms.google.title'), region: 'us', description: $t('oauthProvider.platforms.google.description') },
  { name: 'facebook', title: $t('oauthProvider.platforms.facebook.title'), region: 'us', description: $t('oauthProvider.platforms.facebook.description') },
  { name: 'twitter', title: $t('oauthProvider.platforms.twitter.title'), region: 'us', description: $t('oauthProvider.platforms.twitter.description') },
  { name: 'github', title: $t('oauthProvider.platforms.github.title'), region: 'us', description: $t('oauthProvider.platforms.github.description') },
  { name: 'linkedin', title: $t('oauthProvider.platforms.linkedin.title'), region: 'us', description: $t('oauthProvider.platforms.linkedin.description') },
  { name: 'microsoft', title: $t('oauthProvider.platforms.microsoft.title'), region: 'us', description: $t('oauthProvider.platforms.microsoft.description') },
  { name: 'apple', title: $t('oauthProvider.platforms.apple.title'), region: 'us', description: $t('oauthProvider.platforms.apple.description') },
  { name: 'amazon', title: $t('oauthProvider.platforms.amazon.title'), region: 'us', description: $t('oauthProvider.platforms.amazon.description') },
  { name: 'discord', title: $t('oauthProvider.platforms.discord.title'), region: 'us', description: $t('oauthProvider.platforms.discord.description') },
  { name: 'slack', title: $t('oauthProvider.platforms.slack.title'), region: 'us', description: $t('oauthProvider.platforms.slack.description') },
  { name: 'telegram', title: $t('oauthProvider.platforms.telegram.title'), region: 'us', description: $t('oauthProvider.platforms.telegram.description') },
  { name: 'line', title: $t('oauthProvider.platforms.line.title'), region: 'us', description: $t('oauthProvider.platforms.line.description') },
]

const loading = ref(false)
const submitting = ref(false)
const configDialogVisible = ref(false)
const activeRegion = ref('')
const providers = ref<OAuthProvider[]>([])
const currentProvider = ref<OAuthProvider | null>(null)
const configForm = reactive<Record<string, any>>({})

// 过滤后的平台
const filteredProviders = computed(() => {
  if (!activeRegion.value) return providers.value
  return providers.value.filter(p => p.region === activeRegion.value)
})

// Simple Icons CDN slug映射（本地图标加载失败时使用）
const iconSlugMap: Record<string, string> = {
  wechat: 'wechat',
  qq: 'tencentqq',
  weibo: 'sinaweibo',
  alipay: 'alipay',
  baidu: 'baidu',
  gitee: 'gitee',
  dingtalk: 'dingtalk',
  feishu: 'lark',
  csdn: 'csdn',
  oschina: 'opensourceinitiative',
  tencent_cloud: 'tencentqq',
  aliyun: 'alibabacloud',
  google: 'google',
  facebook: 'facebook',
  twitter: 'x',
  github: 'github',
  linkedin: 'linkedin',
  microsoft: 'microsoft',
  apple: 'apple',
  amazon: 'amazon',
  discord: 'discord',
  slack: 'slack',
  telegram: 'telegram',
  line: 'line',
}

// 获取平台图标
const getProviderIcon = (name: string) => {
  // 优先使用本地SVG图标
  return `/assets/oauth/${name}.svg`
}

// 图标加载失败时使用CDN
const handleIconError = (e: Event, name: string) => {
  const img = e.target as HTMLImageElement
  const slug = iconSlugMap[name] || name
  img.src = `https://cdn.simpleicons.org/${slug}`
  img.onerror = null // 防止无限循环
}

// 获取回调地址占位符
const getCallbackPlaceholder = (name?: string) => {
  const domain = window.location.origin
  return `${domain}/oauth/${name}/callback`
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get({ url: '/api/admin/oauth-providers' })
    
    // 合并已配置的平台和支持的平台列表
    const configured = data || []
    const configuredMap = new Map<string, OAuthProvider>(configured.map((p: OAuthProvider) => [p.name, p]))
    const supportedPlatforms = getSupportedPlatforms()
    
    providers.value = supportedPlatforms.map(platform => {
      const existing = configuredMap.get(platform.name)
      return {
        id: existing?.id || 0,
        name: platform.name,
        title: platform.title,
        description: platform.description,
        region: platform.region,
        icon: getProviderIcon(platform.name),
        config: existing?.config || '{}',
        is_enabled: existing?.is_enabled || false,
      }
    })
  } catch (error) {
    console.error('获取OAuth提供商列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 打开配置对话框
const handleConfig = (provider: OAuthProvider) => {
  currentProvider.value = provider
  
  // 清空配置表单
  Object.keys(configForm).forEach(key => delete configForm[key])
  
  // 解析配置
  if (provider.config) {
    try {
      const config = JSON.parse(provider.config)
      Object.assign(configForm, config)
    } catch (e) {
      console.error('解析配置失败:', e)
    }
  }
  
  configDialogVisible.value = true
}

// 切换状态
const handleToggle = async (provider: OAuthProvider) => {
  try {
    if (provider.id) {
      // 已存在，更新状态
      await request.put({
        url: `/api/admin/oauth-providers/${provider.id}`,
        data: { is_enabled: provider.is_enabled }
      })
    } else {
      // 不存在，创建
      const { data } = await request.post({
        url: '/api/admin/oauth-providers',
        data: {
          name: provider.name,
          title: provider.title,
          description: provider.description,
          icon: provider.icon,
          is_enabled: provider.is_enabled,
        }
      })
      provider.id = data?.id || 0
    }
    ElMessage.success($t('oauthProvider.statusUpdated'))
  } catch (error) {
    provider.is_enabled = !provider.is_enabled
    ElMessage.error($t('oauthProvider.updateStatusFailed'))
  }
}

// 保存配置
const handleSaveConfig = async () => {
  submitting.value = true
  try {
    const configStr = JSON.stringify(configForm)
    
    if (currentProvider.value?.id) {
      await request.put({
        url: `/api/admin/oauth-providers/${currentProvider.value.id}`,
        data: { config: configStr }
      })
    } else {
      const { data } = await request.post({
        url: '/api/admin/oauth-providers',
        data: {
          name: currentProvider.value?.name,
          title: currentProvider.value?.title,
          description: currentProvider.value?.description,
          icon: currentProvider.value?.icon,
          config: configStr,
          is_enabled: false,
        }
      })
      if (currentProvider.value) {
        currentProvider.value.id = data?.id || 0
      }
    }
    
    ElMessage.success($t('oauthProvider.configSaved'))
    configDialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || $t('oauthProvider.saveConfigFailed'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.oauth-providers-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  color: #909399;
  font-size: 14px;
}

.provider-grid {
  margin-top: 20px;
}

.provider-card {
  margin-bottom: 20px;
  transition: all 0.3s;
}

.provider-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.provider-card.is-enabled {
  border-color: #67c23a;
}

.provider-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.provider-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
}

.provider-info {
  flex: 1;
}

.provider-info h3 {
  margin: 0 0 4px 0;
  font-size: 16px;
}

.provider-desc {
  margin-top: 12px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

.provider-actions {
  margin-top: 12px;
  text-align: right;
}

.form-hint {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
</style>
