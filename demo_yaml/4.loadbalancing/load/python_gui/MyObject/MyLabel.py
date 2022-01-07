from Configure.configure import *


class MyLabel(QtWidgets.QLabel):

    def __init__(self, text="", styleNo=1):
        super().__init__()
        self.setText(text)
        self.setStyle(styleNo)

        cssText = Global.readCSS(self.fileName)
        self.setStyleSheet(cssText)
        self.setSizePolicy(
                QtWidgets.QSizePolicy.Expanding,QtWidgets.QSizePolicy.Expanding
            )
    def setStyle(self, styleNo):
        if styleNo == 1:
            self.fileName = "resources/css/label/label_title.css"
        elif styleNo == 2:
            self.fileName = "resources/css/label/label_16pt.css"
        elif styleNo == 3:
            self.fileName = "resources/css/label/label_12pt.css"





