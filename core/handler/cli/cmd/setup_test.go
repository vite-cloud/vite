package cmd

//func TestNewSetupCommand(t *testing.T) {
//	datadir.UseTestHome(t)
//
//	CommandTest{
//		NewCommand: NewSetupCommand,
//		Test: func(console *Expect) {
//			console.Expect("Select your provider:").Enter()
//			console.Expect("Select your protocol:").Enter()
//			/console.Expect("Enter your repository:").Write("vite/org\n")
//
//			console.EOF()
//		},
//	}.Run(t)
//
//	f, err := locator.ConfigStore.Open(locator.ConfigFile, os.O_RDONLY, 0)
//	assert.NilError(t, err)
//	defer f.Close()
//
//	marshaled, err := json.Marshal(locator.Locator{
//		Provider:   "github",
//		Protocol:   "ssh",
//		Repository: "foo/bar",
//		Path:       "/sub-path",
//		Branch:     "main",
//	})
//	assert.NilError(t, err)
//
//	contents, err := io.ReadAll(f)
//	assert.NilError(t, err)
//
//	assert.Equal(t, string(contents), string(marshaled))
//}
