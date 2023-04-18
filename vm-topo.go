/*
__Author__ : Aravind Prabhakar
Description: native go implementation of vNF topology builder. vMX/vSRX/vQFX can be used.
This is for functional testing purposes. The topology will be built on the same compute only.

version:  2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/vishvananda/netlink"
	"gopkg.in/yaml.v3"
	libvirtxml "libvirt.org/libvirt-go-xml"
	//"github.com/libvirt/libvirt-go"
	//"./templates/vqfx"
	//"./templates/vsrx"
	//"./templates/vmx"
)

/* Global variable definition */
var BASE_VMX_LINK_MAC_ADDRESS string = "02:aa:01:10:00:00"
var BASE_VMX_INT_MAC_ADDRESS string = "02:aa:01:11:00:00"
var BASE_VMX_EXT_MAC_ADDRESS string = "02:aa:01:12:00:00"

var BASE_VQFX_LINK_MAC_ADDRESS string = "02:aa:01:20:00:00"
var BASE_VQFX_INT_MAC_ADDRESS string = "02:aa:01:21:00:00"
var BASE_VQFX_EXT_MAC_ADDRESS string = "02:aa:01:22:00:00"
var BASE_VQFX_RES_MAC_ADDRESS string = "02:aa:01:23:00:00"

var BASE_VSRX_LINK_MAC_ADDRESS string = "02:aa:01:30:00:00"
var BASE_VSRX_INT_MAC_ADDRESS string = "02:aa:01:30:00:00"
var BASE_VSRX_EXT_MAC_ADDRESS string = "02:aa:01:31:00:00"

var BASE_LINUX_LINK_MAC_ADDRESS string = "02:aa:01:40:00:00"
var BASE_LINUX_INT_MAC_ADDRESS string = "02:aa:01:40:00:00"
var BASE_LINUX_EXT_MAC_ADDRESS string = "02:aa:01:41:00:00"

const LISTEN_ADDRESS string = "127.0.0.1"
const BRIDGE_MTU int = 9000
const BRIDGE_MGMT string = "BR_MGMT"
const VQFX_BRIDGE_RES string = "vqfx_res"

/*
Struct to parse input yaml file
*/
type Link struct {
	Intf     string `yaml:"intf"`
	PeerIntf string `yaml:"peerintf"`
	Name     string `yaml:"name"`
}

type Node struct {
	Name           string `yaml:"name"`
	ImagePath      string `yaml:"imagePath"`
	Re_Image       string `yaml:"reImage"`
	Pfe_Image      string `yaml:"pfeImage"`
	Metadata_Image string `yaml:"metadataImage"`
	Vhdd_Image     string `yaml:"vhddImage"`
	VnfType        string `yaml:"vnfType"`
	Re_memory      uint   `yaml:"re_memory"`
	Pfe_memory     uint   `yaml:"pfe_memory"`
	Re_Cores       int    `yaml:re_cores"`
	Pfe_Cores      int    `yaml:pfe_cores"`
	Re_port        string `yaml:re_port"`
	Pfe_port       string `yaml:pfe_port"`
	Re_Vcpu        uint   `yaml:re_vcpu"`
	Pfe_Vcpu       uint   `yaml:pfe_vcpu"`
	Links          []Link `yaml:"link"`
}

type Nodes struct {
	Network_nodes []Node `yaml:"network_nodes"`
}

/* Read XM template
Currently not used.
*/
func Readtemplate(file string) string {
	xmlfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlfile.Close()
	val, _ := ioutil.ReadAll(xmlfile)
	return string(val)
}

/* Func create brdges
- Create bridges based on link names and vNF type
- Place the interfaces into the bridges. i.e. while generating XML templates
  load these as target bridges
*/

func CreateBridges(node Nodes) {
	/* Define common mgmt bridg */
	qattr_ext := netlink.NewLinkAttrs()
	qattr_ext.Name = BRIDGE_MGMT
	qattr_ext.MTU = BRIDGE_MTU
	qbattr_ext := &netlink.Bridge{LinkAttrs: qattr_ext}
	err_ext := netlink.LinkAdd(qbattr_ext)
	if err_ext != nil {
		fmt.Printf("Bridge %s not added: %v\n", qattr_ext.Name, err_ext)
	}
	qsattr_ext, _ := netlink.LinkByName(qattr_ext.Name)
	if serr := netlink.LinkSetUp(qsattr_ext); serr != nil {
		fmt.Printf("Bridge not up: %v\n", serr)
	}
	for i := 0; i < len(node.Network_nodes); i++ {
		if node.Network_nodes[i].VnfType == "vqfx" {
			fmt.Printf("Entering vqfx creating bridges\n")
			qattr_int := netlink.NewLinkAttrs()
			qattr_res := netlink.NewLinkAttrs()
			qattr_int.Name = node.Network_nodes[i].Name + "_int"
			qattr_res.Name = VQFX_BRIDGE_RES
			qattr_int.MTU = BRIDGE_MTU
			qattr_res.MTU = BRIDGE_MTU
			qbattr_int := &netlink.Bridge{LinkAttrs: qattr_int}
			qbattr_res := &netlink.Bridge{LinkAttrs: qattr_res}

			err_int := netlink.LinkAdd(qbattr_int)
			if err_int != nil {
				fmt.Printf("Bridge %s not added: %v\n", qattr_int.Name, err_int)
			}
			err_res := netlink.LinkAdd(qbattr_res)
			if err_res != nil {
				fmt.Printf("Bridge %s not added: %v\n", qattr_res.Name, err_res)
			}

			qsattr_int, _ := netlink.LinkByName(qattr_int.Name)
			if serr := netlink.LinkSetUp(qsattr_int); serr != nil {
				fmt.Printf("Bridge not up: %v\n", serr)
			}
			qsattr_res, _ := netlink.LinkByName(qattr_res.Name)
			if serr := netlink.LinkSetUp(qsattr_res); serr != nil {
				fmt.Printf("Bridge not up: %v\n", serr)
			}

		} else if node.Network_nodes[i].VnfType == "vmx" {
			fmt.Print("Creating int and ext bridges for VMX [WIP]\n")
			qattr_int := netlink.NewLinkAttrs()
			qattr_int.Name = "br-int-" + node.Network_nodes[i].Name
			qattr_int.MTU = BRIDGE_MTU
			qbattr_int := &netlink.Bridge{LinkAttrs: qattr_int}
			err_int := netlink.LinkAdd(qbattr_int)
			if err_int != nil {
				fmt.Printf("Bridge %s not added: %v\n", qattr_int.Name, err_int)
			}
			qsattr_int, _ := netlink.LinkByName(qattr_int.Name)
			if serr := netlink.LinkSetUp(qsattr_int); serr != nil {
				fmt.Printf("Bridge not up: %v\n", serr)
			}
		}
		for j := 0; j < len(node.Network_nodes[i].Links); j++ {
			//fmt.Printf("+v\n", node.Network_nodes[i].Links)
			linkattr := netlink.NewLinkAttrs()
			linkattr.Name = node.Network_nodes[i].Links[j].Name
			linkattr.MTU = BRIDGE_MTU
			bridge := &netlink.Bridge{LinkAttrs: linkattr}
			err := netlink.LinkAdd(bridge)
			if err != nil {
				fmt.Printf("Link %s not added: %v\n", linkattr.Name, err)
				continue
			}
			linkname, _ := netlink.LinkByName(linkattr.Name)

			if err := netlink.LinkSetUp(linkname); err != nil {
				fmt.Printf("link not up : %s", err)
			}
		}
	}
}

