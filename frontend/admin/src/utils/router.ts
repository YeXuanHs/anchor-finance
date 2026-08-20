/**
 * 路由工具函数
 *
 * 提供路由相关的工具函数
 *
 * @module utils/router
 */
import { RouteLocationNormalized, RouteRecordRaw } from 'vue-router'
import AppConfig from '@/config'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import i18n, { $t } from '@/locales'

/** 扩展的路由配置类型 */
export type AppRouteRecordRaw = RouteRecordRaw & {
  hidden?: boolean
}

/** 顶部进度条配置 */
export const configureNProgress = () => {
  NProgress.configure({
    easing: 'ease',
    speed: 600,
    showSpinner: false,
    parent: 'body'
  })
}

/**
 * 设置页面标题，根据路由元信息和系统信息拼接标题
 * @param to 当前路由对象
 */
export const setPageTitle = (to: RouteLocationNormalized): void => {
  const { title, languageMap } = to.meta as { title?: string; languageMap?: Record<string, string> }
  if (title) {
    setTimeout(() => {
      const displayTitle = languageMap ? getMenuTitle({ title, languageMap }) : formatMenuTitle(String(title))
      document.title = `${displayTitle} - ${AppConfig.systemInfo.name}`
    }, 150)
  }
}

/**
 * 格式化菜单标题
 * @param title 菜单标题，可以是 i18n 的 key，也可以是字符串
 * @returns 格式化后的菜单标题
 */
export const formatMenuTitle = (title: string): string => {
  if (title) {
    if (title.startsWith('menus.')) {
      // 使用 te() 方法检查翻译键值是否存在，避免控制台警告
      if (i18n.global.te(title)) {
        return $t(title)
      } else {
        // 如果翻译不存在，返回键值的最后部分作为fallback
        return title.split('.').pop() || title
      }
    }
    return title
  }
  return ''
}

/**
 * 获取菜单标题（支持languageMap多语言）
 * @param meta 菜单元数据，包含 title 和 languageMap
 * @param currentLang 当前语言（可选，由调用方传入以保证Vue响应式追踪）
 * @returns 格式化后的菜单标题
 */
export const getMenuTitle = (meta: { title?: string; languageMap?: Record<string, string> }, currentLang?: string): string => {
  if (!meta) return ''
  
  // 如果有languageMap，根据当前语言返回对应的翻译
  if (meta.languageMap) {
    // 如果调用方没有传入 currentLang，从 localStorage 获取
    let lang = currentLang || 'zh'
    if (!currentLang) {
      try {
        const keys = Object.keys(localStorage)
        const userKey = keys.find(k => k.endsWith('-user'))
        if (userKey) {
          const userData = JSON.parse(localStorage.getItem(userKey) || '{}')
          if (userData.language) {
            lang = userData.language
          }
        }
      } catch {
        // fallback to $t-based detection
        const testResult = $t('topBar.search.title')
        lang = testResult === 'Search' ? 'en' : testResult === '搜尋' ? 'zh-TW' : 'zh'
      }
    }
    
    // 按优先级尝试多种 key 格式（兼容 zjmf 的 CN/HK/US 格式和 zh/zh-TW/en 格式）
    const candidates: string[] = []
    if (lang === 'zh' || lang === 'zh-CN' || lang === 'zh_cn') {
      candidates.push('zh-CN', 'CN', 'zh', 'zh_cn')
    } else if (lang === 'zh-TW' || lang === 'zh_tw' || lang === 'HK') {
      candidates.push('zh-TW', 'HK', 'zh_tw', 'zh')
    } else if (lang === 'en' || lang === 'US') {
      candidates.push('en', 'US')
    } else {
      candidates.push(lang)
    }
    
    for (const key of candidates) {
      if (meta.languageMap[key]) {
        return meta.languageMap[key]
      }
    }
  }
  
  // 否则使用原有的formatMenuTitle逻辑
  return formatMenuTitle(meta.title || '')
}
