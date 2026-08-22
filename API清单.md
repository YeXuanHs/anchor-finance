# AnchorFinance API清单

**自动生成时间**: 2026-08-22
**Admin路由数**: 421
**Client路由数**: 93
**总计**: 514

---

## 一、Admin管理后台API

### AI工单

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ai-ticket/config | GetAITicketConfig |
| GET | /api/admin/ai-ticket/queue/stats | GetAITicketQueueStats |
| GET | /api/admin/ai-ticket/queue | GetAITicketQueueList |
| POST | /api/admin/ai-ticket/queue/process | ProcessAITicketQueue |
| GET | /api/admin/ai-ticket/knowledge | GetAITicketKnowledgeList |
| POST | /api/admin/ai-ticket/knowledge | CreateAITicketKnowledge |
| PUT | /api/admin/ai-ticket/knowledge/:id | UpdateAITicketKnowledge |
| DELETE | /api/admin/ai-ticket/knowledge/:id | DeleteAITicketKnowledge |
| GET | /api/admin/ai-ticket/rules | GetAITicketRuleList |
| POST | /api/admin/ai-ticket/rules | CreateAITicketRule |
| PUT | /api/admin/ai-ticket/rules/:id | UpdateAITicketRule |
| DELETE | /api/admin/ai-ticket/rules/:id | DeleteAITicketRule |
| GET | /api/admin/ai-ticket/logs | GetAITicketProcessLogs |
| POST | /api/admin/ai-ticket/tickets/:id/mode | SetAITicketMode |

### AI配置

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ai/config | GetAIConfig |
| POST | /api/admin/ai/test | TestAIConnection |
| POST | /api/admin/ai/generate-description | GenerateProductDescription |
| POST | /api/admin/ai/ticket-reply | AITicketReply |

### Contract Templates

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/contract-templates | GetContractTemplateList |
| POST | /api/admin/contract-templates | CreateContractTemplate |
| PUT | /api/admin/contract-templates/:id | UpdateContractTemplate |
| DELETE | /api/admin/contract-templates/:id | DeleteContractTemplate |

### Coupon Campaigns

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/coupon-campaigns | GetCouponCampaignList |
| GET | /api/admin/coupon-campaigns/summary | GetCouponCampaignSummary |
| POST | /api/admin/coupon-campaigns | CreateCouponCampaign |
| PUT | /api/admin/coupon-campaigns/:id | UpdateCouponCampaign |
| DELETE | /api/admin/coupon-campaigns/:id | DeleteCouponCampaign |
| PATCH | /api/admin/coupon-campaigns/:id/status | UpdateCouponCampaignStatus |

### Coupon Product Groups

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/coupon-product-groups | GetCouponProductGroups |

### Cpu Model Catalog

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/cpu-model-catalog | GetCPUModelCatalog |

### Cron Tasks

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/cron-tasks | GetCronTasks |

### Custom Fields

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/custom-fields | GetCustomFieldList |
| POST | /api/admin/custom-fields | CreateCustomField |
| PUT | /api/admin/custom-fields/:id | UpdateCustomField |
| DELETE | /api/admin/custom-fields/:id | DeleteCustomField |

### Finance

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/finance/new-customer-daily-summary | GetNewCustomerDailySummary |
| GET | /api/admin/finance/product-income-summary | GetProductIncomeSummary |
| GET | /api/admin/finance/ledger | GetFinanceLedger |
| GET | /api/admin/finance/ledger/:id | GetFinanceLedgerDetail |
| GET | /api/admin/finance/ledger/summary | GetFinanceLedgerSummary |
| GET | /api/admin/finance/recharges | GetRechargeList |
| GET | /api/admin/finance/recharges/summary | GetRechargeSummary |
| GET | /api/admin/finance/renewal-orders | GetRenewalOrders |
| GET | /api/admin/finance/upgrade-orders | GetUpgradeOrders |

### Finance Config

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/finance-config | GetFinanceConfig |
| PUT | /api/admin/finance-config | UpdateFinanceConfig |

### Instance Spec Catalog

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/instance-spec-catalog | GetInstanceSpecCatalog |

### Log Cleanups

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/log-cleanups/overview | GetLogCleanupOverview |
| POST | /api/admin/log-cleanups | CleanupLogs |

### Log Summaries

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/log-summaries/:channel | GetLogSummaryByChannel |

### Media Files

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/media-files | GetMediaFileList |
| DELETE | /api/admin/media-files/:id | DeleteMediaFile |
| GET | /api/admin/media-files/:id/references | GetMediaFileReferences |

### Member Levels

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/member-levels | GetMemberLevelList |
| POST | /api/admin/member-levels | CreateMemberLevel |
| PUT | /api/admin/member-levels/:id | UpdateMemberLevel |
| DELETE | /api/admin/member-levels/:id | DeleteMemberLevel |