/* Function to delete bridges when deleting the topology */
func DeleteBridges(node Nodes) {
	qattr_ext := netlink.NewLinkAttrs()
	qattr_ext.Name = BRIDGE_MGMT
	qbattr_ext := &netlink.Bridge{LinkAttrs: qattr_ext}
	qsattr_ext, serr := netlink.LinkByName(qattr_ext.Name)
	if serr != nil {
		fmt.Printf("Link not found \n")
	} else {
		serr := netlink.LinkSetDown(qsattr_ext)
		if serr != nil {
			fmt.Printf("Link %s not present %s\n", qsattr_ext, serr)
		}
		errext := netlink.LinkDel(qbattr_ext)
		if errext != nil {
			fmt.Print("Bridge %s not deleted\n", qattr_ext.Name)
		}
	}
	for i := 0; i < len(node.Network_nodes); i++ {
		if node.Network_nodes[i].VnfType == "vqfx" {
			qattr_int := netlink.NewLinkAttrs()
			qattr_res := netlink.NewLinkAttrs()
			qattr_int.Name = node.Network_nodes[i].Name + "_int"
			qattr_res.Name = VQFX_BRIDGE_RES
			qbattr_int := &netlink.Bridge{LinkAttrs: qattr_int}
			qbattr_res := &netlink.Bridge{LinkAttrs: qattr_res}

			qsattr_int, serr := netlink.LinkByName(qattr_int.Name)
			if serr != nil {
				fmt.Printf("Link not found\n")
			} else {
				serr := netlink.LinkSetDown(qsattr_int)
				if serr != nil {
					fmt.Printf("Link not present %s\n", serr)
				}
				errint := netlink.LinkDel(qbattr_int)
				if errint != nil {
					fmt.Print("Bridge %s not deleted\n", qattr_int.Name)
				}
			}
			qsattr_res, serr := netlink.LinkByName(qattr_res.Name)
			if serr != nil {
				fmt.Printf("Link %s not found\n", qattr_res.Name)
			} else {
				serr := netlink.LinkSetDown(qsattr_res)
				if serr != nil {
					fmt.Printf("Link not present %s\n", serr)
				}
				errres := netlink.LinkDel(qbattr_res)
				if errres != nil {
					fmt.Print("Bridge %s not deleted\n", qattr_res.Name)
				}
			}
		} else if node.Network_nodes[i].VnfType == "vmx" {
			fmt.Printf("Deleting vmx int and vmx ext bridges\n")
			qattr_int := netlink.NewLinkAttrs()
			qattr_int.Name = "br-int-" + node.Network_nodes[i].Name
			qbattr_int := &netlink.Bridge{LinkAttrs: qattr_int}
			qsattr_int, serr := netlink.LinkByName(qattr_int.Name)
			if serr != nil {
				fmt.Printf("link not found\n")
			} else {
				serr := netlink.LinkSetDown(qsattr_int)
				if serr != nil {
					fmt.Printf("Bridge %s not deleted\n", qattr_int.Name)
				}
				errint := netlink.LinkDel(qbattr_int)
				if errint != nil {
					fmt.Printf("Bridge %s not deleted\n", qattr_int.Name)
				}
			}
		}
		/* There is no RE-PFE connecting bridges for vSRX so nothing to delete there */
		for j := 0; j < len(node.Network_nodes[i].Links); j++ {
			linkattr := netlink.NewLinkAttrs()
			linkattr.Name = node.Network_nodes[i].Links[j].Name
			linkname, err := netlink.LinkByName(linkattr.Name)
			if err != nil {
				fmt.Printf("Link not found\n")
				continue
			}
			if err := netlink.LinkSetDown(linkname); err != nil {
				fmt.Printf("Link not present : %s\n", err)
			}
			bridge := &netlink.Bridge{LinkAttrs: linkattr}
			erro := netlink.LinkDel(bridge)
			if erro != nil {
				fmt.Print("Bridge %s not deleted: %v\n", linkattr.Name, erro)
			}
		}
	}
}

/* Func start and stop Domains
- start domain (action create)
- stop domain (action delete)

/* Function to create new device disk to append into array devices list */
func NewDisk(path string, imgtype string, dev string, bus string, name string, slot uint) libvirtxml.DomainDisk {
	var val uint = 0
	if bus == "virtio" {
		return libvirtxml.DomainDisk{
			Device: "disk",
			//BackingStore: {},
			Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: imgtype},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: path}},
			Target: &libvirtxml.DomainDiskTarget{
				Dev: dev,
				Bus: bus},
			Alias: &libvirtxml.DomainAlias{
				Name: imgtype},
			Address: &libvirtxml.DomainAddress{
				PCI: &libvirtxml.DomainAddressPCI{
					Domain:   &val,
					Bus:      &val,
					Slot:     &slot,
					Function: &val,
				},
			},
		}
	} else {
		return libvirtxml.DomainDisk{
			Device: "disk",
			//BackingStore: {},
			Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: imgtype},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: path}},
			Target: &libvirtxml.DomainDiskTarget{
				Dev: dev,
				Bus: bus},
			Alias: &libvirtxml.DomainAlias{
				Name: imgtype},
			Address: &libvirtxml.DomainAddress{
				Drive: &libvirtxml.DomainAddressDrive{
					Controller: &val,
					Bus:        &val,
					Target:     &val,
					Unit:       &val,
				},
			},
		}
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
		return libvirtxml.DomainController{
			Type:  ctype,
			Index: &index,
			Model: model,
			Alias: &libvirtxml.DomainAlias{Name: alias},
		}
	} else {
		return libvirtxml.DomainController{
			Type:  ctype,
			Index: &index,
			Model: model,
			Alias: &libvirtxml.DomainAlias{Name: alias},
			Address: &libvirtxml.DomainAddress{
				PCI: &libvirtxml.DomainAddressPCI{
					Domain:        &pcidomain,
					Bus:           &pcibus,
					Slot:          &pcislot,
					Function:      &pcifunc,
					MultiFunction: pcimultifunc,
				},
			},
		}
	}
}

