import { useNavigate } from 'react-router-dom'
import "./post.css";
import React, { useEffect, useState } from "react";

import AddComment from './components/postAdd/postAdd';
import { ImageGet } from '../../../../utils/ImageControl';


interface voteData {
    votes: number
    target: number
}

export const Post = ({ data, className, addNewComment }: any) => {
    const navigate = useNavigate()
    const [postImage, setPostImage] = useState('')
    const [avatarImage, setAvatarImage] = useState('')
    useEffect(() => {
        ImageGet(`user/${data.postCreator.creatorId}`, setAvatarImage)
        ImageGet(`post/${data.postId}`, setPostImage)
        if (data.postVotesTarget === 1) {
            setSelectedUp(true)
        } else if (data.postVotesTarget === -1) {
            setSelectedDown(true)
        }
        console.log(data)
    }, [data])

    const [selectedUp, setSelectedUp] = useState(false)
    const [selectedDown, setSelectedDown] = useState(false)
    const [postVotes, changeVotes] = useState(data.postVotes)

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
        fetch(`http://localhost:8080/api/posts/vote/${data.postId}`, {
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
    const openFullPost = (id: number) => {
        if (addNewComment) {
            addNewComment()
            return
        }
        navigate(`/post/${id}`)
    }

    return (
        <div className="post" id={`Post${data.postId}`} >
            <div className="post-header">
                <div className="post-header__avatar">
                    <img src={avatarImage ? avatarImage : "/assets/default_avatar.png"} alt='avatar'loading='lazy' />
                </div>
                <div className="post-header__name"
                    id={data.postCreator.creatorId}
                    onClick={() => { navigate(`/users/${data.postCreator.creatorId}`) }}>
                    {data.postCreator.creatorName}
                </div>
                <div className="post-header__time">{data.postTime}</div>
            </div>
            {postImage === '' &&
                <div className="post__content">
                    <p className='post__content-name'>{data.postCreator.creatorName}</p>
                    <div className='post__content-text'> {`${data.postContent}`}</div>
                </div>}
            <div className={`post-main ${className}`} onClick={() => { if (className !== 'full-post') { openFullPost(data.postId) } }}>
                <img className="post-main__image" src={postImage} loading='lazy' />
            </div>
            <div className="post-options">
                <div className="post-options__votes">
                    <img className={`post-options__img up-vote ${selectedUp ? 'voted' : ''}`}
                        id='postUpVote'
                        onClick={(e) => { voteHandler(true, e) }}
                        src='/assets/UpVote.svg' />

                    <div className="post-options__vote-number">{postVotes}</div>
                    <img className={`post-options__img down-vote ${selectedDown ? 'voted' : ''}`}
                        id='postDownVote'
                        onClick={(e) => { voteHandler(false, e) }}
                        src='/assets/DownVote.svg' />
                </div>
                {className !== 'full-post' && <div className="post-options__comments-btn" onClick={() => { openFullPost(data.postId) }}>
                    <img className="post-options__img" src='/assets/Comments.svg' />
                    <div className="post-options__comments-btn__title">COMMENTS</div>
                </div>}
            </div>
            {postImage !== '' && <div className="post__content">
                <p className='post__content-name'>{data.postCreator.creatorName}</p>
                <div className='post__content-text'> {`${data.postContent}`}</div>
            </div>}
            <AddComment id={data.postId} onSubmit={openFullPost} />
        </div >
    );
};

