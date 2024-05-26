package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func mutateResponse(
	logr zerolog.Logger,
) func(context.Context, http.ResponseWriter, protoreflect.ProtoMessage) error {
	logr.Info().Str("middleware", "mutateResponse").Msg(REGISTER)
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
