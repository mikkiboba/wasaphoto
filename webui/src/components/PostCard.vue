<script>

    export default{
        props: {
            post: Object,
            home: Boolean
        },
        data: function() {
            return {
                photoid: this.post.id,
                date: this.post.date,
                hour: this.post.hour,
                photo: null,
                errormsg: null,
                liked: null,
                owner: null,
                ishome: this.home
            }
        },
        methods: {
            getPhoto: async function() {
                this.errormsg = null
                try {
                    let res = await this.$axios.get(`/posts/${this.photoid}`, {
                        responseType: 'arraybuffer'
                    })
                    let blob = new Blob([res.data])
                    this.photo = URL.createObjectURL(blob)
                    this.checkLike()
                } catch (err) {
                    this.errormsg = err.toString()
                }
            },
            likePhoto: async function() {
                try {
                    await this.$axios.put(`/posts/${this.photoid}/likes/${localStorage.getItem("username")}`)
                    this.checkLike()
                } catch (err) {
                    this.errormsg = err.toString()
                }
            },
            unlikePhoto: async function() {
                try {
                    await this.$axios.delete(`/posts/${this.photoid}/likes/${localStorage.getItem("username")}`)
                    this.checkLike()
                } catch (err) {
                    this.errormsg = err.toString()
                }
            },
            checkLike: async function() {
                try{
                    let res = await this.$axios.get(`/posts/${this.photoid}/likes/${localStorage.getItem("username")}`)
                    if (res.status === 200 && res.data.status === true) {
                        this.liked = true
                    } else {
                        this.liked = false
                    }
                } catch (err) {
                    this.errormsg = err.toString()
                }
            },
            deletePhoto: async function() {
                this.errormsg = null
                try {
                    await this.$axios.delete(`/users/${localStorage.getItem("username")}/posts/${this.photoid}`)
                    window.location.reload()
                } catch (err) {
                    this.errormsg = err.toString()
                }
            },
            commentButton: async function() {
                this.$router.push(`/comments/${this.photoid}`)
            },
            goProfile: async function() {
                this.$router.push(`/users/${this.post.user}`)
            }
        },
        mounted() {
            if (localStorage.getItem("username") === this.post.user) {
                this.owner = true
            } else { this.owner = false }
            this.getPhoto()
        }
    }

</script>

<template>

    <div class="col ownimg">
        <div class="card shadow-sm">
            <h3 v-if="ishome" @click="goProfile()" class="ownunderline fw-light">{{ this.post.user }}</h3>
          <img :src="photo" class="bd-placeholder-img card-img-top ownimg">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div class="btn-group">
                <img v-if="!liked && !owner" class="like" height="32" width="32" src="/like.png" @click="likePhoto()">
                <img v-if="owner" class="like" height="32" width="32" src="/bin.png" @click="deletePhoto()">
                <img v-else-if="liked && !owner" class="like" height="32" width="32" src ="/like-2.png" @click="unlikePhoto()">
                &nbsp
                <img class="comment" height="32" width="32" src="/comment.png" @click="commentButton()">
              </div>
              <small class="text-muted">{{ date }} {{ hour }}</small>
            </div>
          </div>
        </div>
      </div>
      <br>

</template>