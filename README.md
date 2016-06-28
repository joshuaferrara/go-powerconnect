# powerconnect

    import "github.com/joshuaferrara/go-powerconnect"
    
## Info

I created this library to eventually use with FRC Team 254s [cheesy-arena](https://github.com/Team254/cheesy-arena) code. They have support for a Cisco switch, though I picked up some cheap Dell PowerConnect 2724s and decided to try my hand at reverse engineering the web server (unfortunately, there was no telnet console). I've only created functions that would be needed for my use-case (managing VLANs), though if any others are needed, feel free to send a PR or let me know. I've included a PDF file in the repo that details my findings throughout the reverse engineering process.
    
## Usage

#### func  Login

```go
func Login(ip, username, password string) (AuthInfo, error)
```
Logs into the PowerConnect 2724 switch. Creates a AuthInfo objec to use for
authentication in other functions.
* `ip` - ip of PowerConnect switch. Ex: 192.168.2.1
* `username` - username to login with.Default: admin
* `password` - password to login with. Default: [blank]

#### func  SetVLAN

```go
func SetVLAN(vlan, portSettings string, aInfo AuthInfo) (bool, error)
```
Sets a VLAN group to defined settings.
* `vlan` - which VLAN you want to edit.
* `portSettings` - 24 digit string. Value at position 1 indicates setting for port 1. Possible values: 0 - not in VLAN, 1 - in VLAN; untagged, 3 - in VLAN; tagged.
* `aInfo` - AuthInfo from Login function.

Ex: `SetVLAN("2", "010000000000100000300000", aInfo);`

#### type AuthInfo

```go
type AuthInfo struct {
}
```

Holds info needed to authenticate requests.