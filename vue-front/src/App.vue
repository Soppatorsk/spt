<template>
  <div>
    <SptHeader></SptHeader>
    <h1>torsk.net/spt</h1>
    <div id="colorLegend">
      <p style="color:red">Energy: █████</p>
      <p style="color:green">Danceability: █████</p>
      <p style="color:blue">Valence (inverse): █████</p>
    </div>
    <CreatePlaylistVue></CreatePlaylistVue>
    <PlaylistGallery></PlaylistGallery>
  </div>
</template>

<script>
import CreatePlaylistVue from './components/CreatePlaylist.vue'
import PlaylistGallery from './components/PlaylistGallery.vue'
import SptHeader from './components/sptHeader.vue'

export default {
  name: 'App',
  components: {
    CreatePlaylistVue,
    PlaylistGallery,
    SptHeader,
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

#colorLegend {
    text-align: right;
    font-size: 0.7em;
    position: absolute;
    top: 7vh;
    left: 90%;
}
</style>


