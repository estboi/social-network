import { useEffect, useState } from "react";
import "./chatChip.css";
import { ImageGet } from "../../../../utils/ImageControl";
import { useNavigate } from "react-router-dom";

const ChatChip = ({ data }: any) => {
    const navigate = useNavigate()
    const [chatterAvatar, setAvatar] = useState('')
    useEffect(() => {
        if (data.type === 'users') {
            ImageGet(`user/${data.otheruser}`, setAvatar)
        } else if (data.type === 'groups') {
            ImageGet(`group/${data.chatId}`, setAvatar)
        }
        if (chatterAvatar == '') {
            setAvatar('/assets/default_avatar.png')
        }
    }, [data])
    return (
        <div className='chat-chip' onClick={() => {
            if (data.type == 'users') {
                navigate(`/chats/users/${data.otheruser}`)
            } else if (data.type == 'groups') {
                navigate(`/chats/groups/${data.chatId}`)
            }
        }}>
            <div className="chat-chip__header">
                <img className="chat-chip__avatar" src={chatterAvatar} />
                <p>{data.firstname}</p>
            </div>
            <p className="chat-chip__time">{data.time}</p>
            <div className="chat-chip__last-message">
                <div className="chat-chip__text">{data.otheruser == data.senderId ? 'you' : data.sendername}: {data.content}</div>
            </div>
        </div>
    );
};

export default ChatChip