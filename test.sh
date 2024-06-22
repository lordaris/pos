# Run sh test.sh to run the scripts
create_product() {
	echo "Creating a new product..."
	curl -X POST http://localhost:8080/product \
		-H "Content-Type: application/json" \
		-d '{
            "name": "Example Product with category 2",
            "brand": "Example Brand",
            "description": "This is an example product",
            "price": 19.99,
            "stock": 100,
            "min_stock": 10,
            "barcode": "1234567890123",
            "plu": 12345,
            "category_id": "6675dceb834a3b2b5e254b65"
        }'
}

create_category() {
	echo "Creating a new category..."
	curl -X POST http://localhost:8080/category \
		-H "Content-Type: application/json" \
		-d '{
            "name": "Example Category"
        }'
}

create_roles() {
	echo "Creating roles..."
	curl -X POST http://localhost:8080/roles \
		-H "Content-Type: application/json" \
		-d '[
            {
                "name": "admin",
                "permissions": ["all_permissions"]
            },
            {
                "name": "manager",
                "permissions": [
                    "manage_users", "manage_inventory", "data_visualization", "management_reports",
                    "purchase_request_authorization", "inventory_transfer_authorization"
                ]
            },
            {
                "name": "cashier",
                "permissions": ["pos"] 
            }
        ]'
	echo "\n"
}
if [ "$1" = "create-product" ]; then
	create_product
elif [ "$1" = "create-category" ]; then
	create_category
elif [ "$1" = "create-roles" ]; then
	create_roles
else
	echo "Usage: $0 {create-product|create-category|create-roles}"
	exit 1
fi
