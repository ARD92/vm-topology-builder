/*
__Author__ : Aravind Prabhakar
native go implementation of vNF topology builder. vMX/vSRX/vQFX can be used.
This is for functional testing purposes. The topology will be built on the same compute only.

The below may be required. if not, packages arent downloaded into go container.
export GO111MODULE=off

+++++++ Usage ++++++++
./vm-topo-builder -action <create/delete> -t <topology.yml>
*/

package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "os"
    //"strconv"
    //"encoding/xml"
    "gopkg.in/yaml.v3"
    "github.com/vishvananda/netlink"
    //"libvirt.org/go/libvirtxml"
    "libvirt.org/libvirt-go-xml"
    //"./templates/vqfx"
    //"./templates/vsrx"
    //"./templates/vmx"
)

/*
Struct to parse input yaml file
*/
type Link struct {
    Intf string `yaml:"intf"`
    PeerIntf string `yaml:"peerintf"`
    Name string `yaml:"name"`
}

type Node struct {
    Name string `yaml:"name"`
    Re_ImagePath string `yaml:"reImagePath"`
    Pfe_ImagePath string `yaml:"pfeImagePath"`
    NumIntf int `yaml:"numIntf"`
    VnfType string `yaml:"vnfType"`
    Re_memory uint `yaml:"re_memory"`
    Pfe_memory uint `yaml:"pfe_memory"`
    Re_Cores int `yaml:re_cores"`
    Pfe_Cores int `yaml:pfe_cores"`
    Links []Link `yaml:"link"`
}

type Nodes struct  {
    Network_nodes [] Node `yaml:"network_nodes"`
}

/* Read template */
func Readtemplate(file string) string {
    xmlfile, err := os.Open(file)
    if err!= nil {
        fmt.Println(err)
    }
    defer xmlfile.Close()
    val,_ := ioutil.ReadAll(xmlfile)
    return string(val)
}


/* Func create brdges 
- Create bridges based on link names 
- Place the interfaces into the bridges. i.e. while generating XML templates 
  load these as target bridges 
*/

func CreateBridges(node Nodes) {
    for i:=0; i<len(node.Network_nodes); i++ {
        for j:=0; j<len(node.Network_nodes[i].Links); j++ {
            //fmt.Printf("+v\n", node.Network_nodes[i].Links)
            linkattr := netlink.NewLinkAttrs() 
            linkattr.Name = node.Network_nodes[i].Links[j].Name
            linkattr.MTU = 9000 
            bridge := &netlink.Bridge{LinkAttrs: linkattr}
            err := netlink.LinkAdd(bridge)
            if err!=nil {
                fmt.Printf("Link %s not added: %v\n", linkattr.Name, err)
                continue
            }
            linkname,_ := netlink.LinkByName(linkattr.Name)

            if err := netlink.LinkSetUp(linkname); err !=nil {
                fmt.Printf("link not up : %s", err)
            }
        }
    }
}

/* Function to delete bridges when deleting the topology */
func DeleteBridges(node Nodes) {
    for i:=0; i<len(node.Network_nodes); i++ {
        for j:=0; j<len(node.Network_nodes[i].Links); j++ {
            linkattr := netlink.NewLinkAttrs()
            linkattr.Name = node.Network_nodes[i].Links[j].Name
            linkname,err := netlink.LinkByName(linkattr.Name)
            if err !=nil {
                fmt.Printf("Link not found")
                continue
            }
            if err := netlink.LinkSetDown(linkname); err !=nil {
                fmt.Printf("Link not present : %s", err)
            }
            bridge := &netlink.Bridge{LinkAttrs: linkattr}
            erro := netlink.LinkDel(bridge)
            if erro !=nil {
                fmt.Print("Bridge %s not deleted: %v\n", linkattr.Name, erro)
            }
        }
    }
}


