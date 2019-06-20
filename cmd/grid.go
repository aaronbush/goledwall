package cmd

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Screen struct {
	Height      int32
	Width       int32
	Grid        Grid
	StatusArea  StatusArea
	ControlArea ControlArea
}

type StatusArea struct {
	Origin rl.Vector2
	Height int32
	Width  int32
}

type ControlArea struct {
	Origin rl.Vector2
	Height int32
	Width  int32
	Info   map[string]interface{}
}

type Grid struct {
	Origin          rl.Vector2
	Height          int32
	Width           int32
	NumRows         uint8
	NumColumns      uint8
	Spacing         int32
	Color           rl.Color
	GridRowByColumn [][]GridSquare // rows x columns
	LastModified    time.Time
}

type GridSquare struct {
	Origin    rl.Vector2
	Row       uint8
	Column    uint8
	CreatedAt time.Time
	Color     rl.Color
}

func NewScreen(width, height int32) Screen {
	screen := Screen{Width: width, Height: height}
	return screen
}

func NewControlArea(origin rl.Vector2, width, height int32) ControlArea {
	controlArea := ControlArea{Width: width, Height: height, Origin: origin, Info: make(map[string]interface{})}
	controlArea.Info["r"] = new(int)
	controlArea.Info["g"] = new(int)
	controlArea.Info["b"] = new(int)
	return controlArea
}

func NewStatusArea(origin rl.Vector2, width, height int32) StatusArea {
	statusArea := StatusArea{Width: width, Height: height, Origin: origin}
	return statusArea
}

func NewGrid(origin rl.Vector2, numColumns, numRows uint8, spacing int32) Grid {
	gridSquares := make([][]GridSquare, numRows) // rows
	for i := range gridSquares {
		gridSquares[i] = make([]GridSquare, numColumns) // columns
		for j := range gridSquares[i] {
			origin := rl.NewVector2(float32(int32(j)*spacing)+origin.X, float32(int32(i)*spacing)+origin.Y)
			gridSquares[i][j] = GridSquare{Column: uint8(j), Row: uint8(i), CreatedAt: time.Now(), Origin: origin}
		}
	}
	grid := Grid{NumRows: numRows, NumColumns: numColumns, Spacing: spacing, GridRowByColumn: gridSquares, Color: rl.Gray, Width: int32(numColumns) * spacing, Height: int32(numRows) * spacing}

	return grid
}
