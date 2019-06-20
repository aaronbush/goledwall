package cmd

import (
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	maxRGB = 255
)

func (controlArea *ControlArea) drawColorInputs(red, green, blue *int) rl.Color {
	position := controlArea.Origin
	spacingFloat := float32(5.0)

	position.X += 5
	drawColorInput("red", red, position)
	position.Y += 45
	drawColorInput("green", green, position)
	position.Y += 45
	drawColorInput("blue", blue, position)
	position.Y += 45

	color := makeColor(*red, *green, *blue, 255)
	rl.DrawRectangleV(position, rl.NewVector2(spacingFloat, spacingFloat), color)
	return color
}

func drawColorInput(name string, colorValue *int, position rl.Vector2) {
	rg.Label(rl.NewRectangle(position.X, position.Y, 50, 20), name)
	color := rg.TextBox(rl.NewRectangle(position.X, position.Y+20, 50, 20), strconv.Itoa(*colorValue))
	*colorValue, _ = strconv.Atoi(color)
	if *colorValue > maxRGB {
		*colorValue = maxRGB
	}
}

func makeColor(red, green, blue, alpha int) rl.Color {
	return rl.NewColor(uint8(red), uint8(green), uint8(blue), uint8(alpha))
}
