# POS

## Endpoints

| Method | URL pattern      | Handler                 | status      |
| ------ | ---------------- | ----------------------- | ----------- |
| GET    | "/user/:id"      | getUser                 | not started |
| POST   | "/user"          | createUser              | in progress |
| PUT    | "/user/:id"      | updateUser              | in progress |
| DELETE | "/user/:id"      | deleteUser              | not started |
| PUT    | "/user/:id/role" | updateUserRole          | in progress |
| GET    | "/product/:id"   | getProduct              | not started |
| POST   | "/product"       | createProduct           | in progress |
| PUT    | "/product/:id"   | updateProduct           | not started |
| DELETE | "/product/:id"   | deleteProduct           | not started |
| GET    | "/category"      | getCategory             | not started |
| POST   | "/category"      | createCategory          | not started |
| PUT    | "/category/:id"  | updateCategory          | not started |
| DELETE | "/category/:id"  | deleteCategory          | not started |
| GET    | "/promotion/:id" | getPromotion            | not started |
| POST   | "/promotion"     | createPromotion         | not started |
| PUT    | "/promotion/:id" | updatePromotion         | not started |
| DELETE | "/promotion/:id" | deletePromotion         | not started |
| GET    | "/inventory/:id" | getInventoryMovement    | not started |
| POST   | "/inventory"     | createInventoryMovement | not started |
| PUT    | "/inventory/:id" | updateInventoryMovement | not started |
| DELETE | "/inventory/:id" | deleteInventoryMovement | not started |
| GET    | "/customer/:id"  | getCustomer             | not started |
| POST   | "/customer"      | createCustomer          | not started |
| PUT    | "/customer/:id"  | updateCustomer          | not started |
| DELETE | "/customer/:id"  | deleteCustomer          | not started |
| GET    | "/invoice/:id"   | getInvoice              | not started |
| POST   | "/invoice"       | createInvoice           | not started |
| PUT    | "/invoice/:id"   | updateInvoice           | not started |
| DELETE | "/invoice/:id"   | deleteInvoice           | not started |

## TODO

- [ ] Create a function to create the default roles and update them from a JSON file. [^1]

## Ideas

- Use cache storing to store the products, and websockets to update it once the product list is modified:
  1. When the app starts, make a request to the server to get the top selling products and store them in a local cache.
  2. When searching the products, a product is not in the cache, get it from the server and store it.
  3. Get updates in real time using websockets
  4. Send the invoice to the server when the sale is done. If the server is not available, store it in a temporal file and send it once the server is available again and its reception is confirmed.

## Projects and documentation to learn from

<https://github.com/mongodb-university/atlas_starter_go/blob/master/main.go>

- search for mflix tutorial

[^1]: Right now it works by modifying the roles directly in the defaultRoles method in cmd/api/roles.go, and going to "/defaultroles" route.
