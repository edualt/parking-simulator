package main

import (
	"simulador/scenes"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/scene"
)

func main() {
	oak.AddScene("tds", scene.Scene{Start: scenes.MainScene})

	oak.Init("tds", func(c oak.Config) (oak.Config, error) {

		c.BatchLoad = true
		c.Assets.ImagePath = "assets/images"
		c.Assets.AudioPath = "assets/audio"
		return c, nil
	})
}
