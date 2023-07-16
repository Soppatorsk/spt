<template>
  <div>
    <button @click="printToken">Print Token Cookie</button>
  </div>
</template>

<script>
// import VueCookies from 'vue-cookies';
export default {
  methods: {
    printToken() {
      // const token = VueCookies.get("token")
 
fetch('http://localhost:3000/createPlaylist/', {
            method: 'POST',
            /*
            headers: {
              'Content-Type': 'application/json'
            },
            */
          })
            .then(response => response.json())
            .then(data => {
              console.log(data);
            })
            .catch(error => {
              console.error(error);
            });
    }
  }
}
</script>
  
  <style>
  /* Your component styling goes here */
  </style>
  