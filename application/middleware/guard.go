package guard

import (
	"antrein/bc-dashboard/model/config"
	"antrein/bc-dashboard/model/dto"
	"antrein/bc-dashboard/model/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type GuardContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type AuthGuardContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Claims         entity.JWTClaim
}

func (g *GuardContext) ReturnError(status int, message string) error {
	g.ResponseWriter.WriteHeader(status)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  status,
		Message: message,
	})
}

func (g *GuardContext) ReturnSuccess(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func (g *GuardContext) ReturnCreated(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Created",
		Data:    data,
	})
}

func (g *GuardContext) ReturnEvent(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(g.ResponseWriter, "data: %s\n\n", jsonData)
	if err != nil {
		return err // Handle writing errors
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	} else {
		return fmt.Errorf("streaming unsupported")
	}

	return nil
}

func (g *AuthGuardContext) ReturnError(status int, message string) error {
	g.ResponseWriter.WriteHeader(status)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  status,
		Message: message,
	})
}

func (g *AuthGuardContext) ReturnSuccess(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func (g *AuthGuardContext) ReturnCreated(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Created",
		Data:    data,
	})
}

func (g *AuthGuardContext) ReturnEvent(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(g.ResponseWriter, "data: %s\n\n", jsonData)
	if err != nil {
		return err // Handle writing errors
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	} else {
		return fmt.Errorf("streaming unsupported")
	}

	return nil
}

func DefaultGuard(handlerFunc func(g *GuardContext) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guardCtx := GuardContext{
			ResponseWriter: w,
			Request:        r,
		}
		if err := handlerFunc(&guardCtx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func AuthGuard(cfg *config.Config, handlerFunc func(g *AuthGuardContext) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			http.Error(w, "Unauthorized - Invalid token format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.Secrets.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authGuardCtx := AuthGuardContext{
			ResponseWriter: w,
			Request:        r,
			Claims: entity.JWTClaim{
				UserID: userID,
			},
		}

		if err := handlerFunc(&authGuardCtx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func BodyParser(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func IsMethod(r *http.Request, method string) bool {
	return r.Method == method
}

func GetParam(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}
