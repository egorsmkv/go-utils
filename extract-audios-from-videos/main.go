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

func checkExistence(path string, checkDir bool) error {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return err
	}

	if checkDir {
		if !fi.IsDir() {
			return errors.New("path is a directory")
		}
	} else {
		if fi.IsDir() {
			return errors.New("path is not a directory")
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func initSlog() {
	logger = slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC3339,
		}),
	)
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
	var (
		err     error
		ptyFile *os.File
	)

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

	if ptyFile, err = pty.Start(command); err == nil {
		defer func(ptyFile *os.File) {
			err = ptyFile.Close()
			if err != nil {
				logger.Error("failed to close pty", "err", err)
			}
		}(ptyFile)
	}

	written, err := io.Copy(&Progress{}, ptyFile)

	logger.Info("copied bytes", "size", byteCountSI(written))

	if err != nil {
		if strings.Contains(err.Error(), "input/output error") {
			return nil
		}
	}

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
				logger.Error("failed to extract audio", "err", err)
			}
		}
	}

	return nil
}

func main() {
	initSlog()

	logger.Info("extract-audios-from-videos has been started")

	parseArgs()

	if argFfmpegBinaryPath == "" {
		logger.Info("check ffmpeg-binary existence by go-which, because the argument is empty")

		argFfmpegBinaryPath = which.Which("ffmpeg")
	} else {
		logger.Info("check ffmpeg-binary existence", "path", argFfmpegBinaryPath)
	}

	if err := checkExistence(argFfmpegBinaryPath, false); err != nil {
		logger.Error("cannot check ffmpeg-binary", "err", err)
	} else {
		logger.Info("ffmpeg-binary has been found", "path", argFfmpegBinaryPath)
	}

	logger.Info("check videos-path existence", "path", argVideosPath)

	if err := checkExistence(argVideosPath, true); err != nil {
		logger.Error("cannot check videos-path", "err", err)
		os.Exit(1)
	} else {
		logger.Info("videos-path has been found", "path", argVideosPath)
	}

	nFiles, err := countFilesInPath(argVideosPath)
	if err != nil {
		logger.Error("cannot count files in path", "err", err)
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
		logger.Error("cannot extract audio", "err", err)
		os.Exit(1)
	}

	logger.Info("extract-audios-from-videos has been finished")
}
