#!/usr/bin/env bats

setup() {
    export TEST_DIR="$(mktemp -d)"
    export ORIGINAL_HOME="$HOME"
    export ORIGINAL_PATH="$PATH"
    export HOME="$TEST_DIR"

    source install.sh

    export PATH="$TEST_DIR/bin:$PATH"
    mkdir -p "$TEST_DIR/bin"
}

teardown() {
    rm -rf "$TEST_DIR"
    export HOME="$ORIGINAL_HOME"
    export PATH="$ORIGINAL_PATH"
}

@test "command_exists returns true for existing command" {
    echo '#!/bin/bash' > "$TEST_DIR/bin/test_cmd"
    chmod +x "$TEST_DIR/bin/test_cmd"

    run command_exists test_cmd
    [ "$status" -eq 0 ]
}

@test "command_exists returns false for non-existing command" {
    run command_exists nonexistent_command
    [ "$status" -eq 1 ]
}

@test "detect_platform returns correct format" {
    run detect_platform
    [ "$status" -eq 0 ]
    [[ "$output" =~ ^[a-z]+_[a-z0-9]+$ ]]
}

@test "print_info outputs colored message" {
    run print_info "test message"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test message" ]]
}

@test "print_success outputs colored message" {
    run print_success "success message"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "success message" ]]
}

@test "print_error outputs colored message" {
    run print_error "error message"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "error message" ]]
}

@test "add_to_path adds to bashrc when using bash" {
    export SHELL="/bin/bash"
    touch "$HOME/.bashrc"

    run add_to_path "/test/path"
    [ "$status" -eq 0 ]

    grep -q "/test/path" "$HOME/.bashrc"
}

@test "add_to_path adds to zshrc when using zsh" {
    export SHELL="/bin/zsh"
    touch "$HOME/.zshrc"

    run add_to_path "/test/path"
    [ "$status" -eq 0 ]

    grep -q "/test/path" "$HOME/.zshrc"
}

@test "setup_shell_completion adds completion to bashrc" {
    export SHELL="/bin/bash"
    echo '#!/bin/bash
echo "/fake/gopath"' > "$TEST_DIR/bin/go"
    chmod +x "$TEST_DIR/bin/go"
    touch "$HOME/.bashrc"

    run setup_shell_completion
    [ "$status" -eq 0 ]

    grep -q "load_ggc_completion" "$HOME/.bashrc"
}

@test "install_with_go fails when go is not available" {
    mkdir -p "$TEST_DIR/empty"
    export PATH="$TEST_DIR/empty"

    run install_with_go
    [ "$status" -eq 1 ]
}

@test "install_from_source fails when go is not available" {
    mkdir -p "$TEST_DIR/empty"
    export PATH="$TEST_DIR/empty"

    run install_from_source
    [ "$status" -eq 1 ]
}

@test "install_from_source fails when git is not available" {
    mkdir -p "$TEST_DIR/bin"
    echo '#!/bin/bash
echo "go version go1.25.0 linux/amd64"' > "$TEST_DIR/bin/go"
    chmod +x "$TEST_DIR/bin/go"

    mkdir -p "$TEST_DIR/empty"
    export PATH="$TEST_DIR/bin:$TEST_DIR/empty"

    run install_from_source
    [ "$status" -eq 1 ]
}
