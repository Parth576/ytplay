# ytplay

Search and play songs right from your terminal!

## Dependencies

- youtube-dl
- ffmpeg (This program currently uses the ```ffplay``` utility to play audio)

## Installation 

- Get API key for the Youtube API by creating a new project [here](https://console.developers.google.com) and then generate credentials
- Install [youtube-dl](https://youtube-dl.org/latest) and [ffmpeg](https://ffmpeg.org/download.html)
- [For Windows] Add ```youtube-dl.exe``` and ```ffplay.exe``` to the path
- Download the binary from the [releases](https://github.com/Parth576/ytplay/releases/latest) section of this repository
- [For Linux] Make the binary executable 
- [For Windows] Add the binary to the path to execute from anywhere

## Usage
- List all available flags ```ytplay -h``` or ```ytplay --help```
- Run ```ytplay -key <your-api-key>``` (API key needs to be set in order to search for songs from YouTube)
- Then just search for songs by keyword ```ytplay search <keyword>``` (Make sure there are no spaces in the keyword)
- Ctrl+C is used for pausing/stopping the song
- A song stopped with ^C can be resumed by typing ```ytplay -resume```

## Note
- The Youtube API has a default quota of 10k units per day
- Each search request consumes 100 units
- That means you can search for 100 songs per day and the quota will be reset for the next day

### TODO
- Tests
- Creating and managing playlists
- Use a Go audio libary like faiface/beep or oto (remove ffplay dependency)
- Improve UI by using TUI libs
