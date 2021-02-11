package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Parth576/ytplay/colors"
	"github.com/spf13/viper"
)

// PrintErr prints the error to logs
func PrintErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

var idMap = make(map[int]string)

func main() {
	homedir, err := os.UserHomeDir()
	PrintErr(err)

	configPath := filepath.Join(homedir, ".ytplay.yaml")
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		PrintErr(err)
		defer file.Close()
	}
	viper.SetConfigFile(configPath)
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Could not find config file at " + configPath)
		}
	}
	viper.WriteConfig()

	apiKey := viper.GetString("YOUTUBE_API_KEY")

	var keyFlag string
	flag.StringVar(&keyFlag, "key", "", "Set Youtube API key")
	flag.Parse()

	if apiKey == "" && keyFlag == "" {
		fmt.Println("Youtube API key not set, please generate API key from https://console.developers.google.com")
		fmt.Println("And then run the command:")
		fmt.Println("ytplay -key=<your-api-key>")
		os.Exit(1)
	}

	if keyFlag != "" {
		viper.Set("YOUTUBE_API_KEY", keyFlag)
		viper.WriteConfig()
		apiKey = keyFlag
		fmt.Println("Youtube API key saved!")
		os.Exit(1)
	}
	//err := godotenv.Load()
	//PrintErr(err)
	//apiKey := os.Getenv("YOUTUBE_API_KEY")

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
	} else if len(argList) > 1 {
		fmt.Printf("Only 1 argument required but %v arguments provided", len(argList))
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
		pprint(items)
		var index int
		fmt.Printf("\n%sEnter choice > %s", colors.Yellow, colors.Reset)
		fmt.Scanln(&index)
		//fmt.Println(idMap)
		videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", idMap[index])
		//fmt.Println(videoUrl)
		//output, err := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", videoUrl, "-o", cachePath).Output()
		//PrintErr(err)
		//fmt.Println(string(output))
		//youtube-dl -x --audio-format mp3 "https://www.youtube.com/watch?v=J_QGZspO4gg" -o ~/Downloads/youtubedl/bruh.mp3

		ytdlExecutable, err := exec.LookPath("youtube-dl")
		if err != nil {
			fmt.Println("Please install youtube-dl")
			fmt.Println(err)
		}

		command := &exec.Cmd{
			Path:   ytdlExecutable,
			Args:   []string{ytdlExecutable, "-x", "--audio-format", "mp3", videoUrl, "-o", tmpFilepath},
			Stdout: os.Stdout,
			Stdin:  os.Stdout,
		}

		if err = command.Run(); err != nil {
			fmt.Println(err)
		}

		ffplayExec, err := exec.LookPath("ffplay")
		if err != nil {
			fmt.Println("Please install ffmpeg for the ffplay command")
			fmt.Println(err)
		}

		playCmd := exec.Cmd{
			Path:   ffplayExec,
			Args:   []string{ffplayExec, tmpFilepath, "-nodisp", "-autoexit"},
			Stdout: os.Stdout,
			Stdin:  os.Stdout,
		}
		if err = playCmd.Run(); err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Println("Some error occurred with fetching details from the Youtube API")
	}

}

func pprint(items interface{}) {
	for index, v := range items.([]interface{}) {
		videoMap := v.(map[string]interface{})
		id := videoMap["id"]
		idMap[index+1] = id.(map[string]interface{})["videoId"].(string)
		info := videoMap["snippet"]
		typedInfo := info.(map[string]interface{})
		fmt.Printf("%v)  %sTitle:%s %s\n", index+1, colors.Cyan, colors.Reset, typedInfo["title"])
		fmt.Printf("    %sChannel:%s %s\n\n", colors.Cyan, colors.Reset, typedInfo["channelTitle"])
	}
}
