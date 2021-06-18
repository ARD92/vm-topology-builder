"""
Description: Script to generate a drawthe.net yaml file
which can then be used on http://go.drawthe.net/ quickly 
to get the topology diagram or save the topology diagram as 
a JPEG/PNG image locally

References:
- https://github.com/cidrblock/drawthe.net
- https://davidban77.hashnode.dev/making-your-network-topology-come-to-virtual-life-with-drawthenet-and-gns3fy-ck9kjsujb004jzss19c5stx6t

THis script reads the drawthenet_boiler.yaml and populates it from the topology file used to spin up the network. The output will 
be stored as <topology name_drawthe_diagram>.yaml
"""
import os
import argparse
import yaml

parser = argparse.ArgumentParser()
parser.add_argument('-t', action='store', dest='TOPOLOGY', help='topology file for which diagram needs to be generated')
args = parser.parse_args()

"""
Generate drawthe.net yaml 
file from topology file 
interface name should be of format <ge/xe/et>-<0.0.0>-<devicename>
example: ge-0.0.0-vmx1
return [devA, devB], [[devA: portA, devB: PortB ], [],[]]
"""
def GenDraw(mapp):
    with open(args.TOPOLOGY, 'r') as r:
        mapp_topo = yaml.safe_load(r)

    icons = []
    connections = []
    
    for i in range(len(mapp_topo["network_nodes"])):
        icons.append(mapp_topo["network_nodes"][i]["name"])
        for j in range(len(mapp_topo["network_nodes"][i]["link"])):
            endpoint = []
            if mapp_topo["network_nodes"][i]["link"][j]["intf"].startswith("ge-" or "et-" or "xe-"):
                conn1 = str(mapp_topo["network_nodes"][i]["link"][j]["intf"][9:])+":"+str(mapp_topo["network_nodes"][i]["link"][j]["intf"][:8])
                endpoint.append(conn1)
                print(conn1)
            if mapp_topo["network_nodes"][i]["link"][j]["intf"].startswith("eth"):
                conn1 = str(mapp_topo["network_nodes"][i]["link"][j]["intf"][4:])+":"+str(mapp_topo["network_nodes"][i]["link"][j]["intf"][:3])
                print(conn1)
                endpoint.append(conn1)
            if mapp_topo["network_nodes"][i]["link"][j]["peerintf"].startswith("ge-" or "et-" or "xe-"):
                conn2 = str(mapp_topo["network_nodes"][i]["link"][j]["peerintf"][9:])+":"+str(mapp_topo["network_nodes"][i]["link"][j]["peerintf"][:8])
                endpoint.append(conn2)
                print(conn2)
            if mapp_topo["network_nodes"][i]["link"][j]["peerintf"].startswith("eth"):
                conn2 = str(mapp_topo["network_nodes"][i]["link"][j]["peerintf"][5:])+":"+str(mapp_topo["network_nodes"][i]["link"][j]["peerintf"][:3])
                print(conn2)
                endpoint.append(conn2)
            connections.append(endpoint)
    return icons, connections
 
def ModifyBoilerFile(icons, connections):
    #print(icons)
    #print(connections)
    print("WIP")
    with open('drawthenet_boiler.yaml','r') as f:
        mapp_boiler = yaml.safe_load(f)
    #print(mapp_boiler)
    # use the same name as topology file by deleting .yaml
    mapp_boiler["title"]["name"] = args.TOPOLOGY[:-5]
    icon_vals_list = []
    icon_vals = {}
    conn_vals = []
    for i in icons:
        #print(i)
        icon_vals_int = {}
        icon_vals_int["icon"]="router"
        icon_vals_int["iconFamily"]="cisco"
        icon_vals_int["fill"]="none"
        icon_vals_int["x"]=2
        icon_vals_int["y"]=8
        icon_vals[i] = icon_vals_int
        #icon_vals_list.append(icon_vals)
    mapp_boiler["icons"] = icon_vals
    for j in connections:
        con_int = {}
        con_int["endpoints"] = j
        conn_vals.append(con_int)
    mapp_boiler["connections"] = conn_vals
    
    with open('drawit_{}.yaml'.format(args.TOPOLOGY[:-5]),'w') as f:    
        f.write(yaml.safe_dump(mapp_boiler, default_flow_style=False))

def main():
    print("Generating yaml file for drawthe.net compatability and saving topology as a file") 
    with open(args.TOPOLOGY, 'r') as f:
        mapp = yaml.safe_load(f)
    icons, connections = GenDraw(mapp)
    ModifyBoilerFile(icons, connections)
    print("end the x and y axis manually ")
if __name__ == "__main__":
    main()
