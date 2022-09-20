# calendar

[![Go Reference](https://pkg.go.dev/badge/github.com/edgelaboratories/calendar.svg)](https://pkg.go.dev/github.com/edgelaboratories/calendar)
![GolangCI Lint](https://github.com/edgelaboratories/calendar/workflows/Golangci-Lint/badge.svg)

Provide basic operations on calendars in Go.

## Install

```bash
go get -u github.com/edgelaboratories/calendar
```

## Requirements

- [go 1.19.x](https://golang.org/dl/)

## Test

Run the following:

```bash
make test
```

Check out the [`Makefile`](Makefile) for more information.

## Purpose

This project aims at offering a unique, simple and fast module to manage calendars with daily granularity. Dates representation is based on [`github.com/fxtlabs/date`](https://github.com/fxtlabs/date).

This project **doesn't aim** at supporting daycount conventions. Have a look at [`github.com/edgelaboratories/daycount`](https://github.com/edgelaboratories/daycount) instead.

## Usage

```go
package main

import (
    "fmt"
    "github.com/edgelaboratories/calendar"
    "github.com/fxtlabs/date"
)

func main() {
    c := calendar.New(calendar.BusinessDays)

    fmt.Println(c.Convention()) // output is "BusinessDays"
    
    fmt.Println(c.DaysBetween(
        date.New(2021, time.October, 13), // Wednesday
        date.New(2021, time.October, 15), // Friday
    )) // output is 2

    fmt.Println(c.Add(
        date.New(2021, time.October, 13), // Wednesday
        2,
    ).String()) // output is Friday "2021-10-15"
}
```
