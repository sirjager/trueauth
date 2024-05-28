package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sirjager/trueauth/config"
)

func customHeaderMatcher(h []string) func(s string) (string, bool) {
	return func(key string) (string, bool) {
		header := strings.ToLower(key)
		for _, v := range h {
			myheader := strings.ToLower(v[1:])
			if myheader == header {

				if strings.HasPrefix(v, "+") {
					return myheader, true
				}
				if strings.HasPrefix(v, "-") {
					return "", false
				}

			}
		}
		return runtime.DefaultHeaderMatcher(key)
	}
}

func addCustomHeaders(
	config config.Config,
) func(context.Context, http.ResponseWriter, protoreflect.ProtoMessage) error {
	return func(ctx context.Context, w http.ResponseWriter, m protoreflect.ProtoMessage) error {
		if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
			for k, v := range md.HeaderMD {
				if strings.HasPrefix(k, "set-cookie:") {
					cookie := http.Cookie{}
					if err := json.Unmarshal([]byte(v[0]), &cookie); err != nil {
						return err
					}
					http.SetCookie(w, &cookie)
				} else {
					w.Header().Set(k, strings.Join(v, ","))
				}
			}
		}

		// adding server name in header to determine which server sent it
		w.Header().Set("x-server-name", config.Server.ServerName)

		return nil
	}
}
