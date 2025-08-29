#!/bin/bash
set -e

# output colors (ANSI)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

REPO_URL="https://github.com/bmf-san/ggc"
INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="ggc"

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

get_latest_version() {
    if command_exists curl; then
        curl -s "https://api.github.com/repos/bmf-san/ggc/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command_exists wget; then
        wget -qO- "https://api.github.com/repos/bmf-san/ggc/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        echo ""
    fi
}

detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux*)   os="linux" ;;
        Darwin*)  os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)        os="unknown" ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        armv7l) arch="arm" ;;
        i386|i686) arch="386" ;;
        *) arch="unknown" ;;
    esac

    echo "${os}_${arch}"
}

install_from_source() {
    print_info "Installing ggc from source (recommended for full version info)..."

    if ! command_exists go; then
        print_error "Go is not installed. Please install Go first."
        return 1
    fi

    if ! command_exists git; then
        print_error "Git is not installed. Please install Git first."
        return 1
    fi

    local temp_dir
    temp_dir=$(mktemp -d)

    print_info "Cloning repository to $temp_dir..."
    if ! git clone "$REPO_URL" "$temp_dir"; then
        print_error "Failed to clone repository"
        rm -rf "$temp_dir"
        return 1
    fi

    cd "$temp_dir"

    print_info "Building ggc with full version information..."

    # Try to build with make first (preferred for version info)
    if [ -f "Makefile" ]; then
        print_info "Using Makefile for build (includes version information)..."
        if make build; then
            print_success "Built successfully with make"
        else
            print_warning "Make build failed, trying go build..."
            if ! go build -o ggc .; then
                print_error "Build failed"
                cd ..
                rm -rf "$temp_dir"
                return 1
            fi
        fi
    else
        print_info "No Makefile found, using go build..."
        if ! go build -o ggc .; then
            print_error "Build failed"
            cd ..
            rm -rf "$temp_dir"
            return 1
        fi
    fi

    # Create install directory
    mkdir -p "$INSTALL_DIR"

    # Move binary to install directory
    if ! mv ggc "$INSTALL_DIR/$BINARY_NAME"; then
        print_error "Failed to move binary to install directory"
        cd ..
        rm -rf "$temp_dir"
        return 1
    fi

    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    print_success "ggc installed successfully to $INSTALL_DIR/$BINARY_NAME"
    print_info "This installation includes full version and commit information"

    # Check if install directory is in PATH
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        print_warning "Install directory not in PATH. Adding it automatically..."
        add_to_path "$INSTALL_DIR"
    fi

    cd ..
    rm -rf "$temp_dir"

    return 0
}

install_with_go() {
    print_warning "Installing ggc using 'go install' (fallback method)..."
    print_warning "Note: This method provides limited version information compared to source installation"
    print_warning "Version may show as development build."

    if ! command_exists go; then
        print_error "Go is not installed. Please install Go first."
        return 1
    fi

    print_info "Running go install..."
    if ! go install github.com/bmf-san/ggc/v4@latest; then
        print_error "go install failed"
        return 1
    fi

    local gobin
    gobin=$(go env GOBIN)
    if [ -z "$gobin" ]; then
        gobin="$(go env GOPATH)/bin"
    fi

    if [ -f "$gobin/ggc" ]; then
        print_success "ggc installed successfully to $gobin/ggc"
        print_warning "This installation has limited version information due to go install limitations"

        if ! echo "$PATH" | grep -q "$gobin"; then
            print_warning "Go bin directory not in PATH. Adding it automatically..."
            add_to_path "$gobin"
        fi

        return 0
    else
        print_error "Installation failed - binary not found"
        return 1
    fi
}

create_fish_completion() {
    local fish_completions_dir="$HOME/.config/fish/completions"
    local fish_completion_file="$fish_completions_dir/ggc.fish"

    # Create completions directory if it doesn't exist
    mkdir -p "$fish_completions_dir"

    # Create the fish completion loader
    cat > "$fish_completion_file" << 'EOF'
# Fish completion loader for ggc

function __ggc_load_completion
    # Try to find the ggc.fish completion file in go modules
    set -l gopath (go env GOPATH 2>/dev/null)
    if test -z "$gopath"
        return 1
    end

    # Look for completion file in go modules
    for completion_file in "$gopath"/pkg/mod/github.com/bmf-san/ggc/v4@*/tools/completions/ggc.fish
        if test -f "$completion_file"
            source "$completion_file"
            return 0
        end
    end

    # Fallback: find completion file
    set -l completion_file (find "$gopath/pkg/mod/github.com/bmf-san" -name "ggc.fish" -path "*/tools/completions/*" 2>/dev/null | head -1)
    if test -n "$completion_file"; and test -f "$completion_file"
        source "$completion_file"
        return 0
    end

    return 1
end

# Load completion if go is available
if command -v go >/dev/null 2>&1
    __ggc_load_completion
end
EOF

    print_success "Created fish completion loader at $fish_completion_file"
}