/* Function to create new Console type device */
func NewConsole(port uint, portnum string, name string) libvirtxml.DomainConsole {
	return libvirtxml.DomainConsole{
		Source: &libvirtxml.DomainChardevSource{
			TCP: &libvirtxml.DomainChardevSourceTCP{
				Mode:    "bind",
				Host:    LISTEN_ADDRESS,
				Service: portnum,
				TLS:     "no",
			},
		},
		Protocol: &libvirtxml.DomainChardevProtocol{
			Type: "telnet",
		},
		Target: &libvirtxml.DomainConsoleTarget{
			Type: "serial",
			Port: &port,
		},
		Alias: &libvirtxml.DomainAlias{
			Name: name,
		},
	}
}

/* Function to creae new Serial type device */
func NewSerial(port uint, portnum string, name string) libvirtxml.DomainSerial {
	return libvirtxml.DomainSerial{
		Source: &libvirtxml.DomainChardevSource{
			TCP: &libvirtxml.DomainChardevSourceTCP{
				Mode:    "bind",
				Host:    LISTEN_ADDRESS,
				Service: portnum,
				TLS:     "no",
			},
		},
		Protocol: &libvirtxml.DomainChardevProtocol{
			Type: "telnet",
		},
		Target: &libvirtxml.DomainSerialTarget{
			Type: "isa-serial",
			Port: &port,
		},
		Alias: &libvirtxml.DomainAlias{
			Name: name,
		},
	}
}

/* Function to create new Input such as PS2 */
func NewInputs(itype string, bus string, name string) libvirtxml.DomainInput {
	return libvirtxml.DomainInput{
		Type: itype,
		Bus:  bus,
		Alias: &libvirtxml.DomainAlias{
			Name: name,
		},
	}
}

/* Function for listeners. This has to be appended into Func Graphics() */
func GraphicListeners(address string) libvirtxml.DomainGraphicListener {
	return libvirtxml.DomainGraphicListener{
		Address: &libvirtxml.DomainGraphicListenerAddress{
			Address: address,
		},
	}
}

/* Function to enable graphics for the VM */
func Graphics(listener []libvirtxml.DomainGraphicListener) libvirtxml.DomainGraphic {
	return libvirtxml.DomainGraphic{
		VNC: &libvirtxml.DomainGraphicVNC{
			AutoPort:  "yes",
			Listen:    LISTEN_ADDRESS,
			Listeners: listener,
			//Listeners is a [] and needs to appended during initialization
		},
	}
}

/* Function for Features array in CPU */
func CpuFeatures(policy string, name string) libvirtxml.DomainCPUFeature {
	return libvirtxml.DomainCPUFeature{
		Policy: policy,
		Name:   name,
	}
}

/* Function for Video */
func Videos(ram uint, vram uint, vgamem uint, alias string,
	pcidomain uint, pcibus uint, pcislot uint,
	pcifunction uint) libvirtxml.DomainVideo {
	return libvirtxml.DomainVideo{
		Model: libvirtxml.DomainVideoModel{
			Type:   "qxl",
			Ram:    ram,
			VRam:   vram,
			VGAMem: vgamem,
			Heads:  1,
		},
		Alias: &libvirtxml.DomainAlias{Name: alias},
		Address: &libvirtxml.DomainAddress{
			PCI: &libvirtxml.DomainAddressPCI{
				Domain:        &pcidomain,
				Bus:           &pcibus,
				Slot:          &pcislot,
				Function:      &pcifunction,
				MultiFunction: "on",
			},
		},
	}
}

/* Function to create new interfaces based on link mentioned in topology.yaml
Since the revenue interfaces dont have the Address params, using a flag type to
determine whether the template for DomainAddress is needed or not.
when Domain address not needed, pass 999 to pcidomain and other params. only PCI
domain is added to check however since it encountered first anyway.
WIP: look for a better way to templatize
*/
func NewNetworkIntf(mac string, link string, intf string,
	model string, pcidomain uint, pcibus uint,
	pcislot uint, pcifunction uint) libvirtxml.DomainInterface {
	if pcidomain == 999 {
		return libvirtxml.DomainInterface{
			MAC: &libvirtxml.DomainInterfaceMAC{Address: mac},
			Source: &libvirtxml.DomainInterfaceSource{
				Bridge: &libvirtxml.DomainInterfaceSourceBridge{
					Bridge: link,
				},
			},
			Target: &libvirtxml.DomainInterfaceTarget{
				Dev: intf,
			},
			Model: &libvirtxml.DomainInterfaceModel{
				Type: model,
			},
		}
	} else {
		return libvirtxml.DomainInterface{
			MAC: &libvirtxml.DomainInterfaceMAC{Address: mac},
			Source: &libvirtxml.DomainInterfaceSource{
				Bridge: &libvirtxml.DomainInterfaceSourceBridge{
					Bridge: link,
				},
			},
			Target: &libvirtxml.DomainInterfaceTarget{
				Dev: intf,
			},
			Model: &libvirtxml.DomainInterfaceModel{
				Type: model,
			},
			Address: &libvirtxml.DomainAddress{
				PCI: &libvirtxml.DomainAddressPCI{
					Domain:   &pcidomain,
					Bus:      &pcibus,
					Slot:     &pcislot,
					Function: &pcifunction,
				},
			},
		}
	}
}

