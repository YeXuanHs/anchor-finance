<template>
  <div class="payment-interface-page">
    <!-- 说明文字 -->
    <el-alert type="info" :closable="false" class="page-desc">
      <template #title>
        <span>选择支付网关点击安装，填写接口相关配置信息后，用户在前台可选择支付方式进行结算付款。</span>
        <el-link type="primary" href="https://www.idcsmart.com/wiki_list/539.html" target="_blank" :underline="false">帮助文档</el-link>
      </template>
    </el-alert>

    <div class="action-bar">
      <el-button type="primary" @click="handleGetMore">获取更多支付接口</el-button>
    </div>

    <!-- 支付接口表格 -->
    <el-table 
      :data="paymentList" 
      v-loading="loading" 
      border 
      stripe
      empty-text="暂无数据"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" sortable />
      
      <el-table-column prop="name" label="插件名称" min-width="150">
        <template #default="{ row }">
          <span class="plugin-name">{{ row.name }}</span>
        </template>
      </el-table-column>
      
      <el-table-column prop="identifier" label="标识" min-width="120" />
      
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      
      <el-table-column prop="author" label="作者" width="120" />
      
      <el-table-column prop="version" label="版本" width="100" />
      
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch 
            v-if="row.installed"
            v-model="row.enabled" 
            @change="handleToggle(row)"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>
      
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <template v-if="row.installed">
            <el-button size="small" type="danger" link @click="handleUninstall(row)">卸载</el-button>
            <el-button 
              size="small" 
              :type="row.enabled ? 'warning' : 'success'" 
              link
              @click="handleToggle(row)"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="primary" link @click="handleConfigure(row)">配置</el-button>
            <el-button v-if="row.allow_copy" size="small" type="info" link @click="handleCopy(row)">复制</el-button>
          </template>
          <template v-else>
            <el-button size="small" type="primary" @click="handleInstall(row)">安装</el-button>
          </template>
        </template>
      </el-table-column>
    </el-table>

    <!-- 配置弹窗 -->
    <el-dialog 
      v-model="configDialogVisible" 
      :title="`${currentItem?.name} - 配置`" 
      width="600px"
      destroy-on-close
    >
      <el-form 
        v-if="currentItem?.config_fields" 
        ref="configFormRef" 
        :model="configFormData" 
        label-position="top"
      >
        <el-form-item 
          v-for="field in currentItem.config_fields" 
          :key="field.key" 
          :label="field.label" 
          :prop="field.key"
        >
          <el-input 
            v-if="field.type === 'text' || field.type === 'password'" 
            v-model="configFormData[field.key]" 
            :type="field.type" 
            :placeholder="field.placeholder" 
            show-password 
          />
          <el-switch 
            v-else-if="field.type === 'switch'" 
            v-model="configFormData[field.key]" 
          />
          <el-select 
            v-else-if="field.type === 'select'" 
            v-model="configFormData[field.key]" 
            style="width: 100%"
          >
            <el-option 
              v-for="opt in field.options" 
              :key="opt.value" 
              :label="opt.label" 
              :value="opt.value" 
            />
          </el-select>
          <el-input
            v-else-if="field.type === 'textarea'"
            v-model="configFormData[field.key]"
            type="textarea"
            :rows="3"
            :placeholder="field.placeholder"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingConfig" @click="handleSaveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const paymentList = ref<any[]>([])
const configDialogVisible = ref(false)
const currentItem = ref<any>(null)
const configFormData = ref<Record<string, any>>({})
const configFormRef = ref<FormInstance>()
const savingConfig = ref(false)

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/payment-gateways' })
    paymentList.value = data || []
  } catch (error) {
    console.error('fetch payment gateways failed:', error)
  } finally {
    loading.value = false
  }
}

const handleGetMore = () => {
  window.open('https://www.idcsmart.com/wiki_list/539.html', '_blank')
}

const handleToggle = async (item: any) => {
  const originalEnabled = item.enabled
  try {
    await request.post({ 
      url: `/api/admin/payment-gateways/${item.id}/toggle`, 
      data: { enabled: item.enabled } 
    })
    ElMessage.success(item.enabled ? '已启用' : '已禁用')
  } catch (error) {
    item.enabled = originalEnabled
    ElMessage.error('操作失败')
  }
}

const handleUninstall = async (item: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要卸载 "${item.name}" 吗？`, 
      '确认卸载', 
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await request.post({ url: `/api/admin/payment-gateways/${item.id}/uninstall` })
    ElMessage.success('已卸载')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('卸载失败')
    }
  }
}

const handleInstall = async (item: any) => {
  try {
    await request.post({ url: `/api/admin/payment-gateways/${item.id}/install` })
    ElMessage.success('已安装')
    fetchList()
  } catch (error) {
    ElMessage.error('安装失败')
  }
}

const handleCopy = async (item: any) => {
  try {
    await request.post({ url: `/api/admin/payment-gateways/${item.id}/copy` })
    ElMessage.success('已复制')
    fetchList()
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const handleConfigure = async (item: any) => {
  currentItem.value = item
  try {
    const data = await request.get({ url: `/api/admin/payment-gateways/${item.id}/config` })
    configFormData.value = data || {}
  } catch (error) {
    configFormData.value = {}
  }
  configDialogVisible.value = true
}

const handleSaveConfig = async () => {
  if (!configFormRef.value) return
  try {
    await configFormRef.value.validate()
    savingConfig.value = true
    await request.post({ 
      url: `/api/admin/payment-gateways/${currentItem.value.id}/config`, 
      data: configFormData.value 
    })
    ElMessage.success('保存成功')
    configDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    savingConfig.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.payment-interface-page {
  padding: 20px;
}

.page-desc {
  margin-bottom: 16px;
  :deep(.el-alert__title) {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.action-bar {
  margin-bottom: 16px;
}

.plugin-name {
  font-weight: 500;
}
</style>
