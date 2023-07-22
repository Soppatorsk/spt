<template>
    <div id="gallery">
      <div :style="'background-color:' + playlist.color" v-for="playlist in playlists" :key="playlist" class="galleryItem">
        <a :href="playlist.collageURL">
          <img id="collage" :src="playlist.collageURL">
        </a>
        <br>
        <a :href="'https://open.spotify.com/playlist/'+playlist.id" target="_blank">{{ playlist.name }} by {{ playlist.user }}</a>
        <p id="ai">{{ playlist.ai }}</p>
        <img id="aiImg" src="../assets/img/sptAi2.png">
          <p :style="'color: ' + playlist.color" id="colorCode"> {{ playlist.color }}</p>
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
    display: flex;
    flex-wrap: wrap;
    align-items: baseline;
    justify-content: center;
  }
  #gallery div {
    width:22vw;
    margin:1vw 1vw;
  }

  #gallery #collage {
    margin:1vw;
    width:20vw
  }
  
  #gallery #ai {
    margin:1vw;
    padding:0.2vw;
    color:black;
    background-color:white;
    border-style: solid;
    border-radius: 1.5vw;
  }

  #aiImg {
    margin-top:-2.5vh;
    /* margin-right:15vw; */
    text-align: left;
    width:4vw;
    float:left
  }

  #colorCode {
    float:right;
    background-color:rgb(20,20,20);
    padding:0.5vw;
    margin-bottom:-1vh;
  }

  #aiImg #colorCode {
    margin:none;
    display: inline-block;
  }
  </style>