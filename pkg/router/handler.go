package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/baileyjm02/veoir/pkg/catalogue"
	"github.com/baileyjm02/veoir/pkg/queue"
	"github.com/baileyjm02/veoir/pkg/types"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var errPathToShort = errors.New("path is too short, list to small")

// EngineHandler is the endpoint for which the Veoir engine events should be sent
// It checks if we support the sent event and handles it accordingly
// func EngineHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
// 	var image types.Image
// 	defer req.Body.Close()

// 	payload, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	image.Encoding = req.Header.Get("X-Veoir-Image-Encoding")
// 	image.Theme = req.Header.Get("X-Veoir-Image-Theme")
// 	// image.User.ID = req.Header.Get("X-Veoir-User-ID")
// 	image.Payload = payload

// 	switch encoding := image.Encoding; encoding {
// 	case "svg":
// 		queue.Queues.Publish("engine.svg", image)
// 		rw.WriteHeader(204)
// 		return
// 	default:
// 		rw.WriteHeader(501) // Return 501 Not Implemented Request as we don't support that function
// 		return
// 	}

// file, _ := json.MarshalIndent(types.Image{
//     Payload: pl,
//     Encodings: []string{"svg"},
//     Hash: "test",
// }, "", " ")

// _ = ioutil.WriteFile("test.json", file, 0644)

// }

// payload, err := ioutil.ReadAll(req.Body)
// if err != nil {
// 	logrus.Error(err)
// }

func EngineHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	var image types.Image
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Error(err)
		rw.WriteHeader(500)
		return
	}

	image.Encodings = []string{"svg", "png"}
	image.Theme = req.Header.Get("X-Veoir-Image-Theme")
	image.Payload = string(payload)
	image.Hash = catalogue.Generate()

	file, err := json.MarshalIndent(image, "", " ")
	if err != nil {
        logrus.Error(err)
		rw.WriteHeader(500)
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("public/%v.json", image.Hash), file, 0644)
	if err != nil {
		logrus.Error(err)
		rw.WriteHeader(500)
		return
	}

	queue.Queues.Publish("engine.svg", image)

	list := []string{
		fmt.Sprintf("/i/%v.svg", image.Hash),
		fmt.Sprintf("/i/%v.png", image.Hash),
	}

	json, err := json.Marshal(list)
	if err != nil {
		logrus.Error("Cannot encode to JSON ", err)
	}

	rw.Write(json)
}

func IndexHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	tmpl, err := template.ParseFiles("templates/editor.html")
	if err != nil {
		// Log the detailed error
		logrus.Error(err.Error())
		// Return a generic "Internal Server Error" message
		rw.WriteHeader(404)
		return
	}

	err = tmpl.Execute(rw, types.Image{})
	if err != nil {
		logrus.Error(err.Error())
		rw.WriteHeader(404)
		return
	}
}

func EditorHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	data, err := ioutil.ReadFile(fmt.Sprintf("public/%v.json", ps.ByName("image")))
	if err != nil {
		logrus.Error(err)
	}

	var fullImage types.Image
	err = json.Unmarshal(data, &fullImage)
	if err != nil {
		logrus.Error(err)
	}

	tmpl, err := template.ParseFiles("templates/editor.html")
	if err != nil {
		// Log the detailed error
		logrus.Error(err.Error())
		// Return a generic "Internal Server Error" message
		rw.WriteHeader(404)
		return
	}

	err = tmpl.Execute(rw, fullImage)
	if err != nil {
		logrus.Error(err.Error())
		rw.WriteHeader(404)
		return
	}
}

func ImageHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	image, err := splitPath(ps.ByName("image"))
	if err != nil && err != errPathToShort {
		logrus.Fatal(err)
	} else if err != nil {
		rw.WriteHeader(404)
		return
	}

	if strings.Contains(req.Header.Get("Accept"), "text/html") {
		http.Redirect(rw, req, fmt.Sprintf("/editor/%v", image.Hash), http.StatusTemporaryRedirect)
		return
	}
	serveImage(rw, req, image)
}

func RawHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	image, err := splitPath(ps.ByName("image"))
	if err != nil && err != errPathToShort {
		logrus.Fatal(err)
	} else if err != nil {
		rw.WriteHeader(404)
		return
	}
	serveImage(rw, req, image)
}

func splitPath(path string) (types.Image, error) {
	list := strings.Split(path, ".")
	if len(list) < 2 {
		return types.Image{}, errPathToShort
	}
	return types.Image{
		Hash:      list[0],
		Encodings: []string{list[1]},
	}, nil
}

func serveImage(rw http.ResponseWriter, req *http.Request, image types.Image) {
	http.ServeFile(rw, req, fmt.Sprintf("public/%v.%v", image.Hash, image.Encodings[0]))
}
