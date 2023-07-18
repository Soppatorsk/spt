<template>
  <div>
    <h1>torsk.net/spt</h1>
    <CreatePlaylistVue></CreatePlaylistVue>
    <PlaylistGallery></PlaylistGallery>
  </div>
</template>

<script>
import CreatePlaylistVue from './components/CreatePlaylist.vue'
import PlaylistGallery from './components/PlaylistGallery.vue'

export default {
  name: 'App',
  components: {
    CreatePlaylistVue,
    PlaylistGallery,
  }
}
</script>

<style>
body {
  background-color:rgb(20,20,20);
  color:rgb(230,230,230)
}

a {
  color:rgb(230,230,230)
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  margin-top: 60px;

}
</style>