/* Function DomTemplate */
func DomTemplate(name string, memory uint, cores int, vcpu uint) *libvirtxml.Domain {
	/* const nint is used to controller the display of certain xmls
	   which need to be declared as null since we cannot pass uint(null)
	   . Check if we can pass optional arguments to the function that we way
	   we can avoid duplication of the objects*/
	ipci := uint(0x00)
	domcfg := &libvirtxml.Domain{
		Type: "kvm",
		Name: name,
		Memory: &libvirtxml.DomainMemory{
			Value: memory,
			Unit:  "MB",
		},
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			Value:     vcpu},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    "x86_64",
				Machine: "pc-i440fx-xenial",
				Type:    "hvm"}},
		Features: &libvirtxml.DomainFeatureList{
			APIC: &libvirtxml.DomainFeatureAPIC{
				EOI: ""},
			ACPI: &libvirtxml.DomainFeature{},
		},
		CPU: &libvirtxml.DomainCPU{
			Mode: "host-model",
			Topology: &libvirtxml.DomainCPUTopology{
				Sockets: 1,
				Cores:   cores,
				Threads: 1}},
		OnPoweroff: "destroy",
		OnReboot:   "restart",
		OnCrash:    "restart",
		PM: &libvirtxml.DomainPM{
			SuspendToMem:  &libvirtxml.DomainPMPolicy{Enabled: "no"},
			SuspendToDisk: &libvirtxml.DomainPMPolicy{Enabled: "no"}},
		Devices: &libvirtxml.DomainDeviceList{
			Emulator: "/usr/bin/qemu-system-x86_64",
			MemBalloon: &libvirtxml.DomainMemBalloon{
				Model: "virtio",
				Alias: &libvirtxml.DomainAlias{Name: "balloon"},
				Address: &libvirtxml.DomainAddress{
					PCI: &libvirtxml.DomainAddressPCI{
						Domain:   &ipci,
						Bus:      &ipci,
						Slot:     &ipci,
						Function: &ipci,
					},
				},
			},
		},
	}
	return domcfg
}

/* Function for vQFX-RE Template generation */
func TemplateVqfx(node Node, devid int) (*libvirtxml.Domain, *libvirtxml.Domain) {
	const nint uint = 999
	domcfg := DomTemplate(node.Name+"-re", uint(node.Re_memory), node.Re_Cores, uint(node.Re_Vcpu))
	dompfecfg := DomTemplate(node.Name+"-pfe", uint(node.Pfe_memory), node.Pfe_Cores, uint(node.Pfe_Vcpu))
	disk_re := NewDisk(node.ImagePath+node.Re_Image, "qcow2", "hda", "ide", "ide-0-0", nint)
	disk_pfe := NewDisk(node.ImagePath+node.Pfe_Image, "qcow2", "hda", "ide", "ide-0-0", nint)

	domcfg.Devices.Disks = append(domcfg.Devices.Disks, disk_re)
	dompfecfg.Devices.Disks = append(dompfecfg.Devices.Disks, disk_pfe)

	/* Reference for function params
	   func NewController(ctype string, index uint, model string,
	                    alias string, pcidomain uint, pcibus uint,
	                    pcislot uint, pcifunc uint, pcimultifunc string,
	                   ) libvirtxml.DomainController {
	*/
	c1 := NewController("usb", uint(0), "ich9-ehci1", "usb", uint(0x0000), uint(0x00), uint(0x0b), uint(0x7), "")
	c2 := NewController("usb", uint(0), "ich9-uhci1", "usb", uint(0x0000), uint(0x00), uint(0x0b), uint(0x0), "on")
	c3 := NewController("usb", uint(0), "ich9-uhci2", "usb", uint(0x0000), uint(0x00), uint(0x0b), uint(0x1), "")
	c4 := NewController("usb", uint(0), "ich9-uhci3", "usb", uint(0x0000), uint(0x00), uint(0x0b), uint(0x2), "")
	c5 := NewController("pci", uint(0), "pci-root", "pci.0", nint, nint, nint, nint, "")
	c6 := NewController("ide", uint(0), "", "ide", uint(0x0000), uint(0x00), uint(0x01), uint(0x01), "")
	c7 := NewController("virtio-serial", uint(0), "", "virtio-serial0", uint(0x0000), uint(0x00), uint(0x0a), uint(0x0), "")

	/* Reference for serial function params
	   func NewSerial(port *uint, portnum string, name string)
	*/
	s1 := NewSerial(uint(0), node.Re_port, "serial0")
	s1_pfe := NewSerial(uint(0), node.Pfe_port, "serial0")

	/* Reference for input function params
	   func NewInputs(itype string, bus string, name string)
	*/
	i1 := NewInputs("mouse", "ps2", "mouse")
	i2 := NewInputs("keyboard", "ps2", "keyboard")

	/* Reference for console function params
	   func NewConsole(port *uint, portnum string, name string)
	*/
	co := NewConsole(uint(0), node.Re_port, "serial0")
	co_pfe := NewConsole(uint(0), node.Pfe_port, "serial0")

	/* reference for Graphics function params
	   func GraphicListeners(address string)
	*/
	gl := GraphicListeners(LISTEN_ADDRESS)
	gla := []libvirtxml.DomainGraphicListener{gl}
	g1 := Graphics(gla)
	/* Reference for Video function params
	   func Videos(ram uint, vram uint, vgamem uint, alias string,
	           pcidomain *uint, pcibus *uint, pcislot *uint,
	           pcifunction *uint)
	*/
	v1 := Videos(65536, 65536, 16384, "video0", uint(0x0000), uint(0x00), uint(0x02), uint(0x0))

	domcfg.Devices.Controllers = append(domcfg.Devices.Controllers, c1, c2, c3, c4, c5, c6, c7)
	dompfecfg.Devices.Controllers = append(dompfecfg.Devices.Controllers, c1, c2, c3, c4, c5, c6, c7)

	domcfg.Devices.Serials = append(domcfg.Devices.Serials, s1)
	dompfecfg.Devices.Serials = append(dompfecfg.Devices.Serials, s1_pfe)

	domcfg.Devices.Inputs = append(domcfg.Devices.Inputs, i1, i2)
	dompfecfg.Devices.Inputs = append(dompfecfg.Devices.Inputs, i1, i2)

	domcfg.Devices.Consoles = append(domcfg.Devices.Consoles, co)
	dompfecfg.Devices.Consoles = append(dompfecfg.Devices.Consoles, co_pfe)

	domcfg.Devices.Graphics = append(domcfg.Devices.Graphics, g1)
	dompfecfg.Devices.Graphics = append(dompfecfg.Devices.Graphics, g1)

	domcfg.Devices.Videos = append(domcfg.Devices.Videos, v1)
	dompfecfg.Devices.Videos = append(dompfecfg.Devices.Videos, v1)

	hdevid := fmt.Sprintf("%02x", devid+1)
	/* Add only vqfx-int and vqfx-mgmt interfaces into vQFX-PFE */
	intf_int := NewNetworkIntf(strings.TrimSuffix(BASE_VQFX_INT_MAC_ADDRESS, "00")+hdevid,
		BRIDGE_MGMT, node.Name+"_pfe-int", "e1000", uint(0x0000),
		uint(0x00), uint(0x03), uint(0x0))
	intf_ext := NewNetworkIntf(strings.TrimSuffix(BASE_VQFX_EXT_MAC_ADDRESS, "00")+hdevid,
		node.Name+"_int", node.Name+"_pfe-ext", "e1000", uint(0x0000),
		uint(0x00), uint(0x04), uint(0x00))
	dompfecfg.Devices.Interfaces = append(dompfecfg.Devices.Interfaces, intf_int, intf_ext)

	/* Add vqfx int, ext and res interfaces for RE. These would be the first 3 interfaces */
	re_int := NewNetworkIntf(strings.TrimSuffix(BASE_VQFX_INT_MAC_ADDRESS, "00:00")+hdevid+":"+hdevid,
		BRIDGE_MGMT, node.Name+"_re-int", "e1000", uint(0x0000), uint(0x00),
		uint(0x03), uint(0x00))
	re_ext := NewNetworkIntf(strings.TrimSuffix(BASE_VQFX_EXT_MAC_ADDRESS, "00:00")+hdevid+":"+hdevid,
		node.Name+"_int", node.Name+"_re-ext", "e1000", uint(0x0000), uint(0x00), uint(0x04),
		uint(0x00))
	re_res := NewNetworkIntf(strings.TrimSuffix(BASE_VQFX_RES_MAC_ADDRESS, "00:00")+hdevid+":"+hdevid,
		VQFX_BRIDGE_RES, node.Name+"_re-res", "e1000", uint(0x0000), uint(0x00), uint(0x05),
		uint(0x00))
	domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, re_int, re_ext, re_res)

	/* Add vqfx xe interfaces. This would be added to teh RE */
	for i := 0; i < len(node.Links); i++ {
		i_to_hex := fmt.Sprintf("%02x", i)
		nintf := NewNetworkIntf(strings.TrimSuffix(BASE_VQFX_LINK_MAC_ADDRESS, "00:00")+hdevid+":"+i_to_hex,
			node.Links[i].Name, node.Links[i].Intf, "e1000", nint, nint, nint, nint)
		domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, nintf)
	}
	return domcfg, dompfecfg
}

