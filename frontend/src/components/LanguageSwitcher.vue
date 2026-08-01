<template>
  <el-dropdown @command="handleCommand" trigger="click">
    <span class="lang-switcher">
      <el-icon><Promotion /></el-icon>
      <span>{{ currentLangName }}</span>
      <el-icon class="el-icon--right"><ArrowDown /></el-icon>
    </span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="lang in languages"
          :key="lang.code"
          :command="lang.code"
          :class="{ active: lang.code === currentLang }"
        >
          <span class="lang-flag">{{ lang.flag }}</span>
          <span>{{ lang.name }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

interface Language {
  code: string
  name: string
  flag: string
}

const languages = ref<Language[]>([
  { code: 'zh-CN', name: '中文简体', flag: 'CN' },
  { code: 'en-US', name: 'English', flag: 'US' },
  { code: 'zh-TW', name: '中文繁體', flag: 'TW' }
])

const currentLang = computed(() => locale.value)
const currentLangName = computed(() => {
  const lang = languages.value.find(l => l.code === locale.value)
  return lang ? lang.name : '中文简体'
})

const handleCommand = (code: string) => {
  locale.value = code
  localStorage.setItem('language', code)
}

onMounted(async () => {
  const saved = localStorage.getItem('language')
  if (saved) {
    locale.value = saved
  }
  
  // 从后端获取可用语言
  try {
    const { default: request } = await import('@/utils/http')
    const res = await request.get('/api/v1/languages')
    if (res.data?.data?.length) {
      languages.value = res.data.data.map((l: any) => ({
        code: l.code,
        name: l.name,
        flag: l.flag
      }))
    }
  } catch (e) {
    // 使用默认语言列表
  }
})
</script>

<style scoped>
.lang-switcher {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.lang-switcher:hover {
  color: var(--el-color-primary);
}

.lang-flag {
  margin-right: 8px;
}

.el-dropdown-item.active {
  color: var(--el-color-primary);
  font-weight: bold;
}
</style>
