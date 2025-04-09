# 🧱 Vertical Slice Architecture in Go

A simple application demonstrating **one way** to structure a Go project using **Vertical Slice Architecture (VSA)**, along with **CQRS**, **Domain-Driven Design (DDD)** and **Clean Architecture (CA)** principles.

---

## 📝 Overview

This project includes two primary domains:

- **Products**
- **Orders**

### Order Flow

When creating an order:

1. The system checks if there is enough product stock.
2. If stock is available, the order is created and the stock is decremented.

This showcases how to keep slices decoupled through **interfaces** and **domain events**, allowing clear communication between the `Orders` and `Products` slices.

---

## 🧩 Architecture & Style

Although VSA encourages simplicity, this implementation maintains a clean internal structure by **separating technology concerns from the domain** and a dependency direction as stated by CA, within each slice. The general flow looks like:

`controller → use case (commands/queries + aggregate) → gateway (persistence)`


For individual use cases, the controller and its associated command/query may be placed in the same file for convenience and clarity.

### ✅ Testing & Repositories

Repositories are abstracted behind interfaces to support **fast and isolated unit testing** and promote loose coupling between layers.

---

## 📁 Project Structure

Here's how the internal directory is organized:

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


---

## 🧠 Domain Modeling & DDD

If a data model involves business logic, I would model it as a **DDD aggregate** from the start. In this app, both `Order` and `Product` are aggregates.

For straightforward CRUD operations, I would just use simple struct-based models.

As business logic increases always consider to refactor to an aggregate.

---

## 🔄 Domain Policies

Order creation involves validating available stock, implemented using a **DDD policy**. While this may introduce a bit more abstraction, it allows logic to remain **highly cohesive** and centralized—making it easier to evolve and maintain over time. 

---

## 🚀 Goals

- Showcase how VSA can be applied in Go.
- Emphasize decoupling, testability, and domain-focused design.
- Provide a solid base for extending into a more complete microservice or monolith.

---

Feel free to fork or adapt this for your own projects!
