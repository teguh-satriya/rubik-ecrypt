package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

type Warna struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

type Block struct {
	X   int
	Y   int
	Col []Warna
}

type Koordinat struct {
	X int
	Y int
}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	imgfile, err := os.Open("tes.jpg")

	if err != nil {
		fmt.Println("img.jpg file not found!")
		os.Exit(1)
	}

	defer imgfile.Close()

	imgCfg, _, err := image.DecodeConfig(imgfile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)

	imgfile.Seek(0, 0)
	img, _, err := image.Decode(imgfile)

	ActImg := make([][]Warna, 1000)
	for y := 0; y < 1000; y++ {
		ActImg[y] = make([]Warna, 1000)
		for x := 0; x < 1000; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			ActImg[y][x].R = r
			ActImg[y][x].G = g
			ActImg[y][x].B = b
			ActImg[y][x].A = a
		}
	}

	var iterasi = 250
	BlockAct := []Block{}
	BlockEnc := []Block{}
	for y := 0; y < 1000; y += iterasi {
		yArr := ActImg[y : y+iterasi]
		for x := 0; x < 1000; x += iterasi {
			var arrPos = []Warna{}
			for z := 0; z < 250; z++ {
				arr := yArr[z][x : x+iterasi]
				arrPos = append(arrPos, arr...)
			}

			BlockAct = append(BlockAct, Block{x, y, arrPos})
			BlockEnc = append(BlockEnc, Block{x, y, RotateClockWise(arrPos, 500)})
		}
	}
	actual := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	enc := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	for _, b := range BlockAct {
		Koord := []Koordinat{}
		for y := b.Y; y < (b.Y + iterasi); y++ {
			for x := b.X; x < (b.X + iterasi); x++ {
				Koord = append(Koord, Koordinat{x, y})
			}
		}

		for k, c := range b.Col {
			actual.Set(Koord[k].X, Koord[k].Y, color.RGBA{
				uint8(c.R >> 8), uint8(c.G >> 8), uint8(c.B >> 8), uint8(c.A >> 8)})
		}
	}

	for _, b := range BlockEnc {
		Koord := []Koordinat{}
		for y := b.Y; y < (b.Y + iterasi); y++ {
			for x := b.X; x < (b.X + iterasi); x++ {
				Koord = append(Koord, Koordinat{x, y})
			}
		}

		for k, c := range b.Col {
			enc.Set(Koord[k].X, Koord[k].Y, color.RGBA{
				uint8(c.R >> 8), uint8(c.G >> 8), uint8(c.B >> 8), uint8(c.A >> 8)})
		}
	}

	a, err := os.OpenFile("Actual.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer a.Close()
	png.Encode(a, actual)

	e, err := os.OpenFile("Enc.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer e.Close()
	png.Encode(e, enc)

}

func RotateClockWise(matrix []Warna, n int) []Warna {
	TempArr := make([]Warna, n)

	TempArr = make([]Warna, n)
	for x := 0; x < n; x++ {
		TempArr[x] = matrix[(n - x - 1)]
	}

	return TempArr
}
