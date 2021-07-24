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
	${Q}GOHOSTOS=linux go build -o ${OUTPUT}/sshr ${SSHR}

sshr_win64: prepare
	${Q}GOHOSTOS=windows go build -o ${OUTPUT}/sshr.exe ${SSHR}

clean:
	${Q}rm -rf ${OUTPUT}
