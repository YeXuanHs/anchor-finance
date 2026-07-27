<template>
  <div class="checkout-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <img src="/logo.png" alt="锚点财务" class="logo-img" />
          <span class="logo-text">锚点财务</span>
        </router-link>
        <h2 class="page-title">确认订单</h2>
        <div class="header-actions">
          <el-button text @click="$router.push('/cart')">返回购物车</el-button>
        </div>
      </div>
    </header>

    <div class="checkout-content" v-loading="loading">
      <div class="checkout-inner">
        <div class="checkout-grid">
          <!-- 左侧订单信息 -->
          <div class="order-section">
            <!-- 订单商品 -->
            <div class="order-card">
              <h3>订单商品</h3>
              <div v-for="item in orderItems" :key="item.id" class="order-item">
                <div class="item-icon" :style="{ background: item.gradient || 'linear-gradient(135deg, #1a73e8, #4a90e2)' }">
                  <el-icon :size="24" color="#fff"><Monitor /></el-icon>
                </div>
                <div class="item-info">
                  <h4>{{ item.product_name }}</h4>
                  <p class="item-specs">
                    <span v-if="item.region">{{ item.region }}</span>
                    <span v-if="item.os">{{ item.os }}</span>
                    <span v-if="item.cpu">{{ item.cpu }}</span>
                    <span v-if="item.memory">{{ item.memory }}</span>
                    <span v-if="item.disk">{{ item.disk }}</span>
                    <span v-if="item.bandwidth">{{ item.bandwidth }}</span>
                  </p>
                  <p class="item-cycle">计费周期：{{ item.billing_cycle_name }}</p>
                </div>
                <div class="item-price">
                  <span class="price-symbol">¥</span>
                  <span class="price-value">{{ item.price?.toFixed(2) }}</span>
                </div>
              </div>
            </div>

            <!-- 支付方式 -->
            <div class="order-card">
              <h3>支付方式</h3>
              <div class="payment-methods">
                <div
                  v-for="method in paymentMethods"
                  :key="method.id"
                  class="payment-item"
                  :class="{ active: selectedPayment === method.id }"
                  @click="selectedPayment = method.id"
                >
                  <img :src="method.icon" :alt="method.name" class="payment-icon" />
                  <span class="payment-name">{{ method.name }}</span>
                  <el-icon v-if="selectedPayment === method.id" class="check-icon" color="#1a73e8"><CircleCheckFilled /></el-icon>
                </div>
              </div>
            </div>

            <!-- 优惠码 -->
            <div class="order-card">
              <h3>优惠码</h3>
              <div class="coupon-input">
                <el-input v-model="couponCode" placeholder="请输入优惠码">
                  <template #append>
                    <el-button @click="applyCoupon" :loading="couponLoading">应用</el-button>
                  </template>
                </el-input>
              </div>
            </div>
          </div>

          <!-- 右侧订单摘要 -->
          <div class="summary-section">
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
              
              <div class="summary-item" v-if="balanceUsed > 0">
                <span>余额抵扣</span>
                <span class="discount">-¥{{ balanceUsed.toFixed(2) }}</span>
              </div>
              
              <el-divider />
              
              <div class="summary-item total">
                <span>应付金额</span>
                <span class="total-price">¥{{ finalAmount.toFixed(2) }}</span>
              </div>
              
              <!-- 余额 -->
              <div class="balance-info" v-if="userBalance > 0">
                <el-checkbox v-model="useBalance">
                  使用余额 (¥{{ userBalance.toFixed(2) }})
                </el-checkbox>
              </div>
              
              <el-button 
                type="primary" 
                size="large" 
                round 
                class="pay-btn"
                :loading="payLoading"
                @click="submitOrder"
              >
                确认支付 ¥{{ finalAmount.toFixed(2) }}
              </el-button>
              
              <div class="summary-tips">
                <p>点击"确认支付"即表示您同意我们的服务条款</p>
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
import { Monitor, CircleCheckFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const payLoading = ref(false)

const orderItems = ref([])
const paymentMethods = ref([])
const userBalance = ref(0)

const selectedPayment = ref('')
const couponCode = ref('')
const couponLoading = ref(false)
const discount = ref(0)
const useBalance = ref(false)

// 获取订单数据
const fetchOrderData = async () => {
  loading.value = true
  try {
    // 获取购物车选中的商品
    const checkoutItems = JSON.parse(localStorage.getItem('checkout_items') || '[]')
    
    const { data } = await request.post('/api/v1/order/preview', {
      items: checkoutItems
    })
    
    if (data?.data) {
      orderItems.value = data.data.items || []
      paymentMethods.value = data.data.payment_methods || []
      userBalance.value = data.data.balance || 0
      
      if (paymentMethods.value.length) {
        selectedPayment.value = paymentMethods.value[0].id
      }
    }
  } catch (error) {
    console.error('获取订单数据失败:', error)
    ElMessage.error('获取订单数据失败')
  } finally {
    loading.value = false
  }
}

// 计算属性
const subtotal = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + (item.price || 0), 0)
})