/* Func start and stop Domains
- start domain (action create)
- stop domain (action delete)

/* Function to create new device disk to append into array devices list */
func NewDisk(node Node) libvirtxml.DomainDisk {
    var val uint = 0
    return libvirtxml.DomainDisk {
        Device:"disk",
        Driver: &libvirtxml.DomainDiskDriver {Name: "qemu", Type: "qcow2"},
        Source: &libvirtxml.DomainDiskSource {
            File: &libvirtxml.DomainDiskSourceFile {
                File: node.Re_ImagePath}},
        Target: &libvirtxml.DomainDiskTarget {
            Dev: "hda",
            Bus: "ide" },
        Alias: &libvirtxml.DomainAlias {
            Name: "ide-0-0"},
        Address: &libvirtxml.DomainAddress {
            Drive: &libvirtxml.DomainAddressDrive {
                Controller: &val ,
                Bus: &val,
                Target: &val,
                Unit: &val,
            },
        },
    }
}

/* Function to create new USB type controller to append into array devices 
Since there is no way to pass null uint , defined the first arg which need to 
be null as 999. Once this is encounterd, the below template kicks in 
which has the Address parameters removed
*/
func NewController(ctype string, index uint, model string,
                     alias string, pcidomain uint, pcibus uint,
                     pcislot uint, pcifunc uint, pcimultifunc string,
                    ) libvirtxml.DomainController {
    if pcidomain == 999 {
        return libvirtxml.DomainController {
            Type: ctype,
            Index: &index, 
            Model: model,
            Alias: &libvirtxml.DomainAlias {Name: alias},
        }
    } else {
        return libvirtxml.DomainController {
            Type: ctype,
            Index: &index, 
            Model: model,
            Alias: &libvirtxml.DomainAlias {Name: alias},
            Address: &libvirtxml.DomainAddress {
                PCI: &libvirtxml.DomainAddressPCI {
                    Domain: &pcidomain,
                    Bus: &pcibus,
                    Slot: &pcislot,
                    Function: &pcifunc,
                    MultiFunction: pcimultifunc,
                },
            },
        }
    }
}

/* Function to create new Serial type device */
func NewSerial(port *uint, portnum string, name string)libvirtxml.DomainConsole {
    return libvirtxml.DomainConsole {
         Source: &libvirtxml.DomainChardevSource {
                     TCP: &libvirtxml.DomainChardevSourceTCP {
                         Mode: "bind",
                         Host: "127.0.0.1",
                         Service: portnum,
                         TLS: "no",
                     },
         },
         Protocol: &libvirtxml.DomainChardevProtocol {
                     Type: "telnet",
         },
         Target: &libvirtxml.DomainConsoleTarget {
                 Type: "serial",
                 Port: port,
         },
         Alias: &libvirtxml.DomainAlias {
             Name: name,
         },
    }
}

/* Function to create new Input such as PS2 */
func NewInputs(itype string, bus string, name string)libvirtxml.DomainInput {
    return libvirtxml.DomainInput {
        Type: itype,
        Bus: bus,
        Alias: &libvirtxml.DomainAlias {
            Name: name,
        },
    }
}

/* Function for listeners. This has to be appended into Func Graphics() */
func GraphicListeners(address string)libvirtxml.DomainGraphicListener {
    return libvirtxml.DomainGraphicListener {
        Address: &libvirtxml.DomainGraphicListenerAddress {
            Address: address,
        },
    }
}

/* Function for Video */
func Videos(ram uint, vram uint, vgamem uint, alias string,
            pcidomain *uint, pcibus *uint, pcislot *uint,
            pcifunction *uint) libvirtxml.DomainVideo {
    return libvirtxml.DomainVideo {
         Model: libvirtxml.DomainVideoModel {
             Type: "qxl",
             Ram: ram,
             VRam: vram,
             VGAMem: vgamem,
             Heads: 1,
         },
         Alias: &libvirtxml.DomainAlias {Name: alias},
         Address: &libvirtxml.DomainAddress {
             PCI: &libvirtxml.DomainAddressPCI {
                 Domain: pcidomain,
                 Bus: pcibus,
                 Slot: pcislot,
                 Function: pcifunction,
             },
         },
    }
}

