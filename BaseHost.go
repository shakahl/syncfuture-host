package host

import (
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/ssecurity"
	"github.com/syncfuture/go/surl"
)

type BaseHost struct {
	// ListenAddr         string
	Debug              bool
	Name               string
	URIKey             string
	RouteKey           string
	PermissionKey      string
	RedisConfig        *sredis.RedisConfig `json:"Redis,omitempty"`
	ConfigProvider     sconfig.IConfigProvider
	URLProvider        surl.IURLProvider
	PermissionProvider ssecurity.IPermissionProvider
	RouteProvider      ssecurity.IRouteProvider
	PermissionAuditor  ssecurity.IPermissionAuditor
}

func (x *BaseHost) BuildBaseHost() {
	// if r.Name == "" {
	// 	log.Fatal("Name cannot be empty")
	// }
	// if r.ListenAddr == "" {
	// 	log.Fatal("ListenAddr cannot be empty")
	// }

	if x.ConfigProvider == nil {
		x.ConfigProvider = sconfig.NewJsonConfigProvider()
	}

	if x.URLProvider == nil && x.URIKey != "" && x.RedisConfig != nil {
		x.URLProvider = surl.NewRedisURLProvider(x.URIKey, x.RedisConfig)
	}

	if x.PermissionProvider == nil && x.PermissionKey != "" && x.RedisConfig != nil {
		x.PermissionProvider = ssecurity.NewRedisPermissionProvider(x.PermissionKey, x.RedisConfig)
	}

	if x.RouteProvider == nil && x.RouteKey != "" && x.RedisConfig != nil {
		x.RouteProvider = ssecurity.NewRedisRouteProvider(x.RouteKey, x.RedisConfig)
	}

	if x.PermissionAuditor == nil && x.PermissionProvider != nil { // RouteProvider 允许为空
		x.PermissionAuditor = ssecurity.NewPermissionAuditor(x.PermissionProvider, x.RouteProvider)
	}

	log.Init(x.ConfigProvider)
	ConfigHttpClient(x.ConfigProvider)

	return
}

func (x BaseHost) GetDebug() bool {
	return x.Debug
}

func (x BaseHost) GetConfigProvider() sconfig.IConfigProvider {
	return x.ConfigProvider
}
func (x BaseHost) GetRedisConfig() *sredis.RedisConfig {
	return x.RedisConfig
}
func (x BaseHost) GetURLProvider() surl.IURLProvider {
	return x.URLProvider
}
func (x BaseHost) GetPermissionAuditor() ssecurity.IPermissionAuditor {
	return x.PermissionAuditor
}
func (x BaseHost) GetPermissionProvider() ssecurity.IPermissionProvider {
	return x.PermissionProvider
}
func (x BaseHost) GetRouteProvider() ssecurity.IRouteProvider {
	return x.RouteProvider
}

type BaseWebHost struct {
	// BaseHost
	ListenAddr        string
	CORS              *CORSOptions
	GlobalPreHandlers []RequestHandler
	GlobalSufHandlers []RequestHandler
	Actions           map[string]*Action
}

func (x *BaseWebHost) BuildBaseWebHost() {
	if x.ListenAddr == "" {
		log.Fatal("ListenAddr cannot be empty")
	}

	x.Actions = make(map[string]*Action)
}

// AddGlobalPreHandlers 添加全局前置中间件, toTail: 是否添加在已有全局前置中间件的尾部
func (x *BaseWebHost) AddGlobalPreHandlers(toTail bool, handlers ...RequestHandler) {
	if toTail {
		x.GlobalPreHandlers = append(x.GlobalPreHandlers, handlers...)
	} else {
		x.GlobalPreHandlers = append(handlers, x.GlobalPreHandlers...)
	}
}

// AppendGlobalSufHandlers 添加全局后置中间件, toTail: 是否添加在已有全局后置中间件的尾部
func (x *BaseWebHost) AppendGlobalSufHandlers(toTail bool, handlers ...RequestHandler) {
	if toTail {
		x.GlobalSufHandlers = append(x.GlobalSufHandlers, handlers...)
	} else {
		x.GlobalSufHandlers = append(handlers, x.GlobalSufHandlers...)
	}
}

func (x *BaseWebHost) AddActionGroups(actionGroups ...*ActionGroup) {
	////////// 添加Actions
	for _, actionGroup := range actionGroups {
		for _, action := range actionGroup.Actions {
			// 添加预先执行中间件
			if len(actionGroup.PreHandlers) > 0 {
				action.Handlers = append(actionGroup.PreHandlers, action.Handlers...)
			}
			// 添加后执行中间件
			if len(actionGroup.AfterHandlers) > 0 {
				action.Handlers = append(action.Handlers, actionGroup.AfterHandlers...)
			}

			_, ok := x.Actions[action.Route]
			if ok {
				log.Fatal("duplicated route found: " + action.Route)
			}
			x.Actions[action.Route] = action
		}
	}
}

func (x *BaseWebHost) AddActions(actions ...*Action) {
	////////// 添加Actions
	for _, action := range actions {
		_, ok := x.Actions[action.Route]
		if ok {
			log.Fatal("duplicated route found: " + action.Route)
		}
		x.Actions[action.Route] = action
	}
}

func (x *BaseWebHost) AddAction(route, routeKey string, handlers ...RequestHandler) {
	////////// 添加Action
	action := NewAction(route, routeKey, handlers...)
	_, ok := x.Actions[action.Route]
	if ok {
		log.Fatal("duplicated route found: " + action.Route)
	}
	x.Actions[action.Route] = action
}
