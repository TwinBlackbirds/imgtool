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

func (i *TrackedImage) resize(modifier uint) error {
	// resize by [modifier]x
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

func (i *TrackedImage) flip(vertical bool) error {
	// flips the image either vertically or horizontally
	// if vertical is false, it flips it horiz. (obviously)

	if vertical { // for now
		// TODO implement
		return errors.New("vertical flipping has not been implemented yet")
	}

	// verify we have an image loaded in memory
	if i.data == nil {
		return errors.New("no image data, you must load the image first")
	}

	// setup our new image
	rect := i.data.Bounds()
	rgba := image.NewRGBA(rect)

	// -- horizontal flip algorithm
	// first, get the pixel grid
	pixels := i.getPixels()
	//fmt.Println(pixels)
	// and then loop through (half of them) (?), and swap them with their inverse pixel
	bounds := i.data.Bounds().Max
	for x, row := range pixels {
		for y, pixel := range row {
			// find the inverse pixel (one to swap with)
			inverseCoordinate := bounds.Y - y - 1
			inverse := pixels[x][inverseCoordinate]
			// TODO fix image being mirrored on both axes (should just be one)
			// probably due to x being the same down here
			// (or inverse having the same x value when calculated)
			rgba.Set(x, y, inverse)               // set the original pixel to the inverse
			rgba.Set(x, inverseCoordinate, pixel) // set the inverse pixel to the original
		}
	}
	i.data = rgba // set our image to the new, flipped image
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
	err = img.flip(false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Flipped image successfully")
	// save
	filepath := "ss_1.jpg"
	img.format = "jpeg"
	err = img.save(&filepath) // uses existing filepath if 'nil'
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved image:", filepath)

	fmt.Println("All commands executed")
}
