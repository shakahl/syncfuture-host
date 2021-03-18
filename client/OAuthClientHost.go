package client

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/muesli/cache2go"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/host"
	"golang.org/x/oauth2"
)

type OAuthClientHost struct {
	host.BaseHost
	OAuthOptions        *host.OAuthOptions `json:"OAuth,omitempty"`
	HashKey             string
	BlockKey            string
	UserJsonSessionKey  string
	UserIDSessionKey    string
	TokenCookieName     string
	SignInPath          string
	SignInCallbackPath  string
	SignOutPath         string
	SignOutCallbackPath string
	AccessDeniedPath    string
	OAuthClientHandler  host.IOAuthClientHandler
	ContextTokenStore   host.IContextTokenStore
	CookieProtoector    *securecookie.SecureCookie
	UserLocks           *cache2go.CacheTable
}

func (x *OAuthClientHost) BuildOAuthClientHost() {
	// if x.BaseWebHost == nil {
	// 	x.BaseWebHost = new(host.BaseWebHost)
	// }
	x.BaseHost.BuildBaseHost()

	if x.OAuthOptions == nil {
		log.Fatal("OAuth secion in configuration is missing")
	}
	x.OAuthOptions.BuildOAuthOptions(x.URLProvider)

	if x.BlockKey == "" {
		log.Fatal("block key cannot be empty")
	}
	if x.HashKey == "" {
		log.Fatal("hash key cannot be empty")
	}
	if x.SignInPath == "" {
		x.SignInPath = "/signin"
	}
	if x.SignInCallbackPath == "" {
		x.SignInCallbackPath = "/signin-oauth"
	}
	if x.SignOutPath == "" {
		x.SignOutPath = "/signout"
	}
	if x.SignOutCallbackPath == "" {
		x.SignOutCallbackPath = "/signout-oauth"
	}
	if x.AccessDeniedPath == "" {
		x.AccessDeniedPath = "/accessdenied"
	}
	if x.UserJsonSessionKey == "" {
		x.UserJsonSessionKey = "USERJSON"
	}
	if x.UserIDSessionKey == "" {
		x.UserIDSessionKey = "USERID"
	}
	if x.TokenCookieName == "" {
		x.TokenCookieName = "go.cookie2"
	}

	////////// user locks
	if x.UserLocks == nil {
		x.UserLocks = cache2go.Cache("UserLocks")
	}

	////////// cookie protoector
	if x.CookieProtoector == nil {
		x.CookieProtoector = securecookie.New(u.StrToBytes(x.HashKey), u.StrToBytes(x.BlockKey))
	}

	////////// context token store
	if x.ContextTokenStore == nil {
		x.ContextTokenStore = host.NewCookieTokenStore(x.TokenCookieName, x.CookieProtoector)
	}

	////////// oauth client handler
	if x.OAuthClientHandler == nil {
		x.OAuthClientHandler = NewOAuthClientHandler(x.OAuthOptions, x.ContextTokenStore, x.UserJsonSessionKey, x.UserIDSessionKey, x.TokenCookieName)
	}

	// ////////// auth middleware
	// if x.authMiddleware == nil {
	// 	x.authMiddleware = newClientAuthMiddleware(x.UserJsonSessionKey, x.AccessDeniedPath, x.OAuthOptions, x.PermissionAuditor)
	// }
}

func (x *OAuthClientHost) GetHttpClient() (*http.Client, error) {
	return x.OAuthOptions.ClientCredential.Client(context.Background()), nil
}

