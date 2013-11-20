package api

import (
	graph "atlantis/manager/api/graph"
	"atlantis/manager/crypto"
	"atlantis/manager/rpc"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cespare/go-apachelog"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	staticDir = "static"
)

var (
	notFoundHTML    = "404 not found"
	serverErrorHTML = "500 internal server error"
	server          *http.Server
	lAddr           = ""
	manager         = new(rpc.Manager)
	HandlerFunc     = func(h http.Handler) http.Handler {
		return h
	}
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, notFoundHTML)
}

func serverError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, serverErrorHTML)
}

func Init(listenAddr string) error {
	gmux := mux.NewRouter() // Use gorilla mux for APIs to make things easier

	gmux.NotFoundHandler = http.HandlerFunc(NotFound)
	// APIs should go here
	gmux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/img/favicon.ico")
	})
	gmux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, staticDir+"/dashboard", 302)
	})

	// Login
	gmux.HandleFunc("/login", Login).Methods("POST")

	// Task Management
	gmux.HandleFunc("/wait", Wait).Methods("GET")

	// Manager Management
	gmux.HandleFunc("/health", Health).Methods("GET")
	gmux.HandleFunc("/managers", ListManagers).Methods("GET")
	gmux.HandleFunc("/managers/{Region}/{Host}", RegisterManager).Methods("PUT")
	gmux.HandleFunc("/managers/{Region}/{Host}", UnregisterManager).Methods("DELETE")

	// Supervisor Management
	gmux.HandleFunc("/supervisors", ListSupervisors).Methods("GET")
	gmux.HandleFunc("/supervisors/{Host}", RegisterSupervisor).Methods("PUT")
	gmux.HandleFunc("/supervisors/{Host}", UnregisterSupervisor).Methods("DELETE")

	// Router Management
	gmux.HandleFunc("/routers", ListRouters).Methods("GET")
	gmux.HandleFunc("/routers/{Zone}/{IP}", GetRouter).Methods("GET")
	gmux.HandleFunc("/routers/{Zone}/{IP}", RegisterRouter).Methods("PUT")
	gmux.HandleFunc("/routers/{Zone}/{IP}", UnregisterRouter).Methods("DELETE")

	// App Management
	gmux.HandleFunc("/apps", ListRegisteredApps).Methods("GET")
	gmux.HandleFunc("/apps/{App}", GetApp).Methods("GET")
	gmux.HandleFunc("/apps/{App}", RegisterApp).Methods("PUT")
	gmux.HandleFunc("/apps/{App}", UnregisterApp).Methods("DELETE")

	// Container Health
	gmux.HandleFunc("/healthz", ContainerHealthzGet).Methods("GET")

	// Router Config Management
	gmux.HandleFunc("/pools", ListPools).Methods("GET")
	gmux.HandleFunc("/pools/{PoolName}", GetPool).Methods("GET")
	gmux.HandleFunc("/pools/{PoolName}", UpdatePool).Methods("PUT")
	gmux.HandleFunc("/pools/{PoolName}", DeletePool).Methods("DELETE")
	gmux.HandleFunc("/rules", ListRules).Methods("GET")
	gmux.HandleFunc("/rules/{RuleName}", GetRule).Methods("GET")
	gmux.HandleFunc("/rules/{RuleName}", UpdateRule).Methods("PUT")
	gmux.HandleFunc("/rules/{RuleName}", DeleteRule).Methods("DELETE")
	gmux.HandleFunc("/tries", ListTries).Methods("GET")
	gmux.HandleFunc("/tries/{TrieName}", GetTrie).Methods("GET")
	gmux.HandleFunc("/tries/{TrieName}", UpdateTrie).Methods("PUT")
	gmux.HandleFunc("/tries/{TrieName}", DeleteTrie).Methods("DELETE")

	// Router Visualizations
	gmux.HandleFunc("/visualize/router", graph.VisualizeIndex)
	gmux.HandleFunc("/visualize/router/tries/{name}/json", graph.TrieJson)
	gmux.HandleFunc("/visualize/router/tries/{name}/dot", graph.TrieDot)
	gmux.HandleFunc("/visualize/router/tries/{name}/svg", graph.TrieSvg)

	// Instance Management
	gmux.HandleFunc("/instances", ListContainers).Methods("GET")
	gmux.HandleFunc("/instances", TeardownContainers).Methods("DELETE")
	gmux.HandleFunc("/instances/{Id}", ContainerIdGet).Methods("GET")
	gmux.HandleFunc("/instances/{Id}", TeardownContainerId).Methods("DELETE")
	gmux.HandleFunc("/instances/{Id}/copy", CopyContainer).Methods("POST")
	gmux.HandleFunc("/instances/{Id}/move", MoveContainer).Methods("POST")
	gmux.HandleFunc("/instances/{Id}/maint", ContainerMaintenance).Methods("POST")
	gmux.HandleFunc("/instances/apps", ListApps).Methods("GET")
	gmux.HandleFunc("/instances/apps/{App}", Teardown).Methods("DELETE")
	gmux.HandleFunc("/instances/apps/{App}/shas", ListShas).Methods("GET")
	gmux.HandleFunc("/instances/apps/{App}/shas/{Sha}", Teardown).Methods("DELETE")
	gmux.HandleFunc("/instances/apps/{App}/shas/{Sha}/envs", DeployListEnvs).Methods("GET")
	gmux.HandleFunc("/instances/apps/{App}/shas/{Sha}/envs/{Env}", Teardown).Methods("DELETE")
	gmux.HandleFunc("/instances/apps/{App}/shas/{Sha}/envs/{Env}/containers", ListContainers).Methods("GET")
	gmux.HandleFunc("/instances/apps/{App}/shas/{Sha}/envs/{Env}/containers", Deploy).Methods("POST")

	// LDAP Management
	gmux.HandleFunc("/users/{User}", GetPermissions).Methods("GET")
	gmux.HandleFunc("/teams", ListTeams).Methods("GET")
	gmux.HandleFunc("/teams/{Team}/apps", ListTeamApps).Methods("GET")
	gmux.HandleFunc("/teams/{Team}/apps/{App}", AllowApp).Methods("PUT")
	gmux.HandleFunc("/teams/{Team}/apps/{App}", DisallowApp).Methods("DELETE")
	gmux.HandleFunc("/teams/{Team}", CreateTeam).Methods("PUT")
	gmux.HandleFunc("/teams/{Team}", DeleteTeam).Methods("DELETE")
	gmux.HandleFunc("/teams/{Team}/admins", ListTeamAdmins).Methods("GET")
	gmux.HandleFunc("/teams/{Team}/admins/{Admin}", AddTeamAdmin).Methods("PUT")
	gmux.HandleFunc("/teams/{Team}/admins/{Admin}", RemoveTeamAdmin).Methods("DELETE")
	gmux.HandleFunc("/teams/{Team}/emails/{Email}", AddTeamEmail).Methods("PUT")
	gmux.HandleFunc("/teams/{Team}/emails/{Email}", RemoveTeamEmail).Methods("DELETE")
	gmux.HandleFunc("/teams/{Team}/members", ListTeamMembers).Methods("GET")
	gmux.HandleFunc("/teams/{Team}/members/{Member}", AddTeamMember).Methods("PUT")
	gmux.HandleFunc("/teams/{Team}/members/{Member}", RemoveTeamMember).Methods("DELETE")

	// Environment Management
	gmux.HandleFunc("/envs/{Env}/resolve/{DepNames}", ResolveDeps).Methods("GET")
	gmux.HandleFunc("/envs/{Env}/deps/{DepName}", GetDep).Methods("GET")
	gmux.HandleFunc("/envs/{Env}/deps/{DepName}", UpdateDep).Methods("PUT")
	gmux.HandleFunc("/envs/{Env}/deps/{DepName}", DeleteDep).Methods("DELETE")
	gmux.HandleFunc("/envs/{Env}/deps", GetDep).Methods("GET")
	gmux.HandleFunc("/envs/{Env}", GetEnv).Methods("GET")
	gmux.HandleFunc("/envs/{Env}", UpdateEnv).Methods("PUT")
	gmux.HandleFunc("/envs/{Env}", DeleteEnv).Methods("DELETE")
	gmux.HandleFunc("/envs", ListEnvs).Methods("GET")

	// Static Assets
	staticPath := "/" + staticDir + "/"
	fileServer := http.StripPrefix(staticPath, http.FileServer(http.Dir("./"+staticDir)))
	gmux.NewRoute().PathPrefix(staticPath).Handler(fileServer)

	handler := apachelog.NewHandler(HandlerFunc(gmux), os.Stderr)
	server = &http.Server{Addr: listenAddr, Handler: handler}
	lAddr = listenAddr
	return nil
}

func listenAndServeTLS(cert, key []byte) error {
	addr := server.Addr
	if addr == "" {
		log.Printf("Current Address: %s", addr)
		panic("[API] Current Address is not HTTPS")
	}

	config := &tls.Config{}
	if server.TLSConfig != nil {
		config = server.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}

	conn, err := net.Listen("tcp", lAddr)
	if err != nil {
		return err
	}

	return server.Serve(tls.NewListener(conn, config))
}

func Output(obj map[string]interface{}, err error) string {
	var bytes []byte
	if err != nil {
		m := map[string]interface{}{"Error": err.Error()}
		bytes, err = json.Marshal(m)
		if err != nil {
			return err.Error()
		}
		return string(bytes)
	}
	bytes, err = json.Marshal(obj)
	return string(bytes)
}

func Listen() {
	if server == nil {
		panic("Not Initialized.")
	}
	log.Println("[API] Listening on", lAddr)
	log.Fatal(listenAndServeTLS(crypto.SERVER_CERT, crypto.SERVER_KEY).Error())
}
