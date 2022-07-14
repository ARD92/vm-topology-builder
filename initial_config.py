import sys
import telnetlib
import argparse
import yaml
import time

parser = argparse.ArgumentParser()
parser.add_argument('-t', action='store', dest='TOPOLOGY', help='topology file for which diagram needs to be generated')
args = parser.parse_args()
HOST = "127.0.0.1"

def basicConfig(host, port, hostname, ip):
    tn = telnetlib.Telnet(host=host,port=port,timeout=20)
    print("***************** opened port {} *************************".format(port))
    config_set = ["delete chassis auto-image-upgrade \n",
                  "delete interfaces fxp0.0 family inet dhcp \n",
                  "set system services netconf ssh \n",
                  "set system services telnet \n",
                  "set system services ssh \n",
                  "set system services ssh root-login allow \n",
                  "set system scripts language python \n",
                  "set system host-name {} \n".format(hostname),
                  "set interfaces fxp0.0 family inet address {} \n".format(ip),
                  ]

    print("setting up root authentication")
    tn.write(b"root \n")
    tn.write(b"cli \n")
    time.sleep(2)
    tn.write(b"configure \n")
    time.sleep(2)
    tn.write(b"set system root-authentication plain-text-password \n")
    time.sleep(2)
    tn.write(b"juniper123\n")
    time.sleep(2)
    tn.write(b"juniper123\n")
    time.sleep(2)
    tn.write(b"commit \n")
    time.sleep(5)

    for list in config_set:
        print(list)
        tn.write(str.encode(list))
        time.sleep(2)

    tn.write(b"commit\n")
    print("completed adding config list")
    print("added the basic necessary config, you can now login to the box using ssh")
    tn.close()

def main():
    num = 0
    with open(args.TOPOLOGY, 'r') as f:
        mapp = yaml.safe_load(f)
    for i in mapp["network_nodes"]:
        num = num + 1
        hostname = i["name"]
        port = int(i["re_port"])
        ip = "192.167.1."+str(num)+"/24"
        basicConfig(HOST, port, hostname, ip)

if __name__ == "__main__" :
    main()
