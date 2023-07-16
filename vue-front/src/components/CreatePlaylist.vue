<template>
    <div>
        <form @submit.prevent="handleSubmit">
            <label>Playlist ID:</label>
            <input type="text" required v-model="id">
            <button>Submit</button>
        </form>
    </div>
</template>

<script>
import VueCookies from 'vue-cookies'
export default {
    methods: {
        handleSubmit() {
            console.log('form submit ' + this.id)
            const parsedId = this.parseSpotifyLink(this.id)
            const token = VueCookies.get("token")
            fetch('http://localhost:3000/playlists/', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id: parsedId, token: token })
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
        id: ''
    }
}
}
</script>

<style>

</style>