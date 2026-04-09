package middleware

import (
	"log/slog"
	"net/http"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/repository"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewMiddleware(jwtValidator *validator.Validator) (*jwtmiddleware.JWTMiddleware, error) {
    return jwtmiddleware.New(
        jwtmiddleware.WithValidator(jwtValidator),
        jwtmiddleware.WithValidateOnOptions(false),
        jwtmiddleware.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
            slog.Error("JWT validation failed", "error", err, "path", r.URL.Path)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"message":"Failed to validate JWT."}`))
        }),
    )
}

func GinMiddleware(jwtMiddleware *jwtmiddleware.JWTMiddleware, pool *pgxpool.Pool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var handled bool

        jwtMiddleware.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Ambil validated claims dari context
            claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](r.Context())
            if err != nil {
                slog.Error("Failed to get claims", "error", err)
                return
            }

            userID := claims.RegisteredClaims.Subject

            // Belum ada → ambil email dari JWT(CustomClaims) lalu insert
            email := ""
            if customClaims, ok := claims.CustomClaims.(*auth.CustomClaims); ok {
                email = customClaims.Email
            }
            
            // Lazy provisioning
            user, err := repository.GetUserBySub(pool, userID)
            if err != nil {
                
                // Email ada di RegisteredClaims jika scope "email" diminta atau ambil dari token langsung
                user, err = repository.CreateUserFromAuth0(pool, userID, email)
                if err != nil {
                    slog.Error("Failed to create user", "error", err)
                    return
                }
            }

            c.Set("user_id", userID)
            c.Set("user", user)
            handled = true
        })).ServeHTTP(c.Writer, c.Request)

        if !handled {
            c.Abort()
            return
        }
        c.Next()
    }
}
