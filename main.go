package main

import (
	chai "github.com/mhamedGd/chai"
)

var sprite_animation_app chai.App = chai.App{
	Title:  "Sprite Animation",
	Width:  600,
	Height: 1000,
	OnStart: func() {
		_sp_anim_scene := chai.NewScene()
		_sp_anim_scene.NewStartSystem(func(_this_scene *chai.Scene) {
			_this_scene.Background = chai.NewRGBA8(255, 100, 200, 255)

			_knight_id := _this_scene.NewEntityId()
			_vt := chai.VisualTransform{
				Dimensions: chai.NewVector2f(120.0, 80.0),
				Rotation:   0.0,
				Scale:      1.0,
				Tint:       chai.NewRGBA8(255, 255, 255, 255),
			}
			// Creating a sprite sheet from the image provided (Basically dividing the sprite sheet based on the image's coloumns and rows)
			_sprite_sheet := chai.NewSpriteSheet(chai.LoadPng("Assets/_AttackCombo.png", &chai.TextureSettings{Filter: chai.TEXTURE_FILTER_NEAREST}), 120, 80)

			// Creating a new sprite sheet animation component
			_sprite_animation := chai.NewSpriteAnimationComponent(&_sprite_sheet)
			// Creating a new animation withtin the component
			_sprite_animation.NewAnimation("Attack")

			// Registering the keyframes of the animation by providing the frames coordinates within the image
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(0, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(1, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(2, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(3, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(4, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(5, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(6, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(7, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(8, 0))
			_sprite_animation.RegisterFrame("Attack", chai.NewVector2i(9, 0))

			// Specifying the current animation
			_sprite_animation.CurrentAnimation = "Attack"
			// Specifying the animation speed (frames per second)
			_sprite_animation.AnimationSpeed = 16

			// Attaching the components to the knight id
			_this_scene.AddComponents(_knight_id, chai.ToComponent(_vt), chai.ToComponent(_sprite_animation))
			_this_scene.AddComponents(_knight_id, chai.ToComponent(chai.NewSpriteComponent(_this_scene, _knight_id, _vt, &_sprite_sheet.Texture, false)))

			// Adding the sprite animation system to the Update Runtime of the scene
			_this_scene.NewUpdateSystem(chai.SpriteAnimationSystem)

			chai.ScaleView(2)
		})
		chai.ChangeScene(&_sp_anim_scene)
	},
}

func main() {
	chai.Run(&sprite_animation_app)
}
