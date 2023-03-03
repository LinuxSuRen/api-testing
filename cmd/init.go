package cmd

import (
	"github.com/linuxsuren/api-testing/pkg/exec"
	"github.com/spf13/cobra"
)

type initOption struct {
	kustomization string
	waitNamespace string
	waitResource  string
}

func CreateInitCommand() (cmd *cobra.Command) {
	opt := &initOption{}
	cmd = &cobra.Command{
		Use:  "init",
		Long: "Support to init Kubernetes cluster with kustomization, and wait it with command: kubectl wait",
		RunE: opt.runE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.kustomization, "kustomization", "k", "", "The kustomization file path")
	flags.StringVarP(&opt.waitNamespace, "wait-namespace", "", "", "")
	flags.StringVarP(&opt.waitResource, "wait-resource", "", "", "")
	return
}

func (o *initOption) runE(cmd *cobra.Command, args []string) (err error) {
	if o.kustomization != "" {
		if err = exec.RunCommand("kubectl", "apply", "-k", o.kustomization, "--wait=true"); err != nil {
			return
		}
	}

	if o.waitNamespace != "" && o.waitResource != "" {
		if err = exec.RunCommand("kubectl", "wait", "-n", o.waitNamespace, o.waitResource, "--for", "condition=Available=True", "--timeout=900s"); err != nil {
			return
		}
	}
	return
}
