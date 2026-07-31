<template>
  <div class="oauth-providers-page">
    <div class="page-header">
      <div class="header-left">
        <h2>OAuth登录管理</h2>
        <span class="subtitle">管理第三方登录平台，支持24个国内外平台</span>
      </div>
    </div>

    <!-- 平台分类 -->
    <el-tabs v-model="activeRegion">
      <el-tab-pane label="全部" name="" />
      <el-tab-pane label="国内平台" name="cn" />
      <el-tab-pane label="海外平台" name="us" />
    </el-tabs>

    <!-- 平台列表 -->
    <el-row :gutter="20" class="provider-grid">
      <el-col :span="6" v-for="provider in filteredProviders" :key="provider.name">
        <el-card class="provider-card" :class="{ 'is-enabled': provider.is_enabled }">
          <div class="provider-header">
            <img :src="getProviderIcon(provider.name)" class="provider-icon" @error="handleIconError($event, provider.name)" />
            <div class="provider-info">
              <h3>{{ provider.title }}</h3>
              <el-tag :type="provider.region === 'cn' ? '' : 'info'" size="small">
                {{ provider.region === 'cn' ? '国内' : '海外' }}
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
              配置
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      :title="`${currentProvider?.title || ''} - 配置`"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form label-width="120px">
        <!-- 通用配置 -->
        <el-form-item label="App ID / Key">
          <el-input v-model="configForm.app_id" placeholder="应用ID" />
        </el-form-item>
        <el-form-item label="App Secret">
          <el-input v-model="configForm.app_secret" type="password" show-password placeholder="应用密钥" />
        </el-form-item>
        <el-form-item label="回调地址">
          <el-input v-model="configForm.callback_url" :placeholder="getCallbackPlaceholder(currentProvider?.name)" />
          <div class="form-hint">请在第三方平台填写此回调地址</div>
        </el-form-item>

        <!-- 特定平台配置 -->
        <template v-if="currentProvider?.name === 'wechat'">
          <el-divider content-position="left">微信特有配置</el-divider>
          <el-form-item label="开放平台类型">
            <el-select v-model="configForm.platform">
              <el-option label="网站应用" value="website" />
              <el-option label="移动应用" value="mobile" />
              <el-option label="公众号" value="mp" />
            </el-select>
          </el-form-item>
        </template>

        <template v-if="currentProvider?.name === 'dingtalk'">
          <el-divider content-position="left">钉钉特有配置</el-divider>
          <el-form-item label="应用类型">
            <el-select v-model="configForm.app_type">
              <el-option label="企业内部应用" value="internal" />
              <el-option label="第三方企业应用" value="third" />
              <el-option label="第三方个人应用" value="personal" />
            </el-select>
          </el-form-item>
        </template>

        <template v-if="currentProvider?.name === 'apple'">
          <el-divider content-position="left">Apple特有配置</el-divider>
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
          <el-divider content-position="left">Telegram特有配置</el-divider>
          <el-form-item label="Bot Token">
            <el-input v-model="configForm.bot_token" placeholder="从 @BotFather 获取" />
          </el-form-item>
          <el-form-item label="Bot Username">
            <el-input v-model="configForm.bot_username" placeholder="MyBot" />
          </el-form-item>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveConfig" :loading="submitting">保存配置</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
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
const supportedPlatforms = [
  // 国内平台
  { name: 'wechat', title: '微信登录', region: 'cn', description: '微信扫码登录，支持网站应用、移动应用、公众号' },
  { name: 'qq', title: 'QQ登录', region: 'cn', description: 'QQ互联登录，支持PC和移动端' },
  { name: 'weibo', title: '微博登录', region: 'cn', description: '新浪微博登录' },
  { name: 'alipay', title: '支付宝登录', region: 'cn', description: '支付宝账号登录' },
  { name: 'baidu', title: '百度登录', region: 'cn', description: '百度账号登录' },
  { name: 'gitee', title: '码云登录', region: 'cn', description: 'Gitee账号登录，适合开发者' },
  { name: 'dingtalk', title: '钉钉登录', region: 'cn', description: '钉钉扫码登录，支持企业内部应用' },
  { name: 'feishu', title: '飞书登录', region: 'cn', description: '飞书扫码登录' },
  { name: 'csdn', title: 'CSDN登录', region: 'cn', description: 'CSDN账号登录' },
  { name: 'oschina', title: '开源中国登录', region: 'cn', description: 'OSChina账号登录' },
  { name: 'tencent_cloud', title: '腾讯云登录', region: 'cn', description: '腾讯云账号登录' },
  { name: 'aliyun', title: '阿里云登录', region: 'cn', description: '阿里云账号登录' },
  
  // 海外平台
  { name: 'google', title: 'Google登录', region: 'us', description: 'Google账号登录，支持Gmail' },
  { name: 'facebook', title: 'Facebook登录', region: 'us', description: 'Facebook账号登录' },
  { name: 'twitter', title: 'Twitter/X登录', region: 'us', description: 'Twitter或X账号登录' },
  { name: 'github', title: 'GitHub登录', region: 'us', description: 'GitHub账号登录，适合开发者' },
  { name: 'linkedin', title: 'LinkedIn登录', region: 'us', description: 'LinkedIn职业社交账号登录' },
  { name: 'microsoft', title: 'Microsoft登录', region: 'us', description: 'Microsoft/Outlook账号登录' },
  { name: 'apple', title: 'Apple登录', region: 'us', description: 'Apple ID登录，iOS应用必备' },
  { name: 'amazon', title: 'Amazon登录', region: 'us', description: 'Amazon账号登录' },
  { name: 'discord', title: 'Discord登录', region: 'us', description: 'Discord账号登录，适合游戏社区' },
  { name: 'slack', title: 'Slack登录', region: 'us', description: 'Slack工作区登录' },
  { name: 'telegram', title: 'Telegram登录', region: 'us', description: 'Telegram账号登录' },
  { name: 'line', title: 'LINE登录', region: 'us', description: 'LINE账号登录，日韩常用' },
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
    const configuredMap = new Map(configured.map((p: OAuthProvider) => [p.name, p]))
    
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
    ElMessage.success('状态已更新')
  } catch (error) {
    provider.is_enabled = !provider.is_enabled
    ElMessage.error('更新状态失败')
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
    
    ElMessage.success('配置已保存')
    configDialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '保存配置失败')
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
