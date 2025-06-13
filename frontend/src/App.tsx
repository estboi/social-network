import { useEffect, useRef, useState } from 'react'
import { Routes, Route, useNavigate } from 'react-router-dom'

import './App.css'

import LoginPage from './Pages/LoginPage/LoginPage'
import RegisterPage from './Pages/RegisterPage/RegisterPage'
import HomePage from './Pages/HomePage/HomePage'
import ProfilePage from './Pages/ProfilePage/ProfilePage'
import CreatePost from './Pages/createPost/CreatePost'
import ListPage from './Pages/listPage/listPage'
import CreateEvent from './Pages/CreateEvent/CreateEvent'
import EventsPage from './Pages/EventPage/EventsPage'
import Chat from './Pages/chatPage/chatPage'
import CreateGroup from './Pages/CreateGroup/CreateGroup'
import GroupPage from './Pages/GroupPage/GroupPage'
import Navbar from './components/navbar/navbar'
import Header from './components/header/header'
import FullPost from './Pages/FullPost/fullPost'
import isAuth from './utils/authCheck'
import GroupOptions from './Pages/GroupOptions/groupOptionsPage'
import { WebsocketProvider } from './websocket/Provider'

function App() {
    const navigate = useNavigate()
    const [logged, setIsLogged] = useState(false)

    useEffect(() => {
        isAuth().then((isLogged) => {
            if (isLogged === false) {
                navigate('/login')
                return
            }
        })
    }, [])

    const PageContainer: React.FC<{ children: React.ReactNode }> = ({ children }) => {
        return (
            <div className='page'>
                <Header className="" />
                <div className='main'>
                    <Navbar />
                    {children}
                </div>
            </div>
        );
    };

    return (
        <WebsocketProvider isLogged={logged}>
            <>
                <Routes>
                    <Route path="/register" element={<RegisterPage setLogin={setIsLogged} />} />
                    <Route path="/login" element={<LoginPage setLogin={setIsLogged} />} />
                    <Route path='/*' element={
                        <PageContainer>
                            <Routes>
                                <Route index element={<HomePage />} />
                                <Route path='/post'>
                                    <Route path=':postID' element={<FullPost />} />
                                </Route>
                                <Route path="/createPost/:groupId?" element={<CreatePost />} />

                                {/* USERS */}
                                <Route path='/users'>
                                    <Route index element={<ListPage signal={['users', 'all']} />} />
                                    <Route path=':userId' element={<ProfilePage />} />
                                    <Route path='followed' element={<ListPage signal={['users', 'followed']} />} />
                                    <Route path='followers' element={<ListPage signal={['users', 'followers']} />} />
                                </Route>
                                d
                                {/* GROUPS */}
                                <Route path='/groups'>
                                    <Route index element={<ListPage signal={['groups', 'all']} />} />
                                    <Route path=':groupId' element={<GroupPage />} />
                                    <Route path=':groupId/members' element={<GroupOptions />} />
                                    <Route path=':groupId/createEvent' element={<CreateEvent />} />
                                    <Route path=':groupId/createPost' element={<CreatePost />} />
                                    <Route path='connected' element={<ListPage signal={['groups', 'connected']} />} />
                                    <Route path='created' element={<ListPage signal={['groups', 'created']} />} />
                                    <Route path='createGroup' element={<CreateGroup />} />
                                </Route>

                                {/* EVENTS */}
                                <Route path='/events' element={<EventsPage />} />
                                <Route path='/events/:groupId' element={<EventsPage />} />

                                {/* CHATS */}
                                <Route path='/chats'>
                                    <Route index element={<ListPage signal={['chats', 'all']} />} />
                                    <Route path='users/:userId' element={<Chat type='user' />} />
                                    <Route path='groups/:groupId' element={<Chat type='group' />} />
                                </Route>
                            </Routes>
                        </PageContainer>
                    } />
                </Routes>
            </>
        </WebsocketProvider>
    )
}

export default App
