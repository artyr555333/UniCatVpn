package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Login struct {
	Password string `json:"password"`
}

func init() {
	f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.SetOutput(f)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	r := gin.Default()

	r.POST("/login", PasswordCheck)
	err := r.Run()
	if err != nil {
		return
	}
}

func PasswordCheck(c *gin.Context) {
	var login Login
	fmt.Println(login.Password)
	if err := c.ShouldBind(&login); err != nil {
		fmt.Println("AZAMAT LOH")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if login.Password != os.Getenv("CORRECT_PASS") {
		fmt.Println("AZAMAT GAY")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный пароль"})
	}

	if login.Password == os.Getenv("CORRECT_PASS") {
		fmt.Println("AZAMAT CHLEN")
		c.JSON(http.StatusOK, gin.H{"status": "Пароль принят"})
	}
}
