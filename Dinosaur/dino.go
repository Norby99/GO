/*
This project uses opencv to play the chrome dino game.
This code won't be perfect, and it will be used for learning purposes only.
*/

//* chrome://dino/

package main

import (
	"fmt"
	"image"
	"image/color"
	"github.com/kbinani/screenshot"	// screenshot handler
	"github.com/micmonay/keybd_event" // key presser
	"gocv.io/x/gocv"
	"os"
	"runtime"
	"time"
)

func imgToMat(img *image.RGBA) gocv.Mat{

	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()

	bytes := make([]byte, 0, x*y)
	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8))
			bytes = append(bytes, byte(g>>8))
			bytes = append(bytes, byte(r>>8))
		}
	}

	mat, err := gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
	if err != nil {
		panic(err)
	}
	return mat
}

func takeScreeShot(area image.Rectangle) gocv.Mat{
	img, err := screenshot.CaptureRect(area)
	if err != nil {
		panic(err)
	}
	mat := imgToMat(img)
	return mat
}

func midCactusDetector(img gocv.Mat) bool{
	var red = color.RGBA{255, 0, 0, 10}
	var rect = image.Rectangle{image.Point{100, 335}, image.Point{200, 350}}
	gocv.Rectangle(&img, rect, red, 2)
	return true
}

func main() {

	fmt.Print("Running!\n")
	bounds := image.Rectangle{image.Point{0, 110}, image.Point{1200, 700}}//screenshot.GetDisplayBounds(0)
	img := takeScreeShot(bounds)

	window := gocv.NewWindow("Dino")
	window.ResizeWindow(2560/4, 1440/4)
	var status int16 = 0

	kup, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kdown, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
	kup.SetKeys(keybd_event.VK_UP)
	kdown.SetKeys(keybd_event.VK_DOWN)

	fmt.Printf("Press the \"a\" button to procede\nIn alternative press ESC, or the X button to exit in any moment\n");

	for status == 0 {
		key := window.WaitKey(1)
		if key == 97{	// Close the window if ESC button was pressed or the x button is pressed
			status = 1
		}
		if key == 27  || window.GetWindowProperty(gocv.WindowPropertyAspectRatio) < 0{	// Close the window if ESC button was pressed or the x button is pressed
			os.Exit(3)
			break
		}
	}
	for {
		img = takeScreeShot(bounds)

		var midCactus bool = midCactusDetector(img)

		if midCactus{
			kup.Press()
			time.Sleep(time.Millisecond)
			kup.Release()
		}

		window.IMShow(img)

		if window.WaitKey(1) == 27  || window.GetWindowProperty(gocv.WindowPropertyAspectRatio) < 0{	// Close the window if ESC button was pressed or the x button is pressed
			os.Exit(3)
			break
		}
	}

}