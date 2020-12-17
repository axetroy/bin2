package bin2

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func Test_getLatestRelease(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "dvm",
			args:    args{owner: "axetroy", repo: "dvm"},
			want:    "v1.3.1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLatestRelease(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getLatestRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	type args struct {
		owner      string
		repo       string
		version    string
		binaryName string
		userAgent  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Script
		wantErr bool
		goos    string
	}{
		{
			name: "dvm",
			args: args{
				owner:      "axetroy",
				repo:       "dvm",
				version:    "v1.3.0",
				binaryName: "dvm",
				userAgent:  "curl/",
			},
			want: &Script{
				FileName: "install.sh",
				ExtName:  ".sh",
				Content:  ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.owner, tt.args.repo, tt.args.version, tt.args.binaryName, tt.args.userAgent)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			bashPath, err := exec.LookPath("bash")

			if err != nil {
				t.Error(err)
				return
			}

			ps := exec.Cmd{
				Path: bashPath,
				Args: []string{},
			}

			shell := bytes.NewBufferString(got.Content)

			ps.Stdin = shell
			ps.Stdout = os.Stdout
			ps.Stderr = os.Stderr

			if err := ps.Run(); err != nil {
				t.Errorf("Run() error = %v", err)
				return
			}

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Generate() = %v, want %v", got, tt.want)
			// }
		})
	}
}
