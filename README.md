# vm-topology-builder

## Description 
Topology builder to build random lab topologies for multiple juniper vNFs. The topology is defined in a yaml file which is fed as input. Once executed, the respective VMs/VNFs are spun up in the same compute.
The VMs/VNFs would be virtio based as this is intended to spin up topologies for functional testing rapidly based on the definition.

The program can also be used to generate only the xml files which could be used later by using virsh define and virsh start commands

Directory structure:
- templates: templates for each vnf based on dir. genxml will store the files in respective directory
- build: creates a directory based on topology names with respective images copied and is referenced in xml
- images: original image
- topologies: store all topologies yaml files which can be referenced
- vm-topo.go : main file

### Sample topology file 
```
---
network_nodes:
  - 
    name: "vqfx1"
    reImagePath: "/home/aprabh/vmx"
    pfeImagePath: "/home/aprabh/vmx/pfe"
    numIntf: 5
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
    numIntf: 5
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
