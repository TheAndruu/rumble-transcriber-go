package main

import (
	"fmt"
	"os"
	"os/exec"
)

func downloadRumbleVideo(url, outputPath string) error {
	cmd := exec.Command("yt-dlp", "-o", outputPath, "--merge-output-format", "mp4", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func extractAudio(videoPath, audioPath string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vn", "-acodec", "mp3", audioPath, "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func transcribeAudio(audioPath, outputTextPath string) error {
	// Assuming whisper.cpp main binary is available as 'whisper'
	cmd := exec.Command("./whisper", "-f", audioPath, "-m", "ggml-base.en.bin", "-otxt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <rumble_video_url>")
		os.Exit(1)
	}

	rumbleURL := os.Args[1]
	videoFile := "video.mp4"
	audioFile := "audio.mp3"
	transcriptFile := "audio.txt" // whisper.cpp outputs to <audio>.txt by default

	// Step 1: Download the video
	fmt.Println("Downloading video from Rumble...")
	if err := downloadRumbleVideo(rumbleURL, videoFile); err != nil {
		fmt.Printf("Error downloading video: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Extract audio
	fmt.Println("Extracting audio...")
	if err := extractAudio(videoFile, audioFile); err != nil {
		fmt.Printf("Error extracting audio: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Transcribe audio
	fmt.Println("Transcribing audio...")
	if err := transcribeAudio(audioFile, transcriptFile); err != nil {
		fmt.Printf("Error transcribing audio: %v\n", err)
		os.Exit(1)
	}

	// Step 4: Read and display the transcription
	transcription, err := os.ReadFile(transcriptFile)
	if err != nil {
		fmt.Printf("Error reading transcription: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nTranscription:")
	fmt.Println(string(transcription))

	// Clean up
	os.Remove(videoFile)
	os.Remove(audioFile)
	os.Remove(transcriptFile)
}