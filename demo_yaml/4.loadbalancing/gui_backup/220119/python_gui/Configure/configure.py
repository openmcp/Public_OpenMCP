from PySide2 import QtGui, QtCore, QtWidgets

class Signal(QtCore.QObject):
    stopBbtnCompleteSignal = QtCore.Signal()

    
class Global(object):
    signal = Signal()
    
    FONT_Family = "MS Sans Serif"
    
    Traffic_Period = 3
    #Traffic_URL = "http://cluster03.keti.wordpress.openmcp.in"
    Traffic_URL = "http://keti.productpage.openmcp.in:8181/productpage"

    Load_Namespace = "bookinfo"
    Load_Svc_Name = "productpage"
    
    winWidth = 700
    winHeight = 400
    winX = None
    winY = None
    
    class image:
        path = 'resources/image/'

        # # Title ToolBar
        # minimize = path + 'minimize.png'
        # maximize = path + 'maximize.png'
        # maximize2 = path + 'maximize2.png'
        # close = path + 'close.png'

        # # ComboBox
        # downarrow = path + 'downarrow.png'

        # # Shape
        # circle = path + 'circle.png'
        # rect = path + 'rect.png'


    def readCSS(fileName):
        f = open(fileName, 'rt', encoding = 'UTF8')
        cssText = f.read()
        f.close()

        list = cssText.split("\n")

        for index, line in enumerate(list):
            if "font-family" in line:
                list[index] = "font-family : "+ Global.FONT_Family + ";"

        cssText = '\n'.join(list)
        return cssText
