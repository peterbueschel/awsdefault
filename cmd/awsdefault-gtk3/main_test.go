package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kylelemons/godebug/pretty"
	"github.com/peterbueschel/awsdefault"
)

func TestMain(m *testing.M) {
	gtk.Init(&os.Args)
	os.Exit(m.Run())
}

func Test_fetchProfiles(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     *profiles
		wantErr  bool
	}{
		{
			name:     "positive — BAT",
			filepath: "testdata/",
			want: &profiles{
				curr:    "dev",
				currIdx: 0,
				list:    []string{"dev", "live", noProfile},
				file:    &awsdefault.CredentialsFile{},
			},
			wantErr: false,
		},
		{
			name:     "positive — no used profile",
			filepath: "testdata/noprofile",
			want: &profiles{
				curr:    noProfile,
				currIdx: 2,
				list:    []string{"dev", "live", noProfile},
				file:    &awsdefault.CredentialsFile{},
			},
			wantErr: false,
		},
		{
			name:     "negative — error getting credentials file",
			filepath: "xxxxxxxx",
			want:     &profiles{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.filepath)
			got, err := fetchProfiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.curr != tt.want.curr {
				t.Errorf("fetchProfiles() curr got = %v, want %v", got.curr, tt.want.curr)
			}
			if got.currIdx != tt.want.currIdx {
				t.Errorf("fetchProfiles() currIdx got = %v, want %v", got.currIdx, tt.want.currIdx)
			}
			if len(got.list) != len(tt.want.list) {
				t.Errorf("fetchProfiles() list got = %v, want %v", got.list, tt.want.list)
			}
			if diff := pretty.Compare(tt.want.list, got.list); diff != "" {
				t.Errorf("fetchProfiles() diff: (-want +got)\n%s", diff)
				return
			}
		})
	}
}

func Test_chooser_setupListStore(t *testing.T) {
	type fields struct {
		selection *gtk.TreeSelection
		view      *gtk.TreeView
		store     *gtk.ListStore
		window    *gtk.Window
		box       *gtk.Box
		err       error
		profiles  *profiles
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			name:    "negative — c.err should return immediately",
			fields:  fields{err: fmt.Errorf("just an error")},
			wantNil: true,
			wantErr: true,
		},
		{
			name:    "positive — no profiles available",
			fields:  fields{profiles: &profiles{list: []string{}}},
			wantNil: false,
			wantErr: false,
		},
		{
			name:    "positive — BAT",
			fields:  fields{profiles: &profiles{list: []string{"dev", "live", noProfile}}},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &chooser{
				selection: tt.fields.selection,
				view:      tt.fields.view,
				store:     tt.fields.store,
				window:    tt.fields.window,
				box:       tt.fields.box,
				err:       tt.fields.err,
				profiles:  tt.fields.profiles,
			}
			c.setupListStore()
			if tt.wantNil && (c.store != nil) {
				t.Errorf("setupListStore() c.view not nil: %+v", c.store)
			}
			if (c.err != nil) != tt.wantErr {
				t.Errorf("setupListStore() error = %v , wantErr =: %v", c.err, tt.wantErr)
			}
		})
	}
}

func Test_chooser_setupTreeView(t *testing.T) {
	testLs, err := gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		t.Fatalf("setupTreeView(): could not create test ListStore: %s", err)
	}
	type fields struct {
		selection *gtk.TreeSelection
		view      *gtk.TreeView
		store     *gtk.ListStore
		window    *gtk.Window
		box       *gtk.Box
		err       error
		profiles  *profiles
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			name:    "negative — c.err should return immediately",
			fields:  fields{err: fmt.Errorf("just an error")},
			wantNil: true,
			wantErr: true,
		},
		{
			name:    "positive — BAT",
			fields:  fields{store: testLs},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &chooser{
				selection: tt.fields.selection,
				view:      tt.fields.view,
				store:     tt.fields.store,
				window:    tt.fields.window,
				box:       tt.fields.box,
				err:       tt.fields.err,
				profiles:  tt.fields.profiles,
			}
			c.setupTreeView()
			if tt.wantNil && (c.view != nil) {
				t.Errorf("setupTreeView() c.view not nil: %+v", c.view)
			}
			if (c.err != nil) != tt.wantErr {
				t.Errorf("setupTreeView() error = %v , wantErr =: %v", c.err, tt.wantErr)
			}
		})
	}
}

