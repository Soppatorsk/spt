<template>
    <div>
        <form @submit.prevent="handleSubmit">
            <label>Playlist Link:</label>
            <input type="text" required v-model="id">
            <button>Submit</button>
        </form>
        <span v-if="requestStatus" class="requestStatus">
            {{ requestStatus }}
        </span>
    </div>
</template>

<script>
import VueCookies from 'vue-cookies'
export default {
    methods: {
        handleSubmit() {
            this.requestStatus = "Loading..."
            console.log('form submit ' + this.id)
            const parsedId = this.parseSpotifyLink(this.id)
            const token = VueCookies.get("token")
            fetch(this.$hostname+'playlists/', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id: parsedId, token: token })
          })
          .then(response => {
            if (response.status == 500) {
                this.requestStatus = "Invalid or duplicate spotify link"
                throw new Error("Invalid or duplicate link")
            } else {
                response = response.json()
            }
            })
          .then(data => {
                console.log(data);
                //TODO auto call /playlists
                this.requestStatus = "Success! Refresh the page"
          })
          .catch(error => {
            console.error(error);
            this.requestStatus = error.message 
          })
        },
        parseSpotifyLink(link) {
  const regex = /^https:\/\/open.spotify.com\/playlist\/([a-zA-Z0-9]+)\?.*$/;
  const match = link.match(regex);
  if (match && match[1]) {
    return match[1];
  }
  return null;
}
    },
data() {
    return {
        id: '',
        requestStatus: '',
        code500: ''
    }
}
}
</script>

<style>

</style>