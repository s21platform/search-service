protogen:
	protoc --go_out=. --go-grpc_out=. ./api/search.proto --experimental_allow_proto3_optional
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/search.proto --experimental_allow_proto3_optional
