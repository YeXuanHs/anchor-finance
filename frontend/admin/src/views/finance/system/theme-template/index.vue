<template>
  <div class="page-container">
    <el-tabs v-model="activeTab">
      <!-- 轮播图管理 -->
      <el-tab-pane :label="$t('systemThemeTemplate.bannerTab')" name="banners">
        <art-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('systemThemeTemplate.homeBanners') }}</span>
              <el-button type="primary" @click="handleAddBanner">
                <el-icon><Plus /></el-icon> {{ $t('systemThemeTemplate.addBanner') }}
              </el-button>
            </div>
          </template>

          <el-table :data="banners" v-loading="bannersLoading" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column :label="$t('systemThemeTemplate.image')" width="120">
              <template #default="{ row }">
                <el-image
                  :src="row.media_url"
                  :preview-src-list="[row.media_url]"
                  fit="cover"
                  style="width: 80px; height: 45px; border-radius: 4px;"
                />
              </template>
            </el-table-column>
            <el-table-column prop="title" :label="$t('systemThemeTemplate.titleCol')" min-width="150" />
            <el-table-column prop="description" :label="$t('common.description')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="link_url" :label="$t('systemThemeTemplate.link')" min-width="150" show-overflow-tooltip />
            <el-table-column :label="$t('common.status')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" :label="$t('common.sort')" width="80" />
            <el-table-column :label="$t('common.action')" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditBanner(row)">{{ $t('common.edit') }}</el-button>
                <el-button
                  :type="row.status === 1 ? 'warning' : 'success'"
                  link size="small"
                  @click="handleToggleBannerStatus(row)"
                >
                  {{ row.status === 1 ? $t('common.disable') : $t('common.enable') }}
                </el-button>
                <el-button type="danger" link size="small" @click="handleDeleteBanner(row)">{{ $t('common.delete') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </art-card>
      </el-tab-pane>

      <!-- 特色区域管理 -->
      <el-tab-pane :label="$t('systemThemeTemplate.featureTab')" name="features">
        <art-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('systemThemeTemplate.homeFeatures') }}</span>
              <el-button type="primary" @click="handleAddFeature">
                <el-icon><Plus /></el-icon> {{ $t('systemThemeTemplate.addFeature') }}
              </el-button>
            </div>
          </template>

          <el-table :data="features" v-loading="featuresLoading" border stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column :label="$t('systemThemeTemplate.icon')" width="80">
              <template #default="{ row }">
                <el-icon :size="24"><component :is="getIconComponent(row.icon)" /></el-icon>
              </template>
            </el-table-column>
            <el-table-column prop="title" :label="$t('systemThemeTemplate.titleCol')" width="120" />
            <el-table-column prop="description" :label="$t('common.description')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="link_url" :label="$t('systemThemeTemplate.link')" min-width="150" show-overflow-tooltip />
            <el-table-column :label="$t('common.status')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" :label="$t('common.sort')" width="80" />
            <el-table-column :label="$t('common.action')" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditFeature(row)">{{ $t('common.edit') }}</el-button>
                <el-button
                  :type="row.status === 1 ? 'warning' : 'success'"
                  link size="small"
                  @click="handleToggleFeatureStatus(row)"
                >
                  {{ row.status === 1 ? $t('common.disable') : $t('common.enable') }}
                </el-button>
                <el-button type="danger" link size="small" @click="handleDeleteFeature(row)">{{ $t('common.delete') }}</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-alert type="info" :closable="false" style="margin-top: 16px;">
            <template #title>
              <span>{{ $t('systemThemeTemplate.featureTip') }}</span>
            </template>
          </el-alert>
        </art-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 轮播图编辑弹窗 -->
    <el-dialog
      v-model="bannerDialogVisible"
      :title="editingBanner ? $t('systemThemeTemplate.editBanner') : $t('systemThemeTemplate.addBanner')"
      width="600px"
    >
      <el-form :model="bannerForm" label-width="100px">
        <el-form-item :label="$t('systemThemeTemplate.titleCol')" required>
          <el-input v-model="bannerForm.title" :placeholder="$t('systemThemeTemplate.enterBannerTitle')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')">
          <el-input v-model="bannerForm.description" type="textarea" :rows="2" :placeholder="$t('systemThemeTemplate.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('systemThemeTemplate.imageUrl')" required>
          <el-input v-model="bannerForm.media_url" :placeholder="$t('systemThemeTemplate.enterImageUrl')" />
        </el-form-item>
        <el-form-item :label="$t('systemThemeTemplate.link')">
          <el-input v-model="bannerForm.link_url" :placeholder="$t('systemThemeTemplate.enterLink')" />
        </el-form-item>
        <el-form-item :label="$t('systemThemeTemplate.buttonText')">
          <el-input v-model="bannerForm.button_text" :placeholder="$t('systemThemeTemplate.enterButtonText')" />
        </el-form-item>
        <el-form-item :label="$t('systemThemeTemplate.openNew')">
          <el-switch v-model="bannerForm.open_new" />
        </el-form-item>
        <el-form-item :label="$t('common.sort')">
          <el-input-number v-model="bannerForm.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="bannerForm.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bannerDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveBanner" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- 特色区域编辑弹窗 -->
    <el-dialog
      v-model="featureDialogVisible"
      :title="editingFeature ? $t('systemThemeTemplate.editFeature') : $t('systemThemeTemplate.addFeature')"
      width="500px"
    >
      <el-form :model="featureForm" label-width="80px">
        <el-form-item :label="$t('systemThemeTemplate.titleCol')" required>
          <el-input v-model="featureForm.title" :placeholder="$t('systemThemeTemplate.enterTitle')" />
        </el-form-item>
        <el-form-item :label="$t('common.description')">
          <el-input v-model="featureForm.description" type="textarea" :rows="2" :placeholder="$t('systemThemeTemplate.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('systemThemeTemplate.icon')">
            <el-select v-model="featureForm.icon" :placeholder="$t('systemThemeTemplate.selectIcon')">
            <el-option :label="$t('systemThemeTemplate.iconWarning')" value="Warning" />
            <el-option :label="$t('systemThemeTemplate.iconTrend')" value="TrendCharts" />
            <el-option :label="$t('systemThemeTemplate.iconSupport')" value="Headset" />
            <el-option :label="$t('systemThemeTemplate.iconConnection')" value="Connection" />
            <el-option :label="$t('systemThemeTemplate.iconStar')" value="Star" />
            <el-option :label="$t('systemThemeTemplate.iconLightning')" value="Lightning" />
            <el-option :label="$t('systemThemeTemplate.iconLock')" value="Lock" />
            <el-option :label="$t('systemThemeTemplate.iconGlobal')" value="Location" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('systemThemeTemplate.link')">
          <el-input v-model="featureForm.link_url" :placeholder="$t('systemThemeTemplate.enterLink')" />
        </el-form-item>
        <el-form-item :label="$t('common.sort')">
          <el-input-number v-model="featureForm.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-switch v-model="featureForm.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="featureDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveFeature" :loading="saving">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, markRaw } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Warning, TrendCharts, Headset, Connection, Star, Lightning, Lock, Location } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const activeTab = ref('banners')
