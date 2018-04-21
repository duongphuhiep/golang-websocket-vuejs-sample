package main

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//log error stacktrace before returning error page
func customHTTPErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	ctx.Logger().Errorf("%+v", err)
	ctx.JSON(code, http.StatusText(code))
}

var (
	upgrader = websocket.Upgrader{}
)

var wsStore = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func createWsConnection(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	wsStore[ws] = true //save the clients connection

	for {
		// Read new message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Infof("connection lost when read %v", err)
			delete(wsStore, ws)
		}

		task := string(msg)

		//save the task to store
		store = append(store, task)

		//and broadcast it to other clients
		broadcast <- task

		// // Write
		// err := ws.WriteMessage(websocket.TextMessage, []byte("database is changed"))
		// if err != nil {
		// 	c.Logger().Error(err)
		// }
	}
}

func broadcastTask(e *echo.Echo) {
	for {
		task := <-broadcast
		for ws := range wsStore {
			err := ws.WriteMessage(websocket.TextMessage, []byte(task))
			if err != nil {
				e.Logger.Infof("connection lost when read %v", err)
				delete(wsStore, ws)
			}
		}
	}
}

// User
type Task struct {
	Task string `json:"task" form:"task" query:"task"`
}

const SERVER_URL string = ":5000"

var store []string

func main() {
	e := echo.New()
	e.Debug = true
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Logger.SetOutput(os.Stdout)

	//e.Use(middleware.Logger())
	currentDir, err := os.Getwd()
	if err != nil {
		e.Logger.Fatal(err)
	}

	//e.Static("assets", "assets")
	e.Logger.Infof("Current folder: %s", currentDir)
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.PUT("/add", func(c echo.Context) error {
		u := new(Task)
		if err = c.Bind(u); err != nil {
			c.Logger().Fatal(err)
		}
		task := u.Task
		c.Logger().Infof("Add task without broadcast to other client: %s", task)
		store = append(store, task)
		return c.String(http.StatusOK, "Added")
	})

	e.GET("/all", func(c echo.Context) error {
		return c.JSON(http.StatusOK, store)
	})

	e.GET("/ws", createWsConnection)
	go broadcastTask(e)
	e.Logger.Infof("server starts listening %s", SERVER_URL)
	e.Logger.Fatal(e.Start(SERVER_URL))
}