const balanceUsed = computed(() => {
  if (!useBalance.value) return 0
  const remaining = subtotal.value - discount.value
  return Math.min(userBalance.value, remaining)
})

const finalAmount = computed(() => {
  return Math.max(0, subtotal.value - discount.value - balanceUsed.value)
})

// 应用优惠码
const applyCoupon = async () => {
  if (!couponCode.value) return
  
  couponLoading.value = true
  try {
    const { data } = await request.post('/api/v1/coupons/verify', {
      code: couponCode.value,
      items: orderItems.value.map(i => i.id)
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

// 提交订单
const submitOrder = async () => {
  if (!selectedPayment.value) {
    ElMessage.warning('请选择支付方式')
    return
  }
  
  payLoading.value = true
  try {
    const checkoutItems = JSON.parse(localStorage.getItem('checkout_items') || '[]')
    
    const { data } = await request.post('/api/v1/orders', {
      items: checkoutItems,
      payment_method: selectedPayment.value,
      coupon_code: couponCode.value,
      use_balance: useBalance.value
    })
    
    if (data?.ok) {
      // 清除购物车选中记录
      localStorage.removeItem('checkout_items')
      
      if (data.payment_url) {
        // 跳转到支付页面
        window.location.href = data.payment_url
      } else {
        // 余额支付成功
        ElMessage.success('订单创建成功')
        router.push(`/payment-result?order=${data.order_no}`)
      }
    } else {
      ElMessage.error(data?.message || '创建订单失败')
    }
  } catch (error) {
    ElMessage.error('创建订单失败')
  } finally {
    payLoading.value = false
  }
}

onMounted(() => {
  fetchOrderData()
})
</script>

<style scoped lang="scss">
.checkout-page {
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

.checkout-content {
  padding: 24px 0;
}

.checkout-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.checkout-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
  
  @media (max-width: 992px) {
    grid-template-columns: 1fr;
  }
}

.order-section {
  .order-card {
    background: #fff;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 16px;
    border: 1px solid #e5e5ea;
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 16px;
    }
  }
  
  .order-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 0;
    border-bottom: 1px solid #f0f0f0;
    
    &:last-child {
      border-bottom: none;
    }
    
    .item-icon {
      width: 48px;
      height: 48px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }
    
    .item-info {
      flex: 1;
      
      h4 {
        font-size: 15px;
        font-weight: 600;
        margin: 0 0 6px;
      }
      
      .item-specs {
        display: flex;
        gap: 10px;
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
        font-size: 20px;
        font-weight: 700;
        color: #1a73e8;
      }
    }
  }
}

.payment-methods {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 12px;
  
  .payment-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 14px 16px;
    border: 1px solid #e5e5ea;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    position: relative;
    
    &:hover {
      border-color: #1a73e8;
    }
    
    &.active {
      border-color: #1a73e8;
      background: rgba(26, 115, 232, 0.05);
    }
    
    .payment-icon {
      width: 24px;
      height: 24px;
      object-fit: contain;
    }
    
    .payment-name {
      font-size: 14px;
      font-weight: 500;
    }
    
    .check-icon {
      position: absolute;
      top: 6px;
      right: 6px;
      font-size: 16px;
    }
  }
}

.coupon-input {
  margin-top: 12px;
}

.summary-section {
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
    
    .balance-info {
      margin: 16px 0;
      padding: 12px;
      background: #f5f7fa;
      border-radius: 8px;
    }
    
    .pay-btn {
      width: 100%;
      height: 48px;
      font-size: 16px;
      margin-top: 16px;
    }
    
    .summary-tips {
      margin-top: 16px;
      
      p {
        font-size: 12px;
        color: #909399;
        text-align: center;
        margin: 0;
      }
    }
  }
}
</style>
