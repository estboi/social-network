
import { useEffect, useState } from "react";
import "./messageChip.css";
import { ImageGet } from "../../../../utils/ImageControl";

const Message = ({ data, myUser, chatType }: { data: any, myUser: number, chatType: string }): JSX.Element => {
    const [chatterAvatar, setAvatar] = useState('')
    useEffect(() => {
        ImageGet(`user/${data.senderId}`, setAvatar)
        if (chatterAvatar === '') {
            setAvatar('/assets/default_avatar.png')
        }
    }, [data])
    const messageClass = myUser === data.senderId ? 'message--user' : 'message--other';
    return (
        <div className={`${messageClass}`}>
            <div className="message__header">
                <img className="message__avatar" src={chatterAvatar} />
                <p className="message__nick-name">{data.sendername}</p>
                <p className="message__time">{data.time}</p>
            </div>
            <div className="message__text">
                {data.content}
            </div>
        </div>
    );
};
export default Message
