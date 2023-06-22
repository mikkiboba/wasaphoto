<script>
	export default{
		data: function() {
			return {
				errormsg: null,
				err: null,
				username: null,
				disableButton: null
			}
		},
		methods: {
			doLogin: async function() {
				this.errormsg = null;
				try {
					if (this.username == null || this.username.localeCompare("") == 0) {
						this.errormsg = "Your credentials are not valid!";
						return;
					}
					this.err = this.username;
					let res = await this.$axios.post("/session", {
						username: this.username
					});
					localStorage.setItem('username', this.username);
					localStorage.setItem('token', res.data.token);
					this.$root.logIn();
					this.$router.push("/home");
				} catch (err) {
				this.errormsg = err.toString();
				this.err = "ciao";
			  	}
			},
			checkInput: async function() {
				let res = !!this.username.match(/^[a-z0-9]+$/i)
				if (!res){
					this.errormsg = "The username is not valid"
					this.disableButton = true
				} else {
					this.errormsg = null
					this.disableButton = false
				}
			}
		}
	}
</script>

<template>
	<link href="../assets/style.css" rel="stylesheet">

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>


	<main class="form-signin w-100 m-auto text-center">
		<form>
			<img class="mb-4" src="../assets/wasa-logo.png" alt="Logo" width="255" height="255">
			<h1 class="h3 mb-3 fw-normal" style="color:white;">Please sign in</h1>
		
			<div class="form-floating">
			  <input type="string" class="form-control" id="floatingInput" v-model="username" placeholder="name@example.com" v-on:input="checkInput()" >
			  <label for="floatingInput">Username</label>
			</div>
			<button class="mt-1 w-100 btn btn-lg btn-primary border-white ownbtn" type="button" @click="doLogin" v-bind:disabled="disableButton">Sign in</button>
		</form>
	  </main>
</template>

<style>
</style>
