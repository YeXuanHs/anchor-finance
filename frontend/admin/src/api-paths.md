# AnchorFinance Admin Frontend API Paths

This document lists all API paths used in the AnchorFinance admin frontend code.

## Base Configuration

The HTTP client is configured in `src/utils/http/index.ts` using Axios with:
- Base URL: `VITE_API_URL` environment variable
- Timeout: 15000ms
- Credentials: `VITE_WITH_CREDENTIALS` environment variable

## API Paths by Module

### 1. Authentication (`/api/auth`)
- `POST /api/auth/login` - User login

### 2. User Management (`/api/user`)
- `GET /api/user/list` - Get user list
- `GET /api/user/info` - Get current user info

### 3. System Menu (`/api/v3/system`)
- `GET /api/v3/system/menus` - Get system menus

### 4. Client Management (`/api/admin/user-manage`)
- `GET /api/admin/user-manage/search` - Search clients
- `GET /api/admin/user-manage/{id}` - Get client details
- `POST /api/admin/user-manage` - Create client
- `PUT /api/admin/user-manage/{id}/profile` - Update client profile
- `POST /api/admin/user-manage/{id}/login-as` - Login as client
- `GET /api/admin/user-manage/cancel-requests` - Get cancellation requests
- `POST /api/admin/user-manage/cancel-requests/{id}` - Process cancellation request

### 5. Client Groups (`/api/admin/client-groups`)
- `GET /api/admin/client-groups` - Get client groups
- `GET /api/admin/client-groups/{id}` - Get client group details
- `POST /api/admin/client-groups` - Create client group
- `PUT /api/admin/client-groups/{id}` - Update client group
- `DELETE /api/admin/client-groups/{id}` - Delete client group

### 6. Client Services (`/api/admin/client-services`)
- `GET /api/admin/client-services` - Get client services
- `GET /api/admin/client-services?user_id={id}` - Get services by user
- `POST /api/admin/client-services/{id}/suspend` - Suspend service
- `POST /api/admin/client-services/{id}/terminate` - Terminate service
- `POST /api/admin/client-services/{id}/renew` - Renew service

### 7. Client Resources (`/api/admin/client-resources`)
- `GET /api/admin/client-resources` - Get client resources
- `POST /api/admin/client-resources` - Create client resource
- `PUT /api/admin/client-resources/{id}` - Update client resource
- `DELETE /api/admin/client-resources/{id}` - Delete client resource

### 8. Client Tracks (`/api/admin/client-tracks`)
- `GET /api/admin/client-tracks` - Get client tracks
- `POST /api/admin/client-tracks` - Create client track
- `PUT /api/admin/client-tracks/{id}` - Update client track
- `DELETE /api/admin/client-tracks/{id}` - Delete client track

### 9. Client Custom Fields (`/api/admin/custom-fields`)
- `GET /api/admin/custom-fields` - Get custom fields
- `POST /api/admin/custom-fields` - Create custom field
- `PUT /api/admin/custom-fields/{id}` - Update custom field
- `DELETE /api/admin/custom-fields/{id}` - Delete custom field

### 10. Client Attachments (`/api/admin/attachments`)
- `GET /api/admin/attachments` - Get attachments
- `POST /api/admin/attachments/upload` - Upload attachment
- `GET /api/admin/attachments/{id}/download` - Download attachment
- `PUT /api/admin/attachments/{id}` - Update attachment
- `DELETE /api/admin/attachments/{id}` - Delete attachment

### 11. Client Emails (`/api/admin/emails`)
- `GET /api/admin/emails` - Get emails
- `GET /api/admin/emails/{id}` - Get email details
- `POST /api/admin/clients/{id}/emails` - Send email to client
- `DELETE /api/admin/emails/{id}` - Delete email
- `GET /api/admin/email-attachments/{id}/download` - Download email attachment

### 12. Client Notifications (`/api/admin/notifications`)
- `GET /api/admin/notifications/logs` - Get notification logs
- `POST /api/admin/clients/{id}/notifications` - Send notification to client

### 13. Client Authentication (`/api/admin/certifications`)
- `GET /api/admin/certifications` - Get certifications
- `POST /api/admin/certifications/{id}/review` - Review certification

