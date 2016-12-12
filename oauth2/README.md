## How to use
```go
package main

import (
	"github.com/tomasky/middleware/oauth2"

	"github.com/RangelReale/osin/example"
	"github.com/kataras/iris"
)

//User name ,email
type User struct {
	Uid   int
	Name  string `json:"name" xml:"name" form:"name"`
	Email string `json:"email" xml:"email" form:"email"`
}

func main() {
	r := iris.Party("/login/oauth")
	g := iris.Party("/api/v1")
	authentication := oauth2.New(oauth2.Config{Debug: true, Storage: example.NewTestStorage()})
	g.Use(authentication)
	r.Use(authentication)
	r.Get("/authorize", nil)
	r.Get("/token", nil)

	g.Get("/user/:id", func(ctx *iris.Context) {

		// take the :id from the path, parse to integer
		// and set it to the new userID local variable.
		//userID, er := ctx.ParamInt("id")

		// userRepo, imaginary database service <- your only job.
		//user := userRepo.GetByID(userID)
		if userID, er := ctx.ParamInt("id"); er == nil {
			user := User{UID: userID}
			ctx.JSON(iris.StatusOK, user)
			return
		}
		// userRepo, imaginary database service <- your only job.
		//user := userRepo.GetByID(userID)
		// send back a response to the client,
		// .JSON: content type as application/json; charset="utf-8"
		// iris.StatusOK: with 200 http status code.
		//
		// send user as it is or make use of any json valid golang type,
		ret := iris.Map{"status": iris.StatusOK, "result": nil}
		ctx.JSON(iris.StatusOK, ret)

	})

	iris.Listen("localhost:5800")
}
```
