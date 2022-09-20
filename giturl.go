package giturl

import "fmt"

type ProtocolType int8

const (
	ProtocolTypeSSH ProtocolType = iota
	ProtocolTypeGit
	ProtocolTypeHTTP
	ProtocolTypeHTTPs
	ProtocolTypeFTP
	ProtocolTypeFTPs
	ProtocolTypeSCP
)

const (
	ImplicitUser = ""
	DefaultUser  = "git"
)

const (
	ImplicitPort     uint16 = 0
	DefaultSSHPort   uint16 = 22
	DefaultGitPort   uint16 = 9418
	DefaultHTTPPort  uint16 = 80
	DefaultHTTPsPort uint16 = 443
	DefaultFTPPort   uint16 = 21
	DefaultFTPsPort  uint16 = 990
)

const (
	prefixSSH   = "ssh://"
	prefixGit   = "git://"
	prefixHTTP  = "http://"
	prefixHTTPs = "https://"
	prefixFTP   = "ftp://"
	prefixFTPs  = "ftps://"
	suffixGit   = ".git"
)

type GitURL struct {
	protocol ProtocolType
	port     uint16
	user     string
	host     string
	path     string
	repo     string
	raw      string
}

func NewGitURL(url string) (*GitURL, error) {
	return parseURL(url)
}

func (g GitURL) Port() uint16 {
	return g.port
}

func (g GitURL) Protocol() ProtocolType {
	return g.protocol
}

func (g GitURL) User() string {
	return g.user
}

func (g GitURL) Host() string {
	return g.host
}

func (g GitURL) Path() string {
	return g.path
}

func (g GitURL) Repo() string {
	return g.repo
}

func (g GitURL) RawURL() string {
	return g.raw
}

// format: ssh://[user@]host.xz[:port]/path/to/repo.git
func (g GitURL) ToSSHFormat(user string, port uint16, withSuffix bool) string {
	var userPart, hostPart, pathPart string
	if user != ImplicitUser {
		userPart = user + "@"
	}
	hostPart = g.host
	if port != ImplicitPort {
		hostPart += fmt.Sprintf(":%d", port)
	}
	pathPart = fmt.Sprintf("/%s/%s", g.path, g.repo)
	if withSuffix {
		pathPart += suffixGit
	}
	return prefixSSH + userPart + hostPart + pathPart
}

// format: git://host.xz[:port]/path/to/repo.git
func (g GitURL) ToGitFormat(port uint16, withSuffix bool) string {
	var hostPart, pathPart string
	hostPart = g.host
	if port != ImplicitPort {
		hostPart += fmt.Sprintf(":%d", port)
	}
	pathPart = fmt.Sprintf("/%s/%s", g.path, g.repo)
	if withSuffix {
		pathPart += suffixGit
	}
	return prefixGit + hostPart + pathPart
}

// format: http[s]://host.xz[:port]/path/to/repo.git
func (g GitURL) ToHTTPFormat(port uint16, isSecure bool, withSuffix bool) string {
	var prefixPart, hostPart, pathPart string
	if isSecure {
		prefixPart = prefixHTTPs
	} else {
		prefixPart = prefixHTTP
	}
	hostPart = g.host
	if port != ImplicitPort {
		hostPart += fmt.Sprintf(":%d", port)
	}
	pathPart = fmt.Sprintf("/%s/%s", g.path, g.repo)
	if withSuffix {
		pathPart += suffixGit
	}
	return prefixPart + hostPart + pathPart
}

// format: ftp[s]://host.xz[:port]/path/to/repo.git
func (g GitURL) ToFTPFormat(port uint16, isSecure bool, withSuffix bool) string {
	var prefixPart, hostPart, pathPart string
	if isSecure {
		prefixPart = prefixFTPs
	} else {
		prefixPart = prefixFTP
	}
	hostPart = g.host
	if port != ImplicitPort {
		hostPart += fmt.Sprintf(":%d", port)
	}
	pathPart = fmt.Sprintf("/%s/%s", g.path, g.repo)
	if withSuffix {
		pathPart += suffixGit
	}
	return prefixPart + hostPart + pathPart
}

// format: [user@]host.xz:path/to/repo.git
func (g GitURL) ToSCPFormat(user string, withSuffix bool) string {
	var userPart, pathPart string
	if user != ImplicitUser {
		userPart = user + "@"
	}

	pathPart = fmt.Sprintf(":%s/%s", g.path, g.repo)
	if withSuffix {
		pathPart += suffixGit
	}
	return userPart + g.host + pathPart
}
