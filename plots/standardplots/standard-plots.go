package standardplots

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"fmt"
	"github.com/skiesel/expsys/plots"
	"github.com/skiesel/expsys/rdb"
	"os"
	"sort"
	"strconv"
	"strings"
)

func PlotSolutionCosts(title, directory string, dss []*rdb.Dataset) {
	PlotXvsY(title, directory, dss, "final sol cost", "Solution Cost", "level", "Instance")
}

func PlotSolutionCostsFactorOfBest(title, directory string, dss []*rdb.Dataset) {
	PlotXvsFactorBestY(title, directory, dss, "final sol cost", "Factor of Best Found Solution", "level", "Instance")
}

func PlotXvsFactorBestY(title, directory string, dss []*rdb.Dataset, yValuesKey, yValuesLabel, xValuesKey, xValuesLabel string) {
	newKey := yValuesKey + " (fact best)"
	newDss := rdb.AddFactorBest(dss, xValuesKey, yValuesKey, newKey)
	PlotXvsY(title, directory, newDss, newKey, yValuesLabel, xValuesKey, xValuesLabel)
}

func savePlot(title, directory string, plot *plot.Plot) {
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}

	plotFilename := strings.Replace(directory+"/"+title+".eps", " ", "", -1)

	err = plot.Save(4, 4, plotFilename)
	if err != nil {
		panic(err)
	}
}

func PlotXvsY(title, directory string, dss []*rdb.Dataset, yValuesKey, yValuesLabel, xValuesKey, xValuesLabel string) {

	var plottingArgs []interface{}

	var targetIds []string

	for i, ds := range dss {

		values, ids := ds.GetDatasetFloatValuesPair(yValuesKey, xValuesKey)

		if i == 0 {
			targetIds = ids
		} else {
			ids, values = datautils.MatchKeys(targetIds, ids, values)
		}

		pts := make(plotter.XYs, len(values))
		for j := range pts {
			pts[j].X = float64(j)
			pts[j].Y = values[j]
		}
		plottingArgs = append(plottingArgs, ds.GetName(), pts)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = title
	p.X.Label.Text = xValuesLabel
	p.Y.Label.Text = yValuesLabel

	err = plotutil.AddLinePoints(p, plottingArgs...)
	if err != nil {
		panic(err)
	}

	savePlot(title, directory, p)
}

func PlotXvsYGroupedByXs(title, directory string, groupeddss map[string][]*rdb.Dataset, yValuesKey, yValuesLabel, xValuesLabel string) {
	var plottingPointArgs []interface{}
	var plottingErrorArgs []interface{}

	names := map[string]bool{}
	xKeys := []float64{}
	for xKey, dssGroup := range groupeddss {
		xKeys = append(xKeys, datautils.ParseFloatOrFail(xKey))
		for _, d := range dssGroup {
			names[d.GetName()] = true
		}
	}

	sort.Float64s(xKeys)
	xKeyStrings := []string{}
	for _, floatVal := range xKeys {
		xKeyStrings = append(xKeyStrings, strconv.FormatFloat(floatVal, 'f', 1, 64))
	}

	for name, _ := range names {

		datasetPointsAcrossGroups := []plotter.XYer{}

		for _, groupKey := range xKeyStrings {

			dssGroup := groupeddss[groupKey]

			var ds *rdb.Dataset

			for _, d := range dssGroup {
				if d.GetName() == name {
					ds = d
					break
				}
			}

			if ds == nil {
				str := fmt.Sprintf("Couldn't find ds (%s) in dss group (%s)", name, groupKey)
				panic(str)
			}

			values := ds.GetDatasetFloatValues(yValuesKey)
			if len(values) == 0 {
				continue
			}
			xys := make(plotter.XYs, len(values))
			datasetPointsAcrossGroups = append(datasetPointsAcrossGroups, xys)

			for i, val := range values {
				xys[i].X = datautils.ParseFloatOrFail(groupKey)
				xys[i].Y = val
			}
		}

		mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, datasetPointsAcrossGroups...)

		if err != nil {
			fmt.Println(mean95)
			panic(err)
		}
		plottingPointArgs = append(plottingPointArgs, name, mean95)
		plottingErrorArgs = append(plottingErrorArgs, mean95)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = title
	p.X.Label.Text = xValuesLabel
	p.Y.Label.Text = yValuesLabel

	plotutil.AddLinePoints(p, plottingPointArgs...)
	plotutil.AddErrorBars(p, plottingErrorArgs...)

	savePlot(title, directory, p)
}
