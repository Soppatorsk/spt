<template>
  <div>
    <img alt="Vue logo" src="./assets/logo.png">
    <SpotifyAuth></SpotifyAuth>
    <ButtonTestVue></ButtonTestVue>
    <CreatePlaylistVue></CreatePlaylistVue>
    <PlaylistGallery></PlaylistGallery>
  </div>
</template>

<script>
import CreatePlaylistVue from './components/CreatePlaylist.vue'
import PlaylistGallery from './components/PlaylistGallery.vue'
import SpotifyAuth from './components/SpotifyAuth.vue'

export default {
  name: 'App',
  components: {
    CreatePlaylistVue,
    SpotifyAuth,
    PlaylistGallery,
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
