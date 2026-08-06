<template>
  <div class="checkout-page">
    <SiteHeader />
    
    <div class="checkout-container">
      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/cart' }">购物车</el-breadcrumb-item>
          <el-breadcrumb-item>确认订单</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      
      <!-- 结算内容 -->
      <div class="checkout-content">
        <!-- 左侧：订单信息 -->
        <div class="checkout-main">
          <!-- 收货信息（域名等产品需要） -->
          <div class="section-card" v-if="needContact">
            <div class="section-header">
              <h3>联系信息</h3>
            </div>
            <div class="section-body">
              <el-form :model="contactForm" label-width="80px">
                <el-form-item label="姓名">
                  <el-input v-model="contactForm.name" placeholder="请输入姓名" />
                </el-form-item>
                <el-form-item label="邮箱">
                  <el-input v-model="contactForm.email" placeholder="请输入邮箱" />
                </el-form-item>
                <el-form-item label="手机">
                  <el-input v-model="contactForm.phone" placeholder="请输入手机号" />
                </el-form-item>
              </el-form>
            </div>
          </div>
          
          <!-- 订单商品 -->
          <div class="section-card">
            <div class="section-header">
              <h3>订单商品</h3>
              <el-button text @click="$router.push('/cart')">返回购物车修改</el-button>
            </div>
            <div class="section-body">
              <div class="order-item" v-for="item in orderItems" :key="item.id">
                <div class="item-icon">
                  <el-icon :size="24"><component :is="getProductIcon(item.type)" /></el-icon>
                </div>
                <div class="item-info">
                  <h4>{{ item.product_name }}</h4>
                  <div class="item-meta">
                    <span v-if="item.config.os">系统：{{ item.config.os }}</span>
                    <span v-if="item.config.area">区域：{{ item.config.area }}</span>
                    <span>周期：{{ item.cycle_name }}</span>
                  </div>
                </div>
                <div class="item-quantity">x{{ item.quantity }}</div>
                <div class="item-price">¥{{ (item.price * item.quantity).toFixed(2) }}</div>
              </div>
            </div>
          </div>
          
          <!-- 优惠码 -->
          <div class="section-card">
            <div class="section-header">
              <h3>优惠码</h3>
            </div>
            <div class="section-body">
              <div class="coupon-input">
                <el-input v-model="couponCode" placeholder="请输入优惠码" :disabled="couponApplied">
                  <template #append>
                    <el-button @click="applyCoupon" :loading="couponLoading">
                      {{ couponApplied ? '已应用' : '应用' }}
                    </el-button>
                  </template>
                </el-input>
                <div class="coupon-info" v-if="couponApplied">
                  <el-tag type="success" closable @close="removeCoupon">
                    优惠 ¥{{ couponDiscount.toFixed(2) }}
                  </el-tag>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 支付方式 -->
          <div class="section-card">
            <div class="section-header">
              <h3>支付方式</h3>
            </div>
            <div class="section-body">
              <div class="payment-methods">
                <div 
                  class="payment-item" 
                  v-for="method in paymentMethods" 
                  :key="method.id"
                  :class="{ active: selectedPayment === method.id }"
                  @click="selectedPayment = method.id"
                >
                  <div class="payment-icon">
                    <img :src="method.icon" :alt="method.name" />
                  </div>
                  <span class="payment-name">{{ method.name }}</span>
                  <el-icon v-if="selectedPayment === method.id" class="check-icon">
                    <CircleCheck />
                  </el-icon>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 备注 -->
          <div class="section-card">
            <div class="section-header">
              <h3>订单备注</h3>
            </div>
            <div class="section-body">
              <el-input 
                v-model="remark" 
                type="textarea" 
                :rows="3" 
                placeholder="请输入备注信息（选填）"
              />
            </div>
          </div>
        </div>
        
        <!-- 右侧：订单摘要 -->
        <div class="checkout-sidebar">
          <div class="order-summary">
            <h3>订单摘要</h3>
            
            <div class="summary-list">
              <div class="summary-item">
                <span>商品金额</span>
                <span>¥{{ subtotal.toFixed(2) }}</span>
              </div>
              <div class="summary-item" v-if="couponDiscount > 0">
                <span>优惠金额</span>
                <span class="discount">-¥{{ couponDiscount.toFixed(2) }}</span>
              </div>
              <div class="summary-item">
                <span>应付金额</span>
                <span class="total">¥{{ totalAmount.toFixed(2) }}</span>
              </div>
            </div>
            
            <!-- 余额支付 -->
            <div class="balance-pay" v-if="userBalance > 0">
              <el-checkbox v-model="useBalance">
                使用余额支付（可用 ¥{{ userBalance.toFixed(2) }}）
              </el-checkbox>
              <div class="balance-amount" v-if="useBalance">
                <span>余额抵扣</span>
                <span class="discount">-¥{{ balanceDeduct.toFixed(2) }}</span>
              </div>
            </div>
            
            <!-- 实际支付 -->
            <div class="actual-pay" v-if="useBalance && balanceDeduct > 0">
              <div class="pay-item">
                <span>还需支付</span>
                <span class="amount">¥{{ actualPay.toFixed(2) }}</span>
              </div>
            </div>
            
            <el-button 
              type="primary" 
              size="large" 
              class="submit-btn"
              :loading="submitting"
              @click="submitOrder"
            >
              提交订单
            </el-button>
            
            <div class="agreement">
              <el-checkbox v-model="agreeTerms">
                我已阅读并同意
                <router-link to="/terms">《服务条款》</router-link>
              </el-checkbox>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 银行转账对话框 -->
    <el-dialog v-model="bankDialogVisible" title="银行转账信息" width="500px" :close-on-click-modal="false">
      <div class="bank-transfer-info" v-if="bankInfo">
        <el-alert type="warning" :closable="false" show-icon>
          <template #title>
            请在24小时内完成转账，否则订单将自动取消
          </template>
        </el-alert>
        
        <div class="bank-info-card">
          <div class="info-row">
            <span class="label">收款银行</span>
            <span class="value">{{ bankInfo.bank_name }}</span>
          </div>
          <div class="info-row">
            <span class="label">开户行</span>
            <span class="value">{{ bankInfo.branch_name }}</span>
          </div>
          <div class="info-row">
            <span class="label">账户名</span>
            <span class="value">{{ bankInfo.account_name }}</span>
          </div>
          <div class="info-row">
            <span class="label">账号</span>
            <span class="value">
              {{ bankInfo.account_no }}
              <el-button type="primary" link @click="copyText(bankInfo.account_no)">复制</el-button>
            </span>
          </div>
          <div class="info-row">
            <span class="label">转账金额</span>
            <span class="value amount">¥{{ bankInfo.amount?.toFixed(2) }}</span>
          </div>
          <div class="info-row" v-if="bankInfo.order_no">
            <span class="label">订单号</span>
            <span class="value">
              {{ bankInfo.order_no }}
              <el-button type="primary" link @click="copyText(bankInfo.order_no)">复制</el-button>
            </span>
          </div>
        </div>
        
        <div class="bank-instructions" v-if="bankInfo.instructions">
          <h4>转账说明</h4>
          <p>{{ bankInfo.instructions }}</p>
        </div>
        
        <div class="bank-qrcode" v-if="bankInfo.qrcode">
          <h4>扫码支付</h4>
          <img :src="bankInfo.qrcode" alt="收款二维码" class="qrcode-img" />
        </div>
      </div>
      <template #footer>
        <el-button @click="bankDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="router.push('/user/orders')">查看订单</el-button>
      </template>
    </el-dialog>
    
    <SiteFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  CircleCheck
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'
import request from '@/utils/request'

