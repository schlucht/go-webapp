import { createRouter, createWebHistory } from "vue-router";
import OTSBody from '@/components/OTSBody.vue';
import OTSLogin from "@/components/OTSLogin.vue";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: OTSBody,
    },
    {
        path: '/login',
        name: 'Login',
        component: OTSLogin,
    }
];

const router = createRouter (
    {
        routes,
        history: createWebHistory(),
    },
);

export default router;