func Test_initializeChooser(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		withErr bool
	}{
		{
			name:    "positive — BAT",
			wantErr: false,
		},
		{
			name:    "negative — index is negative",
			withErr: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", "testdata/")
			testProfiles, err := fetchProfiles()
			if err != nil {
				t.Fatalf("initializeChooser(): could not fetch profiles: %s", err)
			}
			if tt.withErr {
				testProfiles.currIdx = -1
			}

			_, err = initializeChooser(testProfiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("initializeChooser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_showError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping gtk3-error test.")
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive — BAT",
			args: args{
				msg: "this is a test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := showError(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("showError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func Test_selectionChanged(t *testing.T) {

	tests := []struct {
		name            string
		credentialsPath string
		wantCurr        string
		isSelected      bool
		changedFile     bool
		wantErr         bool
	}{
		{
			name:            "positive — BAT",
			credentialsPath: "testdata/",
			wantCurr:        "dev",
			isSelected:      true,
			wantErr:         false,
		},
		{
			name:            "positive — BAT unset profile",
			credentialsPath: "testdata/noprofile/",
			wantCurr:        noProfile,
			isSelected:      true,
			wantErr:         false,
		},
		{
			name:            "positive — nothing changed - single click on same line",
			credentialsPath: "testdata/noprofile/",
			wantCurr:        noProfile,
			isSelected:      false,
			wantErr:         false,
		},
		{
			name:            "negative — profile is not anymore in the credentials file",
			credentialsPath: "testdata/",
			wantCurr:        "dev",
			changedFile:     true,
			isSelected:      true,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.credentialsPath)
			testProfiles, err := fetchProfiles()
			if err != nil {
				t.Fatalf("selectionChanged(): could not fetch profiles: %s", err)
			}
			c := &chooser{
				profiles: testProfiles,
			}
			if tt.changedFile {
				c.profiles.list = append(c.profiles.list, "xxxxxxx")
				c.profiles.currIdx = len(c.profiles.list) - 1
			}

			c.setupListStore()
			c.setupTreeView()
			c.selection, err = c.view.GetSelection()
			if err != nil {
				t.Fatalf("selectionChanged(): c.view.GetSelection() failed with: %s", err)
			}
			if tt.isSelected {
				path, err := gtk.TreePathNewFromString(fmt.Sprintf("%d", c.profiles.currIdx))
				if err != nil {
					t.Fatalf("selectionChanged(): could not get path: %s", err)
				}
				c.selection.SelectPath(path)
				c.view.RowActivated(path, c.view.GetColumn(0))
			}
			if err := c.selectionChanged(); (err != nil) != tt.wantErr {
				t.Errorf("selectionChanged() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if c.profiles.curr != tt.wantCurr {
				t.Errorf("current profile: selectionChanged() got = %v, want %v",
					c.curr, tt.wantCurr)
			}
		})
	}
}

func Test_chooser_setupSelection(t *testing.T) {

	tests := []struct {
		name         string
		negativeIdx  bool
		withTreeView bool
		withErr      bool
		wantNil      bool
		wantErr      bool
	}{
		{
			name:    "negative — c.err should return immediately",
			withErr: true,
			wantNil: true,
			wantErr: true,
		},
		{
			name:    "negative — TreeView is not set",
			wantNil: true,
			wantErr: true,
		},
		{
			name:         "negative — currIdx is negative",
			withTreeView: true,
			negativeIdx:  true,
			wantErr:      true,
		},
		{
			name:         "positive — BAT",
			withTreeView: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProfiles, err := fetchProfiles()
			if err != nil {
				t.Fatalf("selectionChanged(): could not fetch profiles: %s", err)
			}
			c := &chooser{
				profiles: testProfiles,
			}
			if tt.withErr {
				c.err = fmt.Errorf("just an error")
			}
			if tt.negativeIdx {
				c.profiles.currIdx = -1
			}
			if tt.withTreeView {
				c.setupTreeView()
				c.setupListStore()
			}
			c.setupSelection()
			if tt.wantNil && (c.selection != nil) {
				t.Errorf("setupSelection() c.selection not nil: %+v", c.view)
			}
			if (c.err != nil) != tt.wantErr {
				t.Errorf("setupSelection() error = %v , wantErr =: %v", c.err, tt.wantErr)
			}
		})
	}
}

func Test_chooser_setupRootBox(t *testing.T) {
	type fields struct {
		selection *gtk.TreeSelection
		view      *gtk.TreeView
		store     *gtk.ListStore
		window    *gtk.Window
		box       *gtk.Box
		err       error
		profiles  *profiles
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			name:    "negative — c.err should return immediately",
			fields:  fields{err: fmt.Errorf("just an error")},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &chooser{
				selection: tt.fields.selection,
				view:      tt.fields.view,
				store:     tt.fields.store,
				window:    tt.fields.window,
				box:       tt.fields.box,
				err:       tt.fields.err,
				profiles:  tt.fields.profiles,
			}
			c.setupRootBox()
			if tt.wantNil && (c.box != nil) {
				t.Errorf("setupRootBox() c.selection not nil: %+v", c.box)
			}
			if (c.err != nil) != tt.wantErr {
				t.Errorf("setupRootBox() error = %v , wantErr =: %v", c.err, tt.wantErr)
			}
		})
	}
}
