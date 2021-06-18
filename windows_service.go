// +build windows

package wireguardexporter

// Provides necessary functionality to run wireguard_exporter as a
// Windows service.

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mdlayher/wireguard_exporter/log"
	"golang.org/x/sys/windows/svc"
)

const (
	serviceName = "wireguard_exporter"
)

type wireguardExporterService struct {
	stopCh chan<- bool
	server *http.Server
}

func InitService(server *http.Server, stopCh chan<- bool) {
	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running as service: %v", err)
	}
	if isService {
		go func() {
			log.Info("Running as a service")
			err = svc.Run(serviceName, &wireguardExporterService{stopCh: stopCh, server: server})
			if err != nil {
				log.Infof("Failed to start service: %v\n", err)
			}
		}()
	}
}

func (s *wireguardExporterService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Info("Stop or Shutdown signal received")
				s.stopCh <- true
				ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
				defer func() {
					cancel()
				}()
				changes <- svc.Status{State: svc.StopPending}
				if err := s.server.Shutdown(ctx); err != nil {
					log.Fatalf("server shutdown error: %+s", err)
				}
				break loop
			default:
				log.Info(fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	return
}
