import { useEffect, useState } from "react"
import "./eventChip.css"
import { useNavigate } from "react-router-dom"

const EventChip = ({ data }: any) => {
    const [attending, setAttendance] = useState(true)
    const [pending, setPending] = useState(true)
    const navigate = useNavigate()
    useEffect(() => {
        setAttendance(data.attending)
        setPending(data.pending)
        console.log(data)
    }, [data.attending, data.pending])
    const chipClass = pending
        ? 'events-chip'
        : attending ? 'events-chip attending' : 'events-chip notattending'

    const handleAccept = () => {
        fetch(`http://localhost:8080/api/events/attend/${data.ID}`,
            { credentials: "include", method: "POST" })
            .then((resp) => {
                if (resp.status === 202) {
                    setPending(false)
                    setAttendance(true)
                }
            })
    }

    const handleDeny = () => {
        fetch(`http://localhost:8080/api/events/deny/${data.ID}`,
            { credentials: "include", method: "POST" })
            .then((resp) => {
                if (resp.status === 202) {
                    setPending(false)
                    setAttendance(false)
                }
            })
    }

    return (
        <div className={chipClass}>
            <div className="events-chip__main">
                <div className="events-chip__title">{data.name}</div>
                <p className="events-chip__description">{data.about}</p>
                <div className="events-chip__date">{data.date}</div>
                <p className="events-chip__group" onClick={() => { navigate(`/groups/${data.groupId}`) }}>
                    GROUP: {data.groupName}
                </p>
            </div>
            <div className="events-chip__options">
                {pending ?
                    <>
                        {!attending &&
                            <>
                                <button className="events-chip__button"
                                    onClick={handleAccept}>GOING</button>
                                <button className="events-chip__button not"
                                    onClick={handleDeny}>NOT GOING</button>
                            </>
                        }
                    </>
                    :
                    <>
                        {attending ?
                            <p className="events-chip__attending">ATTENDING</p>
                            :
                            <p className="events-chip__attending not">NOT ATTENDING</p>
                        }
                    </>
                }
            </div>
        </div>)
}

export default EventChip
