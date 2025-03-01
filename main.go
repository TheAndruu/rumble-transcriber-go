package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"strconv"
	"regexp"
)

type DiarizationSegment struct {
    Start   float64
    End     float64
    Speaker string
}

type TranscriptionSegment struct {
    Start   float64
    End     float64
    Text    string
    Speaker string // To store the assigned speaker
}

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
	diarization := string(output)
	if diarization == "" {
		fmt.Println("Warning: Diarization output is empty")
	} else {
		fmt.Println("Diarization output:\n", diarization)
	}
	return diarization, nil
}

func transcribeAudio(audioPath, outputFile string) (string, error) {
	cmd := exec.Command("./whisper", "-f", audioPath, "-m", "ggml-base.en.bin", "-otxt", "-of", outputFile, "-t", "4")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("transcription failed: %v, output: %s", err, output)
	}
	fmt.Println("Transcription console output:\n", string(output))
	return string(output), nil
}

func parseDiarization(diarization string) ([]DiarizationSegment, error) {
	lines := strings.Split(diarization, "\n")
	var segments []DiarizationSegment
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "INFO:") || strings.HasPrefix(line, "/usr/local") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("Invalid diarization line: %s\n", line)
			continue
		}
		start, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}
		end, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		segments = append(segments, DiarizationSegment{Start: start, End: end, Speaker: parts[2]})
	}
	return segments, nil
}

func parseTranscription(transcription string) ([]TranscriptionSegment, error) {
	re := regexp.MustCompile(`\[(\d+:\d+:\d+\.\d+)\s*-->\s*(\d+:\d+:\d+\.\d+)\]\s*(.*)`)
	lines := strings.Split(transcription, "\n")
	var segments []TranscriptionSegment
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) != 4 {
			continue
		}
		start := parseTime(matches[1])
		end := parseTime(matches[2])
		text := strings.TrimSpace(matches[3])
		segments = append(segments, TranscriptionSegment{Start: start, End: end, Text: text})
	}
	return segments, nil
}

func parseTime(timeStr string) float64 {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0.0
	}
	hours, _ := strconv.ParseFloat(parts[0], 64)
	minutes, _ := strconv.ParseFloat(parts[1], 64)
	seconds, _ := strconv.ParseFloat(parts[2], 64)
	return hours*3600 + minutes*60 + seconds
}

func combineDiarizationAndTranscription(diarization, transcription string) {
	diarSegments, err := parseDiarization(diarization)
	if err != nil {
		fmt.Printf("Error parsing diarization: %v\n", err)
		return
	}
	transSegments, err := parseTranscription(transcription)
	if err != nil {
		fmt.Printf("Error parsing transcription: %v\n", err)
		return
	}

	// Assign each transcription segment to the diarization segment with the most overlap
	for i := range transSegments {
		maxOverlap := 0.0
		var bestSpeaker string
		for _, dSeg := range diarSegments {
			overlapStart := max(transSegments[i].Start, dSeg.Start)
			overlapEnd := min(transSegments[i].End, dSeg.End)
			overlap := overlapEnd - overlapStart
			if overlap > 0 && overlap > maxOverlap {
				maxOverlap = overlap
				bestSpeaker = dSeg.Speaker
			}
		}
		if bestSpeaker != "" {
			transSegments[i].Speaker = bestSpeaker
		} else {
			transSegments[i].Speaker = "Unknown"
		}
	}

	// Print the transcription with speakers
	fmt.Println("Transcription with Speakers:")
	for _, tSeg := range transSegments {
		fmt.Printf("%s (%.2f-%.2f): %s\n", tSeg.Speaker, tSeg.Start, tSeg.End, tSeg.Text)
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
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

	fmt.Println("Downloading video from Rumble...")
	if err := downloadRumbleVideo(rumbleURL, videoFile); err != nil {
		fmt.Printf("Error downloading video: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Extracting audio...")
	if err := extractAudio(videoFile, audioFile); err != nil {
		fmt.Printf("Error extracting audio: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Diarizing audio...")
	diarizationOutput, err := diarizeAudio(audioFile)
	if err != nil {
		fmt.Printf("Error diarizing audio: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Transcribing audio...")
	transcriptionOutput, err := transcribeAudio(audioFile, transcriptFile)
	if err != nil {
		fmt.Printf("Error transcribing audio: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nTranscription with Speakers:")
	combineDiarizationAndTranscription(diarizationOutput, transcriptionOutput)

	os.Remove(videoFile)
	os.Remove(audioFile)
	os.Remove(transcriptFile + ".txt")
}