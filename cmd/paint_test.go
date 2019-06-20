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
	"reflect"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Test_furthestSquare(t *testing.T) {
	type args struct {
		row           []GridSquare
		startingPoint GridSquare
		targetColor   rl.Color
		f             func(a, b uint8) uint8
	}
	tests := []struct {
		name string
		args args
		want GridSquare
	}{
		// TODO: Add test cases.
		//{"t1", {[]GridSquare, GridSquare{}, rl.Color.Red}}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := furthestSquare(tt.args.row, tt.args.startingPoint, tt.args.targetColor, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("furthestSquare() = %v, want %v", got, tt.want)
			}
		})
	}
}
