package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/fogleman/gg"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	treeArray []*structs.Node
)

func readfile(filename string) {
	// read the json
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// initialize the rootnode
	var rootnode structs.Node

	// unmarshal the json into the rootnode
	err = json.Unmarshal(file, &rootnode)
	if err != nil {
		panic(err)
	}

	treeArray = append(treeArray, &rootnode)
}

// draw the requested tree
func drawtree(treeindex int64, savepath string) {

	// generate a list of all stars
	var starlist []structs.Star2D
	starlist = treeArray[treeindex].GetAllStars()

	log.Println("[   ] Initializing the Plot")
	dc := initializePlot()
	log.Println("[   ] Done Initializing the Plot")

	log.Println("[   ] Drawing the Starlist")
	drawStarlist(dc, starlist)
	log.Println("[   ] Done Drawing the Starlist")

	log.Println("[   ] Drawing the Boxes")
	drawBoxes(dc, treeindex)
	log.Println("[   ] Done Drawing the Boxes")

	log.Println("[   ] Saving the image")
	saveImage(dc, savepath)
	log.Println("[   ] Done Saving the image")
}

func saveImage(dc *gg.Context, path string) {
	err := dc.SavePNG(path)
	if err != nil {
		panic(err)
	}
}

func drawBox(dc *gg.Context, box structs.BoundingBox) {
	x := (box.Center.X / 5e3 * 2.5) - ((box.Width / 5e3 * 2.5) / 4)
	y := (box.Center.Y / 5e3 * 2.5) - ((box.Width / 5e3 * 2.5) / 4)
	w := box.Width / 5e3

	log.Println("[   ] Drawing the Box")
	dc.DrawRectangle(x, y, w, w)
	log.Println("[   ] 0")
	dc.Stroke()
	log.Println("[   ] Done Drawing the Box")
}

func genBoxes(dc *gg.Context, node structs.Node) {

	// if the BoundingBox is not empty, draw it
	if node.Boundry != (structs.BoundingBox{}) {
		drawBox(dc, node.Boundry)
	}

	for i := 0; i < len(node.Subtrees); i++ {
		if node.Subtrees[i] != nil {
			genBoxes(dc, *node.Subtrees[i])
		}
	}
}

func drawBoxes(dc *gg.Context, treeindex int64) {
	log.Println("[   ] before genBoxes")
	root := treeArray[treeindex]
	genBoxes(dc, *root)
	log.Println("[   ] after genBoxes")
}

func drawStar(dc *gg.Context, star structs.Star2D) {
	// scalingFactor := 50
	defaultStarSize := 2.0

	x := star.C.X / 5e3 * 2.5
	y := star.C.Y / 5e3 * 2.5

	fmt.Printf("(%20.3f, %20.3f)\n", x, y)

	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(x, y, defaultStarSize)
	dc.Fill()
	dc.Stroke()
}

func drawStarlist(dc *gg.Context, starlist []structs.Star2D) {
	for _, star := range starlist {
		drawStar(dc, star)
	}
}

// initializePlot generates a new plot and returns the plot context
func initializePlot() *gg.Context {
	// Define the image size
	const imageWidth = 8192 * 2
	const imageHeight = 8192 * 2

	// Initialize the new context
	dc := gg.NewContext(imageWidth, imageHeight)

	// Set the background black
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	// Invert the Y axis (positive values are on the top and right)
	dc.InvertY()

	// Set the coordinate midpoint to the middle of the image
	dc.Translate(imageWidth/2, imageHeight/2)

	return dc
}

func drawallboxes(amount int64) {
	for i := 0; i < int(amount); i++ {
		index := fmt.Sprintf("%d", i)
		readfile(fmt.Sprintf("%s.json", index))
		drawtree(0, fmt.Sprintf("%s.png", index))
	}
}

func main() {
	amount, parseErr := strconv.ParseInt(os.Args[1], 10, 64)
	if parseErr != nil {
		panic(amount)
	}
	drawallboxes(amount)
}
