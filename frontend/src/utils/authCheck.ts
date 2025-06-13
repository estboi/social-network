const isAuth = async () => {
    try {
        const response = await fetch('http://localhost:8080/api/auth', { credentials: "include", method: "GET" });
        if (!response.ok) {
            // Handle non-successful responses here
            if (response.status === 401) {
                throw new Error(`User is not authorised`);
            } else if (response.status === 202) {
                return true;
            }
        }

        // Handle successful response
        return true;
    } catch (error) {
        // Handle network errors or other exceptions here
        console.error("Error during fetch:", error);
        return false;
    }
};

export default isAuth