package giturl

import (
	"errors"
	"testing"
)

var ParseURLTestCases = []struct {
	URL      string
	Protocol ProtocolType
	Port     uint16
	User     string
	Host     string
	Path     string
	Repo     string
	Err      error
}{
	{
		URL:      "ssh://git@gitlab.com/charlie/wto/bomb.git",
		Protocol: ProtocolTypeSSH,
		Port:     22,
		User:     "git",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "ssh://git@gitlab.com:22/charlie/wto/bomb.git",
		Protocol: ProtocolTypeSSH,
		Port:     22,
		User:     "git",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "ssh://gitlab.com:22/charlie/wto/bomb.git",
		Protocol: ProtocolTypeSSH,
		Port:     22,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "git://gitlab.com/charlie/wto/bomb.git",
		Protocol: ProtocolTypeGit,
		Port:     9418,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "git://gitlab.com:12345/charlie/wto/bomb.git",
		Protocol: ProtocolTypeGit,
		Port:     12345,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "http://gitlab.com/charlie/wto/bomb.git",
		Protocol: ProtocolTypeHTTP,
		Port:     80,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "http://gitlab.com:8080/charlie/wto/bomb.git",
		Protocol: ProtocolTypeHTTP,
		Port:     8080,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "https://gitlab.com/charlie/wto/bomb.git",
		Protocol: ProtocolTypeHTTPs,
		Port:     443,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "ftp://gitlab.com/charlie/wto/bomb.git",
		Protocol: ProtocolTypeFTP,
		Port:     21,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "ftps://gitlab.com/charlie/wto/bomb.git",
		Protocol: ProtocolTypeFTPs,
		Port:     990,
		User:     "",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "git@gitlab.com:charlie/wto/bomb.git",
		Protocol: ProtocolTypeSCP,
		Port:     22,
		User:     "git",
		Host:     "gitlab.com",
		Path:     "charlie/wto",
		Repo:     "bomb",
		Err:      nil,
	},
	{
		URL:      "ssh://admin@119.91.33.58:29418/All-Projects",
		Protocol: ProtocolTypeSSH,
		Port:     29418,
		User:     "admin",
		Host:     "119.91.33.58",
		Path:     "",
		Repo:     "All-Projects",
		Err:      nil,
	},
	{
		URL: "ssh://admin@119.91.33.58:65536/All-Projects",
		Err: ErrInvalidURL,
	},
	{
		URL: "git@gitlab.com/charlie/wto/bomb.git",
		Err: ErrInvalidURL,
	},
}

func TestParseURL(t *testing.T) {
	for _, c := range ParseURLTestCases {
		URL, err := NewGitURL(c.URL)
		if err != nil {
			if !errors.Is(err, c.Err) {
				t.Error(err.Error())
			}
			continue
		}

		if actual, expected := URL.Protocol(), c.Protocol; actual != expected {
			t.Errorf("test failed, actual protocol: %d, expected: %d", actual, expected)
		}
		if actual, expected := URL.User(), c.User; actual != expected {
			t.Errorf("test failed, actual user: %s, expected: %s", actual, expected)
		}
		if actual, expected := URL.Host(), c.Host; actual != expected {
			t.Errorf("test failed, actual host: %s, expected: %s", actual, expected)
		}
		if actual, expected := URL.Port(), c.Port; actual != expected {
			t.Errorf("test failed, actual port: %d, expected: %d", actual, expected)
		}
		if actual, expected := URL.Path(), c.Path; actual != expected {
			t.Errorf("test failed, actual path: %s, expected: %s", actual, expected)
		}
		if actual, expected := URL.Repo(), c.Repo; actual != expected {
			t.Errorf("test failed, actual repo: %s, expected: %s", actual, expected)
		}
	}
}

func TestToSSHFormat(t *testing.T) {
	testCases := []struct {
		InURL        string
		InUser       string
		InPort       uint16
		InWithSuffix bool
		OutAddr      string
	}{
		{
			InURL:        "http://gitlab.com/charlie/wto/bomb.git",
			InUser:       ImplicitUser,
			InPort:       ImplicitPort,
			InWithSuffix: false,
			OutAddr:      "ssh://gitlab.com/charlie/wto/bomb",
		},
		{
			InURL:        "http://gitlab.com/charlie/wto/bomb.git",
			InUser:       DefaultUser,
			InPort:       ImplicitPort,
			InWithSuffix: false,
			OutAddr:      "ssh://git@gitlab.com/charlie/wto/bomb",
		},
		{
			InURL:        "http://gitlab.com/charlie/wto/bomb.git",
			InUser:       DefaultUser,
			InPort:       DefaultSSHPort,
			InWithSuffix: false,
			OutAddr:      "ssh://git@gitlab.com:22/charlie/wto/bomb",
		},
		{
			InURL:        "http://gitlab.com/charlie/wto/bomb.git",
			InUser:       DefaultUser,
			InPort:       DefaultSSHPort,
			InWithSuffix: true,
			OutAddr:      "ssh://git@gitlab.com:22/charlie/wto/bomb.git",
		},
	}

	for _, c := range testCases {
		URL, err := NewGitURL(c.InURL)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		if actual, expected := URL.ToSSHFormat(c.InUser, c.InPort, c.InWithSuffix), c.OutAddr; actual != expected {
			t.Errorf("test failed, actual addr: %s, expected: %s", actual, expected)
		}
	}
}

