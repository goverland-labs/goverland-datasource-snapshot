module github.com/goverland-labs/goverland-datasource-snapshot

go 1.22

toolchain go1.22.3

replace github.com/goverland-labs/goverland-datasource-snapshot/protocol => ./protocol

require (
	github.com/Yamashou/gqlgenc v0.14.0
	github.com/akamensky/argparse v1.4.0
	github.com/caarlos0/env/v8 v8.0.0
	github.com/ethereum/go-ethereum v1.13.14
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/goverland-labs/goverland-datasource-snapshot/protocol v0.0.0
	github.com/goverland-labs/goverland-ipfs-fetcher/protocol v0.0.1
	github.com/goverland-labs/goverland-platform-events v0.3.12
	github.com/goverland-labs/snapshot-sdk-go v0.4.2
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.1
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/nats-io/nats.go v1.30.2
	github.com/pereslava/grpc_zerolog v0.0.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.18.0
	github.com/rs/zerolog v1.29.1
	github.com/s-larionov/process-manager v0.0.1
	github.com/samber/lo v1.38.1
	github.com/shutter-network/shutter/shlib v0.1.11
	github.com/smartystreets/goconvey v1.8.0
	github.com/vektah/gqlparser/v2 v2.5.3
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.2
)

require (
	github.com/99designs/gqlgen v0.17.33 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/nats-io/nkeys v0.4.5 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/smartystreets/assertions v1.13.1 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
)
