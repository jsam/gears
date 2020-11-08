package version

import (
	"fmt"
	"runtime"
)

// The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// The main version number that is being run at the moment.
const Version = "0.1.0"

//BuildDate
var BuildDate = ""

//GoVersion
var GoVersion = runtime.Version()

//OsArch
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