### 14. Client CRM (`/api/admin/crm`)
- `GET /api/admin/crm/follow-records` - Get follow records
- `POST /api/admin/crm/follow-records` - Create follow record
- `GET /api/admin/crm/opportunities` - Get opportunities
- `GET /api/admin/crm/contracts` - Get contracts

### 15. Client Remarks (`/api/admin/user-remarks`)
- `GET /api/admin/user-remarks` - Get user remarks
- `POST /api/admin/user-remarks` - Create user remark
- `DELETE /api/admin/user-remarks/{id}` - Delete user remark

### 16. Client Blacklist (`/api/admin/blacklist`)
- `GET /api/admin/blacklist` - Get blacklist
- `POST /api/admin/blacklist` - Add to blacklist
- `PUT /api/admin/blacklist/{id}` - Update blacklist entry
- `DELETE /api/admin/blacklist/{id}` - Remove from blacklist

### 17. Client Contacts (`/api/admin/contacts`)
- `GET /api/admin/contacts` - Get contacts
- `POST /api/admin/contacts` - Create contact
- `PUT /api/admin/contacts/{id}` - Update contact
- `DELETE /api/admin/contacts/{id}` - Delete contact

### 18. Client Resources (`/api/admin/resource-pools`)
- `GET /api/admin/resource-pools` - Get resource pools
- `POST /api/admin/resource-pools` - Create resource pool
- `PUT /api/admin/resource-pools/{id}` - Update resource pool
- `DELETE /api/admin/resource-pools/{id}` - Delete resource pool

### 19. Orders (`/api/admin/orders`)
- `GET /api/admin/orders` - Get orders
- `GET /api/admin/orders/{id}` - Get order details
- `POST /api/admin/orders` - Create order
- `PUT /api/admin/orders/{id}/status` - Update order status
- `POST /api/admin/orders/{id}/activate` - Activate order

### 20. Order Renewal (`/api/admin/multi-renew`)
- `GET /api/admin/multi-renew` - Get renewal orders
- `POST /api/admin/multi-renew` - Create renewal order
- `POST /api/admin/multi-renew/{id}/execute` - Execute renewal
- `POST /api/admin/multi-renew/{id}/cancel` - Cancel renewal

### 21. Traffic Orders (`/api/admin/traffic-orders`)
- `GET /api/admin/traffic-orders` - Get traffic orders
- `POST /api/admin/traffic-orders/{id}/renew` - Renew traffic order

### 22. Invoices (`/api/admin/invoices`)
- `GET /api/admin/invoices` - Get invoices
- `GET /api/admin/invoices/{id}` - Get invoice details
- `GET /api/admin/invoices/{id}/payments` - Get invoice payments
- `POST /api/admin/invoices/{id}/email` - Send invoice email
- `POST /api/admin/invoices/{id}/refund` - Refund invoice
- `POST /api/admin/invoices/{id}/cancel` - Cancel invoice
- `POST /api/admin/invoices/{id}/audit` - Audit invoice
- `POST /api/admin/invoices/{id}/issued` - Mark invoice as issued
- `GET /api/admin/invoices/status/{status}` - Get invoices by status
- `GET /api/admin/invoices/audit` - Get invoices for audit

### 23. Invoice Items (`/api/admin/invoice-items`)
- `GET /api/admin/invoice-items` - Get invoice items

### 24. Invoice Config (`/api/admin/config/invoice`)
- `GET /api/admin/config/invoice` - Get invoice config
- `PUT /api/admin/config/invoice` - Update invoice config
- `GET /api/admin/config/invoice/express` - Get express config
- `POST /api/admin/config/invoice/express` - Create express config
- `PUT /api/admin/config/invoice/express/{id}` - Update express config
- `DELETE /api/admin/config/invoice/express/{id}` - Delete express config

### 25. Products (`/api/admin/products`)
- `GET /api/admin/products` - Get products
- `GET /api/admin/products/{id}` - Get product details
- `POST /api/admin/products` - Create product
- `PUT /api/admin/products/{id}` - Update product
- `DELETE /api/admin/products/{id}` - Delete product

