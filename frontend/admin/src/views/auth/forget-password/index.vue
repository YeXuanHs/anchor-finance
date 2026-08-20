<template>
  <div class="flex w-full h-screen">
    <LoginLeftView />

    <div class="relative flex-1">
      <AuthTopBar />

      <div class="auth-right-wrap">
        <div class="form">
          <h3 class="title">{{ $t('forgetPasswordReset.title') }}</h3>
          <p class="sub-title">{{ $t('forgetPasswordReset.subtitle') }}</p>

          <div class="mt-6 p-4 bg-gray-50 rounded-lg text-sm leading-7 text-gray-700">
            <p class="font-bold mb-2">{{ $t('forgetPasswordReset.stepsTitle') }}</p>
            <p>{{ $t('forgetPasswordReset.step1') }}</p>
            <p>{{ $t('forgetPasswordReset.step2') }}</p>
            <code class="block my-2 p-2 bg-gray-100 rounded text-xs">
              python3 -c "import bcrypt; print(bcrypt.hashpw(b'{{ $t('forgetPasswordReset.newPasswordPlaceholder') }}', bcrypt.gensalt()).decode())"
            </code>
            <p>{{ $t('forgetPasswordReset.step3') }}</p>
            <code class="block my-2 p-2 bg-gray-100 rounded text-xs">
              mysql -u root -p anchorfinance<br>
              UPDATE users SET password='{{ $t('forgetPasswordReset.hashPlaceholder') }}' WHERE username='admin';
            </code>
            <p>{{ $t('forgetPasswordReset.step4') }}</p>
            <code class="block my-2 p-2 bg-gray-100 rounded text-xs">
              systemctl restart anchorfinance
            </code>
          </div>

          <div style="margin-top: 20px">
            <ElButton class="w-full custom-height" type="primary" plain @click="toLogin" v-ripple>
              {{ $t('forgetPasswordReset.backToLogin') }}
            </ElButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { $t } from '@/locales'

  defineOptions({ name: 'ForgetPassword' })

  const router = useRouter()

  const toLogin = () => {
    router.push({ name: 'Login' })
  }
</script>

<style scoped>
  @import '../login/style.css';
</style>
