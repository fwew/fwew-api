module github.com/fwew/fwew-api

go 1.20

require (
	github.com/fwew/fwew-lib/v5 v5.7.1-dev.0.20240218080308-828e3d59d279
	github.com/gorilla/mux v1.7.4
)

//for testing on a local machine's fwew-lib
replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
