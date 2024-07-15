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
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"
)

type ClusterConfig struct {
	Input      io.Reader
	Output     io.Writer
	Indent     int
	Vms        VmDetails
	YamlString string
}

type Cluster struct {
	//For internal use with the VmDetails Ordered Map
	Name        string            `yaml:"-"`
	Description string            `yaml:"vm_description"`
	VCPUs       int               `yaml:"vm_vcpus"`
	RAM         string            `yaml:"vm_ram"`
	OS          string            `yaml:"vm_os"`
	DiskSize    map[string]string `yaml:"vm_disk_size"`

	Team  string `yaml:"vm_request_by_team"`
	Email string `yaml:"vm_requested_by_email"`
}

type option func(*ClusterConfig)

type VmDetails struct {
	VirtualMachines *orderedmap.OrderedMap[string, Cluster] `yaml:"vm_details"`
}

func getSupportedOS() []string {
	return []string{"centos7", "rocky8", "rocky9", "ubuntu20.04", "ubuntu22.04", "ubuntu24.04"}
}

func checkSupportedOS(os string) bool {
	for _, os_supported := range getSupportedOS() {
		if os == os_supported {
			return true
		}
	}
	return false
}

func NewClusterConfig(opts ...option) *ClusterConfig {
	c := &ClusterConfig{
		Input:  os.Stdin,
		Output: os.Stdout,
		Indent: 2,
		Vms:    VmDetails{VirtualMachines: orderedmap.New[string, Cluster]()},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithMapSize(n int) func(*ClusterConfig) {
	return func(c *ClusterConfig) {
		c.Vms.VirtualMachines = orderedmap.New[string, Cluster](n)
	}
}

func WithClusterIndent(n int) func(*ClusterConfig) {
	return func(c *ClusterConfig) {
		if n < 2 {
			n = 2
		}
		c.Indent = n
	}
}

func WithClusterInput(input io.Reader) func(*ClusterConfig) {
	return func(c *ClusterConfig) {
		c.Input = input
	}
}

func WithClusterOutput(output io.Writer) func(*ClusterConfig) {
	return func(c *ClusterConfig) {
		c.Output = output
	}
}

func CreateCluster(name, description, ram, os, team, email, disksize string, vcpu int) (Cluster, error) {
	name_regex, err := regexp.Compile("^[0-9a-zA-Z_]+$")
	if err != nil {
		return Cluster{}, err
	}
	if !name_regex.MatchString(name) {
		return Cluster{}, errors.New("Invalid cluster name. May only contain letters, numbers, and underscores")
	}
	disks := make(map[string]string, 1)
	disks["disk1"] = strings.ToUpper(disksize)
	ram = strings.ToUpper(ram)
	os = strings.ToLower(os)
	if !checkSupportedOS(os) {
		err := fmt.Sprintf("Cluster OS: '%v' invalid", os)
		return Cluster{}, errors.New(err)
	}

	c := Cluster{
		Name:        name,
		Description: description,
		RAM:         ram,
		OS:          os,
		Team:        team,
		Email:       email,
		DiskSize:    disks,
		VCPUs:       vcpu,
	}
	return c, nil
}

func (c *ClusterConfig) AddVM(vm Cluster) (Cluster, error){
	v, exists := c.Vms.VirtualMachines.Get(vm.Name)
	if exists {
		return v, errors.New(fmt.Sprintf("VM '%v' already exists",v.Name))
	}
	c.Vms.VirtualMachines.Set(vm.Name,vm)
	v,_ = c.Vms.VirtualMachines.Get(vm.Name)
	return v, nil
}

func Marshal() {
	var b bytes.Buffer

	vm := VmDetails{}
	vm.VirtualMachines = orderedmap.New[string, Cluster](2)
	c, _ := CreateCluster(
		"jenkins",
		"jenkins cluster",
		"16GB",
		"rocky8",
		"TEAMNAME",
		"fake@email.com",
		"100gb",
		4,
	)

	c2 := Cluster{}
	c2.Name = "jenkins"
	c2.Description = "Test Hue"
	c2.VCPUs = 4
	c2.RAM = "16gb"
	c2.OS = "rocky8"
	c2.DiskSize = make(map[string]string)
	c2.DiskSize["disk1"] = "100GB"
	c2.Team = "TEAM2"
	c2.Email = "fakename@email.com"
	vm.VirtualMachines.Set(c2.Name, c2)
	_, ok := vm.VirtualMachines.Get(c2.Name)
	if !ok {
		vm.VirtualMachines.Set(c.Name, c)
	} else {
		fmt.Println("Virtual machine already exists")
	}

	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	_ = encoder.Encode(&vm)
	fmt.Println(string(b.String()))

}
