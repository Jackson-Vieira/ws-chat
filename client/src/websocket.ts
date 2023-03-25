import { WebsocketMessage } from "./@types/websocket.type";

export default class Websocket {
  ws: WebSocket | null;
  onMessage: ((message: WebsocketMessage) => void) | undefined;
  constructor(url: string) {
    this.ws = new WebSocket(url);
    this.ws.onopen = () => {
      this.send({
        messageType: "text",
        data: "Hello from client",
      });
    };
    this.ws.onmessage = (data) => {
      const message: WebsocketMessage = JSON.parse(data.data);

      switch (message.messageType) {
        case "text":
          if(this.onMessage != null) {
            this.onMessage(message);
          } else {
            console.log("No onMessage handler");          
            console.log('message', message.data)
          }
          break;
        default:
          console.log("Unknown message type");
          console.log(message.data)
      }
    };
  }

  send(message: WebsocketMessage) {
    if (this.ws == null) {
      return;
    }
    this.ws.send(JSON.stringify(message));
  }

  close() {
    if (this.ws == null) {
      return;
    }
    this.ws.close();
    this.ws = null;
  }
}
