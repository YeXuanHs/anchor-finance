<template>
  <div class="ssh-key-manager">
    <div class="toolbar">
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加SSH密钥
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="keyList"
      style="width: 100%"
      empty-text="暂无SSH密钥"
    >
      <el-table-column prop="name" label="密钥名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="fingerprint" label="指纹" min-width="220">
        <template #default="{ row }">
          <span class="mono text-secondary">{{ row.fingerprint || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag size="small" effect="light">
            {{ row.type || 'RSA' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="160" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            size="small"
            @click="handleCopy(row)"
          >
            复制公钥
          </el-button>
          <el-button
            link
            type="danger"
            size="small"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>

    <!-- 添加SSH密钥弹窗 -->
    <el-dialog
      v-model="showCreateDialog"
      title="添加SSH密钥"
      width="560px"
      :close-on-click-modal="false"
      @closed="resetCreateForm"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="密钥名称" required>
          <el-input
            v-model="createForm.name"
            placeholder="请输入密钥名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="公钥内容" required>
          <el-input
            v-model="createForm.public_key"
            type="textarea"
            :rows="6"
            placeholder="粘贴您的SSH公钥（以 ssh-rsa、ssh-ed25519 等开头）"
          />
        </el-form-item>
        <el-alert type="info" :closable="false" show-icon>
          <template #default>
            <p>如何获取SSH公钥？</p>
            <p style="margin-top: 4px; font-size: 12px; color: #909399;">
              运行 <code style="background: #f5f7fa; padding: 2px 6px; border-radius: 4px;">cat ~/.ssh/id_rsa.pub</code>
              或 <code style="background: #f5f7fa; padding: 2px 6px; border-radius: 4px;">ssh-keygen -t ed25519</code> 生成
            </p>
          </template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="confirmCreate">
          确认添加
        </el-button>
      </template>
    </el-dialog>

    <!-- 显示生成的密钥对弹窗 -->
    <el-dialog
      v-model="showKeyPairDialog"
      title="密钥对已生成"
      width="560px"
      :close-on-click-modal="false"
    >
      <el-alert type="warning" :closable="false" show-icon style="margin-bottom: 16px;">
        请立即下载或复制私钥，关闭后将无法再次查看。
      </el-alert>
      <el-form label-width="80px">
        <el-form-item label="公钥">
          <el-input v-model="generatedKeyPair.public_key" type="textarea" :rows="3" readonly />
        </el-form-item>
        <el-form-item label="私钥">
          <el-input v-model="generatedKeyPair.private_key" type="textarea" :rows="6" readonly />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="downloadPrivateKey">下载私钥</el-button>
        <el-button type="primary" @click="copyPrivateKey">复制私钥</el-button>
        <el-button @click="showKeyPairDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()

const loading = ref(false)
const keyList = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showCreateDialog = ref(false)
const createLoading = ref(false)
const createForm = ref({
  name: '',
  public_key: ''
})

const showKeyPairDialog = ref(false)
const generatedKeyPair = ref({
  public_key: '',
  private_key: ''
})

async function fetchList() {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const { data } = await request.get(`/api/v2/hosts/${id}/ssh-keys`, {
      params: { page: page.value, limit: pageSize.value }
    })
    if (data?.data) {
      keyList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取SSH密钥列表失败', error)
  } finally {
    loading.value = false
  }
}

function resetCreateForm() {
  createForm.value = { name: '', public_key: '' }
}

async function confirmCreate() {
  const id = route.params.id
  if (!id) return

  if (!createForm.value.name.trim()) {
    ElMessage.warning('请输入密钥名称')
    return
  }

  createLoading.value = true
  try {
    const { data } = await request.post(`/api/v2/hosts/${id}/ssh-keys`, createForm.value)
    ElMessage.success('SSH密钥添加成功')
    showCreateDialog.value = false

    // 如果返回了生成的密钥对，显示给用户
    if (data?.data?.private_key) {
      generatedKeyPair.value = {
        public_key: data.data.public_key || createForm.value.public_key,
        private_key: data.data.private_key
      }
      showKeyPairDialog.value = true
    }

    fetchList()
  } catch (error: any) {
    ElMessage.error(error.message || '添加SSH密钥失败')
  } finally {
    createLoading.value = false
  }
}

function handleCopy(row: any) {
  if (row.public_key) {
    navigator.clipboard.writeText(row.public_key).then(() => {
      ElMessage.success('公钥已复制到剪贴板')
    }).catch(() => {
      ElMessage.error('复制失败')
    })
  }
}

function copyPrivateKey() {
  navigator.clipboard.writeText(generatedKeyPair.value.private_key).then(() => {
    ElMessage.success('私钥已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

function downloadPrivateKey() {
  const blob = new Blob([generatedKeyPair.value.private_key], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'id_rsa'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  ElMessage.success('私钥已下载')
}

async function handleDelete(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认删除SSH密钥「${row.name}」？删除后使用该密钥的服务器将无法自动认证。`,
      '确认删除',
      { type: 'warning', confirmButtonText: '删除', confirmButtonClass: 'el-button--danger' }
    )
    await request.delete(`/api/v2/hosts/${id}/ssh-keys/${row.id}`)
    ElMessage.success('SSH密钥已删除')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除SSH密钥失败')
    }
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.ssh-key-manager {
  .toolbar {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 16px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.mono {
  font-family: 'Monaco', 'Menlo', monospace;
}

.text-secondary {
  color: #909399;
  font-size: 13px;
}
</style>
