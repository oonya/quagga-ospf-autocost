package ospfd
import (
	"regexp"
	"time"
	expect "github.com/google/goexpect"
	"github.com/ziutek/telnet"
	"fmt"
)

const (
	timeout = 1 * time.Second
)

func ConnectOspfd() (expect.Expecter, expect.Expecter, error) {
	localExp, _, err := telnetSpawn("localhost:2604", time.Second, expect.Verbose(true))
	if err != nil {
		panic(err)
	}

	remoteExp, _, err := telnetSpawn("192.168.1.1:2604", time.Second, expect.Verbose(true))
	if err != nil {
		panic(err)
	}

	// connect ospfd
	localExp.Expect(regexp.MustCompile("Password:"), timeout)
	localExp.Send("zebra\n")
	localExp.Expect(regexp.MustCompile("ospfd>"), timeout)
	localExp.Send("en\n")
	localExp.Expect(regexp.MustCompile("ospfd#"), timeout)

	remoteExp.Expect(regexp.MustCompile("Password:"), timeout)
	remoteExp.Send("zebra\n")
	remoteExp.Expect(regexp.MustCompile("ospfd>"), timeout)
	remoteExp.Send("en\n")
	remoteExp.Expect(regexp.MustCompile("ospfd#"), timeout)

	return localExp, remoteExp, nil
}

func CostSet(exp expect.Expecter, cost int, nic string) error {
	// TODO: error handling
	if err := exp.Send("configure t\n"); err != nil {
		return err
	}
	exp.Expect(regexp.MustCompile("ospfd(config)#"), timeout)
	exp.Send(fmt.Sprintf("interface %s\n", nic))
	exp.Expect(regexp.MustCompile("ospfd(config-if)"), timeout)
	exp.Send(fmt.Sprintf("ip ospf cost %d\n", cost))
	exp.Expect(regexp.MustCompile("ospfd(config-if)"), timeout)
	exp.Send("end\n")
	exp.Expect(regexp.MustCompile("ospfd#"), timeout)

	return nil
}


func telnetSpawn(addr string, timeout time.Duration, opts ...expect.Option) (expect.Expecter, <-chan error, error) {
	conn, err := telnet.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	resCh := make(chan error)

	return expect.SpawnGeneric(&expect.GenOptions{
		In:  conn,
		Out: conn,
		Wait: func() error {
			return <-resCh
		},
		Close: func() error {
			close(resCh)
			return conn.Close()
		},
		Check: func() bool { return true },
	}, timeout, opts...)
}