setup_shell_completion() {
    print_info "Setting up shell completion..."

    # better way to detect shell using basename
    local current_shell
    current_shell=$(basename "$SHELL" 2>/dev/null || echo "unknown")

    local bash_completion_line
    bash_completion_line=$(cat <<-'EOF'
# ggc completion loader
load_ggc_completion() {
	local gopath completion_file
	gopath=$(go env GOPATH 2>/dev/null)
	if [ -z "$gopath" ]; then
		return 1
	fi

	for completion_file in "$gopath"/pkg/mod/github.com/bmf-san/ggc/v4@*/tools/completions/ggc.bash; do
		if [ -f "$completion_file" ]; then
			source "$completion_file"
			return 0
		fi
	done

	completion_file=$(find "$gopath/pkg/mod/github.com/bmf-san" -name "ggc.bash" -path "*/tools/completions/*" 2>/dev/null | head -1)
	if [ -n "$completion_file" ] && [ -f "$completion_file" ]; then
		source "$completion_file"
		return 0
	fi

	return 1
}

# Load completion if go is available
if command -v go >/dev/null 2>&1; then
	load_ggc_completion
fi
EOF
)

    local zsh_completion_line
    zsh_completion_line=$(cat <<-'EOF'
# ggc completion loader for zsh
load_ggc_completion() {
	local gopath completion_file
	gopath=$(go env GOPATH 2>/dev/null)
	if [ -z "$gopath" ]; then
		return 1
	fi

	for completion_file in "$gopath"/pkg/mod/github.com/bmf-san/ggc/v4@*/tools/completions/ggc.zsh; do
		if [ -f "$completion_file" ]; then
			source "$completion_file"
			return 0
		fi
	done

	completion_file=$(find "$gopath/pkg/mod/github.com/bmf-san" -name "ggc.zsh" -path "*/tools/completions/*" 2>/dev/null | head -1)
	if [ -n "$completion_file" ] && [ -f "$completion_file" ]; then
		source "$completion_file"
		return 0
	fi

	return 1
}

# Load completion if go is available
if command -v go >/dev/null 2>&1; then
	load_ggc_completion
fi
EOF
)

    case "$current_shell" in
        bash)
            local bash_profile="$HOME/.bashrc"
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS typically uses ~/.bash_profile
                bash_profile="$HOME/.bash_profile"
                if [ ! -f "$bash_profile" ] && [ -f "$HOME/.bashrc" ]; then
                    bash_profile="$HOME/.bashrc"
                fi
            fi

            if [ -f "$bash_profile" ]; then
                if ! grep -q "load_ggc_completion" "$bash_profile"; then
                    echo "" >> "$bash_profile"
                    echo "$bash_completion_line" >> "$bash_profile"
                    print_success "Added ggc completion to $bash_profile"
                    print_info "Restart your terminal or run 'source $bash_profile' to enable completion"
                else
                    print_info "ggc completion already configured in $bash_profile"
                fi
            else
                print_warning "Could not find $bash_profile. Please add the following manually:"
                echo "$bash_completion_line"
            fi
            ;;
        zsh)
            local zsh_profile="$HOME/.zshrc"
            if [ -f "$zsh_profile" ]; then
                if ! grep -q "load_ggc_completion" "$zsh_profile"; then
                    echo "" >> "$zsh_profile"
                    echo "$zsh_completion_line" >> "$zsh_profile"
                    print_success "Added ggc completion to $zsh_profile"
                    print_info "Restart your terminal or run 'source $zsh_profile' to enable completion"
                else
                    print_info "ggc completion already configured in $zsh_profile"
                fi
            else
                print_warning "Could not find $zsh_profile. Please add the following manually:"
                echo "$zsh_completion_line"
            fi
            ;;
        fish)
            print_info "Fish shell detected. Creating native fish completion..."
            create_fish_completion
            print_info "Fish completion installed. Restart your terminal or run 'fish' to enable completion"
            ;;
        *)
            print_info "Shell completion setup for $current_shell not automated."
            print_info "To enable shell completion, run this ggc command and add the output to your shell profile:"
            echo "ggc completion $current_shell"
            ;;
    esac
}

