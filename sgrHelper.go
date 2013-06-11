package main

import (
	"io"
)

const (
	sgrFG        = "\x1b[38;5;"
	sgrBG        = "\x1b[48;5;"
	sgrRED       = "1m"
	sgrGREEN     = "2m"
	sgrYellow    = "3m"
	sgrBLUE      = "4m"
	sgrMagenta   = "5m"
	sgrCyan      = "6m"
	sgrLightGray = "7m"

	sgrBOLD  = "\x1b[1;4m"
	sgrRESET = "\x1b[0m"
)

var rainbowColors = [...]string{sgrRED, sgrBLUE, sgrGREEN, sgrYellow, sgrMagenta, sgrCyan, sgrLightGray}

func rainbowFG(text string, index int) string {
	return sgrFG + rainbowColors[index%len(rainbowColors)] + text + sgrRESET
}

func rainbowBG(text string, index int) string {
	return sgrBG + rainbowColors[index%len(rainbowColors)] + sgrFG + rainbowColors[index%len(rainbowColors)] + text + sgrRESET
}

func writeBold(w io.Writer, content []byte) {
	w.Write([]byte(sgrBOLD))
	w.Write(content)
	w.Write([]byte(sgrRESET))
}
