# Vertical Slices Architecture
Simple application to see what a vertical slices architecture (VSA) + CQRS + DDD would look like in golang

## Description

We have products and orders.

When creating an order we need to check if there is enough stock.
On order creation, the product stock is decremented.

This project demonstrates how to communicate between the two slices, Order and Products, keeping them decoupled through the use of interfaces and events.

## Style
Although VSA advocates simplicity, but I still like to separate technology from the domain within a slice, meaning that this project style will always use a controller -> domain (commands/queries + entities) -> persistence separation.

Even if there is separation, for a specific use case, the controller and command/query are in the same file.

Regarding repositories, I think the decoupling interfaces here is beneficial to allow faster unit tests.

This is what is what the dir structure looks like:

```
internal
└── features
    ├── orders
    │   ├── commands
    │   │   ├── create_order.go
    │   │   └── delete_order.go
    │   ├── domain
    │   │   └── order.go
    │   ├── queries
    │   │   ├── get_order.go
    │   │   └── list_order.go
    │   └── repository.go
    └── products
        ├── commands
        │   ├── create_product.go
        │   └── delete_product.go
        ├── domain
        │   └── product.go
        ├── eventhandlers
        │   └── order_created.go
        ├── queries
        │   ├── get_product.go
        │   └── list_products.go
        └── repository.go
```

If a specific data model looks like it will contain logic around it, I would start implementing it as DDD aggregate from the start. Both Order and Product are aggregates. For CRUD operations just use plain struct models. 

I implemented the above order creation validation by using a DDD policy, and this can be argued to be too complex by adding unnecessary indirection, but from my perspective, that little extra complexity allows us to put every piece of logic in the same place, increasing cohesion.
