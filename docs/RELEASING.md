# Releasing PDF Joiner

Official releases are built, Developer ID-signed, and notarized locally on the maintainer's Mac. The certificate's private key remains in the login keychain and is never exported to CI.

## One-time setup

Install these certificates, including their private keys, in the login keychain:

- `Developer ID Application: Vinit Kumar (RRVJ5HE22Y)` signs the macOS executables.
- `Developer ID Installer: Vinit Kumar (RRVJ5HE22Y)` is available if a `.pkg` installer is added later; the current archive release does not use it.

Verify the signing identity:

```bash
security find-identity -v -p codesigning
```

Store Apple notarization credentials in Keychain. The Apple ID form securely prompts for an app-specific password:

```bash
xcrun notarytool store-credentials "pdf-joiner-notary" \
  --apple-id "APPLE_ACCOUNT_EMAIL" \
  --team-id RRVJ5HE22Y
```

Verify the saved profile:

```bash
xcrun notarytool history --keychain-profile "pdf-joiner-notary"
```

## Build, sign, and notarize

Start from a clean commit and provide the release version:

```bash
make check
VERSION=1.3.0 make release-local
```

The local release script:

1. Builds macOS and Linux binaries for amd64 and arm64.
2. Signs both macOS binaries with hardened runtime and a secure timestamp.
3. Submits both signed executables to Apple in one notarization request.
4. Waits for an `Accepted` result.
5. Creates four release archives, the installer, and SHA-256 checksums.
6. Verifies the macOS signatures and Gatekeeper assessment.

Artifacts are written to `dist/release/VERSION`. Existing output is never overwritten.

The defaults can be overridden when necessary:

```bash
SIGNING_IDENTITY="Developer ID Application: Example (TEAMID)" \
NOTARY_PROFILE="another-notary-profile" \
OUTPUT_DIR="/absolute/output/path" \
VERSION="1.3.0" \
make release-local
```

## Inspect the release

```bash
cd dist/release/1.3.0
shasum -a 256 -c checksums.txt
codesign -dv --verbose=4 bin/darwin/arm64/pdf-joiner
spctl --assess --type execute --verbose=4 bin/darwin/arm64/pdf-joiner
```

Every archive contains an executable named simply `pdf-joiner`; platform and architecture appear only in the archive filename. The signature should identify `Developer ID Application: Vinit Kumar (RRVJ5HE22Y)` and include hardened runtime and a secure timestamp. Apple must report `Accepted`, `Ready for distribution`, and no issues in `notarization.json` and `notarization-log.json`.

`spctl` can report that a valid standalone command-line executable “does not seem to be an app.” That message describes the assessment context; it does not override an `Accepted` notary-service result. A future `.pkg` installer can provide the conventional stapled Gatekeeper experience.

Apple publishes notarization tickets online for standalone executables. Tickets cannot be stapled directly to standalone executables or ZIP archives. For offline Gatekeeper validation, add a signed, notarized, and stapled `.pkg` or `.dmg` distribution.

## Publish to GitHub

After inspecting the artifacts, create the tag and release:

```bash
git tag -a v1.3.0 -m "v1.3.0"
git push origin v1.3.0
gh release create v1.3.0 dist/release/1.3.0/*.tar.gz \
  dist/release/1.3.0/checksums.txt dist/release/1.3.0/install.sh \
  --generate-notes \
  --title "v1.3.0"
```

Never publish if the local release script fails or the notarization result is anything other than `Accepted`.
