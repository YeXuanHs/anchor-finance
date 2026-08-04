<template>
  <div class="geetest-captcha" ref="captchaRef"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import request from '@/utils/request'

interface GeetestConfig {
  enabled: boolean
  captcha_id: string
}

interface GeetestResult {
  lot_number: string
  captcha_output: string
  pass_token: string
  gen_time: string
}

const props = defineProps<{
  onSuccess?: (result: GeetestResult) => void
  onError?: (error: string) => void
  onClose?: () => void
}>()

const emit = defineEmits<{
  success: [result: GeetestResult]
  error: [error: string]
  close: []
}>()

const captchaRef = ref<HTMLElement | null>(null)
let captchaObj: any = null

// 加载 Geetest JS
function loadGeetestScript(): Promise<void> {
  return new Promise((resolve, reject) => {
    // 检查是否已加载
    if ((window as any).initGeetest4) {
      resolve()
      return
    }

    const script = document.createElement('script')
    script.src = 'https://static.geetest.com/v4/gt4.js'
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Failed to load Geetest script'))
    document.head.appendChild(script)
  })
}

// 初始化 Geetest
async function initGeetest() {
  try {
    // 加载 Geetest JS
    await loadGeetestScript()

    // 获取配置
    const res = await request.get('/api/v1/geetest/register')
    const config: GeetestConfig = res.data

    if (!config.enabled || !config.captcha_id) {
      console.warn('Geetest is not configured')
      return
    }

    // 初始化 Geetest
    const initGeetest4 = (window as any).initGeetest4
    initGeetest4(
      {
        captchaId: config.captcha_id,
        product: 'bind', // bind 模式
        // 可选配置
        riskType: 'slide', // 滑动验证
        hideSuccess: false, // 显示成功提示
        hideClose: false, // 显示关闭按钮
      },
      (captcha: any) => {
        captchaObj = captcha

        // 将验证码插入到 DOM
        if (captchaRef.value) {
          captcha.appendTo(captchaRef.value)
        }

        // 验证成功回调
        captcha.onSuccess(() => {
          const result = captcha.getValidate()
          const geetestResult: GeetestResult = {
            lot_number: result.lot_number,
            captcha_output: result.captcha_output,
            pass_token: result.pass_token,
            gen_time: result.gen_time,
          }

          // 触发成功回调
          if (props.onSuccess) {
            props.onSuccess(geetestResult)
          }
          emit('success', geetestResult)
        })

        // 验证失败回调
        captcha.onError((error: any) => {
          const errorMsg = error?.msg || '验证失败'
          if (props.onError) {
            props.onError(errorMsg)
          }
          emit('error', errorMsg)
        })

        // 关闭回调
        captcha.onClose(() => {
          if (props.onClose) {
            props.onClose()
          }
          emit('close')
        })
      }
    )
  } catch (error) {
    console.error('Failed to initialize Geetest:', error)
  }
}

// 重置验证码
function reset() {
  if (captchaObj) {
    captchaObj.reset()
  }
}

// 获取验证结果
function getValidate(): GeetestResult | null {
  if (captchaObj) {
    const result = captchaObj.getValidate()
    if (result) {
      return {
        lot_number: result.lot_number,
        captcha_output: result.captcha_output,
        pass_token: result.pass_token,
        gen_time: result.gen_time,
      }
    }
  }
  return null
}

// 暴露方法给父组件
defineExpose({
  reset,
  getValidate,
})

onMounted(() => {
  initGeetest()
})

onUnmounted(() => {
  // 清理
  if (captchaObj) {
    try {
      captchaObj.destroy()
    } catch (e) {
      // ignore
    }
    captchaObj = null
  }
})
</script>

<style scoped>
.geetest-captcha {
  width: 100%;
  min-height: 44px;
}
</style>
