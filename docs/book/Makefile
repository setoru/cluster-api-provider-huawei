ROOT_DIR_RELATIVE := ../..
include $(ROOT_DIR_RELATIVE)/common.mk

MDBOOK := $(TOOLS_BIN_DIR)/mdbook
MDBOOK_EMBED := $(TOOLS_BIN_DIR)/mdbook-embed
MDBOOK_RELEASELINK := $(TOOLS_BIN_DIR)/mdbook-releaselink
MDBOOK_TABULATE := $(TOOLS_BIN_DIR)/mdbook-tabulate

PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
export PATH

BOOK_DEPS := $(MDBOOK) $(MDBOOK_EMBED) $(MDBOOK_RELEASELINK) $(MDBOOK_TABULATE)

.PHONY: build
build: $(BOOK_DEPS) ## Build the book
	$(MDBOOK) build

.PHONY: serve
serve: $(BOOK_DEPS) ## Start a local server to preview the book
	$(MDBOOK) serve -p 3000 -n 127.0.0.1

.PHONY: clean
clean: $(BOOK_DEPS) ## Clean the book
	$(MDBOOK) clean
