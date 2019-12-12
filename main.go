package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

var indexHtml = `<html>
	<head>
		<link href="https://vjs.zencdn.net/7.2.3/video-js.css" rel="stylesheet">
	</head>
	<body>
		
		<video id='video'  class="video-js vjs-default-skin" width="1280" height="720" controls>
			<source type="application/x-mpegURL" src="http://35.238.243.208:8080/hls/index.m3u8">
		</video>


		<!-- JS code -->
		<!-- If you'd like to support IE8 (for Video.js versions prior to v7) -->
		<script src="https://vjs.zencdn.net/ie8/ie8-version/videojs-ie8.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/videojs-contrib-hls/5.14.1/videojs-contrib-hls.js"></script>
		<script src="https://vjs.zencdn.net/7.2.3/video.js"></script>

		<script>
		var player = videojs('video');
		player.play();
		</script>
		
	</body>
</html>`

func main() {
	running := false
	var cmd *exec.Cmd;
	router := httprouter.New()
	router.ServeFiles("/hls/*filepath", http.Dir("hls"))
	router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprintf(writer, indexHtml)
	})
	router.GET("/healthz", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(writer, "Healthy!\n")
	})
	router.GET("/stop-transcoding", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		cmd.Process.Kill()
	})
	router.GET("/start-transcoding", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		if running {
			fmt.Fprintf(writer, "Still running")
			return
		}

		go func() {
			cmd = exec.Command("/usr/local/bin/ffmpeg", strings.Split(`-v verbose -i rtmp://35.193.201.151:1935/test -vcodec libx264 -acodec aac -ac 1 -strict -2 -crf 18 -profile:v baseline -maxrate 400k -bufsize 1835k -pix_fmt yuv420p -flags -global_header -hls_time 6 -hls_list_size 6 -start_number 1 ./hls/index1.m3u8 -vcodec libx264 -acodec aac -ac 1 -strict -2 -crf 18 -profile:v baseline -maxrate 700k -bufsize 1835k -pix_fmt yuv420p -flags -global_header -hls_time 6 -hls_list_size 6a -start_number 1 ./hls/index2.m3u8`, " ")...)

			running = true;

			output, err := cmd.CombinedOutput()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(output))

			running = false;

			names, err := ioutil.ReadDir("hls")
			if err != nil {
				fmt.Println(err)
			}
			for _, entery := range names {
				if entery.Name() != "index.m3u8" {
					os.Remove(path.Join([]string{"hls", entery.Name()}...))
				}
			}
		}()
	})

	handler := cors.Default().Handler(router)

	http.ListenAndServe(":8080", handler)
}
