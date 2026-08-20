<template>
  <div class="language-manage">
    <art-card>
      <template #header>
        <div class="card-header">
          <h3>{{ $t('language.title') }}</h3>
          <el-button type="primary" @click="showAddLanguage">
            <el-icon><Plus /></el-icon>
            {{ $t('language.addLanguage') }}
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 语言列表 -->
        <el-tab-pane :label="$t('language.languageList')" name="languages">
          <el-table :data="languages" stripe>
            <el-table-column prop="code" :label="$t('language.languageCode')" width="120" />
            <el-table-column prop="name" :label="$t('language.languageName')" width="150" />
            <el-table-column prop="flag" :label="$t('language.flag')" width="80">
              <template #default="{ row }">
                <span class="flag">{{ row.flag }}</span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('language.default')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.is_default ? 'success' : 'info'" size="small">
                  {{ row.is_default ? $t('language.default') : '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('language.status')" width="80">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="(val: string | number | boolean) => toggleLanguage(row.id, val as number)"
                />
              </template>
            </el-table-column>
            <el-table-column :label="$t('language.operations')" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="editLanguage(row)">{{ $t('language.edit') }}</el-button>
                <el-button size="small" @click="manageTranslations(row)">{{ $t('language.translate') }}</el-button>
                <el-button 
                  v-if="!row.is_default"
                  size="small" 
                  type="warning"
                  @click="setDefault(row.id)"
                >
                  {{ $t('language.setDefault') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 翻译管理 -->
        <el-tab-pane v-if="currentLang" :label="`${$t('language.translationManagement')} - ${currentLang.name}`" name="translations">
          <div class="translation-header">
            <el-select v-model="currentModule" :placeholder="$t('language.selectModule')" clearable @change="loadTranslations">
              <el-option :label="$t('language.allModules')" value="" />
              <el-option :label="$t('language.moduleCommon')" value="common" />
              <el-option :label="$t('language.moduleMenu')" value="menu" />
              <el-option :label="$t('language.moduleAuth')" value="auth" />
              <el-option :label="$t('language.moduleHome')" value="home" />
              <el-option :label="$t('language.moduleService')" value="service" />
              <el-option :label="$t('language.moduleInvoice')" value="invoice" />
              <el-option :label="$t('language.moduleTicket')" value="ticket" />
            </el-select>
            <el-button type="primary" @click="saveTranslations">
              <el-icon><Check /></el-icon>
              {{ $t('language.saveTranslations') }}
            </el-button>
            <el-button @click="importFromZjmf">
              <el-icon><Upload /></el-icon>
              {{ $t('language.importFromZjmf') }}
            </el-button>
          </div>

          <el-table :data="translationList" stripe border>
            <el-table-column prop="key" :label="$t('language.key')" width="250" />
            <el-table-column :label="$t('language.value')">
              <template #default="{ row }">
                <el-input v-model="row.value" :placeholder="$t('language.enterTranslation')" />
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </art-card>

    <!-- 添加/编辑语言对话框 -->
    <el-dialog v-model="langDialogVisible" :title="isEdit ? $t('language.editLanguage') : $t('language.addLanguage')" width="400px">
      <el-form :model="langForm" label-width="80px">
        <el-form-item :label="$t('language.languageCode')">
          <el-input v-model="langForm.code" :placeholder="$t('language.codePlaceholder')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="$t('language.languageName')">
          <el-input v-model="langForm.name" :placeholder="$t('language.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('language.flagCode')">
          <el-input v-model="langForm.flag" :placeholder="$t('language.flagPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="langDialogVisible = false">{{ $t('language.cancel') }}</el-button>
        <el-button type="primary" @click="saveLanguage">{{ $t('language.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Check, Upload } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

interface Language {
  id: number
  code: string
  name: string
  flag: string
  is_default: boolean
  status: number
}

interface Translation {
  key: string
  value: string
}

const activeTab = ref('languages')
const languages = ref<Language[]>([])
const currentLang = ref<Language | null>(null)
const currentModule = ref('')
const translationList = ref<Translation[]>([])
const langDialogVisible = ref(false)
const isEdit = ref(false)
const langForm = ref({ code: '', name: '', flag: '' })
const editingId = ref<number | null>(null)

const loadLanguages = async () => {
  try {
    const res = await request.get({ url: '/api/admin/languages' })
    languages.value = res || []
  } catch (error) {
    console.error('Failed to load languages:', error)
  }
}

const showAddLanguage = () => {
  isEdit.value = false
  langForm.value = { code: '', name: '', flag: '' }
  editingId.value = null
  langDialogVisible.value = true
}

const editLanguage = (lang: any) => {
  isEdit.value = true
  langForm.value = { code: lang.code, name: lang.name, flag: lang.flag }
  editingId.value = lang.id
  langDialogVisible.value = true
}

const saveLanguage = async () => {
  try {
    if (isEdit.value && editingId.value) {
      await request.put({ url: `/api/admin/languages/${editingId.value}`, params: { name: langForm.value.name, flag: langForm.value.flag } })
    } else {
      await request.post({ url: '/api/admin/languages', params: langForm.value })
    }
    ElMessage.success($t('language.operationSuccess'))
    langDialogVisible.value = false
    loadLanguages()
  } catch (error) {
    ElMessage.error($t('language.operationFailed'))
  }
}

const toggleLanguage = async (id: number, status: number) => {
  try {
    await request.put({ url: `/api/admin/languages/${id}`, params: { status } })
    ElMessage.success($t('language.operationSuccess'))
  } catch (error) {
    ElMessage.error($t('language.operationFailed'))
    loadLanguages()
  }
}

const setDefault = async (id: number) => {
  try {
    await ElMessageBox.confirm($t('language.setDefaultConfirm'), $t('language.tips'))
    await request.post({ url: `/api/admin/languages/${id}/default` })
    ElMessage.success($t('language.setDefaultSuccess'))
    loadLanguages()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error($t('language.setDefaultFailed'))
    }
  }
}

const manageTranslations = (lang: any) => {
  currentLang.value = lang
  activeTab.value = 'translations'
  loadTranslations()
}

const loadTranslations = async () => {
  if (!currentLang.value) return
  try {
    const params: any = {}
    if (currentModule.value) params.module = currentModule.value
    const res = await request.get({ url: `/api/admin/languages/${currentLang.value.code}/translations`, params })
    const data = res || {}
    translationList.value = Object.entries(data).map(([key, value]) => ({
      key,
      value: value as string
    }))
  } catch (error) {
    console.error('Failed to load translations:', error)
  }
}

const saveTranslations = async () => {
  if (!currentLang.value) return
  try {
    const data: Record<string, string> = {}
    translationList.value.forEach(item => {
      if (item.key && item.value) {
        data[item.key] = item.value
      }
    })
    await request.post({ url: `/api/admin/languages/${currentLang.value.code}/translations`, params: data })
    ElMessage.success($t('language.saveSuccess'))
  } catch (error) {
    ElMessage.error($t('language.saveFailed'))
  }
}

const importFromZjmf = async () => {
  try {
    await ElMessageBox.confirm($t('language.importConfirm'), $t('language.tips'))
    ElMessage.info($t('language.importAdminTip'))
  } catch (error) {
    // cancelled
  }
}

const handleTabChange = (tab: string | number) => {
  if (tab === 'languages') {
    currentLang.value = null
    loadLanguages()
  }
}

onMounted(() => {
  loadLanguages()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
}

.flag {
  font-size: 18px;
}

.translation-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
