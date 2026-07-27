<template>
  <div class="dcim-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>DCIM 服务器管理</span>
          <el-button type="primary" @click="refreshList">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="物理服务器" name="physical">
          <el-table :data="physicalServers" style="width: 100%">
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="ip" label="IP地址" />
            <el-table-column prop="datacenter" label="机房" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                  {{ row.status === 1 ? '运行中' : '已关机' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="viewDetail(row)">详情</el-button>
                <el-button size="small" type="primary" @click="openConsole(row)">控制台</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="云服务器" name="cloud">
          <el-table :data="cloudServers" style="width: 100%">
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="ip" label="IP地址" />
            <el-table-column prop="cpu" label="CPU" />
            <el-table-column prop="memory" label="内存" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                  {{ row.status === 1 ? '运行中' : '已关机' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="300">
              <template #default="{ row }">
                <el-button size="small" @click="viewDetail(row)">详情</el-button>
                <el-button size="small" @click="bootServer(row)" v-if="row.status !== 1">开机</el-button>
                <el-button size="small" type="warning" @click="shutdownServer(row)" v-if="row.status === 1">关机</el-button>
                <el-button size="small" @click="rebootServer(row)" v-if="row.status === 1">重启</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'

const activeTab = ref('physical')
const physicalServers = ref([])
const cloudServers = ref([])

const refreshList = () => {
  // TODO: 刷新服务器列表
}

const viewDetail = (server: any) => {
  // TODO: 查看服务器详情
}

const openConsole = (server: any) => {
  // TODO: 打开控制台
}

const bootServer = async (server: any) => {
  // TODO: 开机
}

const shutdownServer = async (server: any) => {
  // TODO: 关机
}

const rebootServer = async (server: any) => {
  // TODO: 重启
}
</script>

<style scoped lang="scss">
.dcim-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
