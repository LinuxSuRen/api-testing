import { createRouter, createWebHashHistory } from 'vue-router'
import type {RouteRecordRaw} from 'vue-router'

import { Document, Promotion, Menu, Location, Share } from '@element-plus/icons-vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Index',
    redirect: '/welcome',
    component: () => import('../views/Index.vue'),
    children: [
      {
        path: '/welcome',
        name: 'welcome',
        component: () => import('../views/WelcomePage.vue'),
        meta: {
          isShow: true,
          icon: Share,
          title: '欢迎'
        }
      },
      {
        path: '/test',
        name: 'test',
        component: () => import('../views/test/TestingPanel.vue'),
        meta: {
          isShow: true,
          icon: Promotion,
          title: '测试'
        }
      },
      {
        path: '/mock',
        name: 'Mock',
        component: () => import('../views/mock/Mock.vue'),
        meta: {
          isShow: true,
          icon: Menu,
          title: 'Mock'
        }
      },
      {
        path: '/secret',
        name: 'secret',
        component: () => import('../views/cret/Secret.vue'),
        meta: {
          isShow: true,
          icon: Document,
          title: '凭据'
        }
      },
      {
        path: '/store',
        name: 'store',
        component: () => import('../views/store/Store.vue'),
        meta: {
          isShow: true,
          icon: Location,
          title: '存储'
        }
      }
    ]
  },

  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/error/404.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  // output route info for debug
  // console.log(to, from)
  next()
})

export default router
