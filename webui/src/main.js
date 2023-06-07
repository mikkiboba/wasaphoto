import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import PostCard from './components/PostCard.vue'
import ProfileCard from './components/ProfileCard.vue'
import ProfileSearchCard from './components/ProfileSearchCard.vue'
import CommentCard from './components/CommentCard.vue'

import './assets/dashboard.css'
import './assets/main.css'
import './assets/style.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
axios.interceptors.request.use(function (config) {
        config.headers['Authorization'] = "Bearer "+localStorage.getItem('token');
        return config;
    }, function (error) {
        console.log(error);
        return Promise.reject(error);
    }
)
app.component("ErrorMsg",ErrorMsg)
app.component("ProfileCard",ProfileCard)
app.component("PostCard",PostCard)
app.component("ProfileSearchCard",ProfileSearchCard)
app.component("CommentCard", CommentCard)
app.use(router)
app.mount('#app')
