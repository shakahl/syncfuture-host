package resource

import (
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	oauth2core "github.com/Lukiya/oauth2go/core"
	"github.com/Lukiya/oauth2go/model"
	"github.com/pascaldekloe/jwt"
	"github.com/syncfuture/go/sconv"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/srsautil"
	"github.com/syncfuture/go/sslice"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/host"
)

type OAuthResourceHost struct {
	host.BaseHost
	OAuthOptions     *model.Resource `json:"OAuth,omitempty"`
	PublicKeyPath    string
	SigningAlgorithm string
	PublicKey        *rsa.PublicKey
	// TokenValidator   func(*jwt.Claims) string
}

func (x *OAuthResourceHost) BuildOAuthResourceHost() {
	x.BaseHost.BuildBaseHost()

	if x.OAuthOptions == nil {
		log.Fatal("OAuth secion in configuration is missing")
	}

	if x.PublicKeyPath == "" {
		log.Fatal("public key path cannot be empty")
	}
	if x.OAuthOptions == nil {
		log.Fatal("oauth options cannot be nil")
	}
	if x.OAuthOptions.ValidIssuers == nil || len(x.OAuthOptions.ValidIssuers) == 0 {
		log.Fatal("Issuers cannot be empty")
	}
	if x.OAuthOptions.ValidAudiences == nil || len(x.OAuthOptions.ValidAudiences) == 0 {
		log.Fatal("Audiences cannot be empty")
	}
	if x.SigningAlgorithm == "" {
		x.SigningAlgorithm = jwt.PS256
	}

	if x.URLProvider != nil {
		for i := range x.OAuthOptions.ValidIssuers {
			x.OAuthOptions.ValidIssuers[i] = x.URLProvider.RenderURL(x.OAuthOptions.ValidIssuers[i])
		}
		for i := range x.OAuthOptions.ValidAudiences {
			x.OAuthOptions.ValidAudiences[i] = x.URLProvider.RenderURL(x.OAuthOptions.ValidAudiences[i])
		}
	}

	// read public certificate
	cert, err := srsautil.ReadCertFromFile(x.PublicKeyPath)
	u.LogFaltal(err)
	x.PublicKey = cert.PublicKey.(*rsa.PublicKey)
}

func (x *OAuthResourceHost) AuthHandler(ctx host.IHttpContext) {
	authHeader := ctx.GetHeader(host.Header_Auth)
	if authHeader == "" {
		ctx.SetStatusCode(http.StatusUnauthorized)
		ctx.WriteString("Authorization header is missing")
		return
	}

	// verify authorization header
	array := strings.Split(authHeader, " ")
	if len(array) != 2 || array[0] != host.AuthType_Bearer {
		ctx.SetStatusCode(http.StatusBadRequest)
		log.Warnf("'%s'invalid authorization header format. '%s'", ctx.GetRemoteIP(), authHeader)
		return
	}

	// verify signature
	token, err := jwt.RSACheck([]byte(array[1]), x.PublicKey)
	if err != nil {
		ctx.SetStatusCode(http.StatusUnauthorized)
		log.Warn("'"+ctx.GetRemoteIP()+"'", err)
		return
	}

	// validate time limits
	isNotExpired := token.Valid(time.Now().UTC())
	if !isNotExpired {
		ctx.SetStatusCode(http.StatusUnauthorized)
		msgCode := "current time not in token's valid period"
		ctx.WriteString(msgCode)
		log.Warn("'"+ctx.GetRemoteIP()+"'", msgCode)
		return
	}

	// validate aud
	isValidAudience := x.OAuthOptions.ValidAudiences != nil && sslice.HasAnyStr(x.OAuthOptions.ValidAudiences, token.Audiences)
	if !isValidAudience {
		ctx.SetStatusCode(http.StatusUnauthorized)
		msgCode := "invalid audience"
		ctx.WriteString(msgCode)
		log.Warn("'"+ctx.GetRemoteIP()+"'", msgCode)
		return
	}

	// validate iss
	isValidIssuer := x.OAuthOptions.ValidIssuers != nil && sslice.HasStr(x.OAuthOptions.ValidIssuers, token.Issuer)
	if !isValidIssuer {
		ctx.SetStatusCode(http.StatusUnauthorized)
		msgCode := "invalid issuer"
		ctx.WriteString(msgCode)
		log.Warn("'"+ctx.GetRemoteIP()+"'", msgCode)
		return
	}

	// if x.TokenValidator != nil {
	// 	if msgCode := x.TokenValidator(token); msgCode != "" {
	// 		ctx.SetStatusCode(http.StatusUnauthorized)
	// 		ctx.WriteString(msgCode)
	// 		log.Warn("'"+ctx.GetRemoteIP()+"'", msgCode)
	// 		return
	// 	}
	// }

	var msgCode string
	if token != nil {
		routeKey := ctx.GetItemString(host.Item_RouteKey)
		area, controller, action := host.GetRoutesByKey(routeKey)

		roles := sconv.ToInt64(token.Set[oauth2core.Claim_Role])
		level := sconv.ToInt64(token.Set[oauth2core.Claim_Level])
		if x.PermissionAuditor.CheckRouteWithLevel(area, controller, action, roles, int32(level)) {
			// Has permission, allow
			ctx.SetItem(host.Item_JWT, token)
			ctx.Next()
			return
		} else {
			msgCode = "permission denied"
		}
	}

	// Not allow
	ctx.SetStatusCode(http.StatusUnauthorized)
	ctx.WriteString(msgCode)
}