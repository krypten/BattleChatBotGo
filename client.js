var net = require('net');

var client = new net.Socket();
client.connect(4000, '127.0.0.1', function() {
    console.log('Connected');

    var step;
    for (step = 0; step < 5; step++) {
        // client.write('{}');
        // recipent_id: 1, sender_id: 1, message :{ id : 1 , type : text , body :  asdf sdf asdfasdfas} }');
    }
});

client.on('data', function(data) {
    console.log("Reply" + data.toString());
});

client.on('close', function() {
    console.log('Connection closed');
});
