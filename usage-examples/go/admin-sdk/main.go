package testing

import (
	"admin-sdk/internal"
	"fmt"
	"os"
)

func main() {
	internal.LoadEnv()
	// Verify they are set
	fmt.Println("ATLAS_CLIENT_ID:", os.Getenv("ATLAS_CLIENT_ID"))
	fmt.Println("ATLAS_CLIENT_SECRET:", os.Getenv("ATLAS_CLIENT_SECRET"))

	// Now you can call LoadSecrets() without needing to manually set env vars
	secrets, _ := internal.LoadSecrets()
	fmt.Printf("Loaded Secrets: %+v\n", secrets)

	// Now you can call CreateAtlasClient() without needing to manually set env vars
	internal.CreateAtlasClient()
}
