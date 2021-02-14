package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Parth576/ytplay/colors"
	"github.com/spf13/viper"
)

func main() {
	homedir, err := os.UserHomeDir()
	PrintErr(err)

	InitConfig(homedir)

	apiKey := viper.GetString("YOUTUBE_API_KEY")

	var keyFlag string
	flag.StringVar(&keyFlag, "key", "", "Set Youtube API key")
	flag.Parse()

	CheckAPIKey(apiKey, keyFlag)

	cachePath := filepath.Join(homedir, "ytplay.cache")
	tmpFilepath := filepath.Join(cachePath, "tmp.mp3")
	if _, err = os.Stat(cachePath); os.IsNotExist(err) {
		os.MkdirAll(cachePath, 0755)
		fmt.Printf("Cache directory created at %s\n", cachePath)
	}
	os.Remove(tmpFilepath)

	argList := os.Args[1:]
	if len(argList) == 0 {
		fmt.Println("Please give keyword to search")
		os.Exit(1)
	} else if len(argList) > 1 {
		fmt.Printf("Only 1 argument required but %v arguments provided\n", len(argList))
		os.Exit(1)
	}

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=10&q=%s&type=video&key=%s", argList[0], apiKey)
	res, err := http.Get(url)
	PrintErr(err)
	defer res.Body.Close()

	var response interface{}
	body, err := ioutil.ReadAll(res.Body)
	PrintErr(err)

	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &response)
		items := response.(map[string]interface{})["items"]
		idMap := PrettyPrint(items)
		var index int
		fmt.Printf("\n%sEnter choice > %s", colors.Yellow, colors.Reset)
		fmt.Scanln(&index)

		//youtube-dl -x --audio-format mp3 "https://www.youtube.com/watch?v=J_QGZspO4gg" -o ~/Downloads/youtubedl/bruh.mp3
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", idMap[index])

		Command("youtube-dl", videoURL, tmpFilepath)
		Command("ffplay", "", tmpFilepath)

	} else {
		fmt.Println("Some error occurred with fetching details from the Youtube API")
	}

}
