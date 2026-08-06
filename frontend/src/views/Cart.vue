<template>
  <div class="cart-page">
    <SiteHeader />
    
    <div class="cart-container">
      <!-- 购物车头部 -->
      <div class="cart-header">
        <h1>
          <el-icon><ShoppingCart /></el-icon>
          购物车
        </h1>
        <el-button text @click="$router.push('/products')">
          <el-icon><ArrowLeft /></el-icon>
          继续购物
        </el-button>
      </div>
      
      <!-- 购物车内容 -->
      <div class="cart-content" v-if="cartItems.length > 0">
        <!-- 购物车列表 -->
        <div class="cart-list">
          <div class="cart-item" v-for="(item, index) in cartItems" :key="item.id">
            <!-- 产品图标 -->
            <div class="item-icon">
              <el-icon :size="32"><component :is="getProductIcon(item.type)" /></el-icon>
            </div>
            
            <!-- 产品信息 -->
            <div class="item-info">
              <h3 class="item-name">{{ item.product_name }}</h3>
              <div class="item-specs">
                <span class="spec-tag" v-for="(spec, i) in item.specs" :key="i">{{ spec }}</span>
              </div>
              <div class="item-config">
                <span v-if="item.config.os">系统：{{ item.config.os }}</span>
                <span v-if="item.config.area">区域：{{ item.config.area }}</span>
              </div>
            </div>
            
            <!-- 计费周期 -->
            <div class="item-cycle">
              <el-select v-model="item.cycle" size="small" @change="updateItemCycle(item)">
                <el-option
                  v-for="cycle in item.cycles"
                  :key="cycle.id"
                  :label="cycle.name"
                  :value="cycle.id"
                />
              </el-select>
            </div>
            
            <!-- 数量 -->
            <div class="item-quantity">
              <el-input-number 
                v-model="item.quantity" 
                :min="1" 
                :max="99" 
                size="small"
                @change="updateCartItem(item)"
              />
            </div>
            
            <!-- 价格 -->
            <div class="item-price">
              <div class="price-amount">¥{{ (item.price * item.quantity).toFixed(2) }}</div>
              <div class="price-unit" v-if="item.quantity > 1">单价 ¥{{ item.price.toFixed(2) }}</div>
            </div>
            
            <!-- 删除 -->
            <div class="item-actions">
              <el-button type="danger" text circle @click="removeItem(index)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
        
        <!-- 购物车底部 -->
        <div class="cart-footer">
          <div class="footer-left">
            <el-checkbox v-model="selectAll" @change="handleSelectAll">全选</el-checkbox>
            <el-button text type="danger" @click="removeSelected" :disabled="!hasSelected">
              删除选中
            </el-button>
          </div>
          
          <div class="footer-right">
            <div class="summary">
              <div class="summary-item">
                <span>已选商品</span>
                <span class="value">{{ selectedCount }} 件</span>
              </div>
              <div class="summary-item">
                <span>合计</span>
                <span class="total-price">¥{{ totalPrice.toFixed(2) }}</span>
              </div>
            </div>
            <el-button 
              type="primary" 
              size="large" 
              :disabled="!hasSelected"
              @click="goToCheckout"
            >
              去结算
            </el-button>
          </div>
        </div>
      </div>
      
      <!-- 空购物车 -->
      <div class="cart-empty" v-else>
        <el-empty description="购物车是空的">
          <el-button type="primary" @click="$router.push('/products')">去选购产品</el-button>
        </el-empty>
      </div>
    </div>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  ShoppingCart, ArrowLeft, Delete
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()

interface CartItem {
  id: number
  product_id: number
  product_name: string
  type: string
  specs: string[]
  config: {
    os?: string
    area?: string
    [key: string]: any
  }
  cycles: Array<{ id: string; name: string; price: number }>
  cycle: string
  price: number
  quantity: number
  selected: boolean
}

const cartItems = ref<CartItem[]>([])
const selectAll = ref(true)
const loading = ref(false)

// 计算属性
const hasSelected = computed(() => cartItems.value.some(item => item.selected))
const selectedCount = computed(() => cartItems.value.filter(item => item.selected).length)
const totalPrice = computed(() => {
  return cartItems.value
    .filter(item => item.selected)
    .reduce((sum, item) => sum + item.price * item.quantity, 0)
})

// 获取产品图标
const getProductIcon = (type: string) => {
  const iconMap: Record<string, any> = {
    cloud: 'Monitor',
    vps: 'Monitor',
    server: 'Monitor',
    domain: 'Connection',
    ssl: 'Ticket'
  }
  return iconMap[type] || 'Setting'
}

// 获取购物车数据
const fetchCart = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/cart')
    if (data?.data) {
      cartItems.value = data.data.map((item: any) => ({
        ...item,
        selected: true,
        specs: item.specs || [],
        config: item.config || {},
        cycles: item.cycles || []
      }))
    }
  } catch (error) {
    console.error('获取购物车失败:', error)
  } finally {
    loading.value = false
  }
}

