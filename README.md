# go-SRT

SRT(Super Rapid Train) application golang wrapper

This project was inspired from [korail2](https://github.com/carpedm20/korail2) of [carpedm20](https://github.com/carpedm20) and [SRT](https://github.com/ryanking13/SRT) of myself.

## Installation

```sh
go get -u github.com/ryanking13/go-SRT
```

## Usage

### 1. Login

```go
package main

import srt "github.com/ryanking13/go-SRT"

func main() {
	client := srt.New()
	client.Login("1234567890", YOUR_PASSWORD)        // with membership number
	client.Login("def6488@gmail.com", YOUR_PASSWORD) // with email
	client.Login("010-1234-xxxx", YOUR_PASSWORD)     // with phone number
}
```

use `SetDebug()` to see some debugging messages

```go
client.SetDebug()
```

### 2. Searching trains

use `SearchTrain` method.

Attributes of `srt.SearchParams`:

- Dep : A departure station in Korean ex) '수서'
- Arr : A arrival station in Korean ex) '부산'
- Date : (optional) (default: today) A departure date in yyyyMMdd format 
- Time : (optional) (default: 000000) A departure time in hhmmss format 
- IncludeSoldOut: (optional) (default: False) include trains which are sold out 

```go
trains, err := client.SearchTrain(&srt.SearchParams{
    Dep:            "수서",
    Arr:            "부산",
    Date:           "20200128",
    Time:           "144000",
    IncludeSoldOut: false,
})
if err != nil {
    panic(err)
}
```

### 3. Making a reservation

use `Reserve` method.

Attributes of `srt.ReserveParams`:

- Train: `*Train` returned by `SearchTrain()`
- Passengers (optional, default is one Adult)

```go
// ...
// trains, _ := client.SearchTrain(...)
reservation, err := client.Reserve(&srt.ReserveParams{
    Train:      trains[0],
    Passengers: []*srt.Passenger{srt.Adult(2), srt.Child(1)},
})

if err != nil {
    panic(err)
}

fmt.Println(reservation)
// [SRT] 02월 10일, 수서~부산(15:00~17:34) 129700원(3석), 구입기한 01월 28일 16:40
```

#### Passenger

__WARNING: 충분히 테스트되지 않음__

- Adult
- Child
- Senior
- Disability1To3
- Disability4To6

### 4. Getting reserved tickets

Use `Reservations` method.

```go
// ...
reservations, err := client.Reservations()
for _, r := range reservations {
    fmt.Println(r)
}

// [SRT] 02월 10일, 수서~부산(15:00~17:34) 129700원(3석), 구입기한 01월 28일 16:40
// ...
```

### 5. Canceling reservation

Use `Cancel` method.

```go
reservation, err := client.Reserve(&srt.ReserveParams{
    Train:      trains[0],
    Passengers: []*srt.Passenger{srt.Adult(2), srt.Child(1)},
})
client.Cancel(reservation)

// OR

reservations, err := client.Reservations()
client.Cancel(reservations[0])
```

## Example

Check [srt-reserve](https://github.com/ryanking13/go-SRT/tree/master/cmd/srt-reserve) for sample Usage.

`srt-reserve` is a CLI application for SRT ticketing.

You can also download Windows binary from [Release](https://github.com/ryanking13/go-SRT/releases) section.

```sh
go get -u github.com/ryanking13/go-SRT/...

srt-reserve
```

## TODO

- Add tests for CI

## Testing

```sh
# Linux
# export SRT_USERNAME=<SRT_USERNAME>
# export SRT_PASSWORD=<SRT_PASSWORD>

# Windows
# set SRT_USERNAME=<SRT_USERNAME>
# set SRT_PASSWORD=<SRT_PASSWORD>

go test -failfast -v
```