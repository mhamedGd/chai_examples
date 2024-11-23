package main

import (
	chai "github.com/mhamedGd/chai"
	"github.com/mhamedGd/chai/ecs"
	chmath "github.com/mhamedGd/chai/math"
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
				Dimensions: chmath.NewVector2f(120.0, 80.0),
				Rotation:   0.0,
				Scale:      1.0,
				Tint:       chai.NewRGBA8(255, 255, 255, 255),
			}
			// Creating a sprite sheet from the image provided (Basically dividing the sprite sheet based on the single tile's size)
			_sprite_sheet := chai.NewSpriteSheet(chai.LoadPng("Assets/_AttackCombo.png", chai.TextureSettings{Filter: chai.TEXTURE_FILTER_NEAREST}), 120, 80)

			// Creating a new sprite sheet animation component
			_sprite_animation := chai.NewSpriteAnimationComponent(&_sprite_sheet)
			// Creating a new animation withtin the component
			_sprite_animation.NewAnimation("Attack")

			// Registering the keyframes of the animation by providing the frames coordinates within the image
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(0, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(1, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(2, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(3, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(4, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(5, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(6, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(7, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(8, 0))
			_sprite_animation.RegisterFrame("Attack", chmath.NewVector2i(9, 0))

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

// //////////////////////////////////////////////////////////////////////////////////////
func SceneStartSystem(_this_scene *chai.Scene) {

	_db_entity := _this_scene.NewEntityId()
	_db_vt := chai.VisualTransform{
		Position:   chmath.NewVector2f(2.0, -1.0),
		Z:          0.0,
		Dimensions: chmath.NewVector2f(1.0, 1.0),
		Rotation:   0.0,
		Scale:      0.5,
		Tint:       chai.RED,
		UV1:        chmath.Vector2fZero,
		UV2:        chmath.Vector2fOne,
	}
	_texture := chai.LoadPng("chai_logo.png", chai.TextureSettings{
		Filter: chai.TEXTURE_FILTER_LINEAR,
	})
	_texture.OverridePixelsPerMeter(256)
	_db_sprite_component := chai.NewSpriteComponent(_this_scene, _db_entity, _db_vt, &_texture, false)

	_db_settings := chai.DynamicBodySettings{
		IsTrigger:       false,
		ColliderShape:   chai.SHAPE_CIRCLEBODY,
		StartPosition:   _db_vt.Position,
		Offset:          chmath.Vector2fZero,
		StartDimensions: _db_vt.Dimensions,
		StartRotation:   _db_vt.Rotation,
		Mass:            50, Friction: 0.2, Elasticity: 0.25,
		GravityScale:      0.2,
		ConstrainRotation: false,
		PhysicsLayer:      chai.PHYSICS_LAYER_ALL,
	}
	_db_component := chai.NewDynamicBodyComponent(_db_entity, _db_vt, &_db_settings)

	_ps_component := chai.NewParticlesShapeComponent(500, 0.0, func(f float32, p *chai.Particle) {
		p.Size *= p.LifePercentage
		p.Velocity = p.Velocity.Scale(p.LifePercentage)
		p.Color.SetColorAFloat32(p.LifePercentage)
	})

	_this_scene.AddComponents(_db_entity, chai.ToComponent(_db_vt), chai.ToComponent(_db_sprite_component), chai.ToComponent(_db_component), chai.ToComponent(_ps_component))

	_sb_entity := _this_scene.NewEntityId()
	_sb_vt := chai.VisualTransform{
		Position:   chmath.NewVector2f(-4, -2.5),
		Z:          15.0,
		Dimensions: chmath.NewVector2f(8, 1),
		Rotation:   0.0,
		Scale:      1.0,
		Tint:       chai.WHITE,
		UV1:        chmath.Vector2fZero,
		UV2:        chmath.Vector2fOne,
	}

	_sb_settings := chai.StaticBodySettings{
		IsTrigger:       false,
		ColliderShape:   chai.SHAPE_RECTBODY,
		Offset:          chmath.Vector2fZero,
		StartDimensions: _sb_vt.Dimensions,
		StartRotation:   _sb_vt.Rotation,
		Friction:        0.3, Elasticity: 0.25,
		PhysicsLayer: chai.PHYSICS_LAYER_ALL,
	}
	_sb_component := chai.NewStaticBodyComponent(_sb_entity, _sb_vt, &_sb_settings)

	_sb_quad := chai.NewQuadComponent(_this_scene, _sb_entity, _sb_vt, true)

	_this_scene.AddComponents(_sb_entity, chai.ToComponent(_sb_vt), chai.ToComponent(_sb_component), chai.ToComponent(_sb_quad))

	_levels := chai.ParseLdtk("Assets/Ldtk/test.ldtk")
	_lev := _levels.Get("Level_1")
	_levels.Set("Level_1", _lev)
	chai.LoadTilemapLevel(_this_scene, "Level_1", _levels, 0.0, 1.0, chmath.NewVector2f(-6, 6))

	_this_scene.NewUpdateSystem(chai.DynamicBodySystem)
	_this_scene.NewUpdateSystem(MoveDynamicBodySystem)
	_this_scene.NewUpdateSystem(chai.ParticlesShapeUpdateSystem)

	_this_scene.NewRenderSystem(chai.ParticlesShapeRenderSystem)
	_this_scene.NewRenderSystem(TestDraw)
	_this_scene.NewRenderSystem(chai.DebugBodyDrawSystem)

	chai.ScaleView(float32(chai.GetCanvasWidth()) / float32(8))
}

func MoveDynamicBodySystem(_this_scene *chai.Scene, _dt float32) {
	chai.Iterate2[chai.DynamicBodyComponent, chai.ParticlesShapeComponent](func(i ecs.Id, dbc *chai.DynamicBodyComponent, pc *chai.ParticlesShapeComponent) {
		_x_axis := chai.GetActionStrength("right") - chai.GetActionStrength("left")
		_speed := float32(6.0)
		_velocity := dbc.GetVelocity()

		if chai.IsPressed("jump") {
			_velocity.Y = 12.0
			pc.AddParticles(12, chai.PARTICLES_SHAPE_QUAD, chai.PARTICLES_CIRCLESPREAD, 1.5, 0.1, dbc.GetPosition(), chai.WHITE, 0.4, 0.0)
		}
		dbc.SetVelocity(chmath.NewVector2f(_x_axis*_speed, _velocity.Y))

		chai.ScrollTo(dbc.GetPosition())

		chai.LogF("%v", 1/_dt)
	})
	chai.IncreaseScaleU((chai.GetActionStrength("zoomin") - chai.GetActionStrength("zoomout")) * 12.0)
}

func TestDraw(_thisScene *chai.Scene, _dt float32) {
	chai.DrawLine(chmath.Vector2fZero, chmath.NewVector2f(0, -2), chai.NewRGBA8(255, 255, 0, 100), 10.0)
	chai.DrawRect(chmath.Vector2fRight.Scale(3), chmath.Vector2fOne.Scale(1.0), chai.RED, 0.0, 40.0)
}

var app chai.App = chai.App{
	Title:          "Docs",
	Width:          1080,
	Height:         1080,
	PixelsPerMeter: 8,
	OnStart: func() {
		chai.BindInput("left", chai.KEY_A)
		chai.BindInput("right", chai.KEY_D)
		chai.BindInput("jump", chai.KEY_SPACE)
		chai.BindInput("zoomin", chai.KEY_E)
		chai.BindInput("zoomout", chai.KEY_Q)

		_test_scene := chai.NewScene()
		_test_scene.NewStartSystem(SceneStartSystem)

		// Transition into this scene
		chai.ChangeScene(&_test_scene)
	},
}

func main() {
	chai.Run(&app)
}
