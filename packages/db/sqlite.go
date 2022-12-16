//go:build ( netbsd || ios || freebsd || darwin || ( linux && riscv64 ) || ( linux && ppc64le )  || ( linux && s390x ) || ( linux && amd64 ) || ( linux && arm64 ) || ( linux && i386 ) || android || openbsd || windows || aix || illumos || solaris || plan9)
package Db

import ( _ "modernc.org/sqlite")
