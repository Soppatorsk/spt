<template>
    <div>
        <form @submit.prevent="submitColor">
            <label>Playlist Link:</label>
            <input type="text" required v-model="id">
            <button>Colors!</button>
        </form>
        <span v-if="requestStatusColor" class="requestStatus">
            {{ requestStatusColor }}
        </span>
            <div class="songColor" v-for="(song, index) in songs" :key="index">
                <span :style="'background-color:' + song.color"> {{ song.title }}</span>
            </div>
    </div>
</template>

<script>
export default {
    methods: {
        submitColor() {
            this.requestStatusColor = "Loading..."
            const parsedId = this.parseSpotifyLink(this.id)
            fetch(this.$hostname+'test/'+parsedId, {
            })
            .then(res => res.json())
            .then(data => {
                this.songs = data
                console.log(this.songs)
                this.requestStatusColor = "success"
            })
            .catch(error => {
                this.requestStatusColor = error.message
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
            songs: [],
            requestStatusColor: '',
        }
    },
}

</script>

<style>



.songColor {
    margin:none;
    display: inline-block;
}

.songColor span {
    display:inline-block;
    width:95vw;
}

</style>