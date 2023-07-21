# torsk.net/spt/

This is a hobby project for learning Go/Gin backend with a Vue frontend. The project is available at [https://torsk.net/spt/](https://torsk.net/spt/) and can also be accessed via API calls at [https://torsk.net/spt](https://torsk.net/spt).

## API

{ID} refers to a Spotify playlist ID. For example, if the playlist URL is `https://open.spotify.com/playlist/3okg2NywBjVFkFM9LNrWA2?si=8eaebcef35f74ca7`, the ID is `3okg2NywBjVFkFM9LNrWA2`.

- GET `/playlists`
  - Returns all public playlists that are viewable on the frontend.
- GET `/playlist/{ID}`
  - Returns the full JSON data for the specified playlist ID, including all generated features. This can be viewed on the frontend as well.
- GET `/collage/{ID}`
  - Returns only the image file of the generated collage for the specified playlist ID.
- GET `/ai/{ID}`
  - Returns only the AI roast response for the specified playlist ID.
  
