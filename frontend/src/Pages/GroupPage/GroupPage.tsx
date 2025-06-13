import GroupBar from "./GroupBar"
import './GroupPage.css'
import { useParams } from "react-router-dom"
import List from "../../components/list/list"
import { useState, useEffect } from "react"
import fetchData from "../../utils/fetchData"

function GroupPage() {
    const [data, setData] = useState([])
    const [error, setError] = useState('')

    const { groupId } = useParams()
    if (groupId === undefined) {
        setError('Group doesn\'t exist')
    }
    const groupString = groupId || ''
    const Id = parseInt(groupString, 10)

    useEffect(() => {
        fetchData(`posts/group/${Id}`, setData)
    }, [])

    return (
        <>
            {error ?
                <div className="group-page__error">
                    <h1 className="error">YOU AREN'T MEMBER OF THIS GROUP</h1>
                </div>
                :
                <div className="group-page">
                    <List signal={'posts'} className={'posts'} data={data} />
                    <GroupBar id={Id} error={setError} />
                </div>
            }
        </>

    )
}

export default GroupPage

