package handlers

import (
	"context"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/tsawler/celeritas"
	"github.com/tsawler/celeritas/mailer"
	"github.com/tsawler/celeritas/render"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var cel celeritas.Celeritas
var testSession *scs.SessionManager
var testHandlers Handlers

func TestMain(m *testing.M) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader("../views"),
		jet.InDevelopmentMode(),
	)

	myRenderer := render.Render{
		Renderer: "jet",
		RootPath: "../",
		Port:     "4000",
		JetViews: views,
		Session:  testSession,
	}

	cel = celeritas.Celeritas{
		AppName:       "myapp",
		Debug:         true,
		Version:       "1.0.0",
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		RootPath:      "../",
		Routes:        nil,
		Render:        &myRenderer,
		Session:       testSession,
		DB:            celeritas.Database{},
		JetViews:      views,
		EncryptionKey: cel.RandomString(32),
		Cache:         nil,
		Scheduler:     nil,
		Mail:          mailer.Mail{},
		Server:        celeritas.Server{},
	}

	testHandlers.App = &cel

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cel.SessionLoad)
	mux.Get("/", testHandlers.Home)
	fileServer := http.FileServer(http.Dir("./../public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))
	return mux
}

func getCtx(req *http.Request) context.Context {
	ctx, err := testSession.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
