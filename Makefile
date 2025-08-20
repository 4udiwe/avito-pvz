gen:
	go tool oapi-codegen --config=internal/gen/dto.gen.yaml api/swagger.yaml 
	go tool oapi-codegen --config=internal/gen/handler.gen.yaml api/swagger.yaml 