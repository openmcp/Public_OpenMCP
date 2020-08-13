# omcplog

## Introduction of omcplog

> KETI에서 개발한 OpenMCP 플랫폼의 Log를 출력하기 위한 함수
>

## How to Use
omcplog를 Import후 함수를 호출해서 사용
```
package main

import (
   "openmcp/omcplog" //omcplog import
   "openmcp/util/controller/logLevel" //log Level 동적 변경을 위한 controller import
)

func main() {
  logLevel.KetiLogInit() //log init
  omcplog.V(0).Info("[OpenMCP] omcplog Test")
}
```
