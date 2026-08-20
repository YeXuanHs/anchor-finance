/**
 * 菜单处理器
 *
 * 负责菜单数据的获取、过滤和处理
 *
 * @module router/core/MenuProcessor
 * @author 锚点财务团队
 */

import type { AppRouteRecord } from '@/types/router'
import { fetchGetMenuList } from '@/api/system-manage'
import { RoutesAlias } from '../routesAlias'
import { formatMenuTitle } from '@/utils'
import { nonMenuRoutes } from '../routes/nonMenuRoutes'

export class MenuProcessor {
  // 存储原始菜单树（用于侧边栏渲染）
  private rawMenuTree: any[] = []

  /**
   * 获取菜单数据（扁平路由，用于路由器注册）
   */
  async getMenuList(): Promise<AppRouteRecord[]> {
    // 后端模式：从 API 获取菜单
    const menuList = await this.processBackendMenu()

    // 在规范化路径之前，验证原始路径配置
    this.validateMenuPaths(menuList)

    // 规范化路径（将相对路径转换为完整路径）
    return this.normalizeMenuPaths(menuList)
  }

  /**
   * 获取菜单树（层级结构，用于侧边栏渲染）
   */
  getMenuTree(): any[] {
    return this.rawMenuTree
  }

  /**
   * 处理后端控制模式的菜单
   * 将后端菜单格式 {url, name, icon} 转换为前端路由格式 {path, name, component, meta}
   * 生成扁平路由结构，所有叶子节点直接挂在 Layout 下
   */
  private async processBackendMenu(): Promise<AppRouteRecord[]> {
    const list = await fetchGetMenuList()

    // 保存原始菜单树（用于侧边栏渲染）
    this.rawMenuTree = list

    // 提取所有叶子节点（有 URL 的菜单项）作为扁平路由
    const leafRoutes = this.extractLeafRoutes(list)

    // 添加非菜单路由（如添加客户、订单详情等页面）
    // 这些路由不在菜单中显示，但需要在Layout下注册以获得完整的布局
    const nonMenuRouteNames = new Set(leafRoutes.map(r => r.name))
    console.log('[MenuProcessor] Adding nonMenuRoutes, current leafRoutes count:', leafRoutes.length)
    for (const route of nonMenuRoutes) {
      if (route.name && !nonMenuRouteNames.has(route.name)) {
        console.log('[MenuProcessor] Adding nonMenuRoute:', route.name, route.path)
        leafRoutes.push({
          path: route.path,
          name: route.name,
          component: (route.component || '') as any,
          redirect: (route as any).redirect,
          meta: route.meta || {}
        } as AppRouteRecord)
      }
    }
    console.log('[MenuProcessor] After adding nonMenuRoutes, leafRoutes count:', leafRoutes.length)

    // 将叶子节点包装在 Layout 下，形成扁平结构
    const layoutRoute: AppRouteRecord = {
      path: '/',
      name: 'Layout',
      component: '/index/index' as any,
      meta: { title: '', icon: '' },
      children: leafRoutes
    }

    return this.filterEmptyMenus([layoutRoute])
  }

  /**
   * 从菜单树中提取所有叶子节点（有 URL 的菜单项）
   */
  private extractLeafRoutes(menus: any[]): AppRouteRecord[] {
    const routes: AppRouteRecord[] = []

    for (const item of menus) {
      // is_visible = 0 的菜单不跳过（需要注册路由），只跳过明确禁用的
      // is_active 用于旧格式兼容
      if (item.is_active === false) continue

      // 兼容新旧格式：新格式用 path，旧格式用 url
      const itemUrl = item.path || item.url || ''

      if (itemUrl && (!item.children || item.children.length === 0)) {
        // 叶子节点：创建路由
        const cleanUrl = itemUrl.replace(/^\/finance/, '') || '/'
        const routeName = this.generateRouteName(cleanUrl)
        routes.push({
          path: cleanUrl,
          name: routeName,
          component: itemUrl as any,
          meta: {
            title: item.meta?.title || item.name || '',
            icon: item.meta?.icon || item.icon || '',
            order: item.sort_order || 0,
            isHide: item.meta?.isHide ?? (item.is_visible === false),
            languageMap: item.meta?.languageMap || item.language_map
          }
        })
      } else if (item.children && item.children.length > 0) {
        // 有子节点：递归提取
        routes.push(...this.extractLeafRoutes(item.children))
      }
    }

    return routes
  }

