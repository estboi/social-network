import { useNavigate } from "react-router-dom"
import HeaderBackButton from "./backButton"
import HomeHeader from "./homePageComp"
import { useEffect, useState } from "react"


const GroupProfileHeader = () => {

    const path = window.location.pathname // get the current URL
    const regex = new RegExp('\\d+$') // regex to find if its creation from home page or group page
    const groupId = path.match(regex) // find the matched groupId

    const [error, setError] = useState(false)

    useEffect(() => {
        fetch(`http://localhost:8080/api/groups/profile/${groupId}`, { credentials: 'include' })
            .then((resp) => {
                console.log(resp)
                if (resp.status === 406) {
                    setError(true)
                }
            })
    }, [groupId])


    const CreateEvent = () => {
        const navigate = useNavigate()
        return (
            <div className="posts-page__create" onClick={() => { navigate(`groups/${groupId}/createEvent`) }}>
                <img className="posts-page__create-img" alt="create post" src="/assets/add-file-4.svg" />
                <button className="posts-page__create-btn" >CREATE EVENT</button>
            </div> 
        )
    }
    return (
        <>
            <HeaderBackButton />
            <>
                {error ?
                    <></> : (
                        <>
                            <HomeHeader />
                            <CreateEvent />
                        </>
                    )}
            </>
        </>
    )
}

export default GroupProfileHeader