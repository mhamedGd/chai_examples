package main

import (
	"fmt"
	"time"

	chai "github.com/mhamedGd/chai"
)

var game chai.App

var RenderQuadTree chai.StaticQuadTreeContainer[chai.Pair[chai.VisualTransform, chai.FillRectRenderComponent]]

var Screen_Dims chai.Vector2f

var level chai.Tilemap

var targets chai.List[chai.Vector2f]
var astar_nodes chai.List[Node_AStar]

var fontSettings chai.FontBatchSettings
var font_atlas chai.FontBatchAtlas

func SceneSetup(thisScene *chai.Scene) {
	thisScene.Background = chai.NewRGBA8(23, 28, 57, 255)

	chai.BindInput("left", chai.KEY_A)
	chai.BindInput("right", chai.KEY_D)
	chai.BindInput("up", chai.KEY_W)
	chai.BindInput("down", chai.KEY_S)
	chai.BindInput("zoomin", chai.KEY_E)
	chai.BindInput("zoomout", chai.KEY_Q)

	thisScene.NewRenderSystem(chai.SpriteRenderSystem)
	thisScene.NewRenderSystem(chai.ShapesDrawingSystem)

	logo_id := thisScene.NewEntityId()
	logo_transform := chai.VisualTransform{
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

	thisScene.NewRenderSystem(chai.FontRenderSystem)

	text_id := thisScene.NewEntityId()
	text_transform := chai.VisualTransform{
		Position:   chai.NewVector2f(25.0, Screen_Dims.Y-70.0),
		Dimensions: chai.Vector2fOne,
		Scale:      1.0,
		Rotation:   0.0,
	}
	fontRender := chai.FontRenderComponent{
		Text:            "Hi\nJoe",
		Scale:           3.0,
		Offset:          chai.Vector2fZero,
		Tint:            chai.WHITE,
		Fontbatch_atlas: &font_atlas,
		FontSettings:    &fontSettings,
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

	ldtk_levels := chai.ParseLdtk("Assets/Ldtk/newasset.ldtk")
	level = ldtk_levels.Get("Level0")
	level_size := level.GridSize().Scale(level.Tilesize())

	level_offset := chai.NewVector2f(-float32(level_size.X)/2, float32(level_size.Y)/2)

	chai.LoadTilemapLevel(thisScene, "Level0", ldtk_levels, level_offset)

	for y := 0; y < level.GridSize().Y; y++ {
		for x := 0; x < level.GridSize().X; x++ {
			_, exists := level.SolidTiles.AllItems()[chai.NewVector2i(x, y)]
			astar_nodes.PushBack(Node_AStar{
				neighbours: chai.NewList[*Node_AStar](),
				solid:      exists,
				grid_pos:   chai.NewVector2i(x, y),
				offset:     level_offset,
			})

		}
	}

	// Handle the Neighbours
	///////////////////////////////////////////
	for y := 0; y < level.GridSize().Y; y++ {
		for x := 0; x < level.GridSize().X; x++ {
			if x > 0 {
				astar_nodes.AllItems()[x+level.GridSize().X*y].neighbours.PushBack(&astar_nodes.AllItems()[x-1+level.GridSize().X*y])
			}
			if x < level.GridSize().X-1 {
				astar_nodes.AllItems()[x+level.GridSize().X*y].neighbours.PushBack(&astar_nodes.AllItems()[x+1+level.GridSize().X*y])
			}
			if y > 0 {
				astar_nodes.AllItems()[x+level.GridSize().X*y].neighbours.PushBack(&astar_nodes.AllItems()[x+level.GridSize().X*(y-1)])
			}
			if y > level.GridSize().Y {
				astar_nodes.AllItems()[x+level.GridSize().X*y].neighbours.PushBack(&astar_nodes.AllItems()[x+level.GridSize().X*(y+1)])
			}
		}
	}

	chai.SetGravity(chai.Vector2fZero)

	rb_id := thisScene.NewEntityId()
	rb_vt := chai.VisualTransform{
		Position:   level.Entities.Get("Player")[0].Position,
		Dimensions: chai.NewVector2f(12.0, 12.0),
		Rotation:   0.0,
		Scale:      1.0,
		Tint:       chai.WHITE,
	}
	rbSettings := chai.RigidBodySettings{
		IsTrigger:       false,
		BodyType:        chai.Type_BodyDynamic,
		ColliderShape:   chai.Shape_CircleBody,
		StartPosition:   rb_vt.Position,
		StartDimensions: rb_vt.Dimensions,
		StartRotation:   0.0,
		Mass:            250, Friction: 0.4, Elasticity: 0.2,
		PhysicsLayer:      0x10000000,
		ConstrainRotation: false,
	}
	rbComponent := chai.NewRigidBody(rb_id, &rbSettings)
	thisScene.AddComponents(rb_id, chai.ToComponent(rb_vt), chai.ToComponent(rbComponent), chai.ToComponent(PlayerComp{}))
	frc := chai.NewQuadComponent(thisScene, rb_id, chai.WHITE, false)
	thisScene.AddComponents(rb_id, chai.ToComponent(frc))

	thisScene.NewUpdateSystem(chai.RigidBodySystem)
	thisScene.NewUpdateSystem(CirclePlayerSystem)

	thisScene.NewUpdateSystem(BulletsSystem)

	chai.Shapes.LineWidth = 0.5

	for _, v := range level.Entities.Get("Target") {
		targets.PushBack(v.Position)
	}

	thisScene.NewUpdateSystem(HealthSystem)
	for h := 0; h < 10; h++ {
		thisScene.AddEntity(chai.ToComponent(Health{}))
	}
	hId := thisScene.AddEntity(chai.ToComponent(Health{}))
	thisScene.AddTag(hId, "Health")

	chai.ScaleView(4)
}

func SceneUpdate(thisScene *chai.Scene, dt float32) {

	// chai.ScrollView(chai.NewVector2f(0, y_axis).Scale(0.5))
	chai.IncreaseScaleU(chai.GetActionStrength("zoomin") - chai.GetActionStrength("zoomout"))

	chai.UIShapes.DrawFillRect(chai.NewVector2f(200.0, Screen_Dims.Y-70.0), chai.NewVector2f(400, 100), chai.BLACK)

	for _, t := range targets.AllItems() {
		chai.Shapes.DrawCircle(t, 4, chai.WHITE)
	}

	for _, at := range astar_nodes.Data {
		color := chai.NewRGBA8(255, 255, 255, 100)
		if at.solid {
			color = chai.NewRGBA8(255, 0, 0, 100)
		}
		node_position := chai.NewVector2f(float32(at.grid_pos.X), -float32(at.grid_pos.Y)).Scale(16.0).AddXY(8, 8).Add(at.offset)
		chai.Shapes.DrawFillRect(node_position, chai.NewVector2f(16.0, 16.0), color)
		// for _, n := range at.neighbours.AllItems() {
		// 	neighbour_position := chai.NewVector2f(float32(n.grid_pos.X), -float32(n.grid_pos.Y)).Scale(16).AddXY(8, 8).Add(n.offset)
		// 	chai.Shapes.DrawLine(node_position, neighbour_position, color)
		// }

	}
}

func main() {
	Screen_Dims = chai.NewVector2f(1240, 720)

	game = chai.App{
		Width:  int(Screen_Dims.X),
		Height: int(Screen_Dims.Y),
		Title:  "Test",

		OnStart: func() {

			fontSettings = chai.FontBatchSettings{
				FontSize: 24, DPI: 48, CharDistance: 4, LineHeight: 36, Arabic: false,
			}
			font_atlas = chai.LoadFontToAtlas("Assets/m5x7.ttf", &fontSettings)
			font_atlas.SpriteBatch = &chai.UISprites

			chai.LogF("STARTED\n")
			default_scene := chai.NewScene()
			default_scene.NewStartSystem(SceneSetup)
			default_scene.NewUpdateSystem(SceneUpdate)
			chai.ChangeScene(&default_scene)
		},
		OnUpdate: func(dt float32) {
		},
		OnDraw: func(dt float32) {
		},
		OnEvent: func(ae *chai.AppEvent) {
		},
	}

	chai.Run(&game)

}

type PlayerComp struct {
	flash_alpha   float32
	muzzle_alpha  float32
	target        chai.Vector2f
	target_vector chai.Vector2f
	target_angle  float32
}

func CirclePlayerSystem(_this_scene *chai.Scene, dt float32) {
	const MAX_AIM_DISTANCE = float32(72.0)

	chai.Iterate3[PlayerComp, chai.RigidBodyComponent, chai.VisualTransform](func(i chai.EntId, pc *PlayerComp, rbc *chai.RigidBodyComponent, vt *chai.VisualTransform) {
		// Input
		x_axis := chai.GetActionStrength("right") - chai.GetActionStrength("left")
		y_axis := chai.GetActionStrength("up") - chai.GetActionStrength("down")

		// Player movement set up with a faux-drag to simulate friction in a 2D top-down setting
		counterForceX := (rbc.GetVelocity().X * 0.1 * chai.BoolToFloat32(!chai.IsPressed("right") || !chai.IsPressed("left")))
		counterForceY := (rbc.GetVelocity().Y * 0.1 * chai.BoolToFloat32(!chai.IsPressed("up") || !chai.IsPressed("down")))
		rbc.SetVelocity(chai.NewVector2f(rbc.GetVelocity().X+5.0*x_axis-counterForceX, rbc.GetVelocity().Y+5.0*y_axis-counterForceY))

		// Camera Position
		chai.ScrollTo(chai.LerpVector2f(chai.Cam.GetPosition(), rbc.GetPosition(), 2.0*dt))

		// Stopping the world from affecting the angle of the dynamic object
		rbc.SetAngularVelocity(0.0)

		// Finding the closest target
		for _, t := range targets.AllItems() {
			if t.Distance(vt.Position) < pc.target.Distance(vt.Position) || pc.target == chai.Vector2fZero {
				// chai.LogF("T Distance: %v", t.Distance(vt.Position))
				pc.target = t
			}
		}
		hit := chai.LineCast(vt.Position, pc.target, 0x0ffffff0)
		if pc.target.Distance(vt.Position) <= MAX_AIM_DISTANCE && !hit.HasHit {
			pc.target_vector = pc.target.Subtract(vt.Position).Normalize()
			chai.Shapes.DrawLine(hit.OriginPoint, hit.HitPosition, chai.WHITE)
		} else if x_axis != 0.0 || y_axis != 0.0 {
			pc.target_vector = chai.NewVector2f(x_axis, y_axis)
		}
		// Rotating Linearly towards input
		// pc.target = chai.NewVector2f(x_axis, y_axis)
		pc.target_angle = chai.LerpRot(vt.Rotation, (pc.target_vector.Angle() * chai.Rad2Deg), 7*dt)
		rbc.SetRotation(pc.target_angle)

		muzzle_pos := vt.Position.Add(chai.Vector2fRight.Scale(6.0).Rotate(vt.Rotation, chai.Vector2fZero))
		chai.Shapes.DrawCircle(chai.GetMouseWorldPosition(), 2, chai.WHITE)

		// Muzzle Flash
		chai.Shapes.DrawFillTriangleRotated(muzzle_pos, chai.Vector2fOne.Scale(6.0), chai.NewRGBA8Float(1.0, 1.0, 1.0, pc.muzzle_alpha), vt.Rotation+90)

		// Shooting
		if chai.IsMouseJustPressed() {
			AddBullet(chai.GetCurrentScene(), vt.Position, vt.Rotation, 300)
			pc.flash_alpha = 0.2
			pc.muzzle_alpha = 1.0
		}

		pc.flash_alpha = chai.ClampFloat32(pc.flash_alpha-dt*4.0, 0.0, 1.0)
		pc.muzzle_alpha = chai.ClampFloat32(pc.muzzle_alpha-dt*2, 0.0, 1.0)

		// Screen Flash
		chai.UIShapes.DrawFillRect_Rect(chai.Rect{Position: chai.Vector2fZero, Size: Screen_Dims}, chai.NewRGBA8Float(1.0, 1.0, 1.0, pc.flash_alpha))
	})
}

type BulletComponent struct {
	active    bool
	speed     float32
	direction chai.Vector2f
}

func BulletsSystem(_this_scene *chai.Scene, dt float32) {
	chai.Iterate3[chai.VisualTransform, BulletComponent, chai.RigidBodyComponent](func(i chai.EntId, vt *chai.VisualTransform, bc *BulletComponent, rbc *chai.RigidBodyComponent) {
		rbc.SetVelocity(bc.direction.Scale(bc.speed * chai.BoolToFloat32(bc.active)))

		// rbc.SetRotation(vt.Rotation)
	})
}

func AddBullet(_thisScene *chai.Scene, _center chai.Vector2f, _rotation float32, _speed float32) {
	bull_id := _thisScene.NewEntityId()
	vt := chai.VisualTransform{
		Position:   _center,
		Dimensions: chai.NewVector2f(1.5, 0.75).Scale(2.0),
		Scale:      1.0,
		Rotation:   _rotation,
		Tint:       chai.NewRGBA8(255, 200, 50, 255),
	}

	rbc := chai.NewRigidBody(bull_id, &chai.RigidBodySettings{
		IsTrigger:       true,
		BodyType:        chai.Type_BodyDynamic,
		ColliderShape:   chai.Shape_CircleBody,
		StartPosition:   vt.Position,
		Offset:          chai.Vector2fZero,
		StartDimensions: vt.Dimensions.Scale(0.5),
		Mass:            5, Friction: 0.35, Elasticity: 0.0,
		ConstrainRotation: false,
		PhysicsLayer:      0x00000001,
		StartRotation:     vt.Rotation,
	})

	rbc.OnCollisionBegin.AddListener(func(c chai.Collision) {
		ent := c.EntA
		if chai.GetComponentPtr[PlayerComp](chai.GetCurrentScene(), c.EntB) != nil {
			return
		}
		bc := chai.GetComponentPtr[BulletComponent](chai.GetCurrentScene(), ent)
		if bc != nil {
			bc.active = false
			RemoveBullet(ent)
			chai.Destroy(chai.GetCurrentScene(), ent)
		}
	})

	bc := BulletComponent{
		active:    true,
		speed:     _speed,
		direction: chai.Vector2fRight.Rotate(vt.Rotation, chai.Vector2fZero),
	}

	_thisScene.AddComponents(bull_id, chai.ToComponent(vt), chai.ToComponent(rbc), chai.ToComponent(bc))
	qComp := chai.NewQuadComponent(_thisScene, bull_id, vt.Tint, false)
	_thisScene.AddComponents(bull_id, chai.ToComponent(qComp))
}

func RemoveBullet(_bull_id chai.EntId) {
	// quadComp := chai.GetComponentPtr[chai.FillRectRenderComponent](chai.GetCurrentScene(), _bull_id)

	// chai.DynamicRenderQuadTreeContainer.Remove(chai.DynamicRenderQuadTreeContainer.AllItems().Get())
	chai.DynamicRenderQuadTreeContainer.RemoveWithIndex(int64(_bull_id))
	// vt := chai.GetComponentPtr[chai.VisualTransform](chai.GetCurrentScene(), _bull_id)

	// Make shift delete object
	// Possible solution is to have the index be the entity_id in the Tree
	// q := chai.GetDynamicQuadsInRect(chai.Rect{Position: vt.Position.Subtract(chai.NewVector2f(5.0, 5.0)), Size: vt.Dimensions.Scale(2)})
	// for _, v := range q.Data {
	// 	chai.DynamicRenderQuadTreeContainer.Remove(v)
	// }
}

func Solve_AStar(_start, _end chai.Vector2f) {

}

type Node_AStar struct {
	neighbours chai.List[*Node_AStar]
	solid      bool
	grid_pos   chai.Vector2i
	parent     *Node_AStar
	offset     chai.Vector2f
}

type Health struct {
	value float32
}

func (h *Health) Setup(_this_scene *chai.Scene) {
	h.value = 0.0
}

func HealthSystem(_this_scene *chai.Scene, dt float32) {
	chai.Iterate1[Health](func(i chai.EntId, h *Health) {
		h.value += dt
		if _this_scene.HasTag(i, "Health") {
			chai.LogF("Health (%d): %v", i, h.value)
		}
	})
}
