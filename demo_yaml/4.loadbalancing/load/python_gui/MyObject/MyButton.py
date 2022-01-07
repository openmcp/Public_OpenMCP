from Configure.configure import *

class MyButton(QtWidgets.QPushButton):
    def __init__(self, text, styleNo=1):
        super().__init__()
        self.setText(text)
        self.setStyle(styleNo)
        self.setSizePolicy(
            QtWidgets.QSizePolicy.Expanding,QtWidgets.QSizePolicy.Expanding
        )


    def setStyle(self, styleNo):
        if styleNo == 1:
            self.fileName = "resources/css/button/button1.css"
        elif styleNo == 2:
            self.fileName = "resources/css/button/button2.css"
        elif styleNo == 3:
            self.fileName = "resources/css/button/button3.css"
        elif styleNo == 4:
            self.fileName = "resources/css/button/button4.css"
        elif styleNo == 5:
            self.fileName = "resources/css/button/button5.css"


        btn_css = Global.readCSS(self.fileName)
       
        self.setStyleSheet(btn_css)
