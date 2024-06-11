package portmapping

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Item struct {
	Port       int
	Network    string
	Status     bool
	TargetIP   string
	TargetPort int
	Desc       string
}

func LoadCSV(filename string) ([]*Item, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var d []*Item
	var e error
	for index, record := range records {
		if index == 0 {
			continue
		}
		itemStr := strings.Join(record, ",")
		if len(record) < 6 {
			e = errors.Join(e, fmt.Errorf("item:%s|err:%s", itemStr, errors.New("item does not match")))
			continue
		}
		var item = &Item{}
		port, err := strconv.Atoi(record[0])
		if err != nil {
			e = errors.Join(e, fmt.Errorf("item:%s|err:port %s is invalid", itemStr, record[0]))
			continue
		}
		item.Port = port
		portType := record[1]
		switch portType {
		case "tcp", "udp", "TCP", "UDP":
			item.Network = strings.ToLower(portType)
		default:
			e = errors.Join(e, fmt.Errorf("item:%s|err:type %s must is tcp or udp", itemStr, record[1]))
			continue
		}
		switch record[2] {
		case "1", "true", "True":
			item.Status = true
		case "0", "false", "False":
			item.Status = false
		default:
			e = errors.Join(e, fmt.Errorf("item:%s|err:status %s must is 1 or true or True or 0 or false or False", itemStr, record[2]))
			continue
		}
		targetIP := net.ParseIP(record[3])
		if targetIP == nil || targetIP.To4() == nil {
			e = errors.Join(e, fmt.Errorf("item:%s|err:target ip %s is invalid or is not ipv4", itemStr, record[2]))
			continue
		}
		item.TargetIP = targetIP.String()

		targetPort, err := strconv.Atoi(record[4])
		if err != nil {
			e = errors.Join(e, fmt.Errorf("item:%s|err:target port %s is invalid", itemStr, record[4]))
			continue
		}
		item.TargetPort = targetPort
		item.Desc = record[5]
		d = append(d, item)
	}
	return d, err
}
