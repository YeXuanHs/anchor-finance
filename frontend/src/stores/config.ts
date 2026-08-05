import { defineStore } from 'pinia'
import { ref } from 'vue'
import request from '@/utils/request'

interface PublicConfig {
  // 鍏徃淇℃伅
  company_name: string
  company_email: string
  company_phone: string
  company_address: string
  record_no: string
  system_url: string

  // Logo
  logo_url: string
  logo_url_home: string
  favicon_url: string

  // 鐧诲綍娉ㄥ唽鏂瑰紡
  login_methods: {
    phone: boolean
    email: boolean
    wechat: boolean
    id: boolean
  }
  register_methods: {
    phone: boolean
    email: boolean
    wechat: boolean
  }

  // 鍔熻兘寮€鍏?  affiliate_enabled: boolean
  addfunds_enabled: boolean
  credit_limit: boolean
  show_cancel: boolean
  nologin_send_ticket: boolean
  evaluate_ticket: boolean

  // 鏄剧ず閰嶇疆
  language: string
  allow_user_language: boolean
  default_country: string

  // 娉曞緥鏉℃
  server_clause_url: string
  privacy_clause_url: string

  // SEO
  seo_keywords: string
  seo_desc: string

  // 缁存姢妯″紡
  maintenance_mode: boolean
}

const defaultConfig: PublicConfig = {
  company_name: '',
  company_email: '',
  company_phone: '',
  company_address: '',
  record_no: '',
  system_url: '',
  logo_url: '',
  logo_url_home: '',
  favicon_url: '',
  login_methods: {
    phone: true,
    email: true,
    wechat: false,
    id: false
  },
  register_methods: {
    phone: true,
    email: true,
    wechat: false
  },
  affiliate_enabled: false,
  addfunds_enabled: true,
  credit_limit: false,
  show_cancel: true,
  nologin_send_ticket: false,
  evaluate_ticket: true,
  language: 'zh-cn',
  allow_user_language: true,
  default_country: 'CN',
  server_clause_url: '',
  privacy_clause_url: '',
  seo_keywords: '',
  seo_desc: '',
  maintenance_mode: false
}

export const useConfigStore = defineStore('config', () => {
  const config = ref<PublicConfig>({ ...defaultConfig })
  const loaded = ref(false)

  // 鑾峰彇鍏紑閰嶇疆
  async function fetchPublicConfig() {
    try {
      const res = await request.get('/api/v1/system/settings')
      if (res.data?.data) {
        config.value = { ...defaultConfig, ...res.data.data }
        loaded.value = true
      }
    } catch (error) {
      console.error('Failed to fetch public config:', error)
    }
  }

  // 鑾峰彇鐧诲綍閰嶇疆
  function getLoginMethods() {
    return config.value.login_methods
  }

  // 鑾峰彇娉ㄥ唽閰嶇疆
  function getRegisterMethods() {
    return config.value.register_methods
  }

  // 鏄惁鏄剧ず鏌愪釜鍔熻兘
  function isFeatureEnabled(feature: string): boolean {
    switch (feature) {
      case 'affiliate':
        return config.value.affiliate_enabled
      case 'recharge':
      case 'addfunds':
        return config.value.addfunds_enabled
      case 'credit':
        return config.value.credit_limit
      case 'cancel':
        return config.value.show_cancel
      case 'guest_ticket':
        return config.value.nologin_send_ticket
      case 'ticket_evaluate':
        return config.value.evaluate_ticket
      case 'wechat_login':
        return config.value.login_methods?.wechat || false
      case 'phone_login':
        return config.value.login_methods?.phone || false
      case 'email_login':
        return config.value.login_methods?.email || false
      case 'phone_register':
        return config.value.register_methods?.phone || false
      case 'email_register':
        return config.value.register_methods?.email || false
      default:
        return true
    }
  }

  // 鑾峰彇Logo
  function getLogo(scene: string = 'default'): string {
    switch (scene) {
      case 'home':
      case 'web':
        return config.value.logo_url_home || config.value.logo_url
      case 'admin':
        return config.value.logo_url
      default:
        return config.value.logo_url
    }
  }

  // 鑾峰彇鍏徃淇℃伅
  function getCompanyInfo() {
    return {
      name: config.value.company_name,
      email: config.value.company_email,
      phone: config.value.company_phone,
      address: config.value.company_address,
      record_no: config.value.record_no
    }
  }

  return {
    config,
    loaded,
    fetchPublicConfig,
    getLoginMethods,
    getRegisterMethods,
    isFeatureEnabled,
    getLogo,
    getCompanyInfo
  }
})