  /**
   * 从 URL 生成路由名称
   * /customer-list -> CustomerList
   */
  private generateRouteName(url: string): string {
    return url
      .replace(/^\//, '')
      .split('/')
      .filter(Boolean)
      .map((s: string) => s.charAt(0).toUpperCase() + s.slice(1))
      .join('')
  }

  /**
   * 递归过滤空菜单项
   */
  private filterEmptyMenus(menuList: AppRouteRecord[]): AppRouteRecord[] {
    return menuList
      .map((item) => {
        // 如果有子菜单，先递归过滤子菜单
        if (item.children && item.children.length > 0) {
          const filteredChildren = this.filterEmptyMenus(item.children)
          return {
            ...item,
            children: filteredChildren
          }
        }
        return item
      })
      .filter((item) => {
        // 如果定义了 children 属性（即使是空数组），说明这是一个目录菜单，应该保留
        if ('children' in item) {
          return true
        }

        // 如果有外链或 iframe，保留
        if (item.meta?.isIframe === true || item.meta?.link) {
          return true
        }

        // 如果有有效的 component，保留
        if (item.component && item.component !== '' && item.component !== RoutesAlias.Layout) {
          return true
        }

        // 其他情况过滤掉
        return false
      })
  }

  /**
   * 验证菜单列表是否有效
   */
  validateMenuList(menuList: AppRouteRecord[]): boolean {
    return Array.isArray(menuList) && menuList.length > 0
  }

  /**
   * 规范化菜单路径
   * 将相对路径转换为完整路径，确保菜单跳转正确
   */
  private normalizeMenuPaths(menuList: AppRouteRecord[], parentPath = ''): AppRouteRecord[] {
    return menuList.map((item) => {
      // 构建完整路径
      const fullPath = this.buildFullPath(item.path || '', parentPath)

      // 递归处理子菜单
      const children = item.children?.length
        ? this.normalizeMenuPaths(item.children, fullPath)
        : item.children

      const redirect = item.redirect || this.resolveDefaultRedirect(children)

      return {
        ...item,
        path: fullPath,
        redirect,
        children
      }
    })
  }

  /**
   * 为目录型菜单推导默认跳转地址
   */
  private resolveDefaultRedirect(children?: AppRouteRecord[]): string | undefined {
    if (!children?.length) {
      return undefined
    }

    for (const child of children) {
      if (this.isNavigableRoute(child)) {
        return child.path
      }

      const nestedRedirect = this.resolveDefaultRedirect(child.children)
      if (nestedRedirect) {
        return nestedRedirect
      }
    }

    return undefined
  }

  /**
   * 判断子路由是否可以作为默认落点
   */
  private isNavigableRoute(route: AppRouteRecord): boolean {
    return Boolean(
      route.path &&
        route.path !== '/' &&
        !route.meta?.link &&
        route.meta?.isIframe !== true &&
        route.component &&
        route.component !== ''
    )
  }

  /**
   * 验证菜单路径配置
   * 检测非一级菜单是否错误使用了 / 开头的路径
   */
  /**
   * 验证菜单路径配置
   * 检测非一级菜单是否错误使用了 / 开头的路径
   */
  private validateMenuPaths(menuList: AppRouteRecord[], level = 1): void {
    menuList.forEach((route) => {
      if (!route.children?.length) return

      const parentName = String(route.name || route.path || '未知路由')

      route.children.forEach((child) => {
        const childPath = child.path || ''

        // 跳过合法的绝对路径：外部链接和 iframe 路由
        if (this.isValidAbsolutePath(childPath)) return

        // 检测非法的绝对路径
        if (childPath.startsWith('/')) {
          this.logPathError(child, childPath, parentName, level)
        }
      })

      // 递归检查更深层级的子路由
      this.validateMenuPaths(route.children, level + 1)
    })
  }

  /**
   * 判断是否为合法的绝对路径
   */
  private isValidAbsolutePath(path: string): boolean {
    return (
      path.startsWith('http://') ||
      path.startsWith('https://') ||
      path.startsWith('/outside/iframe/')
    )
  }

  /**
   * 输出路径配置错误日志
   */
  private logPathError(
    route: AppRouteRecord,
    path: string,
    parentName: string,
    level: number
  ): void {
    const routeName = String(route.name || path || '未知路由')
    const menuTitle = route.meta?.title || routeName
    const suggestedPath = path.split('/').pop() || path.slice(1)

    console.error(
      `[路由配置错误] 菜单 "${formatMenuTitle(menuTitle)}" (name: ${routeName}, path: ${path}) 配置错误\n` +
        `  位置: ${parentName} > ${routeName}\n` +
        `  问题: ${level + 1}级菜单的 path 不能以 / 开头\n` +
        `  当前配置: path: '${path}'\n` +
        `  应该改为: path: '${suggestedPath}'`
    )
  }

  /**
   * 构建完整路径
   */
  private buildFullPath(path: string, parentPath: string): string {
    if (!path) return ''

    // 外部链接直接返回
    if (path.startsWith('http://') || path.startsWith('https://')) {
      return path
    }

    // 如果已经是绝对路径，直接返回
    if (path.startsWith('/')) {
      return path
    }

    // 拼接父路径和当前路径
    if (parentPath) {
      // 移除父路径末尾的斜杠，移除子路径开头的斜杠，然后拼接
      const cleanParent = parentPath.replace(/\/$/, '')
      const cleanChild = path.replace(/^\//, '')
      return `${cleanParent}/${cleanChild}`
    }

    // 没有父路径，添加前导斜杠
    return `/${path}`
  }
}
