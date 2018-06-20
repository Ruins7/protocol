package bitcoin

import (
	"time"

	brpc "github.com/Oneledger/protocol/node/chains/bitcoin/rpc"

	"github.com/Oneledger/protocol/node/log"
	"strings"
	"net"
	"encoding/base64"
	"strconv"
	"github.com/Oneledger/protocol/node/convert"
)


func GetBtcClient(address string) *brpc.Bitcoind {
	addr:= strings.Split(address,":")
	if len(addr) < 2 {
		log.Error("address not in correct format", "fullAddress", address)
	}
	ip := net.ParseIP(addr[0])
	if ip == nil {
		log.Error("address can not be parsed", "addr", addr)
	}

	port := convert.GetInt(addr[1], 46688)

	usr, pass := getCredential()
	cli, err :=  brpc.New(ip.String(), port, usr, pass, false)
	if err != nil{
		log.Error("Can't get the btc rpc client at given address", "err", err)
		return nil
	}

	return cli
}

func getCredential() (usr string, pass string){
	//todo: getCredential from database which should be randomly generated when register or import if user already has bitcoin node
	usrBytes, err := base64.StdEncoding.DecodeString("b2x0ZXN0MDE=")
	if err != nil {
		log.Error(err.Error())
		usr = ""
	}
	usr = string(usrBytes)
	passBytes, err := base64.StdEncoding.DecodeString("b2xwYXNzMDE=")
	if err != nil {
		log.Error(err.Error())
		pass = ""
	}
	pass = string(passBytes)
	return
}

func ScheduleBlockGeneration(cli brpc.Bitcoind, interval time.Duration) chan bool {
	ticker := time.NewTicker(interval * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				cli.Generate(1)
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	return stop
}

func StopBlockGeneration(stop chan bool) {
	close(stop)
}
