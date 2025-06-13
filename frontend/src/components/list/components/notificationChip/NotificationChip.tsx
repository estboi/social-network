import { useNavigate } from 'react-router-dom'
import './NotificationChip.css'

interface notification {
    userId: number
    sourceId: number
    sourceType: string
    type: string
    Content: string
}

function NotificationChip({ data }: { data?: notification }) {
    const navigate = useNavigate()
    const handleClick = () => {
        var link
        if (data) {
            switch (data.type) {
                case "invite":
                    navigate(`groups`)
                    break
                case "request":
                    navigate(`groups/${data.sourceId}/members`)
                    break
                case "friendRequest":
                    navigate(`users`)
                    break
                case "newEvent":
                    navigate(`events/${data.sourceId}`)
                    break
            }
        }
    }

    return (
        <div className='notification' onClick={handleClick}>
            {data?.Content}
        </div>
    );
}

export default NotificationChip