<template>
    <div class="playlistInput">
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
export default {
    methods: {
        handleSubmit() {
            this.requestStatus = "Loading..."
            console.log('form submit ' + this.id)
            const parsedId = this.parseSpotifyLink(this.id)
            fetch(this.$hostname+'playlist/'+parsedId, {
          })
          .then(response => {
            if (response.status == 500) {
                this.requestStatus = "Playlist private, invalid or duplicate"
                throw new Error("Invalid or duplicate link")
            } else {
                response = response.json()
            }
            })
          .then(data => {
                console.log(data);
                //TODO auto call /playlists
                this.requestStatus = "Success! Refresh the page"
                window.location.reload()
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
    }
},


}
</script>

<style>

.playlistInput {
  margin-bottom: 5vh;

}
</style>