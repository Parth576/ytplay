package audio

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func printTime(streamer beep.StreamSeekCloser, format beep.Format) {
	for {
		select {
		case <-time.After(time.Second):
			fmt.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
		}

	}
}

func Play(fp string) {
	f, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	speaker.Play(ctrl, beep.Callback(func() {
		go printTime(streamer, format)
	}))

	fmt.Println(format.SampleRate.D(streamer.Len()).Round(time.Second))

	fmt.Print("Press any key to pause/resume. ")
	for {
		fmt.Scanln()

		speaker.Lock()
		ctrl.Paused = !ctrl.Paused
		speaker.Unlock()
	}
}
