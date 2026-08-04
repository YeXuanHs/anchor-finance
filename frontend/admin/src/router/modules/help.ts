import { AppRouteRecord } from '@/types/router'
import { WEB_LINKS } from '@/utils/constants'

export const helpRoutes: AppRouteRecord[] = [
  {
    name: 'Document',
    path: '',
    component: '',
    meta: {
      title: 'menus.help.document',
      icon: 'ri:bill-line',
      link: WEB_LINKS.GITHUB + '/anchor-finance#readme',
      isIframe: false,
      keepAlive: false
    }
  },
  {
    name: 'QQGroup',
    path: '',
    component: '',
    meta: {
      title: 'QQ交流群',
      icon: 'ri:qq-line',
      link: `https://qm.qq.com/q/${WEB_LINKS.QQ_GROUP}`,
      isIframe: false,
      keepAlive: false
    }
  },
  {
    name: 'ChangeLog',
    path: '/change/log',
    component: '/change/log',
    meta: {
      title: 'menus.plan.log',
      showTextBadge: `v${__APP_VERSION__}`,
      icon: 'ri:gamepad-line',
      keepAlive: false
    }
  }
]
