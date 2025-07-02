package constants

const (
	RouteIndex = "/"
	RouteStart = "/start"
	RouteGame  = "/game"
	RouteDice  = "/dice/:value"
	RouteMove  = "/move/:x/:y"

	// CustomDiceDir is checked for user-provided dice images. Files should
	// be named using the format "dice-<value>.svg".
	CustomDiceDir = "static/custom"
	ServerPort    = ":8080"
)
