package main

import (
	"testing"

	"github.com/peterbueschel/awsdefault"
)

func Test_getProfiles(t *testing.T) {
	type args struct {
		file *awsdefault.CredentialsFile
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "positive — BAT",
			args: args{
				file: &awsdefault.CredentialsFile{
					Path: "somewhere",
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProfiles(tt.args.file); tt.wantNil && (got == nil) {
				t.Errorf("getProfiles() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_getUsedProfile(t *testing.T) {
	type args struct {
		file *awsdefault.CredentialsFile
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "positive — BAT",
			args: args{
				file: &awsdefault.CredentialsFile{
					Path: "somewhere",
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUsedProfile(tt.args.file); tt.wantNil && (got == nil) {
				t.Errorf("getUsedProfile() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_setDefaultProfile(t *testing.T) {
	type args struct {
		file *awsdefault.CredentialsFile
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "positive — BAT",
			args: args{
				file: &awsdefault.CredentialsFile{
					Path: "somewhere",
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setDefaultProfile(tt.args.file); tt.wantNil && (got == nil) {
				t.Errorf("setDefaultProfile() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_unsetDefaultProfile(t *testing.T) {
	type args struct {
		file *awsdefault.CredentialsFile
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "positive — BAT",
			args: args{
				file: &awsdefault.CredentialsFile{
					Path: "somewhere",
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unsetDefaultProfile(tt.args.file); tt.wantNil && (got == nil) {
				t.Errorf("unsetDefaultProfile() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}
