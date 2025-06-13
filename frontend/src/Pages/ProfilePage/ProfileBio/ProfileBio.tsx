import { useState, useEffect } from 'react'
import '../ProfilePage.css'
import fetchData from '../../../utils/fetchData'
import { ImageGet } from '../../../utils/ImageControl'

interface ProfileBarProps {
    id: number
    setError: Function
}

function Profilebar({ id, setError }: ProfileBarProps) {

    const [data, setData] = useState<{
        name: string
        nickname?: string
        date: string
        email: string
        avatar?: File
        about?: string
        privacy: string
        canModify: boolean
    }>({
        name: '',
        date: '',
        email: '',
        privacy: '',
        canModify: false
    })
    const [avatar, setAvatar] = useState('')

    useEffect(() => {
        fetchData(`users/profile/${id}`, setData)
            .then((resp) => {
                console.log(resp);
                if (resp === true) {
                    setError(true)
                }
                if (resp === 406) {
                    setError(false)
                }
            })
        ImageGet(`user/${id}`, setAvatar)
    }, [])


    const handlePrivacy = () => {
        fetch('http://localhost:8080/api/users/modify', {
            method: "POST",
            credentials: "include"
        }).then(() => {
            fetchData(`users/profile/${id}`, setData)
        })
    }

    return (
        <div className="profileInfo">
            <div className="infoblock">
                <img alt='' src={avatar ? avatar : '/assets/default_avatar.png'} />
                <p className="group-bio__name">{data.name}</p>
                <>
                    {data.nickname && (
                        <p className='group-bio'>AKA: <strong>{data.nickname}</strong></p>
                    )}
                </>
                <p className="group-bio">BIRHT DATE: <strong>{data.date}</strong></p>
                <p className="group-bio">EMAIL: <strong>{data.email}</strong></p>
                <>
                    {data.about && (<div>
                        <p className="group-bio__about-title">About</p>
                        <div className="group-bio__about">
                            <p>{data.about}</p>
                        </div>
                    </div>)}
                </>
                <p className='profile-privacy'>Privacy: <strong>{data.privacy}</strong></p>
                <>
                    {data.canModify === true && (
                        <button className='change-button' onClick={handlePrivacy}>Change privacy</button>
                    )}
                </>
            </div>
            <button className='change-button-profile' onClick={() => {
                fetch('http://localhost:8080/api/logout', {
                    method: "POST",
                    credentials: "include"
                })
                    .then(() => { window.location.reload(); })
            }}>LOGOUT</button>
        </div>
    )
}

export default Profilebar
