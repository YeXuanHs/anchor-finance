<template>
  <div class="flex w-full h-screen">
    <LoginLeftView />

    <div class="relative flex-1">
      <AuthTopBar />

      <div class="auth-right-wrap">
        <div class="form">
          <h3 class="title">重置管理员密码</h3>
          <p class="sub-title">后台管理员不支持自助找回密码，请联系服务器管理员通过数据库重置</p>

          <div class="mt-6 p-4 bg-gray-50 rounded-lg text-sm leading-7 text-gray-700">
            <p class="font-bold mb-2">重置步骤：</p>
            <p>1. SSH 登录服务器</p>
            <p>2. 执行以下命令生成新密码的 bcrypt 哈希：</p>
            <code class="block my-2 p-2 bg-gray-100 rounded text-xs">
              python3 -c "import bcrypt; print(bcrypt.hashpw(b'你的新密码', bcrypt.gensalt()).decode())"
            </code>
            <p>3. 连接数据库并更新密码：</p>
            <code class="block my-2 p-2 bg-gray-100 rounded text-xs">
              mysql -u root -p anchorfinance<br>
              UPDATE users SET password='哈希值' WHERE username='admin';
            </code>
            <p>4. 重启服务：</p>
            <code class="block my-2 p-2 bg-gray-100 rounded text-xs">
              systemctl restart anchorfinance
            </code>
          </div>

          <div style="margin-top: 20px">
            <ElButton class="w-full custom-height" type="primary" plain @click="toLogin" v-ripple>
              返回登录
            </ElButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  defineOptions({ name: 'ForgetPassword' })

  const router = useRouter()

  const toLogin = () => {
    router.push({ name: 'Login' })
  }
</script>

<style scoped>
  @import '../login/style.css';
</style>
