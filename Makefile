################################################################################
#
# THE NEGRONI-WARE LICENSE [Derived from "THE BEER-WARE LICENSE" (Revision 42)]
#
# Davide "FunkyAss" Del Zompo wrote this file. As long as you retain this notice
# you can do whatever you want with this stuff. If we meet some day, and you
# think this stuff is worth it, you can buy me a Negroni cocktail in return.
#
# Davide "FunkyAss" Del Zompo
#
################################################################################

NETWORK_BIN		= network
NETWORK_ROOT	= .
NETWORK_SRC		= $(addprefix ${NETWORK_ROOT}/,main.go)
NETWORK_ENV		=

MANAGER_BIN		= manager
MANAGER_ROOT	= ./cmd/manager
MANAGER_SRC		= $(addprefix ${MANAGER_ROOT}/,bluetooth.go main.go network.go)
MANAGER_ENV		=


#### cross-compilation options #################################################
TARGET_CONF_DIR	= .targets

ifneq (${TARGET},)
ifeq ($(shell test -f ${TARGET_CONF_DIR}/${TARGET}.mk && printf "yes"), yes)
include ${TARGET_CONF_DIR}/${TARGET}.mk
else
$(error Unknown cross-compile target: "${TARGET}")
endif
endif
################################################################################

.PHONY: all
all: ${NETWORK_BIN} ${MANAGER_BIN}

${NETWORK_BIN}: ${NETWORK_SRC}
	${NETWORK_ENV} go build -o $@ ${NETWORK_ROOT}

${MANAGER_BIN}: ${MANAGER_SRC}
	${MANAGER_ENV} go build -o $@ ${MANAGER_ROOT}

.PHONY: clean
clean: clean-network clean-manager

.PHONY: clean-network
clean-network:
	rm -f ${NETWORK_BIN}

.PHONY: clean-manager
clean-manager:
	rm -f ${MANAGER_BIN}
