# Releasing PDF Joiner

Every tag matching `v*` starts the release workflow. It tests the repository, builds Linux and macOS archives, signs each macOS executable with Developer ID, submits it to Apple's notary service, waits for acceptance, and then publishes the GitHub release.

The workflow deliberately fails before building a release if any Apple credential is absent.

## One-time Apple setup

An active Apple Developer Program membership is required.

### 1. Create the signing certificate

In Apple Developer, open **Certificates, Identifiers & Profiles → Certificates**, create a **Developer ID Application** certificate, and install it in Keychain Access. The certificate must appear under **My Certificates** with its private key beneath it.

Do not use Apple Development, Apple Distribution, or Mac App Distribution certificates.

Confirm the identity locally:

```bash
security find-identity -v -p codesigning
```

Export the certificate and private key from Keychain Access as a password-protected `.p12` file. Keep the file and its password private.

### 2. Create the notarization key

In App Store Connect, open **Users and Access → Integrations → App Store Connect API** and create a team API key. Record its **Issuer ID** and **Key ID**, then download the `.p8` private key. Apple allows the private key to be downloaded only once.

### 3. Configure GitHub Actions secrets

Base64-encode the two credential files on macOS:

```bash
base64 -i DeveloperIDApplication.p12 | pbcopy
base64 -i AuthKey_KEYID.p8 | pbcopy
```

In **GitHub repository → Settings → Secrets and variables → Actions**, add:

| Secret | Value |
| --- | --- |
| `MACOS_SIGN_P12` | Base64 contents of the `.p12` file |
| `MACOS_SIGN_PASSWORD` | Password used when exporting the `.p12` |
| `MACOS_NOTARY_KEY` | Base64 contents of the `.p8` file |
| `MACOS_NOTARY_KEY_ID` | App Store Connect API Key ID |
| `MACOS_NOTARY_ISSUER_ID` | App Store Connect API Issuer ID |

Base64 is an encoding, not encryption. Never commit these values, paste them into workflow files, or expose them in logs.

## Create a release

First verify locally:

```bash
make check
make release-snapshot
```

Then create and push a semantic-version tag:

```bash
git tag -a v1.0.0 -m "v1.0.0"
git push origin v1.0.0
```

The release appears only after Apple accepts both Darwin binaries.

## Verify a downloaded macOS release

Extract the Darwin archive and run:

```bash
codesign --verify --strict --verbose=2 ./pdf-joiner
codesign -dv --verbose=4 ./pdf-joiner
spctl --assess --type execute --verbose=4 ./pdf-joiner
```

The signature details should contain `Developer ID Application`, the expected Team Identifier, the hardened-runtime flag, and a secure timestamp. `spctl` should report that the executable is accepted and notarized.

Apple publishes the notarization ticket online for a standalone executable. Tickets cannot be stapled directly to a standalone executable or ZIP archive. If offline Gatekeeper validation becomes a requirement, add a signed, notarized, and stapled `.pkg` or `.dmg` distribution.

## Credential rotation

If a certificate or API key is lost or exposed, revoke it in Apple Developer or App Store Connect immediately, create a replacement, and update the corresponding GitHub secrets. Existing releases should never be re-signed in place; publish a new version instead.
