import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: { name: 'home' }
    },
    {
      path: '/userContent/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/userContent/runGroup/:runGroupName',
      name: 'runGroup',
      props: true,
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/RunGroupView.vue')
    },
  ]
})

export default router
