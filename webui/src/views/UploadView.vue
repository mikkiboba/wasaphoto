<script>

    export default (await import('vue')).defineComponent({
        data: function () {
            return {
                errormsg: null,
                image: undefined,
                previewimage: undefined
            }
        },
        methods: {
            selectImage() {
                this.errormsg = null
                try {
                    this.image = this.$refs.file.files.item(0)
                    this.previewimage = URL.createObjectURL(this.image)
                } catch (err) {
                    this.errormsg = err.toString()
                }
            },
            uploadImage: async function() {
                this.errormsg = null
                if (this.image == null || this.image == undefined) {
                    this.errormsg = "Please select an image first!"
                    return
                }
                try {
                    let formData = new FormData() 
                    formData.append('file', this.image)
                    for (let i of formData.entries()){
                        console.log(i[0], ',', i[1])
                    }
                    await this.$axios.post(`/users/${localStorage.getItem("username")}/posts`, formData, {
                        headers: {
                            "Content-Type": "multipart/form-data"
                        }
                    })
                    this.$router.push("/profile")
                } catch (err) {
                    this.errormsg = err.toString()
                }
            }
        },
        mounted() {
            this.$root.logIn()
        }
    })

</script>

<template>

    <link rel="stylesheet" href="./src/assets/style.css">
    <br>

    <div class="container-fluid row col-md-9 ms-sm-auto col-lg-10 px-md-2">
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

        <div class="container bg-light py-3 text-center w-25 rounded-3 shadow">
            <h3 class="fw-light">Upload your photo</h3><br>
            <form action="upload.php" method="POST" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="photo" class="form-label">Select a photo:</label>
                    <input type="file" class="form-control" accept="image/*" id="photo" name="photo" ref="file" @change="selectImage()">
                </div><br>
                <button type="button" class="btn btn-primary ownbtn" @click="uploadImage()">Upload</button>
            </form><br>
            <div v-if="previewimage">
                <img class="card-img-top" :src="previewimage" />
            </div>
        </div>
        </div>

</template>