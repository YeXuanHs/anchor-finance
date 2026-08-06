import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/Home.vue')
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue')
    },
    {
      path: '/forgot-password',
      name: 'ForgotPassword',
      component: () => import('@/views/ForgotPassword.vue')
    },
    {
      path: '/products',
      name: 'Products',
      component: () => import('@/views/Products.vue')
    },
    {
      path: '/products/cloud',
      name: 'ProductsCloud',
      component: () => import('@/views/products/Cloud.vue')
    },
    {
      path: '/products/dedicated',
      name: 'ProductsDedicated',
      component: () => import('@/views/products/Dedicated.vue')
    },
    {
      path: '/products/hosting',
      name: 'ProductsHosting',
      component: () => import('@/views/products/Hosting.vue')
    },
    {
      path: '/products/domain',
      name: 'ProductsDomain',
      component: () => import('@/views/products/Domain.vue')
    },
    {
      path: '/products/ssl',
      name: 'ProductsSSL',
      component: () => import('@/views/products/SSL.vue')
    },
    {
      path: '/products/nat',
      name: 'ProductsNAT',
      component: () => import('@/views/products/NAT.vue')
    },
    {
      path: '/products/antiddos',
      name: 'ProductsAntiDDoS',
      component: () => import('@/views/products/AntiDDoS.vue')
    },
    {
      path: '/products/cdn',
      name: 'ProductsCDN',
      component: () => import('@/views/products/CDN.vue')
    },
    {
      path: '/products/:id',
      name: 'ProductDetail',
      component: () => import('@/views/ProductDetail.vue')
    },
    {
      path: '/cart',
      name: 'Cart',
      component: () => import('@/views/Cart.vue')
    },
    {
      path: '/checkout',
      name: 'Checkout',
      component: () => import('@/views/Checkout.vue')
    },
    {
      path: '/payment-result',
      name: 'PaymentResult',
      component: () => import('@/views/PaymentResult.vue')
    },
    {
      path: '/about',
      name: 'About',
      component: () => import('@/views/About.vue')
    },
    {
      path: '/contact',
      name: 'Contact',
      component: () => import('@/views/Contact.vue')
    },
    {
      path: '/help',
      name: 'Help',
      component: () => import('@/views/Help.vue')
    },
    {
      path: '/help/article/:id',
      name: 'HelpArticle',
      component: () => import('@/views/HelpArticle.vue')
    },
    {
      path: '/help/category/:id',
      name: 'HelpCategory',
      component: () => import('@/views/help/HelpCategory.vue')
    },
    {
      path: '/help/content/:id',
      name: 'HelpContent',
      component: () => import('@/views/help/HelpContent.vue')
    },
    {
      path: '/help/cate',
      name: 'HelpCate',
      component: () => import('@/views/help/HelpCate.vue')
    },
    {
      path: '/help/search',
      name: 'HelpSearch',
      component: () => import('@/views/help/HelpSearch.vue')
    },
    {
      path: '/knowledge-base',
      name: 'KnowledgeBase',
      component: () => import('@/views/KnowledgeBase.vue')
    },
    {
      path: '/downloads',
      name: 'Downloads',
      component: () => import('@/views/DownloadCenter.vue')
    },
    {
      path: '/document',
      name: 'Document',
      component: () => import('@/views/Document.vue')
    },
    {
      path: '/news',
      name: 'News',
      component: () => import('@/views/news/NewsList.vue')
    },
    {
      path: '/news/category/:id',
      name: 'NewsCategory',
      component: () => import('@/views/news/NewsCategory.vue')
    },
    {
      path: '/news/search',
      name: 'NewsSearch',
      component: () => import('@/views/news/NewsSearch.vue')
    },
    {
      path: '/news/:id',
      name: 'NewsDetail',
      component: () => import('@/views/NewsDetail.vue')
    },
    {
      path: '/solutions',
      name: 'Solutions',
      component: () => import('@/views/Solutions.vue')
    },
    {
      path: '/solutions/:id',
      name: 'SolutionDetail',
      component: () => import('@/views/SolutionDetail.vue')
    },
    {
      path: '/solutions/game',
      name: 'SolutionGame',
      component: () => import('@/views/solutions/game.vue')
    },
    {
      path: '/solutions/video',
      name: 'SolutionVideo',
      component: () => import('@/views/solutions/video.vue')
    },
    {
      path: '/solutions/edu',
      name: 'SolutionEdu',
      component: () => import('@/views/solutions/edu.vue')
    },
    {
      path: '/solutions/ecommerce',
      name: 'SolutionEcommerce',
      component: () => import('@/views/solutions/ecommerce.vue')
    },
    {
      path: '/solutions/security',
      name: 'SolutionSecurity',
      component: () => import('@/views/solutions/security.vue')
    },
    {
      path: '/solutions/caredisaster',
      name: 'SolutionCaredisaster',
      component: () => import('@/views/solutions/caredisaster.vue')
    },
    {
      path: '/solutions/mixedcloud',
      name: 'SolutionMixedcloud',
      component: () => import('@/views/solutions/mixedcloud.vue')
    },
    {
      path: '/solutions/highcalculation',
      name: 'SolutionHighcalculation',
      component: () => import('@/views/solutions/highcalculation.vue')
    },
    {
      path: '/solutions/high-availability',
      name: 'SolutionHighAvailability',
      component: () => import('@/views/solutions/HighAvailability.vue')
    },
    {
      path: '/solutions/website',
      name: 'SolutionWebsite',
      component: () => import('@/views/solutions/website.vue')
    },
    {
      path: '/solutions/hosting',
      name: 'SolutionHosting',
      component: () => import('@/views/solutions/hosting.vue')
    },
    {
      path: '/solutions/highbuild',
      name: 'SolutionHighbuild',
      component: () => import('@/views/solutions/HighBuild.vue')
    },
    {
      path: '/aup',
      name: 'Aup',
      component: () => import('@/views/Aup.vue')
    },
    {
      path: '/management',
      name: 'Management',
      component: () => import('@/views/Management.vue')
    },
    {
      path: '/relation',
      name: 'Relation',
      component: () => import('@/views/Relation.vue')
    },
    {
      path: '/safeguard',
      name: 'Safeguard',
      component: () => import('@/views/Safeguard.vue')
    },
    {
      path: '/transfer',
      name: 'Transfer',
      component: () => import('@/views/Transfer.vue')
    },
    {
      path: '/dediserver',
      name: 'Dediserver',
      component: () => import('@/views/Dediserver.vue')
    },
    {
      path: '/colocation',
      name: 'Colocation',
      component: () => import('@/views/Colocation.vue')
    },
    {
      path: '/privacy',
      name: 'Privacy',
      component: () => import('@/views/Privacy.vue')
    },
    {
      path: '/terms',
      name: 'Terms',
      component: () => import('@/views/Terms.vue')
    },
    {
      path: '/site-map',
      name: 'SiteMap',
      component: () => import('@/views/SiteMap.vue')
    },
    {
      path: '/user',
      component: () => import('@/layouts/UserLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/user/dashboard'
        },
        {
          path: 'dashboard',
          name: 'UserDashboard',
          component: () => import('@/views/user/Dashboard.vue')
        },
        {
          path: 'products',
          name: 'UserProducts',
          component: () => import('@/views/user/Products.vue')
        },
        {
          path: 'products/:id',
          name: 'UserServiceDetail',
          component: () => import('@/views/user/ServiceDetail.vue')
        },
        {
          path: 'orders',
          name: 'UserOrders',
          component: () => import('@/views/user/Orders.vue')
        },
        {
          path: 'orders/:id',
          name: 'UserOrderDetail',
          component: () => import('@/views/user/OrderDetail.vue')
        },
        {
          path: 'invoices',
          name: 'UserInvoices',
          component: () => import('@/views/user/Invoices.vue')
        },
        {
          path: 'invoice/list',
          name: 'InvoiceList',
          component: () => import('@/views/user/InvoiceList.vue')
        },
        {
          path: 'invoice/apply',
          name: 'InvoiceApply',
          component: () => import('@/views/user/InvoiceApply.vue')
        },
        {
          path: 'invoice/address',
          name: 'InvoiceAddress',
          component: () => import('@/views/user/InvoiceAddress.vue')
        },
        {
          path: 'invoice/company',
          name: 'InvoiceCompany',
          component: () => import('@/views/user/InvoiceCompany.vue')
        },
        {
          path: 'invoices/:id',
          name: 'InvoiceDetail',
          component: () => import('@/views/user/InvoiceDetail.vue')
        },
        {
          path: 'wallet',
          name: 'UserWallet',
          component: () => import('@/views/Wallet.vue')
        },
        {
          path: 'tickets',
          name: 'UserTickets',
          component: () => import('@/views/user/Tickets.vue')
        },
        {
          path: 'tickets/create',
          name: 'UserTicketCreate',
          component: () => import('@/views/user/TicketCreate.vue')
        },
        {
          path: 'tickets/:id',
          name: 'UserTicketDetail',
          component: () => import('@/views/user/TicketDetail.vue')
        },
        {
          path: 'profile',
          name: 'UserProfile',
          component: () => import('@/views/user/Profile.vue')
        },
        {
          path: 'security',
          name: 'UserSecurity',
          component: () => import('@/views/user/Security.vue')
        },
        {
          path: 'verification',
          name: 'UserVerification',
          component: () => import('@/views/user/Verification.vue')
        },
        {
          path: 'contacts',
          name: 'UserContacts',
          component: () => import('@/views/user/Contacts.vue')
        },
        {
          path: 'referral',
          name: 'UserReferral',
          component: () => import('@/views/user/Referral.vue')
        },
        {
          path: 'affiliate/buy-record',
          name: 'AffBuyRecord',
          component: () => import('@/views/user/affiliate/AffBuyRecord.vue')
        },
        {
          path: 'affiliate/user-list',
          name: 'UserAffiList',
          component: () => import('@/views/user/affiliate/UserAffiList.vue')
        },
        {
          path: 'affiliate/withdraw',
          name: 'WithdrawRecord',
          component: () => import('@/views/user/affiliate/WithdrawRecord.vue')
        },
        {
          path: 'oauth-bind',
          name: 'UserOauthBind',
          component: () => import('@/views/user/OauthBind.vue')
        },
        {
          path: 'system-message',
          name: 'UserSystemMessage',
          component: () => import('@/views/user/SystemMessage.vue')
        },
        {
          path: 'record-log',
          name: 'UserRecordLog',
          component: () => import('@/views/user/RecordLog.vue')
        },
        {
          path: 'upgrade',
          name: 'UserUpgrade',
          component: () => import('@/views/user/Upgrade.vue')
        },
        {
          path: 'transaction',
          name: 'UserTransaction',
          component: () => import('@/views/user/transaction/Transaction.vue')
        },
        {
          path: 'transaction/recharge',
          name: 'RechargeRecord',
          component: () => import('@/views/user/transaction/RechargeRecord.vue')
        },
        {
          path: 'transaction/refund',
          name: 'RefundRecord',
          component: () => import('@/views/user/transaction/RefundRecord.vue')
        },
        {
          path: 'transaction/withdraw',
          name: 'TransactionWithdrawRecord',
          component: () => import('@/views/user/transaction/WithdrawRecord.vue')
        },
        {
          path: 'transaction/credit',
          name: 'CreditRecord',
          component: () => import('@/views/user/transaction/CreditRecord.vue')
        },
        {
          path: 'transaction/accounts',
          name: 'AccountsRecord',
          component: () => import('@/views/user/transaction/AccountsRecord.vue')
        },
        {
          path: 'api-manage',
          name: 'UserApiManage',
          component: () => import('@/views/user/ApiManage.vue')
        },
        {
          path: 'api-log',
          name: 'UserApiLog',
          component: () => import('@/views/user/ApiLog.vue')
        },
        {
          path: 'login-log',
          name: 'UserLoginLog',
          component: () => import('@/views/user/LoginLog.vue')
        },
        {
          path: 'system-log',
          name: 'UserSystemLog',
          component: () => import('@/views/user/SystemLog.vue')
        },
        {
          path: 'add-funds',
          name: 'UserAddFunds',
          component: () => import('@/views/user/AddFunds.vue')
        },
        {
          path: 'bind-account',
          name: 'UserBindAccount',
          component: () => import('@/views/user/BindAccount.vue')
        },
        {
          path: 'batch-renew',
          name: 'UserBatchRenew',
          component: () => import('@/views/user/BatchRenew.vue')
        },
        {
          path: 'maintenance',
          name: 'UserMaintenance',
          component: () => import('@/views/user/Maintenance.vue')
        },
        {
          path: 'credit-limit',
          name: 'UserCreditLimit',
          component: () => import('@/views/user/CreditLimit.vue')
        },
        {
          path: 'host',
          name: 'UserHost',
          component: () => import('@/views/user/Host.vue')
        },
        {
          path: 'contract',
          name: 'UserContract',
          component: () => import('@/views/user/Contract.vue')
        },
        {
          path: 'contract-host',
          name: 'UserContractHost',
          component: () => import('@/views/user/ContractHost.vue')
        },
        {
          path: 'combine-billing',
          name: 'UserCombineBilling',
          component: () => import('@/views/user/CombineBilling.vue')
        },
        {
          path: 'login-access-token',
          name: 'UserLoginAccessToken',
          component: () => import('@/views/user/LoginAccessToken.vue')
        },
        {
          path: 'service/sms',
          name: 'UserServiceSMS',
          component: () => import('@/views/user/ServiceSMS.vue')
        },
        {
          path: 'service/soft',
          name: 'UserServiceSoft',
          component: () => import('@/views/user/ServiceSoft.vue')
        },
        {
          path: 'enterprise-verification',
          name: 'UserEnterpriseVerification',
          component: () => import('@/views/user/EnterpriseVerification.vue')
        },
        {
          path: 'marketplace',
          name: 'UserMarketplace',
          component: () => import('@/views/user/Marketplace.vue')
        },
        {
          path: 'marketplace/sell',
          name: 'UserMarketplaceSell',
          component: () => import('@/views/user/MarketplaceSell.vue')
        },
        {
          path: 'marketplace/orders',
          name: 'UserMarketplaceOrders',
          component: () => import('@/views/user/MarketplaceOrders.vue')
        },
        {
          path: 'marketplace/chat',
          name: 'UserMarketplaceChat',
          component: () => import('@/views/user/MarketplaceChat.vue')
        },
        {
          path: 'marketplace/chat/:listing_id/:user_id',
          name: 'UserMarketplaceChatDirect',
          component: () => import('@/views/user/MarketplaceChat.vue')
        },
        {
          path: 'marketplace/earnings',
          name: 'UserMarketplaceEarnings',
          component: () => import('@/views/user/MarketplaceEarnings.vue')
        },
        {
          path: 'marketplace/transactions',
          name: 'UserMarketplaceTransactions',
          component: () => import('@/views/user/MarketplaceTransactions.vue')
        },
        {
          path: 'marketplace/logs',
          name: 'UserMarketplaceLogs',
          component: () => import('@/views/user/MarketplaceLogs.vue')
        },
        {
          path: 'credit-bill/:id',
          name: 'UserCreditBillDetail',
          component: () => import('@/views/user/CreditBillDetail.vue')
        },
        {
          path: 'credit-used-detail',
          name: 'UserCreditUsedDetail',
          component: () => import('@/views/user/CreditUsedDetail.vue')
        },
        {
          path: 'other-server',
          name: 'UserOtherServer',
          component: () => import('@/views/user/OtherServer.vue')
        },
        {
          path: 'dcim-console/:id',
          name: 'DcimConsole',
          component: () => import('@/views/user/DcimConsole.vue')
        },
        {
          path: 'vnc-console/:id',
          name: 'VncConsole',
          component: () => import('@/views/user/VncConsole.vue')
        },
        {
          path: 'product-transfer',
          name: 'ProductTransfer',
          component: () => import('@/views/user/ProductTransfer.vue')
        }
      ]
    },
    {
      path: '/announcements',
      name: 'Announcements',
      component: () => import('@/views/Announcements.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    const userStore = useUserStore()
    if (!userStore.token) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
