package ls

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	iconMap = map[string]string{
		".md":      "󰍔 ", // Markdown file
		".txt":     " ", // Plain text file
		".doc":     "󰈬 ", // Microsoft Word document
		".docx":    "󰈬 ", // Microsoft Word Open XML document
		".pdf":     " ", // Portable Document Format
		".xls":     "󰈛 ", // Microsoft Excel spreadsheet
		".xlsx":    "󰈛 ", // Microsoft Excel Open XML spreadsheet
		".ppt":     "󰈧 ", // Microsoft PowerPoint presentation
		".pptx":    "󰈧 ", // Microsoft PowerPoint Open XML presentation
		".rtf":     " ", // Rich Text Format document
		".odt":     "󰈬 ", // OpenDocument Text document
		".ods":     "󰈛 ", // OpenDocument Spreadsheet
		".odp":     "󰈧 ", // OpenDocument Presentation
		".csv":     " ", // Comma-Separated Values file
		".tsv":     " ", // Tab-Separated Values file
		".html":    " ", // HyperText Markup Language file
		".htm":     " ", // HTML file
		".xml":     "󰗀 ", // eXtensible Markup Language file
		".xhtml":   "󰗀 ", // eXtensible HTML file
		".css":     " ", // Cascading Style Sheets file
		".js":      " ", // JavaScript file
		".json":    " ", // JavaScript Object Notation file
		".php":     " ", // PHP script file
		".jsp":     " ", // Java Server Pages file
		".java":    " ", // Java source code file
		".py":      "󰌠 ", // Python script file
		".rb":      "󰴭 ", // Ruby script file
		".cpp":     "󰙲 ", // C++ source code file
		".c":       " ", // C source code file
		".cs":      " ", // C# source code file
		".go":      " ", // Go programming language source file
		".swift":   " ", // Swift source code file
		".kt":      " ", // Kotlin source code file
		".sql":     " ", // Structured Query Language file
		".db":      " ", // Database file
		".sqlite":  " ", // SQLite database file
		".bak":     "󰁯 ", // Backup file
		".log":     "󱂅 ", // Log file
		".sh":      " ", // Shell script (Unix/Linux)
		".bat":     " ", // Batch file (Windows)
		".ps1":     " ", // PowerShell script (Windows)
		".pl":      " ", // Perl script file
		".r":       " ", // R script file
		".jar":     " ", // Java Archive file
		".war":     " ", // Web Archive file (Java)
		".ear":     " ", // Enterprise Archive file (Java)
		".exe":     "󰨡 ", // Executable file (Windows)
		".dll":     "󰨡 ", // Dynamic Link Library (Windows)
		".sys":     "󰨡 ", // System file (Windows)
		".msi":     "󰨡 ", // Microsoft Installer package
		".deb":     " ", // Debian package file (Linux)
		".rpm":     " ", // RPM Package Manager file (Linux)
		".apk":     "󰀲 ", // Android application package
		".ipa":     " ", // iOS application archive
		".iso":     " ", // Disk image file (ISO)
		".img":     "󰨣 ", // Disk image file
		".bin":     " ", // Binary file
		".cue":     "󱔼",  // Cue sheet file
		".vhd":     "󰋊 ", // Virtual Hard Disk file
		".vmdk":    "󰋊 ", // VMware virtual disk file
		".dmg":     "󰋊 ", // Apple Disk Image file (macOS)
		".zip":     "󰿺 ", // ZIP compressed archive
		".rar":     "󰿺 ", // RAR compressed archive
		".7z":      "󰿺 ", // 7-Zip compressed archive
		".tar":     "󰿺 ", // Tarball archive
		".gz":      "󰿺 ", // Gzip compressed file
		".bz2":     "󰿺 ", // Bzip2 compressed file
		".xz":      "󰿺 ", // XZ compressed file
		".zst":     "󰿺 ", // Zstandard compressed file (Linux/Unix)
		".cpio":    "󰿺 ", // CPIO archive file (Unix/Linux)
		".mp3":     "󱑽 ", // MPEG Audio Layer III file
		".wav":     "󱑽 ", // Waveform Audio File
		".flac":    "󱑽 ", // Free Lossless Audio Codec file
		".aac":     "󱑽 ", // Advanced Audio Coding file
		".ogg":     "󱑽 ", // Ogg Vorbis audio file
		".wma":     "󱑽 ", // Windows Media Audio file
		".m4a":     "󱑽 ", // MPEG-4 Audio file
		".aiff":    "󱑽 ", // Audio Interchange File Format
		".avi":     "󰈫 ", // Audio Video Interleave file
		".mp4":     "󰈫 ", // MPEG-4 video file
		".mkv":     "󰈫 ", // Matroska video file
		".mov":     "󰈫 ", // Apple QuickTime movie file
		".wmv":     "󰈫 ", // Windows Media Video file
		".flv":     "󰈫 ", // Flash video file
		".mpeg":    "󰈫 ", // MPEG video file
		".mpg":     "󰈫 ", // MPEG video file
		".m4v":     "󰈫 ", // iTunes video file
		".3gp":     "󰈫 ", // 3GPP multimedia file
		".3g2":     "󰈫 ", // 3GPP2 multimedia file
		".swf":     " ", // Shockwave Flash file
		".vob":     "󰈫 ", // DVD Video Object file
		".svg":     "󰕠 ", // Scalable Vector Graphics file
		".png":     " ", // Portable Network Graphics image
		".jpg":     " ", // JPEG image file
		".jpeg":    " ", // JPEG image file
		".gif":     " ", // Graphics Interchange Format image
		".bmp":     " ", // Bitmap image file
		".tiff":    " ", // Tagged Image File Format
		".ico":     " ", // Icon file
		".heic":    " ", // High Efficiency Image File Format
		".psd":     " ", // Adobe Photoshop Document
		".ai":      " ", // Adobe Illustrator file
		".sketch":  "󰉼 ", // Sketch design file
		".xcf":     " ", // GIMP image file
		".raw":     "󱨏 ", // Raw image format (camera)
		".cdr":     "󰇞 ", // CorelDRAW file
		".tex":     " ", // LaTeX document file
		".cfg":     " ", // Configuration file
		".ini":     " ", // Initialization file
		".conf":    " ", // Configuration file (commonly used in Unix/Linux)
		".env":     " ", // Environment configuration file
		".cmd":     "",   // Command script (Windows)
		".reg":     " ", // Windows Registry file
		".vcf":     " ", // vCard file (contact information)
		".ics":     " ", // iCalendar file (calendar events)
		".mobi":    " ", // Mobipocket eBook file
		".epub":    " ", // eBook file format
		".azw":     " ", // Amazon Kindle eBook file
		".ttf":     " ", // TrueType Font file
		".otf":     " ", // OpenType Font file
		".fon":     " ", // Generic font file
		".mht":     "󰖟 ", // MHTML web archive file
		".mhtml":   "󰖟 ", // MHTML file
		".part":    "󱑢 ", // Partial download file
		".tmp":     " ", // Temporary file
		".desktop": " ", // Linux desktop entry file (application launcher)
		".service": " ", // Systemd service unit file (Linux service)
		".ko":      " ", // Linux kernel module file
		".run":     " ", // Linux executable installer script
	}
)

