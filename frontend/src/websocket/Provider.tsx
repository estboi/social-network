import React, { createContext, useEffect, useRef, useState } from 'react';

export class Event {
    type: string
    payload: any
    constructor(type: string, payload: any) {
        this.type = type
        this.payload = payload
    }
}

export interface WebsocketContextProps {
    ready: boolean;
    value: any;
    send: ((data: any) => void);
}

export const WebsocketContext = createContext<WebsocketContextProps>({
    ready: false,
    value: null,
    send: () => { },
});

interface WebsocketProviderProps {
    children: React.ReactNode;
    isLogged: boolean
}

export const WebsocketProvider: React.FC<WebsocketProviderProps> = ({ children, isLogged }) => {
    const [isReady, setIsReady] = useState(false);
    const [val, setVal] = useState<any>(null);
    const ws = useRef<WebSocket | null>(null);

    useEffect(() => {
        if (ws.current != null) {
            ws.current.close()
        }
        console.log(isLogged)
        const socket = new WebSocket("ws://localhost:8080/ws");

        socket.onopen = () => setIsReady(true);
        socket.onclose = () => setIsReady(false);
        socket.onmessage = (event) => setVal(event.data);

        ws.current = socket;

        return () => {
            if (ws.current) {
                ws.current.close();
            }
        };
    }, [isLogged]);

    const send = (data: any) => {
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(data);
        }
    };

    const contextValue: WebsocketContextProps = {
        ready: isReady,
        value: val,
        send: send,
    };

    return (
        <WebsocketContext.Provider value={contextValue}>
            {children}
        </WebsocketContext.Provider>
    );
};