### Menu Types

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/menu-types | GetMenuTypeList |

### News Categories

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/news-categories | GetNewsCategories |
| POST | /api/admin/news-categories | CreateNewsCategory |
| PUT | /api/admin/news-categories/:id | UpdateNewsCategory |
| DELETE | /api/admin/news-categories/:id | DeleteNewsCategory |

### Notification Templates

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/notification-templates | GetNotificationTemplates |
| POST | /api/admin/notification-templates | CreateNotificationTemplate |
| POST | /api/admin/notification-templates/test-send | TestNotificationTemplate |
| PUT | /api/admin/notification-templates/:id | UpdateNotificationTemplate |
| DELETE | /api/admin/notification-templates/:id | DeleteNotificationTemplate |

### Order Config

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/order-config | GetOrderConfig |
| PUT | /api/admin/order-config | UpdateOrderConfig |

### Os Options

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/os-options | GetOSOptions |

### Permissions

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/permissions | GetPermissions |

### Product Groups

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/product-groups | GetProductGroups |
| GET | /api/admin/product-groups/tree | GetProductGroupTree |
| GET | /api/admin/product-groups/:id | GetProductGroupDetail |
| GET | /api/admin/product-groups/:id/children | GetProductGroupChildren |
| POST | /api/admin/product-groups | CreateProductGroup |
| PUT | /api/admin/product-groups/:id | UpdateProductGroup |
| DELETE | /api/admin/product-groups/:id | DeleteProductGroup |
| POST | /api/admin/product-groups/reorders | ReorderProductGroups |

### Product Types

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/product-types | GetProductTypeList |
| POST | /api/admin/product-types | CreateProductType |
| PUT | /api/admin/product-types/:id | UpdateProductType |
| DELETE | /api/admin/product-types/:id | DeleteProductType |
| POST | /api/admin/product-types/reorders | ReorderProductTypes |

### Redis配置

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/redis/config | GetRedisConfig |
| GET | /api/admin/redis/health | RedisHealthCheck |

### Reports

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/reports/new-customers | GetNewCustomerStatistics |
| GET | /api/admin/reports/revenue-ranking | GetRevenueRanking |

### Sales Config

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/sales-config | GetSalesConfig |
| PUT | /api/admin/sales-config | UpdateSalesConfig |

### Sales Groups

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/sales-groups | GetSalesGroupList |
| POST | /api/admin/sales-groups | CreateSalesGroup |
| PUT | /api/admin/sales-groups/:id | UpdateSalesGroup |
| DELETE | /api/admin/sales-groups/:id | DeleteSalesGroup |

### Schedule Triggers

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /api/admin/schedule-triggers | TriggerSchedule |

### Send Message

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/send-message/search-params | GetSendMessageSearchParams |
| GET | /api/admin/send-message/send-methods | GetSendMethodList |
| GET | /api/admin/send-message/search | SearchSendMessageList |

### Ticket Delivery Rules

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ticket-delivery-rules | GetTicketDeliveryRules |
| POST | /api/admin/ticket-delivery-rules | CreateTicketDeliveryRule |
| PUT | /api/admin/ticket-delivery-rules/:id | UpdateTicketDeliveryRule |
| DELETE | /api/admin/ticket-delivery-rules/:id | DeleteTicketDeliveryRule |

### Ticket Departments

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ticket-departments | GetTicketDepartments |
| GET | /api/admin/ticket-departments/:id | GetTicketDepartmentDetail |
| POST | /api/admin/ticket-departments/:id/move-up | MoveTicketDepartmentUp |
| POST | /api/admin/ticket-departments/:id/move-down | MoveTicketDepartmentDown |

### Ticket Prereplies

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ticket-prereplies | GetTicketPrereplyList |
| POST | /api/admin/ticket-prereplies | CreateTicketPrereply |
| PUT | /api/admin/ticket-prereplies/:id | UpdateTicketPrereply |
| DELETE | /api/admin/ticket-prereplies/:id | DeleteTicketPrereply |
| POST | /api/admin/ticket-prereplies/search | SearchTicketPrereply |

### Ticket Prereply Categories

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ticket-prereply-categories | GetTicketPrereplyCategoryList |
| POST | /api/admin/ticket-prereply-categories | CreateTicketPrereplyCategory |
| PUT | /api/admin/ticket-prereply-categories/:id | UpdateTicketPrereplyCategory |
| DELETE | /api/admin/ticket-prereply-categories/:id | DeleteTicketPrereplyCategory |

### Ticket Statuses

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ticket-statuses | GetTicketStatuses |
| GET | /api/admin/ticket-statuses/:id | GetTicketStatusDetail |

### Traffic Logs

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/traffic-logs | GetTrafficLogList |

### Transactions

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/transactions | GetTransactionList |

### Upload

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /api/admin/upload | UploadFile |

