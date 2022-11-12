package server

import (
	"context"
	v1 "kratos-shop/api/helloworld/v1"
	"kratos-shop/internal/conf"
	"kratos-shop/internal/service"
	"kratos-shop/pkg/trace"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(

	c *conf.Server,
	greeter *service.GreeterService,
	studentSvc *service.StudentService,
	logger log.Logger,
) *http.Server {

	tp, err := trace.InitTracer("http://host.docker.internal:14268/api/traces")
	if err != nil {
		log.Error(err)
	}

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(
				tracing.WithTracerProvider(tp), // 接入 jaeger
			),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
		),
		// http.Filter(handlers.CORS( // 浏览器跨域
		// 	handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		// 	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		// 	handlers.AllowedOrigins([]string{"*"}),
		// )),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterStudentHTTPServer(srv, studentSvc)

	tr := tp.Tracer("component-main")
	ctx, span := tr.Start(nil, "Foo")
	defer span.End()
	bar(ctx)
	return srv
}

func bar(ctx context.Context) {
	// Use the global TracerProvider.
	tr := otel.Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	// Do bar...
}

// NewWhiteListMatcher 设置白名单，不需要 token 验证的接口
/**
 http.Middleware(
            selector.Server( // jwt 验证
                jwt.Server(func(token *jwt2.Token) (interface{}, error) {
                    return []byte(ac.JwtKey), nil
                }, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
            ).Match(NewWhiteListMatcher()).Build(),

 ),
**/
// func NewWhiteListMatcher() selector.MatchFunc {
//     whiteList := make(map[string]struct{})
//     whiteList["/shop.shop.v1.Shop/Captcha"] = struct{}{}
//     whiteList["/shop.shop.v1.Shop/Login"] = struct{}{}
//     whiteList["/shop.shop.v1.Shop/Register"] = struct{}{}
//     return func(ctx context.Context, operation string) bool {
//         if _, ok := whiteList[operation]; ok {
//             return false
//         }
//         return true
//     }
// }
