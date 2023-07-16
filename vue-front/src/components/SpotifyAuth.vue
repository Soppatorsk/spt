<template>
 <div>
    <a :href="auth">Authorize with spotify</a>
 </div> 
</template>

<script>
export default {

    data() {
        return {
            auth: ""
        }
    },
    mounted() {
        fetch('/auth/')
        .then(res => res.text())
        .then(data => this.auth = data)
        .catch(err => console.log(err.message))
    }
}
</script>

<style>

</style>