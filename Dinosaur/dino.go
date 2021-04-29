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

type points struct {
	points int64
	timeStart int64
	setCounter uint8
}


func speedHandler(arrCm, arrBm *[4]int, p *points) {
	p.points = time.Now().Unix() - p.timeStart
	if p.points > 20 && p.setCounter == 0{
		arrCm[0] += 20
		arrCm[2] += 20
		p.setCounter++
	}else if p.points > 30 && p.setCounter == 1{
		arrCm[0] += 20
		arrCm[2] += 20
		p.setCounter++
	}else if p.points > 40 && p.setCounter == 2{
		arrCm[0] += 30
		arrCm[2] += 30
		p.setCounter++
	}else if p.points > 50 && p.setCounter == 3{
		arrCm[0] += 30
		arrCm[2] += 30
		p.setCounter++
	}else if p.points > 60 && p.setCounter == 4{
		arrCm[0] += 30
		arrCm[2] += 30
		p.setCounter++
	}
	arrBm[0] = arrCm[0] + 40
	arrBm[2] = arrCm[2] - 10
}

func reset(arrCm , arrBm *[4]int, p *points){
	arrCm[0] = 100
	arrCm[1] = 340
	arrCm[2] = 180
	arrCm[3] = 350
	p.timeStart = time.Now().Unix()
	p.points = 0
	p.setCounter = 0
	arrBm[0] = arrCm[0] + 40
	arrBm[2] = arrCm[2] - 10
}

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

func gameOverDetector(img gocv.Mat, split []gocv.Mat) bool{
	var x, y int = 600, 300
	for i := x-50; i < x; i++ {
		for j := y-50; j < y; j++ {
			var pixel = [3]uint8{split[0].GetUCharAt(j, i), split[1].GetUCharAt(j, i), split[2].GetUCharAt(j, i)}
			if (reflect.DeepEqual(pixel, [3]uint8{8, 8, 8})) {
				return true
			}
		}
	}
	return false
}

func objInRegion(x1, y1, x2, y2 int, bgr [3]uint8, reg []gocv.Mat) bool{
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

	return objInRegion(x1, y1, x2, y2, bgr, reg)
}

func midBirdDetector(img gocv.Mat, x1, y1, x2, y2 int, bgr [3]uint8, reg []gocv.Mat) bool{
	var rect = image.Rectangle{image.Point{x1, y1}, image.Point{x2, y2}}

	gocv.Rectangle(&img, rect, blue, 2)

	return objInRegion(x1, y1, x2, y2, bgr, reg)
}

var red = color.RGBA{255, 0, 0, 10}
var blue = color.RGBA{0, 0, 255, 10}
var green = color.RGBA{0, 255, 0, 10}
var yellow = color.RGBA{0, 255, 255, 10}

func main() {

	bounds := image.Rectangle{image.Point{0, 110}, image.Point{1200, 700}}//screenshot.GetDisplayBounds(0)
	img := takeScreeShot(bounds)

	window := gocv.NewWindow("Dino")
	window.ResizeWindow(2560/4, 1440/4)
	var status bool = true
	var dinoRect = image.Rectangle{image.Point{0, 285}, image.Point{85, 380}}
	var midCactusCoords [4]int
	var midBirdCoord [4]int
	midBirdCoord[1] = 300
	midBirdCoord[3] = 310
	var points = points{}

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

	fmt.Printf("Center the dino in the red rectangle and press the \"a\" button, then press on the browser to procede\nIn alternative press ESC, or the X button to exit in any moment\n");

	for {
		img = takeScreeShot(bounds)
		gocv.Rectangle(&img, dinoRect, red, 2)

		window.IMShow(img)

		key := window.WaitKey(1)
		if key == 97{	// Continue if a button was pressed
			break
		}
		if key == 27  || window.GetWindowProperty(gocv.WindowPropertyAspectRatio) < 0{	// Close the window if ESC button was pressed or the x button is pressed
			os.Exit(3)
		}
	}
	for {
		reset(&midCactusCoords, &midBirdCoord, &points)
		for status {

			img = takeScreeShot(bounds)
			var split = gocv.Split(img)
			var backGroundColor = [3]uint8{split[0].GetUCharAt(10, 100), split[1].GetUCharAt(10, 100), split[2].GetUCharAt(10, 100)}

			if midCactusDetector(img, midCactusCoords[0], midCactusCoords[1], midCactusCoords[2], midCactusCoords[3], backGroundColor, split){
				kup.Press()
				time.Sleep(time.Millisecond)
				kup.Release()
			}else if midBirdDetector(img, midBirdCoord[0], midBirdCoord[1], midBirdCoord[2], midBirdCoord[3], backGroundColor, split){
				kdown.Press()
				time.Sleep(400 * time.Millisecond)
				kdown.Release()
				fmt.Printf("Down - ")
			}

			fmt.Printf("points: %d\n", points)
			
			speedHandler(&midCactusCoords, &midBirdCoord, &points)
			status = !gameOverDetector(img, split)

			window.IMShow(img)

			if !status {
				fmt.Printf("Game Over! Press  the \"retry\" button, or press ESC to close...\n")
			}
			if window.WaitKey(1) == 27  || window.GetWindowProperty(gocv.WindowPropertyAspectRatio) < 0{	// Close the window if ESC button was pressed or the x button is pressed
				os.Exit(3)
			}

		}

		img = takeScreeShot(bounds)
		var split = gocv.Split(img)
		status = !gameOverDetector(img, split)
		window.IMShow(img)

		key := window.WaitKey(1)
		if key == 27  || window.GetWindowProperty(gocv.WindowPropertyAspectRatio) < 0{	// Close the window if ESC button was pressed or the x button is pressed
			os.Exit(3)
		}
	}
}