<template>
  <div class="vnc-console-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-page-header @back="$router.back()">
        <template #content>
          <span class="page-title">VNC 控制台</span>
        </template>
        <template #extra>
          <div class="header-actions">
            <el-tag :type="connected ? 'success' : 'info'" size="small" effect="light" round>
              {{ connected ? '已连接' : '未连接' }}
            </el-tag>
            <el-button-group>
              <el-button size="small" @click="toggleFullscreen" :icon="FullScreen">
                {{ isFullscreen ? '退出全屏' : '全屏' }}
              </el-button>
              <el-button size="small" @click="sendCtrlAltDel" :disabled="!connected">
                Ctrl+Alt+Del
              </el-button>
              <el-button size="small" type="danger" plain @click="disconnect" :disabled="!connected">
                断开
              </el-button>
            </el-button-group>
          </div>
        </template>
      </el-page-header>
    </div>

    <!-- 连接信息 -->
    <el-card shadow="never" class="info-card" v-if="!connected && !connecting">
      <div class="connect-info">
        <el-icon :size="48" color="#909399"><Monitor /></el-icon>
        <h3>VNC 远程控制台</h3>
        <p>通过浏览器远程访问您的服务器桌面</p>
        <div class="server-info" v-if="serverInfo">
          <div class="info-row">
            <span class="label">服务器：</span>
            <span class="value">{{ serverInfo.name || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="label">IP 地址：</span>
            <span class="value mono">{{ serverInfo.ip || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="label">VNC 端口：</span>
            <span class="value mono">{{ serverInfo.vncPort || '5900' }}</span>
          </div>
        </div>
        <el-button type="primary" size="large" @click="connectVnc" :loading="connecting">
          <el-icon><Connection /></el-icon>
          连接 VNC
        </el-button>
      </div>
    </el-card>

    <!-- 连接中状态 -->
    <div v-if="connecting" class="connecting-overlay">
      <el-icon class="loading-spinner" :size="48" color="#409eff"><Loading /></el-icon>
      <p>正在建立 VNC 连接...</p>
    </div>

    <!-- VNC 画布容器 -->
    <div
      ref="vncContainer"
      class="vnc-container"
      v-show="connected"
      @contextmenu.prevent
    >
      <canvas ref="vncCanvas" class="vnc-canvas"></canvas>
    </div>

    <!-- 连接状态栏 -->
    <div class="status-bar" v-if="connected">
      <div class="status-left">
        <el-icon :size="14" color="#67c23a"><CircleCheck /></el-icon>
        <span>已连接到 {{ serverInfo?.ip || 'VNC Server' }}</span>
      </div>
      <div class="status-right">
        <span>分辨率：{{ resolution }}</span>
        <span class="separator">|</span>
        <span>延迟：{{ latency }}ms</span>
        <span class="separator">|</span>
        <span>编码：{{ encoding }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Monitor, Connection, FullScreen, Loading, CircleCheck
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()

// VNC 状态
const connecting = ref(false)
const connected = ref(false)
const isFullscreen = ref(false)
const resolution = ref('1024x768')
const latency = ref(0)
const encoding = ref('Tight')

// DOM 引用
const vncContainer = ref<HTMLDivElement>()
const vncCanvas = ref<HTMLCanvasElement>()

// 服务器信息
const serverInfo = ref<any>(null)

// VNC 连接相关
let ws: WebSocket | null = null
let pingTimer: ReturnType<typeof setInterval> | null = null
let latencyTimer: ReturnType<typeof setInterval> | null = null

// 获取服务器信息
async function fetchServerInfo() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  try {
    const { data } = await request.get(`/api/v2/hosts/${hostId}`)
    if (data?.data) {
      serverInfo.value = {
        name: data.data.product_name || data.data.name,
        ip: data.data.dedicated_ip,
        vncPort: data.data.vnc_port || 5900
      }
    }
  } catch (e) {
    console.error('Failed to fetch server info:', e)
  }
}

// 连接 VNC
async function connectVnc() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) {
    ElMessage.warning('缺少服务器ID参数')
    return
  }

  connecting.value = true

  try {
    // 获取 VNC 连接信息
    const { data } = await request.post(`/api/v2/hosts/${hostId}/vnc`)
    const vncInfo = data?.data

    if (!vncInfo?.url) {
      throw new Error('获取VNC连接地址失败')
    }

    // 建立 WebSocket 连接
    await initVncConnection(vncInfo.url, vncInfo.token)
  } catch (e: any) {
    ElMessage.error(e.message || '连接VNC失败')
    connecting.value = false
  }
}

// 初始化 VNC 连接（通过 noVNC WebSocket 代理）
async function initVncConnection(url: string, token?: string) {
  try {
    const wsUrl = token ? `${url}?token=${token}` : url
    ws = new WebSocket(wsUrl)
    ws.binaryType = 'arraybuffer'

    ws.onopen = () => {
      connected.value = true
      connecting.value = false
      ElMessage.success('VNC 连接成功')
      startPing()
      startLatencyMonitor()
      initCanvas()
    }

    ws.onmessage = (event) => {
      handleVncMessage(event.data)
    }

    ws.onerror = () => {
      ElMessage.error('VNC 连接错误')
      disconnect()
    }

    ws.onclose = () => {
      if (connected.value) {
        ElMessage.warning('VNC 连接已断开')
        disconnect()
      }
    }
  } catch (e: any) {
    connecting.value = false
    throw e
  }
}

// 初始化画布
function initCanvas() {
  nextTick(() => {
    if (!vncCanvas.value || !vncContainer.value) return

    const container = vncContainer.value
    vncCanvas.value.width = container.clientWidth
    vncCanvas.value.height = container.clientHeight

    // 监听键盘和鼠标事件
    vncCanvas.value.addEventListener('mousedown', handleMouseDown)
    vncCanvas.value.addEventListener('mouseup', handleMouseUp)
    vncCanvas.value.addEventListener('mousemove', handleMouseMove)
    vncCanvas.value.addEventListener('wheel', handleWheel)
    vncCanvas.value.addEventListener('keydown', handleKeyDown)
    vncCanvas.value.addEventListener('keyup', handleKeyUp)

    resolution.value = `${vncCanvas.value.width}x${vncCanvas.value.height}`
  })
}

// 处理 VNC 消息
function handleVncMessage(data: ArrayBuffer) {
  // 实际实现需根据 noVNC 协议解析帧数据
  // 这里是框架代码，实际项目需要集成 noVNC 客户端库
  if (!vncCanvas.value) return
  const ctx = vncCanvas.value.getContext('2d')
  if (!ctx) return

  // 将接收到的帧数据绘制到画布
  // 实际项目中应使用 noVNC 的 RFB 协议处理
}

// 鼠标事件处理
function handleMouseDown(e: MouseEvent) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return
  const rect = vncCanvas.value!.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top
  sendPointerEvent(x, y, e.button, true)
}

function handleMouseUp(e: MouseEvent) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return
  const rect = vncCanvas.value!.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top
  sendPointerEvent(x, y, e.button, false)
}

