import React, { useEffect, useState } from "react";

import List from "../../components/list/list";
import "./listPage.css";

import fetchData from "../../utils/fetchData";

interface listType {
    signal: string[],
}

const ListPage = React.memo(({ signal }: listType) => {
    const [data, setData] = useState([])
    useEffect(() => {
        fetchData(`${signal[0]}/${signal[1]}`, setData)
    }, [signal])
    return (
        <>
            <List signal={signal[0]} className={signal[0]} data={data} />
        </>
    )
})

export default ListPage