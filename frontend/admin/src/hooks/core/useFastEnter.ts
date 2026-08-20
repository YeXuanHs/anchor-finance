/**
 * useFastEnter - 快速入口管理
 *
 * 管理顶部栏的快速入口功能，提供应用列表和快速链接的配置和过滤。
 * 支持动态启用/禁用、自定义排序、响应式宽度控制等功能。
 *
 * ## 主要功能
 *
 * 1. 应用列表管理 - 获取启用的应用列表，自动按排序权重排序
 * 2. 快速链接管理 - 获取启用的快速链接，支持自定义排序
 * 3. 响应式配置 - 所有配置自动响应变化，无需手动更新
 * 4. 宽度控制 - 提供最小显示宽度配置，支持响应式布局
 * 5. 多语言支持 - 应用名称和描述支持 i18n 翻译
 *
 * @module useFastEnter
 * @author 锚点财务团队
 */

import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import appConfig from '@/config'
import type { FastEnterApplication, FastEnterQuickLink } from '@/types/config'

// 应用名称到 i18n key 的映射
const appNameKeyMap: Record<string, string> = {
  '工作台': 'fastEnter.workspace',
  '分析页': 'fastEnter.analysis',
  '更新日志': 'fastEnter.changeLog',
  '开发者GitHub': 'fastEnter.github',
  'QQ交流群': 'fastEnter.qqGroup'
}

// 应用描述到 i18n key 的映射
const appDescKeyMap: Record<string, string> = {
  '系统概览与数据统计': 'fastEnter.workspaceDesc',
  '数据分析与可视化': 'fastEnter.analysisDesc',
  '版本更新与变更记录': 'fastEnter.changeLogDesc',
  '项目源码与反馈': 'fastEnter.githubDesc',
  '加入QQ群交流': 'fastEnter.qqGroupDesc'
}

// 快速链接名称到 i18n key 的映射
const quickLinkNameKeyMap: Record<string, string> = {
  '登录': 'fastEnter.login',
  '个人中心': 'fastEnter.userCenter'
}

export function useFastEnter() {
  const { t } = useI18n()

  // 获取快速入口配置
  const fastEnterConfig = computed(() => appConfig.fastEnter)

  // 获取启用的应用列表（按排序权重排序，带翻译）
  const enabledApplications = computed<FastEnterApplication[]>(() => {
    if (!fastEnterConfig.value?.applications) return []

    return fastEnterConfig.value.applications
      .filter((app) => app.enabled !== false)
      .sort((a, b) => (a.order || 0) - (b.order || 0))
      .map((app) => ({
        ...app,
        name: appNameKeyMap[app.name] ? t(appNameKeyMap[app.name]) : app.name,
        description: appDescKeyMap[app.description || ''] ? t(appDescKeyMap[app.description || '']) : app.description
      }))
  })

  // 获取启用的快速链接（按排序权重排序，带翻译）
  const enabledQuickLinks = computed<FastEnterQuickLink[]>(() => {
    if (!fastEnterConfig.value?.quickLinks) return []

    return fastEnterConfig.value.quickLinks
      .filter((link) => link.enabled !== false)
      .sort((a, b) => (a.order || 0) - (b.order || 0))
      .map((link) => ({
        ...link,
        name: quickLinkNameKeyMap[link.name] ? t(quickLinkNameKeyMap[link.name]) : link.name
      }))
  })

  // 获取最小显示宽度
  const minWidth = computed(() => {
    return fastEnterConfig.value?.minWidth || 1200
  })

  return {
    fastEnterConfig,
    enabledApplications,
    enabledQuickLinks,
    minWidth
  }
}
