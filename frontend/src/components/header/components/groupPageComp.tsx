import { useNavigate } from "react-router-dom"

const GroupsHeader = () => {
    const navigate = useNavigate()
    return (
        <div className="posts-page__create" onClick={() => { navigate('/groups/createGroup') }}>
            <img className="posts-page__create-img" alt="create post" src="/assets/add-file-4.svg" />
            <button className="posts-page__create-btn" >CREATE GROUP</button>
        </div>
    )
}
export default GroupsHeader