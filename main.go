package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	running := false
	router := httprouter.New()
	router.ServeFiles("/hls/*filepath", http.Dir("hls"))
	router.GET("/healthz", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(writer, "Healthy!\n")
	})
	router.GET("/start-transcoding", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		go func() {
			cmd := exec.Command("ffmpeg", strings.Split(`-v verbose -i rtmp://35.193.201.151:1935/test -vcodec libx264 -acodec aac -ac 1 -strict -2 -crf 18 -profile:v baseline -maxrate 400k -bufsize 1835k -pix_fmt yuv420p -flags -global_header -hls_time 6 -hls_list_size 6 -hls_wrap 4 -start_number 1 ./hls/index1.m3u8 -vcodec libx264 -acodec aac -ac 1 -strict -2 -crf 18 -profile:v baseline -maxrate 700k -bufsize 1835k -pix_fmt yuv420p -flags -global_header -hls_time 6 -hls_list_size 6 -hls_wrap 4 -start_number 1 ./hls/index2.m3u8`, " ")...)

			running = true;

			output, err := cmd.CombinedOutput()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(output))

			running = false;
		}()
	})

	http.ListenAndServe(":8080", router)
}
