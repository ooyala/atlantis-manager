package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	. "atlantis/manager/rpc/client"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	. "atlantis/supervisor/rpc/types"
	"errors"
	"log"
	"time"
)

type HealthCheckExecutor struct {
	arg   ManagerHealthCheckArg
	reply *ManagerHealthCheckReply
}

func (e *HealthCheckExecutor) Request() interface{} {
	return e.arg
}

func (e *HealthCheckExecutor) Result() interface{} {
	return e.reply
}

func (e *HealthCheckExecutor) Description() string {
	return "HealthCheck"
}

func (e *HealthCheckExecutor) Execute(t *Task) error {
	e.reply.Region = Region
	e.reply.Zone = Zone
	if Tracker.UnderMaintenance() {
		e.reply.Status = StatusMaintenance
	} else {
		e.reply.Status = StatusOk
	}
	t.Log("[RPC][HealthCheck] -> region: %s", e.reply.Region)
	t.Log("[RPC][HealthCheck] -> zone: %s", e.reply.Zone)
	t.Log("[RPC][HealthCheck] -> status: %s", e.reply.Status)
	return nil
}

func (e *HealthCheckExecutor) Authorize() error {
	return nil // allow anyone to check health
}

func (e *HealthCheckExecutor) AllowDuringMaintenance() bool {
	return true // allow checking health during maintenance.
}

func (o *Manager) HealthCheck(arg ManagerHealthCheckArg, reply *ManagerHealthCheckReply) error {
	return NewTask("HealthCheck", &HealthCheckExecutor{arg, reply}).Run()
}

const (
	healthCheckPeriod = time.Minute
)

var supervisorStatus = map[string]error{}
var supervisorDie = map[string]chan bool{}
var managerStatus = map[string]error{}
var managerDie = map[string]chan bool{}

func watchSupervisor(host string) {
	dieChan := make(chan bool)
	supervisorDie[host] = dieChan
	tick := time.NewTicker(healthCheckPeriod)
	var err error
	var health *SupervisorHealthCheckReply
	for {
		select {
		case <-tick.C:
			log.Printf("Checking Supervisor Health %s:%s", host, supervisor.Port)
			health, err = supervisor.HealthCheck(host)
			if err == nil && health.Status != StatusOk && health.Status != StatusFull {
				err = errors.New("Status is " + health.Status)
			}
			supervisorStatus[host] = err
		case <-dieChan:
			close(dieChan)
			tick.Stop()
			delete(supervisorDie, host)
			delete(supervisorStatus, host)
			return
		}
	}
}

func managerHealthCheck(host string) (*ManagerHealthCheckReply, error) {
	args := ManagerHealthCheckArg{}
	var reply ManagerHealthCheckReply
	return &reply, NewManagerRPCClient(host+lAddr).Call("HealthCheck", args, &reply)
}

func watchManager(host string) {
	dieChan := make(chan bool)
	managerDie[host] = dieChan
	tick := time.NewTicker(healthCheckPeriod)
	var err error
	var health *ManagerHealthCheckReply
	for {
		select {
		case <-tick.C:
			log.Printf("Checking Manager Health %s%s", host, lAddr)
			health, err = managerHealthCheck(host)
			if err == nil && health.Status != StatusOk {
				err = errors.New("Status is " + health.Status)
			}
			managerStatus[host] = err
		case <-dieChan:
			close(dieChan)
			tick.Stop()
			delete(managerDie, host)
			delete(managerStatus, host)
			return
		}
	}
}
