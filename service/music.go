package svc

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type MusicEntry struct {
	Id         int
	Name       string
	Artist     string
	Path       string
	FileStream *os.File
}

func (me *MusicEntry) Open() {
	var err error
	me.FileStream, err = os.Open(me.Path)
	if err != nil {
		log.Println("Fail to open: ", err)
		return
	}
}

func (me *MusicEntry) Close() {
	err := me.FileStream.Close()
	if err != nil {
		log.Println("Fail to close: ", err)
		return
	}
}

func (me *MusicEntry) Play() {
	fileStream, err := os.Open(me.Path)
	if err != nil {
		log.Println("Fail to open: ", err)
		return
	}
	defer fileStream.Close()
	me.Open()
	defer me.Close()
	streamer, format, err := mp3.Decode(fileStream)
	if err != nil {
		log.Println("Fail to play: ", err)
		return
	}
	defer streamer.Close()
	_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	log.Println("Playing: ", me.Id, "| ", me.Name)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
