package main

import (
	"fmt"
	"time"

	chai "github.com/mhamedGd/chai"
	"github.com/mhamedGd/chai/ecs"
)

var game chai.App

var RenderQuadTree chai.StaticQuadTreeContainer[chai.Pair[chai.Transform, chai.FillRectRenderComponent]]

// var rects []chai.Pair[chai.Transform, chai.FillRectRenderComponent]
var Screen_Dims chai.Vector2f

func main() {
	Screen_Dims = chai.NewVector2f(800, 600)
	// RenderQuadTree = chai.NewStaticQuadTreeContainer[chai.Pair[chai.Transform, chai.FillRectRenderComponent]]()
	// RenderQuadTree.Resize(chai.Rect{Position: chai.NewVector2f(0.0, 0.0), Size: chai.NewVector2f(10000, 10000)})
	// const GRID_WIDTH = 500
	// const GRID_HEIGHT = 500
	// const CELL_SIZE = float32(1)
	// var GRID_OFFSET chai.Vector2f = chai.Vector2f{X: -250.0, Y: -250.0}
	// for x := 0; x < GRID_WIDTH; x++ {
	// 	for y := 0; y < GRID_HEIGHT; y++ {
	// 		rectTransform := chai.Transform{
	// 			// Position: chai.RandVector2f().Scale(2000),
	// 			Position:   chai.NewVector2f(float32(x)*CELL_SIZE+(float32(x)*0.25), float32(y)*CELL_SIZE+(float32(y)*0.25)).Add(GRID_OFFSET),
	// 			Dimensions: chai.NewVector2f(1.0, 1.0).Scale(CELL_SIZE),
	// 			Scale:      1.0,
	// 			Rotation:   0,
	// 		}

	// 		rectComp := chai.FillRectRenderComponent{
	// 			Tint: chai.GetRandomRGBA8(),
	// 		}
	// 		// rects = append(rects, chai.Pair[chai.Transform, chai.FillRectRenderComponent]{First: rectTransform, Second: rectComp})
	// 		RenderQuadTree.Insert(chai.Pair[chai.Transform, chai.FillRectRenderComponent]{First: rectTransform, Second: rectComp}, chai.Rect{Position: rectTransform.Position.Subtract(rectTransform.Dimensions.Scale(0.5)), Size: rectTransform.Dimensions})
	// 	}
	// }

	game = chai.App{
		Width:  800,
		Height: 600,
		Title:  "Test",

		OnStart: func() {

			chai.LogF("STARTED\n")
			default_scene := chai.NewScene()
			default_scene.OnSceneStart = func(thisScene *chai.Scene) {
				thisScene.Background = chai.NewRGBA8(0, 10, 20, 255)

				chai.BindInput("left", chai.KEY_A)
				chai.BindInput("right", chai.KEY_D)
				chai.BindInput("up", chai.KEY_W)
				chai.BindInput("down", chai.KEY_S)
				chai.BindInput("zoomin", chai.KEY_E)
				chai.BindInput("zoomout", chai.KEY_Q)

				thisScene.NewRenderSystem(&chai.SpriteRenderSystem{Sprites: &chai.Sprites, Scale: 0.1})
				thisScene.NewRenderSystem(&chai.ShapesDrawingSystem{Shapes: &chai.Shapes})
				logo_id := thisScene.NewEntityId()
				logo_transform := chai.Transform{
					Position:   chai.Vector2fZero,
					Dimensions: chai.Vector2fOne,
					Scale:      1.0,
					Rotation:   0.0,
				}
				logo_sprite := chai.SpriteComponent{
					Texture: chai.LoadPng("Assets/Chai_Logo_transparent.png", &chai.TextureSettings{Filter: chai.TEXTURE_FILTER_LINEAR}),
					Tint:    chai.WHITE,
				}

				thisScene.AddComponents(logo_id, chai.ToComponent(logo_transform), chai.ToComponent(logo_sprite))

				font_render_system := chai.FontRenderSystem{}
				font_render_system.FontSettings = chai.FontBatchSettings{
					FontSize: 24, DPI: 48, CharDistance: 4, LineHeight: 36, Arabic: false,
				}
				font_render_system.SetFont("Assets/m5x7.ttf")
				font_render_system.SetFontBatchRenderer(&chai.UISprites)
				thisScene.NewRenderSystem(&font_render_system)
				thisScene.NewUpdateSystem(&DebugTransformSystem{})

				text_id := thisScene.NewEntityId()
				text_transform := chai.Transform{
					Position:   chai.NewVector2f(25.0, 570.0),
					Dimensions: chai.Vector2fOne,
					Scale:      1.0,
					Rotation:   0.0,
				}
				fontRender := chai.FontRenderComponent{
					Text:   "Hi\nJoe",
					Scale:  3.0,
					Offset: chai.Vector2fZero,
					Tint:   chai.WHITE,
				}

				thisScene.AddComponents(text_id, chai.ToComponent(text_transform), chai.ToComponent(fontRender))
				text_timer := time.NewTimer(100 * time.Millisecond)
				go func() {
					for {
						<-text_timer.C
						chai.GetComponentPtr[chai.FontRenderComponent](thisScene, text_id).Text = fmt.Sprintf("FPS: %.2f\nRects in View: %d", 1.0/chai.GetDeltaTime(), chai.NumOfQuadsInView())
						text_timer.Reset(100 * time.Millisecond)
					}
				}()

				const GRID_WIDTH = 75
				const GRID_HEIGHT = 75
				const CELL_SIZE = float32(1)
				var GRID_OFFSET chai.Vector2f = chai.Vector2f{X: 0.0, Y: 0.0}

				for x := 0; x < GRID_WIDTH; x++ {
					for y := 0; y < GRID_HEIGHT; y++ {
						rectId := thisScene.NewEntityId()

						rectTransform := chai.Transform{
							Position:   chai.NewVector2f(float32(x)*CELL_SIZE+(float32(x)*0.25), float32(y)*CELL_SIZE+(float32(y)*0.25)).Add(GRID_OFFSET),
							Dimensions: chai.NewVector2f(1.0, 1.0).Scale(CELL_SIZE * chai.RandomRangeFloat32(0.2, 1.0)),
							Scale:      1.0,
							Rotation:   0,
						}

						thisScene.AddComponents(rectId, chai.ToComponent(rectTransform))
						quadComp := chai.NewQuadComponent(thisScene, rectId, chai.GetRandomRGBA8())
						thisScene.AddComponents(rectId, chai.ToComponent(quadComp))

						// rects = append(rects, chai.Pair[chai.Transform, chai.FillRectRenderComponent]{First: rectTransform, Second: rectComp})
						// RenderQuadTree.Insert(chai.Pair[chai.Transform, chai.FillRectRenderComponent]{First: rectTransform, Second: rectComp}, chai.Rect{Position: rectTransform.Position.Subtract(rectTransform.Dimensions.Scale(0.5)), Size: rectTransform.Dimensions})
					}
				}
				chai.ScaleView(800)

				// for i := 0; i < 20000; i++ {
				// 	rectId := thisScene.NewEntityId()

				// 	rectTransform := chai.Transform{
				// 		Position:   chai.RandVector2f().Scale(4000),
				// 		Dimensions: chai.RandPosVector2f().Scale(150),
				// 		Scale:      1.0,
				// 		Rotation:   0,
				// 	}

				// 	rectComp := chai.FillRectRenderComponent{
				// 		Dimensions: rectTransform.Dimensions,
				// 		Tint:       chai.GetRandomRGBA8(),
				// 	}

				// 	thisScene.AddComponents(rectId, chai.ToComponent(rectTransform), chai.ToComponent(rectComp))
				// }

			}
			default_scene.OnSceneUpdate = func(dt float32, thisScene *chai.Scene) {
				x_axis := chai.GetActionStrength("right") - chai.GetActionStrength("left")
				y_axis := chai.GetActionStrength("up") - chai.GetActionStrength("down")

				chai.ScrollView(chai.NewVector2f(x_axis, y_axis).Scale(0.5))
				chai.IncreaseScaleU(chai.GetActionStrength("zoomin") - chai.GetActionStrength("zoomout"))

				if chai.IsMousePressed(0) {
					chai.Shapes.DrawFillRect(chai.GetMouseWorldPosition(), chai.NewVector2f(10, 10), chai.NewRGBA8(255, 255, 255, 100))
					q := chai.GetQuadsInRect(chai.Rect{Position: chai.GetMouseWorldPosition().Subtract(chai.NewVector2f(5.0, 5.0)), Size: chai.NewVector2f(10.0, 10.0)})
					for _, v := range q.Data {
						chai.RenderQuadTreeContainer.Remove(v)
					}
				}
			}

			chai.ChangeScene(&default_scene)
		},
		OnUpdate: func(dt float32) {
			// chai.LogF("FPS: %v", 1.0/dt)
		},
		OnDraw: func(dt float32) {
			// for _, v := range RenderQuadTree.Search(chai.Rect{Position: Screen_Dims.Scale(-0.5).Scale(1 / chai.Cam.GetScale()).Add(chai.Cam.GetPosition()), Size: Screen_Dims.Scale(1 / chai.Cam.GetScale())}) {
			// 	chai.Shapes.DrawFillRect(v.First.Position, v.First.Dimensions, v.Second.Tint)
			// 	rects_count++
			// }
			chai.UIShapes.DrawFillRect(chai.NewVector2f(200.0, 570), chai.NewVector2f(400, 100), chai.BLACK)
		},
		OnEvent: func(ae *chai.AppEvent) {
		},
	}

	chai.Run(&game)

}

type DebugTransformSystem struct {
	chai.EcsSystem
}

func (dts *DebugTransformSystem) Update(dt float32) {
	chai.Iterate2[chai.Transform, chai.FillRectRenderComponent](func(i ecs.Id, t *chai.Transform, frc *chai.FillRectRenderComponent) {
		if i == 4 {
			chai.LogF("Double: %p", t)
		}
		t.Position = t.Position.Add(chai.Vector2fDown.Scale(0.05))
		// t.Rotation += dt * 150.0
		// t.Rotation += dt * 150.0
		// t.Rotation += dt * 150.0
		// t.Rotation += dt * 150.0
		t.Scale += dt
	})

	chai.Iterate1[chai.Transform](func(i ecs.Id, t *chai.Transform) {
		if i == 4 {
			chai.LogF("Single: %p", t)
		}
	})
}
