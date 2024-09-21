package main

import (
	chai "github.com/mhamedGd/chai"
	"github.com/mhamedGd/chai/ecs"
)

var tween_animation_app chai.App = chai.App{
	Title:  "Animation",
	Width:  800,
	Height: 600,
	OnStart: func() {
		_tween_scene := chai.NewScene()
		chai.ChangeScene(&_tween_scene)
		_tween_scene.NewStartSystem(func(_this_scene *chai.Scene) {
			chai.ScaleView(4)

			// Here we will create an new Entity and Attach an AnimationComponent that outputs Vector2f values
			_sprite_moved_id := _this_scene.NewEntityId()
			_sprite_moved_vt := chai.VisualTransform{
				Position:   chai.Vector2fZero,
				Dimensions: chai.NewVector2f(8, 8),
				Z:          5,
				Rotation:   0.0,
				Scale:      1.0,
				Tint:       chai.WHITE,
				UV1:        chai.Vector2fZero,
				UV2:        chai.Vector2fOne,
			}
			_sprite_tex := chai.LoadPng("Assets/kenney.png", &chai.TextureSettings{Filter: chai.TEXTURE_FILTER_NEAREST})
			// Notice how the last boolean is false, that's because we want the Sprite to move around with the Entity (Not Static)
			_sprite_moved_comp := chai.NewSpriteComponent(_this_scene, _sprite_moved_id, _sprite_moved_vt, &_sprite_tex, false)

			// We start by decalring a new animation component
			// We create inside of it a new TweenAnimation of Vector2f with a name and loop value
			// We register wanted keyframes and the play it
			// but wait, How do we actually connect the value of this Animation Component with the Entity's Position?
			_sprite_moved_animation := chai.NewAnimationComponentVector2f()
			_sprite_moved_animation.NewTweenAnimationVector2f("Circle", true)
			_sprite_moved_animation.RegisterKeyframe("Circle", 0.0, chai.NewVector2f(0.0, 16.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 0.5, chai.NewVector2f(12.0, 12.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 1.0, chai.NewVector2f(16.0, 0.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 1.5, chai.NewVector2f(12.0, -12.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 2.0, chai.NewVector2f(0.0, -16.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 2.5, chai.NewVector2f(-12.0, -12.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 3.0, chai.NewVector2f(-16.0, 0.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 3.5, chai.NewVector2f(-12.0, 12.0))
			_sprite_moved_animation.RegisterKeyframe("Circle", 4.0, chai.NewVector2f(0.0, 16.0))
			_sprite_moved_animation.Play("Circle")

			_this_scene.AddComponents(_sprite_moved_id, chai.ToComponent(_sprite_moved_vt), chai.ToComponent(_sprite_moved_comp), chai.ToComponent(_sprite_moved_animation))
		})
		// We Create A Custom System that runs on the OnUpdate runtime and takes the value output by
		// the animation component and assigns it to the Entity's Position
		_tween_scene.NewUpdateSystem(AnimateBuddyPosition)
		// We also need to Add the System TweenAnimatorSystem in order to update the current keyframe
		_tween_scene.NewUpdateSystem(chai.TweenAnimatorSystem)
	},
}

// When declaring a new OnUpdate system, We must have 2 parameters. The first being a reference to the
// scene in which the system is running, The second is the deltaTime.
func AnimateBuddyPosition(_this_scene *chai.Scene, _dt float32) {
	// An ECS Iterator that iterates over all entities that have These 2 Specified components and performs
	// certain commands, specified by you by providing an appropriate function.
	chai.Iterate2[chai.VisualTransform, chai.AnimationComponent[chai.Vector2f]](func(i ecs.Id, vt *chai.VisualTransform, ac *chai.AnimationComponent[chai.Vector2f]) {
		_anim_value := ac.GetCurrentValue("Circle")
		vt.Position = _anim_value
	})
}
