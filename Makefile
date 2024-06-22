# Makefile alternative to the sh file
# which I don't like because of all of the excesive escape symbols

create-product: 
	@echo "Creating a new product..."
	@curl -X POST http://localhost:8080/product \
		-H "Content-Type: application/json" \
		-d "{\"name\":\"Example Product with category 2\",\"brand\":\"Example Brand\",\"description\":\"This is an example product\",\"price\":19.99,\"stock\":100,\"min_stock\":10,\"barcode\":\"1234567890123\",\"plu\":12345,\"category_id\":\"6675dceb834a3b2b5e254b65\"}"


create-category:
	@echo "Creating a new category..."
	@curl -X POST http://localhost:8080/category \
		-H "Content-Type: application/json" \
		-d "{\"name\":\"Example Category\"}"

.PHONY: all
all: create-product
