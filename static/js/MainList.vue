<template>
    <div>
        <h1>User list</h1>
        <p>Welcome {{ user }}</p>
        <ul>
            <li v-for="l in displayList" :key="l.key">
                <a href="#" @click.prevent="sendMessage(l.key)">{{ l.name }}</a>
            </li>
        </ul>
        <div class="row">
            <div class="col-md-6">
                <h3>Send Message</h3>
                <div v-for="(item, index) in send" :key="index">
                    {{ item }}
                </div>
            </div>
            <div class="col-md-6">
                <h3>Received Message</h3>
                <div v-for="(item, index) in rec" :key="index">
                    {{ item }}
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import Axios from 'axios'

export default {
    props: ['user'],
    data() {
        return {
            list: [],
            send:[],
            rec:[],
        }
    },
    computed: {
        displayList() {
            let vm = this
            return this.list.filter(l => {
                return vm.user != l.key
            })
        }
    },
    mounted() {
        let vm = this;
        //ping user 
        new EventSource(`/events/ping/${vm.user}`)

        //listen list of user
        const e = new EventSource(`/events/list`)
        e.onmessage = ev => {
           const list = JSON.parse(ev.data)
           let formatList = []
           for (var key in list) {
               formatList.push({
                   key: key, 
                   name: list[key]
               })
           }
           
           this.list = formatList
        }

        //listen chat 
        const c = new EventSource('/events/chat/' + this.user)
        c.onmessage = ev => {
            const data = ev.data
            this.rec.push(data)
        }
    },
    methods: {
        sendMessage(to) {
            let vm = this;
            const message = prompt("Your message", "message")
            if (message) {
                const params = {
                    from: this.user,
                    to: to,
                    message: message
                }

                Axios.post("/send", params).then(result => {
                    const data = JSON.stringify(result.data)
                    vm.send.push(data)
                }).catch(error => {

                })
            }
        }
    },
}
</script>