package standardplots

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"github.com/skiesel/expsys/rdb"
	"github.com/skiesel/expsys/plots"
	"math"
	"strings"
)

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

	for i, ds := range dss {
		
		values, ids := ds.GetDatasetFloatValuesPair(yValuesKey, xValuesKey)

		if i == 0 {
			ids, values = datautils.SortBothArrays(ids, values)
			targetIds = ids
		} else {
			ids, values = datautils.MatchKeys(targetIds, ids, values)
		}

		pts := make(plotter.XYs, len(values))
		for j := range pts {
			pts[j].X = float64(j)
			pts[j].Y = values[j]
		}
		plottingArgs = append(plottingArgs, ds.GetName())
		plottingArgs = append(plottingArgs, pts)
	}

	err = plotutil.AddLinePoints(p, plottingArgs...)
	if err != nil {
		panic(err)
	}

	plotFilename := strings.Replace(title + ".eps", " ", "", -1)

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
	p.Legend.Top = true
	
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
	targetIds, bests = datautils.SortBothArrays(targetIds, bests)

	for i := range dss {

		ids[i], values[i] = datautils.MatchKeys(targetIds, ids[i], values[i])

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

	plotFilename := strings.Replace(title + ".eps", " ", "", -1)

	err = p.Save(4, 4, plotFilename)
	if err != nil {
		panic(err)
	}
}