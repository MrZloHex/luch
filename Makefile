# Makefile
WH_DIR=third_party/whisper.cpp
WH_BUILD=$(WH_DIR)/build

ABS_WH_DIR := $(abspath $(WH_DIR))
ABS_BUILD  := $(abspath $(WH_BUILD))

# Compile optimization for cgo C/C++
COPT ?= -O2 -DNDEBUG

.PHONY: whisper
whisper:
	@test -d $(WH_BUILD) || cmake -S $(WH_DIR) -B $(WH_BUILD) \
		-DCMAKE_BUILD_TYPE=Release \
		-DBUILD_SHARED_LIBS=OFF -DGGML_STATIC=ON \
		-DWHISPER_BUILD_EXAMPLES=OFF -DWHISPER_BUILD_TESTS=OFF
	cmake --build $(WH_BUILD) -j

# Include dirs for whisper headers
CGO_CFLAGS_COMMON   := $(COPT) -I$(ABS_WH_DIR)/include -I$(ABS_WH_DIR)/ggml/include
CGO_CXXFLAGS_COMMON := $(COPT)
# Library search paths + static archives (order matters)
CGO_LDFLAGS_COMMON  := \
	-L$(ABS_BUILD)/src -L$(ABS_BUILD)/ggml/src \
	$(ABS_BUILD)/src/libwhisper.a \
	$(ABS_BUILD)/ggml/src/libggml.a \
	$(ABS_BUILD)/ggml/src/libggml-base.a \
	$(ABS_BUILD)/ggml/src/libggml-cpu.a \
	-lpthread -lm -lstdc++

.PHONY: build
build: whisper
	CGO_CFLAGS='$(CGO_CFLAGS_COMMON)' CGO_CXXFLAGS='$(CGO_CXXFLAGS_COMMON)' CGO_LDFLAGS='$(CGO_LDFLAGS_COMMON)' \
		go build -o bin/luch ./cmd/luch