func (x *OAuthClientHost) GetUserHttpClient(ctx host.IHttpContext) (*http.Client, error) {
	// goctx := context.Background()
	// userID := host.GetUserID(ctx, x.UserIDSessionKey)
	// if userID == "" {
	// 	return http.DefaultClient, nil
	// }

	// // 获取用户锁
	// userLock := x.getUserLock(userID)

	// // read lock
	// userLock.RLock()
	// t, err := x.ContextTokenStore.GetToken(ctx)
	// defer func() { userLock.RUnlock() }()
	// if err != nil {
	// 	return http.DefaultClient, err
	// }

	// tokenSource := x.OAuthOptions.TokenSource(goctx, t)
	// newToken, err := tokenSource.Token()
	// if err != nil {
	// 	// refresh token failed, sign user out
	// 	host.SignOut(ctx, x.TokenCookieName)
	// 	return http.DefaultClient, err
	// }

	// if newToken.AccessToken != t.AccessToken {
	// 	// token been refreshed, lock
	// 	userLock.Lock()
	// 	// save token to session
	// 	err = x.ContextTokenStore.SaveToken(ctx, newToken)
	// 	// unlock
	// 	defer func() { userLock.Unlock() }()
	// 	if err != nil {
	// 		return http.DefaultClient, err
	// 	}
	// }

	tokenSource, err := x.GetUserToken(ctx)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(context.Background(), *tokenSource), nil
}

func (x *OAuthClientHost) GetClientToken(ctx host.IHttpContext) (*oauth2.Token, error) {
	return x.OAuthOptions.ClientCredential.Token()
}

func (x *OAuthClientHost) GetUserToken(ctx host.IHttpContext) (*oauth2.TokenSource, error) {
	goctx := context.Background()
	userID := host.GetUserID(ctx, x.UserIDSessionKey)
	if userID == "" {
		return nil, errors.New("user isn't authenticated")
	}

	// 获取用户锁
	userLock := x.getUserLock(userID)

	// read lock
	userLock.RLock()
	t, err := x.ContextTokenStore.GetToken(ctx)
	defer func() { userLock.RUnlock() }()
	if err != nil {
		return nil, err
	}

	tokenSource := x.OAuthOptions.TokenSource(goctx, t)
	newToken, err := tokenSource.Token()
	if err != nil {
		// refresh token failed, sign user out
		host.SignOut(ctx, x.TokenCookieName)
		return nil, err
	}

	if newToken.AccessToken != t.AccessToken {
		// token been refreshed, lock
		userLock.Lock()
		// save token to session
		err = x.ContextTokenStore.SaveToken(ctx, newToken)
		// unlock
		defer func() { userLock.Unlock() }()
		if err != nil {
			return nil, err
		}
	}

	return &tokenSource, nil
}

func (x *OAuthClientHost) AuthHandler(ctx host.IHttpContext) {
	routeKey := ctx.GetItemString(host.Item_JWT)
	if routeKey == "" {
		ctx.SetStatusCode(500)
		ctx.WriteString("route key does not exist")
		return
	}

	routes := strings.Split(routeKey, "_")

	var area, controller, action string
	count := len(routes)
	if count == 0 || count > 3 {
		log.Fatal("invalid routes array")
	}

	area = routes[0]
	if count >= 2 {
		controller = routes[1]
	}
	if count == 3 {
		action = routes[2]
	}

	user := host.GetUser(ctx, x.UserJsonSessionKey)

	// 判断请求是否允许访问
	if user != nil {
		if x.PermissionAuditor.CheckRouteWithLevel(area, controller, action, user.Role, user.Level) {
			// 有权限
			ctx.Next()
			return
		} else {
			// 没权限
			ctx.Redirect(x.AccessDeniedPath, http.StatusFound)
			return
		}
	}

	// 未登录
	allow := x.PermissionAuditor.CheckRouteWithLevel(area, controller, action, 0, 0)
	if allow {
		// 允许匿名
		ctx.Next()
		return
	}

	// 记录请求地址，跳转去登录页面
	host.RedirectAuthorizeEndpoint(ctx, x.OAuthOptions, ctx.RequestURL())
}

func (x *OAuthClientHost) getUserLock(userID string) *sync.RWMutex {
	if !x.UserLocks.Exists(userID) {
		x.UserLocks.Add(userID, time.Second*30, new(sync.RWMutex))
	}

	userLockCache, err := x.UserLocks.Value(userID)
	u.LogError(err)
	return userLockCache.Data().(*sync.RWMutex)
}