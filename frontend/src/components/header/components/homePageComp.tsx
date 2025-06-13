import { useNavigate } from "react-router-dom"

const HomeHeader = () => {
    const navigate = useNavigate()

    const path = window.location.pathname // get the current URL
    const regex = new RegExp('\\d+$') // regex to find if its creation from home page or group page
    const groupId = path.match(regex) // find the matched groupId

    let link = groupId ? `/createPost/${groupId}` : '/createPost' 

    return (
        <div className="posts-page__create" onClick={() => { navigate(link) }}>
            <img className="posts-page__create-img" alt="create post" src="/assets/add-file-4.svg" />
            <button className="posts-page__create-btn" >CREATE POST</button>
        </div>
    )
}
export default HomeHeader