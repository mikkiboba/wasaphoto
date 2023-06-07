<script setup>
</script>
<script>
export default {
	data() {
		return {
			logged_in : false
		}
	},
	methods : {
		logIn() {
			this.logged_in = true;
		},
		logOut() {
			this.logged_in = false;
			localStorage.clear();
			this.$router.push("/");
		}
	}
}
</script>

<template>
	<div>
		<header id="header" class="navbar navbar-dark sticky-top flex-md-nowrap p-0 shadow">
			<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="/">WASA PHOTO</a>
			<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
		</header>

		<div class="container-fluid row">
			<div class="row">
				<nav v-if="logged_in" v-bind:key="logged_in" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
					<div id="sidebarMenu" class="position-sticky pt-3 sidebar-sticky">
						<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
							<span>General</span>
						</h6>
						<ul class="nav flex-column">
							<li class="nav-item">
								<RouterLink to="/home" class="nav-link">
									<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
									Home
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink to="/search" class="nav-link">
									<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#eye"/></svg>
									Search user
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink to="/profile" class="nav-link">
									<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#smile"/></svg>
									Personal profile
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink to="/uploadPhoto" class="nav-link">
									<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus-circle"/></svg>
									Upload photo
								</RouterLink>
							</li>
						</ul>
					</div>
				</nav>

				<main>
					<RouterView />
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>