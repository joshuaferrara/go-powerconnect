package powerconnect

import (
	"net/http"
	"net/url"
	"bytes"
	"strings"
	"errors"
	"io/ioutil"
)

// The Dell PowerConnect expects form values to be in a certain order.
// The original Encode() function for a url.Values reorders form values alphabetically
// based on their key. If the form values are not in the correct order, the switch will
// terminate the TCP connection. 
func customEncode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	for _, k := range keys {
		vs := v[k]
		prefix := url.QueryEscape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

// Sets a VLAN group to defined settings.
// vlan - which VLAN you want to edit.
// portSettings - 24 digit string. Value at position 1 indicates setting for port 1.
//                Possible values: 0 - not in VLAN, 1 - in VLAN; untagged, 3 - in VLAN; tagged.
// aInfo - AuthInfo from Login function.
func SetVLAN(vlan, portSettings string, aInfo AuthInfo) (bool, error) {
	form := url.Values{"op": {"select"}, "vlan": {string(vlan)}, "ports": {portSettings}, "trunks": {"000000"}}

	vlanReq, _ := http.NewRequest("POST", "http://" + aInfo.ip + "/tgi/vlan.tgi", strings.NewReader(customEncode(form)))
	vlanReq.Header.Add("Referer", "http://" + aInfo.ip + "/vlan.htm?op=select&vlan=2")
	vlanReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	vlanReq.Header.Add("Accept", "text/html, application/xhtml+xml, image/jxr, */*")
	vlanReq.Header.Add("Connection", "Keep-Alive")
	vlanReq.Header.Add("Cache-Control", "no-cache")
	vlanReq.Header.Add("Accept-Encoding", "gzip, deflate")
	vlanResp, vlanErr := aInfo.client.Do(vlanReq)

	// This usually means the connection was dropped for some reason.
	// At times, the connection seems to be randomly dropped despite a correct request...
	if vlanErr != nil {
		return false, vlanErr
	}

	vlanData, _ := ioutil.ReadAll(vlanResp.Body)
	vlanBody := string(vlanData)

	// Check to see if there was no weird auth issues
	vlanMembersCurrentIndex := strings.Index(vlanBody, "vlanMembersCurrent")
	if vlanMembersCurrentIndex == -1 {
		return false, errors.New("Error setting VLANs. Ensure your AuthInfo is valid.")
	} else {
		// Check to see if the VLAN was actually set
		vlanMembersString := vlanBody[vlanMembersCurrentIndex + 20 : vlanMembersCurrentIndex + 20 + 24]
		if vlanMembersString != portSettings {
			return false, errors.New("Error setting VLANs. Your request was processed, but the switch did not accept the changes.")
		}
	}

	return true, nil
}