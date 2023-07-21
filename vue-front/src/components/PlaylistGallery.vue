<template>
    <div id="gallery">
      <div :style="'background-color:' + playlist.color" v-for="playlist in playlists" :key="playlist" class="galleryItem">
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
        fetch(this.$hostname+"playlists/")
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
    margin:auto;
    display: flex;
    flex-wrap: wrap;
  }
  #gallery div {
    width:22vw;
    margin:1vw 1vw;
  }

  #gallery div p {
    margin:1vw;
  }

  #gallery div img {
    margin:1vw;
    width:20vw
  }
  </style>