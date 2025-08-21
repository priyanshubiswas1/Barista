.PHONY: up down up-db down-db auth-proto-gen
up:
	docker-compose up -d

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


