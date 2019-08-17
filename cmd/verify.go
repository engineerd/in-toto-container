package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/engineerd/in-toto-container/pkg/docker"
	"github.com/engineerd/in-toto-container/pkg/intoto"
)

type verifyCmd struct {
	layout            string
	layoutKey         string
	linkDir           string
	targetFiles       []string
	verificationImage string
}

func newVerifyCmd() *cobra.Command {
	v := verifyCmd{}
	cmd := &cobra.Command{
		Use: "verify",
		RunE: func(cmd *cobra.Command, args []string) error {
			return v.run()
		},
	}
	cmd.Flags().StringVarP(&v.layout, "layout", "", "", "path to the root layout")
	cmd.Flags().StringVarP(&v.layoutKey, "layout-key", "", "", "path to the key used to sign the root layout")
	cmd.Flags().StringVarP(&v.linkDir, "links", "", "", "path to directory with the links")
	cmd.Flags().StringVarP(&v.verificationImage, "image", "", "", "container image to run the verification")
	cmd.Flags().StringArrayVarP(&v.targetFiles, "target", "", nil, "target files to copy in container")

	return cmd
}

func (v *verifyCmd) run() error {
	fmt.Printf("validating layout structure and signatures...\n")
	err := intoto.ValidateFromPath(v.layout)
	if err != nil {
		return err
	}

	fmt.Printf("running in-toto verifications in container based on image %v...\n", v.verificationImage)
	return docker.Run(v.verificationImage, v.layout, v.layoutKey, v.linkDir, v.targetFiles)
}
