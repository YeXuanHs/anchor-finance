<template>
  <div class="cart-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <img src="/logo.png" alt="锚点财务" class="logo-img" />
          <span class="logo-text">锚点财务</span>
        </router-link>
        <h2 class="page-title">购物车</h2>
        <div class="header-actions">
          <el-button text @click="$router.push('/products')">继续选购</el-button>
        </div>
      </div>
    </header>

    <div class="cart-content">
      <div class="cart-inner">
        <div class="cart-grid" v-loading="loading">
          <!-- 左侧购物车列表 -->
          <div class="cart-list">
            <div class="empty-cart" v-if="!cartItems.length">
              <el-icon :size="64" color="#e5e5ea"><ShoppingCart /></el-icon>
              <p>购物车是空的</p>
              <el-button type="primary" @click="$router.push('/products')">去选购</el-button>
            </div>

            <div v-else>
              <div class="cart-header">
                <el-checkbox v-model="selectAll" @change="handleSelectAll">全选</el-checkbox>
                <span class="items-count">共 {{ cartItems.length }} 件商品</span>
              </div>

              <div v-for="item in cartItems" :key="item.id" class="cart-item">
                <el-checkbox v-model="item.selected" @change="updateSelection" />
                
                <div class="item-icon" :style="{ background: item.gradient || 'linear-gradient(135deg, #1a73e8, #4a90e2)' }">
                  <el-icon :size="24" color="#fff"><Monitor /></el-icon>
                </div>
                
                <div class="item-info">
                  <h3>{{ item.product_name }}</h3>
                  <p class="item-specs">
                    <span v-if="item.region">{{ item.region }}</span>
                    <span v-if="item.os">{{ item.os }}</span>
                    <span v-if="item.cpu">{{ item.cpu }}</span>
                    <span v-if="item.memory">{{ item.memory }}</span>
                  </p>
                  <p class="item-cycle">计费周期：{{ item.billing_cycle }}</p>
                </div>
                
                <div class="item-price">
                  <span class="price-symbol">¥</span>
                  <span class="price-value">{{ item.price?.toFixed(2) }}</span>
                  <span class="price-unit">/{{ item.billing_cycle === 'monthly' ? '月' : '年' }}</span>
                </div>
                
                <el-button type="danger" text @click="removeItem(item)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
          </div>

          <!-- 右侧结算区 -->
          <div class="cart-summary" v-if="cartItems.length">
            <div class="summary-card">
              <h3>订单摘要</h3>
              
              <div class="summary-item">
                <span>商品合计</span>
                <span>¥{{ subtotal.toFixed(2) }}</span>
              </div>
              
              <div class="summary-item" v-if="discount > 0">
                <span>优惠金额</span>
                <span class="discount">-¥{{ discount.toFixed(2) }}</span>
              </div>
              
              <!-- 优惠码 -->
              <div class="coupon-input">
                <el-input v-model="couponCode" placeholder="请输入优惠码" size="small">
                  <template #append>
                    <el-button @click="applyCoupon" :loading="couponLoading">应用</el-button>
                  </template>
                </el-input>
              </div>
              
              <el-divider />
              
              <div class="summary-item total">
                <span>应付金额</span>
                <span class="total-price">¥{{ total.toFixed(2) }}</span>
              </div>
              
              <el-button type="primary" size="large" round class="checkout-btn" @click="checkout">
                去结算 ({{ selectedCount }})
              </el-button>
              
              <div class="summary-tips">
                <p><el-icon><CircleCheck /></el-icon> 支持多种支付方式</p>
                <p><el-icon><CircleCheck /></el-icon> 7天无理由退款</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ShoppingCart, Monitor, Delete, CircleCheck } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const cartItems = ref([])
const selectAll = ref(false)
const couponCode = ref('')
const couponLoading = ref(false)
const discount = ref(0)

// 获取购物车数据
const fetchCart = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/cart')
    if (data?.data) {
      cartItems.value = data.data.map((item: any) => ({
        ...item,
        selected: true
      }))
    }
  } catch (error) {
    console.error('获取购物车失败:', error)
  } finally {
    loading.value = false
  }
}

// 计算属性
const subtotal = computed(() => {
  return cartItems.value
    .filter(item => item.selected)
    .reduce((sum, item) => sum + (item.price || 0), 0)
})

const total = computed(() => Math.max(0, subtotal.value - discount.value))

const selectedCount = computed(() => cartItems.value.filter(item => item.selected).length)

