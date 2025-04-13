package config

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/quintans/vertical-slices/internal/lib/eventbus"

	"github.com/quintans/vertical-slices/internal/features/orders"
	ordCmd "github.com/quintans/vertical-slices/internal/features/orders/commands"
	ordQry "github.com/quintans/vertical-slices/internal/features/orders/queries"
	"github.com/quintans/vertical-slices/internal/features/products"
	prdCmd "github.com/quintans/vertical-slices/internal/features/products/commands"
	"github.com/quintans/vertical-slices/internal/features/products/eventhandlers"
	prdQry "github.com/quintans/vertical-slices/internal/features/products/queries"
)

type Config struct {
	Infra
	Repositories
}

type Infra struct {
	EventBus *eventbus.Bus
}

type Repositories struct {
	ProductsRepo *products.Repo
	OrdersRepo   *orders.Repo
}

func WireInfra(c *Config) {
	eb := eventbus.New()
	c.Infra = Infra{
		EventBus: eb,
	}
}

func WireRepositories(c *Config) {
	c.Repositories = Repositories{
		ProductsRepo: products.NewRepository(),
		OrdersRepo:   orders.NewRepository(c.EventBus),
	}
}

func WireProductEventHandlers(c *Config) {
	eventbus.Register(c.EventBus, eventhandlers.NewOrderCreatedHandler(c.ProductsRepo))
}

func WireProductAPI(c *Config, api huma.API) {
	prdCmd.RegisterCreateProductController(api, c.ProductsRepo)
	prdCmd.RegisterDeleteProductController(api, c.ProductsRepo)
	prdQry.RegisterGetProductController(api, c.ProductsRepo)
	prdQry.RegisterListProductsController(api, c.ProductsRepo)
}

func WireOrderAPI(c *Config, api huma.API) {
	ordCmd.RegisterCreateOrderController(api, c.OrdersRepo, c.ProductsRepo)
	ordCmd.RegisterDeleteOrderController(api, c.OrdersRepo)
	ordQry.RegisterGetOrderController(api, c.OrdersRepo)
	ordQry.RegisterListOrdersController(api, c.OrdersRepo)
}
