# POS

1. create a .env file containing a variable called MONGODB_URI, and use your mongodb uri there.
2. run `go run ./cmd/api` to start the server
3. run `sh test.sh create-roles` to create default roles (you can modify this file to create or delete default roles)
4. go to <http://localhost:8080> or

## Endpoints

| Method | URL pattern         | Handler                 | Status      |
| ------ | ------------------- | ----------------------- | ----------- |
| GET    | "/user/:id"         | getUser                 | In progress |
| GET    | "/users/:role"      | getUsersByRole          | In progress |
| POST   | "/user"             | createUser              | In progress |
| PUT    | "/user/:id"         | updateUser              | In progress |
| DELETE | "/user/:id"         | deleteUser              | Done        |
| PUT    | "/user/:id/role"    | updateUserRole          | In progress |
| GET    | "/product/:barcode" | getProduct              | Done        |
| POST   | "/product"          | createProduct           | In progress |
| PUT    | "/product/:barcode" | updateProduct           | Not started |
| DELETE | "/product/:barcode" | deleteProduct           | Done        |
| GET    | "/category"         | getCategory             | Not started |
| POST   | "/category"         | createCategory          | Not started |
| PUT    | "/category/:id"     | updateCategory          | Not started |
| DELETE | "/category/:id"     | deleteCategory          | Not started |
| GET    | "/promotion/:id"    | getPromotion            | Not started |
| POST   | "/promotion"        | productPromotion        | Done        |
| PUT    | "/promotion/:id"    | updatePromotion         | Not started |
| DELETE | "/promotion/:id"    | deletePromotion         | Not started |
| GET    | "/inventory/:id"    | getInventoryMovement    | Not started |
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

TODO: Update this when the middleware is created
The user update doesn't ask for the old password to modify the user nor its password, as it's intended to be used for administrators and not the actual user itself. It should be protected via middleware checking for the role permissions.

---

test.sh have a method to create default roles. Modify it and run it to create new roles and its permissions.

## Projects and Documentation to Learn From

- [Atlas Starter Go](https://github.com/mongodb-university/atlas_starter_go/blob/master/main.go)

- Search for the "mflix tutorial"

---
