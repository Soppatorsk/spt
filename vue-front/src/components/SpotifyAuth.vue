<template>
 <div>
    <a :href="auth">Authorize with spotify {{ $hostDir }}</a>
 </div> 
</template>

<script>
//TODO green if authorized else no 
//Old Auth used for user data, not needed, clean up later
export default {

    data() {
        return {
            auth: "",
        }
    },
    mounted() {
        fetch(this.$hostname+'auth/')
        .then(res => res.text())
        .then(data => this.auth = data)
        .catch(err => {
            console.log(err.message)
            this.errorMessage = err.message
        })
    }
}
</script>

<style>

</style>