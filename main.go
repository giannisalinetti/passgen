// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sethvargo/go-password/password"
)

// genInput defines the allowed characters to use"
var genInput = &password.GeneratorInput{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"1234567890",
	"!@#$%^&*()_+{}:;.<>?/|",
}

// Config defines the configuration parameters
type Config struct {
	length        int
	numDigits     int
	numSymbols    int
	noUpper       bool
	allowRepeat   bool
	numIterations int
	json          bool
}

// NewConfig is the default constructor for the Config type
func NewConfig() *Config {
	return &Config{
		32,    // Default number of characters
		10,    // Default number of digits
		10,    // Default number of symbols
		false, // Avoid uppercase bool
		false, // Allow repetitions bool
		1,     // Default number of iterations
		false, // Print in json format
	}
}

// initParams initializes the configuration parameters
func initParams(r *http.Request) (*Config, error) {

	c := NewConfig()
	q := r.URL.Query()

	values, _ := q["length"]
	if len(values) != 0 {
		c.length, _ = strconv.Atoi(values[0])
	}

	values, _ = q["digits"]
	if len(values) != 0 {
		c.numDigits, _ = strconv.Atoi(values[0])
	}

	values, _ = q["symbols"]
	if len(values) != 0 {
		c.numSymbols, _ = strconv.Atoi(values[0])
	}

	values, _ = q["noupper"]
	if len(values) != 0 {
		c.noUpper, _ = strconv.ParseBool(values[0])
	}

	values, _ = q["allowrepeat"]
	if len(values) != 0 {
		c.allowRepeat, _ = strconv.ParseBool(values[0])
	}

	values, _ = q["iterations"]
	if len(values) != 0 {
		c.numIterations, _ = strconv.Atoi(values[0])
	}

	values, _ = q["json"]
	if len(values) != 0 {
		c.json, _ = strconv.ParseBool(values[0])
	}

	return c, nil
}

// makePasswdSlice generates the slice of passwords
func makePasswdSlice(gen *password.Generator, cfg *Config) ([]string, error) {
	p := make([]string, 0)
	for i := 0; i < cfg.numIterations; i++ {
		res, err := gen.Generate(cfg.length, cfg.numDigits, cfg.numSymbols, cfg.noUpper, cfg.allowRepeat)
		if err != nil {
			return nil, err
		}
		p = append(p, res)
	}
	return p, nil
}

// jsonPrinter handles the json formatting of the output
func jsonPrinter(s []string) ([]byte, error) {
	p := make(map[string]string)
	for i, v := range s {
		p["Password"+strconv.Itoa(i)] = v
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// passwdFunc is the handler function for password generation
func passwdFunc(w http.ResponseWriter, r *http.Request) {

	cfg, err := initParams(r)
	if err != nil {
		log.Fatal(err)
	}

	gen, err := password.NewGenerator(genInput)
	if err != nil {
		log.Fatal(err)
	}

	passwdSlice, err := makePasswdSlice(gen, cfg)
	if err != nil {
		log.Println(err) // Don't exit on password generation errors
		fmt.Fprintf(w, "%v\n", err)
	} else {
		if cfg.json {
			passwdJson, err := jsonPrinter(passwdSlice)
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, "JSON encoding error: %v\n", err)
			}
			log.Printf("Client request: %d new password(s) generated for host %s\n", cfg.numIterations, r.RemoteAddr)
			fmt.Fprintf(w, string(passwdJson))
		} else {
			log.Printf("Client request: %d new password(s) generated for host %s\n", cfg.numIterations, r.RemoteAddr)
			for index, value := range passwdSlice {
				fmt.Fprintf(w, "Password%d: %s\n", index, value)
			}
		}
	}
}

// helpFunc prints an help about allowed parameters
func helpFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Password generator service parameters:\n\n")
	fmt.Fprintf(w, "length:\t\tPassword length. Default 32\n")
	fmt.Fprintf(w, "digits:\t\tNumber of ditigs to use. Default: 10\n")
	fmt.Fprintf(w, "symbols:\tNumber of symbols to use. Default: 10\n")
	fmt.Fprintf(w, "noupper:\tForbid uppercase letters. Default: false\n")
	fmt.Fprintf(w, "allowrepeat:\tAllow repetitions. Default: false\n")
	fmt.Fprintf(w, "iterations:\tNumber of passwords to print. Default 1\n")
	fmt.Fprintf(w, "json:\t\tReturn output in json format. Default false\n")
}

// healthFunc return an HTTP 200 status for liveness probes
func healthFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: OK\n")
}

// verifyCerts tests if the certificate and key file exist
func verifyCerts(crt, key string) error {
	// Verify if the certificate exists
	_, err := os.Stat(crt)
	if err != nil {
		return err
	}

	// Verify if the key exists
	_, err = os.Stat(key)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Define the listening port with a flag
	var (
		port       string
		serverCert string
		serverKey  string
	)

	flag.StringVar(&port, "port", "8443", "Listening port")
	flag.StringVar(&serverCert, "crt", "/etc/passgen/certs/server.crt", "Path to server certificate")
	flag.StringVar(&serverKey, "key", "/etc/passgen/certs/server.key", "Path to private key")
	flag.Parse()

	// Test if provided certificates exist
	err := verifyCerts(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	// Declare the HandleFuncs
	http.HandleFunc("/passwd", passwdFunc)
	http.HandleFunc("/help", helpFunc)
	http.HandleFunc("/health", healthFunc)

	// Start a os.Signal channel to accept signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-sigs
		log.Printf("Shutting down password generator service\n")
		os.Exit(0)
	}()

	log.Printf("Starting password generator service on port %s/tcp\n", port)
	log.Fatal(http.ListenAndServeTLS(":"+port, serverCert, serverKey, nil))
}
