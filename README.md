# Video Storyboard or Seek Thumbnails Generator

This Go application generates video storyboards by extracting thumbnails from a given video file. It uses FFmpeg to create the thumbnails and assemble them into a storyboard image. The program also generates a VTT file for use with the storyboard, allowing each thumbnail to be linked with the corresponding time in the video.

## Features

- Extracts thumbnails from a video at a specified frame rate (FPS).
- Generates a storyboard image by arranging thumbnails in a grid.
- Creates a WebVTT file that links each thumbnail to its timestamp in the video.
- Cleans up temporary thumbnail files after the process is completed.

## Prerequisites

- [Go](https://golang.org/dl/) version 1.18 or higher.
- [FFmpeg](https://ffmpeg.org/download.html) installed and available in your system's `PATH`.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/storyboard-generator.git
   cd storyboard-generator
   ```

2. Ensure Go and FFmpeg are properly installed on your system.

3. Build the Go executable:

   ```bash
   go build -o storyboard
   ```

## Usage

To use the program, run it with the following arguments:

```bash
./storyboard <fps> <input_file> <output_dir>
```

### Arguments:

- `<fps>`: The frame rate (FPS) for extracting thumbnails from the video (e.g., 1 for one thumbnail per second).
- `<input_file>`: The path to the video file (e.g., `example.mp4`).
- `<output_dir>`: The directory where the storyboard image and VTT file will be saved (e.g., `./output`).

### Example:

```bash
./storyboard 1 example.mp4 ./output
```

This will:

- Extract thumbnails from `example.mp4` at a rate of 1 frame per second.
- Save the thumbnails and generated files in the `./output` directory.
- Create a storyboard image named `storyboard.jpg`.
- Create a VTT file named `storyboard.vtt`.

## Output Files

- **storyboard.jpg**: A single image file containing all the thumbnails arranged in a grid.
- **storyboard.vtt**: A WebVTT file that links each thumbnail to its timestamp in the video.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
