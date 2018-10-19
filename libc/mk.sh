set -e
git checkout ../libc_linux_amd64.go
go install -v github.com/cznic/ccgo/v2/ccgo
rm -f log-ccgo
make distclean
make clean
./configure \
	CC=ccgo \
	CFLAGS='-D__typeof=typeof --ccgo-define-values' \
	--target=x86_64 |& tee log-configure
make AR=ar RANLIB=ranlib |& tee log-make
ccgo -ffreestanding -D_XOPEN_SOURCE=700 -I./arch/x86_64 -I./arch/generic -Iobj/src/internal \
	-I./src/internal -Iobj/include -I./include -D__typeof=typeof --ccgo-import os,runtime/debug,sync/atomic \
	--ccgo-pkg-name crt -o ../libc_linux_amd64.go ccgo.c lib/libc.a |& tee -a log-ccgo
go install -v github.com/cznic/ccgo/v2/ccgo
date
