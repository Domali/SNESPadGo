package controller

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	_ "image/png"
	"os"
)

type Button struct {
	BitOffset     int           `json:"bit_offset"`
	XLoc          float64       `json:"x_loc"`
	YLoc          float64       `json:"y_loc"`
	DrawLocation  pixel.Vec     `json:"-"`
	Sprite        *pixel.Sprite `json:"-"`
	ImageFilePath string        `json:"image_file_path"`
	Name          string        `json:"name"`
	Pressed       bool          `json:"-"`
}

type Controller struct {
	Buttons       []Button      `json:"buttons"`
	Sprite        *pixel.Sprite `json:"-"`
	ImageFilePath string        `json:"base_image"`
	Name          string        `json:"name"`
	XSize         float64       `json:"x_size"`
	YSize         float64       `json:"y_size"`
}

func (c *Controller) Update(buttonStatus string) {
	// Because we're reading a string all of these bits are actually backward....
	// So we need to subtract the offset from 15(16 total bits minus 1 for starting
	// at 0). Since we're getting a byte back we check to see if its equal to 48
	// which is the string "0". If its 48 the button is pressed otherwise it will
	// be 49(a string of "1") and not be pressed.
	for i, v := range c.Buttons {
		if buttonStatus[15-v.BitOffset] == 48 {
			c.Buttons[i].Pressed = true
		} else {
			c.Buttons[i].Pressed = false
		}
	}
}

func (c *Controller) DrawController(win *pixelgl.Window) {
	c.Sprite.Draw(win, pixel.IM.Moved(c.Sprite.Picture().Bounds().Center()))
	for _, v := range c.Buttons {
		if v.Pressed {
			v.Sprite.Draw(win, pixel.IM.Moved(v.DrawLocation))
		}
	}
}

func NewController() Controller {
	controller := loadConfiguration()
	err := controller.initializeImages(controller)
	if err != nil {
		fmt.Println("Failed to initialized controller!")
		os.Exit(1)
	}
	controller.generateButtonVectors()
	return controller
}

func (c *Controller) generateButtonVectors() {
	for i, v := range c.Buttons {
		// Since sprite locations have there (0,0) location be the center of the image we need to adjust
		// from the center size of the image.
		vs := v.Sprite.Picture().Bounds().Center()

		// Our image Y assume that 0 starts at the top of the image. Unfortunately the library has
		// 0 starting at the bottom. Here we do some math so the add function correctly subtracts
		// the value in va. We also subtract our Y location from the Y size of the canvas so we can
		// start at the top value and subtract from there.
		vs.Y = vs.Y * -1
		va := pixel.V(v.XLoc, c.YSize-v.YLoc)
		c.Buttons[i].DrawLocation = vs.Add(va)
	}
}

func loadSprite(path string) (*pixel.Sprite, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	picture := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(picture, picture.Bounds())
	return sprite, nil
}

func (c *Controller) initializeImages(controller Controller) error {
	var err error
	c.Sprite, err = loadSprite(controller.ImageFilePath)
	if err != nil {
		fmt.Println("Failed to load controller image.")
		return err
	}
	for i, v := range c.Buttons {
		c.Buttons[i].Sprite, err = loadSprite(v.ImageFilePath)
		if err != nil {
			fmt.Println("Failed to load image for %v button.", v.Name)
			return err
		}
	}
	return nil
}
