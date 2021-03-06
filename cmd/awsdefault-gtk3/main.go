package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/peterbueschel/awsdefault"
)

const (
	noProfile   = "--No Profile--"
	appName     = "awsdefault-ui"
	columnTitle = "Select AWS Profile"
)

var (
	permanent bool
)

type (
	profiles struct {
		curr    string
		currIdx int
		list    []string
		file    *awsdefault.CredentialsFile
	}
	chooser struct {
		selection *gtk.TreeSelection
		view      *gtk.TreeView
		store     *gtk.ListStore
		window    *gtk.Window
		box       *gtk.Box
		err       error
		*profiles
	}
)

func (c *chooser) selectionChanged() error {
	model, iter, ok := c.selection.GetSelected()
	if ok {
		tpath, err := model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			return err
		}
		iter, err := c.store.GetIter(tpath)
		if err != nil {
			return err
		}
		value, err := c.store.GetValue(iter, 0)
		if err != nil {
			return err
		}
		str, err := value.GetString()
		if err != nil {
			return err
		}
		if str == noProfile {
			err = c.profiles.file.UnSetDefault()
		} else {
			err = c.profiles.file.SetDefaultTo(str)
		}
		if err != nil {
			return err
		}
		c.profiles.curr = str
	}
	return nil
}

func showError(msg string) error {
	c := new(chooser)
	c.setupWindow()
	if c.err != nil {
		return fmt.Errorf("Unable to create window: %s", c.err)
	}
	btn, err := gtk.ButtonNewWithLabel(msg)
	if err != nil {
		return fmt.Errorf("Unable to create button: %s", err)
	}
	_, err = btn.Connect("clicked", func() { c.window.Destroy() })
	if err != nil {
		return fmt.Errorf("Unable to connect button: %s", err)
	}
	c.window.Add(btn)
	c.window.ShowAll()
	if msg == "this is a test" {
		return nil
	}
	gtk.Main()
	return nil
}

func (c *chooser) setupWindow() {
	if c.err != nil {
		return
	}
	if c.window, c.err = gtk.WindowNew(gtk.WINDOW_POPUP); c.err != nil {
		return
	}
	c.window.SetTitle(appName)
	if _, c.err = c.window.Connect("destroy", gtk.MainQuit); c.err != nil {
		return
	}
	c.window.SetPosition(gtk.WIN_POS_MOUSE)
	// need this workaround to support also single click to close the Popup;
	// lost Focus or button-release-event not available for TreeViewNewSelection
	if !permanent {
		_, c.err = c.window.ConnectAfter("button-release-event", func() {
			fmt.Println(c.profiles.curr)
			c.window.Destroy()
		})
	}
}

func (c *chooser) setupRootBox() {
	if c.err != nil {
		return
	}
	c.box, c.err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
}

func (c *chooser) setupSelection() {
	if c.err != nil {
		return
	}
	if c.selection, c.err = c.view.GetSelection(); c.err != nil {
		return
	}
	c.selection.SetMode(gtk.SELECTION_SINGLE)
	path, err := gtk.TreePathNewFromString(fmt.Sprintf("%d", c.profiles.currIdx))
	if err != nil {
		c.err = err
		return
	}
	c.selection.SelectPath(path)
	c.view.RowActivated(path, c.view.GetColumn(0))
	_, c.err = c.selection.Connect("changed", func() error { return c.selectionChanged() })
}

func (c *chooser) setupTreeView() {
	if c.err != nil {
		return
	}
	if c.view, c.err = gtk.TreeViewNewWithModel(c.store); c.err != nil {
		return
	}
	r, err := gtk.CellRendererTextNew()
	if err != nil {
		c.err = err
		return
	}
	column, err := gtk.TreeViewColumnNewWithAttribute(columnTitle, r, "text", 0)
	if err != nil {
		c.err = err
		return
	}
	c.view.AppendColumn(column)
}

func (c *chooser) setupListStore() {
	if c.err != nil {
		return
	}
	if c.store, c.err = gtk.ListStoreNew(glib.TYPE_STRING); c.err != nil {
		return
	}
	for _, i := range c.profiles.list {
		if c.err = c.store.SetValue(c.store.Append(), 0, i); c.err != nil {
			return
		}
	}
}

func initializeChooser(p *profiles) (c *chooser, err error) {
	c = &chooser{profiles: p}
	c.setupListStore()
	c.setupTreeView()
	c.setupSelection()
	c.setupRootBox()
	c.setupWindow()
	if c.err != nil {
		return nil, c.err
	}
	c.box.PackStart(c.view, true, true, 0)
	c.window.Add(c.box)
	return c, nil
}

func fetchProfiles() (p *profiles, err error) {
	p = new(profiles)
	if p.file, err = awsdefault.GetCredentialsFile(); err != nil {
		return
	}
	p.list = append(p.file.GetProfilesNames(), noProfile)
	p.curr, p.currIdx, err = p.file.GetUsedProfileNameAndIndex()
	if err != nil || p.currIdx == -2 { // -2 means no default set
		p.curr = noProfile
		p.currIdx = len(p.list) - 1 // last item is noProfile
	}
	return p, nil
}

func init() {
	flag.BoolVar(&permanent, "permanent", false, "the popup will not be closed after you clicked on a profile")
	flag.Parse()
}

func main() {
	gtk.Init(&os.Args)
	p, err := fetchProfiles()
	if err != nil { // only profile related errors
		if e := showError(err.Error()); e != nil {
			log.Println(e)
		}
		log.Fatalln(err)
	}

	c, err := initializeChooser(p)
	if err != nil {
		log.Fatalln(err)
	}
	c.window.ShowAll()
	gtk.Main()
}
