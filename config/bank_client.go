package config

import "os"

func GetBankClientURL() string {
	url, ok := os.LookupEnv("BANK_CLIENT_URL")
	if !ok {
		url = "https://65d38c2f522627d501091ff2.mockapi.io"
	}
	return url

}
