package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	awscloud "github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/scenario"
	"github.com/wallix/awless/scenario/driver"
	"github.com/wallix/awless/scenario/driver/aws"
)

func init() {
	createCmd.AddCommand(createInstanceCmd)
	createCmd.AddCommand(createAliasCmd)

	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create various type of resources by id: users, groups, instances, vpcs, ...",
}

var createAliasCmd = &cobra.Command{
	Use:   "alias [name] [alias of]",
	Short: "Create alias",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("Not enough args, need two args, alias name and the resource id\n")
		}
		if len(args) > 2 {
			return fmt.Errorf("Two many args, need two args, alias name and the resource id\n")
		}
		return statsDB.AddAlias(args[0], args[1])
	},
}

var createInstanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"inst", "i"},
	Short:   "Create an instance",

	RunE: func(cmd *cobra.Command, args []string) error {
		var buff bytes.Buffer

		buff.WriteString("CREATE INSTANCE")

		var count int
		fmt.Print("Number of instances? ")
		_, err := fmt.Scanln(&count)
		if err != nil {
			return err
		}
		buff.WriteString(fmt.Sprintf(" COUNT %d", count))

		types := []string{
			"t2.nano:   vCPU=1, CPU/hour=3, Mem Gio=0,5, EBS only",
			"t2.micro:  vCPU=1, CPU/hour=6, Mem Gio=1, EBS only",
			"t2.small:  vCPU=1, CPU/hour=12, Mem Gio=2, EBS only",
			"t2.medium: vCPU=2, CPU/hour=24, Mem Gio=4, EBS only",
			"t2.large:  vCPU=2, CPU/hour=36, Mem Gio=8, EBS only",
			"t2.xlarge: vCPU=4, CPU/hour=54, Mem Gio=16, EBS only",
			"t2.2xlarge: vCPU=8, CPU/hour=81, Mem Gio=32, EBS only",
		}

		var typ int
		fmt.Println()
		for index, typ := range types {
			fmt.Printf("%d. %s\n", index+1, typ)
		}
		fmt.Print("\nType of instance? ")

		_, err = fmt.Scanln(&typ)
		if err != nil {
			return err
		}

		mytype := strings.Split(types[typ], ":")[0]

		buff.WriteString(fmt.Sprintf(" TYPE %s", mytype))

		scen := buff.String()

		var yesorno string
		fmt.Print("\nDone\n\n")
		fmt.Print(scen)
		fmt.Print("\n\nAbout to run? (y/n): ")
		_, err = fmt.Scanln(&yesorno)
		if err != nil {
			return err
		}

		if strings.TrimSpace(yesorno) == "y" {
			lex := &scenario.Lexer{}
			scen := lex.ParseScenario(scen)

			awsDriver := aws.NewDriver(awscloud.InfraService)
			awsDriver.SetLogger(log.New(os.Stdout, "[aws driver] ", log.Ltime))

			runner := &driver.Runner{Driver: awsDriver}

			return runner.Run(scen)
		}

		return nil
	},
}