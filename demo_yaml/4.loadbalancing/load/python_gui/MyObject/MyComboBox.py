from Configure.configure import *

class MyComboBox(QtWidgets.QComboBox):

    def __init__(self, styleNo=1):
        super().__init__()

        self.setStyle(styleNo)
        self.setView(QtWidgets.QListView())
        cssText = Global.readCSS(self.fileName)
        self.setStyleSheet(cssText)
        self.setSizePolicy(
                QtWidgets.QSizePolicy.Expanding,QtWidgets.QSizePolicy.Expanding
            )
    def setStyle(self, styleNo):
        if styleNo == 1:
            self.fileName = "resources/css/comboBox/combo1.css"
        elif styleNo == 2:
            self.fileName = "resources/css/comboBox/combo2.css"
        elif styleNo == 3:
            self.fileName = "resources/css/comboBox/combo3.css"
