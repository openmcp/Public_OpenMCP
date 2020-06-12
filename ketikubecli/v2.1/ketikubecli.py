# Server - nfs 
# Client Requirement - apt-get install nfs-common
import argparse
import os
import subprocess
import yaml
import pprint
from termcolor import colored
import errno

NfsServer="10.0.3.12"

policy_dict = {
    "resource_priority": 1,
    "affinity_analysis": 0,
    "service_failover": 0,
    "servuce_location_pinning": 1,
    "cluster_regist": 1,
    "cluster_delete": 0,
    "monitoring_resource_select": 1,
    "service_pod_delete": 0,
    "cluster_status_priority": 0,
    "failure_frequency_notification": 1,
    "object_management": 1,
    "replica_auto_scaling": 1,
    "service_load_balancing": 1,
    "geo_service_loadbalancing": 0,
    "multi_metrics_scaling": 0,
    "geo_distance_analysis": 0
}
def ketimkdir(path):
    try:
        os.mkdir(path)
    except OSError as exc:
        if exc.errno != errno.EEXIST:
            raise
        pass

def initMount():
    try:
        mntInfo = subprocess.check_output(["cat /proc/mounts | grep /mnt"], shell=True).decode('utf-8').split()
        os.system("umount -l /mnt")
    except Exception as e:
        pass





def policy(args):
    if args.command == "list":
        print(colored("*** OpenMCP Current Policy List ***", 'yellow'))
        i = 1
        for k, v in policy_dict.items():
            if v == 1 :
                result = "Enabled"
                c = "green"
            else:
                result = "Disabled"
                c = "red"

            print(str(i)+". "+colored(k,'blue')+" is "+ colored(result,c))
            i+=1

    elif args.command == "insert":
        pass

    elif args.command == "delete":
        pass

