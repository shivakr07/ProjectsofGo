package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"` //you can add env-default:"production"
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// we will write the logic to parse this //this function must be executed successfully as it is required as as it is configuration
// means if this function faces some error then we shouln't start our application
func MustLoad() *Config {
	var configPath string
	// we need to pass the location/path of the config file
	// here we will use os package to get that from the normal env variable
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// if not available in environment then we will check other method [it might have been passed via flags/program arguments [we pass this while we run the program]]
		// like go run cmd/students-api/main.go -config-path xyz

		flags := flag.String("config", "", "path to the configuration file") //flagname config, default value "", usage like program -help [gives all the listof commands]
		flag.Parse()

		configPath = *flags
		//it is pointer so we need to dereference

		//if still you didn't get the path
		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	//we need to check the availability of the file on the provided path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	// now if file exist
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	//and on which structure we need to serialize we also need to give so here it is cfg [memory location of that]
	// if there exist some problem then it returns error

	if err != nil {
		log.Fatalf("can't read config files: %s", err.Error())
	}

	return &cfg
}

// and make sure when you are naming these type of things like MustLoad then make sure you are not returning the error from here, as name suggests it must load it is required for application to run
// as if this doesn't work then don'e even need to return error, just close the application from there itself