const saving = ref(false)

// 图标映射
const iconMap: Record<string, any> = {
  Shield: markRaw(Warning),
  TrendCharts: markRaw(TrendCharts),
  Headset: markRaw(Headset),
  Connection: markRaw(Connection),
  Star: markRaw(Star),
  Lightning: markRaw(Lightning),
  Lock: markRaw(Lock),
  Location: markRaw(Location)
}

const getIconComponent = (iconName: string) => {
  return iconMap[iconName] || markRaw(Warning)
}

// ========== 轮播图 ==========
const banners = ref<any[]>([])
const bannersLoading = ref(false)
const bannerDialogVisible = ref(false)
const editingBanner = ref<any>(null)
const bannerForm = ref({
  title: '',
  description: '',
  media_url: '',
  link_url: '',
  button_text: '',
  open_new: false,
  sort_order: 0,
  status: 1
})

const fetchBanners = async () => {
  bannersLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/banners' })
    banners.value = res?.list || res || []
  } catch (error) {
    console.error('获取轮播图失败:', error)
  } finally {
    bannersLoading.value = false
  }
}

const handleAddBanner = () => {
  editingBanner.value = null
  bannerForm.value = {
    title: '',
    description: '',
    media_url: '',
    link_url: '',
    button_text: '',
    open_new: false,
    sort_order: 0,
    status: 1
  }
  bannerDialogVisible.value = true
}

