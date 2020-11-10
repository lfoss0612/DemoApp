package middleware

imports (
	"github.com/lfoss0612/DemoApp/server"
)

// Middleware getter
func Middlewares() []server.MiddlewareFunc {
	return []server.MiddlewareFunc{
		getContextHandler,
		getRecoveryHandler,				
		getLoggingHandler,		
		getTimeoutHandler,
	}
}