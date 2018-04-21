A simple task list application 

I've just wanted to taste the following technologies stack.

* Backend
  * [Golang](https://golang.org/) / [Echo framework](https://echo.labstack.com/)
  * [Websocket](https://github.com/gorilla/websocket) server push info to client
  
* Frontend
  * [VueJs](https://vuejs.org) / [Bulma CSS](https://bulma.io/)
  * [Websocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_client_applications)
  
# Install dependencies

```sh
go get github.com/gorilla/websocket
go get github.com/labstack/echo
go get github.com/labstack/echo/middleware
```

# Run

```sh
go run app.go
```

navigate to `localhost:5000`

# Notes

I experimented various technique on the stack
* Backend database is just a in-memory array store on the backend echo/golang server
* Frontend components communicate with each-other via event bus