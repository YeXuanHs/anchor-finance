<template>
  <div class="captcha-container">
    <div class="captcha-input-wrapper">
      <el-input
        v-model="inputValue"
        :placeholder="placeholder"
        :size="size"
        @input="onInput"
        @keyup.enter="$emit('verify')"
      >
        <template #append>
          <div class="captcha-image" @click="refresh">
            <img v-if="captchaImage" :src="captchaImage" alt="验证码" title="点击刷新" />
            <div v-else class="captcha-loading">
              <el-icon class="is-loading"><Loading /></el-icon>
            </div>
          </div>
        </template>
      </el-input>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import request from '@/utils/request'

const props = defineProps<{
  modelValue?: string
  captchaKey?: string
  placeholder?: string
  size?: 'large' | 'default' | 'small'
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:captchaId': [id: string]
  verify: []
}>()

const inputValue = ref(props.modelValue || '')
const captchaImage = ref('')
const captchaId = ref('')
const captchaKey = ref(props.captchaKey || generateKey())

function generateKey(): string {
  return Math.random().toString(36).substring(2, 15)
}

async function refresh() {
  try {
    captchaKey.value = generateKey()
    const response = await request.get('/api/v1/captcha/image/json', {
      params: { key: captchaKey.value },
      responseType: 'json'
    })
    
    if (response.data) {
      captchaImage.value = response.data.image
      captchaId.value = response.data.captcha_id
      emit('update:captchaId', captchaId.value)
    }
  } catch (error) {
    console.error('Failed to load captcha:', error)
  }
}

function onInput(value: string) {
  inputValue.value = value
  emit('update:modelValue', value)
}

// Verify captcha
async function verify(): Promise<boolean> {
  try {
    await request.post('/api/v1/captcha/image/verify', {
      key: captchaKey.value,
      captcha_id: captchaId.value,
      digits: inputValue.value
    })
    return true
  } catch {
    refresh()
    return false
  }
}

// Expose methods
defineExpose({
  refresh,
  verify,
  captchaKey: captchaKey.value,
  captchaId: captchaId.value
})

// Watch for external value changes
watch(() => props.modelValue, (val) => {
  inputValue.value = val || ''
})

onMounted(() => {
  refresh()
})
</script>

<style scoped lang="scss">
.captcha-container {
  .captcha-input-wrapper {
    width: 100%;
    
    :deep(.el-input-group__append) {
      padding: 0;
      background-color: #fff;
    }
  }
  
  .captcha-image {
    width: 120px;
    height: 32px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    border-radius: 0 4px 4px 0;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    
    &:hover {
      opacity: 0.8;
    }
  }
  
  .captcha-loading {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f7fa;
  }
}
</style>
