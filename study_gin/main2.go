package main

import (
	//"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Users struct {
	Id        int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Firstname string `gorm:"not null" form:"firstname" json:"firstname"`
	Lastname  string `gorm:"not null" form:"lastname" json:"lastname"`
}

type TestUser struct {
	PassWord string `json:"password" binding:"required"` // 密码
	Mobile   string `json:"mobile" binding:"required"`   // 电话
	NickName string `json:"nick_name"`                   // 昵称
	Icon     string `json:"icon"`                        // 头像
}

type NameArr9 struct {
	Name string `json:"name"`
}

type UsersInfo struct {
	Id    string          `gorm:"column:id" json:"id"`                 //
	Users json.RawMessage `gorm:"column:users;type:json" json:"users"` //
}

type UserInfo struct {
	UserId string `gorm:"column:userid" json:"userid"` //
	Remark string `json:"remark"`
}

type NameArr struct {
	Patoks []Reporters `json:"name" binding:"required"`
}

type Meeting struct {
	Date      time.Time   `json:"date"`
	Area      string      `json:"area"`
	Reporters []Reporters `json:"reporters" binding:"required"`
}

type Reporters struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}

/*

{
    "date": "2020-01-29T14:47:43.511Z",
    "area": "Hongkong",
    "reporters": [
        {
            "name":"jack",
            "age": 12
        },
        {
            "name":"rick",
            "age": 78
        }
    ]
}

*/
func handleMeeting(c *gin.Context) {

	var request Meeting
	if err := c.ShouldBindJSON(&request); err != nil {
		// c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	fmt.Println("request=", request)
	//  c.JSON(http.StatusOK, utils.Response("success"))
}

type Rpts struct {
	Reporters []Reporters `json:"reporters"`
}

/*

{
    "reporters": [
        {
            "name":"jack",
            "age": 12
        },
        {
            "name":"rick",
            "age": 78
        }
    ]
}

*/

func handleRpts(c *gin.Context) {

	var request Rpts
	if err := c.ShouldBindJSON(&request); err != nil {
		// c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	fmt.Println("request=", request)
	//  c.JSON(http.StatusOK, utils.Response("success"))
}

type Writers struct {
	Writers []Reporters
}

func handleWriters(c *gin.Context) {

	var request []Reporters

	bodyByts, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("raw_data=", string(bodyByts)) // 获取到原生态的请求数据

	if err := c.ShouldBindJSON(&request); err != nil {
		// c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	fmt.Println("request=", request)
	//  c.JSON(http.StatusOK, utils.Response("success"))
}

// CreateTeast 创建测试用户
func PostUser(c *gin.Context) {
	var postData []NameArr

	bodyByts, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("raw_data=", string(bodyByts)) // 获取到原生态的请求数据

	fmt.Println("postData=", postData)

	if err := c.ShouldBindJSON(&postData); err != nil {
		fmt.Println("err=", err)
		// response.ReturnJSON(c, http.StatusOK, statuscode.InvalidParam.Code,statuscode.InvalidParam.Msg, nil)
		return
	}
	// 走到这里，postData 里面就有数据了
}

// func PostUser(c *gin.Context) {
// 	// The futur code…
// 	fmt.Println("Post_User")
// }

type MsgJson struct {
	Msg string `json:"msg"`
}

func GetUsers(c *gin.Context) { // http://127.0.0.1:8080/api/v1/users
	fmt.Println("----------------get_user-----------")

	//lastname := c.Query("lastname")

	bodyByts, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("raw_data=", string(bodyByts)) // 获取到原生态的请求数据

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
	name := c.Params.ByName("name")
	fmt.Println("name=", name)

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
		v1.POST("/users", handleMeeting)

		v1.GET("/rpts", handleRpts)
		//v1.GET("/writers", handleWriters)  // 这个没有调通， 等下次调通再说

		v1.GET("/users/:id", GetUser)       // http://127.0.0.1:8080/api/v1/users/22    22就是id值
		v1.GET("/users/:id/:name", GetUser) // http://127.0.0.1:8080/api/v1/users/22/jack    22就是id,jack就是name
		v1.DELETE("/users/:id", DeleteUser)
	}

	r.Run(":8080")
}

type Transport struct {
	Time  string
	MAC   string
	Id    string
	Rssid string
}

func main() {
	RouterWithGroup()

}
