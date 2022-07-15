# vm-topology-builder

## Description 
Topology builder to build random lab topologies for multiple juniper vNFs. The topology is defined in a yaml file which is fed as input. Once executed, the respective VMs/VNFs are spun up in the same compute.
The VMs/VNFs would be virtio based as this is intended to spin up topologies for functional testing rapidly based on the definition.

The program can also be used to generate only the xml files which could be used later by using virsh define and virsh start commands

### Packages required for compilation and functioning 
- go get "gopkg.in/yaml.v3"
- go get "github.com/vishvananda/netlink"
- go get "libvirt.org/libvirt-go-xml"
- go get "github.com/digitalocean/go-libvirt"

    Be sure to download libvirt-dev package , else it may not work as expected. 
    apt install -y libvirt-dev 

sometimess, the below may be required. if not, packages arent downloaded into go container.
export GO111MODULE=off

Note: In case libvirt-dev package is not installed, you may see multiple errors such as below when performing "go get" and you would not be able to compile.

```
If use native libvirt package you might see the below errors due to cbindings. 
Perhaps you should add the directory containing `libvirt.pc'
to the PKG_CONFIG_PATH environment variable
No package 'libvirt' found

...
```
Hence moved to using go-libvirt by digital ocean which handles better as a pure go implementation by making RPC calls to libvirt.

### Directory structure:
- templates: templates for each vnf based on dir. genxml will store the files in respective directory
- build: creates a directory based on topology names with respective images copied and is referenced in domain xml. Images from source path is copied to here.
- topologies: store all topologies yaml files which can be referenced. 
- topology-diagram: application to generate topology diagram based on the input topology file. 
- vm-topo.go : main file


```
├── README.md
├── build
├── templates
│   ├── ipfabric-test
│   │   ├── leaf1_vqfx-pfe.xml
│   │   ├── leaf1_vqfx-re.xml
│   │   ├── leaf2_vqfx-pfe.xml
│   │   ├── leaf2_vqfx-re.xml
│   │   ├── leaf3_vqfx-pfe.xml
│   │   ├── leaf3_vqfx-re.xml
│   │   ├── spine1_vqfx-pfe.xml
│   │   ├── spine1_vqfx-re.xml
│   │   ├── spine2_vqfx-pfe.xml
│   │   └── spine2_vqfx-re.xml
│   └── vsrx-test
│       ├── vsrx1_vsrx.xml
│       └── vsrx2_vsrx.xml
├── topologies
│   ├── ipfabric-test.yaml
│   ├── topology.yaml
│   └── vsrx-test.yaml
├── topology-diagram
│   ├── README.md
│   ├── drawthenet_boiler.yaml
│   └── generate_topology_diagram.py
└── vm-topo.go
```
Note: currently relative path doesnt work while executing the application. Place the topology file along with the binary when executing.

### Sample topology file 

* vQFX requires both RE and PFE params 
* vMX requires both RE and PFE params
* vSRX requires only RE params since it is a single image. The PFE params can be left blank.

```
---
network_nodes:
  - 
    name: "vqfx1"
    reImagePath: "/home/aprabh/vmx"
    pfeImagePath: "/home/aprabh/vmx/pfe"
    vnfType: "vqfx"
    re_memory: 1024
    pfe_memory: 2048
    re_cores: 1
    pfe_cores: 1
    re_vcpu: 1
    pfe_vcpu: 4
    re_port: 8610
    pfe_port: 8710
    link: 
      # name is user defined. This will be used to create the bridge name
      - name: "test1_test2"
        intf: "ge-0.0.0-vqfx1"
        peerintf: "ge-0.0.0-vqfx2"
      - name: "test1_test22"
        intf: "ge-0.0.1-vqfx1"
        peerintf: "ge-0.0.1-vqfx2"
  -
    name: "vqfx2"
    reImagePath: "/home/aprabh/vqfx"
    pfeImagePath: "/home/aprabh/vqfx"
    re_memory: 1024
    pfe_memory: 2048
    vnfType: "vqfx"
    re_cores: 1
    pfe_cores: 1
    re_vcpu: 1
    pfe_vcpu: 4
    re_port: 8611
    pfe_port: 8711
    link:
      - name: "test1_test2"
        intf: "ge-0.0.0-vqfx2"
        peerintf: "ge-0.0.0-vqfx1"
      - name: "test1_test22"
        intf: "ge-0.0.1-vqfx2"
        peerintf: "ge-0.0.1-vqfx1"
