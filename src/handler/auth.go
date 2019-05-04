package handler

import (
	"cloud-storage/src/common"
	"cloud-storage/src/util"
	"net/http"
)

func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			// Verify the token is valid or not
			if len(username) < 3 || !IsTokenValid(token) {
				// Token is invalid, return error response message
				resp := util.NewRespMsg(
					int(common.StatusInvalidToken),
					"Invalid Token",
					nil,
				)
				w.Write(resp.JSONBytes())
				return
			}
			h(w, r)
		})
}
