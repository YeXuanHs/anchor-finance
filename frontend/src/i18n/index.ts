import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'
import zhTW from './zh-TW'

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('language') || 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
    'zh-TW': zhTW
  }
})

// 切换语言
export async function switchLanguage(langCode: string) {
  i18n.global.locale.value = langCode as any
  localStorage.setItem('language', langCode)
  document.documentElement.lang = langCode
}

export default i18n
