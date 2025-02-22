BIN_DIR := bin
CMDS := gofortune gocowsay gololcat

build: $(CMDS)

$(CMDS):
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$@ ./cmd/$@

clean:
	rm -rf $(BIN_DIR)
