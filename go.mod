module github.com/fwew/fwew-api

go 1.13

require (
	github.com/fwew/fwew-lib/v6 v6.0.0-20230605225434-7ee72a1e0ab4
	github.com/gorilla/mux v1.7.4
)

//for testing on a local machine's fwew-lib
replace github.com/fwew/fwew-lib/v6 => ../fwew-lib
