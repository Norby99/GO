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
	"reflect"
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

func objInRegion(img gocv.Mat, x1, y1, x2, y2 int, bgr [3]uint8, reg []gocv.Mat) bool{
	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			var pixel = [3]uint8{reg[0].GetUCharAt(j, i), reg[1].GetUCharAt(j, i), reg[2].GetUCharAt(j, i)}
			if !(reflect.DeepEqual(pixel, bgr)) {
				return true
			}
		}
	}
	return false
}

func midCactusDetector(img gocv.Mat, x1, y1, x2, y2 int, bgr [3]uint8, reg []gocv.Mat) bool{
	var rect = image.Rectangle{image.Point{x1, y1}, image.Point{x2, y2}}

	gocv.Rectangle(&img, rect, red, 2)

	return objInRegion(img, x1, y1, x2, y2, bgr, reg)
}

func midBirdDetector(img gocv.Mat, x1, y1, x2, y2 int, bgr [3]uint8, reg []gocv.Mat) bool{
	var rect = image.Rectangle{image.Point{x1, y1}, image.Point{x2, y2}}

	gocv.Rectangle(&img, rect, blue, 2)

	return objInRegion(img, x1, y1, x2, y2, bgr, reg)
}

var red = color.RGBA{255, 0, 0, 10}
var blue = color.RGBA{0, 0, 255, 10}
var green = color.RGBA{0, 255, 0, 10}
var yellow = color.RGBA{0, 255, 255, 10}

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
		var split = gocv.Split(img)
		var backGroundColor = [3]uint8{split[0].GetUCharAt(10, 100), split[1].GetUCharAt(10, 100), split[2].GetUCharAt(10, 100)}

		if midCactusDetector(img, 100, 335, 200, 350, backGroundColor, split){
			kup.Press()
			time.Sleep(time.Millisecond)
			kup.Release()
		}else if midBirdDetector(img, 100, 300, 200, 330, backGroundColor, split){
			kdown.Press()
			time.Sleep(400 * time.Millisecond)
			kdown.Release()
		}

		window.IMShow(img)

		if window.WaitKey(1) == 27  || window.GetWindowProperty(gocv.WindowPropertyAspectRatio) < 0{	// Close the window if ESC button was pressed or the x button is pressed
			os.Exit(3)
			break
		}
	}

}