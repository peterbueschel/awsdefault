//Copyright 2018 Peter Büschel
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.package awsdefault

package awsdefault

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/go-ini/ini"
)

const (
	commentPrefix = "active_profile="
)

type (
	// Profile stored in the AWS shared credentials file consisting of an
	// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
	Profile struct {
		Name            string `ini:"name"`
		AccessKeyID     string `ini:"aws_access_key_id"`
		SecretAccessKey string `ini:"aws_secret_access_key"`
	}

	// CredentialsFile stores the content and path of the AWS credentials file
	CredentialsFile struct {
		Content *ini.File
		Path    string
	}
)

// GetCredentialsFile reads the AWS credentials file either from the HOME directory or
// from a path given by the environment variable AWS_SHARED_CREDENTIALS_FILE
func GetCredentialsFile() (*CredentialsFile, error) {
	home := func() string {
		if runtime.GOOS == "windows" {
			return os.Getenv("USERPROFILE")
		}
		return os.Getenv("HOME")
	}
	path := filepath.Join(home(), ".aws", "credentials")
	if p := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); len(p) > 0 {
		path = p
	}
	ini.DefaultHeader = true
	f, err := ini.InsensitiveLoad(path)
	return &CredentialsFile{f, path}, err
}

// GetProfilesNames returns a sorted list of all available profiles inside the AWS credentials file.
func (f *CredentialsFile) GetProfilesNames() (names []string) {
	if f.Content != nil {
		for _, p := range f.Content.SectionStrings() {
			if strings.ToLower(p) != "default" {
				names = append(names, p)
			}
		}
	}
	sort.Strings(names)
	return
}

// GetUsedProfileNameAndIndex returns the name and the index of the profile currently used as default
// profile.
func (f *CredentialsFile) GetUsedProfileNameAndIndex() (string, int, error) {
	d, err := f.GetProfileBy("default")
	if err != nil {
		return "", -1, err
	}
	if len(d.AccessKeyID) == 0 || len(d.SecretAccessKey) == 0 {
		return "no default set", -2, nil
	}
	for i, n := range f.GetProfilesNames() {
		if n == d.Name {
			return n, i, nil
		}
	}
	return "", -1,
		fmt.Errorf(
			"No profile in %s matches the current configured default-profile '%s'.",
			f.Path, d.Name,
		)

}

// GetUsedID returns the AWS_ACCESS_KEY_ID of the profile currently used as default profile.
func (f *CredentialsFile) GetUsedID() (string, error) {
	d, err := f.GetProfileBy("default")
	return d.AccessKeyID, err
}

// GetUsedKey returns the AWS_SECRET_ACCESS_KEY of the profile currently used as default profile.
func (f *CredentialsFile) GetUsedKey() (string, error) {
	d, err := f.GetProfileBy("default")
	return d.SecretAccessKey, err
}

// GetProfileBy returns the profile by a given name
func (f *CredentialsFile) GetProfileBy(name string) (*Profile, error) {
	p := &Profile{Name: name}
	s, err := f.Content.GetSection(name)
	if err != nil {
		return p, err
	}
	_ = s.MapTo(p) // error cannot happen; p is always a pointer
	return p, nil
}

// SetDefaultTo overwrites/creates the default section inside the AWS credentials file.
// It also adds a comment containing the name of the profile used as current default profile
func (f *CredentialsFile) SetDefaultTo(profileName string) error {
	p, err := f.GetProfileBy(profileName)
	if err != nil {
		return err
	}
	_ = f.Content.Section("default").ReflectFrom(p) // error cannot happen; p is always a pointer
	return f.Content.SaveTo(f.Path)
}

// UnSetDefault deletes the default section inside the AWS credentials file.
func (f *CredentialsFile) UnSetDefault() error {
	f.Content.DeleteSection("default")
	return f.Content.SaveTo(f.Path)
}
