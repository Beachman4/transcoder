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
    <source type="application/x-mpegURL" src="http://35.238.243.208:8080/hls/test/index.m3u8">
</video>


<!-- JS code -->
<!-- If you'd like to support IE8 (for Video.js versions prior to v7) -->
<script src="https://vjs.zencdn.net/ie8/ie8-version/videojs-ie8.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/videojs-contrib-hls/5.14.1/videojs-contrib-hls.js"></script>
<script src="https://vjs.zencdn.net/7.2.3/video.js"></script>

<script>
    var player = videojs('video');
    player.ready(function() {
        var myPlayer = this;

        const urlParams = new URLSearchParams(window.location.search);
        const stream = urlParams.get('stream');

        if (stream) {
            const url = "http://35.238.243.208:8080/hls/" + stream + "/index.m3u8"

            myPlayer.src({type: "application/x-mpegURL", src: url})
        }
    })
    player.play();
</script>

</body>
</html>`

type Info struct {
	Cmd *exec.Cmd
	running bool
}

func main() {
	var mappedTranscoding = map[string]*Info{}

	router := httprouter.New()
	router.ServeFiles("/hls/*filepath", http.Dir("hls"))
	router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprintf(writer, indexHtml)
	})
	router.GET("/healthz", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(writer, "Healthy!\n")
	})
	router.GET("/stop-transcoding/:key", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		mappedTranscoding[params.ByName("key")].Cmd.Process.Kill()

		delete(mappedTranscoding, params.ByName("key"))
	})
	router.GET("/start-transcoding/:key", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		key := params.ByName("key")

		if mappedTranscoding[key] != nil {
			return
		}

		go func(streamKey string) {
			pathHls := path.Join([]string{"hls", streamKey}...)

			os.MkdirAll(pathHls, os.ModePerm)

			copyIndex(path.Join("hls", "index.m3u8"), path.Join(pathHls, "index.m3u8"), streamKey)

			cmd := exec.Command("/usr/local/bin/ffmpeg", strings.Split(fmt.Sprintf(`-v verbose -re -i rtmp://35.193.201.151:1935/%s -vcodec libx264 -acodec aac -ac 1 -strict -2 -crf 18 -profile:v baseline -preset veryfast -maxrate 400k -bufsize 1835k -pix_fmt yuv420p -flags -global_header -hls_time 6 -hls_list_size 6 -start_number 1 ./hls/%s/index1.m3u8 -vcodec libx264 -acodec aac -preset veryfast -ac 1 -strict -2 -crf 18 -profile:v baseline -maxrate 3000k -bufsize 1835k -pix_fmt yuv420p -flags -global_header -hls_time 6 -hls_list_size 6 -start_number 1 ./hls/%s/index2.m3u8`, streamKey, streamKey, streamKey), " ")...)

			mappedTranscoding[streamKey] = &Info{
				Cmd: cmd,
			}

			output, err := cmd.CombinedOutput()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(output))

			names, err := ioutil.ReadDir(pathHls)
			if err != nil {
				fmt.Println(err)
			}
			for _, entery := range names {
				os.Remove(path.Join([]string{"hls", streamKey, entery.Name()}...))
			}
		}(key)
	})

	handler := cors.Default().Handler(router)

	http.ListenAndServe(":8080", handler)
}

func copyIndex(src, dest, key string) {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	contents := string(input)

	contents = strings.Replace(contents, "http://35.238.243.208:8080/hls/index1.m3u8", fmt.Sprintf("http://35.238.243.208:8080/hls/%s/index1.m3u8", key), -1)
	contents = strings.Replace(contents, "http://35.238.243.208:8080/hls/index2.m3u8", fmt.Sprintf("http://35.238.243.208:8080/hls/%s/index2.m3u8", key), -1)

	err = ioutil.WriteFile(dest, []byte(contents), 0644)
	if err != nil {
		fmt.Println("Error creating", dest)
		fmt.Println(err)
		return
	}
}