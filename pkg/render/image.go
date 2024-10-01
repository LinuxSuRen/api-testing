/*
Copyright 2024 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language 24 permissions and
limitations under the License.
*/
package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"time"

	"github.com/linuxsuren/api-testing/pkg/util"
)

func generateRandomImage(width, height int) (data string, err error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	blockSize := int(math.Max(float64(width), float64(height)) / 4)
	rand.Seed(time.Now().UnixNano())

	for y := 0; y < height; y += blockSize {
		for x := 0; x < width; x += blockSize {
			r := uint8(rand.Intn(255))
			g := uint8(rand.Intn(255))
			b := uint8(rand.Intn(255))
			col := color.RGBA{R: r, G: g, B: b, A: 255}

			for iy := y; iy < y+blockSize && iy < height; iy++ {
				for ix := x; ix < x+blockSize && ix < width; ix++ {
					img.Set(ix, iy, col)
				}
			}
		}
	}

	buf := new(bytes.Buffer)
	if err = png.Encode(buf, img); err == nil {
		data = fmt.Sprintf("%s%s", util.ImageBase64Prefix, base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
	return
}
