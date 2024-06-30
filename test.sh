#!/bin/bash

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
	echo
}

create_category() {
	echo "Creating a new category..."
	curl -X POST http://localhost:8080/category \
		-H "Content-Type: application/json" \
		-d '{
            "name": "Example Category"
        }'
	echo
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
	echo
}

create_user() {
	echo "Creating a new user..."
	curl -X POST http://localhost:8080/user \
		-H "Content-Type: application/json" \
		-d '{
            "name": "'$2'",
            "username": "'$3'",
            "password": "'$4'",
            "role_id": "'$5'"
        }'
	echo
}

search_user() {
	if [ -z "$1" ]; then
		echo "User ID is required"
		echo "Usage: $0 search-user <user_id>"
		exit 1
	fi

	echo "Searching user with ID: $1 ..."
	curl -X GET http://localhost:8080/user/$1
	echo
}

if [ "$1" = "create-product" ]; then
	create_product
elif [ "$1" = "create-category" ]; then
	create_category
elif [ "$1" = "create-roles" ]; then
	create_roles
elif [ "$1" = "create-user" ]; then
	if [ -z "$2" ] || [ -z "$3" ] || [ -z "$4" ]; then
		echo "Usage: $0 create-user <name> <username> <password> [role_id]"
		exit 1
	fi
	create_user "$@"
elif [ "$1" = "search-user" ]; then
	search_user $2
else
	echo "Usage: $0 {create-product|create-category|create-roles|create-user <name> <username> <password> [role_id]|search-user <user_id>}"
	exit 1
fi
