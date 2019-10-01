package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const ENV_PREFIX = "WQH_"

func main() {
	newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "wqh",
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				varName := ENV_PREFIX + strings.ToUpper(f.Name)
				if val, ok := os.LookupEnv(varName); !f.Changed && ok {
					f.Value.Set(val)
				}
			})
		},
	}
	cmd.AddCommand(
		newRunCmd(),
		newConvertCmd(),
	)
	return cmd
}

func newRunCmd() *cobra.Command {
	var (
		saveFile   string
		headerFile string
	)
	cmd := &cobra.Command{
		Use:   "create <pic>",
		Short: "convert picture to text and run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				picture     io.Reader
				err         error
				pictureFile string    = args[0]
				output      io.Writer = os.Stdout
			)

			// output
			if saveFile != "" {
				output, err := os.Create(saveFile)
				output = io.MultiWriter(os.Stdout, save)
			}

			// write header
			if headerFile != "" {
				header, err := os.Open(headerFile)
				if err != nil {
					return err
				}
				_, err = io.Copy(output, header)
				if err != nil {
					return err
				}
			}

			// input
			if pictureFile == "-" {
				picture = os.Stdin
			} else {
				picture, err = os.Open(pictureFile)
				if err != nil {
					return err
				}
			}

			ctx := context.Background()

			text, err := convert(ctx, picture)
			if err != nil {
				return err
			}
			_, err = output.Write([]byte(text))
			return err
		},
	}
	cmd.Flags().StringVar(&saveFile, "save", "", "save created output to file")
	cmd.Flags().StringVar(&headerFile, "header", "", "header file which will be prepended to text form picture")
	return cmd
}

func newConvertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert <pic>",
		Short: "convert a picture to text and print text to stdout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				picture     io.Reader
				err         error
				pictureFile string = args[0]
			)

			if pictureFile == "-" {
				picture = os.Stdin
			} else {
				picture, err = os.Open(pictureFile)
				if err != nil {
					return err
				}
			}

			ctx := context.Background()

			text, err := convert(ctx, picture)
			if err != nil {
				return err
			}
			fmt.Println(text)
			return nil
		},
	}
	return cmd
}
