package main

import (
	chai "github.com/mhamedGd/chai"
	"github.com/mhamedGd/chai/ecs"
	. "github.com/mhamedGd/chai/math"
)

// An App is everything in your game
var load_sprites_app chai.App = chai.App{
	Title:  "دكتور ابو طرشمان",
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
		test_scene.NewUpdateSystem(chai.TweenAnimatorSystem)
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
		chai.ScrollView(NewVector2f(_camera_x_axis, _camera_y_axis).Scale(0.5))
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

	chai.ScaleView(4)

	_abu_tarshaman := _this_scene.NewEntityId()

	_transform_comp := chai.VisualTransform{
		Dimensions: NewVector2f(56, 33),
		Scale:      0.0,
		Tint:       chai.BLACK,
		UV1:        Vector2fZero,
		UV2:        Vector2fOne,
	}
	_abu_png := chai.LoadPng("Assets/Chai_Logo_transparent.png", chai.TextureSettings{Filter: chai.TEXTURE_FILTER_LINEAR})
	_abu_sprite := chai.NewSpriteComponent(_this_scene, _abu_tarshaman, _transform_comp, &_abu_png, false)

	_abu_animation := chai.NewAnimationComponentVector2f()
	_abu_animation.NewTweenAnimationVector2f("Vertical", true)
	_abu_animation.RegisterKeyframe("Vertical", 0.0, Vector2fZero)
	_abu_animation.RegisterKeyframe("Vertical", 0.5, Vector2fUp)
	_abu_animation.RegisterKeyframe("Vertical", 1.0, Vector2fZero)
	_abu_animation.RegisterKeyframe("Vertical", 1.5, Vector2fDown)
	_abu_animation.RegisterKeyframe("Vertical", 2.0, Vector2fZero)
	_abu_animation.Play("Vertical")

	_this_scene.AddComponents(_abu_tarshaman, chai.ToComponent(_transform_comp), chai.ToComponent(_abu_sprite))
}

func MoveAbuTarshamanSystem(_this_scene *chai.Scene, _dt float32) {
	chai.Iterate2[chai.VisualTransform, chai.AnimationComponent[Vector2f]](func(i ecs.Id, t *chai.VisualTransform, ac *chai.AnimationComponent[Vector2f]) {
		t.Position = ac.GetCurrentValue("Vertical")
	})
}
