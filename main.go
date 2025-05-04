package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
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
		http.Error(w, "Image not provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Invalid image format", http.StatusBadRequest)
		return
	}

	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	// Optional logging to file
	logFile, _ := os.Create("border_log.txt")
	defer logFile.Close()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if isBlack(img.At(x, y)) {
				logFile.WriteString(fmt.Sprintf("black pixel at: (%d, %d)\n", x, y))
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// Crop image
	rect := image.Rect(minX, minY, maxX+1, maxY+1)
	cropped := image.NewRGBA(rect)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			cropped.Set(x, y, img.At(x, y))
		}
	}

	// Save the output
	outFile, err := os.Create("output.png")
	if err != nil {
		ReturnBadRequest(w, err, "Failed to create output file")
	}
	defer outFile.Close()

	// Send image as response
	if err := png.Encode(outFile, cropped); err != nil {
		ReturnErr(w, err, "Failed to encode output")
	}

	ReturnOK(w, "Successfully cropping the image. Saved as 'output.png'.")
}
