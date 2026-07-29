<template>
  <div class="community-page page-container">
    <div class="page-header">
      <h2>社区管理</h2>
    </div>

    <el-tabs v-model="activeTab" @tab-click="handleTabClick">
      <el-tab-pane label="语言配置" name="language">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>系统语言配置</span>
              <el-button type="primary" size="small" @click="handleSaveLanguage">保存配置</el-button>
            </div>
          </template>
          <el-form :model="languageForm" label-width="150px">
            <el-form-item label="允许用户切换语言">
              <el-switch v-model="languageForm.allow_user_language" />
            </el-form-item>
            <el-form-item label="默认语言">
              <el-select v-model="languageForm.language" placeholder="请选择默认语言">
                <el-option v-for="lang in languages" :key="lang.code" :label="lang.name" :value="lang.code" />
              </el-select>
            </el-form-item>
            <el-form-item label="系统语言">
              <el-select v-model="languageForm.language_system" placeholder="请选择系统语言">
                <el-option v-for="lang in languages" :key="lang.code" :label="lang.name" :value="lang.code" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="语言包管理" name="language-packs">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>语言包列表</span>
              <el-button type="primary" size="small" @click="handleUploadLanguagePack">上传语言包</el-button>
            </div>
          </template>
          <el-table :data="languagePacks" border fit>
            <el-table-column prop="display_name" label="语言名称" />
            <el-table-column prop="file_name" label="文件名" />
            <el-table-column prop="display_flag" label="标识" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="handleEditLanguagePack(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="handleDeleteLanguagePack(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 上传语言包弹窗 -->
    <el-dialog
      v-model="showUploadDialog"
      title="上传语言包"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-upload
        ref="uploadRef"
        class="language-upload"
        drag
        :auto-upload="false"
        :limit="1"
        :on-change="handleFileChange"
        :on-exceed="handleExceed"
        accept=".json,.js"
      >
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            只能上传 .json 或 .js 文件
          </div>
        </template>
      </el-upload>
      <el-form :model="uploadForm" label-width="80px" style="margin-top: 16px;">
        <el-form-item label="语言名称">
          <el-input v-model="uploadForm.display_name" placeholder="例如：中文" />
        </el-form-item>
        <el-form-item label="语言标识">
          <el-input v-model="uploadForm.display_flag" placeholder="例如：CN" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" :loading="uploadLoading" @click="submitUpload">
          上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 编辑语言包弹窗 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑语言包"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="语言名称">
          <el-input v-model="editForm.display_name" placeholder="语言名称" />
        </el-form-item>
        <el-form-item label="语言标识">
          <el-input v-model="editForm.display_flag" placeholder="语言标识" />
        </el-form-item>
        <el-form-item label="语言内容">
          <el-input
            v-model="editForm.content"
            type="textarea"
            :rows="15"
            placeholder="JSON 格式的语言内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" :loading="editLoading" @click="submitEdit">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import request from '@/utils/request'

const activeTab = ref('language')
const uploadRef = ref()

// 上传相关
const showUploadDialog = ref(false)
const uploadLoading = ref(false)
const uploadFile = ref(null)
const uploadForm = reactive({
  display_name: '',
  display_flag: ''
})

// 编辑相关
const showEditDialog = ref(false)
const editLoading = ref(false)
const editingPack = ref(null)
const editForm = reactive({
  display_name: '',
  display_flag: '',
  content: ''
})

const languageForm = reactive({
  allow_user_language: false,
  language: '',
  language_system: ''
})

const languages = ref([
  { code: 'zh', name: '中文' },
  { code: 'en', name: 'English' },
  { code: 'ja', name: '日本語' }
])

const languagePacks = ref([])

const fetchLanguageConfig = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/community/language-config')
    if (data.data) {
      Object.assign(languageForm, data.data)
    }
  } catch (error) {
    ElMessage.error('获取语言配置失败')
  }
}

const fetchLanguagePacks = async () => {
  try {
    const { data } = await request.get('/admin/api/v1/community/language-packs')
    languagePacks.value = data.data?.list || data.data || []
  } catch (error) {
    ElMessage.error('获取语言包列表失败')
  }
}

const handleTabClick = (tab) => {
  if (tab.name === 'language') {
    fetchLanguageConfig()
  } else if (tab.name === 'language-packs') {
    fetchLanguagePacks()
  }
}

const handleSaveLanguage = async () => {
  try {
    await request.put('/admin/api/v1/community/language-config', languageForm)
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const handleUploadLanguagePack = () => {
  uploadFile.value = null
  uploadForm.display_name = ''
  uploadForm.display_flag = ''
  showUploadDialog.value = true
}

const handleFileChange = (file) => {
  uploadFile.value = file.raw
  // 尝试从文件名提取语言标识
  const fileName = file.name.replace(/\.(json|js)$/, '')
  if (!uploadForm.display_name) {
    uploadForm.display_name = fileName
  }
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const submitUpload = async () => {
  if (!uploadFile.value) {
    ElMessage.warning('请选择文件')
    return
  }
  if (!uploadForm.display_name) {
    ElMessage.warning('请输入语言名称')
    return
  }

  uploadLoading.value = true
  try {
    const formData = new FormData()
    formData.append('file', uploadFile.value)
    formData.append('display_name', uploadForm.display_name)
    formData.append('display_flag', uploadForm.display_flag)

    const response = await fetch('/api/v1/admin/language-pack/upload', {
      method: 'POST',
      body: formData
    })
    const data = await response.json()
    if (data.code === 0) {
      showUploadDialog.value = false
      ElMessage.success('上传成功')
      fetchLanguagePacks()
    } else {
      ElMessage.error(data.message || '上传失败')
    }
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploadLoading.value = false
  }
}

const handleEditLanguagePack = async (row) => {
  editingPack.value = row
  editForm.display_name = row.display_name || ''
  editForm.display_flag = row.display_flag || ''
  editForm.content = ''
  
  // 获取语言包内容
  try {
    const { data } = await fetch(`/api/v1/admin/language-pack/${row.id}`).then(r => r.json())
    if (data?.content) {
      editForm.content = typeof data.content === 'string' ? data.content : JSON.stringify(data.content, null, 2)
    }
  } catch (error) {
    console.error('获取语言包内容失败', error)
  }
  
  showEditDialog.value = true
}

const submitEdit = async () => {
  if (!editForm.display_name) {
    ElMessage.warning('请输入语言名称')
    return
  }

  editLoading.value = true
  try {
    const response = await fetch(`/api/v1/admin/language-pack/${editingPack.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        display_name: editForm.display_name,
        display_flag: editForm.display_flag,
        content: editForm.content
      })
    })
    const data = await response.json()
    if (data.code === 0) {
      showEditDialog.value = false
      ElMessage.success('保存成功')
      fetchLanguagePacks()
    } else {
      ElMessage.error(data.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    editLoading.value = false
  }
}

const handleDeleteLanguagePack = (row) => {
  ElMessageBox.confirm('确认删除该语言包？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/admin/api/v1/community/language-packs/${row.id}`)
      ElMessage.success('删除成功')
      fetchLanguagePacks()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

onMounted(() => {
  fetchLanguageConfig()
  fetchLanguagePacks()
})
</script>

<style scoped>
.community-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.language-upload {
  width: 100%;
}

.language-upload :deep(.el-upload-dragger) {
  width: 100%;
}
</style>
