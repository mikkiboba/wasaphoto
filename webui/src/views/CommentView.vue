<script>
import ErrorMsg from '../components/ErrorMsg.vue'


    export default{
    data: function () {
        return {
            username: localStorage.getItem("username"),
            photoid: null,
            photo: null,
            errormsg: null,
            commentText: null,
            comments: null
        }
    },
    methods: {
        getPhoto: async function () {
            this.errormsg = null
            try {
                let res = await this.$axios.get(`/posts/${this.photoid}`, {
                    responseType: 'arraybuffer'
                })
                let blob = new Blob([res.data])
                this.photo = URL.createObjectURL(blob)
                this.loadComments()
            }
            catch (err) {
                this.errormsg = err.toString()
            }
        },
        addComment: async function () {
            this.errormsg = null
            try {
                if(this.commentText.length > 0)
                    await this.$axios.post(`/users/${this.username}/posts/${this.photoid}/comments`, {
                        comment: this.commentText
                    })
                this.commentText = ""
                this.loadComments()
            }
            catch (err) {
                this.errormsg = err.toString()
            }
        },
        loadComments: async function () {
            try {
                let res = await this.$axios.get(`/posts/${this.photoid}/comments`)
                this.comments = res.data
            }
            catch (err) {
                this.errormsg = err.toString()
            }
        },
        back: async function () {
            this.$router.back()
        }
    },
    mounted() {
        this.photoid = this.$route.params.post
        this.getPhoto()
    },
}


</script>


<template>
<br>
<div class="row d-flex justify-content-center">
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div class="col-md-8 col-lg-6">
        <div class="card shadow-0 border" style="background-color: #f0f2f5;">
            <br>
            <div class="ownunderline" @click="back()">
                <img height="32" width="32" class="ownimg" src="/previous.png">Back
            </div>
            <div class="card-body p-4">
                <img :src="photo" class="bd-placeholder-img card-img-top ownimg">
                <div class="form-outline mb-4">
                    <br>
                <input v-model="commentText" type="text" class="form-control" placeholder="Type comment..." />
                <label class="form-label ownunderline" for="addANote" @click="addComment()">+ Add a comment</label>
                </div>

                <CommentCard v-for="comment in comments" v-bind:comment="comment" v-bind:key="comment" v-bind:photoid="this.photoid" v-bind:supera="this"></CommentCard>
            </div>
        </div>
  </div>
</div>


</template>