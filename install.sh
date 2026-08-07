#!/bin/sh

set -eu

repository="vinitkumar/pdf-joiner"
binary_name="pdf-joiner"
install_dir="${PDF_JOINER_INSTALL_DIR:-${HOME:?HOME must be set}/.local/bin}"

fail() {
	printf 'pdf-joiner installer: %s\n' "$*" >&2
	exit 1
}

for command_name in curl tar awk sed uname; do
	command -v "$command_name" >/dev/null 2>&1 || fail "required command not found: $command_name"
done

case "$(uname -s)" in
	Darwin) operating_system="darwin" ;;
	Linux) operating_system="linux" ;;
	*) fail "unsupported operating system: $(uname -s)" ;;
esac

case "$(uname -m)" in
	x86_64 | amd64) architecture="amd64" ;;
	arm64 | aarch64) architecture="arm64" ;;
	*) fail "unsupported architecture: $(uname -m)" ;;
esac

temporary_dir="$(mktemp -d "${TMPDIR:-/tmp}/pdf-joiner-install.XXXXXX")"
cleanup() {
	rm -rf "$temporary_dir"
}
trap cleanup EXIT HUP INT TERM

requested_version="${PDF_JOINER_VERSION:-}"
if [ -n "$requested_version" ]; then
	release_tag="v${requested_version#v}"
else
	release_metadata="$temporary_dir/release.json"
	curl -fsSL --retry 3 --proto '=https' --tlsv1.2 \
		-o "$release_metadata" \
		"https://api.github.com/repos/${repository}/releases/latest"
	release_tag="$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$release_metadata" | awk 'NR == 1 { print; exit }')"
	[ -n "$release_tag" ] || fail "could not determine the latest release"
fi

version="${release_tag#v}"
archive_name="${binary_name}_${version}_${operating_system}_${architecture}.tar.gz"
release_url="https://github.com/${repository}/releases/download/${release_tag}"
archive_path="$temporary_dir/$archive_name"
checksums_path="$temporary_dir/checksums.txt"

printf 'Downloading pdf-joiner %s for %s/%s...\n' "$version" "$operating_system" "$architecture"
curl -fsSL --retry 3 --proto '=https' --tlsv1.2 -o "$archive_path" "$release_url/$archive_name"
curl -fsSL --retry 3 --proto '=https' --tlsv1.2 -o "$checksums_path" "$release_url/checksums.txt"

expected_checksum="$(awk -v archive="$archive_name" '$2 == archive || $2 == "./" archive { print $1; exit }' "$checksums_path")"
[ -n "$expected_checksum" ] || fail "no checksum found for $archive_name"

if command -v shasum >/dev/null 2>&1; then
	actual_checksum="$(shasum -a 256 "$archive_path" | awk '{ print $1 }')"
elif command -v sha256sum >/dev/null 2>&1; then
	actual_checksum="$(sha256sum "$archive_path" | awk '{ print $1 }')"
else
	fail "shasum or sha256sum is required to verify the download"
fi

[ "$actual_checksum" = "$expected_checksum" ] || fail "checksum verification failed for $archive_name"

extract_dir="$temporary_dir/extracted"
mkdir -p "$extract_dir"
tar -xzf "$archive_path" -C "$extract_dir"
[ -f "$extract_dir/$binary_name" ] || fail "release archive does not contain $binary_name"

if [ "$operating_system" = "darwin" ] && command -v codesign >/dev/null 2>&1; then
	codesign --verify --strict "$extract_dir/$binary_name" || fail "the macOS signature is invalid"
fi

mkdir -p "$install_dir"
temporary_binary="$install_dir/.${binary_name}.tmp.$$"
cp "$extract_dir/$binary_name" "$temporary_binary"
chmod 0755 "$temporary_binary"
mv "$temporary_binary" "$install_dir/$binary_name"

path_configured=false
case ":${PATH:-}:" in
	*":$install_dir:"*) path_configured=true ;;
esac

if [ "$path_configured" = false ] && [ "${PDF_JOINER_NO_MODIFY_PATH:-0}" != "1" ]; then
	shell_name="${SHELL##*/}"
	marker="# Added by the pdf-joiner installer: $install_dir"
	escaped_install_dir="$(printf '%s' "$install_dir" | sed "s/'/'\\\\''/g")"

	case "$shell_name" in
		zsh)
			shell_config="${ZDOTDIR:-$HOME}/.zshrc"
			path_line="export PATH='$escaped_install_dir':\"\$PATH\""
			;;
		bash)
			if [ "$operating_system" = "darwin" ]; then
				shell_config="$HOME/.bash_profile"
			else
				shell_config="$HOME/.bashrc"
			fi
			path_line="export PATH='$escaped_install_dir':\"\$PATH\""
			;;
		fish)
			shell_config="$HOME/.config/fish/config.fish"
			path_line="fish_add_path -- '$escaped_install_dir'"
			;;
		*)
			shell_config="$HOME/.profile"
			path_line="export PATH='$escaped_install_dir':\"\$PATH\""
			;;
	esac

	mkdir -p "$(dirname "$shell_config")"
	if [ ! -f "$shell_config" ] || ! grep -F -- "$marker" "$shell_config" >/dev/null 2>&1; then
		printf '\n%s\n%s\n' "$marker" "$path_line" >>"$shell_config"
		printf 'Added %s to PATH in %s.\n' "$install_dir" "$shell_config"
	fi
fi

if [ "$operating_system" = "linux" ] && \
	! command -v pdfunite >/dev/null 2>&1 && \
	! command -v gs >/dev/null 2>&1 && \
	! command -v qpdf >/dev/null 2>&1; then
	printf '%s\n' 'Warning: install pdfunite (recommended), Ghostscript, or qpdf before joining PDFs.' >&2
fi

installed_version="$("$install_dir/$binary_name" -version)"
printf 'Installed pdf-joiner %s at %s.\n' "$installed_version" "$install_dir/$binary_name"
if [ "$path_configured" = false ]; then
	printf 'Open a new terminal, or run: export PATH=%s:%s\n' "$install_dir" "\"\$PATH\""
fi
