package awsdefault

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/go-ini/ini"
	"github.com/kylelemons/godebug/pretty"
)

var (
	testFileContent = []byte("")
	testFilePath    = "testdata/credentials"
)

func setup() {

}

func TestMain(m *testing.M) {
	// setup
	var err error
	testFileContent, err = ioutil.ReadFile("testdata/.aws/credentials")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	code := m.Run()
	// teardown
	err = os.Remove(testFilePath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(code)
}

func TestCredentialsFile_SetDefaultTo(t *testing.T) {
	content, err := ini.InsensitiveLoad(testFileContent)
	if err != nil {
		t.Fatalf("CredentialsFile.SetDefaultTo() error = %v", err)
	}
	type fields struct {
		Content *ini.File
		Path    string
	}
	type args struct {
		profileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "0positiv - add default section and set values to live values",
			fields: fields{
				Content: content,
				Path:    testFilePath,
			},
			args: args{
				profileName: "live",
			},
			wantErr: false,
		},
		{
			name: "1positiv - add default section and set values to dev values",
			fields: fields{
				Content: content,
				Path:    testFilePath,
			},
			args: args{
				profileName: "dev",
			},
			wantErr: false,
		},
		{
			name: "2negativ - given section not existing, but last active_profile is used",
			fields: fields{
				Content: content,
				Path:    testFilePath,
			},
			args: args{
				profileName: "xxxxxxx",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &CredentialsFile{
				Content: tt.fields.Content,
				Path:    tt.fields.Path,
			}
			if err := f.SetDefaultTo(tt.args.profileName); (err != nil) != tt.wantErr {
				t.Errorf("CredentialsFile.SetDefaultTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCredentialsFile(t *testing.T) {
	tests := []struct {
		name    string
		want    *CredentialsFile
		envVar  string
		envVal  string
		wantErr bool
	}{
		{
			name:    "0positiv - read valid credentials file from HOME",
			envVar:  "HOME",
			envVal:  "testdata",
			wantErr: false,
			want: &CredentialsFile{
				Path: "testdata/.aws/credentials",
			},
		},
		{
			name:    "1negativ - read malformed credentials file",
			envVar:  "HOME",
			envVal:  "testdata/malformed",
			wantErr: true,
			want: &CredentialsFile{
				Path:    "testdata/malformed/.aws/credentials",
				Content: nil,
			},
		},
		{
			name:    "2positiv - read valid credentials file from AWS_SHARED_CREDENTIALS_FILE",
			envVar:  "AWS_SHARED_CREDENTIALS_FILE",
			envVal:  "testdata/byEnv/.aws/credentials",
			wantErr: false,
			want: &CredentialsFile{
				Path: "testdata/byEnv/.aws/credentials",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv(tt.envVar, tt.envVal); err != nil {
				t.Errorf("GetCredentialsFile() error = %v", err)
				return
			}
			got, err := GetCredentialsFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentialsFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Path != tt.want.Path {
				t.Errorf("GetCredentialsFile() got path = %v, want path %v", got.Path, tt.want.Path)
				return
			}
			if (got.Content == nil) && (tt.want.Content != nil) {
				t.Errorf("GetCredentialsFile() got empty content")
			}
		})
	}
}

func TestCredentialsFile_GetProfilesNames(t *testing.T) {
	type fields struct {
		Content []byte
		Path    string
	}
	tests := []struct {
		name      string
		fields    fields
		wantNames []string
	}{
		{
			name: "0positiv - simple return the correct profile names",
			fields: fields{
				Content: testFileContent,
				Path:    testFilePath,
			},
			wantNames: []string{"dev", "live"},
		},
		{
			name: "1positiv - multiple profiles with same name",
			fields: fields{
				Content: []byte("[a]\n[A]\n[a]"),
				Path:    "",
			},
			wantNames: []string{"a"},
		},
		{
			name: "2positiv - returns, if malformed content, empty slice",
			fields: fields{
				Content: []byte("[1]\n[]"),
				Path:    "",
			},
			wantNames: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, _ := ini.InsensitiveLoad(tt.fields.Content)
			f := &CredentialsFile{
				Content: content,
				Path:    tt.fields.Path,
			}
			gotNames := f.GetProfilesNames()
			if diff := pretty.Compare(tt.wantNames, gotNames); diff != "" {
				t.Errorf(
					"%s: CredentialsFile.GetProfilesNames() diff: (-want +got)\n%s", tt.name, diff,
				)
				return
			}
		})
	}
}

func TestCredentialsFile_GetUsedProfileNameAndIndex(t *testing.T) {
	ini.DefaultHeader = true
	noMatch := []byte(`
	[default]
	aws_access_key_id=A
	aws_secret_access_key=B
	[live]
	aws_access_key_id=C
	aws_secret_access_key=B
	`)
	noDefault := []byte(`
	[dev]
	aws_access_key_id=A
	aws_secret_access_key=B
	[live]
	aws_access_key_id=C
	aws_secret_access_key=D
	`)
	type fields struct {
		Content []byte
		Path    string
	}
	tests := []struct {
		name      string
		fields    fields
		want      string
		wantIndex int
		wantErr   bool
	}{
		{
			name: "0positiv - get the correct active/default profile",
			fields: fields{
				Content: testFileContent,
				Path:    testFilePath,
			},
			want:      "dev",
			wantIndex: 0,
			wantErr:   false,
		},
		{
			name: "1negativ - default profile does not match to one of the others",
			fields: fields{
				Content: noMatch,
				Path:    testFilePath,
			},
			want:      "",
			wantIndex: -1,
			wantErr:   true,
		},
		{
			name: "2positiv - no default set",
			fields: fields{
				Content: noDefault,
				Path:    testFilePath,
			},
			want:      "no default",
			wantIndex: -2,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, _ := ini.InsensitiveLoad(tt.fields.Content)
			f := &CredentialsFile{
				Content: content,
				Path:    tt.fields.Path,
			}
			got, idx, err := f.GetUsedProfileNameAndIndex()
			if (err != nil) != tt.wantErr {
				t.Errorf("CredentialsFile.GetUsedProfileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CredentialsFile.GetUsedProfileName() = %v, want %v", got, tt.want)
			}
			if idx != tt.wantIndex {
				t.Errorf("CredentialsFile.GetUsedProfileName() = %v, want %v", idx, tt.wantIndex)
			}
		})
	}
}

func TestCredentialsFile_UnSetDefault(t *testing.T) {
	unsetTestFilePath := "testdata/unsetTests"
	rmFile := func() {
		_ = os.Remove(unsetTestFilePath)
	}
	good := []byte(`
	; comment
	[default]
	aws_access_key_id=A
	aws_secret_access_key=B
	[live]
	aws_access_key_id=C
	aws_secret_access_key=B
	`)
	expectGood := []byte(`[live]
aws_access_key_id     = C
aws_secret_access_key = B

`)

	noDefault := []byte(`[dev]
aws_access_key_id     = A
aws_secret_access_key = B

[live]
aws_access_key_id     = C
aws_secret_access_key = D

`)
	type fields struct {
		Path string
	}
	tests := []struct {
		name    string
		content []byte
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "positiv - remove default section and comment",
			content: good,
			fields:  fields{unsetTestFilePath},
			want:    expectGood,
			wantErr: false,
		},
		{
			name:    "positiv - no default section",
			content: noDefault,
			fields:  fields{unsetTestFilePath},
			want:    noDefault,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ini.InsensitiveLoad(tt.content)
			if err != nil {
				t.Errorf("CredentialsFile.UnSetDefault() error = %v", err)
			}
			f := &CredentialsFile{
				Content: content,
				Path:    tt.fields.Path,
			}
			if err := f.UnSetDefault(); (err != nil) != tt.wantErr {
				t.Errorf("CredentialsFile.UnSetDefault() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, err := ioutil.ReadFile(unsetTestFilePath)
			if err != nil {
				t.Errorf("CredentialsFile.UnSetDefault() error = %v", err)
				return
			}
			if string(b) != string(tt.want) {
				t.Errorf("CredentialsFile.UnSetDefault() got:\n%+#v, want:\n%+#v",
					string(b), string(tt.want))
			}
			rmFile()
		})
	}
}

func TestCredentialsFile_GetUsedID(t *testing.T) {
	ini.DefaultHeader = true
	good := []byte(`
	[default]
	aws_access_key_id=A
	aws_secret_access_key=B
	`)
	noDefault := []byte(`
	[dev]
	aws_access_key_id=A
	aws_secret_access_key=B
	`)
	noAWS := []byte(`
	[default]
	foo=bar
	[dev]
	foo=bar
	`)
	type fields struct {
		Content []byte
		Path    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "positiv — BAT",
			fields: fields{
				Content: good,
				Path:    testFilePath,
			},
			want:    "A",
			wantErr: false,
		},
		{
			name: "negativ — no default",
			fields: fields{
				Content: noDefault,
				Path:    testFilePath,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "negativ — no valid key id",
			fields: fields{
				Content: noAWS,
				Path:    testFilePath,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ini.InsensitiveLoad(tt.fields.Content)
			if err != nil {
				t.Errorf("CredentialsFile.UnSetDefault() error = %v", err)
			}
			f := &CredentialsFile{
				Content: content,
				Path:    tt.fields.Path,
			}
			got, err := f.GetUsedID()
			if (err != nil) != tt.wantErr {
				t.Errorf("CredentialsFile.GetUsedID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CredentialsFile.GetUsedID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentialsFile_GetUsedKey(t *testing.T) {
	ini.DefaultHeader = true
	good := []byte(`
	[default]
	aws_access_key_id=A
	aws_secret_access_key=B
	`)
	noDefault := []byte(`
	[dev]
	aws_access_key_id=A
	aws_secret_access_key=B
	`)
	noAWS := []byte(`
	[default]
	foo=bar
	[dev]
	foo=bar
	`)
	type fields struct {
		Content []byte
		Path    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "positiv — BAT",
			fields: fields{
				Content: good,
				Path:    testFilePath,
			},
			want:    "B",
			wantErr: false,
		},
		{
			name: "negativ — empty default",
			fields: fields{
				Content: noDefault,
				Path:    testFilePath,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "negativ — no valid key id",
			fields: fields{
				Content: noAWS,
				Path:    testFilePath,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ini.InsensitiveLoad(tt.fields.Content)
			if err != nil {
				t.Errorf("CredentialsFile.UnSetDefault() error = %v", err)
			}
			f := &CredentialsFile{
				Content: content,
				Path:    tt.fields.Path,
			}
			got, err := f.GetUsedKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("CredentialsFile.GetUsedID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CredentialsFile.GetUsedID() = %v, want %v", got, tt.want)
			}
		})
	}
}