### 26. Product Groups (`/api/admin/product-groups`)
- `GET /api/admin/product-groups` - Get product groups

### 27. Tickets (`/api/admin/tickets`)
- `GET /api/admin/tickets` - Get tickets
- `GET /api/admin/tickets/{id}` - Get ticket details
- `POST /api/admin/tickets/{id}/reply` - Reply to ticket
- `POST /api/admin/tickets/{id}/transfer` - Transfer ticket
- `POST /api/admin/tickets/{id}/close` - Close ticket
- `GET /api/admin/tickets/{id}/attachments` - Get ticket attachments
- `GET /api/admin/tickets/statistics` - Get ticket statistics

### 28. Ticket Departments (`/api/admin/ticket-depts`)
- `GET /api/admin/ticket-depts` - Get ticket departments
- `GET /api/admin/ticket-depts/{id}` - Get ticket department details
- `POST /api/admin/ticket-depts` - Create ticket department
- `PUT /api/admin/ticket-depts/{id}` - Update ticket department
- `DELETE /api/admin/ticket-depts/{id}` - Delete ticket department

### 29. Ticket Pre-reply (`/api/admin/ticket-prereply`)
- `GET /api/admin/ticket-prereply` - Get ticket pre-replies
- `POST /api/admin/ticket-prereply/replies` - Create ticket pre-reply
- `PUT /api/admin/ticket-prereply/replies/{id}` - Update ticket pre-reply
- `DELETE /api/admin/ticket-prereply/replies/{id}` - Delete ticket pre-reply
- `GET /api/admin/ticket-prereply/password-config` - Get password config
- `POST /api/admin/ticket-prereply/password-config` - Update password config

### 30. Ticket Delivery Rules (`/api/admin/ticket-deliver`)
- `GET /api/admin/ticket-deliver/rules` - Get delivery rules
- `POST /api/admin/ticket-deliver/rules` - Create delivery rule
- `PUT /api/admin/ticket-deliver/rules/{id}` - Update delivery rule
- `DELETE /api/admin/ticket-deliver/rules/{id}` - Delete delivery rule

### 31. Ticket Statuses (`/api/admin/ticket-statuses`)
- `GET /api/admin/ticket-statuses` - Get ticket statuses
- `POST /api/admin/ticket-statuses` - Create ticket status
- `PUT /api/admin/ticket-statuses/{id}` - Update ticket status
- `DELETE /api/admin/ticket-statuses/{id}` - Delete ticket status

### 32. AI Ticket (`/api/admin/ai-ticket`)
- `GET /api/admin/ai-ticket/dashboard` - Get AI ticket dashboard
- `PUT /api/admin/ai-ticket/dashboard` - Update AI ticket dashboard
- `GET /api/admin/ai-ticket/knowledge` - Get AI ticket knowledge
- `POST /api/admin/ai-ticket/knowledge` - Create AI ticket knowledge
- `PUT /api/admin/ai-ticket/knowledge/{id}` - Update AI ticket knowledge
- `DELETE /api/admin/ai-ticket/knowledge/{id}` - Delete AI ticket knowledge
- `POST /api/admin/ai-ticket/knowledge/import` - Import AI ticket knowledge
- `GET /api/admin/ai-ticket/rules` - Get AI ticket rules
- `POST /api/admin/ai-ticket/rules` - Create AI ticket rule
- `PUT /api/admin/ai-ticket/rules/{id}` - Update AI ticket rule
- `DELETE /api/admin/ai-ticket/rules/{id}` - Delete AI ticket rule
- `GET /api/admin/ai-ticket/queue` - Get AI ticket queue
- `GET /api/admin/ai-ticket/queue/stats` - Get AI ticket queue stats
- `GET /api/admin/ai-ticket/process-logs` - Get AI ticket process logs
- `POST /api/admin/ai-ticket/test` - Test AI ticket

### 33. Knowledge Base (`/api/admin/knowledge`)
- `GET /api/admin/knowledge/categories` - Get knowledge categories
- `POST /api/admin/knowledge/categories` - Create knowledge category
- `PUT /api/admin/knowledge/categories/{id}` - Update knowledge category
- `DELETE /api/admin/knowledge/categories/{id}` - Delete knowledge category
- `GET /api/admin/knowledge/articles` - Get knowledge articles
- `POST /api/admin/knowledge/articles` - Create knowledge article
- `PUT /api/admin/knowledge/articles/{id}` - Update knowledge article
- `DELETE /api/admin/knowledge/articles/{id}` - Delete knowledge article

