module github.com/fwew/fwew-api

go 1.20

require (
	github.com/fwew/fwew-lib/v5 v5.17.1
	github.com/gorilla/mux v1.8.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
)

//for testing on a local machine's fwew-lib
replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
