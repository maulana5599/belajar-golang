package main

import (
	"encoding/json"
	"fmt"
	"golang-dasar/database"
	"io"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Barang struct {
	Id         int     `json:"id"`
	NamaBarang string  `json:"nama_barang"`
	StokBarang float32 `json:"stok_barang,omitempty"`
}

var mutex sync.Mutex

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	database.ConnectionDatabase()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/test", apiTransaction)
	e.GET("/reflect", detail)
	e.GET("/channel", channelHandler)
	e.GET("/request", getHttpRequest)
	e.Logger.Fatal(e.Start(":3000"))
}

func handlerMutex(ctx echo.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	mutex.Lock()
	var result Barang
	database.DB.QueryRow("SELECT * FROM stok_barang ").Scan(&result.Id, &result.NamaBarang, &result.StokBarang)

	id := 1
	jumlahstok := result.StokBarang - 1
	_, err := database.DB.ExecContext(ctx.Request().Context(), "UPDATE stok_barang SET stok_barang = ? where id = ?", jumlahstok, id)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatalln(err.Error())
	} else {
		fmt.Println("Successfully updated stok_barang")
	}
	mutex.Unlock()
}

func apiTransaction(ctx echo.Context) error {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			handlerMutex(ctx, &wg)
		}()
	}

	wg.Wait()

	return ctx.JSON(200, echo.Map{
		"message": "Success",
	})
}

func detail(ctx echo.Context) error {

	message := map[string]interface{}{
		"Message": "Hello Maulana",
	}

	ext := reflect.ValueOf(message)
	fmt.Println("Tipe data :", ext.MapKeys())
	return ctx.JSON(200, echo.Map{
		"halo": message,
	})
}

func channelHandler(ctx echo.Context) error {

	runtime.GOMAXPROCS(2)

	var message = make(chan interface{})
	// var message = make(chan interface{}, 10) // Buffered Channel

	go messageChannel("Maulana Muhammad Rizky", message)

	var message1 = <-message

	return ctx.JSON(200, echo.Map{
		"message":   message1,
		"JumlahCpu": runtime.NumCPU(),
	})
}

func messageChannel(nama string, messages chan interface{}) {
	var data = fmt.Sprintf("Nama saya : %s", nama)

	// Penggunaan Buffered Channel
	// for i := 0; i < 1000; i++ {
	// 	messages <- echo.Map{
	// 		"Nama": data,
	// 	}

	// 	fmt.Println("Index", i)
	// }

	messages <- echo.Map{
		"Nama": data,
	}
}

type ResponseDummy struct {
	Limit int           `json:"limit"`
	Posts []interface{} `json:"posts"`
	Skip  int           `json:"skip"`
	Total int           `json:"total"`
}

func getHttpRequest(ctx echo.Context) error {
	// var client = &http.Client{}
	response, err := http.Get("https://dummyjson.com/posts?skip=5&limit=10")
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return ctx.JSON(404, err.Error())
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ctx.JSON(404, err.Error())
	}

	// var stringData map[string]interface{}
	var stringData ResponseDummy

	err = json.Unmarshal(body, &stringData)

	if err != nil {
		return ctx.JSON(404, err.Error())
	}

	data := stringData.Posts
	for _, v := range data {
		fmt.Println(v)
	}

	return ctx.JSON(200, echo.Map{
		"Message": "Get message success",
		"Data":    stringData,
	})
}
