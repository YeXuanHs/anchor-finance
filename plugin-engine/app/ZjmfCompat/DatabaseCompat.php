<?php

namespace App\ZjmfCompat;

use PDO;

/**
 * 数据库兼容层 - 纯PDO版本
 * 完整映射：表名 + 字段名 + 内容格式
 */
class DatabaseCompat
{
    private static $pdo = null;

    // ========== 表名映射（zjmf shd_前缀 → 我们的表名）==========
    private static $tableMap = [
        // 核心表
        'shd_clients' => 'users',
        'shd_host' => 'services',
        'shd_invoices' => 'invoices',
        'shd_orders' => 'orders',
        'shd_ticket' => 'tickets',
        'shd_product' => 'products',
        'shd_servers' => 'servers',

        // 插件相关
        'shd_plugin' => 'plugins',
        'shd_hook_plugin' => 'hook_bindings',

        // 配置
        'shd_configuration' => 'settings',

        // 工单
        'shd_ticket_department' => 'ticket_departments',
        'shd_ticket_status' => 'ticket_statuses',

        // 自定义字段
        'shd_customfields' => 'custom_fields',
        'shd_customfieldsvalues' => 'custom_field_values',

        // 优惠码
        'shd_promo_code' => 'promo_codes',

        // 内容
        'shd_news' => 'news',
        'shd_downloads' => 'downloads',
        'shd_downloadcats' => 'download_categories',
        'shd_knowledge_base' => 'knowledge_articles',
        'shd_knowledge_categories' => 'knowledge_categories',

        // 无shd_前缀的兼容（有些插件可能不带前缀）
        'clients' => 'users',
        'host' => 'services',
        'invoices' => 'invoices',
        'orders' => 'orders',
        'ticket' => 'tickets',
        'product' => 'products',
        'servers' => 'servers',
        'plugin' => 'plugins',
        'hook_plugin' => 'hook_bindings',
        'configuration' => 'settings',
        'ticket_department' => 'ticket_departments',
        'ticket_status' => 'ticket_statuses',
        'customfields' => 'custom_fields',
        'customfieldsvalues' => 'custom_field_values',
        'promo_code' => 'promo_codes',
        'news' => 'news',
        'downloads' => 'downloads',
        'downloadcats' => 'download_categories',
        'knowledge_base' => 'knowledge_articles',
        'knowledge_categories' => 'knowledge_categories',
    ];

