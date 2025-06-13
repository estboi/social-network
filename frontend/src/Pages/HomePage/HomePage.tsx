import { useState, useEffect } from 'react'

import List from "../../components/list/list";
import fetchData from '../../utils/fetchData'

const HomePage = () => {
    const [data, setData] = useState([])

    useEffect(() => {
        fetchData('posts/home', setData)
    }, [])


    return <List signal={'posts'} className={'posts'} data={data} />
}
export default HomePage