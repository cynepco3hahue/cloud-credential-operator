package provisioning

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type options struct {
	TargetDir           string
	PublicKeyPath       string
	Region              string
	NamePrefix          string
	CredRequestDir      string
	IdentityProviderARN string
	DryRun              bool
}

var (
	// CreateOpts captures the options that affect creation/updating
	// of the generated objects.
	CreateOpts = options{
		TargetDir: "",
	}
)

// NewCreateCmd implements the "create" command for the credentials provisioning
func NewCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:              "create",
		Short:            "Create credentials objects",
		Long:             "Creating/updating/generating objects related to cloud credentials",
		PersistentPreRun: initEnv,
	}

	createCmd.AddCommand(NewKeyProvision())
	createCmd.AddCommand(NewIdentityProviderSetup())
	createCmd.AddCommand(NewIAMRolesSetup())

	createCmd.AddCommand(NewAllSetup())

	return createCmd
}

// initEnv will ensure the destination directory is ready to receive the generated
// files, and will create the directory if necessary.
func initEnv(cmd *cobra.Command, args []string) {
	if CreateOpts.TargetDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %s", err)
		}

		CreateOpts.TargetDir = pwd
	}

	fPath, err := filepath.Abs(CreateOpts.TargetDir)
	if err != nil {
		log.Fatalf("Failed to resolve full path: %s", err)
	}

	// create target dir if necessary
	ensureDir(fPath)

	// create manifests dir if necessary
	manifestsDir := filepath.Join(fPath, manifestsDirName)
	ensureDir(manifestsDir)

	tlsDir := filepath.Join(fPath, tlsDirName)
	ensureDir(tlsDir)
}

func ensureDir(path string) {
	sResult, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, 0700); err != nil {
			log.Fatalf("Failed to create directory: %s", err)
		}
		sResult, err = os.Stat(path)
	} else if err != nil {
		log.Fatalf("Failed to stat: %+v", err)
	}

	if !sResult.IsDir() {
		log.Fatalf("File %s exists and is not a directory", path)
	}
}
