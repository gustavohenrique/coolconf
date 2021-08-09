# CoolConf

[![Coverage Status](https://coveralls.io/repos/github/gustavohenrique/coolconf/badge.svg?branch=main)](https://coveralls.io/github/gustavohenrique/coolconf?branch=main)

CoolConf shares information across all packages. You can use your custom struct with annotated fields to store da read from environment variables, YAML files, or encrypted strings.

## Installation

```sh
go get github.com/gustavohenrique/coolconf/.../
```

## Usage

### Reading data

#### Environment variables

Initial values:

```sh
export SOME_INT=9
export SOME_STR=hi
export SOME_BOOL=true
export DATABASE_URL_DEV="postgres://user:pass@dev"
export STAGING__DATABASE_URL="postgres://user:pass@staging"
```

The configuration struct must have the `env` annotation with the envvar's name.

```go
package main
import (
    "fmt"
    "github.com/gustavohenrique/coolconf"
)

type MyConfig struct {
    Number int    `env:"SOME_INT"`
    Text   string `env:"SOME_STR"`
    Yes    bool   `env:"SOME_BOOL"`
}

func init() {
    coolconf.New()
}

func main() {
    var myConfig MyConfig
    coolconf.Load(&myConfig)
    fmt.Println(myConfig.Text)
    fmt.Println(myConfig.Number)
    fmt.Println(myConfig.Yes)
}
```

You can use a second argument to add a suffix or prefix with _ as the default separator.

```go
type MyConfig struct {
    DatabaseURL string `env:"DATABASE_URL"`
}

func init() {
    coolconf.New()
}

func main() {
    var myConfig MyConfig
    coolconf.Load(&myConfig, "DEV")
    fmt.Println(myConfig.DatabaseURL)
    // output is postgres://user:pass@dev

    coolconf.New(coolconf.Settings{
        Source: coolconf.ENV,
        Env: coolconf.Option{
            UseGroupAsPrefix: true,
            Separator: "__",
        },
    })
    coolconf.Load(&myConfig, "STAGING")
    fmt.Println(myConfig.DatabaseURL)
    // output is postgres://user:pass@staging
}
```

#### YAML file

File must end with .yaml or .yml.

```sh
echo "database_url: postgres://user:pass@localhost" > /tmp/config.yaml
echo "database_url: postgres://user:pass@dev" > /tmp/config_dev.yaml
```

The configuration struct must have the `YAML` annotation with the envvar's name.

```go
package main
import (
    "fmt"
    "github.com/gustavohenrique/coolconf"
)

type MyConfig struct {
    DatabaseURL string `yaml:"database_url"`
}

func init() {
    coolconf.New(coolconf.Settings{
        Source: "/tmp/config.yaml",
    })
}

func main() {
    var myConfig MyConfig
    coolconf.Load(&myConfig)
    fmt.Println(myConfig.DatabaseURL)
    // output: postgres://user:pass@localhost
}
```
The second argument is the suffix or prefix with _ as default separator to the filename.

```go
func main() {
    var myConfig MyConfig
    coolconf.Load(&myConfig, "dev")
    fmt.Println(myConfig.DatabaseURL)
    // output: postgres://user:pass@dev
}
```

You can configure to create the YAML file using default values, if the file does not exists:

```go
package main
import (
    "fmt"
    "github.com/gustavohenrique/coolconf"
)
type MyConfig struct {
    Host string `yaml:"host" default:"localhost"`
    Port int    `yaml:"port" default:"8080"`
}
func main() {
    coolconf.New(coolconf.Settings{
        Source: "/tmp/server.yaml",
        ShouldCreateDefaultYaml: true,
    })
    var myConfig MyConfig
    coolconf.Load(&myConfig)
    fmt.Println(myConfig.Host, ":", myConfig.Port)
    // output: localhost:8080
}
```

#### Encrypted string

CoolConf provides the `cccli` tool to allow you to encrypt/decrypt files using [AES256-GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode).  
You can encrypt a YAML file and store it in an S3 bucket and decrypt this file on the application's boot, for example.

**Be careful**: keep your password safe and do not allow public access to your encrypted file.

```sh
echo "database_url: postgres://user:pass@localhost" > /tmp/config.yaml

# generate an encrypted file
cccli -encrypt /tmp/config.yaml -output /tmp/encrypted.yaml -secret mystrongpassword

# reverse the file to an decrypted state
cccli -decrypt /tmp/encrypted.yaml -output /tmp/config.yaml -secret mystrongpassword
```

This example read the hex string encrypted from the file:

```go
package main
import (
    "fmt"
    "log"
    "io/ioutil"

    "github.com/gustavohenrique/coolconf"
)

type MyConfig struct {
    DatabaseURL string `yaml:"database_url"`
}

func init() {
    coolconf.New(coolconf.Settings{
        Encrypted: true,
        SecretKey: "mystrongpass",
    })
}

func main() {
    bytes, _ := ioutil.ReadFile("/tmp/encrypted.yaml")
    err := coolconf.DecryptYaml(bytes)
    if err != nil {
        log.Fatalln(err)
    }
    var myConfig MyConfig
    coolconf.Load(&myConfig)
    fmt.Println(myConfig.DatabaseURL)
    // output: postgres://user:pass@localhost
}
```

### Reset configuration

You can reset the configuration clearing the old data and loading the new one.  
For web applications, you can add a secret endpoint to do it, allowing your application gets a new configuration without a redeploy.

```go
package main

import (
    "os"
    "log"
    "net/http"
    "github.com/gustavohenrique/coolconf"
)

type MyConfig struct {
    DatabaseURL string `env:"DATABASE_URL"`
}

func init() {
    coolconf.New()
    os.Setenv("DATABASE_URL", "postgres://user:pass@localhost")
}

func view(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        var myConfig MyConfig
        coolconf.Load(&myConfig)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("value=" + myConfig.DatabaseURL))
    }
}

func reload(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        coolconf.Clear()
        os.Setenv("DATABASE_URL", "postgres://admin:admin@127.0.0.1")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("reloaded!"))
    }
}

func main() {
    http.HandleFunc("/", view)
    http.HandleFunc("/reload", reload)
    log.Println("Listening :8080")
    log.Fatal(http.ListenAndServe(":8080",nil))
}
```

See the magic happening:

```sh
curl http://localhost:8080/
value=postgres://user:pass@localhost

curl http://localhost:8080/reload
reloaded!

curl http://localhost:8080/
value=postgres://admin:admin@127.0.0.1
```

## License

Copyright 2021 Gustavo Henrique

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
