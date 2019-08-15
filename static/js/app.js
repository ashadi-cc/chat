import Vue from 'vue'
import Login from './Login'
import MainList from './MainList'

const app = new Vue({
    components: {
        Login,
        MainList
    },
    el: '#app',
    data: {
        isLogin: false,
        user: ""
    },
    methods: {
        doList(data) {
            this.isLogin = true
            this.user = data.username
        }
    },
    mounted() {

    }
})