<?php

namespace App\ZjmfCompat;

use Illuminate\Support\Facades\DB;

/**
 * 数据库兼容层
 * 兼容zjmf的数据库查询语法
 */
class DatabaseCompat
{
    // 表名映射：zjmf表名 => AnchorFinance表名
    private static $tableMap = [
        'plugin' => 'plugins',
        'hook_plugin' => 'hook_bindings',
        'clients' => 'users',
        'host' => 'services',
        'invoices' => 'invoices',
        'orders' => 'orders',
        'ticket' => 'tickets',
        'servers' => 'servers',
        'product' => 'products',
        'configuration' => 'settings',
        'customfields' => 'custom_fields',
        'customfieldsvalues' => 'custom_field_values',
        'promo_code' => 'promo_codes',
        'ticket_department' => 'ticket_departments',
        'ticket_status' => 'ticket_statuses',
        'news' => 'news',
        'downloads' => 'downloads',
        'downloadcats' => 'download_categories',
        'knowledge_base' => 'knowledge_articles',
        'knowledge_categories' => 'knowledge_categories',
    ];

    // 字段名映射：zjmf字段 => AnchorFinance字段
    private static $fieldMap = [
        'clients' => [
            'id' => 'id',
            'username' => 'username',
            'email' => 'email',
            'password' => 'password_hash',
            'phone' => 'phone',
            'company' => 'company',
            'status' => 'status',
            'balance' => 'balance',
            'credit' => 'credit_limit',
            'groupid' => 'group_id',
            'level' => 'level_id',
        ],
        'host' => [
            'id' => 'id',
            'userid' => 'user_id',
            'productid' => 'product_id',
            'domain' => 'domain',
            'username' => 'username',
            'password' => 'password_encrypted',
            'status' => 'status',
            'billingcycle' => 'billing_cycle',
            'amount' => 'amount',
            'nextduedate' => 'next_due_date',
            'serverid' => 'upstream_id',
        ],
        'invoices' => [
            'id' => 'id',
            'userid' => 'user_id',
            'amount' => 'amount',
            'status' => 'status',
            'paymentmethod' => 'payment_method',
            'duedate' => 'due_date',
        ],
        'plugin' => [
            'id' => 'id',
            'name' => 'slug',
            'title' => 'name',
            'description' => 'description',
            'status' => 'status',
            'config' => 'config_json',
        ],
    ];

    // 内容格式映射：zjmf值 => AnchorFinance值
    private static $valueMap = [
        'status' => [
            'Active' => 'active',
            'Suspended' => 'suspended',
            'Terminated' => 'terminated',
            'Cancelled' => 'cancelled',
            'Pending' => 'pending',
        ],
    ];

    /**
     * 获取表名映射
     * @param string $name zjmf表名
     * @return string AnchorFinance表名
     */
    public static function getMappedTableName($name)
    {
        return self::$tableMap[$name] ?? $name;
    }

    /**
     * 获取查询构建器（兼容db()函数）
     * @param string $name 表名
     * @return \Illuminate\Database\Query\Builder
     */
    public static function table($name)
    {
        $mapped = self::getMappedTableName($name);
        return DB::table($mapped);
    }

    /**
     * 执行原生SQL查询
     * @param string $sql SQL语句
     * @param array $bindings 绑定参数
     * @return array
     */
    public static function query($sql, $bindings = [])
    {
        return DB::select($sql, $bindings);
    }

    /**
     * 执行原生SQL语句
     * @param string $sql SQL语句
     * @param array $bindings 绑定参数
     * @return int 影响行数
     */
    public static function execute($sql, $bindings = [])
    {
        return DB::statement($sql, $bindings);
    }

    /**
     * 映射字段名
     * @param string $table 表名
     * @param array $data 数据
     * @return array 映射后的数据
     */
    public static function mapFields($table, $data)
    {
        $map = self::$fieldMap[$table] ?? [];
        $result = [];
        foreach ($data as $key => $value) {
            $mappedKey = $map[$key] ?? $key;
            $result[$mappedKey] = $value;
        }
        return $result;
    }

    /**
     * 反向映射字段名（AnchorFinance字段 => zjmf字段）
     * @param string $table 表名
     * @param array $data 数据
     * @return array 反向映射后的数据
     */
    public static function reverseMapFields($table, $data)
    {
        $map = self::$fieldMap[$table] ?? [];
        $reverseMap = array_flip($map);
        $result = [];
        foreach ($data as $key => $value) {
            $mappedKey = $reverseMap[$key] ?? $key;
            $result[$mappedKey] = $value;
        }
        return $result;
    }

    /**
     * 映射值
     * @param string $field 字段名
     * @param mixed $value 值
     * @return mixed 映射后的值
     */
    public static function mapValue($field, $value)
    {
        $map = self::$valueMap[$field] ?? [];
        return $map[$value] ?? $value;
    }

    /**
     * 反向映射值
     * @param string $field 字段名
     * @param mixed $value 值
     * @return mixed 反向映射后的值
     */
    public static function reverseMapValue($field, $value)
    {
        $map = self::$valueMap[$field] ?? [];
        $reverseMap = array_flip($map);
        return $reverseMap[$value] ?? $value;
    }

    /**
     * 开始事务
     */
    public static function beginTransaction()
    {
        DB::beginTransaction();
    }

    /**
     * 提交事务
     */
    public static function commit()
    {
        DB::commit();
    }

    /**
     * 回滚事务
     */
    public static function rollback()
    {
        DB::rollBack();
    }
}
