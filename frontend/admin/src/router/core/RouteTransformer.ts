/**
 * 路由转换器
 *
 * 负责将菜单数据转换为 Vue Router 路由配置
 *
 * @module router/core/RouteTransformer
 * @author 锚点财务团队
 */

import type { RouteRecordRaw } from 'vue-router'
import type { AppRouteRecord } from '@/types/router'
import { ComponentLoader } from './ComponentLoader'
import { IframeRouteManager } from './IframeRouteManager'

interface ConvertedRoute extends Omit<RouteRecordRaw, 'children'> {
  id?: number
  children?: ConvertedRoute[]
  component?: RouteRecordRaw['component'] | (() => Promise<any>)
}

export class RouteTransformer {
  private componentLoader: ComponentLoader
  private iframeManager: IframeRouteManager

  constructor(componentLoader: ComponentLoader) {
    this.componentLoader = componentLoader
    this.iframeManager = IframeRouteManager.getInstance()
  }

  /**
   * 转换路由配置
   */
  transform(route: AppRouteRecord, depth = 0): ConvertedRoute {
    const { component, children, ...routeConfig } = route

    console.log('[RouteTransformer] Transforming route:', route.name, route.path, 'depth:', depth)

    // 基础路由配置
    const converted: ConvertedRoute = {
      ...routeConfig,
      component: undefined
    }

    // 处理不同类型的路由
    if (route.meta.isIframe) {
      this.handleIframeRoute(converted, route, depth)
    } else if (this.isFirstLevelRoute(route, depth)) {
      console.log('[RouteTransformer] Handling as first level route:', route.name)
      this.handleFirstLevelRoute(converted, route, component as string)
    } else {
      console.log('[RouteTransformer] Handling as normal route:', route.name)
      this.handleNormalRoute(converted, component as string)
    }

    // 递归处理子路由
    if (children?.length) {
      converted.children = children.map((child) => this.transform(child, depth + 1))
    }

    return converted
  }

  /**
   * 判断是否为一级路由（需要 Layout 包裹）
   */
  private isFirstLevelRoute(route: AppRouteRecord, depth: number): boolean {
    // redirect-only 路由（有 redirect 但没有 component 和 children）不需要 Layout 包裹
    if (depth === 0 && (route as any).redirect && !route.component && (!route.children || route.children.length === 0)) {
      return false
    }
    return depth === 0 && (!route.children || route.children.length === 0)
  }

  /**
   * 处理 iframe 类型路由
   */
  private handleIframeRoute(
    targetRoute: ConvertedRoute,
    sourceRoute: AppRouteRecord,
    depth: number
  ): void {
    if (depth === 0) {
      // 顶级 iframe：用 Layout 包裹
      targetRoute.component = this.componentLoader.loadLayout()
      targetRoute.path = this.extractFirstSegment(sourceRoute.path || '')
      targetRoute.name = ''

      targetRoute.children = [
        {
          ...sourceRoute,
          component: this.componentLoader.loadIframe()
        } as ConvertedRoute
      ]
    } else {
      // 非顶级（嵌套）iframe：直接使用 Iframe.vue
      targetRoute.component = this.componentLoader.loadIframe()
    }

    // 记录 iframe 路由
    this.iframeManager.add(sourceRoute)
  }

  /**
   * 处理一级菜单路由
   */
  private handleFirstLevelRoute(
    converted: ConvertedRoute,
    route: AppRouteRecord,
    component: string | undefined | (() => Promise<any>)
  ): void {
    converted.component = this.componentLoader.loadLayout()
    const routePath = route.path || ''

    // 处理函数组件（如 nonMenuRoutes 中的动态导入）和字符串路径
    let resolvedComponent: (() => Promise<any>) | undefined
    if (typeof component === 'function') {
      resolvedComponent = component as () => Promise<any>
    } else if (component) {
      resolvedComponent = this.componentLoader.load(component)
    }

    // 判断是否为简单路径（如 /add-order, /customer-add）
    if (this.isSimplePath(routePath)) {
      // 简单路径：直接使用完整路径作为子路由
      converted.path = '/'
      converted.name = ''
      route.meta.isFirstLevel = true

      converted.children = [
        {
          ...route,
          component: resolvedComponent
        } as ConvertedRoute
      ]
    } else {
      // 复杂路径（如 /customer-view/abstract）：拆分为父路径和子路径
      const firstSegment = this.extractFirstSegment(routePath)
      converted.path = firstSegment
      converted.name = ''
      route.meta.isFirstLevel = true

      // 计算子路由的相对路径
      const childPath = routePath.substring(firstSegment.length)

      converted.children = [
        {
          ...route,
          path: childPath || '',
          component: resolvedComponent
        } as ConvertedRoute
      ]
    }
  }

  /**
   * 处理普通路由
   */
  private handleNormalRoute(converted: ConvertedRoute, component: string | undefined | (() => Promise<any>)): void {
    if (component) {
      // 如果component是函数（动态导入），直接使用
      if (typeof component === 'function') {
        converted.component = component as () => Promise<any>
      } else {
        converted.component = this.componentLoader.load(component)
      }
    }
  }

  /**
   * 提取路径的第一段
   */
  private extractFirstSegment(path: string): string {
    const segments = path.split('/').filter(Boolean)
    return segments.length > 0 ? `/${segments[0]}` : '/'
  }

  /**
   * 判断是否为简单路径（单段路径，如 /add-order）
   */
  private isSimplePath(path: string): boolean {
    const segments = path.split('/').filter(Boolean)
    return segments.length === 1
  }
}
