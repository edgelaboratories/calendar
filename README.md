# calendar

Provide basic operations on calendars in Go.

## Installation

### Clone

`calendar` is a private project of [`edgelaboratories`](https://github.com/edgelaboratories) and must be installed by setting the `GONOPROXY` and `GONOSUMDB` variables first:

```bash
GONOPROXY=github.com/edgelaboratories/calendar GONOSUMDB=github.com/edgelaboratories/calendar go get github.com/edgelaboratories/calendar
```

### Requirements

[go 1.17.x](https://golang.org/dl/)

## Purpose

This project **aims** at offering a unique, simple and fast module to manage calendars based on `date.Date`.

The natural extension to intra-day durations implies the usage of `time.Time` and `time.Duration`, which is nowadays not needed.

This project **doesn't aim at** supporting daycount conventions; have a look at [`edgelaboratories/daycount`](https://github.com/edgelaboratories/daycount) instead.

## Usage

```go
package main

import (
    "fmt"
    "github.com/edgelaboratories/calendar"
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
