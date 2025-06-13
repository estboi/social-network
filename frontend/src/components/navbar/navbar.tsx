import { useNavigate } from 'react-router-dom'
import { useContext, useEffect, useRef, useState } from 'react'

import './navbar.css'
import './adaptive.css'

import fetchData from '../../utils/fetchData'
import { ImageGet } from '../../utils/ImageControl'
import List from '../list/list'
import { WebsocketContext, WebsocketContextProps } from '../../websocket/Provider'


const Navbar = () => {
    const navigate = useNavigate()

    const { ready, value, send }: WebsocketContextProps = useContext(WebsocketContext);
    const [newNotification, setNewNotification] = useState(false)
    const [notificationData, setNotificationData] = useState([])
    const [isNotification, openNotifications] = useState(false)
    const toggleNotifications = () => {
        setNewNotification(false)
        openNotifications((prevState) => !prevState)
    }

    const path = window.location.pathname
    const isPageActive = (currentPage: string | undefined) => { return path === currentPage }

    const [navBarData, setData] = useState<{ username: string, userID: number }>({ username: '', userID: 0 })
    const [image, setImage] = useState<string>('')

    const NavigateToPage = (event: React.MouseEvent<HTMLElement>) => {
        const element = event.target as HTMLElement
        const currentPage = element.dataset.setRoute
        if (isPageActive(currentPage)) return
        if (currentPage !== undefined) navigate(currentPage)
    }
    useEffect(() => {
        fetchData('navbar', setData)
            .then((isLogged) => {
                if (isLogged === false) {
                    navigate('/login')
                    return
                }
                ImageGet(`user/${navBarData.userID}`, setImage)
                fetchData('notification', setNotificationData)
            })
    }, [])

    useEffect(() => {
        if (value) {
            fetchData('notification', setNotificationData)
            setNewNotification(true)}
    }, [value])
    return (
        <>
            <div className='navbar'>
                <div className='navbar__folders-container'>
                    {renderFolders(NavigateToPage, isPageActive)}
                    <div className="navbar__folder folder" onClick={toggleNotifications}>
                        {newNotification ? <div className='new-notification'></div> : <></>}
                        <img className="navbar__folder__img" alt="folder" src="/assets/folder.png" />
                        <p className="navbar__folder__name">NOTIFICATIONS</p>
                    </div>
                </div>
                <div className="navbar__profile-bar" id={`/users/${navBarData.userID}`} onClick={() => { navigate(`/users/${navBarData.userID}`) }}>
                    <img src={image === '' ? '/assets/default_avatar.png' : image} alt="avatar" className="navbar__profile-avatar avatar" id={`/users/${navBarData.userID}`} />
                    <div className="navbar__profile-name nickname" id={`/users/${navBarData.userID}`}>{navBarData.username}</div>
                </div>
                {isNotification ?
                    (<List signal="notifications" className="notifications-list" data={notificationData} />) : <></>}
            </div >
        </>
    )
}

const renderFolders = (NavigateToPage: (event: any) => void, isPageActive: (currentPage: string) => boolean) => {
    const folders = [{
        name: 'POSTS',
        route: '/'
    },
    {
        name: 'USERS',
        route: '/users',
        subfolders: [
            {
                name: 'ALL',
                route: '/users',
            },
            {
                name: 'FOLLOWED',
                route: '/users/followed',
            },
            {
                name: 'FOLLOWERS',
                route: '/users/followers',
            },
        ],
    },
    {
        name: 'GROUPS',
        route: '/groups',
        subfolders: [
            {
                name: 'ALL',
                route: '/groups',
            },
            {
                name: 'CONNECTED',
                route: '/groups/connected',
            },
            {
                name: 'CREATED',
                route: '/groups/created',
            },
        ],
    },
    {
        name: 'CHATS',
        route: '/chats'
    },
    {
        name: 'EVENTS',
        route: '/events'
    },]
    return (
        <>
            {folders.map((folder, key) => (
                <div className="navbar__folder-pack" key={key}>
                    <div className={`navbar__folder folder ${isPageActive(folder.route) ? 'active-page' : ''}`}
                        data-set-route={folder.route} onClick={NavigateToPage}>
                        <img className="navbar__folder__img" alt="folder" src="/assets/folder.png" data-set-route={folder.route} />
                        <p className="navbar__folder__name" data-set-route={folder.route}>{folder.name}</p>
                    </div>
                    {folder.subfolders && folder.subfolders.map((subfolder, subkey) => (
                        <div className={`navbar__subfolder folder ${isPageActive(subfolder.route) ? 'active-page' : ''}`} data-set-route={subfolder.route}
                            onClick={NavigateToPage} key={subkey}>
                            <img className="navbar__subfolder__img" src="/assets/add-file-4.svg" alt="folder" data-set-route={subfolder.route} />
                            <p className="navbar__subfolder__name" data-set-route={subfolder.route}>{subfolder.name}</p>
                        </div>
                    ))}
                </div>
            ))}
        </>
    )
}
export default Navbar
