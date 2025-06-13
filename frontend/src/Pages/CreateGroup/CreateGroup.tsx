import { HandleImageSend } from "../../utils/ImageControl"
import { useNavigate } from "react-router-dom"
import { useState, useRef } from "react"
import './CreateGroup.css'

function CreateGroup() {
    const navigate = useNavigate()

    const [error, setError] = useState('')
    const imageRef = useRef<HTMLInputElement | null>(null)
    const [imageUrl, setUrl] = useState('')
    const [groupData, setGroupData] = useState<{
        groupName: string;
        GroupAbout: string;
        image?: File;
    }>({
        groupName: '',
        GroupAbout: ''
    })

    const handleChange = (e: any) => {
        const { name, value, type } = e.target
        if (type === 'file') {
            const file = HandleImageSend(e.target.files?.[0])
            if (typeof file === 'string') {
                setError(file)
                return
            }
            setUrl(URL.createObjectURL(file))
            setGroupData((prevData) => ({
                ...prevData,
                image: file,
            }))
        } else {
            setGroupData((prevData) => ({
                ...prevData,
                [name]: value,
            }))
        }
    }

    const handleSubmit = async (e: any) => {
        if (groupData.groupName === '') {
            setError('Please enter group name')
            return
        } else if (groupData.GroupAbout === '') {
            setError('Please enter group description')
            return
        }
        if (groupData.groupName.length > 25) setError('Too long name. Max 25 chars')
        const { image, ...formDataWithoutAvatar } = groupData
        await fetch('http://localhost:8080/api/groups/create', {
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
            const id = await response.json()
            if (image) {
                const data = new FormData() //FormData is for avatar to send this as File 
                data.append('image', image)
                await fetch(`http://localhost:8080/api/image/group/${id}`, {
                    method: "POST",
                    credentials: "include",
                    body: data
                }).then(async (response) => {
                    console.log(response)
                    if (!response.ok) {
                        const message = await response.text()
                        setError(message)
                        return
                    }
                })
            }
            navigate(`/groups/${id}`)
        })
    }

    return (
        <div className="create-group-main">
            <input
                className="event-title event-inputs"
                placeholder="TITLE"
                name="groupName"
                value={groupData.groupName}
                onChange={handleChange} />
            <textarea className="content-text-area"
                placeholder="DESCRIPTION"
                name="GroupAbout"
                value={groupData.GroupAbout}
                onChange={handleChange} />
            <div className="add-image"
                onDrop={handleChange}
                onClick={() => { imageRef.current?.click() }}>
                {groupData.image ?
                    <img className="post-img" alt="" src={imageUrl} /> : <p>Drop image or click here</p>
                }
                <input hidden
                    type="file"
                    className="post-image-input"
                    onDrop={handleChange}
                    onChange={handleChange}
                    ref={imageRef} />
            </div>
            <p className='post-create__error group-error' >{error}</p>
            <button className="create-button group-create-btn" onClick={handleSubmit}>CREATE</button>
        </div>
    )
}

export default CreateGroup