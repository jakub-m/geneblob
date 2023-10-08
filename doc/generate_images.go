package main

import (
	"fmt"
	"geneblob/graph"
	"geneblob/matrix"
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

	n := 30
	iterCount := 100
	edgeProb := 0.10

	g := graph.New(n)
	for i := range g.Vertices {
		g.Vertices[i] = graph.XY{
			X: 50 + 200*rand.Float64(),
			Y: 50 + 200*rand.Float64(),
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
	for i := 0; i < iterCount; i++ {
		fmt.Printf("Frame %d/%d\r", i+1, iterCount)
		g.UpdateForces()
		g.UpdatePoints()
		//g.SavePNG(fmt.Sprintf("tmp_%0.3d.png", i))
		frame := g.DrawImage()
		gifImage.Image = append(gifImage.Image, imageToPaletted(frame))
		gifImage.Delay = append(gifImage.Delay, 5)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return gif.EncodeAll(out, gifImage)
}

func printMatrix[T matrix.Val[T]](m *matrix.Matrix[T]) {
	for it := m.Iter(); it.HasNext(); it.Next() {
		fmt.Printf("%s %v\n", it, m.GetIt(it))
	}
}

func imageToPaletted(orig image.Image) *image.Paletted {
	paletted := image.NewPaletted(orig.Bounds(), palette.Plan9)
	for y := orig.Bounds().Min.Y; y < orig.Bounds().Max.Y; y++ {
		for x := orig.Bounds().Min.X; x < orig.Bounds().Max.X; x++ {
			paletted.Set(x, y, orig.At(x, y))
		}
	}
	return paletted
}
