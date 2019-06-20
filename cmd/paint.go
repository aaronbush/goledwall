// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/spf13/cobra"
)

var (
	StartSentinal    = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	LastBinaryOutput time.Time
)

type BinaryOutput struct {
	StartSentinal [6]byte
	NumLEDs       uint16
}

type BinaryOutputSquare struct {
	Row        uint8
	Column     uint8
	Red        uint8
	Green      uint8
	Blue       uint8
	Brightness uint8
}

// paintCmd represents the paint command
var paintCmd = &cobra.Command{
	Use:   "paint",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: paint,
}

type Game struct {
	PrimaryScreen Screen
	Spacing       float32
	TargetFPS     int32
	FadeMode      bool
	DecayMode     bool
	LogMode       bool
	FillMode      bool
}

func init() {
	rootCmd.AddCommand(paintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paintCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func paint(cmd *cobra.Command, args []string) error {
	cols := uint8(20)
	rows := uint8(36)
	spacing := int32(20)
	fps := int32(10)
	logMode := true
	controlWidth := int32(80)
	statusHeight := int32(80)

	// game screen setup

	screen := NewScreen(int32(cols)*spacing+controlWidth, int32(rows)*spacing+statusHeight)
	controlArea := NewControlArea(rl.NewVector2(float32(int32(cols)*spacing), 0), controlWidth, screen.Height)
	grid := NewGrid(rl.NewVector2(0, 0), cols, rows, spacing)
	statusArea := NewStatusArea(rl.NewVector2(0, float32(int32(rows)*spacing)), screen.Width, statusHeight)

	statusArea.Origin = rl.NewVector2(0, float32(screen.Height-statusArea.Height))

	screen.ControlArea = controlArea
	screen.Grid = grid
	screen.StatusArea = statusArea

	gameContext := Game{TargetFPS: fps, LogMode: logMode}
	gameContext.PrimaryScreen = screen

	// main control loop
	rl.InitWindow(screen.Width, screen.Height, "pixel drawing")
	rl.SetTargetFPS(gameContext.TargetFPS)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blank)

		gameContext.handleInput()
		gameContext.Draw()
		rl.EndDrawing()
	}
	rl.CloseWindow()
	return nil
}

func (game *Game) Draw() {
	game.PrimaryScreen.Draw(*game)

	if game.LogMode {
		if game.PrimaryScreen.Grid.LastModified.After(LastBinaryOutput) {
			data, err := game.PrimaryScreen.Grid.MarshalBinary()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", hex.Dump(data))
			LastBinaryOutput = time.Now()
		}
	}
}

func (screen *Screen) Draw(gameContext Game) {
	// TODO: move this out of here
	redValue := screen.ControlArea.Info["r"].(*int)
	greenValue := screen.ControlArea.Info["g"].(*int)
	blueValue := screen.ControlArea.Info["b"].(*int)

	color := screen.ControlArea.drawColorInputs(redValue, greenValue, blueValue)
	screen.ControlArea.Info["r"] = redValue
	screen.ControlArea.Info["g"] = greenValue
	screen.ControlArea.Info["b"] = blueValue
	screen.ControlArea.Info["color"] = color

	screen.Grid.Draw(gameContext)
	screen.StatusArea.Draw(gameContext)
}

func (grid *Grid) Draw(gameContext Game) {
	for _, aRow := range grid.GridRowByColumn {
		for _, square := range aRow {
			rl.DrawRectangleV(square.Origin, rl.NewVector2(float32(grid.Spacing), float32(grid.Spacing)), square.Color)
		}
	}
	// draw row lines
	for rowNum, rowBegin := uint8(0), grid.Origin; rowNum <= grid.NumRows; rowNum++ {
		rowEnd := rl.NewVector2(rowBegin.X+float32(int32(grid.NumColumns)*grid.Spacing), rowBegin.Y)
		rl.DrawLineEx(rowBegin, rowEnd, 1.0, grid.Color)
		rowBegin.Y += float32(grid.Spacing)
	}
	// draw column lines
	for colNum, colBegin := uint8(0), grid.Origin; colNum <= grid.NumColumns; colNum++ {
		colEnd := rl.NewVector2(colBegin.X, colBegin.Y+float32(int32(grid.NumRows)*grid.Spacing))
		rl.DrawLineEx(colBegin, colEnd, 1.0, grid.Color)
		colBegin.X += float32(grid.Spacing)
	}

}

func (statusArea *StatusArea) Draw(gameContext Game) {
	statusText := fmt.Sprintf("fade:%t, log:%t, decay:%t, fill:%t\nFPS: %.1f (%.03f), Color: %v", gameContext.FadeMode, gameContext.LogMode, gameContext.DecayMode, gameContext.FillMode, rl.GetFPS(), rl.GetFrameTime(), gameContext.PrimaryScreen.ControlArea.Info["color"])
	rl.DrawText(statusText, int32(statusArea.Origin.X+3), int32(statusArea.Origin.Y), 12, rl.Gray)
}

