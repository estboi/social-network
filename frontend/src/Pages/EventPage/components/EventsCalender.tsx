import { useState } from "react";

interface Props {
    data: any,
    handleMonthClick: Function,
    handleDayClick: Function,
}


const Calender = ({ data, handleMonthClick, handleDayClick }: Props) => {
    const [hoveredDate, setHoveredDate] = useState<string | null>(null);

    /** HANDLE CLICK */
    const getEventsAtDay = (date: Date) => {
        if (!data) return []
        const formattedDate = date.toISOString().slice(0, 10)
        const filteredEvents = data.filter((eventData: any) => {
            const eventDate = eventData.date.split(' ')[0]
            return eventDate === formattedDate
        })
        return filteredEvents.length !== 0 ? filteredEvents : [] 
    }
    const changeOpacityToEventsNumber = (date: Date) => {
        const filteredEvents = getEventsAtDay(date)
        const eventsLen = filteredEvents ? filteredEvents.length : 0
        const numberOfEvents = eventsLen
        const incrementNumber = 0.2
        return `rgba(0,255,0,${0.1 + incrementNumber * numberOfEvents})`
    }

    /** Handle hover */
    const handleHover = (date: string) => {
        setHoveredDate(date);
    };

    const handleMouseLeave = () => {
        setHoveredDate(null);
    };

    const getDaysInYear = (year: number) => {
        const startDate = new Date(year, 0, 1); // January 1st
        const endDate = new Date(year, 11, 31); // December 31st

        const days = [];
        for (let date = new Date(startDate); date <= endDate; date.setDate(date.getDate() + 1)) {
            days.push(new Date(date));
        }

        return days;
    };

    const currentYear = new Date().getFullYear();
    const daysInCurrentYear = getDaysInYear(currentYear);
    const months = Array.from({ length: 12 }, (_, i) => {
        const tempDate = new Date(0, i);
        const monthName = tempDate.toLocaleString('en-US', { month: 'long' });
        return (
            <p key={i} className="events-page__calendar__month"
                onClick={(e) => { handleMonthClick(i) }}>{monthName}
            </p>
        );
    })

    return (
        <div className="events-page__calender-div">
            <div className="events-page__calender__months">
                {months}
            </div>
            <div className="events-page__calender">
                {daysInCurrentYear.map((date, index) => (
                    <div
                        key={index}
                        onMouseEnter={() => handleHover(date.toDateString())}
                        onMouseLeave={handleMouseLeave}
                        onClick={() => handleDayClick(date)}
                        style={{ backgroundColor: changeOpacityToEventsNumber(date) }}
                        className="events-page__calender__day"
                    >
                        {hoveredDate === date.toDateString() ?
                            <p className="events-page__calender__day__tooltip">{`${date.toDateString()} Events:${getEventsAtDay(date).length}`}</p>
                            : null}
                    </div>
                ))}
            </div>
        </div>
    )
}

export default Calender