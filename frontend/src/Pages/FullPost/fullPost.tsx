import { useLocation, useNavigate, useParams } from "react-router-dom"

import './fullPost.css'
import AddComment from "../../components/list/components/postChip/components/postAdd/postAdd"
import List from "../../components/list/list"
import { useEffect, useState } from "react"
import fetchData from "../../utils/fetchData"
import { Post } from "../../components/list/components/postChip/post"


interface PostData {
    postId: string
    postContent: string
    postCreator: {
        creatorId: string
        creatorName: string
        creatorImg: string
    }
    postTime: string
    postVotes: number
    postImg: string
}



const FullPost = () => {
    const [postData, setPostData] = useState(undefined)
    const [commentsData, setCommentsData] = useState(undefined)
    const [newComment, setComment] = useState<number>(0)

    const { postID } = useParams()
    useEffect(() => {
        fetchData(`posts/${postID}`, setPostData)
        fetchData(`posts/comments/${postID}`, setCommentsData)
    }, [])

    useEffect(() => {
        fetchData(`posts/comments/${postID}`, setCommentsData)
    }, [newComment])

    if (!postData || !commentsData && commentsData !== null) {
        return (
            <h1 style={{ color: "whitesmoke" }}>LOADING...</h1>
        )
    }


    const addNewComment = () => setComment((p) => p + 1)

    return (
        <div className="full-post--page">
            <Post data={postData} className={'full-post'} addNewComment={addNewComment} />
            <div className="full-post__comment--wrapper">
                <List signal="comments" className="comments" data={commentsData} />
            </div>
        </div>
    )
}

export default FullPost