func (game *Game) handleInput() {
	if rl.IsKeyPressed(rl.KeyB) {
		game.FillMode = !game.FillMode
	}

	mousePos := rl.GetMousePosition()
	gridSquare, err := game.PrimaryScreen.Grid.gridCordFromMouseCord(mousePos)

	if err == nil {
		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			gridSquare.Color = rl.Blank
			game.PrimaryScreen.Grid.LastModified = time.Now()
		} else if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			if game.FillMode {
				game.PrimaryScreen.Grid.floodFill(game, *gridSquare)
			}
			color := game.PrimaryScreen.ControlArea.Info["color"]
			gridSquare.Color = color.(rl.Color) // might be redundant if we just filled it
			gridSquare.CreatedAt = time.Now()
			game.PrimaryScreen.Grid.LastModified = time.Now()
		}
	}

}

func (grid *Grid) gridCordFromMouseCord(mouseVec rl.Vector2) (*GridSquare, error) {
	// adjust x,y to use grids origin
	x := int32(mouseVec.X + grid.Origin.X)
	y := int32(mouseVec.Y + grid.Origin.Y)

	if x <= 0 || y <= 0 ||
		x >= grid.Width || y >= grid.Height {
		return nil, errors.New(fmt.Sprintf("Outside of grid bounds: (%d,%d)  (%d,%d)", x, y, grid.Width, grid.Height))
	}
	xPos := x / grid.Spacing
	yPos := y / grid.Spacing

	gs := &grid.GridRowByColumn[yPos][xPos]

	return gs, nil
}

/*
Flood-fill (node, target-color, replacement-color):
 1. If target-color is equal to replacement-color, return.
 2. If color of node is not equal to target-color, return.
 3. Set Q to the empty queue.
 4. Add node to Q.
 5. For each element N of Q:
 6.     Set w and e equal to N.
 7.     Move w to the west until the color of the node to the west of w no longer matches target-color.
 8.     Move e to the east until the color of the node to the east of e no longer matches target-color.
 9.     For each node n between w and e:
10.         Set the color of n to replacement-color.
11.         If the color of the node to the north of n is target-color, add that node to Q.
12.         If the color of the node to the south of n is target-color, add that node to Q.
13. Continue looping until Q is exhausted.
14. Return.
*/
func (grid *Grid) floodFill(game *Game, square GridSquare) int {
	squaresChanged := 0
	targetColor := square.Color
	newColor := game.PrimaryScreen.ControlArea.Info["color"].(rl.Color)

	if targetColor == newColor {
		return squaresChanged // nothing to do
	}
	var queue []GridSquare
	queue = append(queue, square)

	for i := 0; i < len(queue); i++ {

		west, east := queue[i], queue[i]
		row := grid.GridRowByColumn[west.Row]

		// Go West
		west = furthestSquare(row, west, targetColor, func(a, b uint8) uint8 { return a - b })
		east = furthestSquare(row, east, targetColor, func(a, b uint8) uint8 { return a + b })
		// fmt.Printf("furthest: %v/%v\n", west, east)

		// set nodes in between to newColor
		for wCol, eCol := west.Column, east.Column; wCol <= eCol; wCol++ {
			newSquare := row[wCol]
			newSquare.Color = newColor
			newSquare.CreatedAt = time.Now()
			row[wCol] = newSquare
			squaresChanged++

			// check to noth
			if west.Row > 0 {
				northSquare := grid.GridRowByColumn[west.Row-1][wCol]
				if northSquare.Color == targetColor {
					// fmt.Printf("Adding %v to the north\n", northSquare)
					queue = append(queue, northSquare)
				}
			}
			// check to the south
			if west.Row < uint8(grid.NumRows)-1 {
				southSquare := grid.GridRowByColumn[west.Row+1][wCol]
				if southSquare.Color == targetColor {
					// fmt.Printf("Adding %v to the south\n", southSquare)
					queue = append(queue, southSquare)
				}
			}
		}
	}
	if squaresChanged > 0 {
		grid.LastModified = time.Now()
	}
	return squaresChanged
}

func furthestSquare(row []GridSquare, startingPoint GridSquare, targetColor rl.Color, f func(a, b uint8) uint8) GridSquare {
	result := startingPoint
	for {
		squareNext := result
		squareNext.Column = f(squareNext.Column, 1)
		if int(squareNext.Column) >= len(row) || squareNext.Column < 0 {
			break
		}
		// advance to next column
		tSquare := row[squareNext.Column]
		if tSquare.Color != targetColor {
			break
		}
		result = tSquare
	}
	return result
}

// TODO: this will serialize the output
func (grid *Grid) MarshalBinary() (data []byte, err error) {
	var binBuf bytes.Buffer

	binOut := BinaryOutput{StartSentinal: StartSentinal, NumLEDs: uint16(grid.NumRows) * uint16(grid.NumColumns)}
	if err := binary.Write(&binBuf, binary.BigEndian, binOut); err != nil {
		return nil, err
	}
	for r, row := range grid.GridRowByColumn {
		for c, square := range row {
			binaryOutputSquare := BinaryOutputSquare{Row: uint8(r), Column: uint8(c), Red: square.Color.R, Blue: square.Color.B, Green: square.Color.G, Brightness: square.Color.A}

			if err := binary.Write(&binBuf, binary.BigEndian, binaryOutputSquare); err != nil {
				return nil, err
			}
		}
		fmt.Println()
	}

	return binBuf.Bytes(), nil
}
