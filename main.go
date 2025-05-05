package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"

	resp "github.com/mahendra2448/tech-test/helper"
)

func main() {
	port := 8070
	http.HandleFunc("/crop", cropImage)

	log.Default().Printf("Server running at port: %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("[error-init-http-server] %v", err)
	}
}

func isBlack(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0 && g == 0 && b == 0
}

func cropImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get image file
	file, _, err := r.FormFile("image")
	if err != nil {
		resp.ReturnBadRequest(w, err, "Image not provided")
		return
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		resp.ReturnBadRequest(w, err, "Invalid image format")
		return
	}

	// Crop image
	minX, minY, maxX, maxY := findBorder(img)
	rect := image.Rect(minX, minY, maxX+1, maxY+1)

	cropped := image.NewRGBA(rect)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			cropped.Set(x, y, img.At(x, y))
		}
	}
	fmt.Println("FINAL RECTANGLE:", rect)

	// Save the output
	outFile, err := os.Create("output.png")
	if err != nil {
		resp.ReturnBadRequest(w, err, "Failed to create output file")
	}
	defer outFile.Close()

	// Send image as response
	if err := png.Encode(outFile, cropped); err != nil {
		resp.ReturnErr(w, err, "Failed to encode output")
	}

	resp.ReturnOK(w, "Successfully cropping the image. Saved as 'output.png'.")
}

func findBorder(img image.Image) (int, int, int, int) {
	// logging to file
	logFile, _ := os.Create("border_log.txt")
	defer logFile.Close()

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var minX, maxX, minY, maxY int
	found := false

	// Top border
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isBlack(img.At(x, y)) {
				minY = y
				found = true
				_, _ = logFile.WriteString(fmt.Sprintf("black pixel top found at: (%d, %d)\n", x, y))
				break
			}
		}
		if found {
			break
		}
	}

	// Bottom border
	found = false
	maxY = height / 2
	for y := minY + 1; y < maxY; y++ {
		for x := minX; x < width; x++ {
			if isBlack(img.At(x, y)) {
				maxY += y
				found = true
				_, _ = logFile.WriteString(fmt.Sprintf("black pixel bottom found at: (%d, %d)\n", x, y))
				break
			}
			maxY++
		}

		if found {
			break
		}
	}

	// Left border
	found = false
	for x := 0; x < width; x++ {
		for y := minY; y <= maxY; y++ {
			if isBlack(img.At(x, y)) {
				minX = x
				found = true
				_, _ = logFile.WriteString(fmt.Sprintf("black pixel left found at: (%d, %d)\n", x, y))
				break
			}
		}
		if found {
			break
		}
	}

	// Right border
	found = false
	for x := width - 1; x >= 0; x-- {
		for y := minY; y <= maxY; y++ {
			if isBlack(img.At(x, y)) {
				maxX = x
				found = true
				_, _ = logFile.WriteString(fmt.Sprintf("black pixel right found at: (%d, %d)\n", x, y))
				break
			}
		}
		if found {
			break
		}
	}

	return minX, minY, maxX, maxY
}