func getIcon(ext string) string {
	if icon, ok := iconMap[ext]; ok {
		return icon
	}
	return " "
}

// getFileNameWithIcon returns a filename prefixed with an icon based on its type.
func getFileNameWithIcon(entry os.DirEntry) string {
	if entry.IsDir() {
		// Use a folder icon for directories.
		return " " + entry.Name()
	}

	// Convert file extension to lowercase for case-insensitive matching.
	ext := strings.ToLower(filepath.Ext(entry.Name()))

	// Return file name with an appropriate icon based on its extension or specific file names.
	switch {

	case entry.Name() == "go.mod", entry.Name() == "go.sum":
		return "󰟓 " + entry.Name() // Icon for Go files.

	case entry.Name() == "Dockerfile", entry.Name() == "docker-compose.yml", entry.Name() == ".dockerignore":
		return " " + entry.Name() // Icon for Dockerfile.

	case entry.Name() == "cargo.toml":
		return " " + entry.Name() // Icon for Rust files.

	case entry.Name() == ".github", entry.Name() == ".gitignore":
		return " " + entry.Name() // Icon for Git-related files.

	case entry.Name() == "Makefile":
		return " " + entry.Name() // Icon for Makefile.

	default:
		return getIcon(ext) + entry.Name()
	}
}
