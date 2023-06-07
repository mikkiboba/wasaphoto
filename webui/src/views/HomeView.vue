<script>
    export default (await import('vue')).defineComponent({
		data: function() {
			return {
				errormsg: null,
				time: null,
                stream: []
			}
		},
		methods: {
            refresh: async function () {
                this.errormsg = null
                this.time = new Date().toString().split(" ")[4]
                try {
                    let res = await this.$axios.get(`/users/${localStorage.getItem("username")}/stream`)
                    this.stream = res.data
                } catch (err) {
                    this.errormsg = err.toString()
                }
            }
		},
        mounted() {
            this.$root.logIn()
            this.refresh()
        }
	})
</script>


<template>
    <br>

    <div class="row d-flex justify-content-center">
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div class="col-md-8 col-lg-6">
        <div class="card shadow-0 border" style="background-color: #f0f2f5;">
            <div class="text-center">
            <img class="like" width="32" height="32" src="/refresh.png" @click="refresh()"> Last refresh: {{ time }}
        </div>

        <div class="text-center">
            <div v-if="stream == null">
                No posts
            </div>
            <div v-else>
                <PostCard v-for="post in stream" v-bind:post="post" :key="post" v-bind:home="true"></PostCard>
            </div>
        </div>
        </div>
  </div>
</div>


    <div class="row d-flex justify-content-center">
        
    </div>

</template>