const router = useRouter()

// 银行转账相关
const bankDialogVisible = ref(false)
const bankInfo = ref<any>(null)

// 复制文本
const copyText = (text: string) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 表单数据
const contactForm = ref({
  name: '',
  email: '',
  phone: ''
})

const couponCode = ref('')
const couponApplied = ref(false)
const couponLoading = ref(false)
const couponDiscount = ref(0)
const selectedPayment = ref('')
const remark = ref('')
const useBalance = ref(false)
const userBalance = ref(0)
const agreeTerms = ref(false)
const submitting = ref(false)

// 订单商品
const orderItems = ref<any[]>([])

// 支付方式
const paymentMethods = ref([
  { id: 'alipay', name: '支付宝', icon: '/assets/payment/alipay.png' },
  { id: 'wechat', name: '微信支付', icon: '/assets/payment/wechat.png' },
  { id: 'qqpay', name: 'QQ钱包', icon: '/assets/payment/qqpay.png' },
  { id: 'usdt', name: 'USDT', icon: '/assets/payment/usdt.png' }
])

// 计算属性
const needContact = computed(() => {
  return orderItems.value.some(item => ['domain', 'ssl'].includes(item.type))
})

const subtotal = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
})

const totalAmount = computed(() => {
  return Math.max(0, subtotal.value - couponDiscount.value)
})

