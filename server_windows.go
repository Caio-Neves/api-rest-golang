package main

import (
	"net/http"
	"rest-api-example/config"

	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
)

var logger service.Logger

type program struct {
	server          *http.Server
	serviceSettings config.ServiceSettings
}

func (p *program) Start(s service.Service) error {
	_, err := s.Logger(nil)
	if err != nil {
		log.Error("error logging to service logger", "error", err)
	}
	go p.run()
	return nil
}

func (p *program) run() {
	InitWebApplication(p.server)
}

func (p *program) Stop(s service.Service) error {
	log.Info("Stopping service...")
	return nil
}

func (p *program) SetupServerWindows(action string, args []string) error {
	svcConfig := service.Config{
		Name:        p.serviceSettings.Name,
		DisplayName: p.serviceSettings.DisplayName,
		Description: p.serviceSettings.Description,
		Arguments:   args,
	}

	svc, err := service.New(p, &svcConfig)
	if err != nil {
		return err
	}

	switch action {
	case "install":
		return svc.Install()
	case "uninstall":
		return svc.Uninstall()
	case "run":
		return svc.Run()
	case "stop":
		return svc.Stop()
	default:
		return svc.Run()
	}
}
