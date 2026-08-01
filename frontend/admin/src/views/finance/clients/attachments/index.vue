<template>
  <div class="attachments-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户附件管理</span>
          <el-upload
            :action="uploadUrl"
            :headers="uploadHeaders"
            :data="{ client_id: clientId }"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :show-file-list="false"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              上传附件
            </el-button>
          </el-upload>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="文件名" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="图片" value="image" />
            <el-option label="文档" value="document" />
            <el-option label="压缩包" value="archive" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="文件名" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="100" align="center">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column prop="operator" label="上传人" width="100" />
        <el-table-column prop="created_at" label="上传时间" width="170" />
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDownload(row)">下载</el-button>
            <el-button type="primary" link @click="handleEditRemark(row)">备注</el-button>
            <el-popconfirm title="确定删除该附件吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 备注对话框 -->
    <el-dialog v-model="remarkVisible" title="编辑备注" width="500px">
      <el-form :model="remarkForm" ref="remarkFormRef" label-width="80px">
        <el-form-item label="备注">
          <el-input v-model="remarkForm.remark" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="remarkVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRemarkSubmit" :loading="remarkLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Upload } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const clientId = route.params.id as string

const loading = ref(false)
const remarkLoading = ref(false)

const uploadUrl = '/api/admin/attachments/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))

const searchForm = reactive({
  keyword: '',
  type: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])
const remarkVisible = ref(false)
const remarkFormRef = ref<FormInstance>()
const currentRow = ref<any>(null)

const remarkForm = reactive({
  remark: ''
})

const getTypeText = (type: string) => {
  const map: Record<string, string> = { image: '图片', document: '文档', archive: '压缩包', other: '其他' }
  return map[type] || '其他'
}

const getTypeTag = (type: string) => {
  const map: Record<string, string> = { image: 'success', document: 'primary', archive: 'warning', other: 'info' }
  return (map[type] || 'info') as any
}

const formatSize = (bytes: number | undefined) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return `${size.toFixed(1)} ${units[i]}`
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, client_id: clientId }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.type) params.type = searchForm.type
    const data = await request.get({ url: '/api/admin/attachments', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.keyword = ''; searchForm.type = ''; handleSearch() }

const handleUploadSuccess = (response: any) => {
  ElMessage.success('上传成功')
  fetchData()
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

const handleDownload = (row: any) => {
  window.open(`/api/admin/attachments/${row.id}/download`, '_blank')
}

const handleEditRemark = (row: any) => {
  currentRow.value = row
  remarkForm.remark = row.remark || ''
  remarkVisible.value = true
}

const handleRemarkSubmit = async () => {
  remarkLoading.value = true
  try {
    await request.put({ url: `/api/admin/attachments/${currentRow.value.id}`, params: { remark: remarkForm.remark } })
    ElMessage.success('更新成功')
    remarkVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('更新失败')
  } finally {
    remarkLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/attachments/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.attachments-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
