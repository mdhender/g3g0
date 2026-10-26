// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package genesis

import (
	"bufio"
	"os"
)

// This file ports lib/map_image.c, which renders the --image debugging BMP.

// rgb is an RGB triple, written to the BMP as B, G, R.
type rgb struct {
	r, g, b byte
}

// glyphColor maps a map glyph to its color.
func glyphColor(c byte) rgb {
	switch c {
	// Ocean; each sea-lane variant shares its plain glyph's color.
	case ',', ';':
		return rgb{0, 0, 160}
	case '.', ':':
		return rgb{0, 80, 220}
	case ' ', '~':
		return rgb{40, 140, 255}
	case '\'', '"':
		return rgb{120, 200, 255}

	// Land and features.
	case 'p':
		return rgb{96, 176, 80}
	case 'f':
		return rgb{32, 96, 32}
	case 'm', 'M':
		return rgb{128, 128, 128}
	case 'd':
		return rgb{210, 180, 140}
	case 's':
		return rgb{107, 142, 35}
	case 'O':
		return rgb{255, 0, 0}
	case 'o':
		return rgb{200, 200, 120}
	case '0':
		return rgb{255, 215, 0}

	// Glyphs only mapgen reads.
	case 'P':
		return rgb{96, 176, 80}
	case 'D':
		return rgb{210, 180, 140}
	case 'S':
		return rgb{107, 142, 35}
	case 'F':
		return rgb{32, 96, 32}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return rgb{255, 215, 0}
	case '^', 'v', '{', '}':
		return rgb{128, 128, 128}
	case ']', '[':
		return rgb{107, 142, 35}
	case '*', '%':
		return rgb{96, 176, 80}
	}
	return rgb{255, 0, 255}
}

// writeBMP is map_image_write_bmp: it writes the width x width glyph grid
// (row-major) to path as an uncompressed 24-bit BMP. It reports false where
// the C function returns nonzero.
func writeBMP(path string, grid []byte, width, pixels, gridN int) bool {
	if width < 1 || pixels < 1 || pixels > 64 || gridN < 0 || gridN > 4 {
		return false
	}
	if int64(width)*int64(pixels)+int64(width+1)*int64(gridN) > 1000000 {
		return false
	}
	dim := width*pixels + (width+1)*gridN

	rowBytes := uint32(dim) * 3
	pad := (4 - (rowBytes & 3)) & 3
	padded := rowBytes + pad
	pixelBytes := padded * uint32(dim)
	fileSize := 14 + 40 + pixelBytes

	fp, ok := fopenW(path)
	if !ok {
		return false
	}
	w := bufio.NewWriter(fp)
	putU16 := func(v uint16) {
		w.WriteByte(byte(v))
		w.WriteByte(byte(v >> 8))
	}
	putU32 := func(v uint32) {
		w.WriteByte(byte(v))
		w.WriteByte(byte(v >> 8))
		w.WriteByte(byte(v >> 16))
		w.WriteByte(byte(v >> 24))
	}

	// BITMAPFILEHEADER
	w.WriteString("BM")
	putU32(fileSize)
	putU16(0)
	putU16(0)
	putU32(14 + 40)

	// BITMAPINFOHEADER
	putU32(40)
	putU32(uint32(dim))
	putU32(uint32(dim))
	putU16(1)
	putU16(24)
	putU32(0)
	putU32(pixelBytes)
	putU32(0)
	putU32(0)
	putU32(0)
	putU32(0)

	// Pixel rows run bottom to top.
	for py := dim - 1; py >= 0; py-- {
		for px := 0; px < dim; px++ {
			cellX, cellY := -1, -1

			t := px % (pixels + gridN)
			if px < width*pixels+(width+1)*gridN && t >= gridN && px/(pixels+gridN) < width {
				cellX = px / (pixels + gridN)
			}
			t = py % (pixels + gridN)
			if py < width*pixels+(width+1)*gridN && t >= gridN && py/(pixels+gridN) < width {
				cellY = py / (pixels + gridN)
			}

			col := rgb{0, 0, 0} // grid lines are black
			if cellX >= 0 && cellY >= 0 {
				col = glyphColor(grid[cellY*width+cellX])
			}
			w.WriteByte(col.b)
			w.WriteByte(col.g)
			w.WriteByte(col.r)
		}
		for range pad {
			w.WriteByte(0)
		}
	}

	if w.Flush() != nil {
		fp.Close()
		return false
	}
	return fp.Close() == nil
}

// fopenW opens path the way fopen(path, "w") does.
func fopenW(path string) (*os.File, bool) {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	return fp, err == nil
}