add_to_path() {
    local path_to_add="$1"

    local current_shell
    current_shell=$(basename "$SHELL" 2>/dev/null || echo "unknown")

    local path_line="export PATH=\$PATH:$path_to_add"

    case "$current_shell" in
        bash)
            local bash_profile="$HOME/.bashrc"
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS typically uses ~/.bash_profile
                bash_profile="$HOME/.bash_profile"
                if [ ! -f "$bash_profile" ] && [ -f "$HOME/.bashrc" ]; then
                    bash_profile="$HOME/.bashrc"
                fi
            fi

            if [ -f "$bash_profile" ]; then
                if ! grep -q "PATH.*$path_to_add" "$bash_profile"; then
					{
						echo ""
						echo "# ggc PATH"
						echo "$path_line"
					} >> "$bash_profile"
                    print_success "Added $path_to_add to PATH in $bash_profile"
                    print_info "Restart your terminal or run 'source $bash_profile' to use ggc"
                else
                    print_info "PATH already contains $path_to_add in $bash_profile"
                fi
            else
                print_warning "Could not find $bash_profile. Please add the following manually:"
                echo "  $path_line"
            fi
            ;;
        zsh)
            local zsh_profile="$HOME/.zshrc"
            if [ -f "$zsh_profile" ]; then
                if ! grep -q "PATH.*$path_to_add" "$zsh_profile"; then
					{
						echo ""
						echo "# ggc PATH"
						echo "$path_line"
					} >> "$zsh_profile"
                    print_success "Added $path_to_add to PATH in $zsh_profile"
                    print_info "Restart your terminal or run 'source $zsh_profile' to use ggc"
                else
                    print_info "PATH already contains $path_to_add in $zsh_profile"
                fi
            else
                print_warning "Could not find $zsh_profile. Please add the following manually:"
                echo "  $path_line"
            fi
            ;;
        fish)
            local fish_config="$HOME/.config/fish/config.fish"
            local fish_path_line="set -gx PATH \$PATH $path_to_add"

            if [ -f "$fish_config" ]; then
                if ! grep -q "PATH.*$path_to_add" "$fish_config"; then
					{
						echo ""
						echo "# ggc PATH"
						echo "$fish_path_line"
					} >> "$fish_config"
                    print_success "Added $path_to_add to PATH in $fish_config"
                    print_info "Restart your terminal to use ggc"
                else
                    print_info "PATH already contains $path_to_add in $fish_config"
                fi
            else
                print_warning "Could not find $fish_config. Please add the following manually:"
                echo "  $fish_path_line"
            fi
            ;;
        *)
            print_warning "Automatic PATH setup for $current_shell not supported."
            print_info "Please add the following to your shell profile:"
            echo "  $path_line"
            ;;
    esac
}

check_dependencies() {
    local missing_deps=()

    if ! command_exists go; then
        missing_deps+=("go")
    fi

    if ! command_exists git; then
        missing_deps+=("git")
    fi

    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        print_info "For source installation (recommended), you need both Go and Git"
        print_info "For fallback installation, you only need Go (but with limited version info)"
        return 1
    fi

    return 0
}

main() {
    echo "ggc Installation Script"
    echo "======================="

    print_info "Installation directory: $INSTALL_DIR"

    installation_success=false

    # Check if we have the basic requirements
    if command_exists go && command_exists git; then
        print_info "✅ Go and Git detected - proceeding with source installation"
        print_info "This will provide full version and commit information"

        if install_from_source; then
            installation_success=true
        else
            print_warning "❌ Source installation failed"
            print_info "Falling back to 'go install' method..."

            if install_with_go; then
                installation_success=true
            fi
        fi
    elif command_exists go; then
        print_warning "⚠️  Git not found - source installation unavailable"
        print_info "Falling back to 'go install' method (limited version info)"

        if install_with_go; then
            installation_success=true
        fi
    else
        print_error "❌ Go is not installed"
        print_error "Please install Go first: https://golang.org/doc/install"
        print_info "For best results, also install Git for source installation"
        exit 1
    fi

    if [ "$installation_success" = true ]; then
        echo ""
        print_success "✅ Installation completed successfully!"

        # Test the installation
        local installed_binary
        if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
            installed_binary="$INSTALL_DIR/$BINARY_NAME"
        else
            local gobin
            gobin=$(go env GOBIN 2>/dev/null)
            if [ -z "$gobin" ]; then
                gobin="$(go env GOPATH)/bin"
            fi
            if [ -f "$gobin/ggc" ]; then
                installed_binary="$gobin/ggc"
            fi
        fi

        if [ -n "$installed_binary" ]; then
            print_info "Testing installation..."
            if "$installed_binary" version >/dev/null 2>&1; then
                print_success "✅ ggc is working correctly"
            else
                print_warning "⚠️  ggc installed but version command failed"
            fi
        fi

        print_info "Run 'ggc --help' or just 'ggc' to get started"

        echo ""
        setup_shell_completion
    else
        print_error "❌ Installation failed"
        print_info "Please check the error messages above and try again"
        exit 1
    fi
}

main "$@"
