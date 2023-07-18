import { createApp} from 'vue';
import App from './App.vue';
import VueCookies from 'vue-cookies';

const app = createApp(App);
app.config.globalProperties.$hostname = ''
app.use(VueCookies);
app.mount('#app');