### 下载管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/downloads | GetDownloads |
| POST | /api/admin/downloads | CreateDownload |
| PUT | /api/admin/downloads/:id | UpdateDownload |
| DELETE | /api/admin/downloads/:id | DeleteDownload |
| GET | /api/admin/downloads/categories | GetDownloadCategories |
| POST | /api/admin/downloads/categories | CreateDownloadCategory |
| PUT | /api/admin/downloads/categories/:id | UpdateDownloadCategory |
| DELETE | /api/admin/downloads/categories/:id | DeleteDownloadCategory |

### 主题模板

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/home-hero | GetHomeHero |
| PUT | /api/admin/home-hero | UpdateHomeHero |
| GET | /api/admin/home-hero/assets | GetHomeHeroAssets |

### 主题管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/themes | GetThemeList |
| GET | /api/admin/themes/active | GetActiveTheme |
| POST | /api/admin/themes | CreateTheme |
| PUT | /api/admin/themes/:id | UpdateTheme |
| DELETE | /api/admin/themes/:id | DeleteTheme |
| POST | /api/admin/themes/:id/set-default | SetDefaultTheme |

### 二次验证

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/two-factor-config | GetTwoFactorConfig |
| PUT | /api/admin/two-factor-config | UpdateTwoFactorConfig |

### 产品管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/products | GetProductList |
| GET | /api/admin/products/summary | GetProductSummary |
| GET | /api/admin/products/:id | GetProduct |
| GET | /api/admin/products/:id/owners | GetProductOwners |
| POST | /api/admin/products | CreateProduct |
| PUT | /api/admin/products/:id | UpdateProduct |
| DELETE | /api/admin/products/:id | DeleteProduct |
| POST | /api/admin/products/:id/restorations | RestoreProduct |
| PATCH | /api/admin/products/:id/status | UpdateProductStatus |
| POST | /api/admin/products/reorders | ReorderProducts |
| POST | /api/admin/products/category-batches | BatchUpdateProductCategory |
| DELETE | /api/admin/products/:id/force | ForceDeleteProduct |

### 仪表盘

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/dashboard/stats | GetDashboardStats |
| GET | /api/admin/dashboard/income-trend | GetIncomeTrend |
| GET | /api/admin/dashboard/online-admins | GetOnlineAdmins |
| GET | /api/admin/dashboard/recent-invoices | GetRecentInvoices |
| GET | /api/admin/dashboard/monthly-revenue | GetMonthlyRevenue |

### 任务队列

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/task-queue/overview | GetTaskQueueOverview |
| GET | /api/admin/task-queue | GetTaskQueueList |
| POST | /api/admin/task-queue/:id/retry | RetryTask |
| DELETE | /api/admin/task-queue/:id | DeleteTask |

### 优惠券

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/coupons | GetCouponList |
| GET | /api/admin/coupons/summary | GetCouponSummary |
| POST | /api/admin/coupons | CreateCoupon |
| PUT | /api/admin/coupons/:id | UpdateCoupon |
| DELETE | /api/admin/coupons/:id | DeleteCoupon |
| PATCH | /api/admin/coupons/:id/status | UpdateCouponStatus |

### 优惠码

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/promo-codes | GetPromoCodeList |
| POST | /api/admin/promo-codes | CreatePromoCode |
| PUT | /api/admin/promo-codes/:id | UpdatePromoCode |
| DELETE | /api/admin/promo-codes/:id | DeletePromoCode |

### 供应商管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/suppliers | GetSupplierList |
| GET | /api/admin/suppliers/summary | GetSupplierSummary |
| GET | /api/admin/suppliers/provider-types | GetSupplierProviderTypes |
| GET | /api/admin/suppliers/:id | GetSupplierDetail |
| GET | /api/admin/suppliers/:id/balance | GetSupplierBalance |
| PATCH | /api/admin/suppliers/:id/status | UpdateSupplierStatus |
| POST | /api/admin/suppliers | CreateSupplier |
| PUT | /api/admin/suppliers/:id | UpdateSupplier |
| DELETE | /api/admin/suppliers/:id | DeleteSupplier |
| GET | /api/admin/suppliers/:id/products | GetSupplierProducts |
| POST | /api/admin/suppliers/:id/tasks | RunSupplierTask |
| POST | /api/admin/suppliers/:id/sync-products | SyncSupplierProducts |
| POST | /api/admin/suppliers/:id/sync-prices | SyncSupplierPrices |
| POST | /api/admin/suppliers/:id/sync-stock | SyncSupplierStock |

### 信用额管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/credit-limits | GetCreditLimitList |
| GET | /api/admin/credit-limits/config | GetCreditLimitConfig |
| POST | /api/admin/credit-limits/config | SaveCreditLimitConfig |
| POST | /api/admin/credit-limits | SaveCreditLimit |
| PUT | /api/admin/credit-limits/:id | UpdateCreditLimit |
| DELETE | /api/admin/credit-limits/:id | DeleteCreditLimit |
| GET | /api/admin/credit-limits/logs | GetCreditLimitLogs |

