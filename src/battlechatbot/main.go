package main

import (
    "os"
    "net"
    "bufio"
    "encoding/json"

    "battlechat/lib"
    "./client"
)

const PORT = "4000"
const HOST_NAME = "localhost"

func main() {
    conn, err := net.Dial("tcp", HOST_NAME + ":" + PORT)
    if err != nil {
        lib.LoggerE(err, "Unable to connect to server.")    
        println("")
        return 
    }
    defer conn.Close()

    userId, recipentId := getConfig()

    channel := make(chan *client.Response)
    replyChan := make(chan *client.Reply)
    osChannel := make(chan string)
    errChannel := make(chan error)

    go watchConn(conn, channel, replyChan, errChannel)
    go watchConsole(osChannel, errChannel)

    defer conn.Close()

    for {
        select {
            case data := <-channel: 
                // watch connection
                if data != nil && data.Message != nil {
                    go Print(data.Message.Body)
                }
            case reply := <-replyChan:
                if reply != nil {
                    println("Received : " + string(reply.MessageId))
                }
            case data := <-osChannel:
                // watch console
                go Send(data, conn, userId, recipentId)
            case <-errChannel:
                break;
        }
    }
}

// parse out the arguments to be used when connecting to the chat server
func getConfig() (client.ID, client.ID) { 
  if (len(os.Args) >= 3) {
    userId := client.ID(os.Args[1])
    recipentId := client.ID(os.Args[2])
    return userId, recipentId
  } else {
    println("You must provide the receiver and sender ids as parameters.")
    os.Exit(1)
    return "", ""
  }
}

func watchConn(conn net.Conn, channel chan *client.Response, replyChan chan *client.Reply, errCh chan error) {
    for {
        // Make a buffer to hold incoming data.
        buf := make([]byte, 1024)
        // Read the incoming connection into the buffer.
        reqLen , err := conn.Read(buf)
        if err != nil {
            // send an error if it's encountered
            errCh<- err
            println("Error")
            return
        }
        println(string(buf))
        var response client.Response
        err = json.Unmarshal(buf[:reqLen], &response)
        if err != nil {
            lib.LoggerE(err, "There was an error decoding the response json.")
            return
        }
        if response.Message == nil {
            // try reply
            var reply client.Reply
            err = json.Unmarshal(buf[:reqLen], &reply)
            if err != nil {
                // send an error if it's encountered
                errCh<- err
            }
            replyChan <- &reply
            println("Send to reply channel")
        } else {
            channel <- &response
            println("Send to channel")
        }
    }
}

func watchConsole(channel chan string, errCh chan error) {
    reader := bufio.NewReader(os.Stdout)
    for {
        watch(reader, channel, errCh)
    }
}

func watch(reader *bufio.Reader, channel chan string, errCh chan error) {
    // try to read the data
    data, err := reader.ReadString('\n')
    if err != nil {
        // send an error if it's encountered
        errCh<- err
        return
    }
    // send data if we read some.
    channel<-data
}

func Send(message string, conn net.Conn, userId client.ID, recipentId client.ID) {
    msgBuilder := client.MessageBuilder{ Family: "text", Body: message }
    response, err := json.Marshal(client.Response { Message: msgBuilder.Build(), RecipentId: recipentId, SenderId: userId })
    if err != nil {
        lib.LoggerE(err, "Error occuring while converting to json")
        return
    }
    conn.Write(response)
    println("Message : " + message + " sent from " + string(userId) + " to " + string(recipentId))
}

func Print(message string) {
    println("Received message : " + message)
}
