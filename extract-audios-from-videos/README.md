`extract-audios-from-videos`

> this little utility lets you ability to extract audio tracks from the folder with videos using ffmpeg

build:

    go build -ldflags="-s -w" -o extract-audios-from-videos main.go

usage:

    ./extract-audios-from-videos --videos-path ./videos

output:

    May 14 18:28:14.783 INF extract-audios-from-videos has been started
    May 14 18:28:14.784 INF check ffmpeg-binary existence by go-which, because the argument is empty
    May 14 18:28:14.784 INF ffmpeg-binary has been found path=/usr/bin/ffmpeg
    May 14 18:28:14.784 INF check videos-path existence path=./videos
    May 14 18:28:14.784 INF videos-path has been found path=./videos
    May 14 18:28:14.784 INF videos-path has videos count=2
    May 14 18:28:14.784 INF extract audio part from files path=./videos
    May 14 18:28:14.784 INF processing a file num=0 filename=1.webm
    May 14 18:28:14.784 INF running command cmd="/usr/bin/ffmpeg -y -loglevel quiet -stats -i ./videos/1.webm ./videos/1.wav"
    May 14 18:28:14.847 INF processed size="3.1 kB"
    May 14 18:28:15.346 INF processed size="35.9 MB"
    May 14 18:28:15.846 INF processed size="71.3 MB"
    May 14 18:28:16.346 INF processed size="107.5 MB"
    May 14 18:28:16.846 INF processed size="143.9 MB"
    May 14 18:28:17.346 INF processed size="180.1 MB"
    May 14 18:28:17.364 INF processed size="181.1 MB"
    May 14 18:28:17.367 WRN failed to extract audio err="read /dev/ptmx: input/output error"
    May 14 18:28:17.367 INF processing a file num=1 filename=2.webm
    May 14 18:28:17.367 INF running command cmd="/usr/bin/ffmpeg -y -loglevel quiet -stats -i ./videos/2.webm ./videos/2.wav"
    May 14 18:28:17.418 INF processed size="3.1 kB"
    May 14 18:28:17.917 INF processed size="36.4 MB"
    May 14 18:28:18.417 INF processed size="73.7 MB"
    May 14 18:28:18.917 INF processed size="110.6 MB"
    May 14 18:28:19.417 INF processed size="147.8 MB"
    May 14 18:28:19.917 INF processed size="185.3 MB"
    May 14 18:28:20.417 INF processed size="217.8 MB"
    May 14 18:28:20.493 INF processed size="221.6 MB"
    May 14 18:28:20.497 WRN failed to extract audio err="read /dev/ptmx: input/output error"
    May 14 18:28:20.497 INF extract-audios-from-videos has been finished
