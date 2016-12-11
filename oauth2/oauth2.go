package oauth2

import (
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/RangelReale/osin/example"
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

//Config .
type Config struct {
}
type auth2Middleware struct {
	server *osin.Server
}

func (b *auth2Middleware) init() {
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	b.server = osin.NewServer(sconfig, example.NewTestStorage())
}

//Serve ctx
func (b *auth2Middleware) Serve(ctx *iris.Context) {
	println("request to /products")
	ctx.Next()
}

//ServerHTTP res,req
func (b *auth2Middleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	code := ""
	ret, _ := b.server.Storage.LoadAccess(code)
}
func (b *auth2Middleware) isGrant(ctx *iris.Context) bool {
	return true
}
func (b *auth2Middleware) token(ctx *iris.Context) {
}
func (b *auth2Middleware) authorized(ctx *iris.Context) {
}

//New take config
func New(c Config) iris.HandlerFunc {
	b := &auth2Middleware{}
	b.init()
	//ret := iris.ToHandlerFunc(b.ServeHTTP)
	h := fasthttpadaptor.NewFastHTTPHandlerFunc(b.ServeHTTP)
	return iris.HandlerFunc((func(ctx *iris.Context) {
		h(ctx.RequestCtx)
		b.Serve(ctx)
	}))
}
