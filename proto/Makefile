.PHONY: proto-gen

proto-gen:
	@echo "Generating protobuf code..."
	protoc \
	  --proto_path=proto/contracts \
	  --go_out=proto/gen \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=proto/gen \
	  --go-grpc_opt=paths=source_relative \
	  dice.proto
	@echo "Dice service protobuf code generated"


# Запуск: make proto-gen

# Генерация прото файлов для dice сервиса
# protoc --proto_path=proto/contracts --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative dice.proto
