package main

import (
	"image/color"
	"log"
	"strconv"

	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/ecs"
)

type Santa struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.MouseComponent
}

type myScene struct {}

type MyLabel struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// Type uniquely defines your game type
func (*myScene) Type() string { return "GoClickGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
	err := engo.Files.Load("santastand.png")
	if err != nil {
		log.Println(err)
	}

	err2 := engo.Files.Load("GrinchedRegular.ttf")
	if err2 != nil {
		log.Println(err2)
	}
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	common.SetBackground(color.White)
	//Load Font
	fnt := &common.Font{
		URL:  "GrinchedRegular.ttf",
		FG:   color.Black,
		Size: 32,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}
	var score = 0
	label1 := MyLabel{BasicEntity: ecs.NewBasic()}
	label1.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: strconv.Itoa(score),
	}
	label1.SetShader(common.HUDShader)

// Retrieve a texture
	texture, err := common.LoadedSprite("santastand.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	santa := Santa{BasicEntity: ecs.NewBasic()}


	// Initialize the components, set scale to 8x
	santa.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{1, 1},
	}
	santa.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{240, 120},
		Width:    texture.Width() * santa.RenderComponent.Scale.X,
		Height:   texture.Height() * santa.RenderComponent.Scale.Y,
	}


	// Add it to appropriate systems
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&santa.BasicEntity, &santa.RenderComponent, &santa.SpaceComponent)
			sys.Add(&label1.BasicEntity, &label1.RenderComponent, &label1.SpaceComponent)
			sys.Add(&santa.BasicEntity, &santa.RenderComponent, &santa.SpaceComponent)
			sys.Add(&santa.BasicEntity, &santa.RenderComponent, &santa.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&santa.BasicEntity, &santa.MouseComponent, &santa.SpaceComponent, &santa.RenderComponent, )
		}

	}

}


func main() {
	opts := engo.RunOptions{
		Title: "GoClicker",
		Width:  640,
		Height: 480,
	}
	engo.Run(opts, &myScene{})
}