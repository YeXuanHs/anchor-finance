// Element Plus 主题配置
// 深蓝+浅蓝主题色系，参考 jdidc.cn / coyjs.cn 风格
export const themeColors = {
  primary: '#1a56db',        // 深蓝 - 主色调
  primaryLight: '#3b82f6',   // 浅蓝 - 辅助色
  primaryDark: '#1e40af',    // 更深蓝
  primaryBg: '#eff6ff',      // 极浅蓝背景
  
  success: '#10b981',        // 绿色
  successLight: '#34d399',
  
  warning: '#f59e0b',        // 橙黄
  warningLight: '#fbbf24',
  
  danger: '#ef4444',         // 红色
  dangerLight: '#f87171',
  
  info: '#6b7280',           // 灰色
  infoLight: '#9ca3af',
  
  // 背景色
  bgPrimary: '#ffffff',
  bgSecondary: '#f8fafc',
  bgTertiary: '#f1f5f9',
  
  // 文字色
  textPrimary: '#0f172a',
  textSecondary: '#475569',
  textTertiary: '#94a3b8',
  textWhite: '#ffffff',
  
  // 边框色
  borderLight: '#e2e8f0',
  borderDefault: '#cbd5e1',
}

// Element Plus 主题覆盖
export const elementThemeOverrides = {
  common: {
    primaryColor: themeColors.primary,
    primaryColorHover: themeColors.primaryLight,
    primaryColorPressed: themeColors.primaryDark,
    primaryColorSuppl: themeColors.primaryLight,
    
    successColor: themeColors.success,
    successColorHover: themeColors.successLight,
    
    warningColor: themeColors.warning,
    warningColorHover: themeColors.warningLight,
    
    dangerColor: themeColors.danger,
    dangerColorHover: themeColors.dangerLight,
    
    infoColor: themeColors.info,
    infoColorHover: themeColors.infoLight,
    
    borderRadius: '12px',
    borderRadiusSmall: '8px',
    borderRadiusLarge: '16px',
    
    fontFamily: `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'PingFang SC', 'Microsoft YaHei', sans-serif`,
  },
  Button: {
    borderRadius: '12px',
    borderRadiusLarge: '12px',
    borderRadiusSmall: '8px',
    fontWeight: '500',
  },
  Card: {
    borderRadius: '16px',
    boxShadow: '0 4px 24px rgba(0, 0, 0, 0.06)',
  },
  Input: {
    borderRadius: '12px',
  },
  Tag: {
    borderRadius: '8px',
  },
  Menu: {
    borderRadius: '8px',
  },
  Dialog: {
    borderRadius: '16px',
  },
  Message: {
    borderRadius: '12px',
  },
}

// 渐变色预设
export const gradients = {
  primary: `linear-gradient(135deg, ${themeColors.primary} 0%, ${themeColors.primaryLight} 100%)`,
  primaryReverse: `linear-gradient(135deg, ${themeColors.primaryLight} 0%, ${themeColors.primary} 100%)`,
  dark: 'linear-gradient(135deg, #0f172a 0%, #1e293b 100%)',
  darkBlue: `linear-gradient(135deg, ${themeColors.primaryDark} 0%, ${themeColors.primary} 50%, ${themeColors.primaryLight} 100%)`,
  success: `linear-gradient(135deg, ${themeColors.success} 0%, ${themeColors.successLight} 100%)`,
  warning: `linear-gradient(135deg, ${themeColors.warning} 0%, ${themeColors.warningLight} 100%)`,
  danger: `linear-gradient(135deg, ${themeColors.danger} 0%, ${themeColors.dangerLight} 100%)`,
  heroSlide1: `linear-gradient(135deg, #1a56db 0%, #1e40af 40%, #1e3a8a 100%)`,
  heroSlide2: `linear-gradient(135deg, #0ea5e9 0%, #0284c7 40%, #0369a1 100%)`,
  heroSlide3: `linear-gradient(135deg, #3b82f6 0%, #2563eb 40%, #1d4ed8 100%)`,
  heroSlide4: `linear-gradient(135deg, #6366f1 0%, #4f46e5 40%, #4338ca 100%)`,
  stats: `linear-gradient(135deg, ${themeColors.primary} 0%, ${themeColors.primaryDark} 100%)`,
}

export default {
  themeColors,
  elementThemeOverrides,
  gradients,
}
