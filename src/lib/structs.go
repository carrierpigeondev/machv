package lib

import (
    "strings"
)

type IsoToml struct {
    Entries []IsoEntry
}

type IsoEntry struct {
    Url  string
}

func (iso IsoEntry) FriendlyName() (string) {
    name := strings.Split(iso.Url, "/")

    if len(name) > 0 {
        return name[len(name)-1]
    } else {
        return ""  // to be used in a context where handling errors is a pain, so I'd rather just deal with it being an empty string
    }
}

type CfgToml struct {
    Config struct {
        iso_fetch_url string
    }
}