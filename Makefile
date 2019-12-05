all: sola auth cors graphql logger native proxy rest router swagger ws x.router
	echo "success!"

sola: .
	go build -o ./dist/sola.a .
auth: ./middleware/auth
	go build -o ./dist/middleware/auth.a ./middleware/auth
cors: ./middleware/cors
	go build -o ./dist/middleware/cors.a ./middleware/cors
graphql: ./middleware/graphql
	go build -o ./dist/middleware/graphql.a ./middleware/graphql
logger: ./middleware/logger
	go build -o ./dist/middleware/logger.a ./middleware/logger
native: ./middleware/native
	go build -o ./dist/middleware/native.a ./middleware/native
proxy: ./middleware/proxy
	go build -o ./dist/middleware/proxy.a ./middleware/proxy
rest: ./middleware/rest
	go build -o ./dist/middleware/rest.a ./middleware/rest
router: ./middleware/router
	go build -o ./dist/middleware/router.a ./middleware/router
swagger: ./middleware/swagger
	go build -o ./dist/middleware/swagger.a ./middleware/swagger
ws: ./middleware/ws
	go build -o ./dist/middleware/ws.a ./middleware/ws
x.router: ./middleware/x/router
	go build -o ./dist/middleware/x/router.a ./middleware/x/router
