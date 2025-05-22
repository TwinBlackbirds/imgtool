// tool which allows you to manipulate images of various formats
// png, jpeg, bmp
// save, load images, flip images vertically and horizontally

package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

type TrackedImage struct {
	filepath string
	format   string
	data     image.Image
}

func (i *TrackedImage) load(filepath string) error {
	// open file handle
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	// close file handle
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// image operations
	decodedImage, format, err := image.Decode(file)
	if err != nil {
		return err
	}
	// initialise class vars
	i.filepath = filepath
	i.format = format
	i.data = decodedImage
	return nil
}

func (i *TrackedImage) save(filepath *string) error {
	// use existing filepath if one is not provided
	if filepath == nil {
		filepath = &i.filepath
	}

	// open/create file handle
	file, err := os.Create(*filepath)
	if err != nil {
		return err
	}
	// close file handle
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// format handling
	switch i.format {
	case "png":
		err = png.Encode(file, i.data)
		break
	case "jpeg":
		err = jpeg.Encode(file, i.data, nil)
		break
	default:
		return fmt.Errorf("unsupported format: %s", i.format)
	}

	if err != nil {
		return err
	}
	return nil
}

func (i *TrackedImage) resize(modifier float64) error {
	if modifier <= 0.1 || (modifier > 1.00 && modifier < 1.01) { // modifiers too small (0.09x - 1.01x)
		return errors.New("unsupported image resize modifier")
	}
	// resize by [modifier]x
	// whether it be 0.1x (smallest), 0.5x, 2x, 4x
	// TODO implement
	return nil
}

func (i *TrackedImage) getPixels() [][]color.Color {
	if i.data == nil {
		return nil
	}
	var pixels [][]color.Color
	size := i.data.Bounds().Size()
	// scroll through vertically
	for x := 0; x < size.X; x++ {
		var row []color.Color
		for y := 0; y < size.Y; y++ {
			row = append(row, i.data.At(x, y))
		}
		pixels = append(pixels, row)
	}
	return pixels
}

func (i *TrackedImage) flipHorizontally() error {
	// same idea as the flipVertically() algorithm, but instead of traversing via the x coordinate,
	// we use the y coordinate and go col by col instead of row by row.
	// flip the colors in the exact same way though.
	// TODO implement
	return errors.New("horizontal flipping has not been implemented yet")
}
func (i *TrackedImage) flipVertically() error {
	// -- image vertical flip algorithm

	// verify we have an image loaded in memory
	if i.data == nil {
		return errors.New("no image data, you must load the image first")
	}
	// setup our new image
	rect := i.data.Bounds()
	rgba := image.NewRGBA(rect)

	// first, get the pixel grid
	pixels := i.getPixels()

	// loop through pixels and swap them
	bounds := i.data.Bounds().Max
	for x, row := range pixels {
		for y, pixel := range row {
			// we only need to go through half of the pixels,
			// otherwise we are unnecessarily copying over the same places twice
			if y > bounds.Y/2 {
				continue
			}
			// find the inverse pixel (one to swap with)
			inverseCoordinate := bounds.Y - y - 1
			inverse := pixels[x][inverseCoordinate]
			// swapping
			rgba.Set(x, y, inverse)               // set the original pixel to the inverse
			rgba.Set(x, inverseCoordinate, pixel) // set the inverse pixel to the original
		}
	}
	i.data = rgba // set our image to the new, flipped image
	return nil
}

func (i *TrackedImage) mirror() error {
	// flip an image both ways (rotate 180deg)
	err := i.flipVertically()
	if err != nil {
		return err
	}
	err = i.flipHorizontally()
	if err != nil {
		return err
	}
	return nil
}

func main() {

	// TODO command line argument parsing

	// load image
	img := TrackedImage{}
	err := img.load("ss.png")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loaded image:", img.filepath)

	// print details
	fmt.Println("Format:", img.format)
	rec := img.data.Bounds()
	res := fmt.Sprintf("%dx%d", rec.Max.X, rec.Max.Y)
	fmt.Println("Resolution:", res)

	// flip image
	err = img.flipVertically()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Vertically flipped image successfully")
	// save
	filepath := "ss_1.png"
	err = img.save(&filepath) // uses existing filepath if 'nil'
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved image:", filepath)

	fmt.Println("All commands executed")
}