    // ========== 字段名映射（zjmf字段 => 我们的字段）==========
    // 严格对照 zjmfv376/public/install/thinkcmf.sql 和我们的Go model
    private static $fieldMap = [
        // shd_clients → users
        'shd_clients' => [
            'id' => 'id',
            'username' => 'username',
            'email' => 'email',
            'password' => 'password_hash',
            'phonenumber' => 'phone',
            'companyname' => 'company',
            'status' => 'status',
            'credit' => 'balance',           // zjmf credit = 我们的 balance
            'groupid' => 'group_id',
            'lastlogin' => 'last_login_at',  // int时间戳 → datetime
            'lastloginip' => 'last_login_ip',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
            'country' => 'country',
            'province' => 'province',
            'city' => 'city',
            'address1' => 'address',
            'postcode' => 'postcode',
            'language' => 'language',
            'avatar' => 'avatar',
            'notes' => 'notes',
            'sale_id' => 'sale_id',
        ],
        // shd_host → services
        'shd_host' => [
            'id' => 'id',
            'uid' => 'user_id',              // zjmf uid = 我们的 user_id
            'orderid' => 'order_id',
            'productid' => 'product_id',
            'serverid' => 'upstream_id',
            'domain' => 'domain',
            'domainstatus' => 'status',      // zjmf domainstatus = 我们的 status
            'username' => 'username',
            'password' => 'password_encrypted',
            'billingcycle' => 'billing_cycle',
            'amount' => 'amount',
            'nextduedate' => 'next_due_date', // int时间戳 → datetime
            'regdate' => 'created_at',        // int时间戳 → datetime
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
            'dedicatedip' => 'dedicated_ip',
            'assignedips' => 'assigned_ips',
            'payment' => 'payment_method',
            'firstpaymentamount' => 'first_payment_amount',
            'notes' => 'notes',
            'remark' => 'remark',
            'os' => 'os',
            'port' => 'port',
        ],
        // shd_invoices → invoices
        'shd_invoices' => [
            'id' => 'id',
            'uid' => 'user_id',              // zjmf uid = 我们的 user_id
            'invoice_num' => 'invoice_no',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
            'due_time' => 'due_date',
            'paid_time' => 'paid_at',
            'subtotal' => 'subtotal',
            'total' => 'amount',             // zjmf total = 我们的 amount
            'status' => 'status',
            'payment' => 'payment_method',
            'notes' => 'notes',
            'tax' => 'tax',
            'credit' => 'credit',
            'type' => 'type',
        ],
        // shd_orders → orders
        'shd_orders' => [
            'id' => 'id',
            'uid' => 'user_id',
            'ordernum' => 'order_no',
            'status' => 'status',
            'amount' => 'amount',
            'payment' => 'payment_method',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
            'pay_time' => 'paid_at',
            'promo_code' => 'promo_code',
            'invoiceid' => 'invoice_id',
            'notes' => 'notes',
        ],
        // shd_ticket → tickets
        'shd_ticket' => [
            'id' => 'id',
            'uid' => 'user_id',
            'tid' => 'ticket_no',
            'dptid' => 'department_id',      // zjmf dptid = 我们的 department_id
            'title' => 'subject',
            'content' => 'message',
            'status' => 'status',
            'priority' => 'priority',
            'admin_id' => 'assigned_to',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
            'last_reply_time' => 'last_reply_at',
            'host_id' => 'service_id',
            'name' => 'contact_name',
            'email' => 'contact_email',
            'attachment' => 'attachments',
            'star' => 'rating',
        ],
        // shd_ticket_department → ticket_departments
        'shd_ticket_department' => [
            'id' => 'id',
            'name' => 'name',
            'description' => 'description',
            'email' => 'email',
            'order' => 'sort_order',
            'hidden' => 'is_hidden',
        ],
        // shd_ticket_status → ticket_statuses
        'shd_ticket_status' => [
            'id' => 'id',
            'title' => 'name',
            'sortorder' => 'sort_order',
            'color' => 'color',
        ],
        // shd_product → products
        'shd_product' => [
            'id' => 'id',
            'name' => 'name',
            'description' => 'description',
            'type' => 'type',
            'gid' => 'group_id',
            'paytype' => 'pay_type',
            'stockcontrol' => 'stock_control',
            'qty' => 'stock',
            'status' => 'status',
            'servertype' => 'server_type',
        ],
        // shd_plugin → plugins
        'shd_plugin' => [
            'id' => 'id',
            'name' => 'slug',                // zjmf name = 我们的 slug
            'title' => 'name',               // zjmf title = 我们的 name
            'description' => 'description',
            'status' => 'status',
            'config' => 'config',
            'hooks' => 'hook_list',
            'author' => 'author',
            'version' => 'version',
            'module' => 'domain',
        ],
        // shd_configuration → settings
        'shd_configuration' => [
            'setting' => 'key',              // zjmf setting = 我们的 key
            'value' => 'value',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
        ],
        // shd_customfields → custom_fields
        'shd_customfields' => [
            'id' => 'id',
            'fieldname' => 'name',
            'fieldtype' => 'type',
            'description' => 'description',
            'adminonly' => 'admin_only',
            'required' => 'required',
            'sortorder' => 'sort_order',
            'showorder' => 'display_order',
            'fieldoptions' => 'options',
            'relid' => 'related_id',
            'type' => 'scope',
        ],
        // shd_currencies → currencies
        'shd_currencies' => [
            'id' => 'id',
            'code' => 'code',
            'prefix' => 'prefix',
            'suffix' => 'suffix',
            'rate' => 'exchange_rate',
            'default' => 'is_default',
        ],
        // shd_promo_code → promo_codes
        'shd_promo_code' => [
            'id' => 'id',
            'code' => 'code',
            'type' => 'type',
            'value' => 'value',
            'maxuses' => 'max_uses',
            'uses' => 'used_count',
            'startdate' => 'start_date',
            'enddate' => 'end_date',
            'status' => 'status',
        ],
        // shd_news → news
        'shd_news' => [
            'id' => 'id',
            'title' => 'title',
            'content' => 'content',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
        ],
        // shd_downloads → downloads
        'shd_downloads' => [
            'id' => 'id',
            'title' => 'title',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
        ],
        // shd_downloadcats → download_categories
        'shd_downloadcats' => [
            'id' => 'id',
            'name' => 'name',
            'sortorder' => 'sort_order',
        ],
        // shd_knowledge_base → knowledge_articles
        'shd_knowledge_base' => [
            'id' => 'id',
            'title' => 'title',
            'article' => 'content',          // zjmf article = 我们的 content
            'views' => 'view_count',
            'create_time' => 'created_at',
            'update_time' => 'updated_at',
        ],
        // shd_knowledge_categories → knowledge_categories
        'shd_knowledge_categories' => [
            'id' => 'id',
            'name' => 'name',
            'sortorder' => 'sort_order',
        ],
        // shd_servers → servers (供应商服务器)
        'shd_servers' => [
            'id' => 'id',
            'name' => 'name',
            'hostname' => 'hostname',
            'ip' => 'ip',
            'type' => 'type',
            'username' => 'username',
            'password' => 'password_encrypted',
            'status' => 'status',
            'maxaccounts' => 'max_accounts',
            'activeaccounts' => 'active_accounts',
            'port' => 'port',
            'secure' => 'secure',
        ],
        // shd_contacts → (子账户/联系人，我们暂无此表，保留映射)
        'shd_contacts' => [
            'id' => 'id',
            'uid' => 'user_id',
            'username' => 'username',
            'email' => 'email',
            'phonenumber' => 'phone',
            'status' => 'status',
        ],
    ];

