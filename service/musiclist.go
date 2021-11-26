package svc

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync/atomic"
)

func ReflushMusicFileAddr(musicfileaddr string) string {
	return musicfileaddr
}

func ReflushMusicList(musicpath string, musiclist *[]*MusicEntry) (total int32, err error) {
	*musiclist = (*musiclist)[0:0]
	rd, err := ioutil.ReadDir(musicpath)
	if err != nil {
		fmt.Println("read dir fail:", err)
		panic(err)
	}
	for index, fi := range rd {
		suffix := path.Ext(fi.Name())
		if strings.Compare(suffix, ".mp3") == 0 {
			prefix := strings.TrimSuffix(fi.Name(), suffix)
			prefixSplit := strings.SplitN(prefix, " - ", 2)
			var m *MusicEntry
			if len(prefixSplit) < 2 {
				log.Println(fi.Name(), "Name is not format! Will set name to filename")
				m = &MusicEntry{
					Id:   index,
					Name: prefixSplit[0],
					Path: fmt.Sprint(musicpath, fi.Name()),
				}
			} else {
				m = &MusicEntry{
					Id:     index,
					Name:   prefixSplit[1],
					Artist: prefixSplit[0],
					Path:   fmt.Sprint(musicpath, fi.Name()),
				}
			}
			atomic.AddInt32(&total, 1)
			*musiclist = append(*musiclist, m)
		}
	}
	return total, err
}
