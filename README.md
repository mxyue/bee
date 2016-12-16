# go主要逻辑部分

安装好 go1.7.3

##运行

运行run文件

    ./run

##编译

    ./build


##代码入口

    src/bee/main.go

##配置

    src/config/configration.go

    InDevice = false 
    //表示在电脑里面运行，只会模拟开门，不会直接调用开门代码逻辑。改成true后进行编译放入树莓派

#设备秘钥configs.txt

    将后台下载的设备identifier，secret。命名为configs.txt放入到执行文件所在的目录。
    host 服务器地址，mqtt_host Mqtt服务器地址


##数据存储

嵌入式key value 数据库: boltdb 


#硬件控制部分

- 4x4的键盘
- HC-SR501 人体红外感应模块
- MFRC522 读卡器，可刷写IC卡

##键盘和开门逻辑
    src/driver/raspberry_gpio_ctrl.go

##读卡部分
使用了树莓派的SPI功能（若没开启，需手动设置下）

安装 c实现的python的spi包 [github地址](https://github.com/lthiery/SPI-Py)

运行下面代码

    card_reader/card_reader.py

上面读卡代码读取的是ic卡内存储扇区。根据自己的需要修改读取数据的逻辑。
读取后将会调用go部分的门禁卡开门接口。目前go的包没有包含spi功能。若想将该部分逻辑整合到go的包一起。可以在go中去调用c。



