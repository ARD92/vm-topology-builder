#  Description   : script to gather report from the box

#!/usr/bin/python
import os
import argparse
import yaml
from lxml import etree
from jnpr.junos import Device
from jnpr.junos.utils.config import Config
from jnpr.junos.exception import ConfigLoadError, CommitError, ConnectError, RpcError, ConnectAuthError, ConnectRefusedError, ConnectTimeoutError


yaml_open = open('report.yaml')
map = yaml.load(yaml_open)

def open_dev(dev):
    try:
        dev.open(gather_facts=True)
        print ("device opened")
        model = dev.facts['model']
        print(model)

    except ConnectAuthError:
        print("failed: authentication incorrect")
    except ConnectRefusedError:
        print("failed: connection refused")
    except ConnectTimeoutError:
        print("failed: connection timed out")
    except ConnectError:
        print("failed: connect error")

def gen_report(dev,hst):
    print(hst)
    with open('info_'+str(hst)+'.txt','w') as f:
        for op in map['opcommand']:
            f.write(20*"*"+" "+op+" "+20*"*")
            f.write('\n')
            f.write(dev.cli(op))
            #xml_rpc = dev.display_xml_rpc(op)
            #try:
            #    f.write(dev.execute(xml_rpc))
            #except RpcError:
            #    print("RPC Error")   
            f.write("\n")
            f.write("\n")
        f.close()

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('-user', required=True)
    parser.add_argument('-pwd', required=True)
    args = parser.parse_args()

    for hst in map['hosts']:
        dev = Device(host=hst, user=args.user, passwd=args.pwd,  normalize=True)
        open_dev(dev)
        gen_report(dev,hst)

if __name__ == "__main__":
    main()
