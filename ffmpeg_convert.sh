#!/usr/bin/env bash

ffmpeg -v verbose -i rtmp://35.193.201.151:1935/test \
-vcodec libx264 -acodec aac -ac 1 -strict -2 -crf 18 \
-profile:v baseline -maxrate 400k -bufsize 1835k \
-pix_fmt yuv420p -flags -global_header \
-hls_time 6 -hls_list_size 6 -hls_wrap 4 -hls_flags delete_segments \
-start_number 1 ./hls/index1.m3u8 \
-vcodec libx264 -acodec aac -ac 1 -strict -2 -crf 18 \
-profile:v baseline -maxrate 700k -bufsize 1835k \
-pix_fmt yuv420p -flags -global_header \
-hls_time 6 -hls_list_size 6 -hls_wrap 4 -hls_flags delete_segments \
-start_number 1 ./hls/index2.m3u8