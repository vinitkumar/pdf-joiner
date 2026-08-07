#!/usr/bin/env bash

set -euo pipefail

repository_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
binary_name="pdf-joiner"
signing_identity="${SIGNING_IDENTITY:-Developer ID Application: Vinit Kumar (RRVJ5HE22Y)}"
notary_profile="${NOTARY_PROFILE:-pdf-joiner-notary}"

cd "$repository_root"

if [[ "$(uname -s)" != "Darwin" ]]; then
	echo "Local signed releases must run on macOS." >&2
	exit 1
fi

for command in go codesign security xcrun ditto tar shasum; do
	if ! command -v "$command" >/dev/null 2>&1; then
		echo "Required command not found: $command" >&2
		exit 1
	fi
done

if ! security find-identity -v -p codesigning | grep -F "\"${signing_identity}\"" >/dev/null; then
	echo "Signing identity is unavailable: ${signing_identity}" >&2
	exit 1
fi

echo "Validating notarization profile ${notary_profile}..."
xcrun notarytool history --keychain-profile "$notary_profile" >/dev/null

version="${VERSION:-$(git -C "$repository_root" describe --tags --always --dirty)}"
version="${version#v}"
artifact_version="${version//\//-}"
output_dir="${OUTPUT_DIR:-$repository_root/dist/release/$artifact_version}"

if [[ -e "$output_dir" ]]; then
	echo "Output directory already exists: ${output_dir}" >&2
	echo "Choose another VERSION or OUTPUT_DIR to avoid overwriting a release." >&2
	exit 1
fi

work_dir="$(mktemp -d "${TMPDIR:-/tmp}/pdf-joiner-release.XXXXXX")"
cleanup() {
	rm -rf "$work_dir"
}
trap cleanup EXIT

mkdir -p "$work_dir/build" "$work_dir/notary-payload" "$output_dir/bin"

targets=(
	"darwin amd64"
	"darwin arm64"
	"linux amd64"
	"linux arm64"
)

echo "Building PDF Joiner ${version}..."
for target in "${targets[@]}"; do
	read -r goos goarch <<<"$target"
	build_dir="$work_dir/build/${goos}-${goarch}"
	binary_path="$build_dir/$binary_name"
	mkdir -p "$build_dir"

	GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "-s -w -X main.version=${version}" \
		-o "$binary_path" \
		.

	if [[ "$goos" == "darwin" ]]; then
		echo "Signing darwin/${goarch}..."
		codesign --force --options runtime --timestamp --sign "$signing_identity" "$binary_path"
		codesign --verify --strict --verbose=2 "$binary_path"
		mkdir -p "$work_dir/notary-payload/$goarch"
		cp "$binary_path" "$work_dir/notary-payload/$goarch/$binary_name"
	fi

	mkdir -p "$output_dir/bin/$goos/$goarch"
	cp "$binary_path" "$output_dir/bin/$goos/$goarch/$binary_name"
done

submission_zip="$work_dir/pdf-joiner-notarization.zip"
ditto -c -k --keepParent "$work_dir/notary-payload" "$submission_zip"

echo "Submitting signed macOS binaries to Apple..."
notary_response="$output_dir/notarization.json"
xcrun notarytool submit "$submission_zip" \
	--keychain-profile "$notary_profile" \
	--wait \
	--output-format json >"$notary_response"

notary_status="$(plutil -extract status raw -o - "$notary_response")"
if [[ "$notary_status" != "Accepted" ]]; then
	echo "Apple notarization did not succeed; status: ${notary_status}" >&2
	cat "$notary_response" >&2
	exit 1
fi

echo "Apple accepted the notarization submission. Creating archives..."
for target in "${targets[@]}"; do
	read -r goos goarch <<<"$target"
	package_dir="$work_dir/package-${goos}-${goarch}"
	mkdir -p "$package_dir"
	cp "$output_dir/bin/$goos/$goarch/$binary_name" "$package_dir/$binary_name"
	cp "$repository_root/README.md" "$repository_root/LICENSE" "$package_dir/"
	COPYFILE_DISABLE=1 tar -czf "$output_dir/${binary_name}_${artifact_version}_${goos}_${goarch}.tar.gz" \
		-C "$package_dir" "$binary_name" README.md LICENSE
done

(
	cd "$output_dir"
	shasum -a 256 ./*.tar.gz >checksums.txt
	shasum -a 256 -c checksums.txt
)

for arch in amd64 arm64; do
	binary_path="$output_dir/bin/darwin/${arch}/$binary_name"
	codesign --verify --strict --verbose=2 "$binary_path"
	spctl --assess --type execute --verbose=4 "$binary_path"
done

echo "Release artifacts are ready in ${output_dir}"
