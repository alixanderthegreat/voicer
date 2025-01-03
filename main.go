package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {

	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		return
	}
	defer r.Close()
	defer w.Close()

	cmd := exec.Command("ffmpeg",
		"-vn",
		"-f", "alsa",
		"-i", "default",
		"output.wav", "-y",
		"-ar", "44100",
		"-ac", "2",
		// "-filter_complex", "asetrate=44100*0.89,aresample=44100,atempo=1.176",
	)

	// pass in the pipe writer
	cmd.Stdout = w
	cmd.Stderr = w

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		store_output := make([]byte, 1028)
		for {
			n, err := r.Read(store_output)
			if err != nil && err != io.EOF {
				fmt.Println(err)
				return
			}
			if n == 0 {
				break
			}
			fmt.Println(string(store_output[:n]))
		}
	}()
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
