package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)





func getParameters() {
	router2 := gin.Default() // http://127.0.0.1:8080/welcome?firstname=Jane&lastname=Doe
	router2.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest") // 如果过去到参数firstname, 就用或取到的，否则就用缺省的 "Gest" 作为替代
		lastname := c.Query("lastname")                   // 单纯获取，没有获取到，就没有获取到

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router2.Run(":8080")
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	usrObj := User{}
	c.BindJSON(&usrObj) // 获取到参数，然后map到对象里面去

	fmt.Println("-----------------%v", &usrObj)
	c.JSON(http.StatusOK, gin.H{
		"name":     usrObj.Name,
		"password": usrObj.Password,
	})
}

func LoginAccep() {
	router := gin.Default() // http://127.0.0.1:8080/login
	router.GET("/login", Login)
	router.Run(":8080")
}
func getRawData() {
	/*
	   http://127.0.0.1:8080/login

	   raw 原生态数据请求

	   {
	   "name":"jack"
	   ,"password":"abc123"
	   }

	*/

	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		bodyByts, err := ioutil.ReadAll(c.Request.Body)
		fmt.Println("raw_data=", string(bodyByts)) // 获取到原生态的请求数据

		if err != nil {
			// 返回错误信息
			c.String(http.StatusBadRequest, err.Error())
			// 执行退出
			c.Abort()
		}

		// 返回的 code 和 对应的参数星系
		c.String(http.StatusOK, "%s \n", string(bodyByts))
	})
	r.Run(":8080")
}
func main2() {
	LoginAccep()
}
