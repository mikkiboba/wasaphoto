import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import SearchUserView from '../views/SearchUserView.vue'
import SettingsView from '../views/SettingsView.vue'
import HomeView from '../views/HomeView.vue'
import UploadView from '../views/UploadView.vue'
import CommentView from '../views/CommentView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/login', component: LoginView},
		{path: '/home', component: HomeView},
		{path: '/search', component: SearchUserView},
		{path: '/users/:user', component: ProfileView},
		{path: '/settings', component: SettingsView},
		{path: '/profile', component: ProfileView},
		{path: '/uploadPhoto', component: UploadView},
		{path: '/comments/:post', component: CommentView}
	]
})

export default router
