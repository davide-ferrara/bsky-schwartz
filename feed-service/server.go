package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bsky-schwartz/internal/bsky"
	"bsky-schwartz/internal/models"
	"bsky-schwartz/views"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewServer(cfg *Config, logger *slog.Logger) http.Handler {
	// Crea mux per feed generator (http.ServeMux)
	feedMux := http.NewServeMux()

	// Feed generator endpoints
	feedMux.HandleFunc("GET /.well-known/did.json", handleWellKnown(cfg))
	feedMux.HandleFunc("GET /xrpc/_health", handleHealth)
	feedMux.HandleFunc("GET /xrpc/app.bsky.feed.describeFeedGenerator", handleDescribeFeedGenerator(cfg))
	feedMux.HandleFunc("GET /xrpc/app.bsky.feed.getFeedSkeleton", handleGetFeedSkeleton(cfg, logger))

	// Crea router Gin per web app
	gin.SetMode(gin.ReleaseMode)
	webRouter := gin.New()
	webRouter.Use(gin.Recovery())

	// Session middleware
	sessionSecret := os.Getenv("SESSION_KEY")
	if sessionSecret == "" {
		sessionSecret = "default-secret-change-in-production"
	}
	store := cookie.NewStore([]byte(sessionSecret))
	webRouter.Use(sessions.Sessions("bluesky-session", store))

	// Static files
	webRouter.Static("/static", "./static")
	webRouter.Static("/lexicons", "./lexicons")

	// Web app routes
	webRouter.GET("/", RootHandler)
	webRouter.GET("/login", LoginGetHandler)
	webRouter.POST("/login", LoginPostHandler)
	webRouter.POST("/logout", LogoutHandler)
	webRouter.GET("/delete", DeleteWeightsHanlder)
	webRouter.GET("/values", ValuesHandler)
	webRouter.POST("/preferences", PreferencesHandler)

	// Feed generator ha priorità sulle route specifiche
	// Web app gestisce tutto il resto
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route specifiche del feed generator
		if r.URL.Path == "/.well-known/did.json" ||
			r.URL.Path == "/xrpc/_health" ||
			r.URL.Path == "/xrpc/app.bsky.feed.describeFeedGenerator" ||
			r.URL.Path == "/xrpc/app.bsky.feed.getFeedSkeleton" {
			feedMux.ServeHTTP(w, r)
			return
		}

		// Tutto il resto va alla web app
		webRouter.ServeHTTP(w, r)
	})
}

// ==== Handler per Web App (Gin) =====

func GetSessionString(session sessions.Session, name string) string {
	if val := session.Get(name); val != nil {
		return val.(string)
	}
	return ""
}

func GetUser(session sessions.Session) (User, error) {
	var user User
	handle := GetSessionString(session, "handle")
	appPassword := GetSessionString(session, "appPassword")
	// fmt.Println("DEBUG: ", handle, appPassword)

	if handle == "" || appPassword == "" {
		return user, fmt.Errorf("invalid handle or appPassword")
	}
	user.Handle = handle
	user.AppPassword = appPassword
	return user, nil
}

func RootHandler(c *gin.Context) {
	session := sessions.Default(c)

	user, err := GetUser(session)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	client, err := bsky.NewClient(user.Handle, user.AppPassword)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	ctx := context.Background()
	username, err := bsky.GetProfileInfo(ctx, &client, user.Handle)
	if err != nil {
		fmt.Println(err)
	}

	// Get weights from shared memory (using user DID as key)
	userDID := GetSessionString(session, "userDID")
	weights := GlobalUserWeights.Get(userDID)
	if weights == nil {
		// Fallback: try to get from session (for backwards compatibility)
		weights = make(map[string]float64)
		if w := session.Get("weights"); w != nil {
			if data, ok := w.([]byte); ok {
				if err := json.Unmarshal(data, &weights); err != nil {
					fmt.Printf("Could not Marshal: %v", err)
				}
			}
		}
	}

	// Get flash message
	var message string
	var messageType string
	if msg := session.Get("message"); msg != nil {
		message = msg.(string)
		session.Delete("message")
	}
	if mt := session.Get("messageType"); mt != nil {
		messageType = mt.(string)
		session.Delete("messageType")
	}
	if err := session.Save(); err != nil {
		fmt.Printf("Could not save session, trying clearing it: %v", err)
		session.Clear()
		return
	}

	component := views.SlidersPage(username, messageType, message, weights)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		panic(err)
	}
}

