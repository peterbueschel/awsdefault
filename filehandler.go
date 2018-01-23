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
	// Profile stored in the AWS shared credentials file
	Profile struct {
		AccessKeyID     string `ini:"aws_access_key_id"`
		SecretAccessKey string `ini:"aws_secret_access_key"`
	}

	CredentialsFile struct {
		Content *ini.File
		Path    string
	}
)

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

func (f *CredentialsFile) GetUsedProfileNameAndIndex() (string, int, error) {
	d, err := f.GetProfile("default")
	if err != nil {
		return "", -1, err
	}
	if len(d.AccessKeyID) == 0 || len(d.SecretAccessKey) == 0 {
		return "", -1,
			fmt.Errorf(
				"%s contains no valid default profile or no default profile was set",
				f.Path,
			)
	}
	for i, n := range f.GetProfilesNames() {
		p, err := f.GetProfile(n)
		if err != nil {
			return "", -1, err
		}
		if p.AccessKeyID == d.AccessKeyID && p.SecretAccessKey == d.SecretAccessKey {
			return n, i, nil
		}
	}
	return "", -1,
		fmt.Errorf(
			"No profile in %s matches the current default-profile. "+
				"The AWS keys in the default profile are different to all "+
				"other keys in the other profiles.", f.Path,
		)

}

func (f *CredentialsFile) GetProfile(name string) (*Profile, error) {
	p := new(Profile)
	s, err := f.Content.GetSection(name)
	if err != nil {
		return nil, err
	}
	err = s.MapTo(p)
	return p, err
}

func (f *CredentialsFile) SetDefaultTo(profileName string) error {
	p, err := f.GetProfile(profileName)
	if err != nil {
		return err
	}
	f.Content.Section("default").Comment = commentPrefix + profileName
	err = f.Content.Section("default").ReflectFrom(p)
	if err != nil {
		return err
	}
	return f.Content.SaveTo(f.Path)
}

func (f *CredentialsFile) UnSetDefault() error {
	f.Content.DeleteSection("default")
	return f.Content.SaveTo(f.Path)
}
