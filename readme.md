# Image Resizer with WebP and AVIF Conversion

This Go program allows you to resize images and convert them to WebP and AVIF formats. It leverages the ImageMagick library for image manipulation.

## Requirements

- Golang: 1.20 or later
- ImageMagick: 7.1.1 or later

## Usage

1. Install Golang version 1.20 or later.
2. Install ImageMagick version 7.1.1 or later.
3. Clone this repository.

### Image Resizing and Conversion

The `main` function in the code performs image resizing and conversion. It accepts an input image file path and an output directory path. The image will be resized if its width exceeds the maximum allowed width of 1860 pixels. The resized image will be saved in the temporary directory before converting it to WebP and AVIF formats.

### Quality Settings

The code allows you to set the quality of the output images. The default quality is 80, and it is recommended to use values between 60 and 80 for a good balance between image size and visual appearance.

### Resizing Images

If the input image width exceeds the maximum allowed width, it will be resized proportionally to fit within the limit. The resize percentage is calculated based on the difference between the original width and the maximum allowed width.

### Creating Thumbnails

The code can also create thumbnails for the input image. Thumbnails are saved in the specified output directory. The thumbnail size and quality can be adjusted to suit your needs.

### Conversion to WebP and AVIF

The program converts the resized image and the original image (if not resized) to WebP and AVIF formats concurrently using goroutines. The converted images are saved in the output directory.

## How to Use

1. Update the `inputFilePath` and `outputDirPath` variables in the `main` function with your desired input image path and output directory path, respectively.
2. Adjust the `maxAllowedWidth` variable to set the maximum allowed width for images.
3. (Optional) Modify the default quality value to achieve the desired balance between image size and quality. Adjust the `defaultQuality` variable accordingly.
4. Run the program using the `go run` command:

```bash
go mod tidy
go run main.go
