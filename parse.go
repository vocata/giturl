package giturl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidURL = errors.New("invalid url")

// doc: https://www.git-scm.com/docs/git-clone#URLS
func parseURL(url string) (*GitURL, error) {
	if hasPrefix(url, prefixSSH) {
		return parseSSHURL(url)
	} else if hasPrefix(url, prefixGit) {
		return parseGitURL(url)
	} else if hasPrefix(url, prefixHTTP) || hasPrefix(url, prefixHTTPs) {
		return parseHTTPURL(url)
	} else if hasPrefix(url, prefixFTP) || hasPrefix(url, prefixFTPs) {
		return parseFTPURL(url)
	} else {
		return parseSCPURL(url)
	}
}

func parseSSHURL(url string) (*GitURL, error) {
	gitURL := &GitURL{
		raw:      url,
		protocol: ProtocolTypeSSH,
	}

	// pre-processing
	left := removePrefix(removeSuffix(url, suffixGit), prefixSSH)

	user, after, ok := cut(left, "@")
	if ok {
		gitURL.user = user
		left = after
	}

	before, left, ok := cut(left, "/")
	if !ok {
		return nil, fmt.Errorf("%w, missing path to repo", ErrInvalidURL)
	}

	host, port, ok := cut(before, ":")
	if ok {
		if p, err := strconv.ParseUint(port, 10, 16); err != nil {
			return nil, fmt.Errorf("%w, illegal port '%s'", ErrInvalidURL, port)
		} else {
			gitURL.port = uint16(p)
		}
	} else {
		gitURL.port = DefaultSSHPort
	}
	gitURL.host = host
	gitURL.path, gitURL.repo, _ = lastCut(left, "/")

	return gitURL, nil
}

func parseGitURL(url string) (*GitURL, error) {
	gitURL := &GitURL{
		protocol: ProtocolTypeGit,
		raw:      url,
	}

	left := removePrefix(removeSuffix(url, suffixGit), prefixGit)

	before, left, ok := cut(left, "/")
	if !ok {
		return nil, fmt.Errorf("%w, missing path to repo", ErrInvalidURL)
	}

	host, port, ok := cut(before, ":")
	if ok {
		if p, err := strconv.ParseUint(port, 10, 16); err != nil {
			return nil, fmt.Errorf("%w, illegal port '%s'", ErrInvalidURL, port)
		} else {
			gitURL.port = uint16(p)
		}
	} else {
		gitURL.port = DefaultGitPort
	}
	gitURL.host = host
	gitURL.path, gitURL.repo, _ = lastCut(left, "/")

	return gitURL, nil
}

func parseHTTPURL(url string) (*GitURL, error) {
	gitURL := &GitURL{
		raw: url,
	}
	left := removeSuffix(url, suffixGit)
	if hasPrefix(url, prefixHTTP) {
		gitURL.protocol = ProtocolTypeHTTP
		left = removePrefix(left, prefixHTTP)
	}
	if hasPrefix(url, prefixHTTPs) {
		gitURL.protocol = ProtocolTypeHTTPs
		left = removePrefix(left, prefixHTTPs)
	}

	before, left, ok := cut(left, "/")
	if !ok {
		return nil, fmt.Errorf("%w, missing path to repo", ErrInvalidURL)
	}

	host, port, ok := cut(before, ":")
	if ok {
		if p, err := strconv.ParseUint(port, 10, 16); err != nil {
			return nil, fmt.Errorf("%w, illegal port '%s'", ErrInvalidURL, port)
		} else {
			gitURL.port = uint16(p)
		}
	} else {
		if gitURL.protocol == ProtocolTypeHTTP {
			gitURL.port = DefaultHTTPPort
		}
		if gitURL.protocol == ProtocolTypeHTTPs {
			gitURL.port = DefaultHTTPsPort
		}
	}
	gitURL.host = host
	gitURL.path, gitURL.repo, _ = lastCut(left, "/")

	return gitURL, nil
}

func parseFTPURL(url string) (*GitURL, error) {
	gitURL := &GitURL{
		raw: url,
	}
	left := removeSuffix(url, suffixGit)
	if hasPrefix(url, prefixFTP) {
		gitURL.protocol = ProtocolTypeFTP
		left = removePrefix(left, prefixFTP)
	}
	if hasPrefix(url, prefixFTPs) {
		gitURL.protocol = ProtocolTypeFTPs
		left = removePrefix(left, prefixFTPs)
	}

	before, left, ok := cut(left, "/")
	if !ok {
		return nil, fmt.Errorf("%w, missing path to repo", ErrInvalidURL)
	}

	host, port, ok := cut(before, ":")
	if ok {
		if p, err := strconv.ParseUint(port, 10, 16); err != nil {
			return nil, fmt.Errorf("%w, illegal port '%s'", ErrInvalidURL, port)
		} else {
			gitURL.port = uint16(p)
		}
	} else {
		if gitURL.protocol == ProtocolTypeFTP {
			gitURL.port = DefaultFTPPort
		}
		if gitURL.protocol == ProtocolTypeFTPs {
			gitURL.port = DefaultFTPsPort
		}
	}
	gitURL.host = host
	gitURL.path, gitURL.repo, _ = lastCut(left, "/")

	return gitURL, nil
}

func parseSCPURL(url string) (*GitURL, error) {
	gitURL := &GitURL{
		port:     DefaultSSHPort,
		protocol: ProtocolTypeSCP,
		raw:      url,
	}

	left := removeSuffix(url, suffixGit)

	user, after, ok := cut(left, "@")
	if ok {
		gitURL.user = user
		left = after
	}

	host, left, ok := cut(left, ":")
	if !ok {
		return nil, fmt.Errorf("%w, expected ':'", ErrInvalidURL)
	}
	gitURL.host = host
	gitURL.path, gitURL.repo, _ = lastCut(left, "/")

	return gitURL, nil
}

// misc.
func hasPrefix(url string, prefix string) bool {
	if len(url) < len(prefix) {
		return false
	}
	return strings.EqualFold(url[:len(prefix)], prefix)
}

func removeSuffix(url string, suffix string) string {
	url = strings.TrimSuffix(url, "/")
	if len(url) < len(suffix) {
		return url
	}
	if strings.EqualFold(url[len(url)-len(suffix):], suffix) {
		return url[:len(url)-len(suffix)]
	}
	return url
}

func removePrefix(url string, prefix string) string {
	if len(url) < len(prefix) {
		return url
	}
	if strings.EqualFold(url[:len(prefix)], prefix) {
		return url[len(prefix):]
	}
	return url
}

func cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func lastCut(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return "", s, false
}