### 内容概览

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/content/summary | GetContentSummary |

### 分类日志

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/logs/sms | GetSMSLogs |
| GET | /api/admin/logs/email | GetEmailLogs |
| GET | /api/admin/logs/api | GetAPILogs |
| GET | /api/admin/logs/cron | GetCronLogs |
| GET | /api/admin/logs/admin-login | GetAdminLoginLogs |
| GET | /api/admin/logs/notification | GetNotificationLogs |

### 友情链接

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/friendly-links | GetFriendlyLinkList |
| POST | /api/admin/friendly-links | CreateFriendlyLink |
| PUT | /api/admin/friendly-links/:id | UpdateFriendlyLink |
| DELETE | /api/admin/friendly-links/:id | DeleteFriendlyLink |

### 取消请求

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/cancel-requests | GetCancelRequestList |
| POST | /api/admin/cancel-requests/:id/approve | ApproveCancelRequest |
| POST | /api/admin/cancel-requests/:id/reject | RejectCancelRequest |

### 可配置项

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/configurable-options | GetConfigurableOptionList |
| POST | /api/admin/configurable-options | CreateConfigurableOption |
| PUT | /api/admin/configurable-options/:id | UpdateConfigurableOption |
| DELETE | /api/admin/configurable-options/:id | DeleteConfigurableOption |

### 合同管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/contracts | GetContractList |
| GET | /api/admin/contracts/:id | GetContractDetail |
| POST | /api/admin/contracts | CreateContract |
| PUT | /api/admin/contracts/:id | UpdateContract |
| DELETE | /api/admin/contracts/:id | DeleteContract |
| POST | /api/admin/contracts/:id/sign | SignContract |
| POST | /api/admin/contracts/:id/cancel | CancelContract |

### 员工管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/staff | GetStaffList |
| GET | /api/admin/staff/roles | GetStaffRoles |
| GET | /api/admin/staff/:id | GetStaffDetail |
| POST | /api/admin/staff | CreateStaff |
| PUT | /api/admin/staff/:id | UpdateStaff |
| DELETE | /api/admin/staff/:id | DeleteStaff |
| PATCH | /api/admin/staff/:id/status | UpdateStaffStatus |
| POST | /api/admin/staff/:id/password-resets | ResetStaffPassword |

### 定时任务

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/schedules/overview | GetScheduleOverview |
| GET | /api/admin/schedule-runs | GetScheduleRunList |
| GET | /api/admin/schedule-runs/:id | GetScheduleRunDetail |
| POST | /api/admin/schedule-runs/:id/retry | RetryScheduleRun |

### 实名认证

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/verifications | GetVerificationList |
| GET | /api/admin/verifications/summary | GetVerificationSummary |
| GET | /api/admin/verifications/:id | GetVerificationDetail |
| GET | /api/admin/verifications/:id/history | GetVerificationHistory |
| POST | /api/admin/verifications/:id/approve | ApproveVerification |
| POST | /api/admin/verifications/:id/reject | RejectVerification |
| POST | /api/admin/verifications/:id/unbindings | UnbindVerificationByUser |

### 客户分组

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/customer-groups | GetCustomerGroupList |
| POST | /api/admin/customer-groups | CreateCustomerGroup |
| PUT | /api/admin/customer-groups/:id | UpdateCustomerGroup |
| DELETE | /api/admin/customer-groups/:id | DeleteCustomerGroup |

### 工单管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/tickets | GetTicketList |
| GET | /api/admin/tickets/summary | GetTicketSummary |
| GET | /api/admin/tickets/admin-users | GetTicketAdminUsers |
| GET | /api/admin/tickets/:id | GetTicket |
| GET | /api/admin/tickets/:id/replies | GetTicketReplies |
| POST | /api/admin/tickets/:id/reply | ReplyTicket |
| POST | /api/admin/tickets/:id/close | CloseTicket |
| POST | /api/admin/tickets/:id/reopen | ReopenTicket |
| POST | /api/admin/tickets/:id/receive | ReceiveTicket |
| PUT | /api/admin/tickets/:id/assignment | AssignTicket |
| POST | /api/admin/tickets/:id/replies/:reply_id/recalls | RecallTicketReply |
| GET | /api/admin/tickets/:id/upstream-delivery | GetTicketUpstreamDelivery |
| GET | /api/admin/tickets/:id/upstream-delivery/logs | GetTicketUpstreamDeliveryLogs |

### 工单规则

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/ticket-rules | GetTicketRuleList |
| POST | /api/admin/ticket-rules | CreateTicketRule |
| PUT | /api/admin/ticket-rules/:id | UpdateTicketRule |
| DELETE | /api/admin/ticket-rules/:id | DeleteTicketRule |