func TemplateVsrx(node Node, devid int) *libvirtxml.Domain {
	const nint uint = 999
	domcfg := DomTemplate(node.Name, uint(node.Re_memory), node.Re_Cores, uint(node.Re_Vcpu))
	disk := NewDisk(node.ImagePath+node.Re_Image, "qcow2", "hda", "ide", "ide-0-0", nint)
	domcfg.Devices.Disks = append(domcfg.Devices.Disks, disk)

	c1 := NewController("usb", uint(0), "ich9-ehci1", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x7), "")
	c2 := NewController("usb", uint(0), "ich9-uhci1", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x0), "on")
	c3 := NewController("usb", uint(0), "ich9-uhci2", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x1), "")
	c4 := NewController("usb", uint(0), "ich9-uhci3", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x2), "")
	c5 := NewController("pci", uint(0), "pci-root", "pci.0", nint, nint, nint, nint, "")
	c6 := NewController("virtio-serial", uint(0), "", "virtio-serial", uint(0x0000), uint(0x00), uint(0x08), uint(0x0), "")

	s1 := NewSerial(uint(0), node.Re_port, "serial")

	i1 := NewInputs("mouse", "ps2", "mouse")
	i2 := NewInputs("keyboard", "ps2", "keyboard")

	co := NewConsole(uint(0), node.Re_port, "serial0")

	gl := GraphicListeners(LISTEN_ADDRESS)
	gla := []libvirtxml.DomainGraphicListener{gl}
	g1 := Graphics(gla)

	v1 := Videos(65536, 65536, 16384, "video0", uint(0x0000), uint(0x00), uint(0x02), uint(0x0))
	domcfg.Devices.Controllers = append(domcfg.Devices.Controllers, c1, c2, c3, c4, c5, c6)
	domcfg.Devices.Serials = append(domcfg.Devices.Serials, s1)
	domcfg.Devices.Inputs = append(domcfg.Devices.Inputs, i1, i2)
	domcfg.Devices.Consoles = append(domcfg.Devices.Consoles, co)
	domcfg.Devices.Graphics = append(domcfg.Devices.Graphics, g1)
	domcfg.Devices.Videos = append(domcfg.Devices.Videos, v1)

	hdevid := fmt.Sprintf("%02x", devid+1)
	/* Add FXP0 interface. This is the first interface added */
	intf_ext := NewNetworkIntf(strings.TrimSuffix(BASE_VSRX_EXT_MAC_ADDRESS, "00")+hdevid,
		BRIDGE_MGMT, node.Name+"-ext", "virtio", uint(0x0000),
		uint(0x00), uint(0x03), uint(0x0))
	domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, intf_ext)

	/* Add the rest of the ge interfaces */
	for i := 0; i < len(node.Links); i++ {
		i_to_hex := fmt.Sprintf("%02x", i)
		nintf := NewNetworkIntf(strings.TrimSuffix(BASE_VSRX_LINK_MAC_ADDRESS, "00:00")+hdevid+":"+i_to_hex,
			node.Links[i].Name, node.Links[i].Intf, "virtio", nint, nint, nint, nint)
		domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, nintf)
	}
	return domcfg
}

