<template>
  <div class="invoice-address-page">
    <div class="page-header">
      <h1 class="page-title">发票地址管理</h1>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增地址
      </el-button>
    </div>

    <!-- 搜索筛选 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="收件人">
          <el-input v-model="filterForm.receiver" placeholder="搜索收件人" clearable style="width: 180px;" />
        </el-form-item>
        <el-form-item label="地区">
          <el-input v-model="filterForm.region" placeholder="搜索地区" clearable style="width: 180px;" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleResetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 地址列表 -->
    <div class="address-grid" v-loading="loading">
      <el-card
        v-for="address in filteredAddresses"
        :key="address.id"
        shadow="never"
        class="address-card"
        :class="{ 'is-default': address.isDefault }"
      >
        <div class="address-header">
          <div class="address-tags">
            <el-tag v-if="address.isDefault" type="primary" size="small" effect="dark">默认</el-tag>
            <el-tag v-if="address.tag" size="small" effect="plain">{{ address.tag }}</el-tag>
          </div>
          <el-dropdown @command="(cmd: string) => handleCommand(cmd, address)">
            <el-button type="primary" size="small" link>
              <el-icon><MoreFilled /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit">编辑</el-dropdown-item>
                <el-dropdown-item command="default" :disabled="address.isDefault">设为默认</el-dropdown-item>
                <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="address-info">
          <p class="address-receiver">
            <span class="receiver-name">{{ address.receiver }}</span>
            <span class="receiver-phone">{{ address.phone }}</span>
          </p>
          <p class="address-detail">{{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}</p>
        </div>
      </el-card>

      <el-empty v-if="filteredAddresses.length === 0" description="暂无发票地址" style="grid-column: 1 / -1;" />
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingAddress ? '编辑地址' : '新增地址'"
      width="520px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="收件人" prop="receiver">
          <el-input v-model="form.receiver" placeholder="请输入收件人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="所在地区" prop="region">
          <el-cascader
            v-model="form.region"
            :options="regionOptions"
            placeholder="请选择省/市/区"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="详细地址" prop="detail">
          <el-input v-model="form.detail" type="textarea" :rows="2" placeholder="请输入详细地址" />
        </el-form-item>
        <el-form-item label="地址标签">
          <el-input v-model="form.tag" placeholder="如：公司、家（选填）" maxlength="10" />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.isDefault" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, MoreFilled } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

interface InvoiceAddress {
  id: number
  receiver: string
  phone: string
  province: string
  city: string
  district: string
  detail: string
  tag?: string
  isDefault: boolean
}

const dialogVisible = ref(false)
const editingAddress = ref<InvoiceAddress | null>(null)
const formRef = ref<FormInstance>()

const filterForm = reactive({
  receiver: '',
  region: ''
})

const addresses = ref<InvoiceAddress[]>([])

const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/contacts')
    if (data?.data) {
      addresses.value = data.data.list || data.data || []
    }
  } catch (e) {
    console.error('Failed to fetch addresses:', e)
  } finally {
    loading.value = false
  }
})

const form = reactive({
  receiver: '',
  phone: '',
  region: [] as string[],
  detail: '',
  tag: '',
  isDefault: false
})

