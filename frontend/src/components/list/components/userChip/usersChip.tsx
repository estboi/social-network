import { useNavigate, useParams } from "react-router-dom";
import "./usersChip.css";
import { useEffect, useState } from "react";
import { ImageGet } from "../../../../utils/ImageControl";

const UserChip = ({ data, invited }: any): JSX.Element => {
    const navigate = useNavigate()
    const [FollowStatus, setFollowStatus] = useState({ isFollowing: data.isFollower, isPending: data.isPending })
    const [avatar, setAvatar] = useState('')
    const [groupStatus, setGroupStatus] = useState(false)
    const [isInvited, setInvited] = useState(false)

    const { groupId } = useParams()
    useEffect(() => {
        if (groupId && !invited) {
            setFollowStatus({ isFollowing: false, isPending: true })
        }
        console.log(data)
        ImageGet(`user/${data.id}`, setAvatar)
    }, [])

    const handleUnsubscribe = () => {
        FollowFetch('unsubscribe', data.id)
        setFollowStatus({ isFollowing: false, isPending: false });
    }

    const handleSubscribe = () => {
        FollowFetch('subscribe', data.id)
        setFollowStatus({ isFollowing: true, isPending: false });
    }

    const handleGroupAccept = () => {
        const sendData = { userID: data.id }
        fetch(`http://localhost:8080/api/groups/accept/${groupId}`, {
            credentials: "include",
            method: "POST",
            body: JSON.stringify(sendData)
        }).then((resp) => {
            switch (resp.status) {
                case 202:
                    setGroupStatus(true)
                    break
            }
        })
    }
    const handleGroupDeny = () => {
        const sendData = { userID: data.id }
        fetch(`http://localhost:8080/api/groups/deny/${groupId}`, {
            credentials: "include",
            method: "POST",
            body: JSON.stringify(sendData)
        }).then((resp) => {
            switch (resp.status) {
                case 202:
                    setGroupStatus(true)
                    break
            }
        })
    }

    const handleInvite = () => {
        const sendData = { userID: data.id }
        fetch(`http://localhost:8080/api/groups/invite/${groupId}`, {
            credentials: "include",
            method: "POST",
            body: JSON.stringify(sendData)
        }).then((resp) => {
            switch (resp.status) {
                case 202:
                    setInvited(true)
                    break
            }
        })
    }

    if (groupStatus) {
        return <></>
    }

    return (
        <div className="users-list__user">
            <div className="users-list__user-header" onClick={() => {
                navigate(`/users/${data.id}`, { state: { targetId: data.id } })
            }}>
                <img className="users-list__user-avatar" src={avatar ? avatar : '/assets/default_avatar.png'} loading="lazy" />
                <div className="users-list__user-name" data-type={'users'}>{data.name}</div>
            </div>
            <div className="users-list__options">
                {invited != '' ?
                    (<>
                        {isInvited || FollowStatus.isPending ?
                            (<div className="users-list__button">PENDING...</div>)
                            :
                            (<button className="users-list__button" onClick={handleInvite} >INVITE</button>)
                        }
                    </>)
                    :
                    (<>
                        {FollowStatus.isFollowing ? (
                            <>
                                <button className="users-list__button" onClick={() => { navigate(`/chats/users/${data.id}`) }}>CHAT</button>
                                <button className="users-list__button"
                                    id="unsub"
                                    onClick={handleUnsubscribe}>UNSUBSCRIBE</button>
                            </>
                        ) : (
                            FollowStatus.isPending ? (
                                <>
                                    <button className="users-list__button" onClick={groupId ? handleGroupAccept : handleSubscribe}>ACCEPT</button>
                                    <button className="users-list__button" id="unsub" onClick={groupId ? handleGroupDeny : handleUnsubscribe}>DECLINE</button>
                                </>
                            ) : (
                                <button className="users-list__button"
                                    onClick={handleSubscribe}>SUBSCRIBE</button>
                            )
                        )}
                    </>)}
            </div>
        </div>
    );
};

const FollowFetch = async (signal: string, id: number) => {
    try {
        await fetch(`http://localhost:8080/api/users/${signal}/${id}`, {
            method: 'POST',
            credentials: "include"
        })
    } catch (error) {
        console.error('Error fetching data:', error);
    }
}

export default UserChip