package vmtools_test

import (
	"bufio"
	"bytes"
	"fmt"
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
	got, err := vmtools.CreateUser("alicebrown", "10.90.9.9")
	want := vmtools.User{Username: "alicebrown", Ip: "10.90.9.9"}
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestCreateMultipleUsersSingleIP(t *testing.T) {
	t.Parallel()
	buf := strings.NewReader("bobby\nzoe")
	config := vmtools.NewConfig(vmtools.WithInput(buf))
	ips := []string{"10.90.9.9"}
	config.CreateUsers(ips)
	got := config.Users
	want := []vmtools.User{
		{Username: "bobby", Ip: "10.90.9.9"},
		{Username: "zoe", Ip: "10.90.9.9"},
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}
func TestCreateUserWithUpperCase(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		username, ip string
		want         vmtools.User
	}{
		{
			username: "JohnDoe",
			ip:       "10.90.9.140",
			want:     vmtools.User{Username: "johndoe", Ip: "10.90.9.140"},
		},
		{
			username: "Bobby",
			ip:       "10.90.9.140",
			want:     vmtools.User{Username: "bobby", Ip: "10.90.9.140"},
		},
	}
	for _, tc := range tcs {
		name := fmt.Sprintf("%s to lowercase = %s", tc.username, tc.want.Username)
		t.Run(name, func(t *testing.T) {
			got, err := vmtools.CreateUser(tc.username, tc.ip)
			if err != nil {
				t.Errorf("Create user %v failed. Error: %v", tc.username, err)
			}
			if got != tc.want {
				t.Errorf("Got %v, want %v", got, tc.want)
			}
		})
	}
}
func TestErrorOnUserWithNonAlphaChars(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		username, ip string
		want         error
	}{
		{
			username: "JohnDoe321",
			ip:       "10.90.9.140",
		},
		{
			username: "#^&Bobby",
			ip:       "10.90.9.140",
		},
	}
	for _, tc := range tcs {
		name := fmt.Sprintf("%s should return error", tc.username)
		t.Run(name, func(t *testing.T) {
			_, err := vmtools.CreateUser(tc.username, tc.ip)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestCreateMultipleUsersMultipleIPs(t *testing.T) {
	t.Parallel()
	buf := strings.NewReader("bobby\nzoe")
	config := vmtools.NewConfig(vmtools.WithInput(buf))
	ips := []string{"10.90.9.9", "192.168.1.4"}
	err := config.CreateUsers(ips)
	if err != nil {
		t.Fatal(err)
	}
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
	err = config.CreateUsers(ips)
	if err != nil {
		t.Fatal(err)
	}
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
	err := config.CreateUsers(ips)
	if err != nil {
		t.Fatal(err)
	}
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
	err := config.CreateUsers(ips)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.GenerateYaml()
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
	err := config.CreateUsers(ips)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.GenerateYaml()
	if err != nil {
		t.Fatal(err)
	}

	want := `#This is my header.
#There are many like it, but this one is mine.
additional_users:
  - username: jupiter
    vm_ip: 10.90.9.9`
	//The output contains an extra newline characters, so to match
	//we need to add one here. Only a problem for the tests
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

func TestEmptyYaml(t *testing.T) {
	t.Parallel()

	config := vmtools.NewConfig()
	err := config.WriteYaml()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestEmptyLineNotIncluded(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	ip := []string{"10.90.9.9"}
	output := bufio.NewWriter(&b)
	input := strings.NewReader("bobby\n\nmillybrown")
	config := vmtools.NewConfig(vmtools.WithInput(input), vmtools.WithOutput(output))
	err := config.CreateUsers(ip)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.GenerateYaml()
	if err != nil {
		t.Error(err)
	}
	//The output contains an extra newline characters, so to match
	//we need to add one here. Only a problem for the tests
	want := `additional_users:
  - username: bobby
    vm_ip: 10.90.9.9
  - username: millybrown
    vm_ip: 10.90.9.9`
	want += "\n"

	config.WriteYaml()
	output.Flush()
	got := b.String()
	if got != want {
		t.Errorf("Got:\n%vWant:\n%v", got, want)
	}
}

func TestInputBufferEmpty(t *testing.T) {
	t.Parallel()
	ip := []string{"10.90.9.9"}
	input := strings.NewReader("")
	config := vmtools.NewConfig(vmtools.WithInput(input))

	err := config.CreateUsers(ip)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.GenerateYaml()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
