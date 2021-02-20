package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func generateSNESDefaultConfig() Controller {
	cfg := Controller{
		Name:          "SNES",
		ImageFilePath: "./ArtAssets/controller.png",
		XSize:         470,
		YSize:         220,
		Buttons: []Button{
			{ImageFilePath: "./ArtAssets/circle.png", Name: "A", XLoc: 374, YLoc: 90, BitOffset: 8},
			{ImageFilePath: "./ArtAssets/circle.png", Name: "B", XLoc: 338, YLoc: 122, BitOffset: 0},
			{ImageFilePath: "./ArtAssets/circle.png", Name: "X", XLoc: 334, YLoc: 58, BitOffset: 9},
			{ImageFilePath: "./ArtAssets/circle.png", Name: "Y", XLoc: 298, YLoc: 88, BitOffset: 1},
			{ImageFilePath: "./ArtAssets/up.png", Name: "Up", XLoc: 98, YLoc: 68, BitOffset: 4},
			{ImageFilePath: "./ArtAssets/down.png", Name: "Down", XLoc: 98, YLoc: 112, BitOffset: 5},
			{ImageFilePath: "./ArtAssets/left.png", Name: "Left", XLoc: 68, YLoc: 96, BitOffset: 6},
			{ImageFilePath: "./ArtAssets/right.png", Name: "Right", XLoc: 112, YLoc: 96, BitOffset: 7},
			{ImageFilePath: "./ArtAssets/L.png", Name: "L", XLoc: 58, YLoc: 13, BitOffset: 10},
			{ImageFilePath: "./ArtAssets/R.png", Name: "R", XLoc: 308, YLoc: 13, BitOffset: 11},
			{ImageFilePath: "./ArtAssets/startselect.png", Name: "Select", XLoc: 180, YLoc: 104, BitOffset: 2},
			{ImageFilePath: "./ArtAssets/startselect.png", Name: "Start", XLoc: 226, YLoc: 104, BitOffset: 3},
		},
	}
	return cfg
}

func loadConfiguration() Controller {
	var cfg Controller
	dat, err := ioutil.ReadFile("./controller.cfg")
	if err != nil {
		fmt.Println("Error opening config file - generating new config.")
		cfg = generateSNESDefaultConfig()
		return cfg
	}
	cfg = unmarshallJSON(dat)
	return cfg
}

func unmarshallJSON(b []byte) Controller {
	var cfg Controller
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		fmt.Println(err)
	}
	return cfg
}

func marshallJSON(cfg Controller) []byte {
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return out
}

func saveToFile(b []byte) {
	err := ioutil.WriteFile("controller.cfg", b, 0644)
	if err != nil {
		panic(err)
	}
}

func (c Controller) SaveConfiguration() {
	out := marshallJSON(c)
	saveToFile(out)
}
