package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func init() {
	routes = append(routes, Route{"testHandler", "GET", "/pwd", testHandler})
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//logger.Info("Here")
	//fmt.Println(settings)

	sshConfig := &ssh.ClientConfig{
		User: settings.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(settings.Password),
			ssh.PasswordCallback(getPass),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				// Just send the password back for all questions
				answers := make([]string, len(questions))
				for i, _ := range answers {
					answers[i] = settings.Password // replace this
				}

				return answers, nil
			}),
		},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	host := settings.Host
	fmt.Println(settings.Host, settings.User)
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		logger.Error(err)
		return
	}

	sftp, err := sftp.NewClient(client)
	if err != nil {
		logger.Error(err)
		return
	}
	defer sftp.Close()

	pwd, _ := sftp.Getwd()
	fmt.Fprintf(w, pwd)

}

func getPass() (string, error) {
	return settings.Password, nil
}