### 34. Affiliate (`/api/admin/affiliate`)
- `GET /api/admin/affiliate` - Get affiliates
- `GET /api/admin/affiliate/{id}` - Get affiliate details
- `POST /api/admin/affiliate` - Create affiliate
- `PUT /api/admin/affiliate/{id}` - Update affiliate
- `DELETE /api/admin/affiliate/{id}` - Delete affiliate
- `GET /api/admin/affiliate/withdraw-records` - Get withdraw records
- `POST /api/admin/affiliate/withdraws/{id}/process` - Process withdraw
- `GET /api/admin/affiliate/user-affi-record` - Get user affiliate record
- `POST /api/admin/affiliate/records/{id}/confirm` - Confirm affiliate record

### 35. Sales (`/api/admin/sales`)
- `GET /api/admin/sales` - Get sales
- `GET /api/admin/sales/{id}` - Get sale details
- `POST /api/admin/sales` - Create sale
- `PUT /api/admin/sales/{id}` - Update sale
- `DELETE /api/admin/sales/{id}` - Delete sale
- `GET /api/admin/sales/statistics` - Get sales statistics
- `GET /api/admin/sales/records` - Get sales records
- `GET /api/admin/sales/admin-list` - Get sales admin list

### 36. Agents (`/api/admin/agents`)
- `GET /api/admin/agents` - Get agents
- `GET /api/admin/agents/{id}` - Get agent details
- `POST /api/admin/agents` - Create agent
- `PUT /api/admin/agents/{id}` - Update agent
- `DELETE /api/admin/agents/{id}` - Delete agent

### 37. Contracts (`/api/admin/contracts`)
- `GET /api/admin/contracts` - Get contracts
- `GET /api/admin/contracts/{id}` - Get contract details
- `POST /api/admin/contracts` - Create contract
- `PUT /api/admin/contracts/{id}` - Update contract
- `DELETE /api/admin/contracts/{id}` - Delete contract
- `GET /api/admin/contracts/{id}/pdf` - Get contract PDF

### 38. Client Care (`/api/admin/client-care`)
- `GET /api/admin/client-care/rules` - Get client care rules
- `POST /api/admin/client-care/rules` - Create client care rule
- `PUT /api/admin/client-care/rules/{id}` - Update client care rule
- `DELETE /api/admin/client-care/rules/{id}` - Delete client care rule
- `GET /api/admin/client-care/logs` - Get client care logs

### 39. Upstream Providers (`/api/admin/upstream`)
- `GET /api/admin/upstream/providers` - Get upstream providers
- `GET /api/admin/upstream/providers/{id}` - Get upstream provider details
- `POST /api/admin/upstream/providers` - Create upstream provider
- `PUT /api/admin/upstream/providers/{id}` - Update upstream provider
- `DELETE /api/admin/upstream/providers/{id}` - Delete upstream provider
- `POST /api/admin/upstream/providers/{id}/test` - Test upstream provider

### 40. Batch Sync (`/api/admin/batch-sync`)
- `GET /api/admin/batch-sync` - Get batch sync tasks
- `POST /api/admin/batch-sync/{id}/execute` - Execute batch sync

### 41. Send Messages (`/api/admin/messages`)
- `POST /api/admin/messages/send` - Send message
- `GET /api/admin/messages/batch/records` - Get batch message records

### 42. Link Causes (`/api/admin/link-causes`)
- `GET /api/admin/link-causes/tree` - Get link causes tree
- `GET /api/admin/link-causes/{id}` - Get link cause details
- `POST /api/admin/link-causes` - Create link cause
- `PUT /api/admin/link-causes/{id}` - Update link cause
- `DELETE /api/admin/link-causes/{id}` - Delete link cause