    // ========== 内容格式映射（zjmf值 => 我们的值）==========
    private static $valueMap = [
        // 业务状态映射（shd_host.domainstatus → services.status）
        'domainstatus' => [
            'Pending' => 'pending',
            'Active' => 'active',
            'Suspended' => 'suspended',
            'Cancelled' => 'cancelled',
            'Fraud' => 'fraud',
            'Completed' => 'completed',
            'Deleted' => 'terminated',
        ],
        // 账单状态映射（shd_invoices.status → invoices.status）
        'invoice_status' => [
            'Paid' => 'paid',
            'Unpaid' => 'unpaid',
            'Draft' => 'draft',
            'Overdue' => 'overdue',
            'Cancelled' => 'cancelled',
            'Refunded' => 'refunded',
            'Collections' => 'collections',
        ],
        // 订单状态映射（shd_orders.status → orders.status）
        'order_status' => [
            'Pending' => 'pending',
            'Active' => 'active',
            'Completed' => 'completed',
            'Suspend' => 'suspended',
            'Terminated' => 'terminated',
            'Cancelled' => 'cancelled',
            'Fraud' => 'fraud',
        ],
        // 用户状态映射（shd_clients.status → users.status）
        'client_status' => [
            '1' => 'active',     // zjmf: 1=激活
            '0' => 'inactive',   // zjmf: 0=未激活
            '2' => 'closed',     // zjmf: 2=关闭
        ],
        // 工单状态映射（shd_ticket.status → tickets.status）
        'ticket_status' => [
            '1' => 'open',
            '2' => 'answered',
            '3' => 'customer_reply',
            '4' => 'closed',
        ],
        // 工单优先级映射
        'priority' => [
            'Low' => 'low',
            'Medium' => 'medium',
            'High' => 'high',
        ],
        // 插件状态映射（shd_plugin.status → plugins.status）
        'plugin_status' => [
            '1' => 'active',     // zjmf: 1=启用
            '0' => 'inactive',   // zjmf: 0=禁用
        ],
        // 通用布尔值映射（zjmf tinyint → 我们的bool）
        'is_default' => [
            '1' => true,
            '0' => false,
        ],
        'hidden' => [
            '1' => true,
            '0' => false,
        ],
        'adminonly' => [
            '1' => true,
            '0' => false,
        ],
        'required' => [
            '1' => true,
            '0' => false,
        ],
    ];

