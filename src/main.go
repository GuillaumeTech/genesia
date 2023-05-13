package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const convolutionSize = 3
const ExtendMaxSize = 9

func dilate(source image.Image) image.Image {
	srcBounds := source.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()
	dst := image.NewRGBA(source.Bounds())
	for i := 0; i < int(srcW); i++ {
		for j := 0; j < int(srcH); j++ {
			neighbors := findNeighbors(i, j, source)
			if len(neighbors) > 0 {
				dst.Set(i, j, color.Black)
			}
		}
	}

	return dst
}

func erode(source image.Image) image.Image {
	srcBounds := source.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()
	dst := image.NewRGBA(source.Bounds())
	for i := 0; i < int(srcW); i++ {
		for j := 0; j < int(srcH); j++ {
			neighbors := findNeighbors(i, j, source)
			if len(neighbors) == 8 {
				dst.Set(i, j, color.Black)
			}
		}
	}

	return dst
}

func unitVector(x int, y int) [2]float64 {
	xFloat, yFloat := float64(x), float64(y)
	length := math.Sqrt(xFloat*xFloat + yFloat*yFloat)
	return [2]float64{xFloat / length, yFloat / length}
}

func coordsToString(x int, y int) string {
	return fmt.Sprintf("%v,%v", x, y)
}

func stringCoordsToInts(coords string) (x int, y int) {
	splitedCoords := strings.Split(coords, ",")
	parsedX, _ := strconv.Atoi(splitedCoords[0])
	parsedY, _ := strconv.Atoi(splitedCoords[1])
	return parsedX, parsedY
}

func isEmpty(r uint32, g uint32, b uint32, a uint32) bool {
	var max uint32 = 65535
	if a == 0 {
		return true
	}
	if r != max || g != max || b != max {
		return false
	}
	return true
}

func findNeighbors(x int, y int, image image.Image) [][]int {
	neighbors := [][]int{}
	Xindexes := []int{x - 1, x, x + 1}
	Yindexes := []int{y - 1, y, y + 1}
	for _, Xindex := range Xindexes {
		for _, Yindex := range Yindexes {
			r, g, b, a := image.At(Xindex, Yindex).RGBA()
			if !isEmpty(r, g, b, a) && (Xindex != x || Yindex != y) {
				neighbor := []int{Xindex, Yindex}
				neighbors = append(neighbors, neighbor)
			}
		}
	}
	return neighbors
}

func loadImage(path string) image.Image {
	// Read image from file that already exists
	existingImageFile, err := os.Open(path)
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()

	// Alternatively, since we know it is a png already
	// we can call png.Decode() directly
	loadedImage, err := png.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}
	return loadedImage
}

func writeImage(imageToWrite image.Image, path string) {
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create(path)
	if err != nil {
		// Handle error
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, imageToWrite)

	// Don't forget to close files
	outputFile.Close()
}

func main() {

	var baseImage = loadImage("test.png")
	allNeighbors := make(map[string][][]int)
	bounds := baseImage.Bounds()
	imageWidth := bounds.Dx()
	imageHeight := bounds.Dy()
	result := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	resultLength := len(result.Pix)

	for i := 0; i < int(imageWidth); i++ {
		for j := 0; j < int(imageHeight); j++ {

			r, g, b, a := baseImage.At(i, j).RGBA()
			if !isEmpty(r, g, b, a) {
				neighbors := findNeighbors(i, j, baseImage)
				key := coordsToString(i, j)
				allNeighbors[key] = neighbors
			}
		}
	}

	for currentPx, currentNeighbors := range allNeighbors {
		vec := []int{}
		x, y := stringCoordsToInts(currentPx)
		numberOfNeighbors := len(currentNeighbors)
		if numberOfNeighbors == 0 {
			vec = []int{0, 1}
		}
		if numberOfNeighbors > 0 && numberOfNeighbors < 2 {
			vec = []int{x - currentNeighbors[0][0], y - currentNeighbors[0][1]}
		} else {
			vec = []int{currentNeighbors[0][0] - currentNeighbors[numberOfNeighbors-1][0], currentNeighbors[0][1] - currentNeighbors[numberOfNeighbors-1][1]}
		}
		normal_1 := unitVector(vec[1], -vec[0])

		len_1 := rand.Intn(ExtendMaxSize-5) + 5
		for i := 0; i < len_1; i++ {
			XtoSet := int(math.Floor(float64(i)*normal_1[0] + float64(x)))
			YtoSet := int(math.Floor(float64(i)*normal_1[1] + float64(y)))

			ToSetOnSlice := 4 * (YtoSet*imageWidth + XtoSet)
			if ToSetOnSlice < 0 || ToSetOnSlice >= resultLength {
				continue
			}
			result.Pix[ToSetOnSlice] = 0     // 1st pixel red
			result.Pix[ToSetOnSlice+1] = 0   // 1st pixel green
			result.Pix[ToSetOnSlice+2] = 0   // 1st pixel blue
			result.Pix[ToSetOnSlice+3] = 255 // 1st pixel alpha
		}

		normal_2 := unitVector(-vec[1], vec[0])

		len_2 := rand.Intn(ExtendMaxSize-5) + 5
		for i := 0; i < len_2; i++ {
			XtoSet := int(math.Floor(float64(i)*normal_2[0] + float64(x)))
			YtoSet := int(math.Floor(float64(i)*normal_2[1] + float64(y)))

			ToSetOnSlice := 4 * (YtoSet*imageWidth + XtoSet)
			if ToSetOnSlice < 0 || ToSetOnSlice >= resultLength {
				continue
			}
			result.Pix[ToSetOnSlice] = 0     // 1st pixel red
			result.Pix[ToSetOnSlice+1] = 0   // 1st pixel green
			result.Pix[ToSetOnSlice+2] = 0   // 1st pixel blue
			result.Pix[ToSetOnSlice+3] = 255 // 1st pixel alpha
		}

	}

	dilated := dilate(result)
	eroded := erode(dilated)

	writeImage(eroded, "test-post.png")

}
