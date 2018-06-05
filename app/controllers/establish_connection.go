package controllers

import (
	"os/user"
	"strconv"
	"strings"

	"github.com/revel/revel"
)

func (c App) EstablishSSHConnection() revel.Result {
	sshIPHostname := strings.TrimSpace(c.Params.Get("ssh_ip"))
	sshUser := strings.TrimSpace(c.Params.Get("ssh_user"))
	sshPassword := c.Params.Get("ssh_password")
	sshPort := strings.TrimSpace(c.Params.Get("ssh_port"))

	var err error

	if sshIPHostname == "" {
		response := CompileJSONResult(false, "SSH IP is empty")
		return c.RenderJSON(response)
	}

	if sshUser == "" {
		username, err := user.Current()
		if err != nil {
			response := CompileJSONResult(false, "You didn't specify SSH user and we were not able to determine it from your system")
			return c.RenderJSON(response)
		}
		sshUser = username.Username
	}

	if sshPort == "" {
		sshPort = "22"
	} else {
		if _, err := strconv.Atoi(sshPort); err != nil {
			response := CompileJSONResult(false, "You specified wrong SSH port")
			return c.RenderJSON(response)
		}
	}

	SSHclient, SSHsession, err = createSSHSession(sshIPHostname, sshUser, sshPassword, sshPort)

	if err != nil {
		switch err.Error() {
		case "cannot dial":
			response := CompileJSONResult(false, "We could not reach '"+sshIPHostname+":"+sshPort+"' OR login/password is incorrect")
			return c.RenderJSON(response)
		case "unable to create session":
			response := CompileJSONResult(false, "We reached '"+sshIPHostname+":"+sshPort+"' but could not create a test session")
			return c.RenderJSON(response)
		default:
			response := CompileJSONResult(false, err.Error())
			return c.RenderJSON(response)
		}
	}
	return nil
}