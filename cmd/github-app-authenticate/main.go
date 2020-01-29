/*
Command 'github-app-authenticate' authenticates Github App by private keys and prints installation access token

  $ github-app-authenticate INTEGRATION_ID INSTALLATION_ID GITHUB_RSA_PRIVATE_KEY_PEM_PATH

To install, use go get,

  $ go get github.com/tcnksm/misc/cmd/github-app-authenticate

*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("[Usage] github-app-authenticate INTEGRATION_ID INSTALLATION_ID GITHUB_RSA_PRIVATE_KEY_PEM_PATH")
	}

	var (
		appIntegrationID  int64
		appInstallationID int64

		err error
	)
	appIntegrationIDstr, appInstallationIDstr, rsaPrivateKeyPemPath := os.Args[1], os.Args[2], os.Args[3]

	appIntegrationID, err = strconv.ParseInt(appIntegrationIDstr, 0, 64)
	if err != nil {
		log.Fatalf("[ERROR] INTEGRATION ID must be number: %s", err)
	}

	appInstallationID, err = strconv.ParseInt(appInstallationIDstr, 0, 64)
	if err != nil {
		log.Fatalf("[ERROR] INSTALLATION ID must be number: %s", err)
	}

	itr, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		appIntegrationID,
		appInstallationID,
		rsaPrivateKeyPemPath,
	)
	if err != nil {
		log.Fatalf("[ERRRO] Failed to create new trasport: %s\n", err)
	}

	token, err := itr.Token(context.Background())
	if err != nil {
		log.Fatalf("[ERRRO] Failed to get token: %s\n", err)
	}
	fmt.Printf("%s", token)
}
