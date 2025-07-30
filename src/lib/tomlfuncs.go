package lib

import (
    "github.com/BurntSushi/toml"
    "github.com/chigopher/pathlib"
)

func ParseIsoToml(isoTomlPath *pathlib.Path) ([]IsoEntry, error) {
    tomlFile, err := isoTomlPath.Open()
    if err != nil {
        return nil, err
    }
    defer tomlFile.Close()

    decoder := toml.NewDecoder(tomlFile)
    thisIsoToml := new(IsoToml)
    if _, err := decoder.Decode(thisIsoToml); err != nil {
        return nil, err
    }

    tomlFile.Close()

    return thisIsoToml.Entries, nil
}