func LoginGetHandler(c *gin.Context) {
	session := sessions.Default(c)

	var message string
	var messageType string
	if msg := session.Get("message"); msg != nil {
		message = msg.(string)
		session.Delete("message")
	}
	if mt := session.Get("messageType"); mt != nil {
		messageType = mt.(string)
		session.Delete("messageType")
	}

	if err := session.Save(); err != nil {
		fmt.Printf("Could not save session, trying clearing it: %v", err)
		session.Clear()
		return
	}

	component := views.LoginPage(messageType, message)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		fmt.Println(err)
	}
}

func LoginPostHandler(c *gin.Context) {
	handle := c.PostForm("handle")
	appPassword := c.PostForm("appPassword")

	session := sessions.Default(c)

	if handle == "" || appPassword == "" {
		session.Set("message", "Handle e App Password sono obbligatori")
		session.Set("messageType", "error")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}

	client, err := bsky.NewClient(handle, appPassword)
	if err != nil {
		fmt.Println("ERRORE LOGIN - bsky.NewClient:", err)
		session.Set("message", "Handle o App Password non validi")
		session.Set("messageType", "error")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("LOGIN OK - Handle:", handle)

	session.Set("handle", handle)
	session.Set("appPassword", appPassword)
	session.Save()

	ctx := context.Background()
	username, err := bsky.GetProfileInfo(ctx, &client, handle)
	if err != nil {
		fmt.Println("WARNING - GetProfileInfo:", err)
	} else {
		fmt.Println("USERNAME:", username)
		session.Set("username", username)
	}

	// Get user DID from profile (for weights storage)
	// For now we use handle as DID
	userDID, err := bsky.ResolveDID(ctx, &client, handle)
	if err != nil {
		fmt.Println("Could not resolve handle, using handle as key")
		userDID = handle
	}

	session.Set("userDID", userDID)

	// Read weights from PDS
	weights, err := bsky.GetWeights(ctx, &client, handle)
	if err != nil {
		fmt.Println("ℹ️ INFO - Nessun weights nel PDS, uso default:", err)
		weights = make(map[string]float64)
		for _, v := range models.SwartzValues {
			weights[v.ID] = 0.0
		}
	} else {
		fmt.Println("WEIGHTS caricati dal PDS")
	}

	// Save weights in shared memory
	GlobalUserWeights.Set(userDID, weights)

	// Also save in session for backwards compatibility
	if data, err := json.Marshal(weights); err == nil {
		session.Set("weights", data)
	}

	if err := session.Save(); err != nil {
		fmt.Println("ERRORE salvataggio sessione:", err)
	}

	c.Redirect(http.StatusFound, "/")
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("handle")
	session.Delete("appPassword")
	session.Delete("username")
	session.Delete("userDID")
	session.Delete("weights")

	if err := session.Save(); err != nil {
		fmt.Println(err)
	}

	c.Redirect(http.StatusFound, "/login")
}

func ValuesHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := GetSessionString(session, "username")
	component := views.ValuesPage(username)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		panic(err)
	}
}

