// Package qrgen generates a qr code based on user string and places it on the graph, this is the last item to be added
package qrgen

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/mrmxf/opentsg-modules/opentsg-core/colour"
	"github.com/mrmxf/opentsg-modules/opentsg-core/tsg"
)

const (
	WidgetType = "builtin.qrcode"
)

func (q Config) Handle(resp tsg.Response, req *tsg.Request) {

	message := q.Code
	if message == "" {
		// Return but don't fill up the stdout with errors
		resp.Write(tsg.WidgetWarning, "no qr code message given")
		return
	}
	/*
		@ TODO: utilise this information for metadata in the barcode
			if qrC.Query != nil {
				// Do some more metadata extraction
				for _, q := range *qrC.Query {
					fmt.Println(q)
					fmt.Println(extract(opt[0].(*context.Context), q.Target, q.Keys...))
				}
			}
	*/

	code, err := qr.Encode(message, qr.H, qr.Auto)
	if err != nil {
		resp.Write(tsg.WidgetError, err.Error())
		return
	}

	b := resp.BaseImage().Bounds().Max
	if q.Size != nil {
		width, height := q.Size.Width, q.Size.Height
		if width != 0 && height != 0 {
			w, h := (width/100)*float64(b.X), (height/100)*float64(b.Y)
			code, err = barcode.Scale(code, int(w), int(h))
			if err != nil {
				resp.Write(tsg.WidgetError, err.Error())
				return
			}
		}
	}

	offset, err := q.CalcOffset(b)

	if err != nil {
		resp.Write(tsg.WidgetError, fmt.Sprintf("0DEV error finding the offset :%v", err))
		return
	}

	if offset.X > (b.X - code.Bounds().Max.X) {
		resp.Write(tsg.WidgetError, fmt.Sprintf("0133 the x position %v is greater than the x boundary of %v", offset.X, resp.BaseImage().Bounds().Max.X))
		return
	} else if offset.Y > b.Y-code.Bounds().Max.Y {
		resp.Write(tsg.WidgetError, fmt.Sprintf("0133 the y position %v is greater than the y boundary of %v", offset.Y, resp.BaseImage().Bounds().Max.Y))
		return
	}
	// draw qr code as a mid point, or make colour space agnostic
	colour.Draw(resp.BaseImage(), resp.BaseImage().Bounds().Add(offset), code, image.Point{}, draw.Over)

	resp.Write(tsg.WidgetSuccess, "success")

}
