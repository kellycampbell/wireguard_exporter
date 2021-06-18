// +build !windows

package wireguardexporter

import (
	"net/http"
)

// Non-windows service init placeholder function

func InitService(server *http.Server, stopChan chan<- bool) {

}
