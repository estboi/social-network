import { useNavigate, useParams } from "react-router-dom"
import React, { useEffect, useRef, useState } from "react"

import { HandleImageSend } from "../../utils/ImageControl"
import './CreatePost.css'
import fetchData from "../../utils/fetchData";
interface PostCredentials {
    content: string;
    privacy?: string;
    image?: File;
    friends?: number[]
}

interface Users {
    id: number,
    isFollow: boolean
    isFollower: boolean
    isPending: boolean
    name: string
}

function CreatePost() {
    const navigate = useNavigate()

    const [selectedPrivacy, setSelectedPrivacy] = useState('Public');
    const [usersSelect, setSelectedUsers] = useState(false)
    const [selectedUsers, setSelectedUserIds] = useState<number[]>([])
    const [users, setUsers] = useState<[Users] | null>(null)
    const imageRef = useRef<HTMLInputElement | null>(null)
    const [error, setError] = useState('')
    const [imageUrl, setUrl] = useState('')
    const [postData, setPostData] = useState<PostCredentials>({
        content: '',
        privacy: selectedPrivacy,
        friends: selectedUsers
    })
    let Id = 0
    const { groupId } = useParams()
    if (groupId) {
        const groupIdAsString = groupId || '';
        Id = parseInt(groupIdAsString, 10)
    }

    const handleContentChange = (e: any) => {
        const value = e.target.value;
        setPostData((prevData) => ({
            ...prevData,
            content: value,
        }));
    };

    const handlePrivacyChange = (e: any) => {
        let name = e.target.name;
        setSelectedUsers(name === "Followers" ? true : false)
        setSelectedPrivacy(name)
        setPostData((prevData) => ({
            ...prevData,
            privacy: name,
        }))

    }

    const handleImageDrop = (e: any) => {
        e.preventDefault()
        if (e.target.files.length < 0) return
        const file = HandleImageSend(e.target.files?.[0])
        if (typeof file === 'string') {
            setError(file)
            return
        }
        setUrl(URL.createObjectURL(file))
        setPostData((prevData) => ({
            ...prevData,
            image: file,
        }))
    }

    const handleSubmit = async () => {
        if (postData.content === '') {
            setError('Please enter content')
            return
        }
        console.log(selectedUsers);

        const { image, ...formDataWithoutAvatar } = postData
        console.log(postData);
        let link = `http://localhost:8080/api/posts/create/${Id}`
        await fetch(link, {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            credentials: "include",
            body: JSON.stringify(formDataWithoutAvatar)
        }).then(async (response) => {
            if (!response.ok) {
                const message = await response.text()
                setError(message)
                return
            }
            setError('')
            if (image) {
                const postId = await response.json()
                const data = new FormData() //FormData is for Image to send this as File 
                data.append('image', image)
                await fetch(`http://localhost:8080/api/image/post/${postId}`, {
                    method: "POST",
                    body: data,
                    credentials: "include"
                })
            }
            navigate(-1)
        })
    }

    useEffect(() => {
        if (!usersSelect) return () => { }
        try {
            fetchData('users/followers', setUsers)
        } catch {
            console.error("Something went wrong")
        }
    }, [usersSelect])

    const handleSelectChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedValue: number = parseInt(event.target.value, 10);
        // Check if the user ID is already in the array
        if (!Number.isNaN(selectedValue) && !selectedUsers.includes(selectedValue)) {
            setSelectedUserIds([...selectedUsers, selectedValue]);
            setPostData((prevData) => ({
                ...prevData,
                friends: [...selectedUsers, selectedValue],
            }))
        } else if (selectedUsers.includes(selectedValue)) {
            const updatedUsers = selectedUsers.filter((userId) => userId !== selectedValue);
            setSelectedUserIds(updatedUsers);
            setPostData((prevData) => ({
                ...prevData,
                friends: updatedUsers,
            }));
        }
    };

    return (
        <div className="posts-create--page">
            <div className="post-create-fields">
                <textarea
                    className="content-text-area"
                    placeholder="START WRITING THE CONTENT"
                    name="content"
                    maxLength={500}
                    value={postData.content}
                    onChange={handleContentChange}
                />
                <div className="add-image"
                    onDrop={handleImageDrop}
                    onClick={() => { imageRef.current?.click() }}>
                    {postData.image ?
                        <img className="post-img" alt="" src={imageUrl} /> : <p>Drop image or click here</p>
                    }
                    <input
                        type="file"
                        className="post-image-input"
                        onDrop={handleImageDrop}
                        onChange={handleImageDrop}
                        ref={imageRef}
                        hidden
                    />

                </div>
            </div>
            {!groupId && <div className="privacy">
                <p className="change-privacy">CHANGE PRIVACY:</p>
                <button
                    className="privacy__button"
                    name="Public"
                    onClick={handlePrivacyChange}
                    style={{
                        borderColor: selectedPrivacy === 'Public' ? 'green' : 'transparent',
                    }}>
                    PUBLIC</button>
                <button
                    className="privacy__button"
                    name="Private"
                    onClick={handlePrivacyChange}
                    style={{
                        borderColor: selectedPrivacy === 'Private' ? 'green' : 'transparent',
                    }}>
                    PRIVATE</button>
                <button
                    className="privacy__button"
                    name="Followers"
                    onClick={handlePrivacyChange}
                    style={{
                        borderColor: selectedPrivacy === 'Followers' ? 'green' : 'transparent',
                    }}>FOR FOLLOWERS</button>
                {usersSelect && <>
                    <div className="select-users--wrapper">
                        <select onChange={handleSelectChange}>
                            <option value='OPTION_PLACEHOLDER'>Select followers</option>
                            {users !== null && users.map((user) => (
                                <option key={user.id} value={user.id}>{user.name}</option>
                            ))}
                        </select>
                    </div>
                    <div className="select-users__nubmer">NUMBER OF SELECTED FRIENDS: {selectedUsers.length}</div>
                </>
                }
            </div>}
            <p className='post-create__error' >{error}</p>
            <button className="create-button" onClick={() => { handleSubmit() }}>CREATE</button>
        </div>
    )
}

export default CreatePost