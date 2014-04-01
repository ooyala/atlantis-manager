/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package client

import (
	. "atlantis/manager/rpc/types"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type SSHCommand struct {
	Container string `short:"c" long:"container" description:"the container to ssh into"`
	Identity  string `short:"i" long:"identity" default:"~/.ssh/id_rsa" description:"the identity to use (SSH Key)"`
	PseudoTTY bool   `short:"t" long:"pseudo-tty" description:"force pseudo-tty allocation"`
}

func authorize(auth *ManagerAuthArg, container, publicKey string) (host, port string, err error) {
	arg := ManagerAuthorizeSSHArg{*auth, container, publicKey}
	var reply ManagerAuthorizeSSHReply
	err = rpcClient.Call("AuthorizeSSH", arg, &reply)
	return reply.Host, fmt.Sprintf("%d", reply.Port), err
}

func deauthorize(auth *ManagerAuthArg, container string) (err error) {
	arg := ManagerDeauthorizeSSHArg{*auth, container}
	var reply ManagerDeauthorizeSSHReply
	return rpcClient.Call("DeauthorizeSSH", arg, &reply)
}

func (c *SSHCommand) ExecuteRaw(args []string) error {
	if err := Init(); err != nil {
		return err
	}
	args = ExtractArgs([]*string{&c.Container}, args)
	Log("SSH ...")
	if c.Identity == "" {
		c.Identity = os.Getenv("HOME") + "/.ssh/id_rsa"
	} else {
		c.Identity = strings.Replace(c.Identity, "~", os.Getenv("HOME"), 1)
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := &ManagerAuthArg{User: user, Secret: secret}
	// fetch public key
	publicKeyBytes, err := ioutil.ReadFile(c.Identity + ".pub")
	if err != nil {
		return err
	}
	// ask manager to authorize ssh
	host, port, err := authorize(auth, c.Container, strings.TrimSpace(string(publicKeyBytes)))
	if err != nil {
		return err
	}
	Log("SSHing to %s:%s", host, port)
	cmdArray := []string{"-p", port, "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no"}
	if c.PseudoTTY {
		cmdArray = append(cmdArray, "-t")
	}
	if IsQuiet() || IsJson() {
		cmdArray = append(cmdArray, "-q")
	}
	cmdArray = append(cmdArray, "root@"+host)
	cmdArray = append(cmdArray, args...)
	cmd := exec.Command("ssh", cmdArray...)
	// let the ssh command hijack the terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		deauthorize(auth, c.Container)
		return err
	}
	if err := cmd.Wait(); err != nil {
		deauthorize(auth, c.Container)
		return err
	}
	// ask manager to deauthorize ssh
	deauthorize(auth, c.Container)
	// don't really care about error here
	Log("SSH connection closed.")
	return nil
}

func (c *SSHCommand) Execute(args []string) error {
	err := c.ExecuteRaw(args)
	if err != nil {
		return OutputError(err)
	}
	return OutputEmpty()
}

type TailCommand struct {
	Container string `short:"c" long:"container" description:"the container to ssh into"`
	Identity  string `short:"i" long:"identity" default:"~/.ssh/id_rsa" description:"the identity to use (SSH Key) Default:"`
	AppNumber int    `short:"n" long:"app" default:"0" description:"the app number to get logs for. Default:"`
	Type      string `short:"t" long:"type" default:"all" description:"which log to tail, can be one of (info, error, all) Default:"`
	Year      int    `short:"y" long:"year" description:"the year for which you want logs. Defaults to this year"`
	Month     int    `short:"m" long:"month" description:"the month for which you want logs. Defaults to this month"`
	Day       int    `short:"d" long:"day" description:"the day for which you want logs. Defaults to today"`
}

func (c *TailCommand) Execute(args []string) error {
	cmd := &SSHCommand{c.Container, c.Identity, true}
	year, monthStr, day := time.Now().UTC().Date()
	month := int(monthStr)
	if c.Year == 0 {
		c.Year = year
	}
	if c.Month == 0 {
		c.Month = month
	}
	if c.Day == 0 {
		c.Day = day
	}
	err := cmd.ExecuteRaw(append(args, fmt.Sprintf("tail -f /var/log/atlantis/syslog/app%d/%d/%02d/%02d/%s.log", c.AppNumber, c.Year, c.Month, c.Day, c.Type)))
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if statusErr, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if statusErr.ExitStatus() != 255 {
					Log(fmt.Sprintf("No logs exist for the requested cmd#/day/type combination. Cmd#: %d, Type: %s, Day: %d/%d/%d", c.AppNumber, c.Type, c.Month, c.Day, c.Year))
				}
			}
		}
	}
	return err
}