```
### Example:

#### Create the topology 
```
root@compute21:/home/aprabh/vm-topology-builder# ./vm-topo create ipfabric-test.yaml
Creating bridges
Entering vqfx creating bridges
Entering vqfx creating bridges
Bridge vqfx_res not added: file exists
Entering vqfx creating bridges
Bridge vqfx_res not added: file exists
Link spine1_leaf1 not added: file exists
Link spine2_leaf1 not added: file exists
Entering vqfx creating bridges
Bridge vqfx_res not added: file exists
Link spine1_leaf2 not added: file exists
Link spine2_leaf2 not added: file exists
Entering vqfx creating bridges
Bridge vqfx_res not added: file exists
Link spine1_leaf3 not added: file exists
Link spine2_leaf3 not added: file exists
Copying image from source path to build directory. Image will use the build dir path as source...
Sleeping for 10seconds to ensure file copies successfully
Sleeping for 10seconds to ensure file copies successfully
Sleeping for 10seconds to ensure file copies successfully
Sleeping for 10seconds to ensure file copies successfully
Sleeping for 10seconds to ensure file copies successfully
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
Defining and Starting the Domains...
```
#### Verify
```
root@compute21:/home/aprabh/vm-topology-builder# virsh list
 Id    Name         State
-----------------------------
 125   spine1-re    running
 126   spine1-pfe   running
 127   spine2-re    running
 128   spine2-pfe   running
 129   leaf1-re     running
 130   leaf1-pfe    running
 131   leaf2-re     running
 132   leaf2-pfe    running
 133   leaf3-re     running
 134   leaf3-pfe    running


root@compute21:/home/aprabh/vm-topology-builder# brctl show
bridge name	bridge id		STP enabled	interfaces
leaf1_int		8000.feaa01210002	no		leaf1_pfe-int
							leaf1_re-int
leaf1_vm1		8000.feaa01200202	no		ge-0.0.2-leaf1
leaf1_vm2		8000.feaa01200203	no		ge-0.0.3-leaf1
leaf2_int		8000.feaa01210003	no		leaf2_pfe-int
							leaf2_re-int
leaf2_vm3		8000.feaa01200302	no		ge-0.0.2-leaf2
leaf2_vm4		8000.feaa01200303	no		ge-0.0.3-leaf2
leaf3_int		8000.feaa01210004	no		leaf3_pfe-int
							leaf3_re-int
leaf3_vm5		8000.feaa01200402	no		ge-0.0.2-leaf3
leaf3_vm6		8000.feaa01200403	no		ge-0.0.3-leaf3
mgmt_ext		8000.feaa01220000	no		leaf1_pfe-ext
							leaf1_re-ext
							leaf2_pfe-ext
							leaf2_re-ext
							leaf3_pfe-ext
							leaf3_re-ext
							spine1_pfe-ext
							spine1_re-ext
							spine2_pfe-ext
							spine2_re-ext
spine1_int		8000.feaa01210000	no		spine1_pfe-int
							spine1_re-int
spine1_leaf1		8000.feaa01200000	no		ge-0.0.0-leaf1
							ge-0.0.0-spine1
spine1_leaf2		8000.feaa01200001	no		ge-0.0.0-leaf2
							ge-0.0.2-spine1
spine1_leaf3		8000.feaa01200002	no		ge-0.0.0-leaf3
							ge-0.0.1-spine1
spine2_int		8000.feaa01210001	no		spine2_pfe-int
							spine2_re-int
spine2_leaf1		8000.feaa01200100	no		ge-0.0.1-leaf1
							ge-0.0.1-spine2
spine2_leaf2		8000.feaa01200101	no		ge-0.0.0-spine2
							ge-0.0.1-leaf2
spine2_leaf3		8000.feaa01200102	no		ge-0.0.1-leaf3
							ge-0.0.2-spine2
virbr0		8000.5254009fbc1f	yes		virbr0-nic
vqfx_res		8000.feaa01230000	no		leaf1_re-res
							leaf2_re-res
							leaf3_re-res
							spine1_re-res
							spine2_re-res



root@compute21:/home/aprabh/vm-topology-builder# telnet localhost 8110
Trying ::1...
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.


vqfx4-re (ttyd0)

login: root
Password:

The default user/password is root/Juniper
```

#### Delete topology
```
./vm-topo delete ipfabric-test.yaml
```

#### Generate XML template for the topology
```
./vm-topo genxml ipfabric-test.yaml
```
The XML files are stored under the directory templates with the topology file name as the subdirectory

#### Initial config on the VMs
```
python3 initial_config.py -t <topology.yaml> 
```
This will login to the VMs and configure the root authentication and enable SSH.

### Setup Instructions 
#### Installing QEMU-KVM dependencies 
The below were tried on Ubuntu 18.04 Bionic
```
sudo apt install qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils
```

#### Install Golang
##### Docker
```
docker run -itd golang:latest --name go -v ${PWD}:/usr/src/app bash 
```
using go get, install all necessary packages

##### Build steps
```
go build vm-topo.go
```
Run the application based on the above steps

