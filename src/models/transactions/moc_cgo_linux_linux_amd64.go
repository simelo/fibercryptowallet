package transactions

/*
#cgo CFLAGS: -pipe -O2 -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -pipe -O2 -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -I../../models -I. -I/home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include -I/home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtGui -I/home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtQml -I/home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtNetwork -I/home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtCore -I. -isystem /usr/include/libdrm -I/home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/mkspecs/linux-g++
#cgo LDFLAGS: -O1 -Wl,-rpath,/home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib
#cgo LDFLAGS:  /home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Gui.so /home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Qml.so /home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Network.so /home/lag/Installations/go-1.12/golang/src/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Core.so -lGL -lpthread
#cgo CFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
#cgo CXXFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
*/
import "C"
