package router

import (
	"github.com/julienschmidt/httprouter"
)

// GetRouter returns a list of routes avaliable and being listened too on the server
func GetRouter() *httprouter.Router {
	// Add the endpoint for engine payloads
    router := httprouter.New()

    // Frontend
    router.GET("/", IndexHandler)
    router.GET("/i/:image", ImageHandler)
    router.GET("/editor/:image/", EditorHandler)

	// Raw image
    router.GET("/raw/:image", RawHandler)
    
	// API
	router.POST("/api/v1/create", EngineHandler)


	return router
}