func TestToHTTPFormat(t *testing.T) {
	testCases := []struct {
		InURL        string
		InPort       uint16
		InIsSecure   bool
		InWithSuffix bool
		OutAddr      string
	}{
		{
			InURL:        "ssh://git@gitlab.com/charlie/wto/bomb.git",
			InPort:       ImplicitPort,
			InIsSecure:   false,
			InWithSuffix: false,
			OutAddr:      "http://gitlab.com/charlie/wto/bomb",
		},
		{
			InURL:        "ssh://git@gitlab.com/charlie/wto/bomb.git",
			InPort:       DefaultHTTPPort,
			InIsSecure:   false,
			InWithSuffix: false,
			OutAddr:      "http://gitlab.com:80/charlie/wto/bomb",
		},
		{
			InURL:        "ssh://git@gitlab.com/charlie/wto/bomb.git",
			InPort:       DefaultHTTPsPort,
			InIsSecure:   true,
			InWithSuffix: false,
			OutAddr:      "https://gitlab.com:443/charlie/wto/bomb",
		},
		{
			InURL:        "ssh://git@gitlab.com/charlie/wto/bomb.git",
			InPort:       DefaultHTTPsPort,
			InIsSecure:   true,
			InWithSuffix: true,
			OutAddr:      "https://gitlab.com:443/charlie/wto/bomb.git",
		},
	}

	for _, c := range testCases {
		URL, err := NewGitURL(c.InURL)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		if actual, expected := URL.ToHTTPFormat(c.InPort, c.InIsSecure, c.InWithSuffix), c.OutAddr; actual != expected {
			t.Errorf("test failed, actual addr: %s, expected: %s", actual, expected)
		}
	}
}

func TestToFTPFormat(t *testing.T) {
	testCases := []struct {
		InURL        string
		InPort       uint16
		InIsSecure   bool
		InWithSuffix bool
		OutAddr      string
	}{
		{
			InURL:        "git@gitlab.com:charlie/wto/bomb.git",
			InPort:       ImplicitPort,
			InIsSecure:   false,
			InWithSuffix: false,
			OutAddr:      "ftp://gitlab.com/charlie/wto/bomb",
		},
		{
			InURL:        "git@gitlab.com:charlie/wto/bomb.git",
			InPort:       DefaultFTPPort,
			InIsSecure:   false,
			InWithSuffix: false,
			OutAddr:      "ftp://gitlab.com:21/charlie/wto/bomb",
		},
		{
			InURL:        "git@gitlab.com:charlie/wto/bomb.git",
			InPort:       DefaultFTPsPort,
			InIsSecure:   true,
			InWithSuffix: false,
			OutAddr:      "ftps://gitlab.com:990/charlie/wto/bomb",
		},
		{
			InURL:        "git@gitlab.com:charlie/wto/bomb.git",
			InPort:       DefaultHTTPsPort,
			InIsSecure:   true,
			InWithSuffix: true,
			OutAddr:      "ftps://gitlab.com:443/charlie/wto/bomb.git",
		},
	}

	for _, c := range testCases {
		URL, err := NewGitURL(c.InURL)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		if actual, expected := URL.ToFTPFormat(c.InPort, c.InIsSecure, c.InWithSuffix), c.OutAddr; actual != expected {
			t.Errorf("test failed, actual addr: %s, expected: %s", actual, expected)
		}
	}
}

func TestToSCPFormat(t *testing.T) {
	testCases := []struct {
		InURL        string
		InUser       string
		InWithSuffix bool
		OutAddr      string
	}{
		{
			InURL:        "ftps://gitlab.com:443/charlie/wto/bomb.git",
			InUser:       ImplicitUser,
			InWithSuffix: false,
			OutAddr:      "gitlab.com:charlie/wto/bomb",
		},
		{
			InURL:        "ftps://gitlab.com:443/charlie/wto/bomb.git",
			InUser:       DefaultUser,
			InWithSuffix: false,
			OutAddr:      "git@gitlab.com:charlie/wto/bomb",
		},
		{
			InURL:        "ftps://gitlab.com:443/charlie/wto/bomb.git",
			InUser:       DefaultUser,
			InWithSuffix: true,
			OutAddr:      "git@gitlab.com:charlie/wto/bomb.git",
		},
	}

	for _, c := range testCases {
		URL, err := NewGitURL(c.InURL)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		if actual, expected := URL.ToSCPFormat(c.InUser, c.InWithSuffix), c.OutAddr; actual != expected {
			t.Errorf("test failed, actual addr: %s, expected: %s", actual, expected)
		}
	}
}
