# go-ws-ascendex

Description 
 
It is required to implement connection to ascendex WebSocket exchange (https://ascendex.github.io/ascendex-pro-api/#websocket). The connection interface is contained in the apiclient.go file. 
 
Connection - the function should implement the connection to the WebSocket exchange. If there are issues with the connection, it should return an error. 
 
Disconnect - the function should implement the disconnection from the WebSocket exchange. 
 
SubscribeToChannel - the function should implement listening to the WebSocket channel to receive BBO. If there are issues with listening, it should return an error and disconnect correctly from the WebSocket. 
 
ReadMessagesFromChannel - the function should implement reading the WebSocket channel about BBO, correctly transforming the data, and writing them to a chan. 
 
WriteMessagesToChannel - the function should implement writing to the WebSocket channel to keep the connection open. 