/* Function to enable graphics for the VM */
func Graphics()libvirtxml.DomainGraphic {
    return libvirtxml.DomainGraphic {
        VNC: &libvirtxml.DomainGraphicVNC {
            AutoPort: "yes",
            Listen: "127.0.0.1",
            //Listeners is a [] and needs to appended during initialization
        },
    }
}

/* Function to create new interfaces based on link mentioned in topology.yaml */
func NewNetworkIntf(mac string, link string , intf string,
                     model string)libvirtxml.DomainInterface {
    return libvirtxml.DomainInterface {
        MAC: &libvirtxml.DomainInterfaceMAC { Address: mac },
        Source: &libvirtxml.DomainInterfaceSource {
            Network: &libvirtxml.DomainInterfaceSourceNetwork {
               Bridge: link,
            },
        },
        Target: &libvirtxml.DomainInterfaceTarget {
            Dev: intf,
        },
        Model: &libvirtxml.DomainInterfaceModel {
            Type: model,
        },
    }
}

func TemplateVqfxRe(node Node) *libvirtxml.Domain {
    /* const nint is used to controller the display of certain xmls 
    which need to be declared as null*/
    const nint uint = 999
    ipci := uint(0x00)
    domcfg := &libvirtxml.Domain {
        Type: "kvm",
        Name: node.Name,
        Memory: &libvirtxml.DomainMemory {
            Value: node.Re_memory,
            Unit: "MB",
        },
        VCPU: &libvirtxml.DomainVCPU {
            Placement: "static",
            Value: 1 },
        OS: &libvirtxml.DomainOS {
            Type: &libvirtxml.DomainOSType {
                Arch: "x86_64",
                Machine: "pc-i440fx-xenial",
                Type: "hvm" }},
        /*Features: []libvirtxml.DomainFeatureList {
            APIC: &libvirtxml.DomainFeatureAPIC {
                EOI:""}
        },*/
        CPU: &libvirtxml.DomainCPU{
            Mode: "host-model",
            Topology: &libvirtxml.DomainCPUTopology {
                Sockets: 1,
                Cores: node.Re_Cores,
                Threads: 1}},
        OnPoweroff: "destroy",
        OnReboot: "restart",
        OnCrash: "restart",
        PM: &libvirtxml.DomainPM {
            SuspendToMem: &libvirtxml.DomainPMPolicy {Enabled: "no"},
            SuspendToDisk: &libvirtxml.DomainPMPolicy {Enabled: "no"}},
        Devices: &libvirtxml.DomainDeviceList {
            Emulator: "/usr/bin/qemu-system-x86_64",
            MemBalloon: &libvirtxml.DomainMemBalloon {
                Model: "virtio",
                Alias: &libvirtxml.DomainAlias {Name: "balloon"},
                Address: &libvirtxml.DomainAddress {
                    PCI: &libvirtxml.DomainAddressPCI {
                        Domain: &ipci,
                        Bus: &ipci,
                        Slot: &ipci,
                        Function: &ipci,
                    },
                },
            },
        },
    }
    disk := NewDisk(node)
    domcfg.Devices.Disks = append(domcfg.Devices.Disks, disk)

    /* Reference for function params 
    func NewController(ctype string, index uint, model string,
                     alias string, pcidomain uint, pcibus uint,
                     pcislot uint, pcifunc uint, pcimultifunc string,
                    ) libvirtxml.DomainController {
    */
    c1 := NewController("usb",uint(0),"ich9-ehci1","usb",uint(0x0000),uint(0x00),uint(0x0b),uint(0x7),"")
    c2 := NewController("usb",uint(0),"ich9-uhci1","usb",uint(0x0000),uint(0x00),uint(0x0b),uint(0x0),"on")
    c3 := NewController("usb",uint(0),"ich9-ehci2","usb",uint(0x0000),uint(0x00),uint(0x0b),uint(0x0b),"")
    c4 := NewController("usb",uint(0),"ich9-ehci3","usb",uint(0x0000),uint(0x00),uint(0x0b),uint(0x0b),"")
    c5 := NewController("pci",uint(0),"pci-root","pci.0",nint,nint,nint,nint,"") 
    c6 := NewController("ide",uint(0),"","ide",uint(0x0000),uint(0x00),uint(0x01),uint(0x0),"")
    c7 := NewController("virtio-serial",uint(0),"","virtio-serial0",uint(0x0000),uint(0x00),uint(0x0a),uint(0x0),"")
    domcfg.Devices.Controllers = append(domcfg.Devices.Controllers, c1, c2, c3, c4, c5, c6, c7)
    for i:=0; i<len(node.Links); i++ {
       nintf := NewNetworkIntf("02:aa:01:10:00:01", node.Links[i].Name, node.Links[i].Intf, "e1000")
       domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, nintf)
    }
    return domcfg
}

