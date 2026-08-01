<template>
  <el-dropdown @command="changeLanguage" trigger="click">
    <span class="language-switch">
      <el-icon><Connection /></el-icon>
      <span class="lang-text">{{ currentLangName }}</span>
    </span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item 
          v-for="lang in languages" 
          :key="lang.code" 
          :command="lang.code"
          :class="{ 'is-active': lang.code === currentLang }"
        >
          <span class="lang-flag">{{ lang.flag }}</span>
          <span>{{ lang.name }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Connection } from '@element-plus/icons-vue'

const { locale } = useI18n()

interface Language {
  code: string
  name: string
  flag: string
}

const languages: Language[] = [
  { code: 'zh-CN', name: '中文简体', flag: 'CN' },
  { code: 'en-US', name: 'English', flag: 'US' },
  { code: 'zh-TW', name: '中文繁體', flag: 'TW' }
]

const currentLang = computed(() => locale.value)

const currentLangName = computed(() => {
  const lang = languages.find(l => l.code === currentLang.value)
  return lang?.name || '中文简体'
})

const changeLanguage = (langCode: string) => {
  locale.value = langCode
  localStorage.setItem('language', langCode)
  document.documentElement.lang = langCode
}
</script>

<style scoped>
.language-switch {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.language-switch:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.lang-text {
  font-size: 14px;
}

.lang-flag {
  margin-right: 8px;
  font-weight: bold;
}

.is-active {
  color: var(--el-color-primary);
}
</style>
