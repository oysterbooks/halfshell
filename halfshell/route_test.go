package halfshell

import (
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {

	basicRoute := RouteForName("basic")

	if basicRoute == nil {
		t.Fatal("basicRoute should not be nil")
	}

	req, _ := http.NewRequest("GET", "http://example.com/basic/image.jpg?w=25&h=10", nil)

	if !basicRoute.ShouldHandleRequest(req) {
		t.Fail()
	}

	basicSourceOptions, basicProcessorOptions := basicRoute.SourceAndProcessorOptionsForRequest(req)

	if "/image.jpg" != basicSourceOptions.Path {
		t.Fail()
	}

	if 25 != basicProcessorOptions.Dimensions.Width {
		t.Fail()
	}

	if 10 != basicProcessorOptions.Dimensions.Height {
		t.Fail()
	}
}

func TestRouteWithDimensionsInFilename(t *testing.T) {

	route := RouteForName("dimensions-embedded-in-filename")

	if route == nil {
		t.Fatal("route should not be nil")
	}

	reqForOriginalImage, _ := http.NewRequest("GET", "http://example.com/complex/image.jpg", nil)

	if !route.ShouldHandleRequest(reqForOriginalImage) {
		t.Error("Route should handle request for original image")
	}

	req, _ := http.NewRequest("GET", "http://example.com/complex/image.100x50.jpg", nil)

	if !route.ShouldHandleRequest(req) {
		t.Error("Route should handle rqeuest for resized image")
	}

	complexSourceOptions, complexProcessorOptions := route.SourceAndProcessorOptionsForRequest(req)

	if "/image.jpg" != complexSourceOptions.Path {
		t.Fail()
	}

	if 100 != complexProcessorOptions.Dimensions.Width {
		t.Fail()
	}

	if 50 != complexProcessorOptions.Dimensions.Height {
		t.Fail()
	}

}

func TestRouteWithDimensionsInPath(t *testing.T) {

	route := RouteForName("dimensions-in-path")

	if route == nil {
		t.Fatal("route should not be nil")
	}

	reqForOriginalImage, _ := http.NewRequest("GET", "http://example.com/complex/image.jpg", nil)

	if !route.ShouldHandleRequest(reqForOriginalImage) {
		t.Error("Route should handle request for original image")
	}

	req, _ := http.NewRequest("GET", "http://example.com/complex/100x50/image.jpg", nil)

	if !route.ShouldHandleRequest(req) {
		t.Error("Route should handle rqeuest for resized image")
	}

	complexSourceOptions, complexProcessorOptions := route.SourceAndProcessorOptionsForRequest(req)

	if "/image.jpg" != complexSourceOptions.Path {
		t.Fail()
	}

	if 100 != complexProcessorOptions.Dimensions.Width {
		t.Fail()
	}

	if 50 != complexProcessorOptions.Dimensions.Height {
		t.Fail()
	}

}

func RouteForName(name string) *Route {
	config := NewConfigFromFile("testdata/config.json")

	var targetRoute *Route

	for _, routeConfig := range config.RouteConfigs {

		route := NewRouteWithConfig(routeConfig, config.StatterConfig)

		if name == route.Name {
			targetRoute = route
		}
	}

	return targetRoute
}
