package main

import (
	//"net/http"

	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Users struct {
	Id        int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Firstname string `gorm:"not null" form:"firstname" json:"firstname"`
	Lastname  string `gorm:"not null" form:"lastname" json:"lastname"`
}

func PostUser(c *gin.Context) {
	// The futur code…
	fmt.Println("Post_User")
}

type MsgJson struct {
	Msg string `json:"msg"`
}

func GetUsers(c *gin.Context) { // http://127.0.0.1:8080/api/v1/users
	fmt.Println("----------------get_user-----------")

	//lastname := c.Query("lastname")

	lastname := c.Query("lastname")
	firstname := c.Query("firstname")

	fmt.Println("lastname=", lastname)
	fmt.Println("first=", firstname)

	var users = []Users{
		Users{Id: 1, Firstname: "Oliver", Lastname: "Queen"},
		Users{Id: 2, Firstname: "Malcom", Lastname: "Merlyn"},
	}

	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	user_id, _ := strconv.ParseInt(id, 0, 64)

	if user_id == 1 {
		content := gin.H{"id": user_id, "firstname": "Oliver", "lastname": "Queen"}
		c.JSON(200, content)
	} else if user_id == 2 {
		content := gin.H{"id": user_id, "firstname": "Malcom", "lastname": "Merlyn"}
		c.JSON(200, content)
	} else {
		content := gin.H{"error": "user with id#" + id + " not found"}
		c.JSON(404, content)
	}

	// curl -i http://localhost:8080/api/v1/users/1
}

func UpdateUser(c *gin.Context) {
	// The futur code…
}

func DeleteUser(c *gin.Context) {
	// The futur code…
}

func GetJackUsers(c *gin.Context) {
	fmt.Println("helel jack")

	c.JSON(200, gin.H{
		"message": "hello_jack",
	})
}
func GroupWithRootPath() {
	r := gin.Default()

	v1 := r.Group("") //不带任何分组路径，比如 http://127.0.0.1:8080/jack
	{
		v1.GET("/jack", GetJackUsers)
	}

	r.Run(":8080")
}

func RouterWithGroup() {
	r := gin.Default()

	v1 := r.Group("/api/v1") //不带任何分组路径，比如 http://127.0.0.1:8080/users
	{
		v1.GET("/users", GetUsers)
		v1.POST("/users", PostUser)
		v1.GET("/users/:id", GetUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
	}

	r.Run(":8080")
}
func main8() {
	RouterWithGroup()
}