// 全选/取消全选
const handleSelectAll = (val: boolean) => {
  cartItems.value.forEach(item => {
    item.selected = val
  })
}

const updateSelection = () => {
  selectAll.value = cartItems.value.every(item => item.selected)
}

// 移除商品
const removeItem = async (item: any) => {
  try {
    await ElMessageBox.confirm('确定要移除该商品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.delete(`/api/v1/cart/${item.id}`)
    cartItems.value = cartItems.value.filter(i => i.id !== item.id)
    ElMessage.success('已移除')
  } catch (error) {
    // 取消操作
  }
}

// 应用优惠码
const applyCoupon = async () => {
  if (!couponCode.value) return
  
  couponLoading.value = true
  try {
    const { data } = await request.post('/api/v1/coupons/verify', {
      code: couponCode.value
    })
    if (data?.ok) {
      discount.value = data.discount || 0
      ElMessage.success('优惠码已应用')
    } else {
      ElMessage.error(data?.message || '优惠码无效')
    }
  } catch (error) {
    ElMessage.error('验证优惠码失败')
  } finally {
    couponLoading.value = false
  }
}

// 去结算
const checkout = () => {
  const selectedItems = cartItems.value.filter(item => item.selected)
  if (selectedItems.length === 0) {
    ElMessage.warning('请选择要结算的商品')
    return
  }
  
  // 将选中的商品ID存入本地存储
  localStorage.setItem('checkout_items', JSON.stringify(selectedItems.map(i => i.id)))
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
}

.page-header {
  background: #fff;
  border-bottom: 1px solid #e5e5ea;
  
  .header-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    text-decoration: none;
    
    .logo-img {
      width: 32px;
      height: 32px;
    }
    
    .logo-text {
      font-size: 18px;
      font-weight: 600;
      color: #1d1d1f;
    }
  }
  
  .page-title {
    font-size: 18px;
    font-weight: 600;
    margin: 0;
  }
}

.cart-content {
  padding: 24px 0;
}

.cart-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.cart-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
  
  @media (max-width: 992px) {
    grid-template-columns: 1fr;
  }
}

.cart-list {
  .empty-cart {
    background: #fff;
    border-radius: 12px;
    padding: 80px 0;
    text-align: center;
    
    p {
      margin: 16px 0;
      color: #909399;
    }
  }
  
  .cart-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    
    .items-count {
      color: #909399;
      font-size: 14px;
    }
  }
  
  .cart-item {
    background: #fff;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    gap: 16px;
    border: 1px solid #e5e5ea;
    transition: all 0.2s;
    
    &:hover {
      border-color: #1a73e8;
    }
    
    .item-icon {
      width: 56px;
      height: 56px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }
    
    .item-info {
      flex: 1;
      
      h3 {
        font-size: 16px;
        font-weight: 600;
        margin: 0 0 8px;
      }
      
      .item-specs {
        display: flex;
        gap: 12px;
        font-size: 13px;
        color: #909399;
        margin: 0 0 4px;
      }
      
      .item-cycle {
        font-size: 13px;
        color: #909399;
        margin: 0;
      }
    }
    
    .item-price {
      text-align: right;
      flex-shrink: 0;
      
      .price-symbol {
        font-size: 14px;
        color: #1a73e8;
      }
      
      .price-value {
        font-size: 24px;
        font-weight: 700;
        color: #1a73e8;
      }
      
      .price-unit {
        font-size: 14px;
        color: #909399;
      }
    }
  }
}

.cart-summary {
  .summary-card {
    background: #fff;
    border-radius: 12px;
    padding: 24px;
    border: 1px solid #e5e5ea;
    position: sticky;
    top: 88px;
    
    h3 {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 20px;
    }
    
    .summary-item {
      display: flex;
      justify-content: space-between;
      padding: 8px 0;
      font-size: 14px;
      color: #606266;
      
      .discount {
        color: #67c23a;
      }
      
      &.total {
        font-weight: 600;
        color: #1d1d1f;
        
        .total-price {
          font-size: 24px;
          color: #1a73e8;
        }
      }
    }
    
    .coupon-input {
      margin: 16px 0;
    }
    
    .checkout-btn {
      width: 100%;
      height: 48px;
      font-size: 16px;
      margin-top: 16px;
    }
    
    .summary-tips {
      margin-top: 20px;
      
      p {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        color: #909399;
        margin: 8px 0;
      }
    }
  }
}
</style>
