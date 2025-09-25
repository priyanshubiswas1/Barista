.PHONY: up up-build down up-db down-db auth-proto-gen
up:
	docker-compose up -d

up-build:
	docker-compose up --build -d

down:
	docker-compose down

up-db:
	docker-compose up user_db menu_db order_db payment_db session_db shop_db table_db -d

down-db:
	docker-compose down user_db menu_db order_db payment_db session_db shop_db table_db

auth-proto-gen:
	protoc \
	--go_out=./contracts/gen/auth --go_opt=paths=source_relative \
    --go-grpc_out=./contracts/gen/auth --go-grpc_opt=paths=source_relative \
    --proto_path=./contracts/protos \
    auth.proto

shop-proto-gen:
	protoc \
	--go_out=./contracts/gen/shop --go_opt=paths=source_relative \
    --go-grpc_out=./contracts/gen/shop --go-grpc_opt=paths=source_relative \
    --proto_path=./contracts/protos \
    shop.proto

menu-proto-gen:
	protoc \
	--go_out=./contracts/gen/menu --go_opt=paths=source_relative \
    --go-grpc_out=./contracts/gen/menu --go-grpc_opt=paths=source_relative \
    --proto_path=./contracts/protos \
    menu.proto

table-proto-gen:
	protoc \
	--go_out=./contracts/gen/table --go_opt=paths=source_relative \
    --go-grpc_out=./contracts/gen/table --go-grpc_opt=paths=source_relative \
    --proto_path=./contracts/protos \
    table.proto

session-proto-gen:
	protoc \
	--go_out=./contracts/gen/session --go_opt=paths=source_relative \
    --go-grpc_out=./contracts/gen/session --go-grpc_opt=paths=source_relative \
    --proto_path=./contracts/protos \
    session.proto

order-proto-gen:
	protoc \
	--go_out=./contracts/gen/order --go_opt=paths=source_relative \
    --go-grpc_out=./contracts/gen/order --go-grpc_opt=paths=source_relative \
    --proto_path=./contracts/protos \
    order.proto

