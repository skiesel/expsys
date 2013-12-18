package realtimeplots

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"github.com/skiesel/expsys/rdb"
	"github.com/skiesel/expsys/plots"
	"strconv"
	"fmt"
)

func groupByWfs(groupedByLookahead map[string]*rdb.Dataset, wfs []float64) (grouped map[float64]*rdb.Dataset) {
	for _, wf := range wfs {
		
		var maxLookahead int64
		maxLookahead = 0
		for lookaheadStr, ds := range groupedByLookahead {
			lookahead, err := strconv.ParseInt(lookaheadStr, 10, 0)
			if err != nil {
				fmt.Printf("groupByWfs -- failed to parse lookahead value (%s)\n", lookaheadStr)
				panic(err)
			}

			if lookahead > maxLookahead {
				stepTimes := ds.GetDatasetFloatValues("mean step cpu time")
				satisfiesWfVal := true
				
				for _, steptime := range stepTimes {
					if steptime > wf {
						satisfiesWfVal = false
						break
					}
				}

				if satisfiesWfVal {
					maxLookahead = lookahead
					grouped[wf] = ds
				}
			}
		}
	}
	return
}

func TrafficCollisionWFFiltered(dss []*rdb.Dataset) {
  p, err := plot.New()
  if err != nil {
          panic(err)
  }

	wfs := []float64{1., 0.1, 0.01, 0.001, 0.0001, 0.00001}
	groupedDss := make([]map[float64]*rdb.Dataset, len(dss))


	for i, ds := range dss {
		// Split the dataset up by the lookahead sizes
		// So you would end up with a map with entries like this
		// groupedLookahead["1000"] = dataset containing only datfiles created with a lookahead 1000
		groupedLookahead := datautils.Group(ds, "lookahead")

		// Based on the lookahead map above, pick the best lookahead
		// that still satisfies each individual wf, this will return a map like
		// groupedDss[1.] = a dataset where all instances were solved satisfying this wf (prev grouped by lookahead)
		// groupedDss[0.1] = a dataset where all instances were solved satisfying this wf
		groupedDss[i] = groupByWfs(groupedLookahead, wfs)
	}


	var plottingArgs []interface{}
	var errorBarArgs []interface{}

	for _, groupedDs := range groupedDss {		
		ptsByWfs := make([]plotter.XYer, len(wfs))
		for i, wf := range wfs {
			pts := make(plotter.XYs, groupedDs[wf].GetSize())
			ptsByWfs[i] = pts
			values := groupedDs[wf].GetDatasetIntegerValues("collisions")
	    for j := range pts {
	      pts[j].X = wf
	      pts[j].Y = float64(values[j])
	    }
		}
		mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, ptsByWfs...)
	  if err != nil {
			panic(err)
	  }

	  plottingArgs = append(plottingArgs, groupedDs[0].GetName())
	  plottingArgs = append(plottingArgs, mean95)
	  errorBarArgs = append(errorBarArgs, mean95)
	}

  plotutil.AddLinePoints(p, plottingArgs...)
  plotutil.AddErrorBars(p, errorBarArgs...)

  p.Save(4, 4, "trafficcollisions.png")
}





