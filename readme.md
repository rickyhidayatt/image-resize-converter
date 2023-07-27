# Image Conversion with GO
## Requirement:
- Golang: 1.20 or later
- ImageMagick: 7.1.1 or later

## Task:
### Quality vs Execution Time
Find balance for acceptable quality with acceptable execution time
Image Type|Conversion To|Quality|Speed|Execution Time|Result
--------------|-----------------|--------|------|---------------|---------
JPG|WebP|80|-|3.259043s|Quality: Acceptable, Execution Time: Fast

### Implement resizeImage
Rules:
- If image width > 1920, resize to 1920, else do nothing
- Keep original image type, if png, then the result should also be in png
- Quality must be 80 or over

### Implement createThumbnail
Rules:
- Keep original image type, if png, then the result should also be in png
- Thumbnail must not exceed 50 kB, please find the desired width for this
- Quality must be "acceptable"
- Execution time must be under 1s

### Concurrency
Execute WebP conversion, AVIF conversion, Thumbnail Generation concurrently with Channel.

Rules:
- If image width > 1920, then resize first before doing all of those above.