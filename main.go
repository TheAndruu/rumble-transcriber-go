package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func diarizeAudio(audioPath string) (string, error) {
    cmd := exec.Command("python3", "diarize.py", audioPath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("diarization failed: %v, output: %s", err, output)
    }
	fmt.Println("Diarization output:\n", string(output)) // Debug print
    return string(output), nil
}

func transcribeAudio(audioPath, outputFile string) error {
    cmd := exec.Command("./whisper", "-f", audioPath, "-m", "ggml-base.en.bin", "-otxt", "-of", outputFile, "-t")
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
	transcriptFile := "transcription"

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

	// Step 3: Diarize audio to identify speakers
	fmt.Println("Diarizing audio...")
	diarizationOutput, err := diarizeAudio(audioFile)
	if err != nil {
		fmt.Printf("Error diarizing audio: %v\n", err)
		os.Exit(1)
	}

	// Step 4: Transcribe audio
	fmt.Println("Transcribing audio...")
	if err := transcribeAudio(audioFile, transcriptFile); err != nil {
		fmt.Printf("Error transcribing audio: %v\n", err)
		os.Exit(1)
	}

	// Step 5: Read transcription
	transcriptPath := transcriptFile + ".txt"
	transcription, err := os.ReadFile(transcriptPath)
	if err != nil {
		fmt.Printf("Error reading transcription: %v\n", err)
		os.Exit(1)
	}

	// Step 6: Combine diarization and transcription
	fmt.Println("\nTranscription with Speakers:")
	combineDiarizationAndTranscription(diarizationOutput, string(transcription))

	// Clean up
	os.Remove(videoFile)
	os.Remove(audioFile)
	os.Remove(transcriptPath)
}

func combineDiarizationAndTranscription(diarization, transcription string) {
	diarizationLines := strings.Split(diarization, "\n")
	transcriptionLines := strings.Split(transcription, "\n")

	for i, dLine := range diarizationLines {
		if i < len(transcriptionLines) && dLine != "" {
			parts := strings.Fields(dLine)
			if len(parts) >= 3 {
				speaker := parts[2]
				fmt.Printf("%s: %s\n", speaker, transcriptionLines[i])
			}
		}
	}
}