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

install_with_go() {
    print_info "Installing ggc using 'go install'..."
    
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go first or use binary installation."
        return 1
    fi
    
    go install github.com/bmf-san/ggc@latest
    
    local gobin
    gobin=$(go env GOBIN)
    if [ -z "$gobin" ]; then
        gobin="$(go env GOPATH)/bin"
    fi
    
    if [ -f "$gobin/ggc" ]; then
        print_success "ggc installed successfully to $gobin/ggc"
        
        if ! echo "$PATH" | grep -q "$gobin"; then
            print_warning "Go bin directory not in PATH. Adding it automatically..."
            add_to_path "$gobin"
        fi
        
        return 0
    else
        print_error "Installation failed"
        return 1
    fi
}

install_from_source() {
    print_info "Installing ggc from source..."
    
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
    git clone "$REPO_URL" "$temp_dir"
    
    cd "$temp_dir"
    
    print_info "Building ggc..."
    if [ -f "Makefile" ]; then
        make build
    else
        go build -o ggc .
    fi
    
    mkdir -p "$INSTALL_DIR"
    
    mv ggc "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    print_success "ggc installed successfully to $INSTALL_DIR/$BINARY_NAME"
    
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        print_warning "Install directory not in PATH. Adding it automatically..."
        add_to_path "$INSTALL_DIR"
    fi
    
    cd ..
    rm -rf "$temp_dir"
    
    return 0
}

setup_shell_completion() {
    print_info "Setting up shell completion..."
    
    # better way to detect shell using basename
    local current_shell
    current_shell=$(basename "$SHELL" 2>/dev/null || echo "unknown")

	local completion_line
	completion_line=$(cat <<-'EOF'
	# ggc completion loader
	load_ggc_completion() {
		local gopath completion_file
		gopath=$(go env GOPATH 2>/dev/null)
		if [ -z "$gopath" ]; then
			return 1
		fi

		for completion_file in "$gopath"/pkg/mod/github.com/bmf-san/ggc@*/tools/completions/ggc.bash; do
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
                    echo "$completion_line" >> "$bash_profile"
                    print_success "Added ggc completion to $bash_profile"
                    print_info "Restart your terminal or run 'source $bash_profile' to enable completion"
                else
                    print_info "ggc completion already configured in $bash_profile"
                fi
            else
                print_warning "Could not find $bash_profile. Please add the following manually:"
                echo "$completion_line"
            fi
            ;;
        zsh)
            local zsh_profile="$HOME/.zshrc"
            if [ -f "$zsh_profile" ]; then
                if ! grep -q "load_ggc_completion" "$zsh_profile"; then
                    echo "" >> "$zsh_profile"
                    echo "$completion_line" >> "$zsh_profile"
                    print_success "Added ggc completion to $zsh_profile"
                    print_info "Restart your terminal or run 'source $zsh_profile' to enable completion"
                else
                    print_info "ggc completion already configured in $zsh_profile"
                fi
            else
                print_warning "Could not find $zsh_profile. Please add the following manually:"
                echo "$completion_line"
            fi
            ;;
        fish)
            print_info "Fish shell detected. ggc uses bash completion scripts."
            print_info "You may need to install bash completion support for fish."
            ;;
        *)
            print_info "Shell completion setup for $current_shell not automated."
            print_info "To enable shell completion, add the following to your shell profile:"
            echo ""
            echo "For Bash (~/.bashrc or ~/.bash_profile):"
            echo "$completion_line"
            echo ""
            echo "For Zsh (~/.zshrc):"
            echo "$completion_line"
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

main() {
    echo "ggc Installation Script"
    echo "======================="
    
    if command_exists go; then
        INSTALL_METHOD="go"
    else
        INSTALL_METHOD="source"
    fi
    
    print_info "Installation method: $INSTALL_METHOD"
    print_info "Installation directory: $INSTALL_DIR"
    
    installation_success=false   
    if [ "$INSTALL_METHOD" = "go" ] && command_exists go; then
        print_info "Attempting installation with 'go install'..."
        if install_with_go; then
            installation_success=true
        else
            print_warning "Go install failed. Falling back to source installation..."
            if command_exists git; then
                INSTALL_METHOD="source"
                print_info "Attempting installation from source..."
                if install_from_source; then
                    installation_success=true
                fi
            else
                print_error "Git is not installed. Cannot fall back to source installation."
            fi
        fi
    else
        if install_from_source; then
            installation_success=true
        fi
    fi
    
    if [ "$installation_success" = true ]; then
        echo ""
        print_success "✅ Installation completed successfully!"
        print_info "Run 'ggc --help' or just 'ggc' to get started."
        
        echo ""
        setup_shell_completion
    else
        print_error "Installation failed"
        exit 1
    fi
}

main "$@"
