package rpc

import (
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/crypto"
	"atlantis/manager/datamodel"
	"atlantis/manager/manager"
	"atlantis/manager/supervisor"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strings"
	"time"
)

type ManagerRPC bool

var (
	lAddr                string
	lPort                string
	l                    net.Listener
	server               *rpc.Server
	config               *tls.Config
	CPUSharesIncrement   = uint(1) // default to no increment
	MemoryLimitIncrement = uint(1) // default to no increment
)

func Init(listenAddr string, supervisorPort uint16, cpuIncr, memIncr uint, resDuration time.Duration) error {
	var err error
	err = LoadEnvs()
	if err != nil {
		return err
	}
	CPUSharesIncrement = cpuIncr
	MemoryLimitIncrement = memIncr
	atlantis.Tracker.ResultDuration = resDuration
	// init rpc stuff here
	lAddr = listenAddr
	lPort = strings.Split(lAddr, ":")[1]
	supervisor.Init(fmt.Sprintf("%d", supervisorPort))
	manager.Init(lPort)
	manager := new(ManagerRPC)
	server = rpc.NewServer()
	server.Register(manager)
	config := &tls.Config{}
	config.InsecureSkipVerify = true
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair(crypto.SERVER_CERT, crypto.SERVER_KEY)

	l, err = tls.Listen("tcp", lAddr, config)
	return err
}

func Listen() {
	go selfRegister()
	if l == nil {
		panic("Not Initialized.")
	}
	log.Println("[RPC] Listening on", lAddr)
	server.Accept(l)
}

func selfRegister() {
	log.Println("[SelfRegister] Registering Self.")
	zkManager, err := datamodel.GetManager(Region, Host)
	if err == nil && zkManager != nil {
		// i'm already registered
		log.Println("[SelfRegister] Already Registered")
		return
	}
	mgr, err := manager.Register(Region, Host, "", "")
	if err != nil {
		log.Fatalln("[SelfRegister] Failure: ", err)
	}
	log.Printf("[SelfRegister] Success: %s", mgr.ManagerCName)
}

func checkRole(role string, rType string) error {
	log.Printf("[CheckRole] checking myself (%s:%s) for %s:%s", Region, Host, rType, role)
	zkManager, err := datamodel.GetManager(Region, Host)
	if err != nil {
		return err
	}
	log.Printf("[CheckRole] roles: %v", zkManager.Roles)
	if !zkManager.HasRole(role, rType) {
		log.Printf("[CheckRole] role check fail.")
		managersWithRole := ""
		managers, err := datamodel.ListManagers()
		if err != nil {
			return err
		}
		for region, rManagers := range managers {
			for _, manager := range rManagers {
				zm, err := datamodel.GetManager(region, manager)
				if err != nil {
					continue
				}
				if zm.HasRole(role, rType) {
					managersWithRole = managersWithRole + zm.ManagerCName + "\n"
				}
			}
		}
		return errors.New(fmt.Sprintf("This manager does not have the ability to %s %s. "+
			"Please try one of these:\n%s", rType, role, managersWithRole))
	}
	log.Printf("[CheckRole] role check success.")
	return nil
}
