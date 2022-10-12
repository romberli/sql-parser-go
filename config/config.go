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

package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/go-multierror"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/romberli/sql-parser-go/pkg/message"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var (
	ValidLogLevels                 = []string{"debug", "info", "warn", "warning", "error", "fatal"}
	ValidLogFormats                = []string{"text", "json"}
	ValidLexFiniteAutomata         = []string{NFA, DFA}
	ValidParseParserFiniteAutomata = []string{NFA, LL}
)

// SetDefaultConfig set default configuration, it is the lowest priority
func SetDefaultConfig(baseDir string) {
	// log
	defaultLogFile := filepath.Join(baseDir, DefaultLogDir, log.DefaultLogFileName)
	viper.SetDefault(LogFileNameKey, defaultLogFile)
	viper.SetDefault(LogLevelKey, log.DefaultLogLevel)
	viper.SetDefault(LogFormatKey, log.DefaultLogFormat)
	viper.SetDefault(LogMaxSizeKey, log.DefaultLogMaxSize)
	viper.SetDefault(LogMaxDaysKey, log.DefaultLogMaxDays)
	viper.SetDefault(LogMaxBackupsKey, log.DefaultLogMaxBackups)
	// lex
	viper.SetDefault(LexFiniteAutomataKey, DefaultLexFiniteAutomata)
	// parse
	viper.SetDefault(ParseLexerFiniteAutomataKey, DefaultParseLexerFiniteAutomata)
	viper.SetDefault(ParseParserFiniteAutomataKey, DefaultParseParserFiniteAutomata)
}

// ValidateConfig validates if the configuration is valid
func ValidateConfig() (err error) {
	merr := &multierror.Error{}

	// validate log section
	err = ValidateLog()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	// validate lex
	err = ValidateLex()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	// validate parse
	err = ValidateParse()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	// validate sql
	err = ValidateSQL()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	return merr.ErrorOrNil()
}

// ValidateLog validates if log section is valid.
func ValidateLog() error {
	var valid bool

	merr := &multierror.Error{}

	// validate log.FileName
	logFileName, err := cast.ToStringE(viper.Get(LogFileNameKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	}
	logFileName = strings.TrimSpace(logFileName)
	if logFileName == constant.EmptyString {
		merr = multierror.Append(merr, message.NewMessage(message.ErrEmptyLogFileName))
	}
	isAbs := filepath.IsAbs(logFileName)
	if !isAbs {
		logFileName, err = filepath.Abs(logFileName)
		if err != nil {
			merr = multierror.Append(merr, err)
		}
	}
	valid, _ = govalidator.IsFilePath(logFileName)
	if !valid {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogFileName, logFileName))
	}

	// validate log.level
	logLevel, err := cast.ToStringE(viper.Get(LogLevelKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else {
		valid, err = common.ElementInSlice(ValidLogLevels, logLevel)
		if err != nil {
			merr = multierror.Append(merr, err)
		} else if !valid {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogLevel, logLevel))
		}
	}

	// validate log.format
	logFormat, err := cast.ToStringE(viper.Get(LogFormatKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else {
		valid, err = common.ElementInSlice(ValidLogFormats, logFormat)
		if err != nil {
			merr = multierror.Append(merr, err)
		} else if !valid {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogFormat, logFormat))
		}
	}

	// validate log.maxSize
	logMaxSize, err := cast.ToIntE(viper.Get(LogMaxSizeKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else if logMaxSize < MinLogMaxSize || logMaxSize > MaxLogMaxSize {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogMaxSize, MinLogMaxSize, MaxLogMaxSize, logMaxSize))
	}

	// validate log.maxDays
	logMaxDays, err := cast.ToIntE(viper.Get(LogMaxDaysKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else if logMaxDays < MinLogMaxDays || logMaxDays > MaxLogMaxDays {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogMaxDays, MinLogMaxDays, MaxLogMaxDays, logMaxDays))
	}

	// validate log.maxBackups
	logMaxBackups, err := cast.ToIntE(viper.Get(LogMaxBackupsKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else if logMaxBackups < MinLogMaxDays || logMaxBackups > MaxLogMaxDays {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogMaxBackups, MinLogMaxBackups, MaxLogMaxBackups, logMaxBackups))
	}

	return merr.ErrorOrNil()
}

func ValidateLex() error {
	var valid bool

	merr := &multierror.Error{}

	// validate lex.finiteAutomata
	fa, err := cast.ToStringE(viper.Get(LexFiniteAutomataKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else {
		valid, err = common.ElementInSlice(ValidLexFiniteAutomata, fa)
		if err != nil {
			merr = multierror.Append(merr, err)
		} else if !valid {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLexFiniteAutomata, valid))
		}
	}

	return merr.ErrorOrNil()
}

func ValidateParse() error {
	var valid bool

	merr := &multierror.Error{}

	// validate parse.lexer.finiteAutomata
	lexerFA, err := cast.ToStringE(viper.Get(ParseLexerFiniteAutomataKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else {
		valid, err = common.ElementInSlice(ValidLexFiniteAutomata, lexerFA)
		if err != nil {
			merr = multierror.Append(merr, err)
		} else if !valid {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidParseLexerFiniteAutomata, lexerFA))
		}
	}

	// validate parse.parserFiniteAutomata
	parserFA, err := cast.ToStringE(viper.Get(ParseParserFiniteAutomataKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	} else {
		valid, err = common.ElementInSlice(ValidParseParserFiniteAutomata, parserFA)
		if err != nil {
			merr = multierror.Append(merr, err)
		} else if !valid {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidParseParserFiniteAutomata, parserFA))
		}
	}

	return merr.ErrorOrNil()
}

func ValidateSQL() error {
	merr := &multierror.Error{}

	_, err := cast.ToStringE(viper.Get(SQLKey))
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	return merr.ErrorOrNil()
}

// TrimSpaceOfArg trims spaces of given argument
func TrimSpaceOfArg(arg string) string {
	args := strings.SplitN(arg, constant.EqualString, 2)

	switch len(args) {
	case 1:
		return strings.TrimSpace(args[0])
	case 2:
		argName := strings.TrimSpace(args[0])
		argValue := strings.TrimSpace(args[1])
		return fmt.Sprintf("%s=%s", argName, argValue)
	default:
		return arg
	}
}
