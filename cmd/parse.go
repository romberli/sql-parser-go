/*
Copyright © 2020 Romber Li <romber2001@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/sql-parser-go/config"
	"github.com/romberli/sql-parser-go/pkg/lexer"
	"github.com/romberli/sql-parser-go/pkg/message"
	"github.com/romberli/sql-parser-go/pkg/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// parseCmd represents the start command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse command",
	Long:  `use parse to match input string`,
	Run: func(cmd *cobra.Command, args []string) {
		// init config
		err := initConfig()
		if err != nil {
			fmt.Println(message.NewMessage(message.ErrInitConfig, err.Error()).Error())
			os.Exit(constant.DefaultAbnormalExitCode)
		}
		var l *lexer.Lexer
		var p *parser.Parser

		fa := viper.GetString(config.FiniteAutomataKey)
		switch fa {
		case config.NFA:
			l = lexer.NewLexer(lexer.NewNFAWithDefault())
			tokens := l.Lex(viper.GetString(config.SQLKey))
			p = parser.NewParser(parser.NewNFA(tokens))
		case config.LLOne:
			l = lexer.NewLexer(lexer.NewDFAWithDefault())
			tokens := l.Lex(viper.GetString(config.SQLKey))
			p = parser.NewParser(parser.NewLLOne(tokens))
		default:
			fmt.Println(message.NewMessage(message.ErrNotValidFiniteAutomata, viper.GetString(config.FiniteAutomataKey)).Error())
			os.Exit(constant.DefaultAbnormalExitCode)
		}

		astNode, err := p.Parse()
		if err != nil {
			fmt.Println(err.Error())
		}
		astNode.PrintChildren()
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
