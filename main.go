package main

import (
    "github.com/baileyjm02/veoir/pkg/logger"
    "github.com/baileyjm02/veoir/pkg/router"
    "github.com/baileyjm02/veoir/pkg/engine"
)

func init() {
    // Format the logger so we can import logrus and have the expected output in all areas
    logger.Format()
}

func main() {
    // Start the image queues
    go engine.StartSVGQueue()
    go engine.StartPNGQueue()

    // Start the router package
    router.Start()
}
