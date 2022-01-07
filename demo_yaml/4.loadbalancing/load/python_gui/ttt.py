import multiprocessing
import sys,time
from multiprocessing import Process, Event, Queue
from threading import Thread
from expressvpn import wrapper
from PySide2.QtCore import QCoreApplication, QSize, Qt, Signal, QThread, QProcess
import PySide2.QtGui, PySide2.QtWidgets as QtWidgets
from PySide2.QtWidgets import *

from wrapper import *
import requests
import os
import json
import socket

from Configure.configure import *


from MyObject.MyButton import *
from MyObject.MyComboBox import *
from MyObject.MyLabel import *
from MyObject.MyTextEdit import *


#ClusterList = {"1" : "10.0.0.1", "2" : "10.0.0.2", "3" : "10.0.0.3", "4" : "10.0.0.4"} # 클러스터 리스트
CountryList  = {"USA - Los Angeles - 1" : "usla1", "Japan - Tokyo" : "jpto", "South Korea - 2" : "kr2"} # 나라 리스트

class SysbenchInfo():
    def __init__(self, clusterName, svcIP, svcPort) -> None:
        self.clusterName = clusterName
        self.svcIP = svcIP
        self.svcPort = svcPort
        self.loadLevel = 0
        

class Form(QMainWindow):
    def __init__(self):
        super().__init__()
        self.initVariable()
        self.setupDefaultUI()
        vpnlist()
        self.buttonA_Init()
        
        
        
        
        
    def moveEvent(self, event):
        # print("check")
        Global.winX = event.pos().x()
        Global.winY = event.pos().y()
    
    def initVariable(self):
        self.processlist = [0,0,0,0]
        # self.eventA = Event()
        self.eventC = Event()
        
    def setupDefaultUI(self):
        self.setFixedSize(Global.winWidth,Global.winHeight)
        self.setWindowTitle("OpenMCP LoadRunner")
        
        
        
        self.MainWidget = QWidget(self)
        self.setCentralWidget(self.MainWidget)
        
        self.grid = QtWidgets.QGridLayout()
        self.MainWidget.setLayout(self.grid)
        
        
        
        # Up Left Widget 정의
        self.upleftWidget= QWidget()
        self.upleftLayout = QGridLayout()
        self.upleftWidget.setLayout(self.upleftLayout)
        
        self.amessage = MyLabel('External Client Connector')
        self.amessage.setAlignment(Qt.AlignCenter)
       
        self.alabel = MyLabel('국가:',3)
        self.abuttonmessage = MyTextEdit("[Status] : <font color=\"red\">DisConnected</font>")
        self.abuttonstart = MyButton('Connect',4)
        self.acombobox = MyComboBox(2)

        self.acombobox.addItems(CountryList)


        self.abuttonstart.clicked.connect(self.buttonA_Start_event)
    
        self.upleftLayout.addWidget(self.amessage,0,0,1,3)
        self.upleftLayout.addWidget(self.alabel,1,0,1,1)
        self.upleftLayout.addWidget(self.acombobox,1,1,1,1)
        self.upleftLayout.addWidget(self.abuttonstart,1,2,1,1)
        self.upleftLayout.addWidget(self.abuttonmessage,2,0,4,3)
        
        
      

        
        
        # Up Right Widget 정의
        self.uprightWidget= QWidget()
        self.uprightLayout = QGridLayout()
        self.uprightWidget.setLayout(self.uprightLayout)
        
        self.cmessage = MyLabel('Traffic Sender')
        self.cmessage.setAlignment(Qt.AlignCenter)
        # self.cmessage.setStyleSheet("background-color : black; color : white; font-size: 14pt; font-weight: bold;")
        # self.cmessage.setAlignment(Qt.AlignCenter)
        self.clabel = MyLabel('동시 접속자 수:',3)
        self.cbuttonmessage = MyTextEdit("[Status] : <font color=\"red\">Disable</font>")
        self.cbuttonstart = MyButton('Send',4)
        self.ccombobox = MyComboBox(2)
        self.ccombobox.addItem("1000")
        self.ccombobox.addItem("5000")
        self.ccombobox.addItem("10000")
        self.ccombobox.addItem("50000")
        self.ccombobox.addItem("100000")
        
        #self.setStyle(self.cbuttonstart,1)
        # self.setStyle(self.abuttonstop,2)
        #self.setStyle(self.cmessage,6)
        #self.setStyle(self.cbuttonmessage,7)
        #self.setStyle(self.ccombobox,10)
        self.cbuttonstart.clicked.connect(self.buttonC_Start_event)
        # self.cbuttonstop.clicked.connect(self.buttonC_Stop_event)
        
        self.uprightLayout.addWidget(self.cmessage,0,0,1,3)
        
        self.uprightLayout.addWidget(self.clabel,1,0,1,1)
        self.uprightLayout.addWidget(self.ccombobox,1,1,1,1)
        self.uprightLayout.addWidget(self.cbuttonstart,1,2,1,1)
        
        self.uprightLayout.addWidget(self.cbuttonmessage,2,0,4,3)
       
        # upleftlayout.addWidget(abuttonstop,1,1)
        
        
        # Down Left Widget 정의
        self.downleftWidget = QWidget()
        self.downleftLayout = QGridLayout()
        self.downleftWidget.setLayout(self.downleftLayout)
        
        self.bmessage = MyLabel('Load Generator')
        # self.bmessage.setStyleSheet("background-color : black; color : white; font-size: 14pt; font-weight: bold;")
        self.bmessage.setAlignment(Qt.AlignCenter)
        self.bbuttonmessage = MyTextEdit("")
   
        self.m_label_gif = QLabel()
        self.m_label_gif.setText("aaa")
        self.m_movie_gif = QtGui.QMovie("resources/image/Spin-1s-200px.gif")
        self.m_label_gif.setMovie(self.m_movie_gif)
        self.m_label_gif.setSizePolicy(QSizePolicy.Ignored, QSizePolicy.Ignored )
        self.m_label_gif.setScaledContents(True)
        self.m_gif_flag = False
        
        #Global.signal.stopBbtnCompleteSignal.connect(self.closeUI)
        
        self.bbuttonstart = MyButton('Load', 2)
        self.bbuttonstop = MyButton('Stop', 3)
        
        #self.setStyle(self.bbuttonstart,1)
        #self.setStyle(self.bbuttonstop,2)
        # self.setStyle(self.bbuttonmessage,7)
        # self.bbuttonstart.clicked.connect(self.buttonB_Start_event)
        # self.bbuttonstop.clicked.connect(self.buttonB_Stop_event)
        self.bcombobox = MyComboBox()
        self.sysbenchInfoList = self.getSvcSysbenchs()
        self.BClusterList = [Event() for i in range (0, len(self.sysbenchInfoList))]
        
        
        
        for i in range(0, len(self.sysbenchInfoList)):
            self.bcombobox.addItem(self.sysbenchInfoList[i].clusterName)
        
        self.getLoadGenStatusAndSetText()
        # for i in ClusterList:
        #     self.bcombobox.addItem(i)
        # self.bcombobox.addItems(ClusterList)
        self.bcombobox2 = MyComboBox()
        
        
        self.bcombobox2.addItem("1")
        self.bcombobox2.addItem("2")
        self.bcombobox2.addItem("3")
        self.bcombobox2.addItem("4")
        self.bcombobox2.addItem("5")
        self.bcombobox2.addItem("6")
        # self.bcombobox2.addItem("Max")
        self.blabel = MyLabel('Cluster:',3)
        self.blabel2 = MyLabel('Load Level:',3)
        self.bbuttonstart.clicked.connect(self.buttonB_Start_event)
        self.bbuttonstop.clicked.connect(self.buttonB_Stop_event)
        # self.bbuttonicon = PySide2.QtGui.QPixmap("assets/button3.png")
        # self.bicon = PySide2.QtGui.QIcon(self.bbuttonicon)
        # self.bbuttonstart.setIcon(self.bicon)
        # self.bbuttonstart.setIconSize(self.bbuttonicon.rect().size())
        # self.bbuttonstart.setIconSize(QSize(200,200))
        # self.bbuttonstart.setFixedSize(self.bbuttonicon.rect().size())
        #self.setStyle(self.bbuttonstart,3)
        #self.setStyle(self.bbuttonstop,4)
        #self.setStyle(self.bbuttonmessage,7)
        #self.setStyle(self.bmessage,6)
        #self.setStyle(self.bcombobox,10)
        #self.setStyle(self.bcombobox2,10)
        
        # self.downleftLayout.addWidget(self.bbuttonmessage,0,0,1,2)
        # self.downleftLayout.addWidget(self.bbuttonstart,1,0)
        # self.downleftLayout.addWidget(self.bbuttonstop,1,1)
            
        self.downleftLayout.addWidget(self.bmessage,0,0,1,6)
        self.downleftLayout.addWidget(self.blabel,1,0,1,1)
        self.downleftLayout.addWidget(self.bcombobox,1,1,1,1)
        self.downleftLayout.addWidget(self.blabel2,1,2,1,1)
        self.downleftLayout.addWidget(self.bcombobox2,1,3,1,1)
        
        self.downleftLayout.addWidget(self.bbuttonstart,1,4,1,1)
        self.downleftLayout.addWidget(self.bbuttonstop,1,5,1,1)
        
        self.downleftLayout.addWidget(self.bbuttonmessage,2,0,4,5)
        self.downleftLayout.addWidget(self.m_label_gif,2,5,4,1)
        
            
        
        self.upleftWidget.setSizePolicy(
                QSizePolicy.Expanding,QSizePolicy.Expanding
            )
        self.uprightWidget.setSizePolicy(
                QSizePolicy.Expanding,QSizePolicy.Expanding
            )
        self.downleftWidget.setSizePolicy(
                QSizePolicy.Expanding,QSizePolicy.Expanding
            )
       
        # Main Layout에 각 Sub Widget들 배치
        self.grid.addWidget(self.upleftWidget,0,0,1,1)
        self.grid.addWidget(self.uprightWidget,0,1,1,1)
        self.grid.addWidget(self.downleftWidget,1,0,1,2)
        
       
       
   
    def getToken(self):
        #run_command("echo -n | openssl s_client -connect openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org:8080 | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > server.crt")
        
        
        headers = {
            'Content-type': 'application/json',
        }

        data = '{"username":"openmcp","password":"keti"}'

        response = requests.post('https://openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org:8080/token', headers=headers, data=data, verify=False)

        #response = requests.post(tokenUrl, headers=headers, data=dict_data, verify=False)
        print(response.status_code)
        response.raise_for_status()
       
        TOKEN = json.loads(response.text).get('token')
        return TOKEN
    
    def getAllClusterList(self):
        clusterList = []
        TOKEN = self.getToken()
        headers = {
            'Authorization': 'Bearer '+TOKEN,
        }

        params = (
            ('clustername', 'openmcp'),
        )
        
        url = "https://openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org:8080/apis/core.kubefed.io/v1beta1/namespaces/kube-federation-system/kubefedclusters"
        response = requests.get(url, headers=headers, params=params, verify=False)
        print(response.status_code)
        response.raise_for_status()
        loaded = json.loads(response.text)
        
        for cluster in loaded.get('items'):
            clusterName = cluster.get('metadata').get('name')
            print("clusterName:", clusterName)
            clusterList.append(clusterName)
        
        return clusterList
    
    def getSvcSysbenchs(self):
        allClusterList = self.getAllClusterList()
        
        sysbenchInfoList = []     
          
        TOKEN = self.getToken()
        headers = {
            'Authorization': 'Bearer '+TOKEN,
        }

        queue = Queue()
        procs = []
        for clusterName in allClusterList:
            params = (
                ('clustername', clusterName),
            )
            proc = Process(target=self.getSvcSysbench,args=(clusterName, headers, params, queue))
            procs.append(proc)
            proc.start()
            
        for proc in procs:
            proc.join()

        print("join complete", queue.qsize())
        for i in range(queue.qsize()):
            print("join complete" ,i)
            try:
                sysbenchInfoList.append(queue.get())
            except:
                print(queue.qsize())
            
        print("get complete" )
        return sysbenchInfoList
        
    
    def getSvcSysbench(self, clusterName, headers, params, queue):
        url = "https://openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org:8080/api/v1/namespaces/"+Global.Load_Namespace+"/services/"+Global.Load_Svc_Name
        try:
            response = requests.get(url, headers=headers, params=params, verify=False, timeout=5)
        except requests.exceptions.Timeout:
            print(clusterName+' Timeout')
            return
        except Exception as e:
            print(clusterName+ 'except', e)
            return
        
        print(clusterName, response.status_code)
        response.raise_for_status()
        
        if response.text != "":
            #print(response.text)
            loaded = json.loads(response.text)
            svcIP = loaded.get('status').get('loadBalancer').get('ingress')[0].get('ip')
            svcPort = ""
            for port in loaded.get('spec').get('ports'):
                if port["name"] =="http-sysbench":
                    svcPort = port["port"]
                    break
            
            print(clusterName, svcIP, svcPort)
            p = SysbenchInfo(clusterName=clusterName,svcIP=svcIP,svcPort=svcPort)
            print("q put")
            queue.put(p)
        return
           

        
        
    def buttonA_Init(self):
        if getStatus():
            
            
            out = run_command(VPN_STATUS)
            contry = ""
            for item in out:
                if "Connected to " in item:
                    
                    contry = item.split("Connected to")[1]
                    break
            
            
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.connect(("pwnbit.kr",443))
            self.abuttonmessage.setText("[Status] : <font color=\"green\">Connected</font><br>[Contry] : "+contry+"<br>[IPaddr] : "+sock.getsockname()[0])
            self.abuttonstart.setText('DisConnect')
            self.abuttonstart.setStyle(5)
        else:
            self.abuttonmessage.setText("[Status] : <font color=\"red\">DisConnected</font>")
            self.abuttonstart.setText('Connect')
            self.abuttonstart.setStyle(4)
            
    
    def buttonA_func_Connect(self):
        # while True:
        print("A Connect func Text is :",self.acombobox.currentText())
        alias = CountryList[self.acombobox.currentText()]
        connect_alias(alias)
            # if eventA.is_set():
        #     break

    def buttonA_func_DisConnect(self):
    # while True:
        # print("A Disconnet func")
        time.sleep(1)
        disconnect()
        print("disconnected")
        # if eventA.is_set():
        #     break

    def buttonA_Start_event(self): # 버튼 A 시작 이벤트 처리
    # global eventA
        
        if getStatus():
            self.abuttonstart.setText("Connect")
            self.abuttonstart.setStyle(4)
            self.buttonA_func_DisConnect()
            self.abuttonmessage.setText("[Status] : <font color=\"red\">DisConnected</font>")
            self.processlist[0] = 0
        else:
            self.abuttonstart.setText("DisConnect")
            self.abuttonstart.setStyle(5)
            # if eventA.is_set():
            #     eventA = Event()
            self.buttonA_func_Connect()
            contry = self.acombobox.currentText()
            
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.connect(("pwnbit.kr",443))
            self.abuttonmessage.setText("[Status] : <font color=\"green\">Connected</font><br>[Contry] : "+contry+"\n[IPaddr] : "+sock.getsockname()[0])
            # proc = Process(target=buttonA_func,args=(even))
            # proc.start()
            self.processlist[0] = 1
    
    def requestTraffic(self, i, reqPerTh):
        url = Global.Traffic_URL
        
        for n in range (0, reqPerTh):
            reqNum = (reqPerTh*i)+n
            print("Http Get(",reqNum,") Request Start")
            response = requests.get(url, verify=False, allow_redirects=False, stream=True)
            print("Status code (",reqNum,"):", response.status_code)
            
    
    def buttonC_func_SendTraffic(self, connectedNum):
        
        while True:
            ths = []
            reqPerTh = 100
            thNum = int(connectedNum/reqPerTh)
            
            
            for i in range(0,thNum):
                th = Thread(target=self.requestTraffic, args=(i,reqPerTh))
                #proc = Process(target=self.requestTraffic, args=(i,))
                #proc.daemon = True
                ths.append(th)
                th.start()
                
            for th in ths:
                th.join()

            # if self.eventC.is_set():
            #     print("Traffic Ended")
            #     break
            print("Traffic wait for "+str(Global.Traffic_Period)+"s")
            time.sleep(Global.Traffic_Period)
           

    def buttonC_func_DisConnect(self):
        #self.eventC.set()
        self.procTraffic.kill()
        print("C Disconnet func")
        print("Traffic Ended")

    def buttonC_Start_event(self): # 버튼 C 시작 이벤트 처리
        if self.processlist[2] != 0:
            self.cbuttonstart.setStyle(4)
            #self.setStyle(self.cbuttonstart,1)
            self.buttonC_func_DisConnect()
            
            self.cbuttonstart.setText("Send")
            self.cbuttonmessage.setText("[Status] : <font color=\"red\">Disable</font>")
            self.processlist[2] = 0
        else:
            #self.setStyle(self.cbuttonstart,2)
            self.cbuttonstart.setStyle(5)
            # if self.eventC.is_set():
            #     self.eventC = Event()
            
            connectedNum = int(self.ccombobox.currentText())
            print("C Connect func Text is :", connectedNum)
            
            
            self.procTraffic = Process(target=self.buttonC_func_SendTraffic, args=(connectedNum,))
            self.procTraffic.daemon = True
            self.procTraffic.start()
            self.cbuttonstart.setText("Stop")
            self.cbuttonmessage.setText("[Status] : <font color=\"green\">Enable</font><br>[CCU] : "+str(connectedNum)+"<br>[Period] : "+str(Global.Traffic_Period)+"s<br>[Target] : "+Global.Traffic_URL)
            # proc = Process(target=buttonA_func,args=(even))
            # proc.start()
            self.processlist[2] = 1
        

    def getLoadGenStatusAndSetText(self):
        for i, sysinfo in enumerate(self.sysbenchInfoList):
            response = requests.get("http://"+sysinfo.svcIP+":"+str(sysinfo.svcPort)+"/status")
            response.raise_for_status()

            if response.text != "":
                
 
                maxlevel = 0
                print(response.text)
                findflag = False
                substring = " Status : START / Level : "
                for item in response.text.split("\n"):
                    if substring in item:
                        findflag = True
                        level = int(item.split(substring)[1])
                        maxlevel = max(maxlevel, level)
                
                if findflag:
                    j = self.bcombobox.findText(sysinfo.clusterName)
                    self.sysbenchInfoList[j].loadLevel = maxlevel
                    self.BClusterList[j].set()
                    
        self.setGenText()
    
    def setGenText(self):
        Bmessegetext = ""
        for i in range (len(self.BClusterList)) :
            if self.BClusterList[i].is_set() == True :
                Bmessegetext = Bmessegetext + "["+ self.sysbenchInfoList[i].clusterName + "] : <font color=\"green\">ON</font>, [Load Level] :" + str(self.sysbenchInfoList[i].loadLevel) + "<br>"
        self.bbuttonmessage.setText(Bmessegetext)
               
               
    def buttonB_Start_event(self):
        if self.m_gif_flag == True:
            return
        if self.BClusterList[self.bcombobox.currentIndex()].is_set() :
            print("Cluster '" + self.bcombobox.currentText() + "' is Running")
            return
        self.BClusterList[self.bcombobox.currentIndex()].set()
        
        
        for i, sysinfo in enumerate(self.sysbenchInfoList):
            if self.sysbenchInfoList[i].clusterName == self.bcombobox.currentText():
                self.sysbenchInfoList[i].loadLevel = self.bcombobox2.currentText()

                break
            
        # if self.bcombobox.currentText.eventB.is_set():
        #     self.bcombobox.currentText.eventB = Event()
        # self.bbuttonmessage.setText("Start")
        params = (
                ('v', str(sysinfo.loadLevel)),
            )
        
        requests.get("http://"+sysinfo.svcIP+":"+str(sysinfo.svcPort)+"/cpu/start", params=params)
        requests.get("http://"+sysinfo.svcIP+":"+str(sysinfo.svcPort)+"/memory/start", params=params)
        
        # self.bcombobox.currentText.processlist[1] = 1
        self.setGenText()

       
    
    def hideLoadingUI(self):

        self.m_movie_gif.stop()
        self.m_label_gif.hide()
        self.setGenText()
        self.m_gif_flag = False
        
    
    def buttonB_Stop_event(self):
        if self.m_gif_flag == True:
            return
        self.m_gif_flag = True
        # self.bcombobox.currentText.eventB.set()
        # self.bcombobox.currentText.processlist[1] = 0
        self.BClusterList[self.bcombobox.currentIndex()].clear()
        #self.bbuttonmessage.setText("Stop")
        
        for i, sysinfo in enumerate(self.sysbenchInfoList):
            if self.sysbenchInfoList[i].clusterName == self.bcombobox.currentText():
                break
        
        
        self.m_movie_gif.start()
        self.m_label_gif.show()

        self.th = Task(sysinfo)
        Global.signal.stopBbtnCompleteSignal.connect(self.hideLoadingUI)
        self.th.start()


        

class Task(QThread):
    def __init__(self, sysinfo):
        super().__init__()

        self.sysinfo = sysinfo

        
    def run(self):

        requests.get("http://"+self.sysinfo.svcIP+":"+str(self.sysinfo.svcPort)+"/cpu/stop")
        requests.get("http://"+self.sysinfo.svcIP+":"+str(self.sysinfo.svcPort)+"/memory/stop")

        Global.signal.stopBbtnCompleteSignal.emit()

        
        
if __name__ == '__main__':
    print(sys.path)
    #QCoreApplication.setLibraryPaths([sys.path[5] + '/PySide2/plugins'])
    app = QApplication(sys.argv)
    window = Form()
    window.show()
    app.exec_()
  