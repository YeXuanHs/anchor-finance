<template>
  <div class="page-container">
    <art-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('systemNavigation.title') }}</span>
        </div>
      </template>

      <p class="description">{{ $t('systemNavigation.description') }}</p>

      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('systemNavigation.positionTab')" name="positions">
          <el-form :model="positionForm" label-width="120px" style="max-width: 600px; margin-top: 20px;">
            <el-form-item :label="$t('systemNavigation.memberNav')">
              <el-select v-model="positionForm.member_nav" :placeholder="$t('common.select')" style="width: 100%;">
                <el-option :label="$t('systemNavigation.defaultMenu')" value="default" />
                <el-option :label="$t('systemNavigation.customMenu')" value="custom" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('systemNavigation.topNav')">
              <el-select v-model="positionForm.top_nav" :placeholder="$t('common.select')" style="width: 100%;">
                <el-option :label="$t('systemNavigation.siteHeaderNav')" value="header" />
                <el-option :label="$t('systemNavigation.customNav')" value="custom" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('systemNavigation.bottomNav')">
              <el-select v-model="positionForm.bottom_nav" :placeholder="$t('common.select')" style="width: 100%;">
                <el-option :label="$t('systemNavigation.siteFooterNav')" value="footer" />
                <el-option :label="$t('systemNavigation.customNav')" value="custom" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSave" :loading="saving">{{ $t('systemNavigation.saveChanges') }}</el-button>
              <el-button @click="handleReset">{{ $t('systemNavigation.cancelChanges') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('systemNavigation.allNavTab')" name="all">
          <div class="nav-list-header">
            <el-button type="primary" @click="handleAddNav">
              <el-icon><Plus /></el-icon> {{ $t('systemNavigation.addNav') }}
            </el-button>
          </div>

          <el-table v-loading="navLoading" :data="navList" border stripe>
            <el-table-column prop="id" label="ID" width="60" align="center" />
            <el-table-column prop="name" :label="$t('systemNavigation.navName')" min-width="150" />
            <el-table-column prop="type" :label="$t('common.type')" width="120" />
            <el-table-column prop="items_count" :label="$t('systemNavigation.itemCount')" width="80" align="center" />
            <el-table-column :label="$t('common.action')" width="150" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditNav(row)">{{ $t('common.edit') }}</el-button>
                <el-button type="danger" link size="small" @click="handleDeleteNav(row)">{{ $t('common.delete') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </art-card>

    <el-dialog v-model="navDialogVisible" :title="isEditNav ? $t('systemNavigation.editNav') : $t('systemNavigation.addNav')" width="500px" @close="resetNavForm">
      <el-form ref="navFormRef" :model="navFormData" :rules="navRules" label-width="80px">
        <el-form-item :label="$t('common.name')" prop="name">
          <el-input v-model="navFormData.name" :placeholder="$t('systemNavigation.enterNavName')" />
        </el-form-item>
        <el-form-item :label="$t('common.type')" prop="type">
          <el-select v-model="navFormData.type" :placeholder="$t('systemNavigation.selectType')" style="width: 100%">
            <el-option :label="$t('systemNavigation.headerNav')" value="header" />
            <el-option :label="$t('systemNavigation.footerNav')" value="footer" />
            <el-option :label="$t('systemNavigation.memberCenter')" value="member" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="navDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="navSubmitting" @click="handleNavSubmit">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const activeTab = ref('positions')
const saving = ref(false)

const positionForm = ref({ member_nav: 'default', top_nav: 'header', bottom_nav: 'footer' })

// 导航列表相关
const navLoading = ref(false)
const navList = ref([])
const navDialogVisible = ref(false)
const isEditNav = ref(false)
const navSubmitting = ref(false)
const editingNavId = ref<number | null>(null)
const navFormRef = ref<FormInstance>()

const navFormData = reactive({ name: '', type: 'header' })
const navRules: FormRules = {
  name: [{ required: true, message: () => $t('systemNavigation.enterNavName'), trigger: 'blur' }],
  type: [{ required: true, message: () => $t('systemNavigation.selectType'), trigger: 'change' }]
}

const fetchPositions = async () => {
  try {
    const res = await request.get({ url: '/api/admin/nav-positions' })
    if (res) positionForm.value = { ...positionForm.value, ...res }
  } catch (error) { console.error('获取导航位置失败:', error) }
}

const fetchNavList = async () => {
  navLoading.value = true
  try {
    const data = await request.get({ url: '/api/admin/web-navs' })
    navList.value = data?.list || data || []
  } catch (error) { console.error('获取导航列表失败:', error) } finally { navLoading.value = false }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.put({ url: '/api/admin/nav-positions', data: positionForm.value }); ElMessage.success($t('common.saveSuccess'))
  } catch (error) { console.error('保存失败:', error) } finally { saving.value = false }
}

const handleReset = () => { fetchPositions() }

const resetNavForm = () => { navFormData.name = ''; navFormData.type = 'header'; editingNavId.value = null; navFormRef.value?.resetFields() }
const handleAddNav = () => { isEditNav.value = false; resetNavForm(); navDialogVisible.value = true }
const handleEditNav = (row: any) => {
  isEditNav.value = true; editingNavId.value = row.id
  Object.assign(navFormData, { name: row.name, type: row.type || 'header' })
  navDialogVisible.value = true
}

const handleNavSubmit = async () => {
  if (!navFormRef.value) return
  try {
    await navFormRef.value.validate(); navSubmitting.value = true
    if (isEditNav.value && editingNavId.value) {
      await request.put({ url: `/api/admin/web-navs/${editingNavId.value}`, data: navFormData }); ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/web-navs', data: navFormData }); ElMessage.success($t('common.addSuccess'))
    }
    navDialogVisible.value = false; fetchNavList()
  } catch (error) { console.error('提交失败:', error) } finally { navSubmitting.value = false }
}

const handleDeleteNav = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('systemNavigation.confirmDeleteNav'), $t('common.tips'), { type: 'warning' })
    await request.del({ url: `/api/admin/web-navs/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchNavList()
  } catch (error) { if (error !== 'cancel') console.error('删除失败:', error) }
}

onMounted(() => { fetchPositions(); fetchNavList() })
</script>

<style scoped lang="scss">
.page-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.description { color: #666; margin-bottom: 20px; }
.nav-list-header { margin-bottom: 16px; }
</style>