    // ========== 时间字段列表（需要int↔datetime转换）==========
    private static $timeFields = [
        'create_time', 'update_time', 'lastlogin', 'lastloginip',
        'regdate', 'nextduedate', 'nextinvoicedate', 'due_time',
        'paid_time', 'pay_time', 'last_reply_time', 'datecreated',
        'pwresetexpiry', 'delete_time', 'startdate', 'enddate',
    ];

    // ========== 反向映射（我们的值 => zjmf值）==========
    private static $reverseValueMap = null;

    /**
     * 设置PDO连接
     */
    public static function setConnection(PDO $pdo)
    {
        self::$pdo = $pdo;
    }

    /**
     * 获取PDO连接
     */
    public static function getConnection()
    {
        return self::$pdo;
    }

    /**
     * 获取映射后的表名
     */
    public static function getMappedTableName($name)
    {
        return self::$tableMap[$name] ?? $name;
    }

    /**
     * 获取反向映射后的表名（我们的表名 => zjmf表名）
     */
    public static function getReverseTableName($name)
    {
        $reverse = array_flip(self::$tableMap);
        return $reverse[$name] ?? $name;
    }

    /**
     * 映射字段名（zjmf字段 => 我们的字段）
     */
    public static function mapFieldName($table, $field)
    {
        $zjmfTable = self::getReverseTableName($table);
        $map = self::$fieldMap[$zjmfTable] ?? [];
        return $map[$field] ?? $field;
    }

    /**
     * 反向映射字段名（我们的字段 => zjmf字段）
     */
    public static function reverseFieldName($table, $field)
    {
        $zjmfTable = self::getReverseTableName($table);
        $map = self::$fieldMap[$zjmfTable] ?? [];
        $reverse = array_flip($map);
        return $reverse[$field] ?? $field;
    }

    /**
     * 映射字段数据（整行数据）
     */
    public static function mapFields($table, $data)
    {
        $zjmfTable = self::getReverseTableName($table);
        $map = self::$fieldMap[$zjmfTable] ?? [];
        $result = [];
        foreach ($data as $key => $value) {
            $mappedKey = $map[$key] ?? $key;
            $result[$mappedKey] = $value;
        }
        return $result;
    }

    /**
     * 反向映射字段数据
     */
    public static function reverseMapFields($table, $data)
    {
        $zjmfTable = self::getReverseTableName($table);
        $map = self::$fieldMap[$zjmfTable] ?? [];
        $reverse = array_flip($map);
        $result = [];
        foreach ($data as $key => $value) {
            $mappedKey = $reverse[$key] ?? $key;
            $result[$mappedKey] = $value;
        }
        return $result;
    }

    /**
     * 映射值（zjmf值 => 我们的值）
     * 自动处理：状态值映射、时间戳转换、布尔值转换
     */
    public static function mapValue($field, $value)
    {
        // 1. 时间字段：int时间戳 → datetime字符串
        if (in_array($field, self::$timeFields) && is_numeric($value) && $value > 0) {
            return date('Y-m-d H:i:s', (int)$value);
        }

        // 2. 状态值映射
        $map = self::$valueMap[$field] ?? [];
        if (isset($map[$value])) {
            return $map[$value];
        }

        // 3. 通用status字段映射（如果字段名是status但没有专门的映射）
        if ($field === 'status' && isset(self::$valueMap['domainstatus'][$value])) {
            return self::$valueMap['domainstatus'][$value];
        }

        return $value;
    }

