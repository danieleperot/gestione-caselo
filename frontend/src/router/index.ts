import { createRouter as createVueRouter, createWebHistory, Router, RouterHistory } from 'vue-router';
import Index from '../views/Index.vue';

export function createRouter(history: RouterHistory): Router {
  return createVueRouter({
    history,
    routes: [
      {
        path: '/',
        name: 'index',
        component: Index,
      },
    ],
  });
}

export const router = createRouter(createWebHistory());
