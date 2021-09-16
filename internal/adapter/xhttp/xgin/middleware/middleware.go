package middleware

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/response"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GinLogFormatter() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			return fmt.Sprintf(`address: %s, time: %s, method: %s, path: %s, errMessage: %s, proto: %s, code: %d, 
latency: %s, body: %v %v`, params.ClientIP, params.TimeStamp.Format("2006-01-02 15:04:05"),
				params.Method, params.Path, params.ErrorMessage, params.Request.Proto, params.StatusCode, params.Latency,
				params.Request.Body, "\n")
		})
}

func Tracer() gin.HandlerFunc {
	return ginhttp.Middleware(opentracing.GlobalTracer())
}

type userAuth struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Phone    string `json:"phone"`
	ID       string `json:"id"`
}

func loginResponse(c *gin.Context, i int, s string, t time.Time) {
	response.NewGinResponse(c).SetStatus(i).Success(s).SetMeta("expire", t.Format(time.RFC3339)).Response()
}

func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*userAuth); ok {
		res := jwt.MapClaims{}
		if err := mapstructure.Decode(v, &res); err != nil {
			logrus.Fatal(err)
		}
		return res
	}
	return map[string]interface{}{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	var (
		res = &userAuth{}
	)
	if err := mapstructure.Decode(claims, res); err != nil {
		logrus.Fatal(err)
	}
	return res
}

func authenticator(fc application.UserServiceInterface) func(c *gin.Context) (res interface{}, err error) {
	return func(c *gin.Context) (res interface{}, err error) {
		var (
			user *vo.UserVO
			f    vo.LoginForm
		)
		if err = c.ShouldBind(&f); err != nil {
			return
		}
		if config.GetConfig().Auth.IsCaptcha {
			verify := captcha.VerifyString(f.CaptchaId, f.Digits)
			if !verify {
				return res, errors.New("incorrect captcha")
			}
		}
		if user, err = fc.Get(c.Request.Context(), vo.UserUnique{Username: f.Username, Phone: f.Phone}); err != nil {
			return res, errors.Wrap(err, jwt.ErrFailedAuthentication.Error())
		}
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(f.Password)); err != nil {
			return res, errors.Wrap(err, jwt.ErrFailedAuthentication.Error())
		}
		fmt.Println("aaaaaaaaaaaaaaaa11111111111111111", user)
		res = &userAuth{
			Username: user.Username,
			UUID:     user.UUID,
			Phone:    user.Phone,
			ID:       user.ID,
		}
		return
	}
}

func authorizator(data interface{}, c *gin.Context) (res bool) {
	if val, ok := data.(*userAuth); ok {
		c.Set("id", val.ID)
		c.Set("uuid", val.UUID)
		return true
	}
	return false
}

func unauthorized(c *gin.Context, i int, s string) {
	response.NewGinResponse(c).SetStatus(i).Fail(i, errors.New(s)).Response()
}

func AuthJWT(fc application.UserServiceInterface) (middleware *jwt.GinJWTMiddleware, err error) {
	midd := &jwt.GinJWTMiddleware{}
	cfg := config.GetConfig()
	authCfg := cfg.Auth
	midd.Key = []byte(authCfg.Key)
	midd.Timeout = time.Duration(authCfg.Timeout) * time.Hour
	midd.MaxRefresh = time.Duration(authCfg.MaxRefresh) * time.Hour
	midd.IdentityKey = authCfg.IdentityKey
	midd.PayloadFunc = payload
	midd.IdentityHandler = identityHandler
	midd.LoginResponse = loginResponse
	midd.RefreshResponse = loginResponse
	midd.Authenticator = authenticator(fc)
	midd.Authorizator = authorizator
	midd.Unauthorized = unauthorized
	midd.TokenLookup = "header: Authorization, query: token, cookie: jwt"
	midd.TokenHeadName = "Bearer"
	midd.TimeFunc = time.Now
	return jwt.New(midd)
}
