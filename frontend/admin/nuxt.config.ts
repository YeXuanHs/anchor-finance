export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },

  modules: [
    '@element-plus/nuxt',
    '@pinia/nuxt',
  ],

  css: [
    '~/assets/css/main.css',
  ],

  elementPlus: {
    importStyle: 'css',
    themes: [],
  },

  app: {
    head: {
      title: 'AnchorFinance Admin',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      ],
    },
  },

  ssr: false,

  routeRules: {
    '/admin/**': { ssr: false },
  },
})
