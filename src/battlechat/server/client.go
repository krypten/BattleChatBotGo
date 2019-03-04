package server

import (
    "net"
    "strconv"
    "encoding/json"
    "battlechat/lib"
)

type Client struct {
    // Channel for client created 
    Channel chan Response
    // the client's connection
    Connection net.Conn
    // the client's id
    Id ID
}

var clients = make(map[ID]*Client)

// Register the connection and cache it
func (client *Client) Register() {
    client.Id = ID(strconv.Itoa(len(clients) + 1))
    lib.LoggerF("Registering new client with id : %s .\n", string(client.Id))
    clients[client.Id] = client
}

// Deregister the connection with client
func (client *Client) Close() {
    client.Connection.Close()
    delete(clients, client.Id)
}

func (client *Client) Send(message *Message, sender *Client) {
    response := Response{ Message: message, RecipentId: client.Id, SenderId: sender.Id }
    println(string("Response : " + string(ToJsonResponse(response))))
    client.Connection.Write(ToJsonResponse(response))
    reply := Reply{ RecipentId : client.Id, MessageId: message.Id }
    println(string("Reply : " + string(ToJsonReply(reply))))
    sender.Connection.Write(ToJsonReply(reply))
}

func ToJsonResponse(response Response) []byte {
    obj, _ := json.Marshal(response)
    return obj
}

func ToJsonReply(reply Reply) []byte {
    obj, _ := json.Marshal(reply)
    return obj
}
