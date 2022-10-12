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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lexCmd represents the start command
var lexCmd = &cobra.Command{
	Use:   "lex",
	Short: "lex command",
	Long:  `use lex to match input string`,
	Run: func(cmd *cobra.Command, args []string) {
		// init config
		err := initConfig()
		if err != nil {
			fmt.Println(message.NewMessage(message.ErrInitConfig, err.Error()).Error())
			os.Exit(constant.DefaultAbnormalExitCode)
		}

		var l *lexer.Lexer

		fa := viper.GetString(config.LexFiniteAutomataKey)
		switch fa {
		case config.NFA:
			l = lexer.NewLexer(lexer.NewNFAWithDefault())
		case config.DFA:
			l = lexer.NewLexer(lexer.NewDFAWithDefault())
		default:
			fmt.Println(message.NewMessage(message.ErrNotValidLexFiniteAutomata, viper.GetString(config.LexFiniteAutomataKey)).Error())
			os.Exit(constant.DefaultAbnormalExitCode)
		}

		tokens := l.Lex(viper.GetString(config.SQLKey))

		for _, token := range tokens {
			fmt.Println(token.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(lexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lexCmd.PersistentFlags().String("foo", "", "A help for foo")
	// finite automata

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	lexCmd.Flags().StringVar(&lexFiniteAutomata, "finite-automata", constant.DefaultRandomString, fmt.Sprintf("specify the finite automata(default: %s)", config.DefaultLexFiniteAutomata))
}
