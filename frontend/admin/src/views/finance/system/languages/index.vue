<template>
  <div class="language-manage">
    <art-card>
      <template #header>
        <div class="card-header">
          <h3>语言管理</h3>
          <el-button type="primary" @click="showAddLanguage">
            <el-icon><Plus /></el-icon>
            添加语言
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 语言列表 -->
        <el-tab-pane label="语言列表" name="languages">
          <el-table :data="languages" stripe>
            <el-table-column prop="code" label="语言代码" width="120" />
            <el-table-column prop="name" label="语言名称" width="150" />
            <el-table-column prop="flag" label="国旗" width="80">
              <template #default="{ row }">
                <span class="flag">{{ row.flag }}</span>
              </template>
            </el-table-column>
            <el-table-column label="默认" width="80">
              <template #default="{ row }">
                <el-tag :type="row.is_default ? 'success' : 'info'" size="small">
                  {{ row.is_default ? '默认' : '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="(val: number) => toggleLanguage(row.id, val)"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="editLanguage(row)">编辑</el-button>
                <el-button size="small" @click="manageTranslations(row)">翻译</el-button>
                <el-button 
                  v-if="!row.is_default"
                  size="small" 
                  type="warning"
                  @click="setDefault(row.id)"
                >
                  设为默认
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 翻译管理 -->
        <el-tab-pane v-if="currentLang" :label="`翻译管理 - ${currentLang.name}`" name="translations">
          <div class="translation-header">
            <el-select v-model="currentModule" placeholder="选择模块" clearable @change="loadTranslations">
              <el-option label="全部" value="" />
              <el-option label="通用" value="common" />
              <el-option label="菜单" value="menu" />
              <el-option label="认证" value="auth" />
              <el-option label="首页" value="home" />
              <el-option label="服务" value="service" />
              <el-option label="账单" value="invoice" />
              <el-option label="工单" value="ticket" />
            </el-select>
            <el-button type="primary" @click="saveTranslations">
              <el-icon><Check /></el-icon>
              保存翻译
            </el-button>
            <el-button @click="importFromZjmf">
              <el-icon><Upload /></el-icon>
              从zjmf导入
            </el-button>
          </div>

          <el-table :data="translationList" stripe border>
            <el-table-column prop="key" label="键名" width="250" />
            <el-table-column label="翻译值">
              <template #default="{ row }">
                <el-input v-model="row.value" placeholder="请输入翻译" />
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </art-card>

    <!-- 添加/编辑语言对话框 -->
    <el-dialog v-model="langDialogVisible" :title="isEdit ? '编辑语言' : '添加语言'" width="400px">
      <el-form :model="langForm" label-width="80px">
        <el-form-item label="语言代码">
          <el-input v-model="langForm.code" placeholder="如: en-US" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="语言名称">
          <el-input v-model="langForm.name" placeholder="如: English" />
        </el-form-item>
        <el-form-item label="国旗代码">
          <el-input v-model="langForm.flag" placeholder="如: US" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="langDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveLanguage">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Check, Upload } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

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

// 加载语言列表
const loadLanguages = async () => {
  try {
    const res = await request.get('/api/admin/languages')
    languages.value = res.data.data || []
  } catch (error) {
    console.error('加载语言列表失败:', error)
  }
}

// 显示添加语言对话框
const showAddLanguage = () => {
  isEdit.value = false
  langForm.value = { code: '', name: '', flag: '' }
  editingId.value = null
  langDialogVisible.value = true
}

// 编辑语言
const editLanguage = (lang: Language) => {
  isEdit.value = true
  langForm.value = { code: lang.code, name: lang.name, flag: lang.flag }
  editingId.value = lang.id
  langDialogVisible.value = true
}

// 保存语言
const saveLanguage = async () => {
  try {
    if (isEdit.value && editingId.value) {
      await request.put(`/api/admin/languages/${editingId.value}`, {
        name: langForm.value.name,
        flag: langForm.value.flag
      })
    } else {
      await request.post('/api/admin/languages', langForm.value)
    }
    ElMessage.success('操作成功')
    langDialogVisible.value = false
    loadLanguages()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 切换语言状态
const toggleLanguage = async (id: number, status: number) => {
  try {
    await request.put(`/api/admin/languages/${id}`, { status })
    ElMessage.success('操作成功')
  } catch (error) {
    ElMessage.error('操作失败')
    loadLanguages()
  }
}

// 设置默认语言
const setDefault = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定将此语言设为默认吗？', '提示')
    await request.post(`/api/admin/languages/${id}/default`)
    ElMessage.success('设置成功')
    loadLanguages()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('设置失败')
    }
  }
}

// 管理翻译
const manageTranslations = (lang: Language) => {
  currentLang.value = lang
  activeTab.value = 'translations'
  loadTranslations()
}

// 加载翻译
const loadTranslations = async () => {
  if (!currentLang.value) return
  try {
    const params: any = {}
    if (currentModule.value) params.module = currentModule.value
    const res = await request.get(`/api/admin/languages/${currentLang.value.code}/translations`, { params })
    const data = res.data.data || {}
    translationList.value = Object.entries(data).map(([key, value]) => ({
      key,
      value: value as string
    }))
  } catch (error) {
    console.error('加载翻译失败:', error)
  }
}

// 保存翻译
const saveTranslations = async () => {
  if (!currentLang.value) return
  try {
    const data: Record<string, string> = {}
    translationList.value.forEach(item => {
      if (item.key && item.value) {
        data[item.key] = item.value
      }
    })
    await request.post(`/api/admin/languages/${currentLang.value.code}/translations`, data)
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 从zjmf导入
const importFromZjmf = async () => {
  try {
    await ElMessageBox.confirm('从zjmf语言包导入翻译数据，确定继续吗？', '提示')
    // 这里需要从zjmf的PHP文件解析语言数据
    // 实际实现需要后端配合解析PHP文件
    ElMessage.info('请联系管理员导入zjmf语言包数据')
  } catch (error) {
    // 取消
  }
}

// 标签切换
const handleTabChange = (tab: string) => {
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
