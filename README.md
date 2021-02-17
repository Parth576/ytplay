# ytplay

Search and play songs right from your terminal!

## Dependencies

- youtube-dl
- ffmpeg (This program currently uses the ffplay utility to handle the sound)

## Installation 

- Download the binary from the releases section (Currently supported on linux only) 
- Get API key for the Youtube API by creating a new project [here](https://console.developers.google.com)
- Install ```youtube-dl``` and ```ffmpeg```

## Usage
- List all available flags ```ytplay -h``` or ```ytplay --help```
- Run ```ytplay -key <your-api-key>``` (API key needs to be set in order to search for songs)
- Then just search for songs by keyword ```ytplay search <keyword>``` (Make sure there are no spaces in the keyword)
- Ctrl+C is used for pausing/stopping the song
- A song stopped with ^C can be resumed by typing ```ytplay -resume```

## TODO
- Creating and managing playlists
- Cross-platform compatibility