// Weights webpage
func PreferencesHandler(c *gin.Context) {
	session := sessions.Default(c)

	user, err := GetUser(session)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Parse form data into map
	weights := make(map[string]float64)
	for _, value := range models.SwartzValues {
		if val := c.PostForm(value.ID); val != "" {
			var f float64
			if _, err := fmt.Sscanf(val, "%f", &f); err == nil {
				weights[value.ID] = f
			}
		}
	}

	// Create client
	client, err := bsky.NewClient(user.Handle, user.AppPassword)
	if err != nil {
		fmt.Println("ERRORE - Creazione client fallita:", err)
		session.Set("message", "Errore di connessione a Bluesky")
		session.Set("messageType", "error")
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Save weights to PDS
	ctx := context.Background()
	if err := bsky.SaveWeights(ctx, &client, user.Handle, weights); err != nil {
		fmt.Println("ERRORE - Salvataggio PDS fallito:", err)
		session.Set("message", "Errore nel salvataggio delle preferenze nei server Bluesky")
		session.Set("messageType", "error")
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	fmt.Println("WEIGHTS salvati nel PDS")
	session.Set("message", "Preferenze salvate con successo!")
	session.Set("messageType", "success")

	// Get user DID
	userDID := GetSessionString(session, "userDID")

	// Save weights in shared memory
	GlobalUserWeights.Set(userDID, weights)

	// Also save in session for backwards compatibility
	if data, err := json.Marshal(weights); err == nil {
		session.Set("weights", data)
	}
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func DeleteWeightsHanlder(c *gin.Context) {
	session := sessions.Default(c)

	user, err := GetUser(session)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	client, err := bsky.NewClient(user.Handle, user.AppPassword)
	if err != nil {
		fmt.Println("ERRORE LOGIN - bsky.NewClient:", err)
		session.Set("message", "Handle o App Password non validi")
		session.Set("messageType", "error")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}

	ctx := context.Background()
	err = bsky.DeleteWeights(ctx, &client, user.Handle)
	if err != nil {
		session.Set("message", "Errore nell'eliminare il record")
		session.Set("messageType", "error")
		session.Save()
		c.Redirect(http.StatusPermanentRedirect, "/profile")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/preferences")
	return
}

// ==== Feed Generator Handlers ====

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleWellKnown(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(cfg.ServiceDID, cfg.Hostname) {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, map[string]any{
			"@context": []string{"https://www.w3.org/ns/did/v1"},
			"id":       cfg.ServiceDID,
			"service": []map[string]any{
				{
					"id":              "#bsky_fg",
					"type":            "BskyFeedGenerator",
					"serviceEndpoint": fmt.Sprintf("https://%s", cfg.Hostname),
				},
			},
		})
	}
}

func handleDescribeFeedGenerator(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		feeds := make([]map[string]string, 0, len(Algos))
		for shortname := range Algos {
			feeds = append(feeds, map[string]string{
				"uri": fmt.Sprintf("at://%s/app.bsky.feed.generator/%s", cfg.PublisherDID, shortname),
			})
		}
		writeJSON(w, map[string]any{
			"did":   cfg.ServiceDID,
			"feeds": feeds,
		})
	}
}

func handleGetFeedSkeleton(cfg *Config, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userDID string

		// TODO: IMPROVE JWT PASRING
		if auth := r.Header.Get("Authorization"); auth != "" {
			fmt.Println(auth)

			parts := strings.Split(auth, " ")
			tokenStr := parts[1]

			parser := jwt.NewParser()
			token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
			if err != nil {
				logger.Warn("Error in parsing JWT")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				logger.Warn("Error in claming JWT")
			}

			// Il DID è nel campo "iss" (issuer) del token
			userDID, ok = claims["iss"].(string)
			if !ok {
				logger.Warn("Unauthenticated feed reqeust")
			}
		}

		fmt.Println("REQEUST FROM:", userDID)

		feedParam := r.URL.Query().Get("feed")
		if feedParam == "" {
			writeError(w, http.StatusBadRequest, "missing feed parameter")
			return
		}

		parts := strings.Split(strings.TrimPrefix(feedParam, "at://"), "/")
		if len(parts) != 3 {
			writeError(w, http.StatusBadRequest, "invalid feed URI")
			return
		}
		did, collection, rkey := parts[0], parts[1], parts[2]

		if did != cfg.PublisherDID || collection != "app.bsky.feed.generator" {
			writeError(w, http.StatusBadRequest, "UnsupportedAlgorithm")
			return
		}

		algo, ok := Algos[rkey]
		if !ok {
			writeError(w, http.StatusBadRequest, "UnsupportedAlgorithm")
			return
		}

		limit := 50
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = min(parsed, 200)
			}
		}
		cursor := r.URL.Query().Get("cursor")

		// Chiamata all'algoritmo di generazione del feed
		result, err := algo(r.Context(), limit, cursor, userDID)
		if err != nil {
			logger.Error("algo error", "algo", rkey, "err", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		writeJSON(w, result)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
