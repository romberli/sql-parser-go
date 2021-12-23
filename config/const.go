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
	"github.com/romberli/go-util/constant"
)

// global constant
const (
	DefaultCommandName        = "go-template"
	DefaultErrorHeader        = "GO-TEMPLATE"
	DefaultDaemon             = false
	DefaultBaseDir            = constant.CurrentDir
	DefaultLogDir             = "./log"
	MinLogMaxSize             = 1
	MaxLogMaxSize             = constant.MaxInt
	MinLogMaxDays             = 1
	MaxLogMaxDays             = constant.MaxInt
	MinLogMaxBackups          = 1
	MaxLogMaxBackups          = constant.MaxInt
	DefaultServerAddr         = "0.0.0.0:80"
	DefaultServerReadTimeout  = 5
	DefaultServerWriteTimeout = 10
	MinServerReadTimeout      = 0
	MaxServerReadTimeout      = 60
	MinServerWriteTimeout     = 1
	MaxServerWriteTimeout     = 60
	DaemonArgTrue             = "--daemon=true"
	DaemonArgFalse            = "--daemon=false"
)

// configuration constant
const (
	ConfKey               = "config"
	DaemonKey             = "daemon"
	LogFileNameKey        = "log.fileName"
	LogLevelKey           = "log.level"
	LogFormatKey          = "log.format"
	LogMaxSizeKey         = "log.maxSize"
	LogMaxDaysKey         = "log.maxDays"
	LogMaxBackupsKey      = "log.maxBackups"
	ServerAddrKey         = "server.addr"
	ServerPidFileKey      = "server.pidFile"
	ServerReadTimeoutKey  = "server.readTimeout"
	ServerWriteTimeoutKey = "server.writeTimeout"
)
