package rest

import (
	"bytes"
	"io/ioutil"
	"net"
	"strconv"
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

func updateSecNginxConf(configPath string, nginxConfigStubPath string, destinations []string, openHTTPPort int, openRTMPPort int) error {
	stubConfigFile, err := ioutil.ReadFile(nginxConfigStubPath)
	if err != nil {
		return err
	}

	// Update push block
	var pushBlockBuffer bytes.Buffer

	for _, dst := range destinations {
		pushBlockBuffer.WriteString("push " + dst + ";\n")
	}

	configFileX := strings.Replace(string(stubConfigFile),
		"XPLEX_PUSH_DESTS",
		pushBlockBuffer.String(),
		-1)

	configFileY := strings.Replace(string(configFileX),
		"XPLEX_HTTP_PORT",
		strconv.Itoa(openHTTPPort),
		-1)

	configFileZ := strings.Replace(string(configFileY),
		"XPLEX_RTMP_PORT",
		strconv.Itoa(openRTMPPort),
		-1)

	err = ioutil.WriteFile(configPath, []byte(configFileZ), 0)
	if err != nil {
		panic(err)
	}

	return nil

}

func updateNginxSystemdConf(configPath string, systemdConfigStubPath string, pidPath string, streamKey string) error {
	stubConfigFile, err := ioutil.ReadFile(systemdConfigStubPath)
	if err != nil {
		return err
	}

	configFileX := strings.Replace(string(stubConfigFile),
		"XPLEX_STREAMKEY",
		streamKey,
		-1)

	configFileY := strings.Replace(string(configFileX),
		"XPLEX_CONFIGPATH",
		configPath,
		-1)

	configFileZ := strings.Replace(string(configFileY),
		"XPLEX_PIDPATH",
		pidPath,
		-1)

	err = ioutil.WriteFile("/lib/systemd/system/nginx-"+streamKey+".service", []byte(configFileZ), 0)
	if err != nil {
		panic(err)
	}

	return nil

}
