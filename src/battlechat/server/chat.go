package server

import (
    "battlechat/lib"
    "encoding/json"
    "net"
    "os"
)

var ListenAndServe = func(port string) error {
    println("Hello")
    psock, err := net.Listen("tcp", ":" + port)
    if (err != nil) {
        lib.LoggerE(err, "Unable to start chat server.")
        os.Exit(1)
    }

    println("Start : " + port)
    lib.LoggerF("Chat server started on port %v...\n", port)

    // creating error channel
    errChan := make(chan error)

    defer psock.Close()
    for {
        // accepting client connections
        conn, err := psock.Accept()
        if (err != nil) {
            lib.LoggerE(err, "Unable to accept connections.")
            return err
        }

        // register new client
        channel := make(chan Response)
        client := Client{Connection: conn, Channel: channel}
        client.Register()

        // handle traffic for the client
        go HandleChannel(&client, errChan)
    }
}

func HandleChannel(client *Client, errCh chan error) {
    go watchConn(client, errCh)
    for {
        select {
            case err := <-errCh:
                lib.LoggerE(err, "Error occured for Client id : %s.\n", string(client.Id))
                client.Close()
                return
            case response := <-client.Channel:
                println("Channel sent something")
                recipent := clients[response.RecipentId] // TODO
                if recipent != nil {
                    go recipent.Send(response.Message, client)
                }
        }
    }   
}

func watchConn(client *Client, errCh chan error) {
    for {
        watch(client, errCh)
    }
}

func watch(client *Client, errCh chan error) {
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    // Read the incoming connection into the buffer.
    reqLen , err := client.Connection.Read(buf)
    if err != nil {
        // send an error if it's encountered
        errCh<- err
        println("Error")
        return
    }
    println(string(buf))
    var response Response
    err = json.Unmarshal(buf[:reqLen], &response)
    if err != nil {
        lib.LoggerE(err, "There was an error decoding the json.")
        // send an error if it's encountered
        errCh<- err
        os.Exit(255)
    }
    client.Channel <- response 
    println("Send to channel")
}
