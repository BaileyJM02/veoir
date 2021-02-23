package router

import (
	"github.com/julienschmidt/httprouter"
)

// GetRouter returns a list of routes avaliable and being listened too on the server
func GetRouter() *httprouter.Router {
	// Add the endpoint for engine payloads
    router := httprouter.New()

    // Frontend routes
    router.GET("/", IndexHandler)
    router.GET("/i/:image", ImageHandler)
    router.GET("/editor/:image/", EditorHandler)

	// Raw image endpoint
    router.GET("/raw/:image", RawHandler)
    
	// API endpoint
	router.POST("/api/v1/create", EngineHandler)


	return router
}
