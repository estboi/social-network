import { useEffect, useState } from "react";
import "./comment.css";
import { useNavigate } from 'react-router-dom'
import { ImageGet } from "../../../../utils/ImageControl";

interface voteData {
    votes: number
    target: number
}

const PostComment = ({ data }: any): JSX.Element => {
    const navigate = useNavigate()
    const [avatarImage, setAvatarImage] = useState('')
    const [commentImage, setCommentImage] = useState('')
    const fetchCommentImage = async () => {
        await ImageGet(`comment/${data.commentId}`, setCommentImage)
    }
    const fetchAvatarImage = async () => {
        await ImageGet(`user/${data.commentCreator.creatorId}`, setAvatarImage)
    }
    useEffect(() => {
        fetchCommentImage()
        fetchAvatarImage()
        console.log(data.postVotesTarget)
        if (data.postVotesTarget === 1) {
            setSelectedUp(true)
        } else if (data.postVotesTarget === -1) {
            setSelectedDown(true)
        }
    }, [data])

    const [selectedUp, setSelectedUp] = useState(false)
    const [selectedDown, setSelectedDown] = useState(false)
    const [postVotes, changeVotes] = useState(data.commentVotes)

    const voteHandler = async (signal: boolean, event: React.MouseEvent<HTMLImageElement>) => {
        const target = event.target as HTMLElement
        if (target.id === 'postUpVote' && selectedUp === true) return
        if (target.id === 'postDownVote' && selectedDown === true) return

        const targetSignal = signal ? 1 : -1
        const votes = signal ? postVotes + 1 : postVotes - 1
        let sendData: voteData = {
            votes: votes,
            target: targetSignal
        }
        fetch(`http://localhost:8080/api/posts/comment/vote/${data.commentId}`, {
            credentials: "include",
            method: "POST",
            body: JSON.stringify(sendData)
        })
            .then(async (response) => {
                if (!response.ok) {
                    console.error(response)
                    return
                }
                const recData = await response.json()
                changeVotes(recData.Votes)
                switch (recData.Target) {
                    case -1:
                        setSelectedUp(false)
                        setSelectedDown(true)
                        break
                    case 1:
                        setSelectedUp(true)
                        setSelectedDown(false)
                        break
                    default:
                        setSelectedUp(false)
                        setSelectedDown(false)
                }
            })
    }
    return (
        <div className="post-comment">
            <div className="post-comment__header">
                <div className="post-comment__creator--wrapper">
                    <div className="post-comment__creator-avatar" >
                        <img src={avatarImage} loading="lazy" />
                    </div>
                    <div className="post-comment__creator-name"
                        id={data.commentCreator.creatorId}
                        onClick={() => { navigate(`/users/${data.commentCreator.creatorId}`) }}>
                        {data.commentCreator.creatorName}
                    </div>
                </div>
                <div className="post-comment__time">{data.commentTime}</div>
            </div>
            <div className="post-comment__main">
                <div className="post-comment-main__content-text">{data.commentContent}</div>
                <div className="post-comment-main__content-image">
                    <img src={commentImage} loading="lazy" />
                </div>
            </div>
            <div className="post-comment__options">
                <img className={`post-comment__vote up-vote ${selectedUp ? 'voted' : ''}`}
                    id='postUpVote'
                    onClick={(e) => { voteHandler(true, e) }}
                    src='/assets/UpVote.svg' />
                <div className="post-comment__vote-number">{postVotes}</div>
                <img className={`post-comment__vote down-vote ${selectedDown ? 'voted' : ''}`}
                    id='postDownVote'
                    onClick={(e) => { voteHandler(false, e) }}
                    src='/assets/DownVote.svg' />
            </div>
        </div>
    );
};
export default PostComment