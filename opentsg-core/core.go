// package drawFunc collates the draw modules to be run in order
package core

/*
go get github.com/smartystreets/goconvey@v1.7.2
go install github.com/smartystreets/goconvey
$GOPATH/bin/goconvey

https://golangci-lint.run/usage/linters/

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1

golangci-lint --version

golangci-lint run --enable gocognit

 curl https://mrmxf.com/get/lxclog | bash
*/
/*
// Draw generates the images for each array section of the json array and applies it to the test card grid.
func Draw(configCont context.Context, imageNo int, debug bool, mnt, logType string) {
	// log here
	// logs := errhandle.LogInit(logType, mnt)
	// Sequence generation area

	// initialise the widgets here with the input file from the config  separate out the widget init section
	// widgetInit(debug, logs, &configCont)
	// potential for an array of contexts to be generated and provide all the meta data at the end

	// frames from this point on
	for frameNo := 0; frameNo < imageNo; frameNo++ {
		// new log here for each frame
		frameLog := errhandle.LogInit(logType, mnt)
		// change the log prefix for each image we generate, make a logger for each one for concurrency at a later date
		i4 := intTo4(frameNo)
		frameLog.SetPrefix(fmt.Sprintf("%v_", i4)) // update prefix to just be frame number
		// update metadata to be included in the frame context
		frameConfigCont, errs := core.FrameWidgetsGenerator(configCont, frameNo, debug)
		// this is important for showing missed widget updates
		for _, e := range errs {
			errhandle.PrintErrorMessage(frameLog, "W_CORE_opentsg_", e, true)
		}

		frameContext := widgethandler.MetaDataInit(frameConfigCont)
		errs = canvaswidget.LoopInit(frameContext)

		if len(errs) > 0 {
			// print all the errors
			for _, e := range errs {
				errhandle.PrintErrorMessage(frameLog, "F_CORE_opentsg_", e, debug)
			}

			continue // skip to the next frame number
		}

		// generate the canvas of type image.Image
		canvas, err := gridgen.GridGen(frameContext)
		if err != nil {
			errhandle.PrintErrorMessage(frameLog, "F_CORE_opentsg_", err, debug)

			continue // skip to the next frame number
		}
		fmt.Println("Generating the following file(s)", canvaswidget.GetFileName(*frameContext))
		// generate all the widgets
		widgetGen(frameContext, debug, canvas, frameLog)

		/*transformation station here

		// save the image
		savefile.CanvasSave(canvas, canvaswidget.GetFileName(*frameContext), canvaswidget.GetFileDepth(*frameContext), mnt, i4, debug, frameLog)
		fmt.Printf("\n\n")
	}
}


// each image is added to the base image
func widgetGen(c *context.Context, debug bool, canvas draw.Image, logs *log.Logger) {

	// gridgen.ParamToCanvas can be changed depending on the coordinate system
	canvasChan := make(chan draw.Image, 1)
	// put the canvas in a channel to prevent race conditions as a pointer
	// it should be called to only be added then returned to the
	canvasChan <- canvas

	var wg sync.WaitGroup  // widget waitgroup
	var wgc sync.WaitGroup // context waitgroup
	// add new widgets to the list of widgets
	widgets := []func(chan draw.Image, bool, *context.Context, *sync.WaitGroup, *sync.WaitGroup, *log.Logger){
		stripes.RampGen, zoneplate.ZoneGen, noise.NGenerator, widgethandler.MockCanvasGen,
		addimage.ImageGen, textbox.TBGenerate, bars.BarGen, saturation.SatGen,
		framecount.CountGen, qrgen.QrGen, twosi.SIGenerate, nearblack.NBGenerate,
		luma.Generate,
		// This one should be placed last as it is checking for missed names,
		// however order doesn't matter for concurrent functions with the wait groups.
		widgethandler.MockMissedGen,
	}
	wg.Add(len(widgets))
	wgc.Add(len(widgets))
	// loop through and run the widgets to glean the metadata and set their context for the run
	for _, w := range widgets {
		go w(canvasChan, debug, c, &wg, &wgc, logs)
	}

	wg.Wait()
}

// conver integer to 4 digit frames
func intTo4(num int) string {
	s := strconv.Itoa(num)
	buf0 := strings.Repeat("0", 4-len(s))
	s = buf0 + s

	return s
}
*/
