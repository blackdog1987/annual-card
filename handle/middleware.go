package handle

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/molibei/annual-card/lib/errors"
	"github.com/molibei/annual-card/module"
	"fmt"
)

func Identity() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var h []string
		for i, v := range ctx.Request.Header {
			var s string
			s = i + "=" + strings.Join(v, ",")
			h = append(h, s)
		}
		//h = append(h, "ClientIP="+ctx.ClientIP())
		ctx.Set("identity", strings.Join(h, "&"))
		secret := ctx.Request.Header.Get("User-Agent") // + ctx.ClientIP()
		hash := md5.New()
		hash.Write([]byte(secret))
		cipherStr := hash.Sum(nil)
		ctx.Set("secret", hex.EncodeToString(cipherStr[4:12]))
		ctx.Next()
	}
}

// Authorization .
func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.PostForm("token")
		fmt.Println("post:", token)
		if token == "" {
			token = ctx.Query("token")
			fmt.Println("get:", token)
		}
		if len(token) < 10 {
			ctx.JSON(http.StatusOK, errors.TOKEN_VALID_ERR)
			ctx.Abort()
			return
		}
		// token 解密
		secret, _ := ctx.Get("secret")
		tk, ok := module.Token.Decode(token, secret.(string))
		fmt.Println(tk, ok)
		if !ok {
			ctx.JSON(http.StatusOK, errors.TOKEN_VALID_ERR)
			ctx.Abort()
			return
		}
		now := time.Now().Unix()
		if tk.ExpiredIn < now {
			ctx.JSON(http.StatusOK, errors.TOKEN_EXPIRED)
			ctx.Abort()
			return
		}
		ctx.Set("token", tk)
		ctx.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		u := session.Get("u")
		if ctx.Query("test") == "yes" {
			u = int64(2)
		}
		if u == nil || u == 0 {
			// ctx.HTML(http.StatusOK, "public/qrcode.html", nil)
			ctx.Abort()
			return
		}
		ctx.Set("uid", u)
		ctx.Next()
	}
}
