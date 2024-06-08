package marching_square

// import (
// 	"log"
// 	"sync"

// 	"gonum.org/v1/plot"
// 	"gonum.org/v1/plot/plotter"
// 	"gonum.org/v1/plot/vg"
// )

// // Grid represents the 2D scalar field
// type Grid struct {
// 	data [][]float64
// }

// // ContourPoint represents a point in the contour
// type ContourPoint struct {
// 	x, y, z float64
// }

// // MarchingSquares generates contours for the given grid
// func MarchingSquares(grid Grid, threshold float64, wg *sync.WaitGroup, ch chan ContourPoint) {
// 	defer wg.Done()

// 	for i := 0; i < len(grid.data)-1; i++ {
// 		for j := 0; j < len(grid.data[i])-1; j++ {
// 			// Get the four corner values of the cell
// 			a := grid.data[i][j]
// 			b := grid.data[i+1][j]
// 			c := grid.data[i+1][j+1]
// 			d := grid.data[i][j+1]

// 			// Determine the index based on the threshold
// 			index := 0
// 			if a > threshold {
// 				index |= 1
// 			}
// 			if b > threshold {
// 				index |= 2
// 			}
// 			if c > threshold {
// 				index |= 4
// 			}
// 			if d > threshold {
// 				index |= 8
// 			}

// 			// Determine the contour points for the cell
// 			switch index {
// 			case 1, 14:
// 				ch <- ContourPoint{float64(i), float64(j), interpolate(a, b, threshold)}
// 			case 2, 13:
// 				ch <- ContourPoint{float64(i + 1), float64(j), interpolate(b, c, threshold)}
// 			case 3, 12:
// 				ch <- ContourPoint{float64(i), float64(j), interpolate(a, b, threshold)}
// 				ch <- ContourPoint{float64(i + 1), float64(j), interpolate(b, c, threshold)}
// 			case 4, 11:
// 				ch <- ContourPoint{float64(i + 1), float64(j), interpolate(c, d, threshold)}
// 			case 5, 10:
// 				ch <- ContourPoint{float64(i), float64(j), interpolate(a, d, threshold)}
// 				ch <- ContourPoint{float64(i + 1), float64(j), interpolate(c, d, threshold)}
// 			case 6, 9:
// 				ch <- ContourPoint{float64(i + 1), float64(j), interpolate(b, c, threshold)}
// 				ch <- ContourPoint{float64(i + 1), float64(j), interpolate(c, d, threshold)}
// 			case 7, 8:
// 				ch <- ContourPoint{float64(i), float64(j), interpolate(a, b, threshold)}
// 			}
// 		}
// 	}
// }

// // interpolate linearly interpolates between two values
// func interpolate(a, b, threshold float64) float64 {
// 	return a + (b-a)*(threshold-a)/(b-a)
// }

// func main() {
// 	// Example usage
// 	grid := Grid{
// 		data: [][]float64{
// 			{1, 1, 1, 1},
// 			{1, 2, 2, 1},
// 			{1, 2, 2, 1},
// 			{1, 2, 2, 1},
// 		},
// 	}

// 	threshold := 1.3

// 	var wg sync.WaitGroup
// 	ch := make(chan ContourPoint)

// 	// Number of goroutines can be adjusted based on the grid size
// 	numGoroutines := 4

// 	wg.Add(numGoroutines)

// 	// Start goroutines
// 	for i := 0; i < numGoroutines; i++ {
// 		go MarchingSquares(grid, threshold, &wg, ch)
// 	}

// 	// Close channel when all goroutines are done
// 	go func() {
// 		wg.Wait()
// 		close(ch)
// 	}()

// 	// Collect and store contour points
// 	var contourPoints []ContourPoint
// 	for point := range ch {
// 		contourPoints = append(contourPoints, point)
// 	}

// 	// Plot the contour points using gonum/plot
// 	err := plotContour(contourPoints)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// // plotContour creates a scatter plot of contour points
// func plotContour(contourPoints []ContourPoint) error {
// 	p := plot.New()

// 	points := make(plotter.XYs, len(contourPoints))
// 	for i, point := range contourPoints {
// 		points[i].X = point.x
// 		points[i].Y = point.y
// 	}

// 	s, err := plotter.NewScatter(points)
// 	if err != nil {
// 		return err
// 	}
// 	p.Add(s)

// 	// Set labels for axes
// 	p.X.Label.Text = "X"
// 	p.Y.Label.Text = "Y"

// 	if err := p.Save(4*vg.Inch, 4*vg.Inch, "contour_plot.png"); err != nil {
// 		return err
// 	}

//		return nil
//	}
