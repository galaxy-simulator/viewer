package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ajstarks/svgo"
	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

const (
	width = 1920 * 8
	height = 1920 * 8
)

var (
	treeArray []*structs.Node
)

func drawTree(w http.ResponseWriter, r *http.Request) {
	log.Println("The drawtree handler was accessed")
	w.Header().Set("Content-Type", "image/svg+xml")

	// get the tree index
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	// define the svg
	s := svg.New(w)
	s.Start(width, height)
	s.Rect(0, 0, width, height, s.RGB(0, 0, 0))
	s.Gtransform(fmt.Sprintf("translate(%d,%d)", width/2, height/2))

	getGalaxy(treeindex)
	listOfStars := treeArray[treeindex].GetAllStars()

	// draw the galaxy
	drawStars(s, listOfStars)
	drawBoxes(s, treeindex)

	s.Gend()
	s.End()
}

func drawStars(s *svg.SVG, listOfStars []structs.Star2D) {
	log.Println("[   ] Drawing the stars")
	for _, star := range listOfStars {
		x := int(star.C.X / 2000)
		y := int(star.C.Y / 2000)
		s.Circle(x, y, 1, s.RGB(255, 255, 255))
	}
	log.Println("[   ] Done drawing the stars")
}

func drawBoxes(s *svg.SVG, treeindex int64) {
	log.Println("[   ] Drawing the Boxes")
	drawBox(s, treeArray[treeindex])
	log.Println("[   ] Done drawing the Boxes")
}

func drawBox(s *svg.SVG, node *structs.Node) {
	if node.Boundry != (structs.BoundingBox{}) {
		x := int(node.Boundry.Center.X / 2000)
		y := int(node.Boundry.Center.Y / 2000)
		w := int(node.Boundry.Width / 2000)
		s.CenterRect(x, y, w, w, "fill:none;stroke:white")
	}

	for i := 0; i < len(node.Subtrees); i++ {
		if node.Subtrees[i] != nil {
			drawBox(s, node.Subtrees[i])
		}
	}
}

func getGalaxy(index int64) {
	log.Println("[   ] Getting the Galaxy")
	// make a http-post request to the databse requesting the tree
	requesturl := fmt.Sprintf("http://db.nbg1.emile.space/dumptree/%d", index)
	resp, err := http.Get(requesturl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		panic(readerr)
	}

	tree := &structs.Node{}
	jsonUnmarshalErr := json.Unmarshal(body, tree)
	if jsonUnmarshalErr != nil {
		panic(jsonUnmarshalErr)
	}

	// if the treeArray is not long enough, fill it
	for int(index) > len(treeArray) {
		emptyNode := structs.NewNode(structs.NewBoundingBox(structs.NewVec2(0, 0), 10))
		treeArray = append(treeArray, emptyNode)
	}

	treeArray = append(treeArray, tree)
	log.Println("[   ] Done Getting the galaxy")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/drawtree/{treeindex}", drawTree).Methods("GET")

	log.Fatal(http.ListenAndServe(":80", router))
}
