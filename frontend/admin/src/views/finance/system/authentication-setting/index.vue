<template>
  <div class="authentication-setting-page">
    <!-- 说明文字 -->
    <el-alert type="info" :closable="false" class="page-desc">
      <template #title>
        <span>选择开启实名认证方式，用户提交实名认证后可进行审核。</span>
      </template>
    </el-alert>

    <!-- 认证方式列表 -->
    <el-table 
      :data="plugins" 
      v-loading="loading" 
      border 
      stripe
      empty-text="暂无数据"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" />
      
      <el-table-column prop="title" label="名称" min-width="150">
        <template #default="{ row }">
          <span class="plugin-name">{{ row.title }}</span>
        </template>
      </el-table-column>
      
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      
      <el-table-column prop="version" label="版本" width="100" />
      
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch 
            v-model="row.is_enabled" 
            @change="handleToggle(row)"
          />
        </template>
      </el-table-column>
      
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button 
            size="small" 
            :type="row.is_enabled ? 'warning' : 'success'" 
            link
            @click="handleToggle(row)"
          >
            {{ row.is_enabled ? '禁用' : '启用' }}
          </el-button>
          <el-button v-if="!row.is_system" size="small" type="danger" link @click="handleUninstall(row)">卸载</el-button>
          <el-button size="small" type="primary" link @click="handleConfig(row)">配置</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 配置弹窗 -->
    <el-dialog 
      v-model="configVisible" 
      :title="`${currentPlugin?.title} - 配置`" 
      width="600px"
      destroy-on-close
    >
      <el-form :model="configForm" label-position="top">
        <template v-if="currentPlugin?.name === 'AlipayCertify'">
          <el-form-item label="AppID">
            <el-input v-model="configForm.appid" placeholder="请输入AppID" />
          </el-form-item>
          <el-form-item label="私钥">
            <el-input v-model="configForm.private_key" type="textarea" :rows="3" placeholder="请输入私钥" />
          </el-form-item>
          <el-form-item label="支付宝公钥">
            <el-input v-model="configForm.alipay_public_key" type="textarea" :rows="3" placeholder="请输入支付宝公钥" />
          </el-form-item>
        </template>

        <template v-if="currentPlugin?.name === 'ManualVerify'">
          <el-form-item label="审核说明">
            <el-input v-model="configForm.description" type="textarea" :rows="3" placeholder="请输入审核说明" />
          </el-form-item>
          <el-form-item label="允许证件类型">
            <el-checkbox-group v-model="configForm.id_types">
              <el-checkbox label="id_card">身份证</el-checkbox>
              <el-checkbox label="passport">护照</el-checkbox>
              <el-checkbox label="driver_license">驾驶证</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </template>

        <el-form-item label="备注">
          <el-input v-model="configForm.remark" type="textarea" :rows="2" placeholder="备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingConfig" @click="handleSaveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const plugins = ref<any[]>([])
const configVisible = ref(false)
const currentPlugin = ref<any>(null)
const configForm = ref<any>({
  appid: '',
  private_key: '',
  alipay_public_key: '',
  description: '',
  id_types: [],
  remark: ''
})
const savingConfig = ref(false)

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/certification-plugins' })
    plugins.value = data || []
  } catch (error) {
    console.error('fetch certification plugins failed:', error)
  } finally {
    loading.value = false
  }
}

const handleToggle = async (item: any) => {
  const original = item.is_enabled
  try {
    await request.post({ 
      url: `/api/admin/certification-plugins/${item.id}/toggle`, 
      data: { enabled: item.is_enabled } 
    })
    ElMessage.success(item.is_enabled ? '已启用' : '已禁用')
  } catch (error) {
    item.is_enabled = original
    ElMessage.error('操作失败')
  }
}

const handleUninstall = async (item: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要卸载 "${item.title}" 吗？`, 
      '确认卸载', 
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await request.post({ url: `/api/admin/certification-plugins/${item.id}/uninstall` })
    ElMessage.success('已卸载')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('卸载失败')
    }
  }
}

const handleConfig = async (item: any) => {
  currentPlugin.value = item
  try {
    const data = await request.get({ url: `/api/admin/certification-plugins/${item.id}/config` })
    configForm.value = data || { appid: '', private_key: '', alipay_public_key: '', description: '', id_types: [], remark: '' }
  } catch (error) {
    configForm.value = { appid: '', private_key: '', alipay_public_key: '', description: '', id_types: [], remark: '' }
  }
  configVisible.value = true
}

const handleSaveConfig = async () => {
  savingConfig.value = true
  try {
    await request.post({ 
      url: `/api/admin/certification-plugins/${currentPlugin.value.id}/config`, 
      data: configForm.value 
    })
    ElMessage.success('保存成功')
    configVisible.value = false
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
.authentication-setting-page {
  padding: 20px;
}

.page-desc {
  margin-bottom: 16px;
}

.plugin-name {
  font-weight: 500;
}
</style>
