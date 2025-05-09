# Kontext

A simple CLI tool for managing Kubernetes contexts with ease.

## Overview

Kontext helps you manage your Kubernetes contexts efficiently with a simple command-line interface. It allows you to:

- List all available Kubernetes contexts
- Show your current active context
- Switch between contexts with tab completion
- Interactively select contexts from a menu
- View and change namespaces within contexts

## Installation

```bash
go install github.com/user-cube/kontext@latest
```

## Usage

### List all available contexts

```bash
kontext list
```

### Show current context

```bash
kontext current
```

### Switch to a different context

```bash
kontext switch <context-name>
```
With tab completion for context names!

### Interactive context selection

Simply run:
```bash
kontext switch
```
Without any arguments to get an interactive selection menu of all available contexts.

You can also run just:
```bash
kontext
```
To access the same interactive context selection.

### Namespace Management

View or change the current namespace:
```bash
kontext namespace
```
or use the shorter alias:
```bash
kontext ns
```

Show the current namespace without the selector:
```bash
kontext ns -s
```

Switch to a specific namespace without using the interactive selector:
```bash
kontext ns my-namespace
```

### Switch Context and Namespace Together

Switch context and then set namespace in one command:
```bash
kontext switch -n
kontext switch my-context -n
```

## Shell Completion

To enable shell completion:

### Bash

```bash
echo 'source <(kontext completion bash)' >> ~/.bashrc
```

### Zsh

```bash
echo 'source <(kontext completion zsh)' >> ~/.zshrc
```

### Fish

```bash
kontext completion fish > ~/.config/fish/completions/kontext.fish
```