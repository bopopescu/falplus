// +build darwin freebsd linux solaris

package console

import (
	"os"

	"golang.org/x/sys/unix"
)

// NewPty creates a new pty pair
// The main is returned as the first console and a string
// with the path to the pty subordinate is returned as the second
func NewPty() (Console, string, error) {
	f, err := os.OpenFile("/dev/ptmx", unix.O_RDWR|unix.O_NOCTTY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, "", err
	}
	subordinate, err := ptsname(f)
	if err != nil {
		return nil, "", err
	}
	if err := unlockpt(f); err != nil {
		return nil, "", err
	}
	m, err := newMain(f)
	if err != nil {
		return nil, "", err
	}
	return m, subordinate, nil
}

type main struct {
	f        *os.File
	original *unix.Termios
}

func (m *main) Read(b []byte) (int, error) {
	return m.f.Read(b)
}

func (m *main) Write(b []byte) (int, error) {
	return m.f.Write(b)
}

func (m *main) Close() error {
	return m.f.Close()
}

func (m *main) Resize(ws WinSize) error {
	return tcswinsz(m.f.Fd(), ws)
}

func (m *main) ResizeFrom(c Console) error {
	ws, err := c.Size()
	if err != nil {
		return err
	}
	return m.Resize(ws)
}

func (m *main) Reset() error {
	if m.original == nil {
		return nil
	}
	return tcset(m.f.Fd(), m.original)
}

func (m *main) getCurrent() (unix.Termios, error) {
	var termios unix.Termios
	if err := tcget(m.f.Fd(), &termios); err != nil {
		return unix.Termios{}, err
	}
	return termios, nil
}

func (m *main) SetRaw() error {
	rawState, err := m.getCurrent()
	if err != nil {
		return err
	}
	rawState = cfmakeraw(rawState)
	rawState.Oflag = rawState.Oflag | unix.OPOST
	return tcset(m.f.Fd(), &rawState)
}

func (m *main) DisableEcho() error {
	rawState, err := m.getCurrent()
	if err != nil {
		return err
	}
	rawState.Lflag = rawState.Lflag &^ unix.ECHO
	return tcset(m.f.Fd(), &rawState)
}

func (m *main) Size() (WinSize, error) {
	return tcgwinsz(m.f.Fd())
}

func (m *main) Fd() uintptr {
	return m.f.Fd()
}

func (m *main) Name() string {
	return m.f.Name()
}

// checkConsole checks if the provided file is a console
func checkConsole(f *os.File) error {
	var termios unix.Termios
	if tcget(f.Fd(), &termios) != nil {
		return ErrNotAConsole
	}
	return nil
}

func newMain(f *os.File) (Console, error) {
	m := &main{
		f: f,
	}
	t, err := m.getCurrent()
	if err != nil {
		return nil, err
	}
	m.original = &t
	return m, nil
}

// ClearONLCR sets the necessary tty_ioctl(4)s to ensure that a pty pair
// created by us acts normally. In particular, a not-very-well-known default of
// Linux unix98 ptys is that they have +onlcr by default. While this isn't a
// problem for terminal emulators, because we relay data from the terminal we
// also relay that funky line discipline.
func ClearONLCR(fd uintptr) error {
	return setONLCR(fd, false)
}

// SetONLCR sets the necessary tty_ioctl(4)s to ensure that a pty pair
// created by us acts as intended for a terminal emulator.
func SetONLCR(fd uintptr) error {
	return setONLCR(fd, true)
}
