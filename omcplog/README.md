# omcplog

## Introduction of omcplog

> Function for printing Log of OpenMCP platform developed by KETI
>

## How to Use
import omcplog and call function to use
> Enable Static Log Level
```
package main

import (
   "flag"
   "openmcp/omcplog" //omcplog import
)

func main() {
   omcplog.InitFlags(nil)
   flag.Set("omcpv", "0") //Log Level Init
   flag.Parse()
   omcplog.V(0).Info("[OpenMCP] omcplog Test")
}
```
> Enable Dynamic Log Level
```
package main

import (
   "openmcp/omcplog" //omcplog import
   "openmcp/util/controller/logLevel" //controller import for dynamic change of log level
)

func main() {
   logLevel.KetiLogInit() //log init
   omcplog.V(0).Info("[OpenMCP] omcplog Test")
}
```
