package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	svc "player/service"

	"github.com/gorilla/mux"
)

var musiclist []*svc.MusicEntry
var musicfileaddr string = "C:/TongyiWang/music/"

func NewHTTPHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/player", Hello).Methods("GET")
	r.HandleFunc("/player/getlist", GetMusicList).Methods("GET")
	r.HandleFunc("/player/flushlist", ReflushMusicList).Methods("GET")
	r.HandleFunc("/player/flushaddr", ReflushMusicFileAddr).Methods("POST")
	r.HandleFunc("/player/{musicid:[0-9]+}", MusicPlay).Methods("GET")
	return r
}

func Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello TY player\r\n"))
}

func GetMusicList(w http.ResponseWriter, _ *http.Request) {
	var statusCode int = 200
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	for i, music := range musiclist {
		index := i + 1
		var info string
		if music.Artist == "" {
			info = fmt.Sprintf(music.Name)
		} else {
			info = fmt.Sprintf("%s - %s", music.Artist, music.Name)
		}
		json.NewEncoder(w).Encode(map[int]string{index: info})
	}
}

func ReflushMusicList(w http.ResponseWriter, _ *http.Request) {
	var statusCode int = 200
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	svc.ReflushMusicList(musicfileaddr, &musiclist)
	for i, music := range musiclist {
		index := i + 1
		var info string
		if music.Artist == "" {
			info = fmt.Sprintf(music.Name)
		} else {
			info = fmt.Sprintf("%s - %s", music.Artist, music.Name)
		}
		json.NewEncoder(w).Encode(map[int]string{index: info})
	}
}

func ReflushMusicFileAddr(w http.ResponseWriter, r *http.Request) {
	var statusCode int = 200
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	type request struct {
		Addr string `json:"addr"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"errno": 1, "errmsg": err})
		return
	}
	musicfileaddr = svc.ReflushMusicFileAddr(req.Addr)
	w.Write([]byte("Change music file success\r\n"))
}

func MusicPlay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	edgeno := vars["musicid"]
	var statusCode int = 200
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	edgeno_i, err := strconv.Atoi(edgeno)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"errno": 1, "errmsg": err})
		return
	}
	music := *musiclist[edgeno_i-1]
	// music.Open()
	// defer music.Close()
	music.Play()
}
