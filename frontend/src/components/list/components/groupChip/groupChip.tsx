import React, { useEffect, useState } from "react";
import "./groupChip.css";
import { useNavigate } from "react-router-dom";
import { ImageGet } from "../../../../utils/ImageControl";


const GroupChip = ({ data }: any): JSX.Element => {
    const navigate = useNavigate()
    const [avatar, setAvatar] = useState('')
    const [optionButtons, setButtons] = useState<JSX.Element>(<></>)
    const notMember = (id: number) => {
        const requestHandle = (e: React.MouseEvent<HTMLButtonElement>) => {
            fetch(`http://localhost:8080/api/groups/request/${id}`, {
                method: "POST",
                credentials: "include"
            }).then((resp) => {
                if (resp.status === 202) {
                    setButtons(isPending())
                }
            })
        }
        return <button className="group__buttons" onClick={requestHandle}>REQUEST</button>
    }

    const isPending = () => {
        return <div className="group__buttons">PENDING</div>
    }

    const isMember = () => {
        return (
            <>
                <button className="group__buttons" onClick={() => { navigate(`/chats/groups/${data.id}`) }} >CHAT</button>
            </>
        )
    }

    const isCreator = () => {
        return <button className="group__buttons" onClick={() => { navigate(`/chats/groups/${data.id}`) }}>CHAT</button>
    }
    const isInvite = (id: number) => {
        const handleInviteAccept = () => {
            fetch(`http://localhost:8080/api/groups/inviteAccept/${id}`, {
                method: "POST",
                credentials: "include"
            }).then((resp) => {
                if (resp.status === 202) {
                    setButtons(isMember())
                }
            })
        }
        const handleInviteDeny = () => {
            fetch(`http://localhost:8080/api/groups/inviteDeny/${id}`, {
                method: "POST",
                credentials: "include"
            }).then((resp) => {
                if (resp.status === 202) {
                    setButtons(notMember(id))
                }
            })
        }
        return <>
            <button className="group__buttons" onClick={handleInviteAccept}>ACCEPT INVITATION</button>
            <button className="group__buttons" id="unconnect" onClick={handleInviteDeny}>DENY</button>
        </>
    }
    useEffect(() => {
        ImageGet(`group/${data.id}`, setAvatar)
        switch (data.status) {
            case 0:
                setButtons(notMember(data.id))
                break
            case 1:
                setButtons(isPending())
                break
            case 2:
                setButtons(isMember())
                break
            case 3:
                setButtons(isCreator())
                break
            case 4:
                setButtons(isInvite(data.id))
                break
        }
    }, [data])

    if (data.status == -1) return (<></>)

    return (
        <div className="group" id={data.id} >
            <div className="group__credentials" onClick={() => { navigate(`/groups/${data.id}`) }}>
                <img className="group__avatar" src={avatar ? avatar : '/assets/group-avatar.png'} loading="lazy" />
                <div className="group__name" data-type={'groups'}>{data.groupName}</div>
            </div>
            <div className="group__description">{data.groupAbout}</div>
            <div className="group__options">
                {optionButtons}
            </div>
        </div>
    );
};



export default GroupChip