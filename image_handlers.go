package main

import (
	"archive/zip"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func ConvertDicomToPNG(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")

	dataset, err := dicom.ParseFile(serverStorageDir+fileName, nil)
	if err != nil {
		http.Error(w, "Unable to read DICOM file.", http.StatusBadRequest)
		return
	}

	pixelData, err := getPixelData(dataset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"images.zip\"")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	for index, frame := range pixelData.Frames {
		image, err := frame.GetImage()
		if err != nil {
			fmt.Printf("Error while getting image, for index: %d %v.\n", index, err)
			continue
		}

		// Create a new entry in the ZIP archive for this image
		zipEntry, err := zipWriter.Create(getImageName(fileName, index))
		if err != nil {
			fmt.Printf("Error creating ZIP entry for %s: %v", fileName, err)
			continue
		}

		grayImg := createContrastInImage(image)
		// Encode the image as PNG and write it to the ZIP entry
		if err := png.Encode(zipEntry, grayImg); err != nil {
			fmt.Printf("Failed to encode to png, for index: %d %v.\n", index, err)
		}
	}
}

func getImageName(name string, index int) string {
	return fmt.Sprintf("image_%s_%d.png", name, index)
}

func getPixelData(dataset dicom.Dataset) (dicom.PixelDataInfo, error) {
	elem, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		return dicom.PixelDataInfo{}, fmt.Errorf("no pixel data found")
	}
	pixelData := dicom.MustGetPixelDataInfo(elem.Value)
	if len(pixelData.Frames) == 0 {
		return dicom.PixelDataInfo{}, fmt.Errorf("no pixel frame found")
	}
	return pixelData, nil
}

// For now, we don't need to store the png file. This is for debug purposes.
func storePngFile(filePath string, index int, img image.Image) {
	name := fmt.Sprintf(serverStorageDir+"png/%s%d%s", filePath, index, ".png")
	// better file name, in case of clash.
	file, err := os.Create(name)
	if err != nil {
		fmt.Printf("Error while creating file: %s.", err.Error())
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		fmt.Printf("Failed to encode to PNG: %v.", err)
	}

	if err = file.Close(); err != nil {
		fmt.Printf("Unable to properly close file: %v.", file.Name())
	}
}

func createContrastInImage(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray8 := image.NewGray(bounds)

	// Find min and max intensity
	var min, max uint16 = 65535, 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray16 := color.Gray16Model.Convert(img.At(x, y)).(color.Gray16)

			if gray16.Y < min {
				min = gray16.Y
			}
			if gray16.Y > max {
				max = gray16.Y
			}
		}
	}

	// Stretch contrast by windowing
	windowCenter := (max + min) / 2
	windowWidth := max - min

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray16 := color.Gray16Model.Convert(img.At(x, y)).(color.Gray16)
			// Apply windowing formula: scale to 0-255 range
			scaled := uint8(((int(gray16.Y) - int(windowCenter)) * 255 / int(windowWidth)) + 128)

			gray8.SetGray(x, y, color.Gray{Y: scaled})
		}
	}

	return gray8
}

// Debug: This function is for debugging purposes only.
func checkImageContrast(image image.Image) {
	// Check image bounds
	bounds := image.Bounds()
	log.Printf("Image bounds: %v", bounds)

	// Inspect some pixel values
	for y := bounds.Min.Y; y < bounds.Max.Y; y += bounds.Dy() / 10 { // Sample 10 rows
		for x := bounds.Min.X; x < bounds.Max.X; x += bounds.Dx() / 10 { // Sample 10 columns
			log.Printf("Pixel (%d, %d): %v", x, y, image.At(x, y))
		}
	}
}
