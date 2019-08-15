<template>
    <div>
        <div class="row">
                <div class="col-md-6 login-form-1">
                    <h3>Check Your Name</h3>
                    <form @submit.prevent="submitForm">
                        <div class="form-group">
                            <input type="text" class="form-control" placeholder="Your Name" v-model="name"/>
                        </div>
                        <div class="form-group">
                            <input type="submit" class="btn btn-primary">
                        </div>
                    </form>
                </div>
        </div>
    </div>
</template>

<script>
import Axios from 'axios'

export default {
    data() {
        return {
            name: ""
        }
    },
    methods: {
        submitForm() {
            if (this.name.trim() == "") {
                alert("empty name")
                return
            }
            const user = {name: this.name}
            const vm = this

            Axios.post("/join", user).then( response => {
                vm.$emit('login', response.data)
            }).catch(err => {
                alert(err.response.data.error)
            })
        }
    },
}
</script>