// 更新购物车项
const updateCartItem = async (item: CartItem) => {
  try {
    await request.put(`/api/v1/cart/${item.id}`, {
      quantity: item.quantity,
      cycle: item.cycle
    })
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

// 更新计费周期
const updateItemCycle = (item: CartItem) => {
  const cycle = item.cycles.find(c => c.id === item.cycle)
  if (cycle) {
    item.price = cycle.price
    updateCartItem(item)
  }
}

// 删除单个商品
const removeItem = async (index: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个商品吗？', '提示', {
      type: 'warning'
    })
    
    const item = cartItems.value[index]
    await request.delete(`/api/v1/cart/${item.id}`)
    cartItems.value.splice(index, 1)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 全选/取消全选
const handleSelectAll = (val: boolean) => {
  cartItems.value.forEach(item => {
    item.selected = val
  })
}

// 删除选中商品
const removeSelected = async () => {
  try {
    await ElMessageBox.confirm('确定要删除选中的商品吗？', '提示', {
      type: 'warning'
    })
    
    const selectedIds = cartItems.value
      .filter(item => item.selected)
      .map(item => item.id)
    
    await request.post('/api/v1/cart/batch-delete', { ids: selectedIds })
    cartItems.value = cartItems.value.filter(item => !item.selected)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 去结算
const goToCheckout = () => {
  const selectedItems = cartItems.value.filter(item => item.selected)
  if (selectedItems.length === 0) {
    ElMessage.warning('请选择要结算的商品')
    return
  }
  
  // 将选中的商品ID存入本地存储
  localStorage.setItem('checkout_items', JSON.stringify(selectedItems.map(item => item.id)))
  router.push('/checkout')
}

onMounted(() => {
  fetchCart()
})
</script>

<style scoped lang="scss">
.cart-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding-top: 64px;
  
  .cart-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px 20px;
  }
  
  // 购物车头部
  .cart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    
    h1 {
      display: flex;
      align-items: center;
      gap: 10px;
      font-size: 24px;
      font-weight: 600;
      color: #1a2332;
      margin: 0;
      
      .el-icon {
        color: #409eff;
      }
    }
  }
  
  // 购物车内容
  .cart-content {
    background: #fff;
    border-radius: 12px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
    overflow: hidden;
  }
  
  // 购物车列表
  .cart-list {
    .cart-item {
      display: flex;
      align-items: center;
      padding: 24px;
      border-bottom: 1px solid #f0f0f0;
      transition: background 0.2s;
      
      &:hover {
        background: #fafafa;
      }
      
      &:last-child {
        border-bottom: none;
      }
      
      .item-icon {
        width: 56px;
        height: 56px;
        background: linear-gradient(135deg, #e3f2fd, #bbdefb);
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #1976d2;
        margin-right: 16px;
        flex-shrink: 0;
      }
      
      .item-info {
        flex: 1;
        min-width: 0;
        
        .item-name {
          font-size: 16px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 8px;
        }
        
        .item-specs {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
          margin-bottom: 6px;
          
          .spec-tag {
            padding: 2px 8px;
            background: #f5f7fa;
            border-radius: 4px;
            font-size: 12px;
            color: #666;
          }
        }
        
        .item-config {
          font-size: 13px;
          color: #909399;
          display: flex;
          gap: 16px;
        }
      }
      
      .item-cycle {
        width: 120px;
        margin: 0 16px;
      }
      
      .item-quantity {
        margin: 0 16px;
      }
      
      .item-price {
        width: 120px;
        text-align: right;
        margin-right: 16px;
        
        .price-amount {
          font-size: 18px;
          font-weight: 600;
          color: #f56c6c;
        }
        
        .price-unit {
          font-size: 12px;
          color: #909399;
        }
      }
      
      .item-actions {
        flex-shrink: 0;
      }
    }
  }
  
  // 购物车底部
  .cart-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px;
    background: #fff;
    border-top: 1px solid #f0f0f0;
    position: sticky;
    bottom: 0;
    z-index: 10;
    
    .footer-left {
      display: flex;
      align-items: center;
      gap: 16px;
    }
    
    .footer-right {
      display: flex;
      align-items: center;
      gap: 24px;
      
      .summary {
        text-align: right;
        
        .summary-item {
          font-size: 14px;
          color: #666;
          margin-bottom: 4px;
          
          .value {
            margin-left: 8px;
          }
          
          .total-price {
            font-size: 24px;
            font-weight: 700;
            color: #f56c6c;
            margin-left: 8px;
          }
        }
      }
    }
  }
  
  // 空购物车
  .cart-empty {
    background: #fff;
    border-radius: 12px;
    padding: 80px 0;
    text-align: center;
  }
}

// 响应式
@media (max-width: 768px) {
  .cart-page {
    .cart-item {
      flex-wrap: wrap;
      gap: 12px;
      
      .item-cycle,
      .item-quantity {
        margin: 0;
      }
      
      .item-price {
        margin-left: auto;
      }
    }
    
    .cart-footer {
      flex-direction: column;
      gap: 16px;
      
      .footer-right {
        width: 100%;
        justify-content: space-between;
      }
    }
  }
}
</style>