    /**
     * 反向映射值（我们的值 => zjmf值）
     * 自动处理：状态值反向映射、datetime→时间戳
     */
    public static function reverseMapValue($field, $value)
    {
        // 1. 时间字段：datetime字符串 → int时间戳
        if (in_array($field, self::$timeFields) && is_string($value)) {
            $ts = strtotime($value);
            if ($ts !== false) return $ts;
        }

        // 2. 状态值反向映射
        if (self::$reverseValueMap === null) {
            self::$reverseValueMap = [];
            foreach (self::$valueMap as $f => $map) {
                self::$reverseValueMap[$f] = array_flip($map);
            }
        }

        $map = self::$reverseValueMap[$field] ?? [];
        return $map[$value] ?? $value;
    }

    /**
     * 查询构建器（兼容db()函数）
     */
    public static function table($name)
    {
        $mapped = self::getMappedTableName($name);
        return new QueryBuilder(self::$pdo, $mapped, $name);
    }

    /**
     * 执行原生SQL查询
     */
    public static function query($sql, $bindings = [])
    {
        $stmt = self::$pdo->prepare($sql);
        $stmt->execute($bindings);
        return $stmt->fetchAll(PDO::FETCH_OBJ);
    }

    /**
     * 执行原生SQL语句
     */
    public static function execute($sql, $bindings = [])
    {
        $stmt = self::$pdo->prepare($sql);
        return $stmt->execute($bindings);
    }
}

/**
 * PDO查询构建器 - 带完整字段/值映射
 */
class QueryBuilder
{
    private $pdo;
    private $realTable;   // 实际表名（映射后）
    private $zjmfTable;   // zjmf表名（原始）
    private $wheres = [];
    private $bindings = [];
    private $orders = [];
    private $limitVal = null;
    private $offsetVal = null;
    private $selects = ['*'];

    public function __construct(PDO $pdo, string $realTable, string $zjmfTable = '')
    {
        $this->pdo = $pdo;
        $this->realTable = $realTable;
        $this->zjmfTable = $zjmfTable ?: $realTable;
    }

    /**
     * 映射WHERE条件中的字段名和值
     */
    private function mapWhereField($column)
    {
        return DatabaseCompat::mapFieldName($this->zjmfTable, $column);
    }

    private function mapWhereValue($column, $value)
    {
        return DatabaseCompat::mapValue($column, $value);
    }

    public function select($columns)
    {
        if (is_string($columns)) {
            $columns = array_map('trim', explode(',', $columns));
        }
        $this->selects = $columns;
        return $this;
    }

    public function where($column, $operator, $value = null)
    {
        if ($value === null) {
            $value = $operator;
            $operator = '=';
        }
        $mappedColumn = $this->mapWhereField($column);
        $mappedValue = $this->mapWhereValue($column, $value);
        $this->wheres[] = "{$mappedColumn} {$operator} ?";
        $this->bindings[] = $mappedValue;
        return $this;
    }

    public function whereIn($column, array $values)
    {
        $mappedColumn = $this->mapWhereField($column);
        $placeholders = implode(',', array_fill(0, count($values), '?'));
        $this->wheres[] = "{$mappedColumn} IN ({$placeholders})";
        $mappedValues = array_map(fn($v) => $this->mapWhereValue($column, $v), $values);
        $this->bindings = array_merge($this->bindings, $mappedValues);
        return $this;
    }

    public function whereNull($column)
    {
        $mappedColumn = $this->mapWhereField($column);
        $this->wheres[] = "{$mappedColumn} IS NULL";
        return $this;
    }

    public function whereNotNull($column)
    {
        $mappedColumn = $this->mapWhereField($column);
        $this->wheres[] = "{$mappedColumn} IS NOT NULL";
        return $this;
    }

    public function orderBy($column, $direction = 'ASC')
    {
        $mappedColumn = $this->mapWhereField($column);
        $this->orders[] = "{$mappedColumn} {$direction}";
        return $this;
    }

