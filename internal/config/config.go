package config

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/quintans/vertical-slices/internal/lib/eventbus"

	"github.com/quintans/vertical-slices/internal/features/orders"
	ordersC "github.com/quintans/vertical-slices/internal/features/orders/commands"
	ordersQ "github.com/quintans/vertical-slices/internal/features/orders/queries"
	"github.com/quintans/vertical-slices/internal/features/products"
	productC "github.com/quintans/vertical-slices/internal/features/products/commands"
	"github.com/quintans/vertical-slices/internal/features/products/eventhandlers"
	productQ "github.com/quintans/vertical-slices/internal/features/products/queries"
)

type Config struct {
	Infra
	Repositories
	ProductUsecases
	OrderUsecases
}

type Infra struct {
	EventBus *eventbus.Bus
}

type Repositories struct {
	ProductsRepo *products.Repo
	OrdersRepo   *orders.Repo
}

type ProductUsecases struct {
	CreateProduct productC.CreateProductHandler
	DeleteProduct productC.DeleteProductHandler
	GetProduct    productQ.GetProductHandler
	ListProducts  productQ.ListProductsHandler
}

type OrderUsecases struct {
	CreateOrder ordersC.CreateOrderHandler
	DeleteOrder ordersC.DeleteOrderHandler
	GetOrder    ordersQ.GetOrderHandler
	ListOrders  ordersQ.ListOrdersHandler
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

func WireProductUsecases(c *Config) {
	c.ProductUsecases = ProductUsecases{
		CreateProduct: productC.NewCreateProductHandler(c.ProductsRepo),
		DeleteProduct: productC.NewDeleteProductHandler(c.ProductsRepo),
		GetProduct:    productQ.NewGetProductHandler(c.ProductsRepo),
		ListProducts:  productQ.NewListProductsHandler(c.ProductsRepo),
	}
}

func WireOrderUsecases(c *Config) {
	c.OrderUsecases = OrderUsecases{
		CreateOrder: ordersC.NewCreateOrderHandler(c.OrdersRepo, c.ProductsRepo),
		DeleteOrder: ordersC.NewDeleteOrderHandler(c.OrdersRepo),
		GetOrder:    ordersQ.NewGetOrderHandler(c.OrdersRepo),
		ListOrders:  ordersQ.NewListOrdersHandler(c.OrdersRepo),
	}
}

func WireProductEventHandlers(c *Config) {
	eventbus.Register(c.EventBus, eventhandlers.NewOrderCreatedHandler(c.ProductsRepo))
}

func WireProductAPI(c *Config, api huma.API) {
	// Register products operations
	createProductOp, createProductHandler := productC.NewCreateProductController(c.CreateProduct)
	huma.Register(api, createProductOp, createProductHandler)

	deleteProductOp, deleteProductHandler := productC.NewDeleteProductController(c.DeleteProduct)
	huma.Register(api, deleteProductOp, deleteProductHandler)

	getProductOp, getProductHandler := productQ.NewGetProductController(c.GetProduct)
	huma.Register(api, getProductOp, getProductHandler)

	listProductsOp, listProductsHandler := productQ.NewListProductsController(c.ListProducts)
	huma.Register(api, listProductsOp, listProductsHandler)
}

func WireOrderAPI(c *Config, api huma.API) {
	createOrdOp, createOrdHandler := ordersC.NewCreateOrderController(c.CreateOrder)
	huma.Register(api, createOrdOp, createOrdHandler)

	deleteOrdOp, deleteOrdHandler := ordersC.NewDeleteOrderController(c.DeleteOrder)
	huma.Register(api, deleteOrdOp, deleteOrdHandler)

	getOrdOp, getOrdHandler := ordersQ.NewGetOrderController(c.GetOrder)
	huma.Register(api, getOrdOp, getOrdHandler)

	listOrdsOp, listOrdsHandler := ordersQ.NewListOrdersController(c.ListOrders)
	huma.Register(api, listOrdsOp, listOrdsHandler)
}
