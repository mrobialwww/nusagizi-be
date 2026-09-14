package middleware

import (
	"log/slog"
	"net/http"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/repository"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/gin-gonic/gin"
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

func GinMiddleware(jwtMiddleware *jwtmiddleware.JWTMiddleware, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var handled bool

		jwtMiddleware.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil validated claims dari context
			claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](r.Context())
			if err != nil {
				slog.Error("Failed to get claims", "error", err)
				return
			}

			auth0ID := claims.RegisteredClaims.Subject

			// Lazy provisioning
			ctx := r.Context()
			user, err := userRepo.GetUserBySub(ctx, auth0ID)

			// Jika user belum ada di table Users, lakukan lazy provisioning.
			// Email & username dibaca dari JWT custom claims yang ditambahkan via Auth0 Post Login Actions — tidak perlu hit /userinfo.
			if err != nil {
				customClaims, ok := claims.CustomClaims.(*auth.CustomClaims)
				if !ok {
					slog.Error("Failed to parse custom claims")
					return
				}

				email := customClaims.Email
				username := customClaims.Username

				// Fallback: jika username belum ada di custom claims
				// (user belum punya Auth0 Action untuk username), gunakan nickname dari sub (prefix sebelum '|')
				if username == "" {
					username = email
				}

				user, err = userRepo.CreateUserFromAuth0(ctx, auth0ID, email, username)
				if err != nil {
					slog.Error("Failed to create user", "error", err)
					return
				}
			}

			c.Set("user_id", user.ID)
			c.Set("auth0_id", auth0ID)
			c.Set("user", user)
			handled = true
		})).ServeHTTP(c.Writer, c.Request)

		if !handled {
			// Prevent sending 200 OK if the request is forcefully aborted
			// and ensure jwtMiddleware hasn't written a response (e.g., 401 Unauthorized)
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user authentication. Please try again."})
			}
			c.Abort()
			return
		}
		c.Next()
	}
}
