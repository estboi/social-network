import "./chatPage.css";
import { useParams } from "react-router-dom";
import React, { useState, useEffect, useRef, useContext } from "react";
import fetchData from "../../utils/fetchData";
import EmojiPicker from "../../utils/emojiFuncs";
import { WebsocketContext, WebsocketContextProps } from '../../websocket/Provider'; // Adjust the path
import { Event } from "../../websocket/Provider";
import ListPage from "../listPage/listPage";
import Message from "../../components/list/components/messageChip/messageChip";

interface messageVM {
    chatterId: number
    chatterName: string
    lastMessage: string
    messageTime: string
    type: string
}

interface ChatType {
    type: string
}

const Chat = ({ type }: ChatType) => {
    const [openEmoji, setEmoji] = useState(false)
    const newMessageRef = useRef<HTMLInputElement>(null)
    const [getData, getChatData] = useState<messageVM[]>([])
    const { ready, value, send }: WebsocketContextProps = useContext(WebsocketContext);
    const [sendData, setData] = useState<{
        chatId: number;
        content: string;
        type: string
    }>({ chatId: 0, content: '', type: '' })
    const [error, setError] = useState('')
    const { userId, groupId } = useParams()
    let chatString: string = ''
    if (type === 'group' && groupId !== undefined) {
        chatString = groupId
    } else if (type === 'user' && userId !== undefined) {
        chatString = userId
    } else {
        console.log(`Chat doesn't exist: ${userId}, ${groupId}`)
    }
    const id = parseInt(chatString, 10)

    const [myUser, setMyUser] = useState<{ username: string, userID: number }>({ username: '', userID: 0 })
    const fetchDataAndSetData = async () => {
        if (type === 'user') {
            await fetchData(`chats/user/${id}`, getChatData);
        } else if (type === 'group') {
            await fetchData(`chats/group/${id}`, getChatData);
        } else {
            console.log('wtf bro');
        }
        await fetchData('navbar', setMyUser)
    };

    useEffect(() => {
        fetchDataAndSetData();
    }, [value])

    const handleMessageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        e.preventDefault();
        const emojiText = newMessageRef.current?.value || '';
        sendData.chatId = id
        sendData.content = emojiText
        sendData.type = type
    };

    const handleSendMessage = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent> | React.MouseEvent<HTMLDivElement> | React.FormEvent) => {
        e.preventDefault();
        console.log(sendData)
        if (newMessageRef.current && sendData.content.trim() !== '') {
            const event = new Event('New_Message', sendData);
            send(JSON.stringify(event));
            fetchDataAndSetData(); // Move this line if you want to fetch data after sending a message
            newMessageRef.current.value = '';
        }
    };

    const handleClick = (event: React.MouseEvent<HTMLDivElement>) => {
        event.preventDefault()
        const parent = event.currentTarget.parentElement?.parentElement?.parentElement
        if (parent) {
            setEmoji(false)
        }
    }
    const closeBtnEmoji = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault()
        setEmoji(false)
    }

    return (
        <div className="chat__page">
            <div className="chat-window">
                <div className='chat-feed'>
                    {getData.map((dataset: any, key: any) => (
                        <Message key={key} data={dataset} myUser={myUser.userID} chatType={dataset.type} />
                    ))}
                </div>
                <form className="chat-options" onSubmit={handleSendMessage}>
                    <input className="chat-input" placeholder="Enter your message" ref={newMessageRef} onChange={handleMessageChange} />
                    <div className='post__comment-btn--wrapper' onClick={(e) => {
                        handleClick(e)
                        setEmoji(true)
                    }}>
                        <img className='post__new-comment-img' src='/assets/AddSmile.svg' />
                    </div>
                    <div className="emoji--wrapper">
                        {openEmoji && <EmojiPicker input={newMessageRef} closeBtn={closeBtnEmoji} className={'chat__emoji-list'} setData={setData} />}
                        <div className='post__comment-btn--wrapper chat' onClick={handleSendMessage}>
                            <img className='post__new-comment-img' src='/assets/Submit.svg' />
                        </div>
                    </div>
                </form>
            </div>
            <ListPage signal={['chats', 'all']} />
        </div>
    );
};
export default Chat