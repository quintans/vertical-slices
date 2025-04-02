package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/quintans/vertical-slices/internal/features/orders"
	ordersC "github.com/quintans/vertical-slices/internal/features/orders/commands"
	ordersQ "github.com/quintans/vertical-slices/internal/features/orders/queries"
	"github.com/quintans/vertical-slices/internal/features/products"
	productC "github.com/quintans/vertical-slices/internal/features/products/commands"
	productQ "github.com/quintans/vertical-slices/internal/features/products/queries"
	"github.com/quintans/vertical-slices/internal/lib/bus"
)

func main() {
	b := bus.New()

	productsRepo := products.NewRepository()
	ordersRepo := orders.NewRepository(b)

	// Configure the API routes
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	// Register your operations here
	createProductOp, createProductHandler := productC.NewCreateProductController(productC.NewCreateProductHandler(productsRepo))
	huma.Register(api, createProductOp, createProductHandler)

	deleteProductOp, deleteProductHandler := productC.NewDeleteProductController(productC.NewDeleteProductHandler(productsRepo))
	huma.Register(api, deleteProductOp, deleteProductHandler)

	getProductOp, getProductHandler := productQ.NewGetProductController(productQ.NewGetProductHandler(productsRepo))
	huma.Register(api, getProductOp, getProductHandler)

	listProductsOp, listProductsHandler := productQ.NewListProductsController(productQ.NewListProductsHandler(productsRepo))
	huma.Register(api, listProductsOp, listProductsHandler)

	createOrdOp, createOrdHandler := ordersC.NewCreateOrderController(ordersC.NewCreateOrderHandler(ordersRepo, productsRepo))
	huma.Register(api, createOrdOp, createOrdHandler)

	deleteOrdOp, deleteOrdHandler := ordersC.NewDeleteOrderController(ordersC.NewDeleteOrderHandler(ordersRepo))
	huma.Register(api, deleteOrdOp, deleteOrdHandler)

	getOrdOp, getOrdHandler := ordersQ.NewGetOrderController(ordersQ.NewGetOrderHandler(ordersRepo))
	huma.Register(api, getOrdOp, getOrdHandler)

	listOrdsOp, listOrdsHandler := ordersQ.NewListOrdersController(ordersQ.NewListOrdersHandler(ordersRepo))
	huma.Register(api, listOrdsOp, listOrdsHandler)

	// Start the server!
	http.ListenAndServe("127.0.0.1:8888", router)
}
