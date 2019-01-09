package configrepo

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/plugins"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var SyntaxCmd = &cobra.Command{
	Use:   "syntax",
	Short: "Checks one or more definition files for syntactical correctness",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		syntax.Run(args)
	},
}

var syntax = &SyntaxRunner{}

type SyntaxRunner struct {
	Raw bool
}

func (sr *SyntaxRunner) Run(args []string) {
	if "" == PluginId {
		utils.DieLoudly(1, "You must provide a --plugin-id")
	}

	PluginJar = plugins.LocatePlugin(PluginId, PluginDir)

	cmdArgs := append([]string{"-jar", PluginJar, "syntax"}, args...)
	cmd := exec.Command("java", cmdArgs...)

	var success bool

	if sr.Raw {
		success = utils.ExecQ(cmd)
	} else {
		stdout := &strings.Builder{}
		stderr := &strings.Builder{}

		if success = utils.Exec(cmd, os.Stdin, stdout, stderr); success {
			// should really make plugins print machine-readable stdoutput instead of "OK"
			utils.Echofln(stdout.String())
		} else {
			resp := api.CrResponse{}

			if err := json.Unmarshal([]byte(stderr.String()), &resp); err != nil {
				utils.AbortLoudly(err)
			}

			utils.Echofln(resp.DisplayErrors())
		}
	}

	if !success {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(SyntaxCmd)
	SyntaxCmd.Flags().BoolVar(&syntax.Raw, "json", false, "machine-readable output (JSON)")
}
