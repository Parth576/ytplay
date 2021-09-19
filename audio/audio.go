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

func printTime(streamer beep.StreamSeekCloser, format beep.Format, songFinished chan bool) {
    onesecCtr := 0
	for {
		select {
		case <-time.After(time.Second):
            currTime := format.SampleRate.D(streamer.Position()).Round(time.Second)
            var d time.Duration = 1*time.Second
            if currTime == d {
                onesecCtr++
            }
            if onesecCtr == 2 {
                songFinished <- true
                fmt.Printf("\r")
                os.Exit(1)
            }
			fmt.Printf("\r%s/%s",format.SampleRate.D(streamer.Position()).Round(time.Second),format.SampleRate.D(streamer.Len()).Round(time.Second))
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

    songFinished := make(chan bool, 1)

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	speaker.Play(ctrl, beep.Callback(func() {
		go printTime(streamer, format, songFinished)
	}))

	fmt.Println("Press ENTER to pause/resume. ")
    // correct way to do channels but scanln blocks main thread so not channel input blocked as well
	//for {
    //    select {
    //    case <- songFinished:
    //        fmt.Println("song fin")
    //        speaker.Close()
    //        os.Exit(1)
    //    default: 
    //        fmt.Scanln()
    //        speaker.Lock()
    //        ctrl.Paused = !ctrl.Paused
    //        speaker.Unlock()
    //	}
    //}
    for {
        fmt.Scanln()
        speaker.Lock()
        ctrl.Paused = !ctrl.Paused
        speaker.Unlock()
    }
}
