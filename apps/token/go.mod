module mandacode.com/accounts/token

go 1.24.4

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	go.uber.org/mock v0.5.2
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.73.0
	mandacode.com/accounts/proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace mandacode.com/accounts/proto => ../proto
