# 🐧 Feature Request: Add Linux Support for PDF Joining

## Overview

This project currently provides a simple and efficient PDF joining tool for **macOS only**, leveraging the native macOS PDF utility at `/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join`. We're looking for community contributions to extend this functionality to **Linux systems** using stable, native Linux PDF processing tools.

## Current State

The tool currently:
- ✅ Works seamlessly on macOS using built-in system utilities
- ✅ Joins multiple PDF files with simple command-line interface
- ✅ Supports custom output paths and automatic timestamped naming
- ✅ Validates input files and creates output directories
- ❌ **Linux support is missing**

## What We Need

We're seeking contributions to add Linux support while maintaining the same design principles:

### 🎯 Goals
1. **Use native/stable Linux tools** - Just like we use macOS's built-in utility, we want to leverage well-established Linux PDF tools
2. **Maintain the same CLI interface** - The user experience should be identical across platforms
3. **Cross-platform detection** - Automatically detect the OS and use the appropriate backend
4. **No external dependencies** - Avoid requiring users to install additional libraries

### 🛠️ Technical Implementation Ideas

We need platform detection and multiple backend support. Here are some proven Linux PDF joining tools to consider:

#### Option 1: `pdfunite` (from poppler-utils)
```bash
# Most Linux distributions include this by default
pdfunite input1.pdf input2.pdf input3.pdf output.pdf
```

#### Option 2: `ghostscript` (gs)
```bash
# Very stable, widely available
gs -dBATCH -dNOPAUSE -q -sDEVICE=pdfwrite -sOutputFile=output.pdf input1.pdf input2.pdf
```

#### Option 3: `pdftk` 
```bash
# Popular choice, though may need installation
pdftk input1.pdf input2.pdf cat output output.pdf
```

#### Option 4: `qpdf`
```bash
# Another stable option
qpdf --empty --pages input1.pdf input2.pdf -- output.pdf
```

### 📋 Implementation Checklist

The ideal contribution should include:

- [ ] **OS Detection Logic** - Use `runtime.GOOS` to detect the operating system
- [ ] **Backend Selection** - Implement a strategy pattern for different PDF joining backends
- [ ] **Linux Tool Detection** - Check which PDF tools are available on the system (priority order)
- [ ] **Error Handling** - Graceful fallbacks and clear error messages if no suitable tool is found
- [ ] **Tests** - Unit tests for the new functionality (mocked where appropriate)
- [ ] **Documentation** - Update README with Linux requirements and installation instructions
- [ ] **CI/CD** - Ensure GitHub Actions test on Linux environments

### 🏗️ Suggested Code Structure

```go
type PDFJoiner interface {
    Join(inputFiles []string, outputPath string) error
    IsAvailable() bool
}

type MacOSJoiner struct{}
type LinuxJoiner struct {
    backend string // "pdfunite", "gs", "pdftk", etc.
}

func NewPDFJoiner() PDFJoiner {
    switch runtime.GOOS {
    case "darwin":
        return &MacOSJoiner{}
    case "linux":
        return NewLinuxJoiner()
    default:
        return nil // unsupported OS
    }
}
```

### 🐧 Linux Distribution Considerations

The solution should work across major Linux distributions:
- **Ubuntu/Debian** - `poppler-utils` (pdfunite) is usually available
- **RHEL/CentOS/Fedora** - Similar package availability
- **Arch Linux** - Most PDF tools available via pacman
- **Alpine Linux** - Lightweight options preferred

### 📚 Resources for Contributors

- **Current macOS implementation**: [`main.go:53-77`](main.go#L53-L77)
- **Command execution pattern**: Uses `exec.Command()` - maintain this approach
- **Testing examples**: See [`main_test.go`](main_test.go) for testing patterns
- **Error handling style**: Follow existing patterns in the codebase

## How to Contribute

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/linux-support`
3. **Implement the changes** following the checklist above
4. **Test thoroughly** on different Linux distributions if possible
5. **Update documentation** (README.md, comments)
6. **Submit a pull request** with:
   - Clear description of your implementation choice
   - Testing details (which distros/tools you tested)
   - Any limitations or known issues

## Questions for Contributors

- Which Linux PDF tool would you recommend as the primary choice? Why?
- Should we implement a fallback chain (try pdfunite, then gs, then pdftk)?
- Any experience with PDF joining tools on specific Linux distributions?

## Recognition

Contributors who successfully implement Linux support will be:
- Added to the project's contributors list
- Mentioned in release notes
- Given co-maintainer consideration for ongoing development

---

**Let's make this tool truly cross-platform! 🚀**

*This is a great opportunity for developers familiar with Linux systems and Go programming to contribute to an open-source project used by the community.*

---

### Labels
`enhancement` `help wanted` `good first issue` `linux` `cross-platform` 