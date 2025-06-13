import { Post } from "./components/postChip/post"
import GroupChip from "./components/groupChip/groupChip"
import UserChip from "./components/userChip/usersChip"
import EventChip from "./components/eventChip/eventChip"
import PostComment from "./components/postChip/comment"
import ChatChip from "./components/chatChip/chatChip"
import NotificationChip from "./components/notificationChip/NotificationChip"
import "./list.css"
import React from "react"

interface listType {
    signal: string,
    className: string,
    data: any[]
}

type listComponentsType = {
    [key: string]: JSX.Element
}
const listComponents: listComponentsType = {
    'posts': <Post />,
    'comments': <PostComment />,
    'users': <UserChip />,
    'usersToInvite': <UserChip />,
    'groups': <GroupChip />,
    'events': <EventChip />,
    'chats': <ChatChip />,
    'notifications': <NotificationChip />
}

const List = React.forwardRef<HTMLDivElement, listType>(({ signal, className, data }, ref) => {
    const renderComponent = listComponents[signal]
    const invited = signal === 'usersToInvite' ? 'invited' : ''
    return (
        <>
            {!data || data.length === 0 ? (<>
                {signal === 'users' && (<p className="list__no-events">NO USERS</p>)}
                {signal === 'posts' && (<p className="list__no-events">NO POSTS YET</p>)}
                {signal === 'events' && (<p className="list__no-events">NO EVENTS YET</p>)}
                {signal === 'groups' && (<p className="list__no-events">NO GROUPS YET</p>)}
                {signal === 'chats' && (<p className="list__no-chats">NO CHATS YET</p>)}
                {signal === 'notifications' && (<p className="list__no-events">NO NOTIFICATIONS</p>)}
            </>
            ) : (
                <div className={`list ${className}`} ref={ref}>
                    {data.map((dataset: any, key: any) => (
                        React.cloneElement(renderComponent, { data: dataset, key: key, invited: invited })
                    ))}
                </div>
            )}
        </>)

})

export default List