### 43. Content Knowledge (`/api/admin/knowledge/articles`)
- `GET /api/admin/knowledge/articles` - Get knowledge articles
- `POST /api/admin/knowledge/articles` - Create knowledge article
- `PUT /api/admin/knowledge/articles/{id}` - Update knowledge article
- `DELETE /api/admin/knowledge/articles/{id}` - Delete knowledge article

### 44. Marketing - Promo Codes (`/api/admin/promo-codes`)
- `GET /api/admin/promo-codes` - Get promo codes
- `GET /api/admin/promo-codes/{id}` - Get promo code details
- `POST /api/admin/promo-codes` - Create promo code
- `PUT /api/admin/promo-codes/{id}` - Update promo code
- `DELETE /api/admin/promo-codes/{id}` - Delete promo code

### 45. Marketing - Vouchers (`/api/admin/voucher`)
- `GET /api/admin/voucher-list` - Get vouchers
- `GET /api/admin/voucher-detail/{id}` - Get voucher details
- `GET /api/admin/voucher-rate` - Get voucher rate
- `POST /api/admin/voucher-status` - Update voucher status
- `POST /api/admin/voucher-rate` - Update voucher rate

### 46. Marketing - Coupons (`/api/admin/coupons`)
- `GET /api/admin/coupons` - Get coupons
- `GET /api/admin/coupons/{id}` - Get coupon details
- `POST /api/admin/coupons` - Create coupon
- `PUT /api/admin/coupons/{id}` - Update coupon
- `DELETE /api/admin/coupons/{id}` - Delete coupon

### 47. Finance - Credit (`/api/admin/credit`)
- `GET /api/admin/credit/index` - Get credit index
- `POST /api/admin/credit/users/{id}/settings` - Update user credit settings
- `POST /api/admin/credit/users/{id}/adjust` - Adjust user credit

### 48. Finance - Accounts (`/api/admin/accounts`)
- `GET /api/admin/accounts` - Get accounts
- `GET /api/admin/payment-gateways` - Get payment gateways

### 49. Finance - Withdraw (`/api/admin/affiliate/withdraws`)
- `GET /api/admin/affiliate/withdraw-records` - Get withdraw records
- `POST /api/admin/affiliate/withdraws/{id}/process` - Process withdraw

### 50. Statistics (`/api/admin/reports`)
- `GET /api/admin/reports/dashboard` - Get dashboard statistics
- `GET /api/admin/reports/revenue` - Get revenue statistics
- `GET /api/admin/reports/users` - Get user statistics
- `GET /api/admin/reports/product-income` - Get product income statistics
- `GET /api/admin/reports/product-income/trend` - Get product income trend
- `GET /api/admin/reports/product-income/comparison` - Get product income comparison
- `GET /api/admin/reports/revenue-ranking` - Get revenue ranking
- `GET /api/admin/reports/revenue-ranking/comparison` - Get revenue ranking comparison
- `GET /api/admin/reports/tickets` - Get ticket statistics
- `GET /api/admin/reports/new-client-statistics` - Get new client statistics
- `GET /api/admin/reports/year-income-statistics` - Get year income statistics

### 51. System - Admins (`/api/admin/users`)
- `GET /api/admin/users` - Get users
- `GET /api/admin/users/{id}` - Get user details
- `POST /api/admin/users` - Create user
- `PUT /api/admin/users/{id}` - Update user
- `DELETE /api/admin/users/{id}` - Delete user

### 52. System - RBAC (`/api/admin/rbac`)
- `GET /api/admin/rbac/roles` - Get roles
- `GET /api/admin/rbac/roles/{id}` - Get role details
- `POST /api/admin/rbac/roles` - Create role
- `PUT /api/admin/rbac/roles/{id}` - Update role
- `DELETE /api/admin/rbac/roles/{id}` - Delete role
- `GET /api/admin/rbac/permissions` - Get permissions
- `PUT /api/admin/rbac/permissions` - Update permissions

### 53. System - Languages (`/api/admin/languages`)
- `GET /api/admin/languages` - Get languages
- `POST /api/admin/languages` - Create language
- `PUT /api/admin/languages/{id}` - Update language
- `POST /api/admin/languages/{id}/default` - Set default language
- `GET /api/admin/languages/{code}/translations` - Get translations
- `POST /api/admin/languages/{code}/translations` - Update translations

