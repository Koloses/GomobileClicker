package main

import (
	"image/color"
	"log"
	"strconv"

	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/ecs"
)
var fnt *common.Font
type Santa struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.MouseComponent
	ClickComponent
}

type myScene struct {}

type MyLabel struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	ClickComponent
}

type ClickComponent struct {
	label string
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
	world.AddSystem(&ClickSystem{})

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
	label1 := MyLabel{BasicEntity: ecs.NewBasic()}
	label1.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "LMB: (0, 0)",
	}
	label1.SetShader(common.HUDShader)
	label1.ClickComponent.label = "left click"


// Retrieve a texture
	texture, err := common.LoadedSprite("santastand.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	santa := Santa{BasicEntity: ecs.NewBasic()}


	// Initialize the components
	santa.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{1, 1},
	}
	santa.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{240, 120},
		Width:    texture.Width() * santa.RenderComponent.Scale.X,
		Height:   texture.Height() * santa.RenderComponent.Scale.Y,
	}
		santa.ClickComponent.label = "left click"


	// Add it to appropriate systems
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&santa.BasicEntity, &santa.RenderComponent, &santa.SpaceComponent)
			sys.Add(&label1.BasicEntity, &label1.RenderComponent, &label1.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&santa.BasicEntity, &santa.MouseComponent, &santa.SpaceComponent, &santa.RenderComponent, )
		case *ClickSystem:
			sys.Add(&label1.BasicEntity,  &label1.RenderComponent, &label1.SpaceComponent, &label1.ClickComponent)
		}

	}

}


type mouseState uint

const (
	up mouseState = iota
	down
	justPressed
)

type clickEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
	*ClickComponent
}

type ClickSystem struct {
	entities []clickEntity

	left, right mouseState
}

func (c *ClickSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, space *common.SpaceComponent, click *ClickComponent) {
	c.entities = append(c.entities, clickEntity{basic, render, space, click})
}

func (c *ClickSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ClickSystem) Update(dt float32) {
	//setup mouse state
	if c.left == justPressed {
		c.left = down
	}
	if c.right == justPressed {
		c.right = down
	}
	if engo.Input.Mouse.Action == engo.Press {
		if engo.Input.Mouse.Button == engo.MouseButtonLeft {
			if c.left == up {
				c.left = justPressed
			}
		} else if engo.Input.Mouse.Button == engo.MouseButtonRight {
			if c.right == up {
				c.right = justPressed
			}
		}
	} else if engo.Input.Mouse.Action == engo.Release {
		if engo.Input.Mouse.Button == engo.MouseButtonLeft {
			c.left = up
		} else if engo.Input.Mouse.Button == engo.MouseButtonRight {
			c.right = up
		}
	}

	//loop through entities
	for _, e := range c.entities {
		switch e.ClickComponent.label {
		case "left click":
			if c.left == justPressed {
				txt := "test"
				e.RenderComponent.Drawable.Close()
				e.RenderComponent.Drawable = common.Text{
					Font: fnt,
					Text: txt,
				}
			}
		case "right click":
			if c.right == justPressed {
				txt := "RMB: (" + strconv.FormatFloat(float64(engo.Input.Mouse.X), 'f', 1, 32) + ", " + strconv.FormatFloat(float64(engo.Input.Mouse.Y), 'f', 1, 32) + ")"
				e.RenderComponent.Drawable.Close()
				e.RenderComponent.Drawable = common.Text{
					Font: fnt,
					Text: txt,
				}
			}
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