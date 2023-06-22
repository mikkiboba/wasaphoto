<script>

    export default{
        props: {
            profile: Object,
            isOwner: Boolean,
            isFollowing: Boolean
        },
        data: function() {
            return {
                id: this.profile.id,
                username: this.profile.username,
                posts: this.profile.posts,
                numPosts: this.profile.postcounter,
                followers: this.profile.followcounter,
                following: this.profile.followingcounter,
                banned: this.profile.youbanned,
                errormsg: null,
            }
        },
        methods: {
          banUser: async function() {
            try {
              await this.$axios.put(`/users/${localStorage.getItem("username")}/bans/${this.profile.username}`)
              this.banned = true
              window.location.reload()
            } catch (err) {
              this.errormsg = err.toString()
            } 
          },
          followUser: async function() {
            try {
              await this.$axios.put(`/users/${localStorage.getItem("username")}/follows/${this.profile.username}`)
              window.location.reload()
            } catch (err) {
              this.errormsg = err.toString()
            }
          },
          unfollowUser: async function() {
            try {
              await this.$axios.delete(`/users/${localStorage.getItem("username")}/follows/${this.profile.username}`)
              window.location.reload()
            } catch (err) {
              this.errormsg = err.toString()
            }
          },
          changeUsername() {
            this.$router.push("/settings")
          }
        }
    }

</script>

<template>

    <link href="../assets/style.css" rel="stylesheet">

    <section class="py-5 text-center container bg-light rounded-5 shadow">
      <div class="row py-lg-5">
        <div class="col-lg-3 col-md-8 mx-auto">
          <h1 class="fw-light">{{ username }}</h1>
          <div class="d-flex justify-content-center rounded-3 p-2 mb-2"
                  style="background-color: #efefef;">
                  <div>
                    <p class="small text-muted mb-1">Following</p>
                    <p class="mb-0">{{ following }}</p>
                  </div>
                  <div class="px-3">
                    <p class="small text-muted mb-1">Followers</p>
                    <p class="mb-0">{{ followers }}</p>
                  </div>
                  <div>
                    <p class="small text-muted mb-1">Photos</p>
                    <p class="mb-0">{{ numPosts }}</p>
                  </div>
                </div>
                
            <button class="btn btn-primary my-2 ownbtn" v-if="!isOwner && !isFollowing" @Click="followUser()">Follow</button>
            <button class="btn btn-primary my-2 ownunfollowbtn" v-if="!isOwner && isFollowing" @Click="unfollowUser()">Unfollow</button>
            <button class="btn btn-secondary my-2 ownbanbtn" v-if="!isOwner" @Click="banUser()">Ban</button>
            <button class="btn btn-secondary my-2 ownbtnname" v-if="isOwner" @click="changeUsername()">Change Username</button>


        </div>
      </div>
    </section>

</template>