const balanceDeduct = computed(() => {
  if (!useBalance.value) return 0
  return Math.min(userBalance.value, totalAmount.value)
})

const actualPay = computed(() => {
  return Math.max(0, totalAmount.value - balanceDeduct.value)
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

// 获取订单数据
const fetchOrderData = async () => {
  try {
    const checkoutItems = JSON.parse(localStorage.getItem('checkout_items') || '[]')
    
    if (checkoutItems.length === 0) {
      ElMessage.warning('请先选择商品')
      router.push('/cart')
      return
    }
    
    const { data } = await request.post('/api/v1/orders/preview', {
      items: checkoutItems
    })
    
    if (data?.data) {
      orderItems.value = data.data.items || []
      userBalance.value = data.data.balance || 0
    }
  } catch (error) {
    console.error('获取订单数据失败:', error)
    ElMessage.error('获取订单数据失败')
  }
}

// 应用优惠码
const applyCoupon = async () => {
  if (!couponCode.value) {
    ElMessage.warning('请输入优惠码')
    return
  }
  
  couponLoading.value = true
  try {
    const { data } = await request.post('/api/v1/promo-codes/validate', {
      code: couponCode.value
    })
    if (data?.data) {
      couponDiscount.value = data.data.discount || 0
      couponApplied.value = true
      ElMessage.success('优惠码已应用')
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '优惠码无效')
  } finally {
    couponLoading.value = false
  }
}

// 移除优惠码
const removeCoupon = () => {
  couponCode.value = ''
  couponApplied.value = false
  couponDiscount.value = 0
}

// 提交订单
const submitOrder = async () => {
  if (!agreeTerms.value) {
    ElMessage.warning('请先同意服务条款')
    return
  }
  
  if (!selectedPayment.value) {
    ElMessage.warning('请选择支付方式')
    return
  }
  
  submitting.value = true
  try {
    const checkoutItems = JSON.parse(localStorage.getItem('checkout_items') || '[]')
    
    const { data } = await request.post('/api/v1/orders', {
      items: checkoutItems,
      payment: selectedPayment.value,
      coupon: couponApplied.value ? couponCode.value : '',
      use_balance: useBalance.value,
      contact: needContact.value ? contactForm.value : undefined,
      remark: remark.value
    })
    
    if (data?.data) {
      // 清除购物车中已结算的商品
      localStorage.removeItem('checkout_items')
      
      // 判断是否是银行转账
      if (data.data.type === 'bank_transfer') {
        try {
          bankInfo.value = JSON.parse(data.data.bank_info)
        } catch (e) {
          bankInfo.value = data.data.bank_info
        }
        bankDialogVisible.value = true
      } else if (data.data.pay_url) {
        // 其他支付方式跳转
        window.location.href = data.data.pay_url
      } else {
        router.push(`/user/orders/${data.data.order_id}`)
        ElMessage.success('订单已创建')
      }
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '创建订单失败')
  } finally {
    submitting.value = false
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
  padding-top: 64px;
  
  .checkout-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px 20px;
  }
  
  // 面包屑
  .breadcrumb {
    margin-bottom: 24px;
  }
  
  // 结算内容
  .checkout-content {
    display: flex;
    gap: 24px;
  }
  
  // 左侧主内容
  .checkout-main {
    flex: 1;
    
    .section-card {
      background: #fff;
      border-radius: 12px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
      margin-bottom: 20px;
      
      .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 16px 20px;
        border-bottom: 1px solid #f0f0f0;
        
        h3 {
          font-size: 16px;
          font-weight: 600;
          color: #1a2332;
          margin: 0;
        }
      }
      
      .section-body {
        padding: 20px;
      }
    }
    
    // 订单商品
    .order-item {
      display: flex;
      align-items: center;
      padding: 16px 0;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .item-icon {
        width: 48px;
        height: 48px;
        background: #e3f2fd;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #1976d2;
        margin-right: 12px;
      }
      
      .item-info {
        flex: 1;
        
        h4 {
          font-size: 15px;
          font-weight: 600;
          color: #1a2332;
          margin: 0 0 6px;
        }
        
        .item-meta {
          font-size: 13px;
          color: #909399;
          display: flex;
          gap: 12px;
        }
      }
      
      .item-quantity {
        font-size: 14px;
        color: #666;
        margin: 0 20px;
      }
      
      .item-price {
        font-size: 16px;
        font-weight: 600;
        color: #f56c6c;
      }
    }
    
    // 优惠码
    .coupon-input {
      display: flex;
      gap: 12px;
      align-items: center;
      
      .el-input {
        max-width: 300px;
      }
      
      .coupon-info {
        margin-left: 12px;
      }
    }
    
    // 支付方式
    .payment-methods {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      
      .payment-item {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 12px 20px;
        background: #f5f7fa;
        border: 2px solid transparent;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s;
        
        &:hover {
          border-color: #ddd;
        }
        
        &.active {
          background: #ecf5ff;
          border-color: #409eff;
        }
        
        .payment-icon {
          width: 28px;
          height: 28px;
          
          img {
            width: 100%;
            height: 100%;
            object-fit: contain;
          }
        }
        
        .payment-name {
          font-size: 14px;
          color: #333;
        }
        
        .check-icon {
          color: #409eff;
          margin-left: 4px;
        }
      }
    }
  }
  
  // 右侧边栏
  .checkout-sidebar {
    width: 360px;
    flex-shrink: 0;
    
    .order-summary {
      background: #fff;
      border-radius: 12px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
      padding: 24px;
      position: sticky;
      top: 88px;
      
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a2332;
        margin: 0 0 20px;
      }
      
      .summary-list {
        margin-bottom: 20px;
        
        .summary-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 10px 0;
          font-size: 14px;
          color: #666;
          
          .discount {
            color: #67c23a;
          }
          
          .total {
            font-size: 20px;
            font-weight: 700;
            color: #f56c6c;
          }
        }
      }
      
      .balance-pay {
        padding: 16px;
        background: #f5f7fa;
        border-radius: 8px;
        margin-bottom: 16px;
        
        .balance-amount {
          display: flex;
          justify-content: space-between;
          margin-top: 10px;
          font-size: 13px;
          color: #666;
          
          .discount {
            color: #67c23a;
          }
        }
      }
      
      .actual-pay {
        padding: 16px;
        background: #ecf5ff;
        border-radius: 8px;
        margin-bottom: 20px;
        
        .pay-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          
          .amount {
            font-size: 24px;
            font-weight: 700;
            color: #f56c6c;
          }
        }
      }
      
      .submit-btn {
        width: 100%;
        height: 48px;
        font-size: 16px;
      }
      
      .agreement {
        margin-top: 16px;
        font-size: 13px;
        color: #909399;
        
        a {
          color: #409eff;
          text-decoration: none;
        }
      }
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .checkout-page {
    .checkout-content {
      flex-direction: column;
    }
    
    .checkout-sidebar {
      width: 100%;
    }
    
    .payment-methods {
      flex-direction: column;
    }
  }
}

