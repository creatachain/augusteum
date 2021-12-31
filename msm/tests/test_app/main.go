package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/creatachain/augusteum/msm/example/code"
	"github.com/creatachain/augusteum/msm/types"
)

var msmType string

func init() {
	msmType = os.Getenv("MSM")
	if msmType == "" {
		msmType = "socket"
	}
}

func main() {
	testCounter()
}

const (
	maxMSMConnectTries = 10
)

func ensureMSMIsUp(typ string, n int) error {
	var err error
	cmdString := "msm-cli echo hello"
	if typ == "grpc" {
		cmdString = "msm-cli --msm grpc echo hello"
	}

	for i := 0; i < n; i++ {
		cmd := exec.Command("bash", "-c", cmdString)
		_, err = cmd.CombinedOutput()
		if err == nil {
			break
		}
		<-time.After(500 * time.Millisecond)
	}
	return err
}

func testCounter() {
	msmApp := os.Getenv("MSM_APP")
	if msmApp == "" {
		panic("No MSM_APP specified")
	}

	fmt.Printf("Running %s test with msm=%s\n", msmApp, msmType)
	subCommand := fmt.Sprintf("msm-cli %s", msmApp)
	cmd := exec.Command("bash", "-c", subCommand)
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		log.Fatalf("starting %q err: %v", msmApp, err)
	}
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("error on process kill: %v", err)
		}
		if err := cmd.Wait(); err != nil {
			log.Printf("error while waiting for cmd to exit: %v", err)
		}
	}()

	if err := ensureMSMIsUp(msmType, maxMSMConnectTries); err != nil {
		log.Fatalf("echo failed: %v", err) //nolint:gocritic
	}

	client := startClient(msmType)
	defer func() {
		if err := client.Stop(); err != nil {
			log.Printf("error trying client stop: %v", err)
		}
	}()

	setOption(client, "serial", "on")
	commit(client, nil)
	deliverTx(client, []byte("abc"), code.CodeTypeBadNonce, nil)
	commit(client, nil)
	deliverTx(client, []byte{0x00}, types.CodeTypeOK, nil)
	commit(client, []byte{0, 0, 0, 0, 0, 0, 0, 1})
	deliverTx(client, []byte{0x00}, code.CodeTypeBadNonce, nil)
	deliverTx(client, []byte{0x01}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x02}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x03}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x00, 0x04}, types.CodeTypeOK, nil)
	deliverTx(client, []byte{0x00, 0x00, 0x06}, code.CodeTypeBadNonce, nil)
	commit(client, []byte{0, 0, 0, 0, 0, 0, 0, 5})
}
