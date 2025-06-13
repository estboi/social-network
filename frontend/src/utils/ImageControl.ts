export function HandleImageSend(file: File) {
    //File type check
    const fileTypes = ['png', 'jpg', 'gif', 'jpeg']
    const fileName = file.name
    const type = fileName.split('.').pop()?.toLowerCase()
    if (!type || !fileTypes.includes(type)) {
        return 'Use the correct format of image: JPEG, PNG, GIF'
    }
    //File size check
    if (file.size / 1024 / 1024 > 10) {
        return 'Too big image. Not greater than 10 mb'
    }
    return file
}

export const ImageGet = async (endpoint: string, setData: Function) => {
    try {
        const response = await fetch(`http://localhost:8080/api/image/${endpoint}`, { credentials: "include" })
        const blob = await response.blob()
        if (blob.size !== 0 && blob.type !== '') {
            setData(URL.createObjectURL(blob))
        }
    } catch (error) {
        console.error(error)
    }
}

export function parseNotification(notification: string): { type: string; payload: string } | null {
    const regexPattern = /\{"type":"([^"]*)","payload":"([^"]*)"\}/;
  
    const match = notification.match(regexPattern);
  
    if (match) {
      const [, type, payload] = match;
      return { type, payload };
    }
  
    return null;
  }