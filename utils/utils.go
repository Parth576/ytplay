package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

// PrintErr prints the error to logs
func PrintErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type music struct {
	Name    string
	Channel string
	Index   int
}

func PrettyPrint(items interface{}) (map[int]string, int) {
	idMap := make(map[int]string)
	var musicList = []music{}
	for index, v := range items.([]interface{}) {
		videoMap := v.(map[string]interface{})
		id := videoMap["id"]
		idMap[index+1] = id.(map[string]interface{})["videoId"].(string)
		info := videoMap["snippet"]
		typedInfo := info.(map[string]interface{})
		musicList = append(musicList, music{typedInfo["title"].(string), typedInfo["channelTitle"].(string), index + 1})
		//fmt.Printf("%v)  %sTitle:%s %s\n", index+1, colors.Cyan, colors.Reset, typedInfo["title"])
		//fmt.Printf("    %sChannel:%s %s\n\n", colors.Cyan, colors.Reset, typedInfo["channelTitle"])
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\u279C {{ .Name | cyan }} ({{ .Channel | red }})",
		Inactive: "{{ .Name | cyan }} ({{ .Channel | red }})",
		Selected: "{{ .Name | yellow }}",
	}

	prompt := promptui.Select{
		Label:     "Select Song",
		Items:     musicList,
		Templates: templates,
		Size:      5,
	}

	index, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	//fmt.Printf("You choose %q\n", result)
	//fmt.Println(index)
	return idMap, index + 1
}

func Command(params ...string) {

	// 0 -> command name
	// 1 -> video url
	// 2 -> tmp file path

	executable, err := exec.LookPath(params[0])
	if err != nil {
		fmt.Printf("Please install %s", params[0])
		log.Fatalln(err)
	}

	argList := []string{}

	seekTime := viper.GetString("SEEK_TIME")

	switch params[0] {
	case "youtube-dl":
		argList = []string{executable, "-q", "-x", "--audio-format", "mp3", params[1], "-o", params[2], "--no-continue"}
	case "ffplay":
		argList = []string{executable, params[2], "-nodisp", "-autoexit", "-ss", seekTime}
	}

	command := &exec.Cmd{
		Path:   executable,
		Args:   argList,
		Stdout: os.Stdout,
		Stdin:  os.Stdout,
	}

	seekTimeFloat, err := strconv.ParseFloat(seekTime, 64)
	PrintErr(err)

	if params[0] == "ffplay" {
		startTime := time.Now()
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			elapsed := time.Since(startTime)
			total := elapsed.Seconds() + seekTimeFloat
			seek := fmt.Sprintf("%f", total)
			viper.Set("SEEK_TIME", seek)
			viper.WriteConfig()
			os.Exit(1)
		}()
	}

	err = command.Run()
	PrintErr(err)

}
