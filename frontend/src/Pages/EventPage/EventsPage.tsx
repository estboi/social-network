import { useEffect, useState } from "react"
import { useNavigate, useParams } from 'react-router-dom'

import List from "../../components/list/list"
import Calender from "./components/EventsCalender"
import "./EventsPage.css"
import fetchData from "../../utils/fetchData"


const EventsPage = () => {
    const [data, setData] = useState([])
    const [filteredData, setFilteredData] = useState([])
    const [isChanged, setChanged] = useState(false)
    const { groupId } = useParams()

    useEffect(() => {
        const link = groupId ? `events/group/${groupId}` : 'events/all'
        fetchData(link, setData)
        if (groupId) { getGroupEvents(groupId) }
    }, [])

    useEffect(() => {
        setFilteredData(data)
    }, [data])


    const getGroupEvents = (id: string) => {
        const filteredEvents = data.filter((eventData: any) => {
            return eventData.groupId == id
        })
        setChanged(true)
        setFilteredData(filteredEvents)
    }

    const getEventsAtDay = (date: Date) => {
        const formattedDate = date.toISOString().slice(0, 10)
        const filteredEvents = data.filter((eventData: any) => {
            const eventDate = eventData.date.split(' ')[0]
            return eventDate === formattedDate
        })
        return filteredEvents
    }
    const handleMonthClick = (month: number) => {
        if (!data) return
        const filteredEvents = data.filter((eventData: any) => new Date(eventData.date).getMonth() === month)
        setChanged(true)
        setFilteredData(filteredEvents)
    }
    const handleDayClick = (date: Date) => {
        if (!data) return
        const filteredEvents = getEventsAtDay(date)
        setChanged(true)
        setFilteredData(filteredEvents)
    }

    return (
        <div className="events-page">
            {isChanged && <button className="events-page__back" onClick={() => {
                setFilteredData(data)
                setChanged(false)
            }}>RESET</button>}
            <Calender data={filteredData} handleMonthClick={handleMonthClick} handleDayClick={handleDayClick} />
            <List className="events-page__list" signal="events" data={filteredData} />
        </div>
    )
}
export default EventsPage