from Configure.configure import *


class MyTextEdit(QtWidgets.QTextEdit):

    def __init__(self, text, styleNo=1):
        super().__init__()

        
        self.setText(text)
        self.setStyle(styleNo)
        
        self.setSizePolicy(
            QtWidgets.QSizePolicy.Expanding,QtWidgets.QSizePolicy.Expanding
        )
        self.setReadOnly(1)

    def setStyle(self, styleNo):
        if styleNo == 1:
            self.fileName = "resources/css/textEdit/text_edit1.css"
            
        cssText = Global.readCSS(self.fileName)
        self.setStyleSheet(cssText)