package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func mutateResponse() func(context.Context, http.ResponseWriter, protoreflect.ProtoMessage) error {
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

		return nil
	}
}
