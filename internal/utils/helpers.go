package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/dicedb/dicedb-go/wire"
	"google.golang.org/protobuf/encoding/protojson"
)

// parseHostAndPort splits a URL string in format "host:port" and returns the host and port
func parseHostAndPort(url string) (string, int) {
	// If the URL is not in the "host:port" format, treat
	// the URL as the host and use the default port 7379
	host := url
	port := 7379

	// If the URL contains a colon, try to split it into host and port
	if strings.Contains(url, ":") {
		var err error
		var portStr string

		host, portStr, err = net.SplitHostPort(url)
		if err == nil {
			portInt, err := strconv.Atoi(portStr)
			if err == nil {
				port = portInt
			}
		}
	}

	return host, port
}

// FormatDiceDBResponse formats the DiceDB response
func FormatDiceDBResponse(resp *wire.Result) string {
	if resp.Status == wire.Status_ERR {
		return fmt.Sprintf("%s %s\n", ("ERR"), resp.Message)
	}

	var result strings.Builder

	// Copied from: https://github.com/DiceDB/dicedb-cli/blob/ironhawk/main.go#L136
	result.WriteString(resp.Message)
	if resp.Fingerprint64 != 0 {
		result.WriteString(fmt.Sprintf("[fingerprint=%d] ", resp.Fingerprint64))
	}

	switch resp.Response.(type) {
	case *wire.Result_GETRes:
		result.WriteString((fmt.Sprintf("\"%s\"\n", resp.GetGETRes().Value)))
	case *wire.Result_GETDELRes:
		result.WriteString(fmt.Sprintf("\"%s\"\n", resp.GetGETDELRes().Value))
	case *wire.Result_SETRes:
		result.WriteString("\n")
	case *wire.Result_FLUSHDBRes:
		result.WriteString("\n")
	case *wire.Result_DELRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetDELRes().Count))
	case *wire.Result_DECRRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetDECRRes().Value))
	case *wire.Result_INCRRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetINCRRes().Value))
	case *wire.Result_DECRBYRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetDECRBYRes().Value))
	case *wire.Result_INCRBYRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetINCRBYRes().Value))
	case *wire.Result_ECHORes:
		result.WriteString(fmt.Sprintf("%s\n", resp.GetECHORes().Message))
	case *wire.Result_EXISTSRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetEXISTSRes().Count))
	case *wire.Result_EXPIRERes:
		result.WriteString(fmt.Sprintf("%v\n", resp.GetEXPIRERes().IsChanged))
	case *wire.Result_EXPIREATRes:
		result.WriteString(fmt.Sprintf("%v\n", resp.GetEXPIREATRes().IsChanged))
	case *wire.Result_EXPIRETIMERes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetEXPIRETIMERes().UnixSec))
	case *wire.Result_TTLRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetTTLRes().Seconds))
	case *wire.Result_GETEXRes:
		result.WriteString(fmt.Sprintf("\"%s\"\n", resp.GetGETEXRes().Value))
	case *wire.Result_GETSETRes:
		result.WriteString(fmt.Sprintf("\"%s\"\n", resp.GetGETSETRes().Value))
	case *wire.Result_HANDSHAKERes:
		result.WriteString("\n")
	case *wire.Result_HGETRes:
		result.WriteString(fmt.Sprintf("\"%s\"\n", resp.GetHGETRes().Value))
	case *wire.Result_HSETRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetHSETRes().Count))
	case *wire.Result_HGETALLRes:
		result.WriteString("\n")
		for i, e := range resp.GetHGETALLRes().Elements {
			result.WriteString(fmt.Sprintf("%d) %s=\"%s\"\n", i, e.Key, e.Value))
		}
	case *wire.Result_KEYSRes:
		result.WriteString("\n")
		for i, key := range resp.GetKEYSRes().Keys {
			result.WriteString(fmt.Sprintf("%d) %s\n", i, key))
		}
	case *wire.Result_PINGRes:
		result.WriteString(fmt.Sprintf("\"%s\"\n", resp.GetPINGRes().Message))
	case *wire.Result_TYPERes:
		result.WriteString(fmt.Sprintf("%s\n", resp.GetTYPERes().Type))
	case *wire.Result_ZADDRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetZADDRes().Count))
	case *wire.Result_ZCOUNTRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetZCOUNTRes().Count))
	case *wire.Result_ZRANGERes:
		result.WriteString("\n")
		for _, e := range resp.GetZRANGERes().Elements {
			printZElement(e)
		}
	case *wire.Result_ZPOPMAXRes:
		result.WriteString("\n")
		for _, e := range resp.GetZPOPMAXRes().Elements {
			printZElement(e)
		}
	case *wire.Result_ZPOPMINRes:
		result.WriteString("\n")
		for _, e := range resp.GetZPOPMINRes().Elements {
			printZElement(e)
		}
	case *wire.Result_ZREMRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetZREMRes().Count))
	case *wire.Result_ZCARDRes:
		result.WriteString(fmt.Sprintf("%d\n", resp.GetZCARDRes().Count))
	case *wire.Result_ZRANKRes:
		printZElement(resp.GetZRANKRes().Element)
	case *wire.Result_GETWATCHRes:
		result.WriteString("\n")
	case *wire.Result_HGETWATCHRes:
		result.WriteString("\n")
	case *wire.Result_HGETALLWATCHRes:
		result.WriteString("\n")
	case *wire.Result_ZRANGEWATCHRes:
		result.WriteString("\n")
	case *wire.Result_ZCARDWATCHRes:
		result.WriteString("\n")
	case *wire.Result_ZCOUNTWATCHRes:
		result.WriteString("\n")
	case *wire.Result_ZRANKWATCHRes:
		result.WriteString("\n")
	case *wire.Result_UNWATCHRes:
		result.WriteString("\n")
	default:
		fmt.Println("note: this response is JSON serialized version of the response because it is not supported by this version of the CLI. You can upgrade the CLI to the latest version to get a formatted response.")
		b, err := protojson.Marshal(resp)
		if err != nil {
			log.Fatalf("failed to marshal to JSON: %v", err)
		}

		var m map[string]interface{}
		_ = json.Unmarshal(b, &m)

		nb, _ := json.MarshalIndent(m, "", "  ")
		result.WriteString((fmt.Sprintf("%s", string(nb))))
	}

	return strings.TrimSpace(result.String())
}

func printZElement(e *wire.ZElement) string {
	return fmt.Sprintf("%d) %d, %s\n", e.Rank, e.Score, e.Member)
}