### 推介系统

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/referral/overview | GetReferralOverview |
| GET | /api/admin/referral/rewards | GetReferralRewards |
| GET | /api/admin/referral-withdrawals | GetReferralWithdrawals |
| POST | /api/admin/referral-withdrawals/:id/approve | ApproveReferralWithdrawal |
| POST | /api/admin/referral-withdrawals/:id/reject | RejectReferralWithdrawal |

### 插件域

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/payment-gateways | GetPaymentGateways |
| GET | /api/admin/sms-providers | GetSMSProviders |
| GET | /api/admin/mail-providers | GetMailProviders |
| GET | /api/admin/certification-providers | GetCertificationProviders |
| GET | /api/admin/server-modules | GetServerModules |

### 插件管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/plugins | GetPluginList |
| POST | /api/admin/plugins/install | InstallPlugin |
| POST | /api/admin/plugins/scan | ScanPlugins |
| GET | /api/admin/plugins/:id | GetPluginDetail |
| POST | /api/admin/plugins/:id/enable | EnablePlugin |
| POST | /api/admin/plugins/:id/disable | DisablePlugin |
| DELETE | /api/admin/plugins/:id | UninstallPlugin |
| GET | /api/admin/plugins/:id/config | GetPluginConfig |
| PUT | /api/admin/plugins/:id/config | UpdatePluginConfig |
| POST | /api/admin/plugins/:id/health | PluginHealthCheck |

### 数据库管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/database/status | GetDatabaseStatus |
| POST | /api/admin/database/optimizations | OptimizeDatabase |
| POST | /api/admin/database/backups | BackupDatabase |

### 新闻管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/news | GetNewsList |
| GET | /api/admin/news/:id | GetNewsDetail |
| POST | /api/admin/news | CreateNews |
| PUT | /api/admin/news/:id | UpdateNews |
| DELETE | /api/admin/news/:id | DeleteNews |

### 日志管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/system-logs | GetSystemLogs |
| GET | /api/admin/operation-logs | GetOperationLogs |
| GET | /api/admin/login-logs | GetLoginLogs |

### 服务管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/services | GetServiceList |
| GET | /api/admin/services/:id | GetService |
| PUT | /api/admin/services/:id | UpdateService |
| POST | /api/admin/services/:id/suspend | SuspendService |
| POST | /api/admin/services/:id/unsuspend | UnsuspendService |
| POST | /api/admin/services/:id/terminate | TerminateService |

### 流量包

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/traffic-packages | GetTrafficPackageList |
| POST | /api/admin/traffic-packages | CreateTrafficPackage |
| PUT | /api/admin/traffic-packages/:id | UpdateTrafficPackage |
| DELETE | /api/admin/traffic-packages/:id | DeleteTrafficPackage |

### 用户管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/users | GetUserList |
| GET | /api/admin/users/:id | GetUser |
| POST | /api/admin/users | CreateUser |
| PUT | /api/admin/users/:id | UpdateUser |
| DELETE | /api/admin/users/:id | DeleteUser |
| PATCH | /api/admin/users/:id/status | UpdateUserStatus |
| GET | /api/admin/users/:id/orders | GetUserOrders |
| GET | /api/admin/users/:id/invoices | GetUserInvoices |
| GET | /api/admin/users/:id/tickets | GetUserTickets |
| GET | /api/admin/users/:id/services | GetUserServices |
| GET | /api/admin/users/:id/balance-logs | GetUserBalanceLogs |
| GET | /api/admin/users/:id/operation-logs | GetUserOperationLogs |
| GET | /api/admin/users/:id/email-logs | GetUserEmailLogs |
| GET | /api/admin/users/:id/sms-logs | GetUserSmsLogs |
| GET | /api/admin/users/:id/invoices/:invoice_id | GetUserInvoiceDetail |
| POST | /api/admin/users/:id/recharges | RechargeUser |
| POST | /api/admin/users/:id/services/:service_id/refunds | RefundUserService |
| GET | /api/admin/users/:id/services/:service_id/connection | AdminGetServiceConnection |
| GET | /api/admin/users/:id/services/:service_id/remote-status | AdminGetServiceRemoteStatus |
| PUT | /api/admin/users/:id/services/:service_id/meta | AdminUpdateServiceMeta |
| POST | /api/admin/users/:id/services/:service_id/manual-provision | AdminManualProvision |
| POST | /api/admin/users/:id/services/:service_id/power-actions | AdminServicePowerAction |
| POST | /api/admin/users/:id/services/:service_id/password-resets | AdminResetServicePassword |
| GET | /api/admin/users/:id/remarks | GetUserRemarks |
| POST | /api/admin/users/:id/remarks | AddUserRemark |
| POST | /api/admin/users/:id/login-as | LoginAsUser |
| POST | /api/admin/users/:id/services/refresh-statuses | RefreshUserServicesStatus |
| POST | /api/admin/users/:id/unbind-verification | UnbindVerification |

