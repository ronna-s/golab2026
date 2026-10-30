package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func Serve(l net.Listener, handle func(net.Conn) error) error {
	var (
		err error
		wg  sync.WaitGroup
	)
	for {
		var conn net.Conn
		if conn, err = l.Accept(); err != nil {
			break
		}
		wg.Go(func() {
			if handleErr := handle(conn); handleErr != nil {
				log.Println(handleErr)
			}
		})
	}
	wg.Wait()
	return err
}

func HandleConnection(conn net.Conn) error {
	conn.SetDeadline(time.Now().Add(time.Second))
	b, err := io.ReadAll(conn)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func main() {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		if err != nil {
			log.Printf("unexpected error establishing a connection: %s", err.Error())
		}
		conn.Write([]byte(output))
		conn.Close()

		time.Sleep(time.Second * 1)
		l.Close()
	}()
	err = Serve(l, HandleConnection)
	if err != nil && !errors.Is(err, net.ErrClosed) {
		log.Fatal(err)
	}
}

const output = "        ,_---~~~~~----._\n  _,,_,*^____      _____``*g*\\\"*,\n / __/ /'     ^.  /      \\ ^@q   f\n[  @f | @))    |  | @))   l  0 _/\n \\`/   \\~____ / __ \\_____/    \\\n  |           _l__l_           I\n  }          [______]           I\n  ]            | | |            |\n  ]             ~ ~             |\n  |                            |\n   |                           |"
