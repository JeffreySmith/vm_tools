package vmtools_test

import (
	"testing"

	"github.com/JeffreySmith/vmtools"
	"github.com/google/go-cmp/cmp"
)

func TestCreateBasicVM(t *testing.T) {
	t.Parallel()

	got,err := vmtools.CreateCluster(
		"jenkins",
		"jenkins cluster",
		"16GB",
		"rocky8",
		"TEAMNAME",
		"fake@email.com",
		"100GB",
		4,
	)
	if err != nil {
		t.Error(err)
	}

	want := vmtools.Cluster{
		Name:        "jenkins",
		Description: "jenkins cluster",
		RAM:         "16GB",
		OS:          "rocky8",
		Team:        "TEAMNAME",
		Email:       "fake@email.com",
		DiskSize:    map[string]string{"disk1": "100GB"},
		VCPUs:       4,
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

func TestOSinUppercase(t *testing.T){
		t.Parallel()

	got,err := vmtools.CreateCluster(
		"jenkins",
		"jenkins cluster",
		"16GB",
		"ROCKY8",
		"TEAMNAME",
		"fake@email.com",
		"100GB",
		4,
	)
	if err != nil {
		t.Error(err)
	}

	want := vmtools.Cluster{
		Name:        "jenkins",
		Description: "jenkins cluster",
		RAM:         "16GB",
		OS:          "rocky8",
		Team:        "TEAMNAME",
		Email:       "fake@email.com",
		DiskSize:    map[string]string{"disk1": "100GB"},
		VCPUs:       4,
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

func TestConfigSetIndentLessThan2(t *testing.T) {
	t.Parallel()

	c := vmtools.NewClusterConfig(vmtools.WithClusterIndent(-2))
	if c.Indent != 2 {
		t.Errorf("Expected 2, got %v", c.Indent)
	}
}

func TestCreateVMWithLowercaseValues(t *testing.T){
	t.Parallel()

	
	got,err := vmtools.CreateCluster(
		"jenkins",
		"jenkins cluster",
		"16gb",
		"rocky8",
		"TEAMNAME",
		"fake@email.com",
		"100gb",
		4,
	)
	if err != nil {
		t.Error(err)
	}

	want := vmtools.Cluster{
		Name:        "jenkins",
		Description: "jenkins cluster",
		RAM:         "16GB",
		OS:          "rocky8",
		Team:        "TEAMNAME",
		Email:       "fake@email.com",
		DiskSize:    map[string]string{"disk1": "100GB"},
		VCPUs:       4,
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

func TestVMWithInvalidName(t *testing.T){
	t.Parallel()
	
	_,err := vmtools.CreateCluster(
		"jenkins#$\\",
		"jenkins cluster",
		"16gb",
		"rocky8",
		"TEAMNAME",
		"fake@email.com",
		"100gb",
		4,
	)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestVmWithInvalidOS(t *testing.T){
	t.Parallel()
	_,err := vmtools.CreateCluster(
		"jenkins",
		"jenkins cluster",
		"16gb",
		"fake_linux_distro",
		"TEAMNAME",
		"fake@email.com",
		"100gb",
		4,
	)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestAddClusterToOrdMap(t *testing.T){
	t.Parallel()
	c := vmtools.NewClusterConfig()
	input := vmtools.Cluster{
		Name:        "jenkins",
		Description: "jenkins cluster",
		RAM:         "16GB",
		OS:          "rocky8",
		Team:        "TEAMNAME",
		Email:       "fake@email.com",
		DiskSize:    map[string]string{"disk1": "100GB"},
		VCPUs:       4,
	}
	vm, err := c.AddVM(input)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(vm, input){
		t.Error(cmp.Diff(vm, input))
	}
	
}

func TestAddClusterWithExistingName(t *testing.T) {
	t.Parallel()
	c := vmtools.NewClusterConfig()
	input := vmtools.Cluster{
		Name:        "jenkins",
		Description: "jenkins cluster",
		RAM:         "16GB",
		OS:          "rocky8",
		Team:        "TEAMNAME",
		Email:       "fake@email.com",
		DiskSize:    map[string]string{"disk1": "100GB"},
		VCPUs:       4,
	}
	_, err := c.AddVM(input)
	if err != nil {
		t.Error(err)
	}
	_,err = c.AddVM(input)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
