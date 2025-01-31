// package widgets wraps core functions to be used by widgethandler
package widgets

import (
	"context"

	"github.com/mrmxf/opentsg-modules/opentsg-core/config/core"
)

// ExtractAllWidgets returns every widget used in a frame. Each widget contains its
// properties and is the
func ExtractAllWidgetsHandle(c *context.Context) map[string]core.AliasIdentityHandle {

	frameWidgets := core.GetFrameWidgetsHandle(*c)
	tagBytes := make(map[string]core.AliasIdentityHandle)
	for k, wf := range frameWidgets {
		// skip factories that don't have types
		tagBytes[k] = core.AliasIdentityHandle{FullName: k, ZPos: wf.Pos, WidgetEssentials: wf.WidgetEssentials, Contents: wf.Data}

	}

	return tagBytes
}

// MissingWidgetCheck compares all the widgets that were assigned to every widget generated for a frame, it enumerates
// a list of every missed widget and their zpos. This is so the missed ones can still be ran and not lead to blocking
// of the image generation.
func MissingWidgetCheck(c context.Context) map[core.AliasIdentity]string {
	bases := core.GetFrameWidgets(c)
	appliedWidgets := core.GetAliasMap(c)
	// update name map
	missed := make(map[core.AliasIdentity]string)
	for widgetName, content := range bases {
		if appliedWidgets.Data[widgetName] == "" { // if it wasn't assigned a tag and isn't a factory
			missed[core.AliasIdentity{FullName: widgetName, ZPos: content.Pos}] = widgetName
		}
	}

	return missed
}
