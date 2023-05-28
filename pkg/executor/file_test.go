package executor

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"runtime"
	"testing"
)

func Test_findExecutable(t *testing.T) {
	type args struct {
		executableName string
	}
	type want struct {
		err error
	}
	dir := t.TempDir()
	script := []byte("echo 'hello'")
	var scriptFileName string

	if runtime.GOOS == "windows" {
		scriptFileName = "test.cmd"
	} else {
		scriptFileName = "test.sh"
	}
	testScriptPath := path.Join(dir, scriptFileName)
	err := os.WriteFile(testScriptPath, script, 0777)
	if err != nil {
		t.FailNow()
	}

	tests := []struct {
		name  string
		args  args
		want  want
		setup func(t *testing.T)
	}{
		{
			name: "missing executable",
			args: args{executableName: "__missing"},
			want: want{
				err: ErrExecutableNotFound,
			},
		},
		{
			name: "executable in filepath",
			args: args{executableName: testScriptPath},
			want: want{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}
			_, err := newFileExecutor(tt.args.executableName)
			require.ErrorIs(t, tt.want.err, err)
		})
	}
}
