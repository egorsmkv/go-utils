`extract-audios-from-videos`

> this little utility lets you ability to extract audio tracks from the folder with videos using ffmpeg

build:

    go build -ldflags="-s -w" -o extract-audios-from-videos main.go

usage:

    ./extract-audios-from-videos --videos-path ./videos

output:

    2024-05-14T21:23:01+03:00 INF extract-audios-from-videos has been started
    2024-05-14T21:23:01+03:00 INF check ffmpeg-binary existence by go-which, because the argument is empty
    2024-05-14T21:23:01+03:00 INF ffmpeg-binary has been found path=/usr/bin/ffmpeg
    2024-05-14T21:23:01+03:00 INF check videos-path existence path=./videos
    2024-05-14T21:23:01+03:00 INF videos-path has been found path=./videos
    2024-05-14T21:23:01+03:00 INF videos-path has videos count=1
    2024-05-14T21:23:01+03:00 INF extract audio part from files path=./videos
    2024-05-14T21:23:01+03:00 INF processing a file num=0 filename=1.webm
    2024-05-14T21:23:01+03:00 INF running command cmd="/usr/bin/ffmpeg -y -loglevel quiet -stats -i ./videos/1.webm ./videos/1.wav"
    2024-05-14T21:23:02+03:00 INF processed size="3.1 kB"
    2024-05-14T21:23:02+03:00 INF processed size="36.4 MB"
    2024-05-14T21:23:03+03:00 INF processed size="73.4 MB"
    2024-05-14T21:23:03+03:00 INF processed size="108.5 MB"
    2024-05-14T21:23:04+03:00 INF processed size="145.0 MB"
    2024-05-14T21:23:04+03:00 INF processed size="182.5 MB"
    2024-05-14T21:23:05+03:00 INF processed size="219.7 MB"
    2024-05-14T21:23:05+03:00 INF processed size="221.6 MB"
    2024-05-14T21:23:05+03:00 INF copied bytes size="569 B"
    2024-05-14T21:23:05+03:00 INF extract-audios-from-videos has been finished
