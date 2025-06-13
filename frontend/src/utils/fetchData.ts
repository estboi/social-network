
const fetchData = async (endpoint: string, setData: any) => {
    return await fetch(`http://localhost:8080/api/${endpoint}`, { credentials: "include" })
        .then(async (response) => {
            if (response.ok) {
                const fetchData = await response.json()
                setData(fetchData)
                return true
            }
            if (response.status === 401) {
                throw ("user is not logged")
            }
            if (response.status === 406) {
                return response.status
            }
        })
        .catch(() => {
            return false
        })
}


export default fetchData