const handleEditBanner = (row: any) => {
  editingBanner.value = row
  bannerForm.value = {
    title: row.title,
    description: row.description || '',
    media_url: row.media_url,
    link_url: row.link_url || '',
    button_text: row.button_text || '',
    open_new: row.open_new || false,
    sort_order: row.sort_order || 0,
    status: row.status
  }
  bannerDialogVisible.value = true
}

const handleSaveBanner = async () => {
  saving.value = true
  try {
    if (editingBanner.value) {
      await request.put({ url: `/api/admin/banners/${editingBanner.value.id}`, data: bannerForm.value })
      ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/banners', data: bannerForm.value })
      ElMessage.success($t('common.addSuccess'))
    }
    bannerDialogVisible.value = false
    fetchBanners()
  } catch (error) {
    console.error('保存轮播图失败:', error)
  } finally {
    saving.value = false
  }
}

const handleToggleBannerStatus = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/banners/${row.id}/status` })
    ElMessage.success($t('common.operateSuccess'))
    fetchBanners()
  } catch (error) {
    console.error('更新状态失败:', error)
  }
}

const handleDeleteBanner = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('systemThemeTemplate.confirmDeleteBanner'), $t('common.tips'), { type: 'warning' })
    await request.del({ url: `/api/admin/banners/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchBanners()
  } catch (error) {
    if (error !== 'cancel') console.error('删除失败:', error)
  }
}

// ========== 特色区域 ==========
const features = ref<any[]>([])
const featuresLoading = ref(false)
const featureDialogVisible = ref(false)
const editingFeature = ref<any>(null)
const featureForm = ref({
  title: '',
  description: '',
  icon: '',
  link_url: '',
  sort_order: 0,
  status: 1
})

const fetchFeatures = async () => {
  featuresLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/homepage-features' })
    features.value = res?.list || res || []
  } catch (error) {
    console.error('获取特色区域失败:', error)
  } finally {
    featuresLoading.value = false
  }
}

const handleAddFeature = () => {
  editingFeature.value = null
  featureForm.value = {
    title: '',
    description: '',
    icon: '',
    link_url: '',
    sort_order: 0,
    status: 1
  }
  featureDialogVisible.value = true
}

const handleEditFeature = (row: any) => {
  editingFeature.value = row
  featureForm.value = {
    title: row.title,
    description: row.description || '',
    icon: row.icon || '',
    link_url: row.link_url || '',
    sort_order: row.sort_order || 0,
    status: row.status
  }
  featureDialogVisible.value = true
}

const handleSaveFeature = async () => {
  saving.value = true
  try {
    if (editingFeature.value) {
      await request.put({ url: `/api/admin/homepage-features/${editingFeature.value.id}`, data: featureForm.value })
      ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({ url: '/api/admin/homepage-features', data: featureForm.value })
      ElMessage.success($t('common.addSuccess'))
    }
    featureDialogVisible.value = false
    fetchFeatures()
  } catch (error) {
    console.error('保存特色项失败:', error)
  } finally {
    saving.value = false
  }
}

const handleToggleFeatureStatus = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/homepage-features/${row.id}/status` })
    ElMessage.success($t('common.operateSuccess'))
    fetchFeatures()
  } catch (error) {
    console.error('更新状态失败:', error)
  }
}

const handleDeleteFeature = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('systemThemeTemplate.confirmDeleteFeature'), $t('common.tips'), { type: 'warning' })
    await request.del({ url: `/api/admin/homepage-features/${row.id}` })
    ElMessage.success($t('common.deleteSuccess'))
    fetchFeatures()
  } catch (error) {
    if (error !== 'cancel') console.error('删除失败:', error)
  }
}

onMounted(() => {
  fetchBanners()
  fetchFeatures()
})
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
