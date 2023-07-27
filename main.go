package main

import (
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var dirPath string

func main() {
	inputFilePath := "C:\\Users\\User\\Desktop\\imageResizer\\samples\\ke2.jpg"
	outputDirPath := "C:\\Users\\User\\Desktop\\imageResizer\\outputs"

	// cek dulu apakah file nya png atau bukan
	size, err := getImageWidth(inputFilePath)
	if err != nil {
		panic(err)
	}

	getFormat := filepath.Ext(inputFilePath)
	if getFormat == ".png" {

		if *size > 1920 {
			resizedOutputPath := fmt.Sprintf("%s%s%s.png", outputDirPath, string(os.PathSeparator), "temp")

			resizeDir, err := ResizeImage(inputFilePath, resizedOutputPath, 80)
			if err != nil {
				log.Fatal("Error resizing image:", err)
			}
			inputFilePath = *resizeDir
		}

		thumbnailCh := make(chan *string)
		go func() {
			dirThumbnail := "C:\\Users\\User\\Desktop\\imageResizer\\outputs\\thumnail"
			outputThumnail := fmt.Sprintf("%s%s%s.png", dirThumbnail, string(os.PathSeparator), uuid.New())
			thumnailOutputPath, err := CreateThumbnail(inputFilePath, outputThumnail, 280)
			if err != nil {
				log.Println("Error converting to WebP:", err)
				thumbnailCh <- nil
				return
			}
			thumbnailCh <- thumnailOutputPath
		}()
		thumbnailResult := <-thumbnailCh

		if thumbnailResult != nil {
			fmt.Println("WebP Conversion completed successfully:", *thumbnailResult)
		}

		return
	}

	if *size > 1920 {
		resizedOutputPath := fmt.Sprintf("%s%s%s.jpg", outputDirPath, string(os.PathSeparator), uuid.New())
		resizeDir, err := ResizeImage(inputFilePath, resizedOutputPath, 80)
		if err != nil {
			log.Fatal("Error resizing image:", err)
		}
		fmt.Println(*resizeDir)
		inputFilePath = *resizeDir
	}
	concurrentConvert(inputFilePath, outputDirPath)
}

func concurrentConvert(inputFilePath, outputDirPath string) {
	webpCh := make(chan *string)
	avifCh := make(chan *string)
	thumbnailCh := make(chan *string)

	// Create Thumnail
	go func() {
		dirThumbnail := "C:\\Users\\User\\Desktop\\imageResizer\\outputs\\thumnail"
		outputThumnail := fmt.Sprintf("%s%s%s.webp", dirThumbnail, string(os.PathSeparator), uuid.New())
		thumnailOutputPath, err := CreateThumbnail(inputFilePath, outputThumnail, 280)
		if err != nil {
			log.Println("Error converting to WebP:", err)
			thumbnailCh <- nil
			return
		}
		thumbnailCh <- thumnailOutputPath
	}()

	// Convert to WebP
	go func() {
		webpOutputPath, err := convertToWebP(inputFilePath, outputDirPath, 80, false)
		if err != nil {
			log.Println("Error converting to WebP:", err)
			webpCh <- nil
			return
		}
		webpCh <- webpOutputPath
	}()

	// Convert to AVIF
	go func() {
		avifOutputPath, err := convertToAvif(inputFilePath, outputDirPath, 80, 3)
		if err != nil {
			log.Println("Error converting to AVIF:", err)
			avifCh <- nil
			return
		}
		avifCh <- avifOutputPath
	}()

	webpResult := <-webpCh
	avifResult := <-avifCh
	thumbnailResult := <-thumbnailCh

	if webpResult != nil {
		fmt.Println("WebP Conversion completed successfully:", *webpResult)
	}
	if avifResult != nil {
		fmt.Println("AVIF Conversion completed successfully:", *avifResult)
	}
	if thumbnailResult != nil {
		fmt.Println("Create Thumnail completed successfully:", *thumbnailResult)
	}
}

func ResizeImage(filePath string, outDir string, quality int) (*string, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}

	startTime := time.Now()
	_, err := exec.Command("magick", filePath, "-resize", "1920x", "-quality", fmt.Sprintf("%d", quality), outDir).Output()
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime).Seconds()
	fmt.Println(fmt.Sprintf("Resize Image time duration : %fs", duration))
	return &outDir, nil
}

func CreateThumbnail(filePath string, outDir string, thumbSize int) (*string, error) {
	startTime := time.Now()
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}

	_, err := exec.Command("magick", filePath, "-resize", fmt.Sprintf("%d", thumbSize), "-quality", "80", outDir).Output()
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime).Seconds()
	fmt.Println(fmt.Sprintf("Thumnail Image time duration : %fs", duration))
	return &outDir, nil
}

func convertToAvif(filePath string, outDir string, quality int, speed int) (*string, error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil, os.ErrNotExist
	}

	_, err := os.Stat(outDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(outDir, 0755)
	}

	fileName := uuid.New()

	startTime := time.Now()
	outputPath := fmt.Sprintf("%s%s%s.avif", outDir, string(os.PathSeparator), fileName)
	_, err = exec.Command(
		"magick",
		filePath,
		"-quality",
		strconv.FormatInt(int64(quality), 10),
		"-define",
		fmt.Sprintf("heic:speed=%d", speed),
		outputPath,
	).Output()

	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime).Seconds()
	fmt.Println(fmt.Sprintf("Converting Image to AVIF took %fs", duration))
	return &outputPath, nil
}

func convertToWebP(filePath string, outDir string, quality int, isLossy bool) (*string, error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil, os.ErrNotExist
	}

	_, err := os.Stat(outDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(outDir, 0755)
	}

	fileName := uuid.New()
	outputPath := fmt.Sprintf("%s%s%s.webp", outDir, string(os.PathSeparator), fileName)

	startTime := time.Now()
	_, err = exec.Command(
		"magick",
		filePath,
		"-quality",
		strconv.FormatInt(int64(quality), 10),
		"-define",
		fmt.Sprintf("webp:lossless=%s", strconv.FormatBool(isLossy)),
		outputPath,
	).Output()

	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime).Seconds()
	fmt.Println(fmt.Sprintf("Converting Image to WebP took %fs", duration))
	return &outputPath, nil
}

func getImageWidth(filePath string) (*int, error) {

	var size int
	getFormat := filepath.Ext(filePath)
	if getFormat == ".jpg" || getFormat == ".jpeg" {

		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		size = img.Bounds().Dx()
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		size = img.Bounds().Dx()
	}

	return &size, nil
}
