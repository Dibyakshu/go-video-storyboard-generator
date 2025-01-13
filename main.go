package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		showUsage()
	}

	fps, inputFile, outputDir := parseArguments()

	ensureDirectoryExists(outputDir)

	startTime := time.Now()
	
	thumbnailsDir := filepath.Join(outputDir, "thumbnails")
	ensureDirectoryExists(thumbnailsDir)

	generateThumbnails(fps, inputFile, thumbnailsDir)

	storyboardImage := filepath.Join(outputDir, "storyboard.jpg")
	generateStoryboard(thumbnailsDir, storyboardImage)

	vttFile := filepath.Join(outputDir, "storyboard.vtt")
	generateVTT(fps, thumbnailsDir, vttFile)

	cleanup(thumbnailsDir)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Process completed in %v.\n", elapsedTime)
}

func showUsage() {
	fmt.Println("Usage: <executable> <fps> <input_file> <output_dir>")
	fmt.Println("./storyboard 1 example.mp4 ./output")
	os.Exit(1)
}

func parseArguments() (int, string, string) {
	fps, err := strconv.Atoi(os.Args[1])
	if err != nil || fps <= 0 {
		fmt.Println("Error: Invalid FPS value. Please provide a positive integer.")
		os.Exit(1)
	}
	return fps, os.Args[2], os.Args[3]
}

func ensureDirectoryExists(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", path, err)
		os.Exit(1)
	}
}

func generateThumbnails(fps int, inputFile, thumbnailsDir string) {
	fmt.Println("Generating thumbnails...")
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vf", fmt.Sprintf("fps=1/%d,scale=384:160", fps), 
		filepath.Join(thumbnailsDir, "thumb%05d.jpg")) // Zero-padded to 5 digits
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error generating thumbnails: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Thumbnails generated successfully.")
}

func generateStoryboard(thumbnailsDir, storyboardImage string) {
	fmt.Println("Creating storyboard.jpg...")
	thumbnails, err := filepath.Glob(filepath.Join(thumbnailsDir, "thumb*.jpg"))
	if err != nil {
		fmt.Printf("Error reading thumbnails: %v\n", err)
		os.Exit(1)
	}
	if len(thumbnails) == 0 {
		fmt.Println("No thumbnails found. Exiting.")
		os.Exit(1)
	}

	columns := 10
	rows := (len(thumbnails) + columns - 1) / columns
	tileFilter := fmt.Sprintf("tile=%dx%d", columns, rows)

	cmd := exec.Command("ffmpeg", "-pattern_type", "glob", "-i", filepath.Join(thumbnailsDir, "*.jpg"), "-filter_complex", tileFilter, storyboardImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error creating storyboard image: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Storyboard image created successfully.")
}

func generateVTT(fps int, thumbnailsDir, vttFile string) {
	fmt.Println("Generating storyboard.vtt...")
	thumbnails, err := filepath.Glob(filepath.Join(thumbnailsDir, "thumb*.jpg"))
	if err != nil {
		fmt.Printf("Error reading thumbnails: %v\n", err)
		os.Exit(1)
	}

	durationPerThumb := 1.0 * float64(fps)
	vttContent := "WEBVTT\n\n"
	for i := range thumbnails {
		start := time.Duration(float64(i) * durationPerThumb * float64(time.Second))
		end := start + time.Duration(durationPerThumb*float64(time.Second))

		x := (i % 10) * 384
		y := (i / 10) * 160
		vttContent += fmt.Sprintf("%s --> %s\n", formatDuration(start), formatDuration(end))
		vttContent += fmt.Sprintf("https://insertlinkhere.com/storyboard.jpg#xywh=%d,%d,384,160\n\n", x, y)
	}

	if err := os.WriteFile(vttFile, []byte(vttContent), 0644); err != nil {
		fmt.Printf("Error writing VTT file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Storyboard VTT file generated successfully.")
}

func cleanup(thumbnailsDir string) {
	fmt.Println("Cleaning up temporary files...")
	if err := os.RemoveAll(thumbnailsDir); err != nil {
		fmt.Printf("Error cleaning up thumbnails: %v\n", err)
		os.Exit(1)
	}
}

func formatDuration(d time.Duration) string {
	seconds := d.Seconds()
	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	s := int(seconds) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
