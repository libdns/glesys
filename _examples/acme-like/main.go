// SPDX-FileCopyrightText: 2024 Peter Magnusson <me@kmpm.se>
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/libdns/glesys"
	"github.com/libdns/libdns"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func main() {
	user := flag.String("user", "", "User/Project name")
	key := flag.String("key", "", "API-Key")
	zone := flag.String("zone", "", "domainname / zone")
	timeout := flag.Int("timeout", 60, "seconds to wait before delete")

	flag.Parse()

	if *zone == "" {
		fmt.Fprintf(os.Stderr, "GLESYS_ZONE not set\n")
		os.Exit(1)
	}

	if *user == "" {
		exitOnError(fmt.Errorf("user is not set"))
	}

	if *key == "" {
		exitOnError(fmt.Errorf("key is not set"))
	}

	fmt.Printf("zone: %s, user: %s\n", zone, user)
	p := &glesys.Provider{
		Project: *user,
		ApiKey:  *key,
	}
	ctx := context.TODO()
	fmt.Println("appending")
	res, err := p.AppendRecords(ctx, *zone,
		[]libdns.Record{
			{Name: "_libdns-challenge", Type: "TXT", Value: "Zgu7tw287LB-LpXyTHYLeROag9-4CLHnM77zvTEvH6o"},
		})
	exitOnError(err)
	printRecords("after append", res)

	fmt.Printf("Will sleep for %d seconds... then delete", timeout)
	sleep := time.Duration(*timeout) * time.Second
	time.Sleep(sleep)

	resAll, err := p.GetRecords(ctx, *zone)
	exitOnError(err)
	printRecords("after all", resAll)

	// delete first
	res, err = p.DeleteRecords(ctx, *zone, res)
	exitOnError(err)
	printRecords("after delete", res)

	// check final result
	resAll, err = p.GetRecords(ctx, *zone)
	exitOnError(err)
	printRecords("after all", resAll)

	fmt.Println("Done!")
}

func printRecords(title string, records []libdns.Record) {
	fmt.Println(title)
	for i, r := range records {
		fmt.Printf("  [%d] %+v\n", i, r)
	}
}
