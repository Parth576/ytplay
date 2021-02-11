# ytplay

Search and play songs right from your terminal!

## Dependencies

- youtube-dl
- ffmpeg (This program currently uses the ffplay utility to handle the sound)

## Installation 

- Download the binary from the releases section (Currently supported on linux only) 
- Get API key for the Youtube API from [here](https://console.developers.google.com)
- Run ```ytplay -key=<your-api-key>```
- Then just search for songs by keyword ```ytplay search <keyword>``` (Make sure there are no spaces)