### 54. System - Config Servers (`/api/admin/config/servers`)
- `GET /api/admin/config/servers` - Get config servers
- `POST /api/admin/config/servers` - Create config server
- `PUT /api/admin/config/servers/{id}` - Update config server
- `DELETE /api/admin/config/servers/{id}` - Delete config server
- `GET /api/admin/config-servers/test-link/{id}` - Test server connection

### 55. System - Config Options (`/api/admin/config-options`)
- `GET /api/admin/config-options/groups-list` - Get config option groups
- `GET /api/admin/config-options/search-page` - Search config options
- `POST /api/admin/config-options/add-options` - Add config options
- `POST /api/admin/config-options/edit-config` - Edit config
- `DELETE /api/admin/config-options/options` - Delete config options

### 56. System - Cron Tasks (`/api/admin/cron-tasks`)
- `GET /api/admin/cron-tasks` - Get cron tasks
- `GET /api/admin/cron-tasks/{id}` - Get cron task details
- `POST /api/admin/cron-tasks` - Create cron task
- `PUT /api/admin/cron-tasks/{id}` - Update cron task
- `DELETE /api/admin/cron-tasks/{id}` - Delete cron task
- `POST /api/admin/cron-tasks/{id}/run` - Run cron task
- `GET /api/admin/cron-tasks/{id}/logs` - Get cron task logs

### 57. System - Log Records (`/api/admin/log-records`)
- `GET /api/admin/log-records` - Get log records
- `POST /api/admin/log-records/export` - Export log records

### 58. System - Log Cleaner (`/api/admin/log-cleaner`)
- `GET /api/admin/log-cleaner/rules` - Get log cleaner rules
- `GET /api/admin/log-cleaner/stats` - Get log cleaner stats
- `POST /api/admin/log-cleaner/rules` - Create log cleaner rule
- `PUT /api/admin/log-cleaner/rules/{id}` - Update log cleaner rule
- `DELETE /api/admin/log-cleaner/rules/{id}` - Delete log cleaner rule
- `POST /api/admin/log-cleaner/rules/{id}/run` - Run log cleaner rule

### 59. System - Operation Logs (`/api/admin/system-logs`)
- `DELETE /api/admin/system-logs/clear-level` - Clear system logs by level

### 60. System - Uploads (`/api/admin/uploads`)
- `GET /api/admin/uploads` - Get uploads
- `POST /api/admin/uploads` - Upload file
- `GET /api/admin/uploads/{id}/download` - Download file
- `DELETE /api/admin/uploads/{id}` - Delete file

### 61. System - Site Settings (`/api/admin/setting/site`)
- `GET /api/admin/setting/site` - Get site settings
- `PUT /api/admin/setting/site` - Update site settings

### 62. System - Security Config (`/api/admin/config/certifi`)
- `GET /api/admin/config/certifi` - Get security config
- `PUT /api/admin/config/certifi` - Update security config

### 63. System - Fund Config (`/api/admin/config/recharge`)
- `GET /api/admin/config/recharge` - Get fund config
- `PUT /api/admin/config/recharge` - Update fund config

### 64. System - Affiliate Config (`/api/admin/config/affiliate`)
- `GET /api/admin/config/affiliate` - Get affiliate config
- `PUT /api/admin/config/affiliate` - Update affiliate config

### 65. System - SMS (`/api/admin/sms`)
- `GET /api/admin/sms/logs` - Get SMS logs
- `GET /api/admin/sms/templates` - Get SMS templates
- `POST /api/admin/sms/templates` - Create SMS template
- `PUT /api/admin/sms/templates/{id}` - Update SMS template
- `DELETE /api/admin/sms/templates/{id}` - Delete SMS template
- `POST /api/admin/sms/send` - Send SMS

### 66. System - SMS Config (`/api/admin/config/messages/mobile`)
- `GET /api/admin/config/messages/mobile` - Get SMS config
- `PUT /api/admin/config/messages/mobile` - Update SMS config
- `POST /api/admin/config/messages/mobile/test` - Test SMS config

