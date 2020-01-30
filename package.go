package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/howeyc/crc16"
)

type Package map[string]interface{}

func CheckIntegrity(p []byte) bool {
	lastCommaIndex := bytes.LastIndex(p, []byte(","))
	if lastCommaIndex == -1 {
		return false
	}

	computedCRC := crc16.ChecksumCCITTFalse(p[:lastCommaIndex])

	lastColonIndex := bytes.LastIndex(p, []byte(":"))

	if lastColonIndex == -1 {
		return false
	}

	if string(p[lastColonIndex-5:lastColonIndex]) != `"crc"` {
		return false
	}

	crcStr := string(p[lastColonIndex+2 : len(p)-2])

	if strings.HasPrefix(crcStr, "0x") {
		crcStr = crcStr[2:]
	}

	crc, err := strconv.ParseUint(crcStr, 16, 16)
	if err != nil {
		return false
	}

	return computedCRC == uint16(crc)
}

func ParsePackage(p []byte) (*Package, error) {
	var pkt Package
	if err := json.Unmarshal(p, &pkt); err != nil {
		return nil, fmt.Errorf("Can't parse package: %s", err)
	}

	delete(pkt, "crc")

	return &pkt, nil
}

func (p *Package) IsValid() bool {
	_, ok := (*p)["mac"]
	return ok
}

func (p *Package) IsLamp() bool {
	tlevel, ok := (*p)["tlevel"]

	if !ok {
		return false
	}

	return int(tlevel.(float64)) != 0
}

func (p *Package) IsOurController() bool {
	mac, ok := (*p)["mac"]

	if !ok {
		return false
	}

	tlevel, ok := (*p)["tlevel"]

	if !ok {
		return false
	}

	smac, ok := (*p)["smac"]

	if !ok {
		return false
	}

	return int(tlevel.(float64)) == 0 && mac.(string) == smac.(string)
}

func (p *Package) IsController() bool {

	tlevel, ok := (*p)["tlevel"]

	if !ok {
		return false
	}

	return int(tlevel.(float64)) == 0
}

func (p *Package) IsOtherController() bool {
	mac, ok := (*p)["mac"]

	if !ok {
		return false
	}

	tlevel, ok := (*p)["tlevel"]

	if !ok {
		return false
	}

	smac, ok := (*p)["smac"]

	if !ok {
		return false
	}

	return int(tlevel.(float64)) == 0 && mac.(string) != smac.(string)
}
