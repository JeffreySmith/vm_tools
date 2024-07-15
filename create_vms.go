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
	"fmt"
	"io"
	"os"

	"github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"

)

type ClusterConfig struct {
	Input io.Reader
	Output io.Writer
	indent int
	Vms VmDetails
	YamlString string
}

type Cluster struct {
	//For internal use with the VmDetails Map
	Name string `yaml:"-"`
	Description string `yaml:"vm_description"`
	VCPUs int `yaml:"vm_vcpus"`
	RAM string `yaml:"vm_ram"`
	OS string `yaml:"vm_os"`
	DiskSize map[string]string `yaml:"vm_disk_size"`

	Team string `yaml:"vm_request_by_team"`
	Email string `yaml:"vm_requested_by_email"`
}
type ClusterOpt func(*ClusterConfig)

type VmDetails struct {
	VirtualMachines *orderedmap.OrderedMap[string, Cluster] `yaml:"vm_details"`
}

func NewCluster(opts ...ClusterOpt) *ClusterConfig {
	c := &ClusterConfig{
		Input: os.Stdin,
		Output: os.Stdout,
		indent: 2,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}


func Marshal() {
	var b bytes.Buffer
	
	vm := VmDetails{}
	vm.VirtualMachines = orderedmap.New[string, Cluster](2)
	c := Cluster{}
	c.Description = "jenkins cluster"
	c.Name = "jenkins"
	c.VCPUs = 4
	c.RAM = "16GB"
	c.OS = "rocky8"
	c.DiskSize = make(map[string]string)
	c.DiskSize["disk1"] = "200GB"
	c.Team = "TEAM"
	c.Email = "billybob@fake.com"


	c2 := Cluster{}
	c2.Name = "TestNode"
	c2.Description = "Test Hue"
	c2.VCPUs = 4
	c2.RAM = "16gb"
	c2.OS = "rocky8"
	c2.DiskSize = make(map[string]string)
	c2.DiskSize["disk1"] = "100GB"
	c2.Team = "TEAM2"
	c2.Email = "fakename@email.com"
	vm.VirtualMachines.Set(c2.Name, c2)
	vm.VirtualMachines.Set(c.Name, c)

	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	_ = encoder.Encode(&vm)
	fmt.Println(string(b.String()))

}
