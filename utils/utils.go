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

	"github.com/Parth576/ytplay/colors"
	"github.com/spf13/viper"
)

// PrintErr prints the error to logs
func PrintErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func PrettyPrint(items interface{}) map[int]string {
	idMap := make(map[int]string)
	for index, v := range items.([]interface{}) {
		videoMap := v.(map[string]interface{})
		id := videoMap["id"]
		idMap[index+1] = id.(map[string]interface{})["videoId"].(string)
		info := videoMap["snippet"]
		typedInfo := info.(map[string]interface{})
		fmt.Printf("%v)  %sTitle:%s %s\n", index+1, colors.Cyan, colors.Reset, typedInfo["title"])
		fmt.Printf("    %sChannel:%s %s\n\n", colors.Cyan, colors.Reset, typedInfo["channelTitle"])
	}
	return idMap
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
		argList = []string{executable, "-x", "--audio-format", "mp3", params[1], "-o", params[2], "--no-continue"}
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
