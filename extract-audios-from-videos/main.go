package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/hairyhenderson/go-which"
	"github.com/integrii/flaggy"
	"github.com/lmittmann/tint"
	"github.com/oriser/regroup"
)

var (
	logger *slog.Logger
	re     = regroup.MustCompile(`(?m)(?P<processed>\d+)kB`)
)

var (
	argVideosPath       string
	argFfmpegBinaryPath string
)

func parseArgs() {
	flaggy.SetName("extract-audios-from-videos")
	flaggy.SetDescription("A program that extracts audios from video files using ffmpeg")

	flaggy.String(&argVideosPath, "vp", "videos-path", "A path to the folder with videos")
	flaggy.String(&argFfmpegBinaryPath, "fbp", "ffmpeg-binary", "A path to the ffmpeg binary")

	flaggy.SetVersion("1.0")
	flaggy.Parse()
}

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

type Progress struct {
	Processed int64 `regroup:"processed"`
}

func (progress *Progress) Write(p []byte) (int, error) {
	n := len(p)

	err := re.MatchToTarget(string(p), progress)

	logger.Info("processed", "size", byteCountSI(progress.Processed*1024))

	return n, err
}

func checkFileExistence(path string) error {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return err
	}

	if fi.IsDir() {
		return errors.New("path is a directory")
	}

	if err != nil {
		return err
	}

	return nil
}

func checkDirExistence(path string) error {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return err
	}

	if !fi.IsDir() {
		return errors.New("path is not a directory")
	}

	if err != nil {
		return err
	}

	return nil
}

func initSlog() {
	w := os.Stderr

	logger = slog.New(tint.NewHandler(w, nil))

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}

func countFilesInPath(path string) (int, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	files := 0

	for _, e := range entries {
		if !e.IsDir() {
			files++
		}
	}

	return files, nil
}

func extractAudio(path string) error {
	var err error

	inFilePath := fmt.Sprintf("%s/%s", argVideosPath, path)

	renamedFilename := strings.ReplaceAll(path, filepath.Ext(path), ".wav")
	outFilePath := fmt.Sprintf("%s/%s", argVideosPath, renamedFilename)

	args := []string{
		"-y",
		"-loglevel",
		"quiet",
		"-stats",
		"-i",
		inFilePath,
		outFilePath,
	}

	cmd := fmt.Sprintf("%s %s", argFfmpegBinaryPath, strings.Join(args, " "))
	logger.Info("running command", "cmd", cmd)

	command := exec.Command(argFfmpegBinaryPath, args...)

	var ptyFile *os.File
	if ptyFile, err = pty.Start(command); err == nil {
		defer func(ptyFile *os.File) {
			err = ptyFile.Close()
			if err != nil {
				logger.Error("failed to close pty", "error", err)
			}
		}(ptyFile)
	}

	_, err = io.Copy(&Progress{}, ptyFile)

	return err
}

func doAudioExtraction(path string) error {
	logger.Info("extract audio part from files", "path", path)

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for n, e := range entries {
		if !e.IsDir() {
			logger.Info("processing a file", "num", n, "filename", e.Name())

			err = extractAudio(e.Name())
			if err != nil {
				logger.Warn("failed to extract audio", "err", err)
			}
		}
	}

	return nil
}

func main() {
	var err error

	initSlog()

	logger.Info("extract-audios-from-videos has been started")

	parseArgs()

	if argFfmpegBinaryPath == "" {
		logger.Info("check ffmpeg-binary existence by go-which, because the argument is empty")

		argFfmpegBinaryPath = which.Which("ffmpeg")
	} else {
		logger.Info("check ffmpeg-binary existence", "path", argFfmpegBinaryPath)
	}

	if err = checkFileExistence(argFfmpegBinaryPath); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else {
		logger.Info("ffmpeg-binary has been found", "path", argFfmpegBinaryPath)
	}

	logger.Info("check videos-path existence", "path", argVideosPath)

	if err = checkDirExistence(argVideosPath); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else {
		logger.Info("videos-path has been found", "path", argVideosPath)
	}

	nFiles, err := countFilesInPath(argVideosPath)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if nFiles == 0 {
		logger.Error("videos-path does not have videos")
		os.Exit(1)
	} else {
		logger.Info("videos-path has videos", "count", nFiles)
	}

	err = doAudioExtraction(argVideosPath)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("extract-audios-from-videos has been finished")
}
