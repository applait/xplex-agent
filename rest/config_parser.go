package rest

import (
	"bytes"
	"io/ioutil"
	"net"
	"strings"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func updateSecNginxConf(configPath string, configStubPath string, destinations []string) error {
	stubConfigFile, err := ioutil.ReadFile(configStubPath)
	if err != nil {
		return err
	}

	// Update push block
	var pushBlockBuffer bytes.Buffer

	for _, dst := range destinations {
		pushBlockBuffer.WriteString("push " + dst + ";")
	}

	configFileX := strings.Replace(string(stubConfigFile),
		"XPLEX_PUSH_DESTS",
		pushBlockBuffer.String(),
		-1)

	// Get available open port
	openPort, err := GetFreePort()
	if err != nil {
		return err
	}

	configFileY := strings.Replace(string(configFileX),
		"XPLEX_PORT",
		string(openPort),
		-1)

	err = ioutil.WriteFile(configPath, []byte(configFileY), 0)
	if err != nil {
		panic(err)
	}

	return nil

}