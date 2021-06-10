/*
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
    //"os"
    //"github.com/clagraff/argparse"
    //"encoding/xml"
    "gopkg.in/yaml.v3"
    "github.com/vishvananda/netlink"
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
    Re_memory string `yaml:"re_memory"`
    Pfe_memory string `yaml:"pfe_memory"`
    Links []Link `yaml:"link"`
}

type Nodes struct  {
    Network_nodes [] Node `yaml:"network_nodes"`
}

/* Read template 
func Readtemplate(file string) {
    xmlfile, err := os.Open(file)
    if err!= nil {
        fmt.Println(err)
    }
    defer xmlFile
}
*/

/* Func create brdges 
- Create bridges based on link names 
- Place the interfaces into the bridges. i.e. while generating XML templates 
  load these as target bridges 
*/

func CreateBridges(node Nodes) {
    for i:=0; i<len(node.Network_nodes); i++ {
        for j:=0; j<len(node.Network_nodes[i].Links); j++ {
            fmt.Printf("+v\n", node.Network_nodes[i].Links)
            linkattr := netlink.NewLinkAttrs() 
            linkattr.Name = node.Network_nodes[i].Links[j].Name
            linkattr.MTU = 9000 
            bridge := &netlink.Bridge{LinkAttrs: linkattr}
            err := netlink.LinkAdd(bridge)
            if err!=nil {
                fmt.Printf("Link %s not added: %v\n", linkattr.Name, err)
            }
            linkname,_ := netlink.LinkByName(linkattr.Name)

            if err := netlink.LinkSetUp(linkname); err !=nil {
                fmt.Printf("link not up : %s", err)
            }
        }
    }
}


/* Func XML unmarshal and Marshal 
WIP:
- Import template struct from another file 
- depending on vNF type place values with nodes struct
- generate XMl using xmlMarshal
- save it as xml in templates folder with respective domain name 
*/

/* Func start and stop Domains
- start domain (action create)
- stop domain (action delete)

Create/delete respective domain names and bridges only based on topology  
*/

/* Pass struct to function to handle topo spin up*/ 
func Genxml(node Nodes) {
    for i:=0; i<len(node.Network_nodes);i++ {
        if node.Network_nodes[i].VnfType == "vmx" {
            fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
            //Generate RE and PFE xml
        } else if node.Network_nodes[i].VnfType == "vqfx" {
            fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
            // Generate RE and PFE xml
        } else if node.Network_nodes[i].VnfType == "vsrx" {
            fmt.Printf("%+v\n", node.Network_nodes[i].VnfType)
            // Generate xml. vSRX uses a single image 
        }
    }
    //return XML
}


/*
Main function definition
*/
func main() {
    rfile, err := ioutil.ReadFile("topology.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // parse yaml into structs
    var nodes Nodes
    err = yaml.Unmarshal(rfile, &nodes)
    if  err != nil {
        log.Fatal(err)
    }

    CreateBridges(nodes)
    // Parse yaml generate xml based on topology templates
    Genxml(nodes)

    // create bridges accordingly on based link names 
    // use libvirt api to start the domains . Create and delete. when deleting, delete bridges also
    // write into log files 
}