### 知识库

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/knowledge/categories | GetKnowledgeCategories |
| POST | /api/admin/knowledge/categories | CreateKnowledgeCategory |
| PUT | /api/admin/knowledge/categories/:id | UpdateKnowledgeCategory |
| DELETE | /api/admin/knowledge/categories/:id | DeleteKnowledgeCategory |
| GET | /api/admin/knowledge/articles | GetKnowledgeArticles |
| GET | /api/admin/knowledge/articles/:id | GetKnowledgeArticleDetail |
| POST | /api/admin/knowledge/articles | CreateKnowledgeArticle |
| PUT | /api/admin/knowledge/articles/:id | UpdateKnowledgeArticle |
| DELETE | /api/admin/knowledge/articles/:id | DeleteKnowledgeArticle |

### 短信模板

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/sms-templates | GetSMSTemplateList |
| POST | /api/admin/sms-templates | CreateSMSTemplate |
| PUT | /api/admin/sms-templates/:id | UpdateSMSTemplate |
| DELETE | /api/admin/sms-templates/:id | DeleteSMSTemplate |

### 第三方登录

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/oauth-providers | GetOAuthProviderList |
| POST | /api/admin/oauth-providers | CreateOAuthProvider |
| PUT | /api/admin/oauth-providers/:id | UpdateOAuthProvider |
| DELETE | /api/admin/oauth-providers/:id | DeleteOAuthProvider |

### 管理员

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/admins | GetAdminList |
| POST | /api/admin/admins | CreateAdmin |
| PUT | /api/admin/admins/:id | UpdateAdmin |

### 系统信息

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/system/info | GetSystemInfo |
| GET | /api/admin/system/modules | GetSystemModules |

### 自定义字段

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/custom-template-fields | GetCustomTemplateFieldList |
| POST | /api/admin/custom-template-fields | CreateCustomTemplateField |
| PUT | /api/admin/custom-template-fields/:id | UpdateCustomTemplateField |
| DELETE | /api/admin/custom-template-fields/:id | DeleteCustomTemplateField |

### 菜单管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/menus | GetMenus |
| POST | /api/admin/menus | CreateMenu |
| PUT | /api/admin/menus/:id | UpdateMenu |
| DELETE | /api/admin/menus/:id | DeleteMenu |

### 营销推送

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/marketing/pushes | GetMarketingPushList |
| POST | /api/admin/marketing/pushes | CreateMarketingPush |
| POST | /api/admin/marketing/pushes/:id/send | SendMarketingPush |
| DELETE | /api/admin/marketing/pushes/:id | DeleteMarketingPush |

### 角色管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/roles | GetRoleList |
| GET | /api/admin/roles/:id | GetRoleDetail |
| POST | /api/admin/roles | CreateRole |
| PUT | /api/admin/roles/:id | UpdateRole |
| DELETE | /api/admin/roles/:id | DeleteRole |
| POST | /api/admin/roles/:id/copies | CopyRole |

### 订单管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/orders | GetOrderList |
| POST | /api/admin/orders/search | SearchOrders |
| GET | /api/admin/orders/:id | GetOrder |
| POST | /api/admin/orders | CreateOrder |
| PUT | /api/admin/orders/:id | UpdateOrder |
| POST | /api/admin/orders/:id/activate | ActivateOrder |
| POST | /api/admin/orders/:id/cancel | CancelOrder |
| POST | /api/admin/orders/:id/notes | AddOrderNote |

### 认证

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /api/admin/login | authHandler |
| GET | /api/admin/auth/info | authHandler |
| PUT | /api/admin/auth/profile | authHandler |
| PUT | /api/admin/auth/password | authHandler |
| POST | /api/admin/logout | authHandler |
| POST | /api/admin/auth/reset-password | authHandler |