def cluster_unjoin(args):
    openmcpIP = subprocess.check_output(['hostname -I'], shell=True).decode('utf-8').split()[0]

    if args.command == "list":
        os.system("mount -t nfs " + NfsServer + ":/home/nfs /mnt")

        try:
            member_list = subprocess.check_output(['ls /mnt/openmcp/' + openmcpIP + '/members/unjoin'], shell=True).split()
        except Exception as e:
            print(colored("Failed", "red") + colored(" UnJoin List '", "yellow") + "' in " + colored("OpenMCP Master: " + openmcpIP, "blue"))
            print(colored("=> Not Yet Register OpenMCP.", "cyan"))
            print(colored(
                "=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : ketikubecli regist openmcp"))
            os.system("umount -l /mnt")
            return

        print(colored("OpenMCP '" + openmcpIP + "' Current UnJoin List", "yellow"))
        print("#\tCluster Name\tIP")

        for i, unjoin_memeber_ip in enumerate(member_list):
            with open("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + unjoin_memeber_ip + "/config/config", 'r') as stream:
                try:
                    cluster_data_dict = yaml.safe_load(stream)

                    clusters = cluster_data_dict["clusters"]
                    cluster_name = clusters[0]["name"]

                except yaml.YAMLError as exc:
                    print(exc)

            print(str(i + 1) + ".\t" + cluster_name + "\t" + unjoin_memeber_ip)

        os.system("umount -l /mnt")
        return

    elif args.command =="cluster":

        if not args.ip:
            print("Must have cluster ip")
            return

        memberIP = args.ip
        os.system("mount -t nfs "+NfsServer+":/home/nfs /mnt")

        if not os.path.exists("/mnt/openmcp/" + openmcpIP):
            print(colored("Failed", "red") + colored(" UnJoin Cluster '"+memberIP+"'", "yellow") + "' in " + colored("OpenMCP Master: " + openmcpIP, "blue"))
            print(colored("=> Not Yet Register OpenMCP.", "cyan"))
            print(colored("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : ketikubecli regist openmcp"))

            os.system("umount -l /mnt")
            return

        if not os.path.exists("/mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP):
            print(colored("Failed", "red") + colored(" UnJoin Cluster '" + memberIP + "'", "yellow") + "' in " + colored("OpenMCP Master: " + openmcpIP, "blue"))
            print(colored("=> '"+memberIP+"' is Not Joined Cluster in OpenMCP.", "cyan"))

            os.system("umount -l /mnt")
            return

        with open("/mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP + "/config/config", 'r') as stream:
            try:
                cluster_data_dict = yaml.safe_load(stream)

                contexts = cluster_data_dict["contexts"]
                context = contexts[0]
                #print(context)

                clusters = cluster_data_dict["clusters"]
                cluster = clusters[0]
                #print(cluster)

                users = cluster_data_dict["users"]
                user = users[0]
                #print(user)

            except yaml.YAMLError as exc:
                print(exc)

        with open("/root/.kube/config", 'r') as stream:
            try:
                openmcp_data_dict = yaml.safe_load(stream)
                #openmcp_data_dict['clusters'].append(cluster)
                #openmcp_data_dict['contexts'].append(context)
                #openmcp_data_dict['users'].append(user)
                for i, cluster in enumerate(openmcp_data_dict['clusters']):
                    if memberIP in cluster['cluster']['server']:
                        target_name = cluster['name']
                        break


                for j, context in enumerate(openmcp_data_dict['contexts']):
                    if target_name == context['context']['cluster']:
                        target_user = context['context']['user']
                        break

                for k, user in enumerate(openmcp_data_dict['users']):
                    if target_user == user['name']:
                        break

                os.system("kubefedctl unjoin "+target_name+" --cluster-context "+target_name+" --host-cluster-context openmcp --v=2")
                del openmcp_data_dict['clusters'][i]
                del openmcp_data_dict['contexts'][j]
                del openmcp_data_dict['users'][k]


            except yaml.YAMLError as exc:
                print(exc)

        with open("/root/.kube/config", 'w') as stream:
            yaml.dump(openmcp_data_dict, stream, default_flow_style=False)
            os.system("mv /mnt/openmcp/"+openmcpIP+"/members/join/"+memberIP+" /mnt/openmcp/"+openmcpIP+"/members/unjoin/"+memberIP)

        os.system("umount -l /mnt")
        result = subprocess.check_output(["kubectl -n kube-federation-system get kubefedclusters"], shell=True)
        print(result)

def cluster_join(args):
    openmcpIP = subprocess.check_output(['hostname -I'], shell=True).decode('utf-8').split()[0]
    if args.command == "list":
        os.system("mount -t nfs "+NfsServer+":/home/nfs /mnt")

        try:
            member_list = subprocess.check_output(['ls /mnt/openmcp/'+openmcpIP+'/members/join'], shell=True).split()
        except Exception as e:
            print(colored("Failed", "red") + colored(" Join List '", "yellow") + "' in " + colored("OpenMCP Master: " + openmcpIP, "blue"))
            print(colored("=> Not Yet Register OpenMCP.", "cyan"))
            print(colored("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : ketikubecli regist openmcp"))
            os.system("umount -l /mnt")
            return

        print(colored("OpenMCP '"+openmcpIP+"' Current Join List","yellow"))
        print("#\tCluster Name\tIP")

        for i, join_member_ip in enumerate(member_list):
            with open("/mnt/openmcp/"+openmcpIP+"/members/join/"+join_member_ip+"/config/config", 'r') as stream:
                try:
                    cluster_data_dict = yaml.safe_load(stream)

                    clusters = cluster_data_dict["clusters"]
                    cluster_name = clusters[0]["name"]

                except yaml.YAMLError as exc:
                    print(exc)

            print(str(i+1) +".\t"+ cluster_name+"\t"+join_member_ip)

        os.system("umount -l /mnt")
        return

    elif args.command == "cluster":
        if not args.ip:
            print("Must have cluster ip")
            return

        memberIP = args.ip
        os.system("mount -t nfs " + NfsServer + ":/home/nfs /mnt")

        if not os.path.exists("/mnt/openmcp/"+openmcpIP):
            print(colored("Failed", "red") + colored(" Join List '", "yellow") + "' in " + colored("OpenMCP Master: " + openmcpIP, "blue"))
            print(colored("=> Not Yet Register OpenMCP.", "cyan"))
            print(colored("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : ketikubecli regist openmcp"))

            os.system("umount -l /mnt")
            return

        if not os.path.exists("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP):
            print(colored("Failed", "red") + colored(" UnJoin Cluster '" + memberIP + "'", "yellow") + "' in " + colored("OpenMCP Master: " + openmcpIP, "blue"))
            print(colored("=> '"+memberIP+"' is Not Joinable Cluster in OpenMCP.", "cyan"))

            os.system("umount -l /mnt")
            return


        with open("/mnt/openmcp/"+openmcpIP+"/members/unjoin/"+memberIP+"/config/config", 'r') as stream:
            try:
                cluster_data_dict = yaml.safe_load(stream)

                contexts = cluster_data_dict["contexts"]
                context = contexts[0]
                #print(context)

                clusters = cluster_data_dict["clusters"]
                cluster = clusters[0]
                #print(cluster)

                users = cluster_data_dict["users"]
                user = users[0]
                #print(user)

            except yaml.YAMLError as exc:
                print(exc)


        with open("/root/.kube/config", 'r') as stream:
            try:
                openmcp_data_dict = yaml.safe_load(stream)
                openmcp_data_dict['clusters'].append(cluster)
                openmcp_data_dict['contexts'].append(context)
                openmcp_data_dict['users'].append(user)
                #pprint.pprint(openmcp_data_dict)

            except yaml.YAMLError as exc:
                print(exc)

        with open("/root/.kube/config", 'w') as stream:
            yaml.dump(openmcp_data_dict, stream, default_flow_style=False)
            os.system("mv /mnt/openmcp/"+openmcpIP+"/members/unjoin/"+memberIP+" /mnt/openmcp/"+openmcpIP+"/members/join/"+memberIP)

        os.system("umount -l /mnt")

        os.system("kubefedctl join "+cluster["name"]+" --cluster-context "+cluster["name"]+" --host-cluster-context openmcp --v=2")
        result = subprocess.check_output(["kubectl -n kube-federation-system get kubefedclusters"], shell=True)
        print(result)

def cluster_regist(args):


    if args.command == "openmcp":
        os.system("mount -t nfs " + NfsServer + ":/home/nfs/ /mnt")
        openmcpIP = subprocess.check_output(['hostname -I'], shell=True).decode('utf-8').split()[0]

        if os.path.exists("/mnt/openmcp/"+openmcpIP):
            print(colored("Failed", "red") + colored(" Regist '", "yellow") + colored(" OpenMCP Master", "blue"))
            print(colored("=> Already Registered OpenMCP.", "cyan") + colored(":"+openmcpIP, "yellow"))

            os.system("umount -l /mnt")
            return


        ketimkdir("/mnt/openmcp")
        ketimkdir("/mnt/openmcp/" + openmcpIP)
        ketimkdir("/mnt/openmcp/" + openmcpIP + "/master")
        ketimkdir("/mnt/openmcp/" + openmcpIP + "/master/config")
        ketimkdir("/mnt/openmcp/" + openmcpIP + "/master/pki")
        ketimkdir("/mnt/openmcp/" + openmcpIP + "/members")
        ketimkdir("/mnt/openmcp/" + openmcpIP + "/members/join")
        ketimkdir("/mnt/openmcp/" + openmcpIP + "/members/unjoin")

        os.system("cp ~/.kube/config /mnt/openmcp/" + openmcpIP + "/master/config/config")
        os.system("cp /etc/kubernetes/pki/etcd/ca.crt /mnt/openmcp/" + openmcpIP + "/master/pki/ca.crt")
        os.system("cp /etc/kubernetes/pki/etcd/server.crt /mnt/openmcp/" + openmcpIP + "/master/pki/server.crt")
        os.system("cp /etc/kubernetes/pki/etcd/server.key /mnt/openmcp/" + openmcpIP + "/master/pki/server.key")

        # SSH Public Key Copy
        os.system("cat /mnt/ssh/id_rsa.pub >> /root/.ssh/authorized_keys")

        print(colored("Success", "green") + colored(" OpenMCP Master Regist '" + openmcpIP, "yellow"))

        os.system("umount -l /mnt")
        return



    elif args.command == "member":
        os.system("mount -t nfs " + NfsServer + ":/home/nfs/ /mnt")
        if not args.ip:
            print("Must have cluster ip")
            os.system("umount -l /mnt")
            return


        openmcpIP = args.ip
        memberIP = subprocess.check_output(['hostname -I'], shell=True).decode('utf-8').split()[0]

        if not os.path.exists("/mnt/openmcp/" + openmcpIP + "/master"):
            print(colored("Failed", "red") + colored(" Regist '" + memberIP, "yellow") + "' in " + colored("OpenMCP Master: " + openmcpIP, "blue"))
            print(colored("=> Not Yet Register OpenMCP.", "cyan"))
            print(colored("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : ketikubecli regist openmcp"))
            os.system("umount -l /mnt")
            return

        if memberIP == openmcpIP:
            print(colored("Failed","red")+ colored(" Regist '"+memberIP,"yellow")+"' in "+colored("OpenMCP Master: "+openmcpIP, "blue"))
            print(colored("=> Can Not Self regist. [My_IP '"+memberIP+"', OpenMCP_IP '"+openmcpIP+"']","cyan"))
            os.system("umount -l /mnt")
            return

        # Already Regist
        if os.path.exists("/mnt/openmcp/"+openmcpIP+"/members/unjoin/"+memberIP) :
            print(colored("Failed","red")+ colored(" Regist '"+memberIP,"yellow")+"' in "+colored("OpenMCP Master: "+openmcpIP, "blue"))
            print(colored("=> Already Regist","cyan"))
            os.system("umount -l /mnt")
            return

        elif os.path.exists("/mnt/openmcp/"+openmcpIP+"/members/join/"+memberIP):
            print(colored("Failed","red")+ colored(" Regist '"+memberIP,"yellow")+"' in "+colored("OpenMCP Master: "+openmcpIP, "blue"))
            print(colored("=> Already Joined by OpenMCP '"+openmcpIP+"'","cyan"))
            os.system("umount -l /mnt")
            return

        else:
            ketimkdir("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP)
            ketimkdir("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/config")
            ketimkdir("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki")

            os.system("cp ~/.kube/config /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/config/config")
            os.system("cp /etc/kubernetes/pki/etcd/ca.crt /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki/ca.crt")
            os.system("cp /etc/kubernetes/pki/etcd/server.crt /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki/server.crt")
            os.system("cp /etc/kubernetes/pki/etcd/server.key /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki/server.key")

            # SSH Public Key Copy
            os.system("cat /mnt/ssh/id_rsa.pub >> /root/.ssh/authorized_keys")

            print(colored("Success","green")+ colored(" Regist '"+memberIP,"yellow")+"' in "+colored("OpenMCP Master: "+openmcpIP, "blue"))
            os.system("umount -l /mnt")
            return




def main():
    initMount()

    parser = argparse.ArgumentParser(prog='ketikubecli')
    #parser.add_argument('command', help='foo help')

    subparsers = parser.add_subparsers(help='sub-command help')

    parser_a = subparsers.add_parser('join', help='join help')
    parser_a.set_defaults(func=cluster_join)

    parser_a.add_argument('command', choices=['cluster','list'], help='join kind help')
    parser_a.add_argument('--ip', required=False, help='IP Address')


    parser_b = subparsers.add_parser('regist', help='regist help')
    parser_b.set_defaults(func=cluster_regist)

    parser_b.add_argument('command', choices=['openmcp', 'member'], help='join kind help')
    parser_b.add_argument('--ip', required=False, help='ip help')

    parser_c = subparsers.add_parser('unjoin', help='unjoin help')
    parser_c.set_defaults(func=cluster_unjoin)

    parser_c.add_argument('command', choices=['cluster', 'list'], help='unjoin kind help')
    parser_c.add_argument('--ip', required=False, help='IP address')

    parser_d = subparsers.add_parser('policy', help='policy help')
    parser_d.set_defaults(func=policy)

    parser_d.add_argument('command', choices=['insert','delete','list'], help='policy kind help')

    args = parser.parse_args()
    args.func(args)
    #if args.cluster:
        #cluster_join(args)
        #print(parser_a.parse_args())

if __name__=="__main__":
    main()
