package shell

import _ "embed"

//go:embed install.sh
var BashScript []byte

//go:embed install.ps1
var PowerShellScript []byte