/* Function for defining vMX template */
func TemplateVmx(node Node, devid int) (*libvirtxml.Domain, *libvirtxml.Domain) {
	const nint uint = 999
	domcfg := DomTemplate(node.Name+"-re", uint(node.Re_memory), node.Re_Cores, uint(node.Re_Vcpu))
	dompfecfg := DomTemplate(node.Name+"-pfe", uint(node.Pfe_memory), node.Pfe_Cores, uint(node.Pfe_Vcpu))
	disk_re := NewDisk(node.ImagePath+node.Re_Image, "qcow2", "vda", "virtio", "virtio-disk0", uint(0x7))
	disk_hdd := NewDisk(node.ImagePath+node.Vhdd_Image, "qcow2", "vdb", "virtio", "virtio-disk1", uint(0x8))
	disk_usb := NewDisk(node.ImagePath+node.Metadata_Image, "raw", "vdc", "virtio", "virtio-disk2", uint(0x9))
	disk_pfe := NewDisk(node.ImagePath+node.Pfe_Image, "raw", "hda", "ide", "ide0-0-0", nint)

	domcfg.Devices.Disks = append(domcfg.Devices.Disks, disk_re, disk_hdd, disk_usb)
	dompfecfg.Devices.Disks = append(dompfecfg.Devices.Disks, disk_pfe)

	//dompfecfg.Devices.Disks = append(dompfecfg.Devices.Disks, disk_pfe)
	c1 := NewController("usb", uint(0), "none", "usb", nint, nint, nint, nint, "")
	c2 := NewController("pci", uint(0), "pci-root", "pci.0", nint, nint, nint, nint, "")

	c1p := NewController("pci", uint(0), "pci-root", "pci.0", nint, nint, nint, nint, "")
	//c2p := NewController("usb", uint(0), "piix3-uhci", "usb",uint(0x0000), uint(0x00),uint(0x01), uint(0x2), "")
	//c3p := NewController("ide", uint(0), "", "ide",uint(0x0000), uint(0x00),uint(0x01), uint(0x1), "")

	s1 := NewSerial(uint(0), node.Re_port, "serial0")
	s1p := NewSerial(uint(0), node.Pfe_port, "serial0")
	co := NewConsole(uint(0), node.Re_port, "serial0")
	cop := NewConsole(uint(0), node.Pfe_port, "serial0")

	i1 := NewInputs("mouse", "ps2", "input0")
	i2 := NewInputs("keyboard", "ps2", "input1")

	gl := GraphicListeners(LISTEN_ADDRESS)
	gla := []libvirtxml.DomainGraphicListener{gl}
	g1 := Graphics(gla)

	v1 := Videos(65536, 65536, 16384, "video", uint(0x0000), uint(0x01), uint(0x03), uint(0x07))

	domcfg.Devices.Controllers = append(domcfg.Devices.Controllers, c1, c2)
	domcfg.Devices.Serials = append(domcfg.Devices.Serials, s1)
	domcfg.Devices.Inputs = append(domcfg.Devices.Inputs, i1, i2)
	domcfg.Devices.Consoles = append(domcfg.Devices.Consoles, co)
	domcfg.Devices.Graphics = append(domcfg.Devices.Graphics, g1)
	domcfg.Devices.Videos = append(domcfg.Devices.Videos, v1)

	dompfecfg.Devices.Controllers = append(dompfecfg.Devices.Controllers, c1p)
	dompfecfg.Devices.Serials = append(dompfecfg.Devices.Serials, s1p)
	dompfecfg.Devices.Consoles = append(dompfecfg.Devices.Consoles, cop)
	dompfecfg.Devices.Inputs = append(dompfecfg.Devices.Inputs, i1, i2)
	dompfecfg.Devices.Graphics = append(dompfecfg.Devices.Graphics, g1)
	dompfecfg.Devices.Videos = append(dompfecfg.Devices.Videos, v1)

	hdevid := fmt.Sprintf("%02x", devid+1)

	/* Add EXT, INT bridges for VCP */
	intf_ext := NewNetworkIntf(strings.TrimSuffix(BASE_VMX_EXT_MAC_ADDRESS, "00")+hdevid,
		BRIDGE_MGMT, node.Name+"_re-ext", "virtio", uint(0x0000),
		uint(0x00), uint(0x03), uint(0x0))
	intf_int := NewNetworkIntf(strings.TrimSuffix(BASE_VMX_INT_MAC_ADDRESS, "00")+hdevid,
		"br-int-"+node.Name, node.Name+"_re-int", "virtio", uint(0x0000),
		uint(0x00), uint(0x04), uint(0x00))
	domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, intf_ext, intf_int)

	/* Add EXT, INT briges for VFP */
	pintf_ext := NewNetworkIntf(strings.TrimSuffix(BASE_VMX_EXT_MAC_ADDRESS, "00:00")+hdevid+":"+hdevid,
		BRIDGE_MGMT, node.Name+"_pfe-ext", "virtio", uint(0x0000),
		uint(0x00), uint(0x03), uint(0x0))
	pintf_int := NewNetworkIntf(strings.TrimSuffix(BASE_VMX_INT_MAC_ADDRESS, "00:00")+hdevid+":"+hdevid,
		"br-int-"+node.Name, node.Name+"_pfe-int", "virtio", uint(0x0000),
		uint(0x00), uint(0x04), uint(0x00))
	dompfecfg.Devices.Interfaces = append(dompfecfg.Devices.Interfaces, pintf_ext, pintf_int)

	/* Add regular ge interfaces to the PFE */
	for i := 0; i < len(node.Links); i++ {
		i_to_hex := fmt.Sprintf("%02x", i)
		nintf := NewNetworkIntf(strings.TrimSuffix(BASE_VMX_LINK_MAC_ADDRESS, "00:00")+hdevid+":"+i_to_hex,
			node.Links[i].Name, node.Links[i].Intf, "virtio", nint, nint, nint, nint)
		dompfecfg.Devices.Interfaces = append(dompfecfg.Devices.Interfaces, nintf)
	}
	return domcfg, dompfecfg
}