    public function limit($limit)
    {
        $this->limitVal = $limit;
        return $this;
    }

    public function offset($offset)
    {
        $this->offsetVal = $offset;
        return $this;
    }

    /**
     * 构建并执行SELECT查询，返回结果时反向映射字段名+值格式
     */
    public function get()
    {
        $columns = implode(', ', $this->mapSelectColumns());
        $sql = "SELECT {$columns} FROM {$this->realTable}";
        if (!empty($this->wheres)) {
            $sql .= ' WHERE ' . implode(' AND ', $this->wheres);
        }
        if (!empty($this->orders)) {
            $sql .= ' ORDER BY ' . implode(', ', $this->orders);
        }
        if ($this->limitVal !== null) {
            $sql .= ' LIMIT ' . $this->limitVal;
        }
        if ($this->offsetVal !== null) {
            $sql .= ' OFFSET ' . $this->offsetVal;
        }

        $stmt = $this->pdo->prepare($sql);
        $stmt->execute($this->bindings);
        $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

        // 反向映射字段名 + 值格式（让zjmf插件看到它期望的格式）
        return array_map(function ($row) {
            $mapped = DatabaseCompat::reverseMapFields($this->zjmfTable, $row);
            // 反向映射值（datetime→时间戳等）
            foreach ($mapped as $key => $value) {
                $mapped[$key] = DatabaseCompat::reverseMapValue($key, $value);
            }
            return $mapped;
        }, $rows);
    }

    public function first()
    {
        $this->limitVal = 1;
        $results = $this->get();
        return $results[0] ?? null;
    }

    public function insert(array $data)
    {
        // 映射字段名
        $mappedData = DatabaseCompat::mapFields($this->zjmfTable, $data);

        $columns = implode(', ', array_keys($mappedData));
        $placeholders = implode(', ', array_fill(0, count($mappedData), '?'));
        $sql = "INSERT INTO {$this->realTable} ({$columns}) VALUES ({$placeholders})";
        $stmt = $this->pdo->prepare($sql);
        $stmt->execute(array_values($mappedData));
        return $this->pdo->lastInsertId();
    }

    public function update(array $data)
    {
        // 映射字段名
        $mappedData = DatabaseCompat::mapFields($this->zjmfTable, $data);

        $sets = [];
        $values = [];
        foreach ($mappedData as $key => $value) {
            $sets[] = "{$key} = ?";
            $values[] = $value;
        }
        $sql = "UPDATE {$this->realTable} SET " . implode(', ', $sets);
        if (!empty($this->wheres)) {
            $sql .= ' WHERE ' . implode(' AND ', $this->wheres);
        }
        $stmt = $this->pdo->prepare($sql);
        $stmt->execute(array_merge($values, $this->bindings));
        return $stmt->rowCount();
    }

    public function delete()
    {
        $sql = "DELETE FROM {$this->realTable}";
        if (!empty($this->wheres)) {
            $sql .= ' WHERE ' . implode(' AND ', $this->wheres);
        }
        $stmt = $this->pdo->prepare($sql);
        $stmt->execute($this->bindings);
        return $stmt->rowCount();
    }

    public function count()
    {
        $sql = "SELECT COUNT(*) as cnt FROM {$this->realTable}";
        if (!empty($this->wheres)) {
            $sql .= ' WHERE ' . implode(' AND ', $this->wheres);
        }
        $stmt = $this->pdo->prepare($sql);
        $stmt->execute($this->bindings);
        $result = $stmt->fetch(PDO::FETCH_OBJ);
        return (int)($result->cnt ?? 0);
    }

    public function exists()
    {
        return $this->count() > 0;
    }

    /**
     * 映射SELECT字段
     */
    private function mapSelectColumns()
    {
        if ($this->selects === ['*']) return ['*'];
        return array_map(function ($col) {
            $mapped = $this->mapWhereField($col);
            return "{$mapped} AS \"{$col}\"";
        }, $this->selects);
    }
}
