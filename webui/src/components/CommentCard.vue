<script>

  export default (await import('vue')).defineComponent({
          props: {
              comment: Object,
              photoid: String,
              supera: Object
          },
          data: function () {
            return {
              errormsg: null,
              isOwner: null
            }
          },
          methods: {
            deleteComment: async function() {
              this.errormsg = null
              try {
                let res = await this.$axios.delete(`/users/${localStorage.getItem("username")}/posts/${this.photoid}/comments/${this.comment.id}`)
                this.supera.loadComments()
              } catch (err) {
                this.errormsg = err.toString()
              }
            }
          },
          mounted() {
            if (this.comment.user === localStorage.getItem("username")) {
              this.isOwner = true
            } else {
              this.isOwner = false
            }
          }
  })

</script>


<template>

  <div class="card mb-4">
    <div class="card-body">
      <p>{{ comment.comment }}</p>

      <div class="d-flex justify-content-between">
        <div class="d-flex flex-row align-items-center">
          <p class="small mb-0 ms-2 ownunderline" @click="this.$router.push(`/users/${comment.user}`)">{{ comment.user }}</p>
        </div>
        <div class="d-flex flex-row align-items-center">
          <label class="small text-muted mb-0 ownunderline" @click="deleteComment()" v-if="isOwner">Delete</label>
        </div>
      </div>
    </div>
  </div>

</template>