package api

import (
	"conf"
	"fmt"
	"net/http"
	"net/url"
)

const (
	isFromWeixinParam = "__is_from_weixin__"
	isFromWeixinValue = "true"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	finalUrl := ""

	isFromWeixin := r.Form[isFromWeixinParam]
	if len(isFromWeixin) > 0 && contains(isFromWeixin, isFromWeixinValue) {
		if key := r.Form[conf.Conf.KeyParam]; len(key) > 0 {
			if redirectUrl := conf.Conf.RedirectUrls[key[0]]; len(redirectUrl) > 0 {

				u, _ := url.Parse("")

				q := r.URL.Query()

				q.Del(isFromWeixinParam)
				q.Del(conf.Conf.KeyParam)
				for key, value := range u.Query() {
					q.Add(key, value[0])
				}

				u.RawQuery = q.Encode()

				finalUrl = redirectUrl + u.String()

			}
		}
	} else {
		state := "STATE"
		if passedState := r.Form["state"]; len(passedState) > 0 {
			state = passedState[0]
		}

		scope := "snsapi_base"
		if passedScope := r.Form["scope"]; len(passedScope) > 0 {
			scope = passedScope[0]
		}

		passedQuery := r.URL.Query()
		passedQuery.Del("state")
		passedQuery.Add(isFromWeixinParam, isFromWeixinValue)
		r.URL.RawQuery = passedQuery.Encode()
		redirectUri := conf.Conf.Host + r.URL.String()

		authUrl, _ := url.Parse("https://open.weixin.qq.com/connect/oauth2/authorize")

		authQuery := authUrl.Query()

		authQuery.Set("appid", conf.Conf.AppId)
		authQuery.Set("redirect_uri", redirectUri)
		authQuery.Set("response_type", "code")
		authQuery.Set("scope", scope)
		authQuery.Set("state", state)

		authUrl.RawQuery = authQuery.Encode()
		authUrl.Fragment = "wechat_redirect"

		finalUrl = authUrl.String()
	}

	if len(finalUrl) > 0 {
		fmt.Println(finalUrl)
		http.Redirect(w, r, finalUrl, http.StatusFound)
	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
