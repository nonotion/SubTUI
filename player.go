package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/blang/mpv"
)

var (
	mpvClient *mpv.Client
	mpvCmd    *exec.Cmd
)

type PlayerStatus struct {
	Title    string
	Artist   string
	Album    string
	Current  float64
	Duration float64
	Paused   bool
}

func initPlayer() error {
	socketPath := "/tmp/depthtui_mpv_socket"

	args := []string{
		"--idle",
		"--no-video",
		"--input-ipc-server=" + socketPath,
	}

	mpvCmd = exec.Command("mpv", args...)
	if err := mpvCmd.Start(); err != nil {
		return fmt.Errorf("failed to start mpv: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	client := mpv.NewClient(mpv.NewIPCClient(socketPath))
	mpvClient = client

	return nil
}

func playSong(songID string) error {
	if mpvClient == nil {
		return fmt.Errorf("player not initialized")
	}

	url := subsonicStream(songID)
	if err := mpvClient.Loadfile(url, mpv.LoadFileModeReplace); err != nil {
		return err
	}

	mpvClient.SetProperty("pause", false)

	return nil
}

func shutdownPlayer() {
	if mpvCmd != nil {
		mpvCmd.Process.Kill()
	}
}

func getPlayerStatus() PlayerStatus {
	if mpvClient == nil {
		return PlayerStatus{}
	}

	getStr := func(prop string) string {
		val, err := mpvClient.GetProperty(prop)
		if err != nil {
			return ""
		}

		s := fmt.Sprintf("%v", val)

		if s == "<nil>" || s == "nil" || s == "" {
			return ""
		}
		return s
	}

	getFloat := func(prop string) float64 {
		val, err := mpvClient.GetProperty(prop)
		if err != nil {
			return 0.0
		}

		s := fmt.Sprintf("%v", val)

		if s == "<nil>" || s == "nil" {
			return 0.0
		}

		f, _ := strconv.ParseFloat(s, 64)
		return f
	}

	pausedVal, _ := mpvClient.GetProperty("pause")
	pausedStr := fmt.Sprintf("%v", pausedVal)
	isPaused := (pausedStr == "yes" || pausedStr == "true")

	return PlayerStatus{
		Title:    getStr("media-title"),
		Artist:   getStr("metadata/by-key/artist"),
		Album:    getStr("metadata/by-key/album"),
		Current:  getFloat("time-pos"),
		Duration: getFloat("duration"),
		Paused:   isPaused,
	}

}