/* Pass struct to function to handle topo spin up*/ 
func Genxml(node Nodes) {
    for i:=0; i<len(node.Network_nodes);i++ {
        if node.Network_nodes[i].VnfType == "vmx" {
            fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
            //Generate RE and PFE xml
        } else if node.Network_nodes[i].VnfType == "vqfx" {
            domcfg := TemplateVqfxRe(node.Network_nodes[i])
            rexmldoc, err := domcfg.Marshal()
            if err != nil {
                fmt.Printf("Marshing failed %v\n", err)
            }
            fmt.Printf("%v", rexmldoc)
            /*
            if os.Args[1] == "action" {
                //write into file
            }*/
        } else if node.Network_nodes[i].VnfType == "vsrx" {
            fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
        }
    }
}

/*
Main function definition
*/
func main() {
    if  len(os.Args) > 1 {
        var nodes Nodes
        if len(os.Args) > 2 {
            // Parse yaml generate xml based on topology templates
            rfile, err := ioutil.ReadFile(os.Args[2])
            if err != nil {
                log.Fatal(err)
            }
           // parse yaml into structs
            err = yaml.Unmarshal(rfile, &nodes)
            if  err != nil {
                log.Fatal(err)
            }
        } else {
            fmt.Printf("Missing argument Topology yaml file\n")
        }
        if os.Args[1] == "create" {
            fmt.Printf("Creating bridges\n")
            //Parse through yaml and create all bridges needed to define topology
            CreateBridges(nodes)
            //copy original image to build directory and refer that xml
            //Genxml(nodes)
        } else if os.Args[1] == "delete" {
            fmt.Printf("Deleting bridges\n")
            DeleteBridges(nodes)
        } else if os.Args[1] == "genxml" {
            fmt.Printf(" Generating xml files..  \n")
            Genxml(nodes)
        } else if os.Args[1] == "help" {
            fmt.Printf(` 
            -------------- VM Topology Builder V1.0 --------------------------------------------
            Topology builder for Juniper VNFs such as vMX, vQFX and vSRX 
            This is used for functional testing only and is virtio based. 
            The topology will be spun up on the same compute
            The topology is defined in a yaml file. Example topologies are placed
            under ./topologies path

            genxml: generates xml for the topology defined and stores it under templates directory
            with the topology name
            create: creates the topology
            delete: deletes the topology

            To use this below is the usage. 
            ./vm-topo create topology.yaml
            ./vm-topo delete topology.yaml
            ./vm-topo genxml topology.yaml
            ----------------------------------------------------------------------------------------`)
            fmt.Printf("\n")
        } else {
            fmt.Printf(" Unknown first argument. Please check \n")
        }
    } else {
        fmt.Printf("\n INVALID OPTION. FOLLOW THE BELOW USAGE\n")
        fmt.Printf("\n+++++++++++++++++ usage +++++++++++++++++++++ \n ")
        fmt.Printf(" ./vm-topo <create/delete> <path to topo.yml>\n")
        fmt.Print("+++++++++++++++++++++++++++++++++++++++++++++\n")
    }
}
