package vmtools_test

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/JeffreySmith/vmtools"
	"github.com/google/go-cmp/cmp"
)

func TestIndent(t *testing.T) {
	t.Parallel()
	conf := vmtools.NewConfig(vmtools.SetIndent(4))
	want := 4
	got := conf.GetIndent()
	if want != got {
		t.Errorf("Got %v, want %v", got, want)
	}
}
func TestCreateUser(t *testing.T) {
	t.Parallel()
	got := vmtools.CreateUser("alicebrown", "10.90.9.9")
	want := vmtools.User{Username: "alicebrown", Ip: "10.90.9.9"}
	if want != got {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestCreateMultipleUsersSingleIP(t *testing.T) {
	t.Parallel()
	buf := strings.NewReader("bobby\nzoe")
	config := vmtools.NewConfig(vmtools.WithInput(buf))
	ips := []string{"10.90.9.9"}
	config.GetUsers(ips)
	got := config.Users
	want := []vmtools.User{
		{Username: "bobby", Ip: "10.90.9.9"},
		{Username: "zoe", Ip: "10.90.9.9"},
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

func TestCreateMultipleUsersMultipleIPs(t *testing.T) {
	t.Parallel()
	buf := strings.NewReader("bobby\nzoe")
	config := vmtools.NewConfig(vmtools.WithInput(buf))
	ips := []string{"10.90.9.9", "192.168.1.4"}
	config.GetUsers(ips)
	got := config.Users
	want := []vmtools.User{
		{Username: "bobby", Ip: "10.90.9.9"},
		{Username: "zoe", Ip: "10.90.9.9"},
		{Username: "bobby", Ip: "192.168.1.4"},
		{Username: "zoe", Ip: "192.168.1.4"},
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

func TestCreateUsersFromFile(t *testing.T) {
	t.Parallel()
	buf, err := os.Open("testdata/users")
	defer buf.Close()
	if err != nil {
		t.Fatal(err)
	}
	config := vmtools.NewConfig(vmtools.WithInput(buf))
	ips := []string{"10.90.9.9"}
	config.GetUsers(ips)
	got := config.Users
	want := []vmtools.User{
		{Username: "bobby", Ip: "10.90.9.9"},
		{Username: "zoe", Ip: "10.90.9.9"},
		{Username: "alice", Ip: "10.90.9.9"},
		{Username: "johndoe", Ip: "10.90.9.9"},
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

func TestYamlOutput(t *testing.T) {
	t.Parallel()
	users := strings.NewReader("johndoe\nalice")

	config := vmtools.NewConfig(vmtools.WithInput(users))
	ips := []string{"10.90.9.9"}
	config.GetUsers(ips)
	yaml_string, err := config.GenerateYaml()
	if err != nil {
		t.Fatal(err)
	}
	want := `additional_users:
  - username: johndoe
    vm_ip: 10.90.9.9
  - username: alice
    vm_ip: 10.90.9.9`
	want += "\n"

	got := yaml_string

	if got != want {
		t.Errorf("\nGot:\n%v\nWant:\n%v", got, want)
	}
}
func TestWriteYaml(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	input := strings.NewReader("jupiter")
	output := bufio.NewWriter(&b)
	config := vmtools.NewConfig(vmtools.WithInput(input), vmtools.WithOutput(output))
	ips := []string{"10.90.9.9"}
	config.GetUsers(ips)
	_, err := config.GenerateYaml()
	if err != nil {
		t.Fatal(err)
	}

	want := `additional_users:
  - username: jupiter
    vm_ip: 10.90.9.9`
	want += "\n"

	err = config.WriteYaml()
	output.Flush()
	got := b.String()
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("Got:\n%v, want:\n%v", got, want)
	}
}
func TestWriteYamlWithHeader(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	input := strings.NewReader("jupiter")
	output := bufio.NewWriter(&b)
	header := "#This is my header.\n#There are many like it, but this one is mine."
	config := vmtools.NewConfig(vmtools.WithInput(input), vmtools.WithOutput(output), vmtools.WithHeader(header))
	ips := []string{"10.90.9.9"}
	config.GetUsers(ips)
	_, err := config.GenerateYaml()
	if err != nil {
		t.Fatal(err)
	}

	want := `#This is my header.
#There are many like it, but this one is mine.
additional_users:
  - username: jupiter
    vm_ip: 10.90.9.9`
	want += "\n"

	err = config.WriteYaml()
	output.Flush()
	got := b.String()
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("Got:\n%v, want:\n%v", got, want)
	}
}
