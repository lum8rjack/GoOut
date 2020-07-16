module github.com/lum8rjack/GoOut/server

replace github.com/lum8rjack/GoOut/server/modules/writefile => ./modules/writefile

replace github.com/lum8rjack/GoOut/server/modules/loadconfig => ./modules/loadconfig

replace github.com/lum8rjack/GoOut/server/modules/http => ./modules/http

replace github.com/lum8rjack/GoOut/server/modules/https => ./modules/https

replace github.com/lum8rjack/GoOut/server/modules/udp => ./modules/udp

replace github.com/lum8rjack/GoOut/server/modules/tcp => ./modules/tcp

go 1.14

require (
	github.com/lum8rjack/GoOut/server/modules/http v0.0.0-00010101000000-000000000000
	github.com/lum8rjack/GoOut/server/modules/https v0.0.0-00010101000000-000000000000
	github.com/lum8rjack/GoOut/server/modules/loadconfig v0.0.0-00010101000000-000000000000
	github.com/lum8rjack/GoOut/server/modules/tcp v0.0.0-00010101000000-000000000000
	github.com/lum8rjack/GoOut/server/modules/udp v0.0.0-00010101000000-000000000000
	github.com/lum8rjack/GoOut/server/modules/writefile v0.0.0-00010101000000-000000000000
)
