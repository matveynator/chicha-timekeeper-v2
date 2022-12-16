//go:build ( dragonfly || ios || freebsd || darwin || ( linux && ppc64 ) || ( linux && ppc64le )  || ( linux && s390x ) || ( linux && amd64 ) || ( linux && mips64 ) || ( linux && mips64le ) || ( linux && arm64 ) || ( linux && i386 ) || android || windows || aix || illumos || solaris || plan9 )

package Db

import (  _ "github.com/genjidb/genji/driver" )