### 设置管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/settings | GetSettings |
| GET | /api/admin/settings/:group | GetSettingsByGroup |
| PUT | /api/admin/settings | UpdateSettings |
| GET | /api/admin/settings/email | GetEmailConfig |
| PUT | /api/admin/settings/email | UpdateEmailConfig |
| GET | /api/admin/settings/sms | GetSMSConfig |
| PUT | /api/admin/settings/sms | UpdateSMSConfig |
| GET | /api/admin/settings/register-login | GetRegisterLoginConfig |
| PUT | /api/admin/settings/register-login | UpdateRegisterLoginConfig |
| GET | /api/admin/settings/captcha | GetCaptchaConfig |
| PUT | /api/admin/settings/captcha | UpdateCaptchaConfig |
| GET | /api/admin/settings/security | GetSecurityConfig |
| PUT | /api/admin/settings/security | UpdateSecurityConfig |
| GET | /api/admin/settings/general | GetGeneralConfig |
| PUT | /api/admin/settings/general | UpdateGeneralConfig |
| GET | /api/admin/settings/display | GetDisplayConfig |
| PUT | /api/admin/settings/display | UpdateDisplayConfig |
| GET | /api/admin/settings/invoice | GetInvoiceConfig |
| PUT | /api/admin/settings/invoice | UpdateInvoiceConfig |
| GET | /api/admin/settings/contract | GetContractConfig |
| PUT | /api/admin/settings/contract | UpdateContractConfig |
| GET | /api/admin/settings/credit-setting | GetCreditSettingConfig |
| PUT | /api/admin/settings/credit-setting | UpdateCreditSettingConfig |
| GET | /api/admin/settings/payment-gateway | GetPaymentGatewayConfig |
| PUT | /api/admin/settings/payment-gateway | UpdatePaymentGatewayConfig |

### 账单管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/invoices | GetInvoiceList |
| GET | /api/admin/invoices/:id | GetInvoice |
| POST | /api/admin/invoices/:id/cancel | CancelInvoice |
| POST | /api/admin/invoices/:id/notes | AddInvoiceNote |

### 货币管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/currencies | GetCurrencyList |
| POST | /api/admin/currencies | CreateCurrency |
| PUT | /api/admin/currencies/:id | UpdateCurrency |

### 邮件模板

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/email-templates | GetEmailTemplateList |
| GET | /api/admin/email-templates/:id | GetEmailTemplateDetail |
| POST | /api/admin/email-templates | CreateEmailTemplate |
| PUT | /api/admin/email-templates/:id | UpdateEmailTemplate |
| DELETE | /api/admin/email-templates/:id | DeleteEmailTemplate |

### 销售管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/sales/statistics | GetSalesStatistics |
| GET | /api/admin/sales/records | GetSalesRecords |

### 黑名单

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/admin/blacklist | GetBlacklist |
| POST | /api/admin/blacklist | CreateBlacklist |
| DELETE | /api/admin/blacklist/:id | DeleteBlacklist |

---

## 二、Client用户前台API

### Password

| 方法 | 路径 | Handler |
|------|------|---------|
| PUT | /api/client/password | authHandler |

### Payment

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /api/client/payment/notify/:gateway | PaymentNotify |

### Register

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /api/client/register | authHandler |

### 主题模板

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/home-hero | GetClientHomeHero |

### 产品管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/products | GetClientProducts |
| GET | /api/client/products/categories | GetClientProductCategories |
| GET | /api/client/products/:id | GetClientProductDetail |

### 优惠券

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/coupons | GetUserCoupons |
| POST | /api/client/coupons/:id/claim | ClaimCoupon |

### 余额日志

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/balance-logs | GetBalanceLogs |
| GET | /api/client/balance-logs/summary | GetBalanceLogsSummary |

### 充值

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/recharge/gateways | GetRechargeGateways |
| POST | /api/client/recharge | CreateRecharge |
| GET | /api/client/recharge/:paymentNo/status | GetRechargeStatus |

### 公告

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/notices | GetNotices |
| GET | /api/client/notices/:id | GetNoticeDetail |
| GET | /api/client/notices/unread-count | GetNoticesUnreadCount |
| POST | /api/client/notices/mark-all-read | MarkAllNoticesRead |

### 内容概览

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/content/overview | GetContentOverview |

### 实名认证

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/verification/status | GetVerificationStatus |
| POST | /api/client/verification/submit | SubmitVerification |

### 工单管理

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /api/client/tickets/upstream/replies | TicketUpstreamReply |
| GET | /api/client/tickets | GetUserTickets |
| GET | /api/client/tickets/service-options | GetTicketServiceOptions |
| GET | /api/client/tickets/:id | GetUserTicket |
| GET | /api/client/tickets/:id/replies | GetUserTicketReplies |
| POST | /api/client/tickets | CreateUserTicket |
| POST | /api/client/tickets/upload-images | UploadTicketImages |
| POST | /api/client/tickets/:id/reply | ReplyUserTicket |
| POST | /api/client/tickets/:id/close | CloseUserTicket |

### 帮助文章

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/help-articles | GetHelpArticles |
| GET | /api/client/help-articles/:id | GetHelpArticleDetail |

### 推介系统

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/referral/overview | GetUserReferralOverview |
| GET | /api/client/referral/rewards | GetUserReferralRewards |
| GET | /api/client/referral/direct-referrals | GetUserDirectReferrals |
| GET | /api/client/referral/account-logs | GetUserReferralAccountLogs |
| POST | /api/client/referral/withdrawals | ApplyReferralWithdrawal |
| GET | /api/client/referral/withdrawals | GetUserReferralWithdrawals |