func TemplateLinux(node Node, devid int) *libvirtxml.Domain {
	const nint uint = 999
	domcfg := DomTemplate(node.Name, uint(node.Re_memory), node.Re_Cores, uint(node.Re_Vcpu))
	disk := NewDisk(node.ImagePath+node.Re_Image, "qcow2", "hda", "ide", "ide-0-0", nint)
	domcfg.Devices.Disks = append(domcfg.Devices.Disks, disk)

	c1 := NewController("usb", uint(0), "ich9-ehci1", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x7), "")
	c2 := NewController("usb", uint(0), "ich9-uhci1", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x0), "on")
	c3 := NewController("usb", uint(0), "ich9-uhci2", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x1), "")
	c4 := NewController("usb", uint(0), "ich9-uhci3", "usb", uint(0x0000), uint(0x00), uint(0x09), uint(0x2), "")
	c5 := NewController("pci", uint(0), "pci-root", "pci.0", nint, nint, nint, nint, "")
	c6 := NewController("virtio-serial", uint(0), "", "virtio-serial", uint(0x0000), uint(0x00), uint(0x08), uint(0x0), "")

	s1 := NewSerial(uint(0), node.Re_port, "serial")

	i1 := NewInputs("mouse", "ps2", "mouse")
	i2 := NewInputs("keyboard", "ps2", "keyboard")

	co := NewConsole(uint(0), node.Re_port, "serial0")

	gl := GraphicListeners(LISTEN_ADDRESS)
	gla := []libvirtxml.DomainGraphicListener{gl}
	g1 := Graphics(gla)

	v1 := Videos(65536, 65536, 16384, "video0", uint(0x0000), uint(0x00), uint(0x02), uint(0x0))
	domcfg.Devices.Controllers = append(domcfg.Devices.Controllers, c1, c2, c3, c4, c5, c6)
	domcfg.Devices.Serials = append(domcfg.Devices.Serials, s1)
	domcfg.Devices.Inputs = append(domcfg.Devices.Inputs, i1, i2)
	domcfg.Devices.Consoles = append(domcfg.Devices.Consoles, co)
	domcfg.Devices.Graphics = append(domcfg.Devices.Graphics, g1)
	domcfg.Devices.Videos = append(domcfg.Devices.Videos, v1)

	hdevid := fmt.Sprintf("%02x", devid+1)
	/* Add FXP0 interface. This is the first interface added */
	intf_ext := NewNetworkIntf(strings.TrimSuffix(BASE_LINUX_EXT_MAC_ADDRESS, "00")+hdevid,
		BRIDGE_MGMT, node.Name+"-ext", "virtio", uint(0x0000),
		uint(0x00), uint(0x03), uint(0x0))
	domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, intf_ext)

	/* Add the rest of the ge interfaces */
	for i := 0; i < len(node.Links); i++ {
		i_to_hex := fmt.Sprintf("%02x", i)
		nintf := NewNetworkIntf(strings.TrimSuffix(BASE_LINUX_LINK_MAC_ADDRESS, "00:00")+hdevid+":"+i_to_hex,
			node.Links[i].Name, node.Links[i].Intf, "virtio", nint, nint, nint, nint)
		domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, nintf)
	}
	return domcfg
}

/* Pass struct to function to handle topo spin up
Generate domain xml, define and start
*/
func Genxml(node Nodes) {
	for i := 0; i < len(node.Network_nodes); i++ {
		if node.Network_nodes[i].VnfType == "vmx" {
			fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
			domcfg, dompfecfg := TemplateVmx(node.Network_nodes[i], i)
			rexmldoc, err := domcfg.Marshal()
			pfexmldoc, err1 := dompfecfg.Marshal()
			if err != nil {
				fmt.Printf("Marshaling of VMX RE Failed %v\n", err)
			}
			if err1 != nil {
				fmt.Printf("Marshaling of VMX PFE Failed %v\n", err1)
			}
			if os.Args[1] == "genxml" {
				direc := "./templates/" + strings.TrimSuffix(os.Args[2], ".yaml")
				_, err = os.Stat(direc)
				if os.IsNotExist(err) {
					os.Mkdir(direc, 0777)
				}
				WriteToFile(node.Network_nodes[i].Name+"_vmx-re.xml", direc+"/", rexmldoc)
				WriteToFile(node.Network_nodes[i].Name+"_vmx-pfe.xml", direc+"/", pfexmldoc)
			} else {
				DomainStart(rexmldoc)
				DomainStart(pfexmldoc)
			}
		} else if node.Network_nodes[i].VnfType == "vqfx" {
			domcfg, dompfecfg := TemplateVqfx(node.Network_nodes[i], i)
			rexmldoc, err := domcfg.Marshal()
			pfexmldoc, err1 := dompfecfg.Marshal()
			if err != nil {
				fmt.Printf("Marshing of RE xml failed %v\n", err)
			}
			if err1 != nil {
				fmt.Printf("marshaling of PFE xml failed %v\n", err1)
			}
			/* if genxml save the xml into a directory */
			if os.Args[1] == "genxml" {
				direc := "./templates/" + strings.TrimSuffix(os.Args[2], ".yaml")
				_, err = os.Stat(direc)
				if os.IsNotExist(err) {
					os.Mkdir(direc, 0777)
				}
				WriteToFile(node.Network_nodes[i].Name+"_vqfx-re.xml", direc+"/", rexmldoc)
				WriteToFile(node.Network_nodes[i].Name+"_vqfx-pfe.xml", direc+"/", pfexmldoc)
			} else {
				DomainStart(rexmldoc)
				DomainStart(pfexmldoc)
			}
		} else if node.Network_nodes[i].VnfType == "vsrx" {
			fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
			domcfg := TemplateVsrx(node.Network_nodes[i], i)
			rexmldoc, err := domcfg.Marshal()
			if err != nil {
				fmt.Printf("Marshing of vSRX XML failed %v\n", err)
			}
			if os.Args[1] == "genxml" {
				// add or condition "||" to trim .yaml or .yml
				direc := "./templates/" + strings.TrimSuffix(os.Args[2], ".yaml")
				_, err = os.Stat(direc)
				if os.IsNotExist(err) {
					os.Mkdir(direc, 0777)
				}
				WriteToFile(node.Network_nodes[i].Name+"_vsrx.xml", direc+"/", rexmldoc)
			} else {
				DomainStart(rexmldoc)
			}
		} else if node.Network_nodes[i].VnfType == "linux" {
			fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
			domcfg := TemplateLinux(node.Network_nodes[i], i)
			rexmldoc, err := domcfg.Marshal()
			if err != nil {
				fmt.Printf("Marshing of Linux XML failed %v\n", err)
			}
			if os.Args[1] == "genxml" {
				// add or condition "||" to trim .yaml or .yml
				direc := "./templates/" + strings.TrimSuffix(os.Args[2], ".yaml")
				_, err = os.Stat(direc)
				if os.IsNotExist(err) {
					os.Mkdir(direc, 0777)
				}
				WriteToFile(node.Network_nodes[i].Name+"_linux.xml", direc+"/", rexmldoc)
			} else {
				DomainStart(rexmldoc)
			}
		}
	}
}

