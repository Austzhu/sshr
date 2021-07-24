TOPDIR      = ${PWD}
SSHR        = ${TOPDIR}/cmd
OUTPUT      = ${TOPDIR}/_out
Q           = @
GO111MODULE = on

export TOPDIR SSHR OUTPUT GO111MODULE Q


all: sshr sshr_win64

prepare:
	${Q}mkdir -p ${OUTPUT}

sshr: prepare
	${Q}CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT}/sshr ${SSHR}

sshr_win64: prepare
	${Q}CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${OUTPUT}/sshr.exe ${SSHR}

clean:
	${Q}rm -rf ${OUTPUT}