function handleMouseMove(e: MouseEvent) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return
  const rect = vncCanvas.value!.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top
  sendPointerEvent(x, y, -1, false)
}

function handleWheel(e: WheelEvent) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return
  e.preventDefault()
  const rect = vncCanvas.value!.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top
  // 发送滚轮事件
  const button = e.deltaY > 0 ? 5 : 4
  sendPointerEvent(x, y, button, true)
  sendPointerEvent(x, y, button, false)
}

// 键盘事件处理
function handleKeyDown(e: KeyboardEvent) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return
  e.preventDefault()
  sendKeyEvent(e.keyCode, true)
}

function handleKeyUp(e: KeyboardEvent) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return
  e.preventDefault()
  sendKeyEvent(e.keyCode, false)
}

// 发送指针事件到 VNC 服务器
function sendPointerEvent(x: number, y: number, button: number, down: boolean) {
  if (!ws) return
  const msg = JSON.stringify({
    type: 'pointer',
    x: Math.round(x),
    y: Math.round(y),
    button,
    down
  })
  ws.send(msg)
}

// 发送键盘事件到 VNC 服务器
function sendKeyEvent(keyCode: number, down: boolean) {
  if (!ws) return
  const msg = JSON.stringify({
    type: 'key',
    keyCode,
    down
  })
  ws.send(msg)
}

