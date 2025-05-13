package opentsgwidgets

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/ebu3373/bars"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/ebu3373/luma"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/ebu3373/nearblack"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/ebu3373/saturation"
	"github.com/mrmxf/opentsg-modules/opentsg-widgets/ebu3373/twosi"
)

func TestHandlers(t *testing.T) {

	// run the current handler methods
	otsgh, err := tsg.BuildOpenTSG("./testdata/handlerLoader.json", "", true, &tsg.RunnerConfiguration{RunnerCount: 6, ProfilerEnabled: true})
	fmt.Println(err)
	//	otsgh.HandleFunc("builtin.canvas", func(r1 tsg.Response, r2 *tsg.Request) { fmt.Println("ring a ding") })
	otsgh.Handle(bars.WidgetType, bars.Schema, bars.BarJSON{})
	otsgh.Handle(luma.WidgetType, luma.Schema, luma.LumaJSON{})
	otsgh.Handle(nearblack.WidgetType, nearblack.Schema, nearblack.Config{})
	otsgh.Handle(saturation.WidgetType, saturation.Schema, saturation.Config{})
	otsgh.Handle(twosi.WidgetType, twosi.Schema, twosi.Config{})

	tsg.LogToFile(otsgh, slog.HandlerOptions{Level: slog.LevelDebug}, &tsg.LogOptions{Folder: "./testdata/", JobID: "handler"})
	tsg.AddBaseEncoders(otsgh)
	otsgh.Run("")
}
