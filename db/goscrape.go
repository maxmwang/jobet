package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"
)

const (
	savepointSite    string = ""
	savepointCompany string = ""
)

type company struct {
	name string
	site string
}

func MigrateFromGoscrape() {
	dbGoscrape, err := sql.Open("sqlite3", "companies.db")
	if err != nil {
		panic(err)
	}
	defer dbGoscrape.Close()

	dbJobet, err := sql.Open("sqlite3", "jobet.db")
	if err != nil {
		panic(err)
	}
	defer dbJobet.Close()

	goscrapeRows := getRows(dbGoscrape)
	jobetRows := getRows(dbJobet)

	goscrapeCompanies := make([]company, 0)
	for goscrapeRows.Next() {
		var c company
		err := goscrapeRows.Scan(&c.name, &c.site)
		if err != nil {
			panic(err)
		}
		goscrapeCompanies = append(goscrapeCompanies, c)
	}

	jobetCompanies := make([]company, 0)
	for jobetRows.Next() {
		var c company
		var void interface{}
		err := jobetRows.Scan(&c.name, &void, &c.site, &void, &void)
		if err != nil {
			panic(err)
		}
		jobetCompanies = append(jobetCompanies, c)
	}

	jobetCompanyNames := make(map[string]struct{})
	for _, c := range jobetCompanies {
		jobetCompanyNames[c.name] = struct{}{}
	}

	fmt.Printf("goscrape: %d\n", len(goscrapeCompanies))
	fmt.Printf("jobet: %d\n", len(jobetCompanies))

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		dbGoscrape.Close()
		dbJobet.Close()
		panic("received interrupt")
	}()

	for _, c := range goscrapeCompanies {
		if _, ok := jobetCompanyNames[c.name]; ok {
			fmt.Printf("company=%s already in\n", c.name)
			continue
		}

		fmt.Printf("company=%s, site=%s\n", c.name, c.site)

		fmt.Println("alias:")
		alias := c.name
		_, err = fmt.Scanln(&alias)
		if err != nil && err.Error() != "unexpected newline" {
			panic(err)
		}

		fmt.Println("priority:")
		priority := 5
		_, err = fmt.Scanf("%d\n", &priority)
		if err != nil && err.Error() != "unexpected newline" {
			panic(err)
		}

		fmt.Printf("INSERT INTO companies (name, alias, site, priority) VALUES (%s, %s, %s, %d)\n", c.name, alias, c.site, priority)

		prepare, err := dbJobet.Prepare("INSERT INTO companies (name, alias, site, priority) VALUES (?, ?, ?, ?)")
		if err != nil {
			panic(err)
		}

		_, err = prepare.Exec(c.name, alias, c.site, priority)
		if err != nil {
			panic(err)
		}
	}
}

func getRows(db *sql.DB) *sql.Rows {
	prepare, err := db.Prepare("SELECT * FROM companies WHERE (site = ? AND name >= ?) OR site > ? ORDER BY site, name")
	if err != nil {
		panic(err)
	}
	rows, err := prepare.Query(savepointSite, savepointCompany, savepointSite)
	if err != nil {
		panic(err)
	}

	return rows
}
