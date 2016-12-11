package oauth2

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RangelReale/osin"
	"github.com/RangelReale/osin/example"
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

//Config .
type Config struct {
	Debug   bool
	Storage osin.Storage
}
type auth2Middleware struct {
	server *osin.Server
	config Config
}

//CheckBearerAuth req
func CheckBearerAuth(ctx *iris.Context) string {
	authHeader := string(ctx.RequestCtx.Request.Header.Peek("Authorization"))
	authForm := ctx.FormValueString("code")
	if authHeader == "" && authForm == "" {
		return ""
	}
	token := authForm
	if authHeader != "" {
		s := strings.SplitN(authHeader, " ", 2)
		if (len(s) != 2 || strings.ToLower(s[0]) != "bearer") && token == "" {
			return ""
		}
		//Use authorization header token only if token type is bearer else query string access token would be returned
		if len(s) > 0 && strings.ToLower(s[0]) == "bearer" {
			token = s[1]
		}
	}
	return token
}
func (b *auth2Middleware) init(config Config) {
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	b.config = config
	b.server = osin.NewServer(sconfig, config.Storage)
}

//Serve ctx
func (b *auth2Middleware) Serve(ctx *iris.Context) {
	if b.isGranted(ctx) {

		ctx.Next()
	}
}

//ServerHTTP res,req
func (b *auth2Middleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	authorize := "/authorize"
	token := "/token"
	if strings.HasSuffix(reqPath, authorize) {
		b.authorize(res, req)
	}
	if strings.HasSuffix(reqPath, token) {
		b.token(res, req)
	}
	println(reqPath)
}
func (b *auth2Middleware) isGranted(ctx *iris.Context) bool {
	code := CheckBearerAuth(ctx)
	ret, _ := b.server.Storage.LoadAccess(code)
	return ret != nil || b.config.Debug

}
func (b *auth2Middleware) token(w http.ResponseWriter, r *http.Request) {
	resp := b.server.NewResponse()
	defer resp.Close()

	if ar := b.server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			if ar.Username == "test" && ar.Password == "test" {
				ar.Authorized = true
			}
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
		case osin.ASSERTION:
			if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
				ar.Authorized = true
			}
		}
		b.server.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["custom_parameter"] = 19923
	}
	osin.OutputJSON(resp, w, r)
}

func (b *auth2Middleware) authorize(w http.ResponseWriter, r *http.Request) {
	resp := b.server.NewResponse()
	defer resp.Close()

	if ar := b.server.HandleAuthorizeRequest(resp, r); ar != nil {
		if !example.HandleLoginPage(ar, w, r) {
			return
		}
		ar.UserData = struct{ Login string }{Login: "test"}
		ar.Authorized = true
		b.server.FinishAuthorizeRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["custom_parameter"] = 187723
	}
	osin.OutputJSON(resp, w, r)
}

//New take config
func New(c Config) iris.HandlerFunc {
	b := &auth2Middleware{}
	b.init(c)
	//ret := iris.ToHandlerFunc(b.ServeHTTP)
	h := fasthttpadaptor.NewFastHTTPHandlerFunc(b.ServeHTTP)
	return iris.HandlerFunc((func(ctx *iris.Context) {
		h(ctx.RequestCtx)
		b.Serve(ctx)
	}))
}
