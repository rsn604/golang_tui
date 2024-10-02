SHELL=/bin/bash
GOBUILD=go build
MAIN_PROG=listdbg
LOAD_PROG=loaddb
MAIN=cmd/main.go
LOAD=cmd/loaddb.go
TARGET=target

# ---------------------------------------------------
build:
	@mkdir -p $(TARGET)
	$(GOBUILD) -o $(MAIN_PROG) $(MAIN)
	$(GOBUILD) -o $(LOAD_PROG) $(LOAD)
	@cp scripts/linux/listdbg.sh .
	@cp scripts/linux/loaddb.sh .
	tar czvf $(TARGET)/$(MAIN_PROG).tgz $(MAIN_PROG) $(LOAD_PROG) csv db listdbg.sh loaddb.sh

# ---------------------------------------------------
build-win64:
	@mkdir -p $(TARGET)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(MAIN_PROG)64.exe $(MAIN)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(LOAD_PROG)64.exe $(LOAD) 
	@cp scripts/win64/listdbg64.bat .
	@cp scripts/win64/loaddb64.bat .
	tar czvf $(TARGET)/$(MAIN_PROG)_win64.tgz $(MAIN_PROG)64.exe $(LOAD_PROG)64.exe csv db listdbg64.bat loaddb64.bat
	@rm listdbg64.bat
	@rm loaddb64.bat
	@rm loaddb64.exe
	@rm listdbg64.exe

build-win32:
	@mkdir -p $(TARGET)
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(MAIN_PROG)32.exe $(MAIN)
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(LOAD_PROG)32.exe $(LOAD)
	@cp scripts/win32/listdbg32.bat .
	@cp scripts/win32/loaddb32.bat .
	tar czvf $(TARGET)/$(MAIN_PROG)_win32.tgz $(MAIN_PROG)32.exe $(LOAD_PROG)32.exe csv db listdbg32.bat loaddb32.bat
	@rm listdbg32.bat
	@rm loaddb32.bat
	@rm loaddb32.exe
	@rm listdbg32.exe

# ---------------------------------------------------
build-arm6:
	@mkdir -p $(TARGET)
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(MAIN_PROG)_arm6 $(MAIN)
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(LOAD_PROG)_arm6 $(LOAD)
	@cp scripts/linux_arm6/listdbg_arm6.sh .
	@cp scripts/linux_arm6/loaddb_arm6.sh .
	tar czvf $(TARGET)/$(MAIN_PROG)_arm6.tgz $(MAIN_PROG)_arm6 $(LOAD_PROG)_arm6 csv db listdbg_arm6.sh loaddb_arm6.sh
	@rm listdbg_arm6.sh
	@rm loaddb_arm6.sh
	@rm listdbg_arm6
	@rm loaddb_arm6
# ---------------------------------------------------
build-arm7:
	@mkdir -p $(TARGET)
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(MAIN_PROG)_arm7 $(MAIN)
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(LOAD_PROG)_arm7 $(LOAD)
	@cp scripts/linux_arm7/listdbg_arm7.sh .
	@cp scripts/linux_arm7/loaddb_arm7.sh .
	tar czvf $(TARGET)/$(MAIN_PROG)_arm7.tgz $(MAIN_PROG)_arm7 $(LOAD_PROG)_arm7 csv db listdbg_arm7.sh loaddb_arm7.sh
	@rm listdbg_arm7.sh
	@rm loaddb_arm7.sh
	@rm listdbg_arm7
	@rm loaddb_arm7

# ---------------------------------------------------
build-arm64:
	@mkdir -p $(TARGET)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(MAIN_PROG)_arm64 $(MAIN)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(LOAD_PROG)_arm64 $(LOAD)
	@cp scripts/linux_arm64/listdbg_arm64.sh .
	@cp scripts/linux_arm64/loaddb_arm64.sh .
	tar czvf $(TARGET)/$(MAIN_PROG)_arm64.tgz $(MAIN_PROG)_arm64 $(LOAD_PROG)_arm64 csv db listdbg_arm64.sh loaddb_arm64.sh
	@rm listdbg_arm64.sh
	@rm loaddb_arm64.sh
	@rm listdbg_arm64
	@rm loaddb_arm64
# ---------------------------------------------------
clean:
	-@rm $(TARGET)/*
	-@rm $(MAIN_PROG)
	-@rm $(LOAD_PROG)
	-@rm $(MAIN_PROG).sh
	-@rm $(LOAD_PROG).sh

# ---------------------------------------------------
gofmt:
	@gofmt -l -s -w .