// 银行转账对话框样式
.bank-transfer-info {
  .bank-info-card {
    margin: 20px 0;
    padding: 20px;
    background: #f8f9fa;
    border-radius: 8px;
    
    .info-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px dashed #e9ecef;
      
      &:last-child {
        border-bottom: none;
      }
      
      .label {
        color: #666;
        font-size: 14px;
      }
      
      .value {
        color: #333;
        font-weight: 500;
        display: flex;
        align-items: center;
        gap: 8px;
        
        &.amount {
          color: #f56c6c;
          font-size: 20px;
          font-weight: 600;
        }
      }
    }
  }
  
  .bank-instructions {
    margin: 20px 0;
    padding: 16px;
    background: #fdf6ec;
    border-radius: 8px;
    
    h4 {
      margin: 0 0 12px 0;
      color: #e6a23c;
      font-size: 14px;
    }
    
    p {
      margin: 0;
      color: #666;
      font-size: 14px;
      line-height: 1.6;
    }
  }
  
  .bank-qrcode {
    margin: 20px 0;
    text-align: center;
    
    h4 {
      margin: 0 0 16px 0;
      color: #333;
      font-size: 14px;
    }
    
    .qrcode-img {
      max-width: 200px;
      border: 1px solid #eee;
      border-radius: 8px;
    }
  }
}
</style>
