package sfasthttp

import (
	"net/http"
	"sync"

	"github.com/fasthttp/session/v2"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/host/shttp"
	"github.com/valyala/fasthttp"
)

var (
	_ctxPool = &sync.Pool{
		New: func() interface{} {
			return new(FastHttpContext)
		},
	}
)

type FastHttpContext struct {
	ctx       *fasthttp.RequestCtx
	sess      *session.Session
	sessStore *session.Store
}

func NewFastHttpContext(ctx *fasthttp.RequestCtx, sess *session.Session) shttp.IHttpContext {
	r := _ctxPool.Get().(*FastHttpContext)
	r.ctx = ctx
	r.sess = sess
	var err error
	r.sessStore, err = r.sess.Get(ctx)
	u.LogFaltal(err)
	return r
}

func PutFastHttpContext(ctx shttp.IHttpContext) {
	_ctxPool.Put(ctx)
}

func (x *FastHttpContext) SetCookie(cookie *http.Cookie) {
	c := new(fasthttp.Cookie)
	c.SetKey(cookie.Name)
	c.SetValue(cookie.Value)
	c.SetMaxAge(cookie.MaxAge)
	c.SetDomain(cookie.Domain)
	c.SetPath(cookie.Path)
	c.SetSecure(cookie.Secure)
	c.SetHTTPOnly(cookie.HttpOnly)
	c.SetExpire(cookie.Expires)
	x.ctx.Response.Header.SetCookie(c)
}
func (x *FastHttpContext) GetCookieString(key string) string {
	r := x.ctx.Request.Header.Cookie(key)
	return u.BytesToStr(r)
}
func (x *FastHttpContext) RemoveCookie(key string) {
	x.ctx.Response.Header.DelClientCookie(key)
}

func (x *FastHttpContext) SetSession(key, value string) {
	store, err := x.sess.Get(x.ctx)
	if u.LogError(err) {
		return
	}
	defer func() {
		u.LogError(x.sess.Save(x.ctx, store))
	}()
	store.Set(key, value)
}
func (x *FastHttpContext) GetSessionString(key string) string {
	store, err := x.sess.Get(x.ctx)
	if u.LogError(err) {
		return ""
	}
	defer func() {
		u.LogError(x.sess.Save(x.ctx, store))
	}()

	if r, ok := store.Get(key).(string); ok {
		return r
	}

	return ""
}
func (x *FastHttpContext) RemoveSession(key string) {
	store, err := x.sess.Get(x.ctx)
	if u.LogError(err) {
		return
	}
	defer func() {
		u.LogError(x.sess.Save(x.ctx, store))
	}()
	store.Delete(key)
}
func (x *FastHttpContext) EndSession() {
	x.sess.Destroy(x.ctx)
	// store, err := x.sess.Get(x.ctx)
	// if u.LogError(err) {
	// 	return
	// }

	// store.Reset()
}

func (x *FastHttpContext) GetFormString(key string) string {
	r := x.ctx.FormValue(key)
	return u.BytesToStr(r)
}

func (x *FastHttpContext) GetBodyString() string {
	return x.ctx.Request.String()
}

func (x *FastHttpContext) GetBodyBytes() []byte {
	return x.ctx.Request.Body()
}

func (x *FastHttpContext) Redirect(url string, statusCode int) {
	x.ctx.Redirect(url, statusCode)
}
func (x *FastHttpContext) WriteString(body string) (int, error) {
	return x.ctx.WriteString(body)
}
func (x *FastHttpContext) WriteBytes(body []byte) (int, error) {
	return x.ctx.Write(body)
}

func (x *FastHttpContext) CopyBodyAndStatusCode(resp *http.Response) {
	x.ctx.SetStatusCode(resp.StatusCode)
	x.ctx.SetBodyStream(resp.Body, -1)
}
func (x *FastHttpContext) SetStatusCode(statusCode int) {
	x.ctx.SetStatusCode(statusCode)
}
func (x *FastHttpContext) SetContentType(cType string) {
	x.ctx.SetContentType(cType)
}

func (x *FastHttpContext) RequestURL() string {
	return x.ctx.URI().String()
}

func (x *FastHttpContext) SetHeader(key, value string) {
	x.ctx.Response.Header.Set(key, value)
}
