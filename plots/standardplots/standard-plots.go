package standardplots

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"github.com/skiesel/expsys/rdb"
	"math"
	"strings"
)

func matchKeys(targetIds []string, ids []string, values []float64) ([]string, []float64) {
	for i := range targetIds {
		for j := i; j < len(ids); j++ {
			if targetIds[i] == ids[j] {
				if i != j {
					temp := ids[i]
					ids[i] = ids[j]
					ids[j] = temp
					temp2 := values[i]
					values[i] = values[j]
					values[j] = temp2
				}
				break;
			}
		}
	}

	return ids, values
}

func sortBothArrays(ids []string, values []float64) ([]string, []float64) {
	for i := range values {
		min := values[i]
		minIndex := i
		for j := i; j < len(ids); j++ {
			if values[j] < min {
				min = values[j]
				minIndex = j
			}
		}
		temp := ids[i]
		ids[i] = ids[minIndex]
		ids[minIndex] = temp
		temp2 := values[i]
		values[i] = values[minIndex]
		values[minIndex] = temp2
	}

	return ids, values
}

func PlotSolutionCosts(title string, dss []*rdb.Dataset) {
	PlotXvsY(title, dss, "final sol cost", "Solution Cost", "level", "Instance")
}

func PlotXvsY(title string, dss []*rdb.Dataset, yValuesKey, yValuesLabel , xValuesKey, xValuesLabel string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = title
	p.X.Label.Text = xValuesLabel
	p.Y.Label.Text = yValuesLabel

	var plottingArgs []interface{}

	var targetIds []string;

	for i := range dss {
		
		values, ids := dss[i].GetDatasetFloatValuesPair(yValuesKey, xValuesKey)

		if i == 0 {
			ids, values = sortBothArrays(ids, values)
			targetIds = ids
		} else {
			ids, values = matchKeys(targetIds, ids, values)
		}

		pts := make(plotter.XYs, len(values))
		for j := range pts {
			pts[j].X = float64(j)
			pts[j].Y = values[j]
		}
		plottingArgs = append(plottingArgs, dss[i].GetName())
		plottingArgs = append(plottingArgs, pts)
	}

	err = plotutil.AddLinePoints(p, plottingArgs...)
	if err != nil {
		panic(err)
	}

	plotFilename := strings.Replace(title + ".png", " ", "", -1)

	err = p.Save(4, 4, plotFilename)
	if err != nil {
		panic(err)
	}
}

func PlotSolutionCostsFactorOfBest(title string, dss []*rdb.Dataset) {
	PlotXvsFactorBestY(title, dss, "final sol cost", "Factor of Best Found Solution", "level", "Instance")
}

func PlotXvsFactorBestY(title string, dss []*rdb.Dataset, yValuesKey, yValuesLabel , xValuesKey, xValuesLabel string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = title
	p.X.Label.Text = xValuesLabel
	p.Y.Label.Text = yValuesLabel

	var plottingArgs []interface{}

	bests := make([]float64, dss[0].GetSize())

	for i := range bests {
		bests[i] = math.MaxFloat64
	}

	values := make([][]float64, len(dss))
	ids := make([][]string, len(dss))

	for i := range values {
		values[i], ids[i] = dss[i].GetDatasetFloatValuesPair(yValuesKey, xValuesKey)

		for j := range values[i] {
			if values[i][j] < bests[j] {
				bests[j] = values[i][j]
			}
		}
	}

	targetIds := make([]string, len(ids[0]))
	for i := range ids[0] {
		targetIds[i] = ids[0][i]
	}
	targetIds, bests = sortBothArrays(targetIds, bests)

	for i := range dss {

		ids[i], values[i] = matchKeys(targetIds, ids[i], values[i])

		pts := make(plotter.XYs, len(values[i]))
		for j := range pts {
			pts[j].X = float64(j)
			pts[j].Y = values[i][j] / bests[j]
		}
		plottingArgs = append(plottingArgs, dss[i].GetName())
		plottingArgs = append(plottingArgs, pts)
	}

	err = plotutil.AddLinePoints(p, plottingArgs...)
	if err != nil {
		panic(err)
	}

	plotFilename := strings.Replace(title + ".png", " ", "", -1)

	err = p.Save(4, 4, plotFilename)
	if err != nil {
		panic(err)
	}
}