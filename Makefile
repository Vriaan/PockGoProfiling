TYPE ?=""
OPEN ?=""

BINARY_REPO := ./bin
BINARY_NAME := profilingSample
PPROF_FILES_DIR := ./prof
GENERATED_FROM_PPROF_FILES_DIR := ./files
CURRENT_DIR := $(GENERATED_FROM_PPROF_FILES_DIR)/$(TYPE)

#./bin/programToProfile -solus -memprofile=prof/samplesolusmemprofile.prof
build:
	go build -o  $(BINARY_REPO)/$(BINARY_NAME) cmd/manualProfiling/main.go

profiling:
	$(BINARY_REPO)/$(BINARY_NAME) -cpuprofile=prof/samplecpuprofile.prof
	$(BINARY_REPO)/$(BINARY_NAME) -solus -cpuprofile=prof/samplesoluscpuprofile.prof
	$(BINARY_REPO)/$(BINARY_NAME) -memprofile=prof/samplememprofile.prof
	$(BINARY_REPO)/$(BINARY_NAME) -solus -memprofile=prof/samplesolusmemprofile.prof
	$(BINARY_REPO)/$(BINARY_NAME) -trace=prof/sampletrace.prof
	$(BINARY_REPO)/$(BINARY_NAME) -solus -trace=prof/samplesolustrace.prof

get_profile:
ifeq ($(TYPE),"")
	@echo DO nothing if TYPE is not specified
else
	mkdir -p $(CURRENT_DIR)
	go tool pprof -$(TYPE) $(PPROF_FILES_DIR)/samplecpuprofile.prof > $(CURRENT_DIR)/samplecpuprofile.pdf
	go tool pprof -$(TYPE) $(PPROF_FILES_DIR)/samplesoluscpuprofile.prof > $(CURRENT_DIR)/samplesoluscpuprofile.pdf
	go tool pprof -$(TYPE) $(PPROF_FILES_DIR)/samplememprofile.prof > $(CURRENT_DIR)/samplememprofile.pdf
	go tool pprof -$(TYPE) $(PPROF_FILES_DIR)/samplesolusmemprofile.prof > $(CURRENT_DIR)/samplesolusmemprofile.pdf

	xterm -e "go tool trace $(PPROF_FILES_DIR)/sampletrace.prof"&
	xterm -e "go tool trace $(PPROF_FILES_DIR)/samplesolustrace.prof"&

	gvfs-open $(CURRENT_DIR)/sample*
endif


clean:
	rm -rf $(GENERATED_FROM_PPROF_FILES_DIR)/*  $(BINARY_REPO)/* $(PPROF_FILES_DIR)/*

all: clean build profiling get_profile

libprofile:
	go build -o $(BINARY_REPO)/profilingUsingLib cmd/libProfiling/mainProfileLib.go
	$(BINARY_REPO)/profilingUsingLib -cpuprofile=p
	$(BINARY_REPO)/profilingUsingLib -solus -cpuprofile=pp
	$(BINARY_REPO)/profilingUsingLib -memprofile=p
	$(BINARY_REPO)/profilingUsingLib -solus -memprofile=p

