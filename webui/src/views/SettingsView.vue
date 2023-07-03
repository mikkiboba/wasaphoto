<script>

    export default{
        data: function() {
            return {
                errormsg1: null,
                errormsg2: null,
                currentUsername: localStorage.getItem("username"),
                newUsername: null,
                disableButton: null
            }
        },
        methods: {
            changeUsername: async function() {
                this.errormsg = null
                try {
                    await this.$axios.put(`/users/${this.currentUsername}`, {
                        username: this.newUsername
                    })
                    localStorage.setItem("username", this.newUsername)
                    this.$router.push(`/users/${localStorage.getItem("username")}`)
                } catch (err) {
                    this.errormsg = err.toString()
                    if (err.response.status === 409) {
                        this.errormsg1 = "The username is already taken"
                    } else if (err.response.status === 400) {
                        this.errormsg1 = "The username is not valid"
                    }
                }
            },
            checkInput: async function() {
				let res = !!this.newUsername.match(/^[a-z0-9]+$/i)
				let l_check = this.newUsername.length < 4

				if(!res) {
					this.errormsg1 = "The username is not valid"
					this.disableButton = true
				} else {
					this.errormsg1 = null
				}
				if(l_check) {
					this.errormsg2 = "The username should at least be 4 characters long"
					this.disableButton = true
				} else {
					this.errormsg2 = null
				}

				if (res && !l_check) {
					this.disableButton = false
				}
			}
        }
    }

</script>


<template>

    <br>


    <div class="container-fluid row col-md-9 ms-sm-auto col-lg-10 px-md-2">
        <ErrorMsg v-if="errormsg1" :msg="errormsg1"></ErrorMsg>
        <ErrorMsg v-if="errormsg2" :msg="errormsg2"></ErrorMsg>
        <div class="container bg-light rounded-5 shadow w-50 text-center">
            <h2 class="fw-light">Change Username</h2>
            <br>
            <h5>Current username: {{ currentUsername }}</h5>

            <div class="text-center">
                New username
                <input type="text" v-model="newUsername" v-on:input="checkInput()">
            </div><br>
            <button class="btn ownbtn" @click="changeUsername()" v-bind:disabled="disableButton">Apply Changes</button>
        </div>
    </div>
</template>