### 67. System - User Tastes (`/api/admin/user-tastes`)
- `GET /api/admin/user-tastes` - Get user tastes
- `PUT /api/admin/user-tastes` - Update user tastes

### 68. System - API Management (`/api/admin/api-keys`)
- `GET /api/admin/api-keys` - Get API keys
- `POST /api/admin/api-keys` - Create API key
- `DELETE /api/admin/api-keys/{id}` - Delete API key

### 69. System - Maintenance (`/api/admin/maintenance`)
- `GET /api/admin/maintenance/status` - Get maintenance status
- `PUT /api/admin/maintenance/settings` - Update maintenance settings

### 70. System - System Info (`/api/admin/system`)
- `GET /api/admin/system/info` - Get system info
- `GET /api/admin/system/check-update` - Check for updates
- `GET /api/admin/system/database` - Get database info
- `POST /api/admin/system/database/optimize` - Optimize database
- `POST /api/admin/system/database/backup` - Backup database

### 71. System - Data Migration (`/api/admin/system/data-migration`)
- `GET /api/admin/system/data-migration` - Get data migration tasks
- `POST /api/admin/system/data-migration` - Create data migration task
- `POST /api/admin/system/data-migration/tasks/{id}/start` - Start migration task
- `POST /api/admin/system/data-migration/tasks/{id}/pause` - Pause migration task
- `DELETE /api/admin/system/data-migration/tasks/{id}` - Delete migration task

### 72. System - Run Map (`/api/admin/run-map`)
- `GET /api/admin/run-map` - Get run map
- `POST /api/admin/run-map/{id}/repeat` - Repeat run map task

### 73. System - Captcha (`/api/admin/captcha-config`)
- `GET /api/admin/captcha-config` - Get captcha config
- `PUT /api/admin/captcha-config/basic` - Update captcha config
- `POST /api/admin/system/captcha/test` - Test captcha

### 74. System - Second Verify (`/api/admin/config/second-verify`)
- `GET /api/admin/config/second-verify` - Get second verify config
- `PUT /api/admin/config/second-verify` - Update second verify config

### 75. System - DCIM (`/api/admin/dcim`)
- `GET /api/admin/dcim/servers` - Get DCIM servers
- `GET /api/admin/dcim/servers/{id}` - Get DCIM server details
- `POST /api/admin/dcim/servers` - Create DCIM server
- `PUT /api/admin/dcim/servers/{id}` - Update DCIM server
- `DELETE /api/admin/dcim/servers/{id}` - Delete DCIM server
- `POST /api/admin/dcim/servers/{id}/refresh` - Refresh DCIM server

### 76. System - DCIM Cloud (`/api/admin/dcim-cloud`)
- `GET /api/admin/dcim-cloud/servers` - Get DCIM cloud servers
- `GET /api/admin/dcim-cloud/servers/{id}` - Get DCIM cloud server details
- `POST /api/admin/dcim-cloud/servers` - Create DCIM cloud server
- `PUT /api/admin/dcim-cloud/servers/{id}` - Update DCIM cloud server
- `DELETE /api/admin/dcim-cloud/servers/{id}` - Delete DCIM cloud server
- `POST /api/admin/dcim-cloud/servers/{id}/test` - Test DCIM cloud server
- `POST /api/admin/dcim-cloud/servers/{id}/sync` - Sync DCIM cloud server

### 77. Cancel Reasons (`/api/admin/cancel-reasons`)
- `GET /api/admin/cancel-reasons` - Get cancel reasons
- `POST /api/admin/cancel-reasons` - Create cancel reason
- `DELETE /api/admin/cancel-reasons/{id}` - Delete cancel reason

### 78. Users (`/api/admin/users`)
- `GET /api/admin/users` - Get users
- `GET /api/admin/users/{id}` - Get user details
- `POST /api/admin/users` - Create user
- `PUT /api/admin/users/{id}` - Update user
- `DELETE /api/admin/users/{id}` - Delete user

## Summary

Total unique API paths found: ~200+

The API follows a RESTful pattern with:
- Base path: `/api/admin/`
- Standard CRUD operations (GET, POST, PUT, DELETE)
- Nested resources for related entities
- Query parameters for filtering and pagination
- Path parameters for specific resource IDs
