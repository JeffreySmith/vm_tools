/*BSD 3-Clause License

Copyright (c) 2024, Jeffrey Smith

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package vmtools

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/go-faster/yaml"
)

type Config struct {
	Input      io.Reader
	Output     io.Writer
	indent     int
	encoder    yaml.Encoder
	Users      []User
	Header     string
	YamlString string
}
type Opt func(*Config)

type User struct {
	Username string `yaml:"username"`
	Ip       string `yaml:"vm_ip"`
}

type AdditionalUsers struct {
	Users []User `yaml:"additional_users"`
}

func NewConfig(opts ...Opt) *Config {
	c := &Config{
		Input:  os.Stdin,
		Output: os.Stdout,
		indent: 2,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func SetIndent(indent int) func(*Config) {
	return func(c *Config) {
		c.indent = indent
	}
}

func WithInput(buf io.Reader) func(*Config) {
	return func(c *Config) {
		c.Input = buf
	}
}
func WithOutput(buf io.Writer) func(*Config) {
	return func(c *Config) {
		c.Output = buf
	}
}
func WithHeader(header string) func(*Config) {
	return func(c *Config) {
		c.Header = header
	}
}
func (c *Config) GetIndent() int {
	return c.indent
}

func CreateUser(username string, ip string) User {
	u := User{Username: username, Ip: ip}
	return u
}

func (c *Config) GetUsers(ips []string) {
	var usernames []string
	var users []User
	scanner := bufio.NewScanner(c.Input)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			usernames = append(usernames, text)
		}
	}
	user_length := len(usernames)
	users = make([]User, user_length*len(ips))

	for i := 0; i < len(users); i++ {
		username := usernames[i%user_length]
		ip := ips[i/user_length]
		users[i] = CreateUser(username, ip)
	}

	c.Users = users
}
func (c *Config) GenerateYaml() (string, error) {
	var b bytes.Buffer
	additionalusers := AdditionalUsers{Users: c.Users}

	if len(additionalusers.Users) == 0 {
		return "", errors.New("No users detected, empty output")
	}

	encoder := yaml.NewEncoder(&b)
	defer encoder.Close()
	encoder.SetIndent(c.indent)
	err := encoder.Encode(&additionalusers)
	if err != nil {
		return "", err
	}
	c.YamlString = b.String()

	return c.YamlString, nil
}

func (c *Config) WriteYaml() error {
	if len(c.YamlString) == 0 {
		return errors.New("Uninitialized yaml string")
	}
	if len(c.Header) > 0 {
		c.Output.Write([]byte(c.Header))
		c.Output.Write([]byte{'\n'})
	}
	c.Output.Write([]byte(c.YamlString))
	return nil
}
