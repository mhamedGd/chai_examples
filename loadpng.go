package main

import chai "github.com/mhamedGd/chai"

// An App is everything in your game
var load_sprites_app chai.App = chai.App{
	Title:  "Helloo",
	Width:  1920,
	Height: 1080,
	OnStart: func() {
		// Binding a refrence string to a certain Key
		// Binding an input is global to an app, Which means different scenes
		// of the same app will shape the same inputs map
		chai.BindInput("left", chai.KEY_A)
		chai.BindInput("right", chai.KEY_D)
		chai.BindInput("up", chai.KEY_W)
		chai.BindInput("down", chai.KEY_S)

		// Creating a new Scene and Transitioning into it
		test_scene := chai.NewScene()
		// Scenes have 3 runtimes
		//	- OnStart(_this_scene *chai.Scene)
		//	- OnUpdate(_this_scene *chai.Scene, _dt float32)
		//	- OnDraw(_this_scene *chai.Scene, _dt float32)
		//	You can subscribe a function to any of those runtimes
		test_scene.NewStartSystem(LoadSpritesSceneStart)
		chai.ChangeScene(&test_scene)
	},
	OnUpdate: func(f float32) {
		// You can access the values of pressed bound inputs by accessing
		// chai.GetActionStrength(_input_name string)
		// chai.IsPressed(_input_name string)
		// chai.IsJustPressed(_input_name string)
		// chai.IsJustReleased(_input_name string)
		_camera_x_axis := chai.GetActionStrength("right") - chai.GetActionStrength("left")
		_camera_y_axis := chai.GetActionStrength("up") - chai.GetActionStrength("down")

		// This offests the camera by the vector provided
		// ScrollTo(_new_pos chai.Vector2f) is the fuction used to set the camera's new position
		chai.ScrollView(chai.NewVector2f(_camera_x_axis, _camera_y_axis).Scale(0.5))
	},
	OnDraw: func(f float32) {},
	OnEvent: func(ae *chai.AppEvent) {
		// You can access the console in the web browser by pressing (Ctrl/Cmd + L.Shift + C)
		if ae.Type == chai.JS_KEYDOWN {
			// Printing to the console
			// You can aslo call chai.WarningF(_message string)
			// or chai.ErrorF(_message string) - This will also cause the WebApp to panic
			chai.LogF(ae.Key)
		}
	},
}

func LoadSpritesSceneStart(_this_scene *chai.Scene) {
	// Setting the Background Color
	_this_scene.Background = chai.NewRGBA8(255, 255, 255, 255)

	// Scaling the Camera Viewport
	chai.ScaleView(16)

	// Chai is an ECS based engine, Which means it operates by separting the concept of object into 3 differet
	// parts:
	//		- Entity: Which is basically an Identifier
	//		- Component: A Struct of any type holding information
	//		- System: A function that uses the Componenets attached to Entities to do perform various behaviours
	_buddy_logo_id := _this_scene.NewEntityId()

	// VisualTransform is a component that holds Positional and Visual information of an entity
	_buddy_logo_vt := chai.VisualTransform{
		Position:   chai.Vector2fZero,
		Z:          1,
		Dimensions: chai.NewVector2f(8, 8),
		Rotation:   0.0,
		Scale:      1.0,
		Tint:       chai.WHITE,
		UV1:        chai.Vector2fZero,
		UV2:        chai.Vector2fOne,
	}

	// Currently Chai only supports PNGs
	// TextureSettings only contains the type of filter, in the it will hold more info
	_buddy_logo_texture := chai.LoadPng("Assets/kenney.png", &chai.TextureSettings{Filter: chai.TEXTURE_FILTER_NEAREST})
	// Creating a new sprite component in the world by using the information provided in the entity id and the visual trasform component
	_buddy_logo_sprite := chai.NewSpriteComponent(_this_scene, _buddy_logo_id, _buddy_logo_vt, &_buddy_logo_texture, true)
	// Attaching the previous 2 components (VisualTransform, SpriteComponent) to the _buddy_logo_id Entity
	_this_scene.AddComponents(_buddy_logo_id, chai.ToComponent(_buddy_logo_vt), chai.ToComponent(_buddy_logo_sprite))

	// Creating a new Entity with similar info to the previous, except this one loads the Chai Logo
	_chai_logo_id := _this_scene.NewEntityId()
	_chai_logo_vt := chai.VisualTransform{
		Position:   chai.Vector2fZero,
		Z:          5,
		Dimensions: chai.NewVector2f(56, 33),
		Rotation:   0.0,
		Scale:      1.0,
		Tint:       chai.BLACK,
		UV1:        chai.Vector2fZero,
		UV2:        chai.Vector2fOne,
	}
	_chai_logo_texture := chai.LoadPng("Assets/Chai_Logo_transparent.png", &chai.TextureSettings{Filter: chai.TEXTURE_FILTER_NEAREST})
	_chai_logo_sprite := chai.NewSpriteComponent(_this_scene, _chai_logo_id, _chai_logo_vt, &_chai_logo_texture, true)
	_this_scene.AddComponents(_chai_logo_id, chai.ToComponent(_chai_logo_vt), chai.ToComponent(_chai_logo_sprite))
}