const rules: FormRules = {
  receiver: [{ required: true, message: '请输入收件人姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  region: [{ required: true, message: '请选择所在地区', trigger: 'change', type: 'array' }],
  detail: [{ required: true, message: '请输入详细地址', trigger: 'blur' }]
}

const regionOptions = [
  {
    value: '广东省',
    label: '广东省',
    children: [
      { value: '深圳市', label: '深圳市', children: [{ value: '南山区', label: '南山区' }, { value: '福田区', label: '福田区' }, { value: '龙岗区', label: '龙岗区' }] },
      { value: '广州市', label: '广州市', children: [{ value: '天河区', label: '天河区' }, { value: '越秀区', label: '越秀区' }] }
    ]
  },
  {
    value: '北京市',
    label: '北京市',
    children: [
      { value: '北京市', label: '北京市', children: [{ value: '海淀区', label: '海淀区' }, { value: '朝阳区', label: '朝阳区' }] }
    ]
  },
  {
    value: '上海市',
    label: '上海市',
    children: [
      { value: '上海市', label: '上海市', children: [{ value: '浦东新区', label: '浦东新区' }, { value: '徐汇区', label: '徐汇区' }] }
    ]
  }
]

const filteredAddresses = computed(() => {
  let result = addresses.value
  if (filterForm.receiver) {
    result = result.filter(a => a.receiver.includes(filterForm.receiver))
  }
  if (filterForm.region) {
    result = result.filter(a =>
      a.province.includes(filterForm.region) ||
      a.city.includes(filterForm.region) ||
      a.district.includes(filterForm.region)
    )
  }
  return result
})

function handleSearch() {
  ElMessage.success('查询完成')
}

function handleResetFilter() {
  filterForm.receiver = ''
  filterForm.region = ''
}

function handleAdd() {
  editingAddress.value = null
  form.receiver = ''
  form.phone = ''
  form.region = []
  form.detail = ''
  form.tag = ''
  form.isDefault = false
  dialogVisible.value = true
}

function handleEdit(address: InvoiceAddress) {
  editingAddress.value = address
  form.receiver = address.receiver
  form.phone = address.phone
  form.region = [address.province, address.city, address.district]
  form.detail = address.detail
  form.tag = address.tag || ''
  form.isDefault = address.isDefault
  dialogVisible.value = true
}

function handleCommand(cmd: string, address: InvoiceAddress) {
  switch (cmd) {
    case 'edit':
      handleEdit(address)
      break
    case 'default':
      addresses.value.forEach(a => a.isDefault = false)
      address.isDefault = true
      ElMessage.success('已设为默认地址')
      break
    case 'delete':
      ElMessageBox.confirm('确定删除该地址吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        addresses.value = addresses.value.filter(a => a.id !== address.id)
        ElMessage.success('已删除')
      }).catch(() => {})
      break
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate((valid) => {
    if (valid) {
      if (form.isDefault) {
        addresses.value.forEach(a => a.isDefault = false)
      }

      if (editingAddress.value) {
        Object.assign(editingAddress.value, {
          receiver: form.receiver,
          phone: form.phone,
          province: form.region[0],
          city: form.region[1],
          district: form.region[2],
          detail: form.detail,
          tag: form.tag,
          isDefault: form.isDefault
        })
        ElMessage.success('地址已更新')
      } else {
        addresses.value.push({
          id: Date.now(),
          receiver: form.receiver,
          phone: form.phone,
          province: form.region[0],
          city: form.region[1],
          district: form.region[2],
          detail: form.detail,
          tag: form.tag,
          isDefault: form.isDefault
        })
        ElMessage.success('地址已添加')
      }
      dialogVisible.value = false
    }
  })
}
</script>

<style scoped>
.invoice-address-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.filter-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.filter-card :deep(.el-card__body) {
  padding: 16px 20px 0;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 0;
}

.address-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.address-card {
  border-radius: 12px;
  border: 2px solid #e8ecf1;
  transition: all 0.2s;
  cursor: default;
}

.address-card:hover {
  border-color: #c0c4cc;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.address-card.is-default {
  border-color: #0056FF;
}

.address-card :deep(.el-card__body) {
  padding: 18px 20px;
}

.address-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.address-tags {
  display: flex;
  gap: 8px;
}

.address-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.address-receiver {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
}

.receiver-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.receiver-phone {
  font-size: 14px;
  color: #606266;
}

.address-detail {
  font-size: 14px;
  color: #909399;
  margin: 0;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .page-header { flex-direction: column; gap: 12px; align-items: flex-start; }
  .address-grid { grid-template-columns: 1fr; }
}
</style>
