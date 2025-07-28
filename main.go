package main

import (
	"os"
	"strconv"

	"github.com/chigopher/pathlib"
	log "github.com/sirupsen/logrus"

	"carrierpigeondev/machv/src/lib"
	"carrierpigeondev/machv/src/options"
)

var (
	userHome  	string
	mainDir   	*pathlib.Path
	isoDir    	*pathlib.Path
	disksDir  	*pathlib.Path
	staticDir 	*pathlib.Path
	configDir 	*pathlib.Path
	isoTomlPath *pathlib.Path
	cfgTomlPath *pathlib.Path
)

func init() {
	var err error

	// set userhome, e.g., /home/x
	userHome, err = os.UserHomeDir()
	if err != nil {
		log.WithError(err).Fatal("A fatal error has occurred while getting user home directory")
	}

	// all dirs that will be initialized in initializeDirectories(), from globals in var
	mainDir = pathlib.NewPath(userHome).Join(".local", "share", "machv")
	isoDir = mainDir.Join("iso")
	disksDir = mainDir.Join("disks")
	staticDir = mainDir.Join("static")
	configDir = pathlib.NewPath(userHome).Join(".config", "machv")

	isoTomlPath = configDir.Join("iso.toml")
	cfgTomlPath = configDir.Join("cfg.toml")
}

func initializeDirectories() {
	// loop over globals and do logging
	for _, dir := range []*pathlib.Path{mainDir, isoDir, disksDir, configDir} {
		log.WithField("dir", dir).Debug("Checking dir if it exists")
		doesExist, err := dir.Exists()
		if err != nil {
			log.WithError(err).WithField("dir", dir).Fatal("A fatal error has occurred while checking dir")
		} else if !doesExist {
			log.WithField("dir", dir).Info("Creating dir as it did not exist")
			dir.MkdirAll()
		} else {
			log.WithField("dir", dir).Debug("The dir already exists")
		}
	}
}

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("Initializing machv...")
	initializeDirectories()  // does not return errors as errors are handled inside the function

	args := os.Args[1:]
	var opt int

	var verb string
	var noun string = ""

	if len(args) == 0 {
		selectOptions := []string{
			"Create new static virtual machine disk",
			"Create new usable virtual machine disk",
			"Load virtual machine disk",
		}
		iopt := lib.SelectOption(selectOptions, "Options")

		switch iopt {
			case "Create new static virtual machine disk": opt = 0
			case "Create new usable virtual machine disk": opt = 1
			case "Load virtual machine disk": 			   opt = 2
		}
	} else {
		if len(args) > 0 {
			verb = args[0]
		}
		if len(args) > 1 {
			noun = args[1]
		}
		
		verbAsInt, err := strconv.Atoi(verb)
		if err == nil {
			opt = verbAsInt
		} else {
			switch verb {
				case "create": {
					switch noun {
						case "static": opt = 0
						case "usable": opt = 1
						default:       opt = -3  // set invalid; invalid create noun
					}
				}
				case "load": {
					if noun != "" {
						opt = -2  // set invalid; load does not take a noun
					} else {
						opt = 2
					}
				}
				default: opt = -1  // set invalid; invalid verb
			}
		}
	}
	
	switch opt {
		// success cases
		case 0: options.OptionCreateNewStaticQCOW2(staticDir, isoTomlPath, isoDir)
		case 1: options.OptionCreateNewUsableQCOW2(staticDir, disksDir)
		case 2: options.OptionLaunchVirtualMachineFromUsableQCOW2(disksDir)

		// error cases
		case -1: log.Error("invalid verb")
		case -2: log.Error("'load' does not take a noun")
		case -3: log.Error("invalid noun for 'create'")
	}
}
