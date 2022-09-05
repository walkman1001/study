package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloWord() { // 不带路径
	//请求参数  http://127.0.0.1:8080/    或者    http://127.0.0.1:8080

	fmt.Println("helloword")
	//1.创建路由
	r := gin.Default()
	//2.绑定路由规则，执行的函数
	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello World!")
	})
	//3.监听端口，默认8080
	r.Run(":8080")
}

func getPing(c *gin.Context) {
	//输出json结果给调用方, 200 表示请求成功
	// 请求参数 http://127.0.0.1:8080/ping?firstname=jack&lastname=xu&ci=jiangxi

	lastname := c.Query("lastname")
	fmt.Println("lastname=", lastname)
	firstname := c.Query("firstname")
	fmt.Println("firstname=", firstname)
	defaulValue := c.DefaultQuery("city", "hongkang") //如果去到city字段，就用取到的字段，如果取不到，就用缺省的hongkang
	fmt.Println("defaulValue=", defaulValue)

	c.JSON(200, gin.H{
		"message": "pong8",
	})
}

func SimpleRouter() { //带简单路径

	r := gin.Default()
	r.GET("/ping", getPing)

	r.Run(":8080")
}

func Router() {
	r := gin.Default()

	//http://127.0.0.1:8080/user/jack9

	r.GET("/user/:name", func(c *gin.Context) { // :name, 表示name就是一个值，不是路径,比如上面的 jack就是值
		name := c.Param("name")
		c.String(http.StatusOK, "---Hello %s", name)
	})

	//http://127.0.0.1:8080/user/jack/   这个参数请求， 表示 jack是个路径，不是参数 name，如果路径完全匹配下面的，优先，否则就匹配有：的参数url

	r.GET("/user/jack", func(c *gin.Context) { // :name, 表示name就是一个值，不是路径,比如上面的 jack就是值
		name := c.Param("name")
		fmt.Println("name6=", name)  // name6=  是个空值
		c.String(http.StatusOK, "---Hello %s", name)

	})

	// http://127.0.0.1:8080/staff/23     // id=/23 注意多了一个/,而不是 id=23
	// http://127.0.0.1:8080/staff/23/qty   //id=/23/qty , 相当于后面被截断了
	// http://127.0.0.1:8080/staff/23/qty/88  //id= /23/qty/88

	r.GET("/staff/*id", func(c *gin.Context) { // :name, 表示name就是一个值，不是路径,比如上面的 jack就是值
		id := c.Param("id")
		fmt.Println("id=", id) // id=/23 注意多了一个/,而不是 id=23

		c.String(http.StatusOK, "---Hello_id= %s", id)
	})

	r.Run(":8080")
}
