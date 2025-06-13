import { useParams } from 'react-router-dom'
import { useState, useEffect } from 'react'
import List from '../../components/list/list'
import fetchData from '../../utils/fetchData'
import Profilebar from './ProfileBio/ProfileBio'

import './ProfilePage.css'

function ProfilePage() {
    const [data, setData] = useState([])
    const [error, setError] = useState(true)
    const { userId } = useParams()
    let userIdInt: number
    if (userId === undefined) {
        console.log("user with does not exist")
    }
    const chatString = userId || ''
    userIdInt = parseInt(chatString, 10)

    useEffect(() => {
        fetchData(`posts/users/${userIdInt}`, setData)
    }, [userIdInt])

    return (
        <>
            {error === false ? (
                <div className="group-page__error">
                    <h1 className="error">SUBSCRIBE TO SEE THIS PRIVATE PROFILE</h1>
                </div>
            ):
            (
            <div className="profile-page">
                <List signal={'posts'} className={'posts'} data={data} />
                <Profilebar id={userIdInt} setError={setError}/>
            </div>
        )}
        </>
    )
}

export default ProfilePage
