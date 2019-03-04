package server

type ID string

type Message struct {
    // id of the message
    Id ID `json:"id"`
    // type of the message
    ContentType string `json:"type"`
    // content of the message
    Body string `json:"body"`
}

type MessageBuilder struct {
    Family string
    Body string
}

func (msgBuilder *MessageBuilder) Build() *Message {
    return &Message { Id: "1", ContentType: msgBuilder.Family, Body: msgBuilder.Body }
}

type Response struct {
    RecipentId ID `json:"recipent_id"`
    SenderId ID `json:"sender_id"`
    Message *Message `json:"message"`
}

type Reply struct {
    RecipentId ID `json:"recipent_id"`
    MessageId ID `json:"message_id"`
}
