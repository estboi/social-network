import React, { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import fetchData from "../../utils/fetchData"
import { ImageGet } from "../../utils/ImageControl"

interface GroupBarProps {
    id: number
    error: any
}

interface GroupData {
    groupName: string,
    groupAbout: string,
    status: number
}

function GroupBar({ id, error }: GroupBarProps) {
    const navigate = useNavigate()
    const [data, setData] = useState<GroupData>({ groupName: '', groupAbout: '', status: 2 })
    const [avatar, setAvatar] = useState('')
    useEffect(() => {
        fetchData(`groups/profile/${id}`, setData)
            .then((resp) => {
                if (resp === 406) {
                    error('NOTALLOWED')
                }
            })
            .then(() => {
                ImageGet(`group/${id}`, setAvatar)
            })
    }, [])

    const openGroupOption = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        navigate(`members`) // Navigating to group options page. Where you can Invite or Response to attendance
    }

    return (
        <div className="profileInfo group-bio">
            <div className="infoblock">
                <img alt='' src={avatar ? avatar : '/assets/group-avatar.png'} />
                <p className="group-bio__name">{data.groupName}</p>
                <div>
                    <p className="group-bio__about-title">About</p>
                    <div className="group-bio__about">
                        <p>{data.groupAbout}</p>
                    </div>
                </div>
                <button className="change-button" onClick={() => { navigate(`/events/${id}`) }}>EVENTS</button>
                {data.status === 3 && <button className="change-button" onClick={openGroupOption}>MEMBERS</button>}
                <button className="change-button" onClick={() => { navigate(`/chats/groups/${id}`) }}>CHAT</button>
            </div>
        </div>
    )
}
export default GroupBar