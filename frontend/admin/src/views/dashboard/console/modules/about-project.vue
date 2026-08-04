<template>
  <div class="art-card p-5 flex-b mb-5 max-sm:mb-4">
    <div>
      <h2 class="text-2xl font-medium">关于项目</h2>
      <p class="text-g-700 mt-1">{{ devInfo }}</p>
      <p class="text-g-700 mt-1">技术栈：Vue3、TypeScript、Vite、Element Plus、Go、Gin、GORM</p>

      <div class="flex flex-wrap gap-3.5 max-w-150 mt-9">
        <div
          class="w-60 flex-cb h-12.5 px-3.5 border border-g-300 c-p rounded-lg text-sm bg-g-100 duration-300 hover:-translate-y-1 max-sm:w-full"
          v-for="link in linkList"
          :key="link.label"
          @click="goPage(link.url)"
        >
          <span class="text-g-700">{{ link.label }}</span>
          <ArtSvgIcon icon="ri:arrow-right-s-line" class="text-lg text-g-600" />
        </div>
        <!-- 开发者QQ（纯文本，不可点击） -->
        <div
          class="w-60 flex-cb h-12.5 px-3.5 border border-g-300 rounded-lg text-sm bg-g-100 max-sm:w-full"
        >
          <span class="text-g-700">开发者QQ：{{ devQQ }}</span>
        </div>
      </div>
    </div>
    <img class="w-75 max-md:!hidden" src="@imgs/draw/draw1.png" alt="draw1" />
  </div>
</template>

<script setup lang="ts">
  import AppConfig from '@/config'
  import { WEB_LINKS } from '@/utils/constants'

  const systemName = AppConfig.systemInfo.name
  const devInfo = WEB_LINKS.DEV_INFO
  const devQQ = WEB_LINKS.DEV_QQ

  const linkList = [
    { label: '开发者GitHub', url: WEB_LINKS.GITHUB },
    { label: 'Telegram', url: WEB_LINKS.DEV_TG },
    { label: 'QQ交流群', url: WEB_LINKS.QQ_GROUP_LINK }
  ]

  /**
   * 在新标签页中打开指定 URL
   * @param url 要打开的网页地址
   */
  const goPage = (url: string): void => {
    window.open(url, '_blank', 'noopener,noreferrer')
  }
</script>
