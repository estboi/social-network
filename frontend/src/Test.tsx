import { useEffect, useState } from "react";

const Test = () => {
    const [socket, setSocket] = useState<WebSocket | null>(null);

    useEffect(() => {
        const newSocket = new WebSocket('ws://localhost:8080/ws');

        // Set up event listeners for the WebSocket
        newSocket.addEventListener('open', () => {
            console.log('WebSocket connection opened');
        });

        newSocket.addEventListener('message', (event) => {
            console.log('WebSocket message received:', (event as MessageEvent).data);
        });

        newSocket.addEventListener('close', () => {
            console.log('WebSocket connection closed');
        });

        // Save the WebSocket instance to the state
        setSocket(newSocket);

        // Clean up the WebSocket connection on component unmount
        return () => {
            newSocket.close();
        };
    }, []); // Empty dependency array means this effect runs only once on mount

    return (
        <div>
            <button>
                Send WebSocket Message
            </button>
        </div>
    );
}

export default Test