## How to use
```go
package main

import (
	"oauth2"

	"github.com/RangelReale/osin/example"
	"github.com/kataras/iris"
)

//User name ,email
type User struct {
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
		user := new(User)
		// send back a response to the client,
		// .JSON: content type as application/json; charset="utf-8"
		// iris.StatusOK: with 200 http status code.
		//
		// send user as it is or make use of any json valid golang type,
		// like the iris.Map{"username" : user.Username}.
		ctx.JSON(iris.StatusOK, user)

	})

	iris.Listen("localhost:5800")
}```
