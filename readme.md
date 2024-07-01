# POS

## Endpoints

| Method | URL pattern      | Handler                 | Status      |
| ------ | ---------------- | ----------------------- | ----------- |
| GET    | "/user/:id"      | getUser                 | Not started |
| POST   | "/user"          | createUser              | In progress |
| PUT    | "/user/:id"      | updateUser              | In progress |
| DELETE | "/user/:id"      | deleteUser              | Not started |
| PUT    | "/user/:id/role" | updateUserRole          | In progress |
| GET    | "/product/:id"   | getProduct              | Not started |
| POST   | "/product"       | createProduct           | In progress |
| PUT    | "/product/:id"   | updateProduct           | Not started |
| DELETE | "/product/:id"   | deleteProduct           | Not started |
| GET    | "/category"      | getCategory             | Not started |
| POST   | "/category"      | createCategory          | Not started |
| PUT    | "/category/:id"  | updateCategory          | Not started |
| DELETE | "/category/:id"  | deleteCategory          | Not started |
| GET    | "/promotion/:id" | getPromotion            | Not started |
| POST   | "/promotion"     | createPromotion         | Not started |
| PUT    | "/promotion/:id" | updatePromotion         | Not started |
| DELETE | "/promotion/:id" | deletePromotion         | Not started |
| GET    | "/inventory/:id" | getInventoryMovement    | Not started |
| POST   | "/inventory"     | createInventoryMovement | Not started |
| PUT    | "/inventory/:id" | updateInventoryMovement | Not started |
| DELETE | "/inventory/:id" | deleteInventoryMovement | Not started |
| GET    | "/customer/:id"  | getCustomer             | Not started |
| POST   | "/customer"      | createCustomer          | Not started |
| PUT    | "/customer/:id"  | updateCustomer          | Not started |
| DELETE | "/customer/:id"  | deleteCustomer          | Not started |
| GET    | "/invoice/:id"   | getInvoice              | Not started |
| POST   | "/invoice"       | createInvoice           | Not started |
| PUT    | "/invoice/:id"   | updateInvoice           | Not started |
| DELETE | "/invoice/:id"   | deleteInvoice           | Not started |

## Ideas

- Use cache storage to store the products, and websockets to update it once the product list is modified:
  1. When the app starts, make a request to the server to get the top-selling products and store them in a local cache.
  2. When searching for products, if a product is not in the cache, retrieve it from the server and store it.
  3. Get updates in real-time using websockets.
  4. Send the invoice to the server when the sale is completed. If the server is not available, store it in a temporary file and send it once the server is available again and its reception is confirmed.

## API Documentation

To create a promotion, the user should select one of three types: `DiscountPrice`, `DiscountPercentage`, or `BuyGet`. Even if more than one of the promotion fields are being used, the type will be the field to select which promotion to apply.

- **DiscountPercentage:** This applies a percentage discount to the price of a product. When selected, the system should calculate the discounted price and use this as the product price instead of the regular price. It receives an int as a percentage
- **DiscountPrice:** This sets a new fixed price for the product. If active, this price should be used instead of the regular price.
- **BuyGet:** This promotion has two fields: the number of items the customer needs to buy (`BuyQuantity`) and the number of items the customer gets for free (`GetQuantity`). The system should detect the quantity of products bought and adjust the total price accordingly by discounting the price of the free products.

---

test.sh have a method to create default roles. Modify it and run it to create new roles and its permissions.

## Projects and Documentation to Learn From

- [Atlas Starter Go](https://github.com/mongodb-university/atlas_starter_go/blob/master/main.go)

- Search for the "mflix tutorial"

---
