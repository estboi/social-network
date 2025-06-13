import React, { useEffect, useState } from "react"
import fetchData from "../../utils/fetchData"
import { useParams } from "react-router-dom"
import List from "../../components/list/list"
import './groupOptionsPage.css'
interface Data {
    id: number
    isFollow: boolean
    isFollower: boolean
    isPending: boolean
    lastname: string
    name: string
}

const GroupOptions = () => {
    const [requestData, setData] = useState<Data[]>([])
    const [inviteData, setInvite] = useState<Data[]>([])
    const [toInvite, setNotMembers] = useState<Data[]>([])

    const [page, setPage] = useState(0)

    const { groupId } = useParams()
    useEffect(() => {
        fetchData(`groups/requested/${groupId}`, setData)
        fetchData(`groups/notmembers/${groupId}`, setNotMembers)
    }, [])

    return (
        <div className="group__option-page">
            <div className="group__option__pages">
                <button className={`group__option__pages-link ${page == 0 ? 'selected' : ''}`} onClick={() => { setPage(0) }}>INVITE</button>
                <button className={`group__option__pages-link ${page == 1 ? 'selected' : ''}`} onClick={() => { setPage(1) }}>REQUESTS</button>
            </div>
            {page == 0 && (
                <div className="group__option__invited">
                    <List className="" signal="usersToInvite" data={toInvite} />
                </div>
            )
            }
            {page == 1 && (<List className="" signal="users" data={requestData} />)}
        </div >
    )
}

export default GroupOptions