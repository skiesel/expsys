package realtimeplots

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"github.com/skiesel/expsys/rdb"
	"strconv"
	"fmt"
)

func group(ds *rdb.Dataset, key string) (grouped map[string]*rdb.Dataset) {
	values := ds.GetDatasetStringValues(key)
	valueSet := make(map[string]string, 0)
	for _, value := range values {
		_, exists := valueSet[value]
		if !exists {
			valueSet[value] = value
		}
	}

	for value, _ := range valueSet {
		filter := func(str string)bool {
			return str == value
		}
		grouped[value] = ds.FilterDataset(filter, key)
	}
	return
}

func groupByWfs(ds *rdb.Dataset, wfs []float64) (grouped map[float64]*rdb.Dataset) {
	groupedLookahead := group(ds, "lookahead")
	for _, wf := range wfs {
		
		var maxLookahead int64
		maxLookahead = 0
		for lookaheadStr, ds := range groupedLookahead {
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
		groupedDss[i] = groupByWfs(ds, wfs)
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





