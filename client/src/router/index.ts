import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
// import Home from '../views/Home.vue';
import NotFound from '../views/NotFound.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    // component: Home
    // route level code-splitting
    // this generates a separate chunk (chat.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "chat" */ '../views/Home.vue'),
  },
  {
    path: '/chat',
    name: 'Chat',
    // route level code-splitting
    // this generates a separate chunk (chat.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "chat" */ '../views/Chat.vue'),
    // components: {
    //   default: A,
    //   B,
    // },
    // children: [
    //   {
    //     path: 'a',
    //     component: A
    //   },
    //   {
    //     path: 'b',
    //     components: {
    //       default: B,
    //       Helper
    //     }
    //   }
    // ]
  },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFound },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
