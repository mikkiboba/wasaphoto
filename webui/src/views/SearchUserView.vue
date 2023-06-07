<script>

 export default (await import('vue')).defineComponent({
    data: function() {
        return {
            errormsg: null,
            searchBox: null,
            userslist: []
        }
    },
    methods: {
        async search() {
            this.errormsg = null
            try {
                if (this.searchBox.length >= 3) {
                    let res = await this.$axios.get(`/users?username=${this.searchBox}`)
                    res.data.forEach(user => {
                        this.userslist.push(user)
                    });
                } else {
                    this.errormsg = "Search query must be at least 3 characters long"
                }
            } catch (err) {
                this.errormsg = "Can't find the user"
            }
        },
        async searchNew() {
            this.userslist = []
            await this.search()
        }
    },
    mounted() {
        this.errormsg = null
    }
 })


</script>
<template>

    <link href="../style/style.css" rel="stylesheet">
    <br>
    <div class="container-fluid row col-md-9 ms-sm-auto col-lg-10 px-md-2">
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        
        <main>

            <section class="text-center">
                Search username: <input v-on:input="searchNew" v-model="searchBox"  class="rounded-3" type="text" id="search" name="search" placeholder="Username">
            </section>
            <br><br><br>
            <ProfileSearchCard v-for="user in userslist" v-bind:key="user" :username="user"></ProfileSearchCard>

        </main>
    </div>

</template>