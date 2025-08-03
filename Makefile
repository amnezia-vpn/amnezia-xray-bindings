ifneq ($(shell where cmd 2>nul || which cmd 2>/dev/null),)
	OS = windows
else
	OS = unix
endif

BUILD_DIR = build

SOURCES := main.go

LIB_A 	= amnezia_xray.a
LIB_DLL = amnezia_xray.dll
LIB_LIB = amnezia_xray.lib

LIB_HDR = amnezia_xray.h
LIB_DEF = amnezia_xray.def

ifeq ($(OS),windows)
all: $(BUILD_DIR)/$(LIB_LIB)
else
all: $(BUILD_DIR)/$(LIB_ARC)
endif

$(BUILD_DIR)/$(LIB_ARC): $(SOURCES)
	CGO_ENABLED=1 go build -ldflags=-w -o $(BUILD_DIR)/$(LIB_ARC) -buildmode=c-archive

$(BUILD_DIR)/$(LIB_DLL): $(SOURCES)
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -ldflags=-w -o $(BUILD_DIR)/$(LIB_DLL) -buildmode=c-shared

$(BUILD_DIR)/$(LIB_LIB): $(BUILD_DIR)/$(LIB_DLL)
	cd $(BUILD_DIR) && gendef $(LIB_DLL)
	cd $(BUILD_DIR) && dlltool -d $(LIB_DEF) -l $(LIB_LIB) -D $(LIB_DLL)
	rm $(BUILD_DIR)/$(LIB_DEF)

PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