### 支付记录

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/payments | GetPaymentList |
| GET | /api/client/payments/summary | GetPaymentSummary |
| GET | /api/client/payments/:id | GetPaymentDetail |

### 服务管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/services | GetUserServices |
| GET | /api/client/services/grouped-overview | GetUserServicesGroupedOverview |
| GET | /api/client/services/:id | GetUserService |
| GET | /api/client/services/:id/connection | GetServiceConnection |
| GET | /api/client/services/:id/runtime | GetServiceRuntime |
| PUT | /api/client/services/:id/name | UpdateServiceName |
| PUT | /api/client/services/:id/remark | UpdateServiceRemark |
| GET | /api/client/services/:id/renewals | GetServiceRenewPreview |
| POST | /api/client/services/:id/renewals | CreateRenewOrder |
| POST | /api/client/services/:id/power-actions | PowerService |
| POST | /api/client/services/:id/password-resets | ResetServicePassword |
| POST | /api/client/services/:id/reinstallations | ReinstallService |
| GET | /api/client/services/:id/module-status | GetServiceStatus |
| GET | /api/client/services/:id/operation-logs | GetServiceOperationLogs |
| GET | /api/client/services/:id/config | GetServiceConfig |
| GET | /api/client/services/:id/reinstallations/options | GetServiceReinstallOptions |
| GET | /api/client/services/:id/upgrades | GetServiceUpgradePreview |
| POST | /api/client/services/:id/upgrades/quotes | QuoteServiceUpgrade |
| POST | /api/client/services/:id/upgrades/orders | CreateServiceUpgradeOrder |
| PUT | /api/client/services/:id/renewals/auto | UpdateAutoRenew |

### 订单管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/orders | GetUserOrders |
| GET | /api/client/orders/summary | GetOrderSummary |
| GET | /api/client/orders/:id | GetUserOrder |
| POST | /api/client/orders/:id/cancel | CancelUserOrder |

### 认证

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /api/client/login | authHandler |
| POST | /api/client/auth/reset-password | authHandler |
| POST | /api/client/auth/captcha | authHandler |
| POST | /api/client/auth/login-by-code | authHandler |
| GET | /api/client/auth/info | authHandler |
| POST | /api/client/auth/logout | authHandler |
| PUT | /api/client/auth/profile | authHandler |
| PUT | /api/client/auth/phone | authHandler |
| PUT | /api/client/auth/email | authHandler |
| GET | /api/client/auth/notification-preferences | GetNotificationPreferences |
| PUT | /api/client/auth/notification-preferences | UpdateNotificationPreferences |

### 账单管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/invoices | GetUserInvoices |
| GET | /api/client/invoices/summary | GetInvoiceSummary |
| GET | /api/client/invoices/:id | GetUserInvoice |
| POST | /api/client/invoices/:id/cancellations | CancelUserInvoice |
| POST | /api/client/invoices/:id/pay/balance | PayInvoiceByBalance |
| POST | /api/client/invoices/combines | CombineInvoices |
| POST | /api/client/invoices/:id/fund | FundInvoice |

### 购物车

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/cart | GetCart |
| POST | /api/client/cart | AddToCart |
| PUT | /api/client/cart/:id | UpdateCartItem |
| DELETE | /api/client/cart/:id | RemoveCartItem |
| DELETE | /api/client/cart | ClearCart |
| POST | /api/client/cart/checkout | Checkout |

### 通知管理

| 方法 | 路径 | Handler |
|------|------|---------|
| GET | /api/client/notifications | GetUserNotifications |
| GET | /api/client/notifications/unread-count | GetNotificationUnreadCount |
| PUT | /api/client/notifications/:id/read-state | MarkNotificationRead |
| POST | /api/client/notifications/mark-all-read | MarkAllNotificationsRead |

---

## 三、API格式规范

### 成功响应
```json
{"code": 0, "message": "success", "data": {}}
```

### 错误响应
```json
{"code": 400, "message": "错误描述", "data": null}
```

### 分页响应
```json
{"code": 0, "message": "success", "data": {"list": [], "total": 100, "page": 1, "page_size": 20}}
```

### 错误码
| code | 含义 |
|------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器错误 |
| 502 | 插件引擎离线 |

---

## 四、安全措施

| 措施 | 说明 |
|------|------|
| JWT认证 | Bearer Token，支持黑名单 |
| Admin/Client分离 | AdminRequired/ClientRequired中间件 |
| IDOR防护 | 所有Client端查询都有WHERE user_id = ? |
| 0元购防护 | 价格服务端计算，不信任前端 |
| 登录防爆破 | 账号级冻结（5次错密码锁6小时）+ IP软锁 |
| 验证码限流 | 目标维度（configurable_rate/秒）+ IP维度（10次/分钟） |
| SQL注入防护 | 全部GORM参数化查询 |
| 余额操作原子化 | gorm.Expr原子操作防并发 |