package main

import (
	chai "github.com/mhamedGd/chai"
	"github.com/mhamedGd/chai/ecs"
)

var ldtk_app chai.App = chai.App{
	Title:  "LDTK",
	Width:  720,
	Height: 720,
	OnStart: func() {
		chai.BindInput("left", chai.KEY_ARROWLEFT)
		chai.BindInput("right", chai.KEY_ARROWRIGHT)
		chai.BindInput("jump", chai.KEY_SPACE)

		_ldtk_scene := chai.NewScene()
		_ldtk_scene.NewStartSystem(func(_this_scene *chai.Scene) {
			_this_scene.Background = chai.NewRGBA8(23, 28, 57, 255)
			chai.ScaleView(3)

			// Loading and Parsing the .ldtk file with its image
			_ldtk_levels := chai.ParseLdtk("Assets/Ldtk/test.ldtk")

			// Loading the level labeled "Level_1"
			_level := _ldtk_levels.Get("Level_1")

			// Setting the Level_1 to the name we want it to be
			_ldtk_levels.Set("Ldtk Level", _level)
			_level_offset := chai.NewVector2f(-112, 105)

			// Loading the tiles and adding them to the StaticQuadTreeContainer responsible for rendering and culling
			chai.LoadTilemapLevel(_this_scene, "Ldtk Level", _ldtk_levels, 5, _level_offset)

			// Adding 2 Update Runtime systems,
			//	- Responsible for updating the movement of Dynamic Bodies
			_this_scene.NewUpdateSystem(chai.DynamicBodySystem)
			//	- A Custom System for player movement
			_this_scene.NewUpdateSystem(PlayerMoveSystem)

			// Creating the player as a box
			_box_id := _this_scene.NewEntityId()
			_box_vt := chai.VisualTransform{
				Position:   _level.Entities.Get("Player")[0].Position,
				Dimensions: chai.NewVector2f(6.0, 6.0),
				Rotation:   0.0,
				Scale:      1.0,
				Tint:       chai.NewRGBA8(100, 255, 40, 255),
			}

			// Creatig the Dynamic Body Component
			_box_rb := chai.NewDynamicBodyComponent(_box_id, _box_vt, &chai.DynamicBodySettings{
				ColliderShape:   chai.SHAPE_RECTBODY,
				StartPosition:   _box_vt.Position,
				StartDimensions: _box_vt.Dimensions,
				StartRotation:   _box_vt.Rotation,
				Mass:            10, Friction: 0.2, Elasticity: 0.1,
				ConstrainRotation: true,
				// Making the player collide with the level while also being distinct so that the Raycast doesn't catch it
				PhysicsLayer: chai.PHYSICS_LAYER_1 | chai.PHYSICS_LAYER_2,
			})

			// Creating the Quad Render Component
			_box_dr := chai.NewQuadComponent(_this_scene, _box_id, _box_vt, false)

			// Attaching the Components to the Box Id in this Scene
			_this_scene.AddComponents(_box_id, chai.ToComponent(_box_vt), chai.ToComponent(_box_dr), chai.ToComponent(_box_rb))

		})
		chai.ChangeScene(&_ldtk_scene)
	},
}

// Declaring an Update Runtime system
func PlayerMoveSystem(_this_scene *chai.Scene, _dt float32) {

	// Iterating over every object with a DynamicBodyComponent in the Scene
	chai.Iterate1[chai.DynamicBodyComponent](func(i ecs.Id, dbc *chai.DynamicBodyComponent) {
		// Creating an input movement axis
		_x_axis := chai.GetActionStrength("right") - chai.GetActionStrength("left")
		const _SPEED = 20

		// Setting the DynamicBodyComponent velocity (_x_axis * _SPEED, velocity.y)
		dbc.SetVelocity(chai.Vector2fRight.Scale(_SPEED * _x_axis).Add(chai.Vector2fUp.Scale(dbc.GetVelocity().Y)))

		// Raycasting to only the First Layer (All ldtk loaded levels are on the first phyiscal layer)
		_hit := chai.RayCast(dbc.GetPosition(), chai.Vector2fDown, 8, chai.PHYSICS_LAYER_1)
		if chai.IsJustPressed("jump") && _hit.HasHit {
			dbc.ApplyImpulse(chai.Vector2fUp.Scale(17_000), chai.Vector2fZero)
		}
	})
}
