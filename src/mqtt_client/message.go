package mqtt_client

import (
	"config"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"meta"
	"os"
	"time"
)

var c MQTT.Client

var identifier = config.Identifier
var secret = config.Secret

// message handle
var onMessage MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	msg_encoding := &meta.Message{}
	err := proto.Unmarshal(msg.Payload(), msg_encoding)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
	}
	requestId := msg_encoding.GetControl().GetRequestId()
	method := msg_encoding.GetControl().GetRequest().GetMethod()
	params := msg_encoding.GetControl().GetRequest().GetParameters()
	fmt.Printf(" Method: %s ,RequestId: %s \n Params: %s\n", method, requestId, params)
	err = nil
	switch method {
	case "openGate":
		err = openGate()
	case "cardOperation":
		err = cardOperation([]byte(params))
	case "passcodeOperation":
		err = passcodeOperation([]byte(params))
	case "ping":
		fmt.Println("ping test")
	case "reboot":
		fmt.Println("reboot")
	}
	if err == nil {
		response := proto.String(fmt.Sprintf("{\"response\": %d }", time.Now().Unix()))
		message_control_response := &meta.Message_Control_Response{Response: response}
		sendResponse(msg_encoding.GetFrom(), requestId, message_control_response)
	}
}

var onConnect MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("连接到mqtt服务器")
}
var disConnect MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Println("错误：", err)
	fmt.Println("断开 mqtt服务器")
	connect(client)
}

func connect(client MQTT.Client) {
	var flag = true
	for flag {
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			flag = false
			if token := c.Subscribe("/clients/"+identifier, 0, nil); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}
		} else {
			fmt.Println("重新连接出错")
		}
		time.Sleep(5 * time.Second)
	}
}

// publish response message
func sendResponse(from string, requestId string, message_control_response *meta.Message_Control_Response) {
	messageControl := &meta.Message_Control{
		Type:            meta.Message_Control_Type.Enum(meta.Message_Control_RESPONSE),
		RequestId:       proto.String(requestId),
		ControlPayloads: &meta.Message_Control_Response_{Response: message_control_response},
	}

	messagePayloads := &meta.Message_Control_{Control: messageControl}

	responseMsg := &meta.Message{
		Id:        proto.String(fmt.Sprintf("%s", uuid.NewV4())),
		Type:      meta.Message_Type.Enum(meta.Message_CONTROL),
		From:      proto.String("mock"),
		Timestamp: proto.Int64(0),
		Payloads:  messagePayloads,
	}
	outResponseMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		fmt.Println("Failed to encode response message:", err)
	}
	if token := c.Publish("/clients/"+from, 0, false, outResponseMsg); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func Start() {
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcps://%s:%s", config.MqttHost, config.MqttPort))
	opts.SetClientID(fmt.Sprintf("%s:[1]:%s", identifier, uuid.NewV4()))
	opts.SetDefaultPublishHandler(onMessage)
	opts.SetOnConnectHandler(onConnect)
	opts.SetConnectionLostHandler(disConnect)
	opts.SetUsername(identifier)
	opts.SetPassword(secret)
	opts.SetAutoReconnect(false)
	opts.SetTLSConfig(&tls.Config{
		RootCAs:            x509.NewCertPool(),
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          x509.NewCertPool(),
		InsecureSkipVerify: true,
	})
	c = MQTT.NewClient(opts)
	connect(c)
}
