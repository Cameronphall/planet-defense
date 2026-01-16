package main

import rl "github.com/gen2brain/raylib-go/raylib"

//rl "github.com/gen2brain/raylib-go/raylib"

type Planet struct {
	Pos rl.Vector2
	SrcRec rl.Rectangle
	DestRec rl.Rectangle
	Texture rl.Texture2D
	Radius float32
	MaxHealth float32
	CurrentHealth float32
	Alive bool
}

func DrawHealthBar(x, y, width, height, health, maxHealth float32) {
    
    rl.DrawRectangle(int32(x), int32(y), int32(width), int32(height), rl.DarkGray)

 
    ratio := health / maxHealth
    var barColor rl.Color

    if ratio > 0.5 {
        t := (1 - ratio) * 2 
        barColor = rl.ColorLerp(rl.Green, rl.Yellow, t)
    } else {
        
        t := (0.5 - ratio) * 2 
        barColor = rl.ColorLerp(rl.Yellow, rl.Red, t)
    }

 
    barWidth := ratio * width
    rl.DrawRectangle(int32(x), int32(y), int32(barWidth), int32(height), barColor)
    rl.DrawRectangleLines(int32(x), int32(y), int32(width), int32(height), rl.Black)
}

func NewPlanet(pos rl.Vector2, src rl.Rectangle, dest rl.Rectangle, text rl.Texture2D, radius float32, maxHealth float32, health float32) Planet{
	return Planet{
		Pos: pos,
		SrcRec: src,
		DestRec: dest,
		Texture: text,
		Radius: radius,
		MaxHealth: maxHealth,
		CurrentHealth: health,
		Alive: true,
	}
}

func (p Planet)DrawPlanet(){
	rl.DrawTexturePro(p.Texture, p.SrcRec, p.DestRec, rl.NewVector2(0, 0), 0, rl.White)
}