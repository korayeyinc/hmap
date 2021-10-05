/*
 * A simple command line tool for heightmap generation.
 *
 */

package main

import (
	"errors"
	"flag"
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/convolution"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/histogram"
	"github.com/anthonynsimon/bild/segment"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/fs"
	"log"
	"os"
	"strings"
)

// image definition
type imgfx struct {
	src    image.Image
	dst    *image.Gray
	format string
	width  int
	height int
}

// cmdline options definition
type options struct {
	input    *string
	in       *string
	output   *string
	out      *string
	blur     *float64
	emboss   *string
	gauss    *float64
	contrast *float64
	invert   *string
	mono     *uint
	blend    *float64
	hist     *string
}

var err error

func main() {
	log.SetFlags(0)
	log.SetPrefix("[info] ")
	opt := new(options)
	opt.input = flag.String("input", "", "Input image -- supported image formats include: [BMP, GIF, JPG, PNG, TIFF, WEBP, VP8]")
	opt.in = flag.String("in", "", "Input image -- supported image formats include: [BMP, GIF, JPG, PNG, TIFF, WEBP, VP8]")
	opt.output = flag.String("output", "output.png", "Output image -- supported image formats include: [BMP, GIF, JPG, PNG, TIFF]")
	opt.out = flag.String("out", "output.png", "Output image -- supported image formats include: [BMP, GIF, JPG, PNG, TIFF]")
	opt.blend = flag.Float64("blend", 0.5, "Blend opacity -- percent must be of range 0.0 to 1.0")
	opt.blur = flag.Float64("blur", 0, "Box blur -- radius must be larger than 0.")
	opt.emboss = flag.String("emboss", "low", "Emboss level -- must be [high/mid/low]")
	opt.gauss = flag.Float64("gauss", 0, "Gaussian blur -- radius must be larger than 0.")
	opt.contrast = flag.Float64("contrast", 0, "Contrast level -- must be of the range -100 to 100.")
	opt.invert = flag.String("invert", "off", "Invert on/off -- reverses the colors of the grayscale image.")
	opt.mono = flag.Uint("mono", 0, "Monochrome level -- must be of the range 0 to 255.")
	opt.hist = flag.String("hist", "hist.png", "Histogram output -- produce histogram output to analyze the frequency distribution of colors in the image")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 && args[0] == "help" {
		doc()
		flag.Usage()
		os.Exit(1)
	}

	var input, output string

	if len(*opt.input) > 0 {
		input = *opt.input
	} else if len(*opt.in) > 0 {
		input = *opt.in
	}

	if *opt.output != "output.png" {
		output = *opt.output
	} else if *opt.out != "output.png" {
		output = *opt.out
	} else {
		output = *opt.output
	}

	img := new(imgfx)
	img.load(input)
	img.proc(opt)
	img.save(output)

	genHist(*opt.hist, img.dst)
	log.Println("Done!")
}

// Loads and decodes the input image.
func (img *imgfx) load(filename string) {
	input := fopen(filename)
	defer input.Close()

	log.Printf("Loading input image: %s \n", filename)

	img.src, img.format, err = image.Decode(input)
	check(err)

	size := img.src.Bounds()
	img.width = size.Max.X
	img.height = size.Max.Y

	log.Printf("Image format: %s \n", img.format)
	log.Printf("Image size: %dx%d \n", img.width, img.height)
}

// Applies image filters and processes image.
func (img *imgfx) proc(opt *options) {
	log.Println("Processing image...")

	if *opt.emboss == "low" {
		img.emboss("low")
	} else if *opt.emboss == "mid" {
		img.emboss("mid")
	} else if *opt.emboss == "high" {
		img.emboss("high")
	}

	if *opt.contrast > 0 {
		img.src = adjust.Contrast(img.src, *opt.contrast)
	}

	if *opt.invert == "on" {
		img.src = effect.Invert(img.src)
	}

	if *opt.blur > 0 {
		img.src = blur.Box(img.src, *opt.blur)
	}

	if *opt.gauss > 0 {
		img.src = blur.Gaussian(img.src, *opt.gauss)
	}

	img.dst = effect.Grayscale(img.src)

	if *opt.mono > 0 {
		mono := segment.Threshold(img.src, uint8(*opt.mono))
		chrome := blend.Opacity(img.dst, mono, *opt.blend)
		img.dst = effect.Grayscale(chrome)
	}
}

// Applies custom convolution kernels to image for depth depiction.
func (img *imgfx) emboss(level string) {
	var kernel convolution.Kernel
	if level == "low" {
		kernel = convolution.Kernel{
			Matrix: []float64{
				-1, -1, 0,
				-1, 0, 1,
				0, 1, 1,
			},
			Width:  3,
			Height: 3,
		}
	} else if level == "mid" {
		kernel = convolution.Kernel{
			Matrix: []float64{
				-1, -1, -1, -1, 0,
				-1, -1, -1, 0, 1,
				-1, -1, 0, 1, 1,
				-1, 0, 1, 1, 1,
				0, 1, 1, 1, 1,
			},
			Width:  5,
			Height: 5,
		}
	} else if level == "high" {
		kernel = convolution.Kernel{
			Matrix: []float64{
				-1, -1, -1, -1, -1, -1, 0,
				-1, -1, -1, -1, -1, 0, 1,
				-1, -1, -1, -1, 0, 1, 1,
				-1, -1, -1, 0, 1, 1, 1,
				-1, -1, 0, 1, 1, 1, 1,
				-1, 0, 1, 1, 1, 1, 1,
				0, 1, 1, 1, 1, 1, 1,
			},
			Width:  7,
			Height: 7,
		}
	}

	img.src = convolution.Convolve(img.src, &kernel, &convolution.Options{Bias: 128, Wrap: false, KeepAlpha: true})
}

// Encodes and saves the output image.
func (img *imgfx) save(filename string) {
	output := fcreate(filename)
	defer output.Close()

	log.Printf("Saving output image as: %s \n", filename)

	if strings.Contains(filename, ".") {
		img.format = strings.Split(filename, ".")[1]
	} else {
		img.format = "png"
	}

	// encode the image using the original format
	switch img.format {
	case "bmp":
		bmp.Encode(output, img.dst)
	case "jpg":
	case "jpeg":
		opt := new(jpeg.Options)
		opt.Quality = 100
		jpeg.Encode(output, img.dst, opt)
	case "png":
		png.Encode(output, img.dst)
	case "gif":
		gif.Encode(output, img.dst, nil)
	case "tiff":
		tiff.Encode(output, img.dst, nil)
	default:
		msg := "Can't encode image file!"
		log.Fatalln(msg)
	}
}

// Generates histogram output from the provided image data.
func genHist(filename string, img *image.Gray) {
	log.Println("Generating histogram for the output image...")
	hist := histogram.NewRGBAHistogram(img)
	result := hist.Image()

	output := fcreate(filename)
	defer output.Close()

	png.Encode(output, result)
}

// Helper function for opening a file.
func fopen(filename string) *os.File {
	file, err := os.Open(filename)
	check(err)
	return file
}

// Helper function for creating a file.
func fcreate(filename string) *os.File {
	file, err := os.Create(filename)
	check(err)
	return file
}

// Helper function for checking errors.
func check(err error) {
	if err != nil {
		var pathError *fs.PathError

		if errors.As(err, &pathError) {
			log.Println("Path error:", pathError.Path)
		} else {
			log.Println(err)
		}
	}
}

// Helper function as part of the embedded documentation.
func doc() {
	log.SetPrefix("")
}
