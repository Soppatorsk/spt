<template>
    <div id="gallery">
      <div v-for="playlist in playlists" :key="playlist" class="image-item">
        <a :href="playlist.collageURL">
          <img :src="playlist.collageURL">
        </a>
        <br>
        <a :href="'https://open.spotify.com/playlist/'+playlist.id" target="_blank">{{ playlist.name }} by {{ playlist.user }}</a>
        <p>{{ playlist.ai }}</p>
      </div>
    </div>
  </template>
  
  <script>
  export default {

    data() {
      return {
        playlists: [],
      };
    },
    mounted() {
      this.loadPlaylists();
    },
    methods: {
      loadPlaylists() {
        fetch("/playlists/")
        .then(res => res.json())
        .then(data => { 
        this.playlists = data;
        console.log(this.playlists)
      })
        .catch(err => console.log(err.message));

      },
    },
  };
  </script>
  
  <style>
  #gallery {
    display: flex;
    flex-wrap: wrap;
  }
  #gallery div {
    width:20vw;
    margin:0 2vw;
  }

  #gallery div img {
    width:20vw
  }
  </style>