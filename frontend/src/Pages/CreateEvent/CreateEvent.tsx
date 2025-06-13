import { useNavigate, useParams } from "react-router-dom"
import './CreateEvent.css'
import { useState } from "react"
import { time } from "console"

function CreateEvent() {
    const navigate = useNavigate()
    const [error, setError] = useState('')

    const { groupId } = useParams()
    const groupIdAsString = groupId || '';
    const Id = parseInt(groupIdAsString, 10)

    const [eventData, setEventData] = useState({
        name: '',
        about: '',
        time: '',
        groupId: Id
    })

    const handleChange = (e: any) => {
        const { name, value } = e.target
        setEventData((prevData) => ({
            ...prevData,
            [name]: value,
        }))
    }

    const handleSubmit = async (e: any) => {
        if (eventData.name === '') {
            setError('Please enter event name')
            return;
        } else if (eventData.name.length > 25) {
            setError('Name max length is 25 chars')
            return;
        } else if (eventData.about === '') {
            setError('Please enter event description')
            return
        } else if (eventData.time === '') {
            setError('Please enter event date')
            return;
        } else if (eventData.time < formattedTomorrow) {
            setError('You cannot create event in history')
            return;
        } else if (eventData.time.split('-')[0] > String(tomorrow.getFullYear())) {
            setError('You cannot create in future year')
            return;
        } else {
            setError('')
        }

        await fetch(`http://localhost:8080/api/events/create/${Id}`, {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            credentials: "include",
            body: JSON.stringify(eventData)
        }).then(async (response) => {
            if (!response.ok) {
                const message = await response.text()
                setError(message)
                return
            } else {
                navigate(`/events/${groupId}`)
            }
        })
    }

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    // Format tomorrow's date in MM-dd-yyyy format
    const formattedTomorrow =
        tomorrow.getFullYear() + '-' +
        (tomorrow.getMonth() + 1).toString().padStart(2, '0') + '-' +
        tomorrow.getDate().toString().padStart(2, '0')



    return (
        <>
            <div className="event-create-main">
                <input
                    className="event-title event-inputs"
                    placeholder="Title"
                    name="name"
                    value={eventData.name}
                    onChange={handleChange} />
                <textarea
                    className="event-description event-inputs"
                    placeholder="description"
                    name="about"
                    value={eventData.about}
                    maxLength={250}
                    onChange={handleChange} />
                <input
                    className="event-date event-inputs"
                    type="date"
                    name="time"
                    value={eventData.time}
                    onChange={handleChange} />
                <p className='post-create__error' >{error}</p>
                <button className="header__back events-submit" onClick={handleSubmit}>CREATE</button>
            </div>
        </>
    )
}

export default CreateEvent