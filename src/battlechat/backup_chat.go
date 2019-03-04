package server

import (
    "battlechat/lib"
    "bufio"
    "encoding/json"
    "net"
    "io"
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
        go func(c net.Conn) {
            io.Copy(os.Stdout, c)
        }(conn)

        // register new client
        channel := make(chan *Response)
        client := Client{Connection: conn, Channel: channel}
        client.Register()

        // handle traffic for the client
        go HandleChannel(&client, errChan)
    }
}

func HandleChannel(client *Client, errCh chan error) {
    watchConn(client, errCh)
    for {
        select {
            case err := <-errCh:
                lib.LoggerE(err, "Error occured for Client id : %s.\n", string(client.Id))
            case response := <-client.Channel:
                sender := clients[response.SenderId] // TODO
                client.Send(response.Message, sender)
            // case client.Channel<- :
               //  response := client.Receive()
                // clients[response.recipentId].Channel <- response.message
        }
    }
}

func watchConn(client * Client, errCh chan error) {
    defer client.Connection.Close()
    for {
        reader := bufio.NewReader(client.Connection)
        watch(reader, client.Channel, errCh)
    }
}

func watch(reader *bufio.Reader, channel chan *Response, errCh chan error) {
    // try to read the data
    data, err := reader.ReadString('\n')
    if err != nil {
        // send an error if it's encountered
        errCh<- err
        println("Error")
        return
    }
    println(data)
    // send data if we read some.
    var response Response
    err = json.Unmarshal([] byte(data), &response)
    if err != nil {
        errCh<- err
        return
    }
    x, err := json.Marshal(response)
    println(string(x))
    channel<- &response
}
