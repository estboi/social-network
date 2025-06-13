import React, { useRef, useState } from "react";
import './postAdd.css'
import { HandleImageSend } from "../../../../../../utils/ImageControl";
import EmojiPicker from "../../../../../../utils/emojiFuncs";
import { useNavigate } from "react-router-dom";

interface newComment {
    content: string
    image?: File
}




const AddComment = ({ id, onSubmit, className }: any) => {
    const [inputValue, setInputValue] = useState('')
    const [imageFile, setFile] = useState<File | undefined>(undefined)
    const [newImage, setImage] = useState<string | null>(null)
    const newCommentRef = useRef<HTMLTextAreaElement>(null)
    const [openEmoji, setEmoji] = useState(false)

    const autoExpand = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        const textarea = event.target
        const newHeight = `${textarea.scrollHeight}px`

        textarea.value === '' ? textarea.style.height = '2rem' : textarea.style.height = newHeight
        setInputValue(textarea.value)
    }
    const charactersLeft = 0 + inputValue.length;
    const handleFocus = (event: React.FocusEvent<HTMLTextAreaElement>) => {
        const parent = event.currentTarget.parentElement?.parentElement
        if (parent) {
            changeCommentSize(parent)
        }
    }
    const handleClick = (event: React.MouseEvent) => {
        const parent = event.currentTarget.parentElement?.parentElement?.parentElement
        if (parent) {
            changeCommentSize(parent)
            setEmoji(false)
        }
    }
    const closeBtnEmoji = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault()
        setEmoji(false)
    }
    const changeCommentSize = (parent: HTMLElement) => {
        parent.style.transition = 'height 0.3s ease-out'
        const options = parent.querySelector('.post__comment-options')
        const optionsHTML = options as HTMLElement
        optionsHTML.style.bottom = '-1.75rem'
        optionsHTML.style.left = '0.5rem'

        const submitBtn = parent.querySelector('.post__comment-submit')
        const child = submitBtn as HTMLElement
        child.style.bottom = '-1.75rem'
        parent.style.minHeight = '5rem'

        parent.style.paddingBottom = '2rem'
    }
    const removeImg = () => {
        const input = document.querySelector('#addImage')
        if (input) {
            const elem = input as HTMLInputElement
            elem.value = ''
        }
        setImage(null)
    }
    const handleBlur = (event: React.FocusEvent<HTMLTextAreaElement>) => {
        const parent = event.currentTarget.parentElement?.parentElement
        if (parent) {
            const childs = parent.querySelector('.post__comment-input__chars ')
            const child = childs as HTMLElement
            child.style.display = 'block'

            parent.style.height = parent.style.height
        }
    }
    // handling the draggin over
    const handleDragging = (event: React.DragEvent<HTMLDivElement>) => {
        event.preventDefault()
    }
    const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
        event.preventDefault()
        const droppedFiles = event.dataTransfer.files;

        handleDroppedFiles(droppedFiles);
    }
    //Handling the adding of files
    const handleDroppedFiles = (files: FileList) => {
        if (files.length > 0) {
            const uploadedFile = files[0];
            const reader = new FileReader();
            const response = HandleImageSend(uploadedFile)
            if (typeof response === 'string') {
                return
            }
            reader.onload = () => {
                setFile(uploadedFile)
                // Set the uploaded image URL
                setImage(reader.result as string);
            };

            // Read the contents of the file
            reader.readAsDataURL(uploadedFile);
        }
    }
    const addImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = event.target.files && event.target.files[0];
        if (selectedFile) {
            const reader = new FileReader();
            const response = HandleImageSend(selectedFile)
            if (typeof response === 'string') {
                return
            }
            reader.onload = () => {
                setFile(selectedFile)
                // Set the uploaded image URL
                setImage(reader.result as string);
            }
            // Read the contents of the file
            reader.readAsDataURL(selectedFile);
        }
    }

    // Reload input values method
    const resetValues = () => {
        removeImg()
        const input = document.querySelector('#commentInput')
        if (input) {
            const elem = input as HTMLInputElement
            elem.value = ''
        }
        setInputValue('')
        setEmoji(false)
    }

    // SUBMIT event
    const submitNewComment = async (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault()
        const sendData: newComment = {
            content: inputValue,
            image: imageFile
        }
        // Validate values
        if (sendData.content.length === 0 && !sendData.image || sendData.content.length > 250) return
        await fetch(`http://localhost:8080/api/posts/comments/create/${id}`,
            {
                method: "POST",
                credentials: "include",
                body: JSON.stringify(sendData.content)
            })
            .then(async (response) => {
                if (response.ok) {
                    const commentId = await response.json()
                    if (sendData.image) {
                        const data = new FormData()
                        data.append('image', sendData.image)
                        await fetch(`http://localhost:8080/api/image/comment/${commentId}`,
                            {
                                method: "POST",
                                credentials: "include",
                                body: data
                            }).then((response) => {
                                if (!response.ok) {
                                    return
                                }
                            })
                    }
                    resetValues()
                    onSubmit(id)
                }
            })
    }

    return (
        <div className={`post__add-comment ${className}`}
            onDragEnter={handleDragging}
            onDragLeave={handleDragging}
            onDragOver={handleDragging}
            onDrop={handleDrop}>
            <form className='post__comment__form' id='commentForm'>
                <textarea className='post__comment-input'
                    id='commentInput'
                    placeholder='ADD COMMENT...'
                    maxLength={250}
                    form='commentForm'
                    ref={newCommentRef}
                    onInput={autoExpand}
                    onBlur={handleBlur}
                    onFocus={handleFocus} />
                <div className='post__comment-input__chars' data-max={charactersLeft === 250 && 'max'}>
                    {charactersLeft >= 0 ? `${charactersLeft}/250` : '250/250'}
                </div>
                {newImage &&
                    <div className='post__comment__img-show--wrapper'>
                        <img className='post__comment__img-show' src={newImage} />
                        <button onClick={removeImg} className='post__comment__img-show--delete'>X</button>
                    </div>}
                <div className='post__comment-options'>
                    <label className='post__comment-btn--wrapper' htmlFor={`addImage${id}`} onClick={handleClick}>
                        <img className='post__new-comment-img' src="/assets/AddImage.svg" />
                        <input type="file" hidden onChange={addImageChange}
                            id={`addImage${id}`} accept="image/png, image/jpeg, image/gif" />
                    </label>
                    <div className='post__comment-btn--wrapper' onClick={(e) => {
                        handleClick(e)
                        setEmoji(true)
                    }}>
                        <img className='post__new-comment-img' src='/assets/AddSmile.svg' />
                    </div>
                    {openEmoji && <EmojiPicker input={newCommentRef} closeBtn={closeBtnEmoji} className={'full-post__emoji'} />}
                </div>
                <button className='post__comment-submit' type='submit' form='commentForm' onClick={submitNewComment}>
                    <div className='post__comment-btn--wrapper'>
                        <img className='post__new-comment-img' src='/assets/Submit.svg' />
                    </div>
                </button>
            </form>
        </div >
    )
}
export default AddComment