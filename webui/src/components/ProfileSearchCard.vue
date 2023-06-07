<script>

    export default (await import('vue')).defineComponent({
        props: {
            username: String,
        },
        data: function() {
            return{}
        },
        methods: {
            showProfile() {
                this.$router.push(`/users/${this.username}`)
            }
        }
    })

</script>

<template>

    <div class="w-50 container-sm bg-light rounded-5 shadow">
        <div class="row py-sm-2">
          <div class="col-8">
            <h2 class="fw-light ml-auto">{{this.username}}</h2> 
          </div>
          <div class="col-4 text-right">
            <button @click="showProfile()" class="btn ownviewbtn mt-1">View Profile</button>
          </div>
        </div>
    </div>  
    <br>

</template>