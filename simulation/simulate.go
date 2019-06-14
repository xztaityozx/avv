package simulation

//type HSPICE struct {
//	Path    string
//	Options string
//}
//
//type Task struct {
//	task.Task
//	HSPICE
//}
//
//// getCommand generate simulation command with hspice
//// returns:
////  - string: command string
//func (t Task) getCommand() string {
//	return fmt.Sprintf("cd %s && %s %s -i %s -o ./hspice &> ./hspice.log",
//		t.Files.Directories.DstDir,
//		t.HSPICE.Path, t.HSPICE.Options,
//		t.Files.SPIScript)
//}
//
//// Invoke start simulation with context
//func (t Task) Invoke(ctx context.Context) error {
//
//	ch := make(chan error)
//
//	go func() {
//		_, err := exec.Command("bash", "-c", t.getCommand()).Output()
//		ch <- err
//	}()
//
//	select {
//	case <-ctx.Done():
//		return nil
//	case err := <-ch:
//		return err
//
//	}
//}
