/**
 * 类型安全的标签工具函数
 *
 * 解决 Element Plus el-tag type 属性的 TypeScript 类型问题
 *
 * @module utils/tag
 */

type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'

/**
 * 获取状态标签类型
 * 返回 Element Plus 兼容的标签类型
 */
export function getStatusTagType(status: string): TagType {
  const map: Record<string, TagType> = {
    active: 'success',
    inactive: 'info',
    disabled: 'danger',
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    completed: 'success',
    cancelled: 'info',
    paid: 'success',
    unpaid: 'warning',
    refunded: 'danger',
    open: 'warning',
    closed: 'info',
    processing: 'primary',
    success: 'success',
    failed: 'danger',
    error: 'danger',
    warning: 'warning',
    info: 'info'
  }
  return map[status] || 'info'
}

/**
 * 获取交易类型标签
 */
export function getTransactionTypeTag(type: string): TagType {
  const map: Record<string, TagType> = {
    income: 'success',
    expense: 'danger',
    refund: 'warning'
  }
  return map[type] || 'info'
}

/**
 * 获取支付方式标签
 */
export function getPaymentMethodTag(method: string): TagType {
  const map: Record<string, TagType> = {
    balance: 'primary',
    offline: 'warning',
    online: 'success',
    free: 'info'
  }
  return map[method] || 'info'
}

/**
 * 获取优先级标签
 */
export function getPriorityTag(priority: number | string): TagType {
  const map: Record<string, TagType> = {
    1: 'info',
    2: 'primary',
    3: 'warning',
    4: 'danger',
    low: 'info',
    normal: 'primary',
    high: 'warning',
    urgent: 'danger'
  }
  return map[String(priority)] || 'info'
}
