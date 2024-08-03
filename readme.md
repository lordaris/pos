# POS

1. Install Gin with:

```bash
go get -u github.com/gin-gonic/gin
```

2. create a .env file containing a variable called MONGODB_URI, and use your mongodb uri there.
3. run `go run ./cmd/api` to start the server
4. run `sh test.sh create-roles` to create default roles (you can modify this file to create or delete default roles)
5. go to <http://localhost:8080> or

## Endpoints

| Method | URL pattern         | Handler                 | Status      |
| ------ | ------------------- | ----------------------- | ----------- |
| GET    | "/user/:id"         | getUser                 | Done        |
| GET    | "/users/:role"      | getUsersByRole          | Done        |
| POST   | "/user"             | createUser              | Done        |
| PUT    | "/user/:id"         | updateUser              | Done        |
| DELETE | "/user/:id"         | deleteUser              | Done        |
| PUT    | "/user/:id/role"    | updateUserRole          | Done        |
| GET    | "/product/:barcode" | getProduct              | Done        |
| POST   | "/product"          | createProduct           | Done        |
| PUT    | "/product/:barcode" | updateProduct           | Done        |
| DELETE | "/product/:barcode" | deleteProduct           | Done        |
| GET    | "/category"         | getCategories           | Done        |
| POST   | "/category"         | createCategory          | Done        |
| PUT    | "/category/:id"     | updateCategory          | Done        |
| DELETE | "/category/:id"     | deleteCategory          | Done        |
| GET    | "/promotion/:id"    | getPromotion            | Done        |
| POST   | "/promotion"        | productPromotion        | Done        |
| PUT    | "/promotion/:id"    | updatePromotion         | Done        |
| DELETE | "/promotion/:id"    | deletePromotion         | Done        |
| POST   | "/inventory"        | createInventoryMovement | Not started |
| PUT    | "/inventory/:id"    | updateInventoryMovement | Not started |
| DELETE | "/inventory/:id"    | deleteInventoryMovement | Not started |
| GET    | "/customer/:id"     | getCustomer             | Not started |
| POST   | "/customer"         | createCustomer          | Not started |
| PUT    | "/customer/:id"     | updateCustomer          | Not started |
| DELETE | "/customer/:id"     | deleteCustomer          | Not started |
| GET    | "/invoice/:id"      | getInvoice              | Not started |
| POST   | "/invoice"          | createInvoice           | In progress |
| PUT    | "/invoice/:id"      | updateInvoice           | Not started |
| DELETE | "/invoice/:id"      | deleteInvoice           | Not started |

## Ideas

## API Documentation

To create a promotion, the user should select one of three types: `DiscountPrice`, `DiscountPercentage`, or `BuyGet`. Even if more than one of the promotion fields are being used, the type will be the field to select which promotion to apply.

- **DiscountPercentage:** This applies a percentage discount to the price of a product. When selected, the system should calculate the discounted price and use this as the product price instead of the regular price. It receives an int as a percentage
- **DiscountPrice:** This sets a new fixed price for the product. If active, this price should be used instead of the regular price.
- **BuyGet:** This promotion has two fields: the number of items the customer needs to buy (`BuyQuantity`) and the number of items the customer gets for free (`GetQuantity`). The system should detect the quantity of products bought and adjust the total price accordingly by discounting the price of the free products.
- **StartDate** and **EndDate** should be sent as ISO 8601 (`2024-07-01t14:00:00Z`). That can be achieved using "toISOString()" method in the client.

---

test.sh have a method to create default roles. Modify it and run it to create new roles and its permissions.

## Projects and Documentation to Learn From

- [Atlas Starter Go](https://github.com/mongodb-university/atlas_starter_go/blob/master/main.go)

- Search for the "mflix tutorial"

---
