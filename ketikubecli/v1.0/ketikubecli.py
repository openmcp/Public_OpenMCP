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
    myip = subprocess.check_output(['hostname -I'], shell=True).decode('utf-8').split()[0]
    if args.command == "list":
        os.system("mount -t nfs "+NfsServer+":/home/nfs /mnt")
        config_list = subprocess.check_output(['ls /mnt/kubeconfig/'+myip+'/unjoin'], shell=True).split()
        ketimkdir("/mnt/kubeconfig")
        ketimkdir("/mnt/kubeconfig/"+myip)
        ketimkdir("/mnt/kubeconfig/"+myip+"/unjoin")
        ketimkdir("/mnt/kubeconfig/"+myip+"/join")

        print(colored("OpenMCP '"+myip+"' Current UnJoin List","yellow"))
        print("#\tCluster Name\tIP")
        for i, config in enumerate(config_list):
            with open("/mnt/kubeconfig/"+myip+"/unjoin/"+config, 'r') as stream:
                try:
                    cluster_data_dict = yaml.safe_load(stream)

                    clusters = cluster_data_dict["clusters"]
                    cluster_name = clusters[0]["name"]

                except yaml.YAMLError as exc:
                    print(exc)

            unjoin_ip = config.replace("config_", "")
            print(str(i+1) +".\t"+ cluster_name+"\t"+unjoin_ip)


        os.system("umount -l /mnt")

    elif args.command =="cluster":
		if not args.ip:
			print("Must have cluster ip")
			return

		os.system("mount -t nfs "+NfsServer+":/home/nfs /mnt")

		ketimkdir("/mnt/kubeconfig")
                ketimkdir("/mnt/kubeconfig/"+myip)
		ketimkdir("/mnt/kubeconfig/"+myip+"/unjoin")
		ketimkdir("/mnt/kubeconfig/"+myip+"/join")


		with open("/mnt/kubeconfig/"+myip+"/join/config_"+args.ip, 'r') as stream:
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
					if args.ip in cluster['cluster']['server']:
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
			os.system("mv /mnt/kubeconfig/"+myip+"/join/config_"+args.ip+" /mnt/kubeconfig/"+myip+"/unjoin/config_"+args.ip)

		os.system("umount -l /mnt")
		result = subprocess.check_output(["kubectl -n kube-federation-system get kubefedclusters"], shell=True)
		print(result)



def cluster_join(args):
	myip = subprocess.check_output(['hostname -I'], shell=True).decode('utf-8').split()[0]
	if args.command == "list":
		os.system("mount -t nfs "+NfsServer+":/home/nfs /mnt")

                ketimkdir("/mnt/kubeconfig")
		ketimkdir("/mnt/kubeconfig/"+myip)
		ketimkdir("/mnt/kubeconfig/"+myip+"/unjoin")
		ketimkdir("/mnt/kubeconfig/"+myip+"/join")

		config_list = subprocess.check_output(['ls /mnt/kubeconfig/'+myip+'/join'], shell=True).split()
		print(colored("OpenMCP '"+myip+"' Current Join List","yellow"))
		print("#\tCluster Name\tIP")

		for i, config in enumerate(config_list):
			with open("/mnt/kubeconfig/"+myip+"/join/"+config, 'r') as stream:
				try:
					cluster_data_dict = yaml.safe_load(stream)

					clusters = cluster_data_dict["clusters"]
					cluster_name = clusters[0]["name"]

				except yaml.YAMLError as exc:
					print(exc)

			join_ip = config.replace("config_", "")
			print(str(i+1) +".\t"+ cluster_name+"\t"+join_ip)
		os.system("umount -l /mnt")

	elif args.command == "cluster":
		if not args.ip:
			print("Must have cluster ip")
			return

		os.system("mount -t nfs "+NfsServer+":/home/nfs /mnt")

                ketimkdir("/mnt/kubeconfig")
		ketimkdir("/mnt/kubeconfig/"+myip)
		ketimkdir("/mnt/kubeconfig/"+myip+"/unjoin")
		ketimkdir("/mnt/kubeconfig/"+myip+"/join")
	
		with open("/mnt/kubeconfig/"+myip+"/unjoin/config_"+args.ip, 'r') as stream:
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
			os.system("mv /mnt/kubeconfig/"+myip+"/unjoin/config_"+args.ip+" /mnt/kubeconfig/"+myip+"/join/config_"+args.ip)

		os.system("umount -l /mnt")
	
		os.system("kubefedctl join "+cluster["name"]+" --cluster-context "+cluster["name"]+" --host-cluster-context openmcp --v=2")
		result = subprocess.check_output(["kubectl -n kube-federation-system get kubefedclusters"], shell=True)
		print(result)



def cluster_regist(args):
	myip = subprocess.check_output(['hostname -I'], shell=True).decode('utf-8').split()[0]
	openmcp_ip = args.ip
	if myip == openmcp_ip:
		print(colored("Failed","red")+ colored(" Regist '"+myip,"yellow")+"' in "+colored("OpenMCP Master: "+openmcp_ip, "blue"))
		print(colored("=> Can Not Self regist. [My_IP '"+myip+"', OpenMCP_IP '"+openmcp_ip+"']","cyan"))
		return

	os.system("mount -t nfs "+NfsServer+":/home/nfs/ /mnt")

        ketimkdir("/mnt/kubeconfig")
	ketimkdir("/mnt/kubeconfig/"+openmcp_ip)
	ketimkdir("/mnt/kubeconfig/"+openmcp_ip+"/unjoin")
	ketimkdir("/mnt/kubeconfig/"+openmcp_ip+"/join")

	# Already Regist
	if os.path.exists("/mnt/kubeconfig/"+openmcp_ip+"/unjoin/config_"+myip):
		print(colored("Failed","red")+ colored(" Regist '"+myip,"yellow")+"' in "+colored("OpenMCP Master: "+openmcp_ip, "blue"))
		print(colored("=> Already Regist","cyan"))

	elif os.path.exists("/mnt/kubeconfig/"+openmcp_ip+"/join/config_"+myip):
		print(colored("Failed","red")+ colored(" Regist '"+myip,"yellow")+"' in "+colored("OpenMCP Master: "+openmcp_ip, "blue"))
		print(colored("=> Already Joined by OpenMCP '"+openmcp_ip+"'","cyan"))

	else:
		os.system("cp ~/.kube/config /mnt/kubeconfig/"+openmcp_ip+"/unjoin/config_"+myip)

		print(colored("Success","green")+ colored(" Regist '"+myip,"yellow")+"' in "+colored("OpenMCP Master: "+openmcp_ip, "blue"))
		os.system("umount -l /mnt")

def main():
    parser = argparse.ArgumentParser(prog='ketikubecli')
    #parser.add_argument('command', help='foo help')

    subparsers = parser.add_subparsers(help='sub-command help')

    parser_a = subparsers.add_parser('join', help='join help')
    parser_a.set_defaults(func=cluster_join)

    parser_a.add_argument('command', choices=['cluster','list'], help='join kind help')
    parser_a.add_argument('--ip', required=False, help='IP Address')

    parser_b = subparsers.add_parser('regist', help='regist help')
    parser_b.set_defaults(func=cluster_regist)

    parser_b.add_argument('--ip', required=True, help='ip help')

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
