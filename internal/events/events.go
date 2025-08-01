package events

import "github.com/catouberos/geoloc/base"

func RegisterEvents(app *base.App) {
	registerGeolocationInsertHandler(app)
}
