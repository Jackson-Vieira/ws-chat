import "./style.css";


import { WebsocketMessage } from "./@types/websocket.type";
import Websocket from "./websocket";

const WS_URL = "ws://localhost:8080/ws";
const websocket = new Websocket(WS_URL);