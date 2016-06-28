package powerconnect

import "testing"
import "fmt"

func Test(t *testing.T) {
	aInfo, err := Login("172.16.0.95", "admin", "")

	if err != nil {
		fmt.Println("Login err:", err)
	} else {
		success, vlanErr := SetVLAN("2", "100000000000000000000000", aInfo)

		if vlanErr != nil {
			fmt.Println("SetVLAN Err:", vlanErr)
			return
		} else {
			fmt.Println("SetVLAN Success:", success)
		}
	}
}