/* Function write to file */
func WriteToFile(name string, path string, data string) {
	cat := path + name
	f, err := os.Create(cat)
	if err != nil {
		fmt.Printf("Failed to create the file\n")
	}
	defer f.Close()
	_, err2 := f.WriteString(data)
	if err2 != nil {
		fmt.Printf(" Failed to write content into file\n")
	}
}

/* Function to copy images to build directory
Pass pointer to function and modify the source Image path
under node.Re_ImagePath and node.Pfe_ImagePath
*/
func CopyImage(nodes *Nodes) {
	// create directory based on node name
	fmt.Printf("Copying image from source path to build directory. Image will use the build dir path as source... \n")
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to obtain current working directory. Exiting program\n")
		os.Exit(1)
	}
	for i := 0; i < len(nodes.Network_nodes); i++ {
		dir := cwd + "/build/" + nodes.Network_nodes[i].Name + "/"
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			fmt.Printf("Directory doesnt exist under build, creating it!..\n")
			os.Mkdir(dir, 0777)
		}
		cpcmd_re := exec.Command("cp", "-rf", nodes.Network_nodes[i].ImagePath+nodes.Network_nodes[i].Re_Image, dir)
		err1 := cpcmd_re.Run()
		if err1 != nil {
			fmt.Printf("FAILED COPYING RE IMAGE TO BUILD DIRECTORY!!: %v \n", err1)
		}
		fmt.Printf("Sleeping for 10seconds to ensure file copies successfully\n")
		time.Sleep(10)
		// vSRX/Ubuntu would be single image. Use ony Re_ImagePath and leave Pfe_ImagePath blank
		if nodes.Network_nodes[i].Pfe_Image != "" {
			cpcmd_pfe := exec.Command("cp", "-rf", nodes.Network_nodes[i].ImagePath+nodes.Network_nodes[i].Pfe_Image, dir)
			err2 := cpcmd_pfe.Run()
			if err2 != nil {
				fmt.Printf("FAILED COPYING PFE IMAGE TO BUILD DIRECTORY !!: %v \n", err2)
			}
		}
		if nodes.Network_nodes[i].Metadata_Image != "" {
			cpcmd_meta := exec.Command("cp", "-rf", nodes.Network_nodes[i].ImagePath+nodes.Network_nodes[i].Metadata_Image, dir)
			err3 := cpcmd_meta.Run()
			if err3 != nil {
				fmt.Printf("FAILED COPYING METADATA IMAGE TO BUILD DIRECTORY !!: %v \n", err3)
			}
		}
		if nodes.Network_nodes[i].Vhdd_Image != "" {
			cpcmd_vhdd := exec.Command("cp", "-rf", nodes.Network_nodes[i].ImagePath+nodes.Network_nodes[i].Vhdd_Image, dir)
			err4 := cpcmd_vhdd.Run()
			if err4 != nil {
				fmt.Printf("FAILED COPYING VHDD IMAGE TO BUILD DIRECTORY !!: %v \n", err4)
			}
		}
		nodes.Network_nodes[i].ImagePath = dir
	}
}

func DomainStop(name string) {
	fmt.Printf("Destroying VMs and undefining.. \n")
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt: %v", err)
	}
	conn := libvirt.New(c)
	if err := conn.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	if err != nil {
		fmt.Printf("LIBVIRT CONNECTION FAILED! :%v", err)
	}
	dom, err := conn.DomainLookupByName(name)
	if err == nil {
		conn.DomainDestroy(dom)
		conn.DomainUndefine(dom)
	}
	fmt.Printf("Destroyed and undefined all domains under %s \n", os.Args[2])
}

func DomainStart(xml string) {
	fmt.Printf("Defining and Starting the Domains... \n")
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt: %v", err)
	}
	conn := libvirt.New(c)
	if err := conn.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	domain, err := conn.DomainDefineXML(xml)
	if err != nil {
		fmt.Printf("DOMAIN %v DEFINITION FAILED!!: %v\n", domain, err)
	}
	err = conn.DomainCreate(domain)
	if err != nil {
		fmt.Printf("DOMAIN CREATION FAILED!! %v", err)
	}
}

/*
Main function definition
*/
func main() {
	if len(os.Args) > 1 {
		var nodes Nodes
		if len(os.Args) > 2 {
			// Parse yaml generate xml based on topology templates
			rfile, err := ioutil.ReadFile(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}
			// parse yaml into structs
			err = yaml.Unmarshal(rfile, &nodes)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Missing argument Topology yaml file\n")
		}
		//var nodesptr *Nodes = &nodes //pass pointer so that we can operate on same address to modify vals
		if os.Args[1] == "create" {
			fmt.Printf("Creating bridges\n")
			//Parse through yaml and create all bridges needed to define topology
			CreateBridges(nodes)
			//copy original image to build directory and refer that xml, hence pass pointer.
			CopyImage(&nodes)
			Genxml(nodes)
		} else if os.Args[1] == "delete" {
			for i := 0; i < len(nodes.Network_nodes); i++ {
				if nodes.Network_nodes[i].VnfType == "vqfx" {
					DomainStop(nodes.Network_nodes[i].Name + "-re")
					DomainStop(nodes.Network_nodes[i].Name + "-pfe")
				} else if nodes.Network_nodes[i].VnfType == "vmx" {
					DomainStop(nodes.Network_nodes[i].Name + "-re")
					DomainStop(nodes.Network_nodes[i].Name + "-pfe")
				} else if nodes.Network_nodes[i].VnfType == "vsrx" {
					DomainStop(nodes.Network_nodes[i].Name)
				} else if nodes.Network_nodes[i].VnfType == "linux" {
					DomainStop(nodes.Network_nodes[i].Name)
				}
			}
			fmt.Printf("Deleting bridges\n")
			DeleteBridges(nodes)
		} else if os.Args[1] == "genxml" {
			fmt.Printf(" Generating xml files..  \n")
			Genxml(nodes)
		} else if os.Args[1] == "help" {
			fmt.Printf(` 
            -------------- VM Topology Builder V1.1 --------------------------------------------
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
		fmt.Printf(`
        INVALID OPTION. Follow the below usage

        +++++++++++++++++ usage +++++++++++++++++++++
        ./vm-topo <create/delete> <path to topo.yml>
        +++++++++++++++++++++++++++++++++++++++++++++`)
	}
}
