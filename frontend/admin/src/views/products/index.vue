<template>
  <div class="products-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="产品名称" clearable />
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="searchForm.group" placeholder="全部" clearable>
            <el-option label="云服务器" value="cloud" />
            <el-option label="独立服务器" value="dedicated" />
            <el-option label="虚拟主机" value="hosting" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="art-card">
      <div class="table-header">
        <h3>产品列表</h3>
        <el-button type="primary" @click="showAddDialog = true">
          <el-icon><Plus /></el-icon>
          添加产品
        </el-button>
      </div>
      
      <el-table :data="products" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="产品名称" />
        <el-table-column prop="group" label="分组" />
        <el-table-column prop="price" label="价格">
          <template #default="{ row }">
            <span class="amount">¥{{ row.price?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="billing_cycle" label="计费周期" />
        <el-table-column prop="stock" label="库存" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <span class="status-tag" :class="row.status">
              {{ row.status === 'active' ? '上架' : '下架' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editProduct(row)">编辑</el-button>
            <el-button type="primary" link @click="configProduct(row)">配置</el-button>
            <el-button type="danger" link @click="deleteProduct(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const loading = ref(false)
const products = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showAddDialog = ref(false)

const searchForm = ref({
  keyword: '',
  group: ''
})

const handleSearch = () => {
  currentPage.value = 1
}

const resetSearch = () => {
  searchForm.value = { keyword: '', group: '' }
  handleSearch()
}

const editProduct = (product: any) => {
  // TODO: 编辑产品
}

const configProduct = (product: any) => {
  // TODO: 配置产品
}

const deleteProduct = async (product: any) => {
  // TODO: 删除产品
}
</script>

<style scoped lang="scss">
.products-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
    }
  }
  
  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
