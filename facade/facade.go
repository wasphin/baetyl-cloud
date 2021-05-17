package facade

import (
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/facade/facade.go -package=facade github.com/baetyl/baetyl-cloud/v2/facade Facade

type Facade interface {
	CreateApp(ns string, baseApp, app *specV1.Application, configs []specV1.Configuration) (*specV1.Application, error)
	UpdateApp(ns string, oldApp, app *specV1.Application, configs []specV1.Configuration) (*specV1.Application, error)
	DeleteApp(ns, name string, app *specV1.Application) error

	CreateConfig(ns string, config *specV1.Configuration) (*specV1.Configuration, error)
	UpdateConfig(ns string, config *specV1.Configuration) (*specV1.Configuration, error)
	DeleteConfig(ns, name string) error
}

type facade struct {
	node      service.NodeService
	app       service.ApplicationService
	config    service.ConfigService
	index     service.IndexService
	txFactory plugin.TransactionFactory
	log       *log.Logger
}

func NewFacade(config *config.CloudConfig) (Facade, error) {
	node, err := service.NewNodeService(config)
	if err != nil {
		return nil, err
	}
	app, err := service.NewApplicationService(config)
	if err != nil {
		return nil, err
	}
	cfg, err := service.NewConfigService(config)
	if err != nil {
		return nil, err
	}
	index, err := service.NewIndexService(config)
	if err != nil {
		return nil, err
	}
	tx, err := plugin.GetPlugin(config.Plugin.Tx)
	if err != nil {
		return nil, err
	}

	return &facade {
		node:      node,
		app:       app,
		config:    cfg,
		index:     index,
		txFactory: tx.(plugin.TransactionFactory),
		log:       log.L().With(log.Any("level", "facade")),
	}, nil
}