<template>
  <div class="marketplace-sell">
    <div class="page-header">
      <el-button @click="$router.back()" text>
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>挂售主机</h2>
    </div>

    <el-card class="sell-form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="选择主机" prop="host_id">
          <el-select
            v-model="form.host_id"
            placeholder="请选择要挂售的主机"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="host in myHosts"
              :key="host.id"
              :label="`${host.product_name} - ${host.dedicated_ip || 'ID:' + host.id}`"
              :value="host.id"
            >
              <div class="host-option">
                <span class="host-name">{{ host.product_name }}</span>
                <span class="host-ip">{{ host.dedicated_ip || '' }}</span>
                <span class="host-expire" v-if="host.expires_at">
                  到期: {{ formatDate(host.expires_at) }}
                </span>
              </div>
            </el-option>
          </el-select>
          <div class="form-tip">只能挂售持有超过 {{ config.minHoldDays }} 天的主机</div>
        </el-form-item>

        <el-form-item label="售价" prop="sell_price">
          <el-input-number
            v-model="form.sell_price"
            :min="0.01"
            :precision="2"
            :step="1"
            style="width: 200px"
          />
          <span class="price-unit">元/月</span>
          <div class="form-tip" v-if="selectedHost">
            官网月付: ¥{{ selectedHost.product?.price?.toFixed(2) || '-' }}
          </div>
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="描述一下主机的使用情况、配置特点等..."
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            确认挂售
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 我的挂售列表 -->
    <el-card class="my-listings-card" v-if="myListings.length > 0">
      <template #header>
        <div class="card-header">
          <span>我的挂售</span>
        </div>
      </template>

      <el-table :data="myListings" style="width: 100%">
        <el-table-column prop="product_name" label="主机" />
        <el-table-column label="售价" width="120">
          <template #default="{ row }">
            <span class="price">¥{{ row.sell_price?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="浏览" width="80" prop="view_count" />
        <el-table-column label="挂售时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 1"
              type="primary"
              size="small"
              text
              @click="editListing(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="row.status === 1"
              type="danger"
              size="small"
              text
              @click="removeListing(row)"
            >
              下架
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" title="编辑挂售" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="售价">
          <el-input-number
            v-model="editForm.sell_price"
            :min="0.01"
            :precision="2"
            :step="1"
          />
          <span class="price-unit">元/月</span>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="editing" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()

const myHosts = ref<any[]>([])
const myListings = ref<any[]>([])
const config = ref({
  minHoldDays: 7
})

const formRef = ref()
const submitting = ref(false)
const form = ref({
  host_id: undefined as number | undefined,
  sell_price: 10,
  description: ''
})

const rules = {
  host_id: [{ required: true, message: '请选择主机', trigger: 'change' }],
  sell_price: [{ required: true, message: '请输入售价', trigger: 'blur' }]
}

const selectedHost = computed(() => {
  return myHosts.value.find(h => h.id === form.value.host_id)
})

// 编辑
const editVisible = ref(false)
const editing = ref(false)
const editForm = ref({
  id: 0,
  sell_price: 0,
  description: ''
})

onMounted(() => {
  fetchMyHosts()
  fetchMyListings()
})

async function fetchMyHosts() {
  try {
    const res = await request.get('/api/v2/hosts')
    myHosts.value = res.data?.list || res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function fetchMyListings() {
  try {
    const res = await request.get('/api/v1/marketplace/listings/mine')
    myListings.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate()

  submitting.value = true
  try {
    await request.post('/api/v1/marketplace/listings', form.value)
    ElMessage.success('挂售成功')
    form.value = { host_id: undefined, sell_price: 10, description: '' }
    fetchMyListings()
  } catch (e: any) {
    ElMessage.error(e.message || '挂售失败')
  } finally {
    submitting.value = false
  }
}

function editListing(row: any) {
  editForm.value = {
    id: row.id,
    sell_price: row.sell_price,
    description: row.description || ''
  }
  editVisible.value = true
}

async function saveEdit() {
  editing.value = true
  try {
    await request.put(`/api/v1/marketplace/listings/${editForm.value.id}`, {
      sell_price: editForm.value.sell_price,
      description: editForm.value.description
    })
    ElMessage.success('更新成功')
    editVisible.value = false
    fetchMyListings()
  } catch (e: any) {
    ElMessage.error(e.message || '更新失败')
  } finally {
    editing.value = false
  }
}

async function removeListing(row: any) {
  try {
    await ElMessageBox.confirm('确定要下架此挂售吗？', '提示', {
      type: 'warning'
    })
    await request.delete(`/api/v1/marketplace/listings/${row.id}`)
    ElMessage.success('下架成功')
    fetchMyListings()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '下架失败')
    }
  }
}

function getStatusType(status: number) {
  const map: Record<number, string> = {
    1: 'success',
    2: 'info',
    3: 'warning'
  }
  return map[status] || 'info'
}

function getStatusText(status: number) {
  const map: Record<number, string> = {
    1: '在售',
    2: '已售',
    3: '已下架'
  }
  return map[status] || '未知'
}

function formatDate(date: string): string {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped lang="scss">
.marketplace-sell {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;

  h2 {
    margin: 0;
  }
}

.sell-form-card {
  margin-bottom: 24px;
}

.price-unit {
  margin-left: 8px;
  color: #666;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.host-option {
  display: flex;
  align-items: center;
  gap: 12px;

  .host-name {
    flex: 1;
  }

  .host-ip {
    color: #666;
  }

  .host-expire {
    color: #999;
    font-size: 12px;
  }
}

.my-listings-card {
  .card-header {
    font-weight: 600;
  }

  .price {
    color: #ff4757;
    font-weight: 600;
  }
}
</style>
