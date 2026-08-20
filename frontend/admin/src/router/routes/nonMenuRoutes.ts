/**
 * 非菜单路由配置
 *
 * 这些路由不在菜单中显示，但需要被注册以便直接访问
 * 例如：添加客户、创建订单、查看详情等页面
 *
 * @module router/routes/nonMenuRoutes
 * @author 锚点财务团队
 */

import type { AppRouteRecordRaw } from '@/utils/router'

/**
 * 非菜单路由列表
 * 这些路由需要登录才能访问，但不在菜单中显示
 */
export const nonMenuRoutes: AppRouteRecordRaw[] = [
  // 常规设置重定向
  {
    path: '/general-settings',
    name: 'GeneralSettings',
    redirect: '/general-settings/general'
  },

  // 工单部门添加重定向
  {
    path: '/work-order-dept-add',
    name: 'WorkOrderDeptAdd',
    redirect: '/work-order-dept'
  },

  // 客户相关
  {
    path: '/add-customer',
    name: 'AddCustomer',
    component: () => import('@views/finance/clients/add/index.vue'),
    meta: { title: '添加客户', isHide: true }
  },
  {
    path: '/customer-detail/:id',
    name: 'CustomerDetail',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '客户详情', isHide: true }
  },
  {
    path: '/customer-view',
    name: 'CustomerView',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '客户详情', isHide: true }
  },
  {
    path: '/customer-view/abstract',
    name: 'CustomerViewAbstract',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '客户摘要', isHide: true }
  },
  {
    path: '/customer-view/person',
    name: 'CustomerViewPerson',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '个人资料', isHide: true }
  },
  {
    path: '/customer-view/product-list',
    name: 'CustomerViewProductList',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '产品/服务', isHide: true }
  },
  {
    path: '/customer-view/bill',
    name: 'CustomerViewBill',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '账单', isHide: true }
  },
  {
    path: '/customer-view/transactions',
    name: 'CustomerViewTransactions',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '交易记录', isHide: true }
  },
  {
    path: '/customer-view/credit',
    name: 'CustomerViewCredit',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '信用管理', isHide: true }
  },
  {
    path: '/customer-view/tickets',
    name: 'CustomerViewTickets',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '工单', isHide: true }
  },
  {
    path: '/customer-view/log',
    name: 'CustomerViewLog',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '日志', isHide: true }
  },
  {
    path: '/customer-view/api-overview',
    name: 'CustomerViewApiOverview',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: 'API概览', isHide: true }
  },
  {
    path: '/customer-view/noticelog',
    name: 'CustomerViewNoticelog',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '通知日志', isHide: true }
  },
  {
    path: '/customer-view/annex',
    name: 'CustomerViewAnnex',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '附件', isHide: true }
  },
  {
    path: '/customer-view/promotion_plan',
    name: 'CustomerViewPromotionPlan',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '推介计划', isHide: true }
  },
  {
    path: '/customer-view/follow-status',
    name: 'CustomerViewFollowStatus',
    component: () => import('@views/finance/clients/detail/index.vue'),
    meta: { title: '跟进状态', isHide: true }
  },
  {
    path: '/edit-customer/:id',
    name: 'EditCustomer',
    component: () => import('@views/finance/clients/add/index.vue'),
    meta: { title: '编辑客户', isHide: true }
  },
  {
    path: '/customer-contacts/:id',
    name: 'CustomerContacts',
    component: () => import('@views/finance/clients/contacts/index.vue'),
    meta: { title: '客户联系人', isHide: true }
  },
  {
    path: '/customer-remarks/:id',
    name: 'CustomerRemarks',
    component: () => import('@views/finance/clients/remarks/index.vue'),
    meta: { title: '客户备注', isHide: true }
  },
  {
    path: '/customer-attachments/:id',
    name: 'CustomerAttachments',
    component: () => import('@views/finance/clients/attachments/index.vue'),
    meta: { title: '客户附件', isHide: true }
  },
  {
    path: '/customer-binds/:id',
    name: 'CustomerBinds',
    component: () => import('@views/finance/clients/binds/index.vue'),
    meta: { title: '客户绑定', isHide: true }
  },
  {
    path: '/customer-track/:id',
    name: 'CustomerTrack',
    component: () => import('@views/finance/clients/track/index.vue'),
    meta: { title: '客户跟踪', isHide: true }
  },

  // 订单相关
  {
    path: '/add-order',
    name: 'AddOrder',
    component: () => import('@views/finance/orders/create/index.vue'),
    meta: { title: '添加订单', isHide: true }
  },
  {
    path: '/order-detail/:id',
    name: 'OrderDetail',
    component: () => import('@views/finance/orders/order-detail/index.vue'),
    meta: { title: '订单详情', isHide: true }
  },
  {
    path: '/edit-order/:id',
    name: 'EditOrder',
    component: () => import('@views/finance/orders/create/index.vue'),
    meta: { title: '编辑订单', isHide: true }
  },
  {
    path: '/multi-renew',
    name: 'MultiRenew',
    component: () => import('@views/finance/orders/multi-renew/index.vue'),
    meta: { title: '批量续费', isHide: true }
  },
  {
    path: '/order-upgrades/:id',
    name: 'OrderUpgrades',
    component: () => import('@views/finance/orders/upgrades/index.vue'),
    meta: { title: '订单升级', isHide: true }
  },

  // 工单相关
  {
    path: '/ticket-detail/:id',
    name: 'TicketDetail',
    component: () => import('@views/finance/tickets/detail/index.vue'),
    meta: { title: '工单详情', isHide: true }
  },
  {
    path: '/add-support-ticket',
    name: 'AddSupportTicket',
    component: () => import('@views/finance/tickets/create/index.vue'),
    meta: { title: '新建工单', isHide: true }
  },

  // 服务详情
  {
    path: '/service-details/:id',
    name: 'ServiceDetails',
    component: () => import('@views/finance/system/service-details/index.vue'),
    meta: { title: '服务详情', isHide: true }
  },

  // 发票相关
  {
    path: '/invoice-items/:id',
    name: 'InvoiceItems',
    component: () => import('@views/finance/orders/invoice-items/index.vue'),
    meta: { title: '发票项目', isHide: true }
  },

  // 新闻相关
  {
    path: '/news-list',
    name: 'NewsList',
    component: () => import('@views/finance/content/news-list/index.vue'),
    meta: { title: '新闻列表', isHide: true }
  },
  {
    path: '/news-category',
    name: 'NewsCategory',
    component: () => import('@views/finance/content/news-category/index.vue'),
    meta: { title: '新闻分类', isHide: true }
  },

  // 帮助相关
  {
    path: '/help-list',
    name: 'HelpList',
    component: () => import('@views/finance/content/help-list/index.vue'),
    meta: { title: '帮助列表', isHide: true }
  },
  {
    path: '/help-category',
    name: 'HelpCategory',
    component: () => import('@views/finance/content/help-category/index.vue'),
    meta: { title: '帮助分类', isHide: true }
  },

  // 下载相关
  {
    path: '/download-group',
    name: 'DownloadGroup',
    component: () => import('@views/finance/content/download-groups/index.vue'),
    meta: { title: '下载分组', isHide: true }
  },

  // 主题模板
  {
    path: '/theme-template',
    name: 'ThemeTemplate',
    component: () => import('@views/finance/system/theme-template/index.vue'),
    meta: { title: '主题模板', isHide: true }
  },
  {
    path: '/custom-template-fields',
    name: 'CustomTemplateFields',
    component: () => import('@views/finance/system/custom-fields/index.vue'),
    meta: { title: '自定义模板字段', isHide: true }
  },
  {
    path: '/menu_manage',
    name: 'MenuManage',
    component: () => import('@views/finance/system/navigation/index.vue'),
    meta: { title: '导航管理', isHide: true }
  },
  {
    path: '/friendly_link',
    name: 'FriendlyLink',
    component: () => import('@views/finance/content/friendly-links/index.vue'),
    meta: { title: '友情链接', isHide: true }
  },

  // 其他
  {
    path: '/balance-management',
    name: 'BalanceManagement',
    component: () => import('@views/finance/finance-core/balance/index.vue'),
    meta: { title: '余额管理', isHide: true }
  },
  {
    path: '/affiliate-records',
    name: 'AffiliateRecords',
    component: () => import('@views/finance/affiliate/records/index.vue'),
    meta: { title: '推介记录', isHide: true }
  },
  {
    path: '/password-config',
    name: 'PasswordConfig',
    component: () => import('@views/finance/tickets/password-config/index.vue'),
    meta: { title: '密码配置', isHide: true }
  },
  {
    path: '/statistics-overview',
    name: 'StatisticsOverview',
    component: () => import('@views/finance/statistics/overview/index.vue'),
    meta: { title: '统计概览', isHide: true }
  },
  {
    path: '/vouchers',
    name: 'Vouchers',
    component: () => import('@views/finance/marketing/vouchers/index.vue'),
    meta: { title: '优惠券', isHide: true }
  }
]
