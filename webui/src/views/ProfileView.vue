<script>

    export default{
        data: function() {
            return {
                errormsg: null,
                profile: null,
                posts: [],
                user: null,
                youbanned: null,
                theybanned: null,
                isOwner: null,
                isFollowing: null,
            }
        },
    methods: {
        getProfile: async function() {
            this.errormsg = null

            try {
                if (!this.isOwner) {
                    let banRes1 = await this.$axios.get(`/users/${localStorage.getItem("username")}/bans/${this.user}`)
                    if (banRes1.status === 200 && banRes1.data.status === true){
                        this.youbanned = true
                        this.errormsg = `You banned ${this.user}`
                    }

                    let banRes2 = await this.$axios.get(`/users/${this.user}/bans/${localStorage.getItem("username")}`)
                    if (banRes2.status === 200 && banRes2.data.status === true){
                        this.theybanned = true
                        this.errormsg = `${this.user} banned you`
                    }
                    await this.checkFollow()

                }

                if (!this.theybanned && !this.youbanned){
                    let res1 = await this.$axios.get(`/users/${this.user}`, null)
                    this.profile = res1.data
                    let postsres = await this.$axios.get(`/users/${this.user}/posts`)
                    this.posts = postsres.data
                }
            } catch (err) {
                this.errormsg = err.toString()
            }
        },
        unBan: async function (){
            this.errormsg = null
            try {
                await this.$axios.delete(`/users/${localStorage.getItem("username")}/bans/${this.user}`)
                this.youbanned = false
                window.location.reload()
            } catch (err) {
                this.errormsg = err.toString()
            }
        },
        deletePhoto: async function(photoid, index) {
            try {
                await this.$axios.delete(`/users/${localStorage.getItem("username")}/posts/${photoid}`)
                this.posts = []
                this.getProfile()
            } catch (err) {
                this.errormsg = err.toString()
            }
        },
        checkFollow: async function() {
            try {
              let res = await this.$axios.get(`/users/${localStorage.getItem("username")}/follows/${this.user}`)
              if (res.status === 200 && res.data.status === true) {
                this.isFollowing = true
              } else if (res.data.status === false) {
                this.isFollowing = false
              }
            } catch (err) {
              this.errormsg = err.toString()
            }
        }
    },
    mounted() {
        this.$root.logIn()
        if (this.$route.params.user === undefined || this.$route.params.user === localStorage.getItem("username")){
            this.user = localStorage.getItem("username")
            this.isOwner = true
            this.theybanned = false
            this.youbanned = false
        } else {
            this.user = this.$route.params.user
            this.isOwner = false
        }
        this.getProfile()
    }
}

</script>


<template>
    <br>
    <div class="container-fluid row col-md-9 ms-sm-auto col-lg-10 px-md-2">
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <div class="text-center" v-if="youbanned">
            <button class="w-25 btn ownunbanbtn" @Click="unBan()">Unban</button>
        </div>
        <div v-if="profile">
            <ProfileCard v-bind:home="false" v-bind:profile="profile" v-bind:key="profile" v-bind:isOwner="isOwner" v-bind:isFollowing="isFollowing"></ProfileCard>
            <br><br>
            <div class="album py-5 ">
                <div class="container">
                <div class="row row-cols-1 row-cols-sm-2 row-cols-md-3 g-3">
                    <div v-if="this.posts.length == 0">
                        <div class="col-md-9">
                            No posts
                        </div>
                    </div>
                    <div v-for="post in this.posts" v-bind:key="post">

                        <PostCard v-bind:post="post"></PostCard>

                    </div>
                </div>
                </div>
            </div>
        </div>
    </div>

</template>