// 发送 Ctrl+Alt+Del
function sendCtrlAltDel() {
  if (!ws || ws.readyState !== WebSocket.OPEN) return
  // 先按下 Ctrl, Alt, Del，再释放
  sendKeyEvent(17, true)  // Ctrl down
  sendKeyEvent(18, true)  // Alt down
  sendKeyEvent(46, true)  // Delete down
  sendKeyEvent(46, false) // Delete up
  sendKeyEvent(18, false) // Alt up
  sendKeyEvent(17, false) // Ctrl up
  ElMessage.success('已发送 Ctrl+Alt+Del')
}

// 切换全屏
function toggleFullscreen() {
  if (!vncContainer.value) return

  if (!document.fullscreenElement) {
    vncContainer.value.requestFullscreen().then(() => {
      isFullscreen.value = true
    }).catch(() => {
      ElMessage.warning('无法进入全屏模式')
    })
  } else {
    document.exitFullscreen().then(() => {
      isFullscreen.value = false
    })
  }
}

// 监听全屏变化
function handleFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement
}

// 开始心跳
function startPing() {
  pingTimer = setInterval(() => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'ping' }))
    }
  }, 30000)
}

// 开始延迟监控
function startLatencyMonitor() {
  let lastPing = Date.now()
  latencyTimer = setInterval(() => {
    lastPing = Date.now()
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'latency', ts: lastPing }))
    }
  }, 5000)
}

// 断开连接
function disconnect() {
  if (ws) {
    ws.close()
    ws = null
  }
  if (pingTimer) {
    clearInterval(pingTimer)
    pingTimer = null
  }
  if (latencyTimer) {
    clearInterval(latencyTimer)
    latencyTimer = null
  }
  connected.value = false
  connecting.value = false

  // 移除事件监听
  if (vncCanvas.value) {
    vncCanvas.value.removeEventListener('mousedown', handleMouseDown)
    vncCanvas.value.removeEventListener('mouseup', handleMouseUp)
    vncCanvas.value.removeEventListener('mousemove', handleMouseMove)
    vncCanvas.value.removeEventListener('wheel', handleWheel)
    vncCanvas.value.removeEventListener('keydown', handleKeyDown)
    vncCanvas.value.removeEventListener('keyup', handleKeyUp)
  }
}

onMounted(() => {
  fetchServerInfo()
  document.addEventListener('fullscreenchange', handleFullscreenChange)
})

onBeforeUnmount(() => {
  disconnect()
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
})
</script>

<style scoped lang="scss">
.vnc-console-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  gap: 16px;
}

.page-header {
  :deep(.el-page-header) {
    width: 100%;
  }

  :deep(.el-page-header__content) {
    flex: 1;
  }

  :deep(.el-page-header__extra) {
    flex: none;
  }
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.info-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  flex: 1;

  :deep(.el-card__body) {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.connect-info {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;

  h3 {
    font-size: 20px;
    font-weight: 600;
    color: #303133;
    margin: 8px 0 0;
  }

  p {
    font-size: 14px;
    color: #909399;
    margin: 0;
  }
}

.server-info {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px 24px;
  margin: 8px 0;
  text-align: left;
  min-width: 300px;
}

.info-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;

  .label {
    font-size: 13px;
    color: #909399;
    width: 80px;
    flex-shrink: 0;
  }

  .value {
    font-size: 14px;
    color: #303133;
    font-weight: 500;

    &.mono {
      font-family: 'Monaco', 'Menlo', monospace;
    }
  }
}

.connecting-overlay {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  background: #f5f7fa;
  border-radius: 12px;

  p {
    font-size: 14px;
    color: #606266;
  }
}

.loading-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.vnc-container {
  flex: 1;
  background: #000;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  cursor: default;
}

.vnc-canvas {
  width: 100%;
  height: 100%;
  display: block;
}

.status-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e8ecf1;
  font-size: 12px;
  color: #606266;
}

.status-left,
.status-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.separator {
  color: #dcdfe6;
}

@media (max-width: 768px) {
  .header-actions {
    flex-wrap: wrap;
  }

  .server-info {
    min-width: auto;
    width: 100%;
  }

  .status-bar {
    flex-direction: column;
    gap: 4px;
  }
}
</style>
