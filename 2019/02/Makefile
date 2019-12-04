CC=clang
CFLAGS=
OBJDIR=obj
BINDIR=bin
SRCDIR=.
TARGET=main
SOURCES := $(wildcard $(SRCDIR)/*.c)
OBJECTS := $(SOURCES:$(SRCDIR)/%.c=$(OBJDIR)/%.o)

.PHONY: all build
all: build

build: $(OBJDIR) $(BINDIR) $(BINDIR)/$(TARGET)

$(OBJECTS) : $(OBJDIR)/%.o : $(SRCDIR)/%.c
	$(CC) $(CFLAGS) -c $< -o $@

$(BINDIR)/$(TARGET): $(OBJECTS)
	$(CC) $(CFLAGS) $^ -o $@

$(OBJDIR) $(BINDIR):
	mkdir -p $@
