package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	_ "image/jpeg"
	"image/png"

	"github.com/jmiguelff/generative-art/sketch"
)

var (
	sourceImgName   = "source2.jpg"
	totalCycleCount = 5000
)

func main() {

	//img, err := loadImage(sourceImgName)
	img, err := loadRandomUnsplashImage(2000, 2000)
	if err != nil {
		log.Panicln(err)
	}
	saveOutput(img, "in2.png")

	destWidth := 2000
	s := sketch.NewSketch(img, sketch.UserParams{
		StrokeRatio:              0.75,
		DestWidth:                destWidth,
		DestHeight:               2000,
		InitialAlpha:             0.1,
		StrokeReduction:          0.002,
		AlphaIncrease:            0.06,
		StrokeInversionThreshold: 0.05,
		StrokeJitter:             int(0.1 * float64(destWidth)),
		MinEdgeCount:             3,
		MaxEdgeCount:             8,
	})

	rand.Seed(time.Now().Unix())

	// the main loop
	for i := 0; i < totalCycleCount; i++ {
		s.Update()
	}

	saveOutput(s.Output(), "out2.png")
}

func cpuProf(fn func()) {
	f, err := os.Create("cpuprof.out")
	if err != nil {
		fmt.Println("Error in creating file for writing cpu profile: ", err)
		return
	}
	defer f.Close()

	if e := pprof.StartCPUProfile(f); e != nil {
		fmt.Println("Error in starting CPU profile: ", e)
		return
	}

	fn()
	defer pprof.StopCPUProfile()
}

func loadRandomUnsplashImage(width, height int) (image.Image, error) {
	url := fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	return img, err
}

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("source image could not be loaded: %w", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("source image format could not be decoded: %w", err)
	}

	return img, nil
}

func saveOutput(img image.Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
