// Copyright
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bosun-monitor/snmp"
)

var (
	community = flag.String("c", "public", "community to use")
	v1        = flag.Bool("v1", false, "use snmpv1 instead of v2c")
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Println("Usage: snmpwalk [-c=public] host oid\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	host := flag.Args()[0]
	oid := flag.Args()[1]

	h, err := snmp.New(host, *community)
	if *v1 {
		h.SetVersion(snmp.V1)
	}
	s, err := h.Walk(oid)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	for s.Next() {
		var a interface{}
		_, err := s.Scan(&a)
		if err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
		switch i := a.(type) {
		case int, int64, int32:
			fmt.Printf("INT: %d\n", i)
		case []byte:
			fmt.Printf("BYTE: %s\n", string(i))
		default:
			fmt.Printf("dont know the type: %+v (%+v)\n", i, a)
		}
	}
	if err = s.Err(); err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
