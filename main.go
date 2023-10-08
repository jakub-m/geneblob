package main

import (
	"fmt"
	"geneblob/graph"
	"image"
	"image/color/palette"
	"image/gif"
	"log"
	"math/rand"
	"os"
)

// The command generates all the images for the documentation.
func main() {
	// Underscore in the image name informs that the file is generated.
	if err := generateAnimatedRandom("doc/_random_graph_forces.gif"); err != nil {
		log.Fatal(err)
	}
	log.Print("Done")
}

func generateAnimatedRandom(path string) error {
	log.Println("Generate", path)
	rand.Seed(0)

	s := 10
	dd := float64(200) / float64(s)
	iterCount := 100
	edgeProb := 0.03

	g := graph.New(s * s)

	//for i := range g.Vertices {
	//	g.Vertices[i] = graph.XY{
	//		X: 50 + 200*rand.Float64(),
	//		Y: 50 + 200*rand.Float64(),
	//	}
	//}
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			g.Vertices[s*j+i] = graph.XY{
				X: 50 + float64(i)*dd,
				Y: 50 + float64(j)*dd,
			}
		}
	}

	for i := range g.Vertices {
		for k := range g.Vertices {
			if i < k {
				if rand.Float64() < edgeProb {
					g.Edges.SetSym(i, k, true)
				}
			}
		}
	}

	gifImage := &gif.GIF{
		Image: []*image.Paletted{},
	}

	iFrame := 0
	saveFrame := func() {
		fmt.Printf("Frame %d/%d\r", iFrame, iterCount)
		pngName := fmt.Sprintf("tmp_%0.3d.png", iFrame)
		log.Println("save", pngName)
		g.SavePNG(pngName)
		frame := g.DrawImage()
		gifImage.Image = append(gifImage.Image, imageToPaletted(frame))
		gifImage.Delay = append(gifImage.Delay, 5)
		iFrame++
	}

	saveFrame()
	for i := 0; i < iterCount; i++ {
		g.UpdateForces()
		g.UpdatePoints()
		saveFrame()
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return gif.EncodeAll(out, gifImage)
}

//func printMatrix[T matrix.Val[T]](m *matrix.Matrix[T]) {
//	for it := m.Iter(); it.HasNext(); it.Next() {
//		fmt.Printf("%s %v\n", it, m.GetIt(it))
//	}
//}

func imageToPaletted(orig image.Image) *image.Paletted {
	paletted := image.NewPaletted(orig.Bounds(), palette.Plan9)
	for y := orig.Bounds().Min.Y; y < orig.Bounds().Max.Y; y++ {
		for x := orig.Bounds().Min.X; x < orig.Bounds().Max.X; x++ {
			paletted.Set(x, y, orig.At(x, y))
		}
	}
	return paletted
}
